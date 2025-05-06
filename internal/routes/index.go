package routes

import (
	"os"
	"time"

	"github.com/msdevbytes/go-microkit/internal/handler"
	"github.com/msdevbytes/go-microkit/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/shahbazkrispx/pkgcommon"
)

func Register(app *fiber.App, svc *service.Container) {

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(pkgcommon.ResponseBuilder(true, "Welcome to Event Service API", time.Now().Format(time.RFC3339), nil))
	})

	api := app.Group(os.Getenv("API_ROUTE_VERSION"))

	handler.NewDefaultHandler().Register(api.Group("/"))
}
