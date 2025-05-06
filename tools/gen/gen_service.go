package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	repoTpl = `package repository

import "gorm.io/gorm"

type {{.Service}}Repository interface {
	// Define repository methods here
}

type {{.Repository}}Repository struct {
	db *gorm.DB
}

func New{{.Service}}Repository(db *gorm.DB) {{.Service}}Repository {
	return &{{.Repository}}Repository{db}
}
`

	serviceTpl = `package service

import "{{.Module}}/internal/repository"

type {{.Receiver}}Service struct {
	repo repository.{{.Service}}Repository
}

func New{{.Service}}Service(r repository.{{.Service}}Repository) *{{.Receiver}}Service {
	return &{{.Receiver}}Service{repo: r}
}
`

	handlerTpl = `package handler

import (
	"{{.Module}}/internal/service"
	"github.com/gofiber/fiber/v2"
)

type {{.Receiver}}Handler struct {
	svc *service.{{.Service}}Service
}

func New{{.Service}}Handler(svc *service.{{.Service}}Service) *{{.Receiver}}Handler {
	return &{{.Receiver}}Handler{svc: svc}
}

func (h *{{.Receiver}}Handler) Register(router fiber.Router) {
	router.Get("/", h.list)
	router.Post("/", h.create)
	router.Get("/:id", h.get)
	router.Put("/:id", h.update)
	router.Delete("/:id", h.delete)
}

func (h *{{.Receiver}}Handler) list(c *fiber.Ctx) error   {
	return c.SendStatus(fiber.StatusOK) 
}
func (h *{{.Receiver}}Handler) create(c *fiber.Ctx) error { 
	return c.SendStatus(fiber.StatusCreated) 
}
func (h *{{.Receiver}}Handler) get(c *fiber.Ctx) error    { 
	return c.SendStatus(fiber.StatusOK) 
}
func (h *{{.Receiver}}Handler) update(c *fiber.Ctx) error { 
	return c.SendStatus(fiber.StatusOK) 
}
func (h *{{.Receiver}}Handler) delete(c *fiber.Ctx) error { 
	return c.SendStatus(fiber.StatusNoContent) 
}
`

	dtoTpl = `package dto

//- todo))  Define DTOs for {{.Service}} operations
`

	testDTOTpl = `package dto_test

import "testing"

func Test{{.Service}}DTO(t *testing.T) {
	t.Log(" todo))  write dto test")
}
`
)

func main() {
	name := flag.String("name", "", "Name of the service (e.g. 'rsvp')")
	overwrite := flag.Bool("force", false, "Force overwrite if file exists")
	dryRun := flag.Bool("dry-run", false, "Preview generated files without writing them")
	customModule := flag.String("module", "", "Override Go module name")

	flag.Parse()

	if *name == "" {
		fmt.Println("❌ Please provide a service name with --name")
		os.Exit(1)
	}
	if !isValidGoIdent(*name) {
		fmt.Println("❌ Invalid service name. Must be a valid Go identifier.")
		os.Exit(1)
	}

	serviceName := strings.ToUpper(string((*name)[0])) + (*name)[1:]
	receiver := cases.Title(language.Und).String(strings.ToUpper(serviceName[:1])) + serviceName[1:]
	repoName := strings.ToLower(serviceName[:1]) + serviceName[1:]
	module := *customModule
	if module == "" {
		module = getGoModule()
	}

	data := map[string]string{
		"Service":    serviceName,
		"Receiver":   receiver,
		"Repository": repoName,
		"Module":     module,
	}

	fileName := strings.ToLower(*name)

	write("internal/repository", fmt.Sprintf("%s_repository.go", fileName), repoTpl, data, *overwrite, *dryRun)
	write("internal/service", fmt.Sprintf("%s_service.go", fileName), serviceTpl, data, *overwrite, *dryRun)
	write("internal/handler", fmt.Sprintf("%s_handler.go", fileName), handlerTpl, data, *overwrite, *dryRun)
	write("internal/dto", fmt.Sprintf("%s.go", fileName), dtoTpl, data, *overwrite, *dryRun)
	write("test/unit/dto", fmt.Sprintf("%s_input_test.go", fileName), testDTOTpl, data, *overwrite, *dryRun)

	if !*dryRun {
		updateContainer(serviceName)
		updateRoutes(serviceName)

		generatedFiles := []string{
			fmt.Sprintf("internal/repository/%s_repository.go", fileName),
			fmt.Sprintf("internal/service/%s_service.go", fileName),
			fmt.Sprintf("internal/handler/%s_handler.go", fileName),
			fmt.Sprintf("internal/dto/%s.go", fileName),
			fmt.Sprintf("test/unit/dto/%s_input_test.go", fileName),
		}
		writeHistory(*name, generatedFiles)
		fmt.Println("📝 History updated in .gen_history.json")
		fmt.Printf("✅ Service '%s' generated and registered in container.\n", fileName)
	} else {
		fmt.Println("💡 Dry run mode: no files written, preview only.")
	}
}

func write(dir, filename, tpl string, data map[string]string, overwrite, dry bool) {
	fullPath := filepath.Join(dir, filename)
	content := executeTemplate(tpl, data)

	if dry {
		fmt.Println("🔍 Preview:", fullPath)
		fmt.Println(content)
		fmt.Println(strings.Repeat("-", 60))
		return
	}

	if _, err := os.Stat(fullPath); err == nil && !overwrite {
		fmt.Printf("⚠️ File %s already exists. Use --force to overwrite.\n", fullPath)
		return
	}

	must(os.MkdirAll(dir, 0755))
	must(os.WriteFile(fullPath, []byte(content), 0644))
	fmt.Println("📄 Created:", fullPath)
}

func executeTemplate(tpl string, data map[string]string) string {
	t := template.Must(template.New("").Parse(tpl))
	var b strings.Builder
	must(t.Execute(&b, data))
	return b.String()
}

func updateContainer(serviceName string) {
	path := "internal/service/container.go"
	lines := []string{}

	file, err := os.Open(path)
	if err != nil {
		fmt.Println("❌ container.go not found")
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	insertedField, insertedAssign, insertedImport := false, false, false
	importBlockStarted := false
	importLines := []string{}
	alreadyImportedRepository := false

	module := getGoModule()

	for scanner.Scan() {
		line := scanner.Text()

		// Check for existing import
		if strings.HasPrefix(strings.TrimSpace(line), `"`+module+`/internal/repository"`) {
			alreadyImportedRepository = true
		}

		// Handle import block start
		if strings.HasPrefix(strings.TrimSpace(line), "import (") {
			importBlockStarted = true
		}

		// Handle import block end
		if importBlockStarted && strings.HasPrefix(strings.TrimSpace(line), ")") && !alreadyImportedRepository {
			importLines = append(importLines, `  "`+module+`/internal/repository"`)
			insertedImport = true
		}

		// Collect import lines separately
		if importBlockStarted {
			importLines = append(importLines, line)
			if strings.HasPrefix(strings.TrimSpace(line), ")") {
				importBlockStarted = false
				lines = append(lines, importLines...)
				continue
			}
			continue
		}

		// Insert service field into Container struct
		if strings.Contains(line, "type Container struct {") && !insertedField {
			lines = append(lines, line)
			lines = append(lines, fmt.Sprintf("\t%s *%sService", serviceName, serviceName))
			insertedField = true
			continue
		}

		// Insert service initialization inside NewContainer
		if strings.Contains(line, "return &Container{") && !insertedAssign {
			lines = append(lines, line)
			lines = append(lines, fmt.Sprintf("\t\t%s: New%sService(repository.New%sRepository(db)),", serviceName, serviceName, serviceName))
			insertedAssign = true
			continue
		}

		lines = append(lines, line)
	}

	if !alreadyImportedRepository && !insertedImport {

		// Add import outside block if block not found
		for i, l := range lines {
			if strings.HasPrefix(strings.TrimSpace(l), "import") {
				lines = append(lines[:i+1], append([]string{` "` + module + `/internal/repository"`}, lines[i+1:]...)...)
				break
			}
		}
	}

	output := strings.Join(lines, "\n")
	must(os.WriteFile(path, []byte(output), 0644))
	fmt.Println("📦 Updated: internal/service/container.go")
}

func getGoModule() string {
	data, err := os.ReadFile("go.mod")
	must(err)

	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module "))
		}
	}
	return "your-module-name"
}

func isValidGoIdent(s string) bool {
	if s == "" || !isLetter(rune(s[0])) {
		return false
	}
	for _, r := range s {
		if !isLetter(r) && !isDigit(r) {
			return false
		}
	}
	return true
}

func isLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func must(err error) {
	if err != nil {
		fmt.Println("❌", err)
		os.Exit(1)
	}
}

func updateRoutes(serviceName string) {
	// module := getGoModule()
	routesFile := "internal/routes/index.go"

	handlerLine := fmt.Sprintf("\thandler.New%sHandler(svc.%s).Register(api.Group(\"/%ss\"))", serviceName, serviceName, strings.ToLower(serviceName))

	// Check if already registered
	data, err := os.ReadFile(routesFile)
	must(err)
	if strings.Contains(string(data), handlerLine) {
		fmt.Println("📍 Route already exists in index.go")
		return
	}

	lines := []string{}
	file, err := os.Open(routesFile)
	must(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		// Insert route just after api := app.Group(...)
		if strings.Contains(scanner.Text(), "api := app.Group") {
			lines = append(lines, handlerLine)
		}
	}
	must(scanner.Err())

	must(os.WriteFile(routesFile, []byte(strings.Join(lines, "\n")), 0644))
	fmt.Println("📍 Updated: internal/routes/routes.go")
}

func writeHistory(name string, files []string) {
	const historyFile = ".gen_history.json"

	history := map[string]struct {
		CreatedAt string   `json:"created_at"`
		Files     []string `json:"files"`
	}{}

	// Load existing history
	data, _ := os.ReadFile(historyFile)
	_ = json.Unmarshal(data, &history)

	history[strings.ToLower(name)] = struct {
		CreatedAt string   `json:"created_at"`
		Files     []string `json:"files"`
	}{
		CreatedAt: time.Now().Format(time.RFC3339),
		Files:     files,
	}

	historyData, _ := json.MarshalIndent(history, "", "  ")
	_ = os.WriteFile(historyFile, historyData, 0644)
}
