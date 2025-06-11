
---

## ⚙️ Requirements

- Go 1.20+
- [mockgen](https://github.com/golang/mock)  
  `go install github.com/golang/mock/mockgen@latest`
- [richgo](https://github.com/kyoh86/richgo)  
  `go install github.com/kyoh86/richgo@latest`
- [air (live reload)](https://github.com/cosmtrek/air)  
  `go install github.com/cosmtrek/air@latest`
- [golangci-lint](https://golangci-lint.run)  
  `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`

---

## 🚀 Quickstart

```bash
# 1. Clone the repo
git clone https://github.com/msdevbytes/go-microkit.git
cd go-microkit

# 2. Install required tools (first time)
go install github.com/golang/mock/mockgen@latest
go install github.com/kyoh86/richgo@latest
go install github.com/cosmtrek/air@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 3. Generate a new service (e.g., RSVP)
make gen-service name=RSVPService

# 4. Run the server with hot reload
make run

# 5. Run tests and lint checks
make test
make lint
