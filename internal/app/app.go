package app

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/msdevbytes/go-microkit/internal/db"
	"github.com/msdevbytes/go-microkit/internal/routes"
	"github.com/msdevbytes/go-microkit/internal/service"
	"github.com/msdevbytes/go-microkit/pkg/logger"
)

func Start() error {

	time.Local = time.UTC

	if err := db.OpenConnection(); err != nil {
		logger.Danger("Please check mysql server")
		return err
	}

	svc := service.NewContainer(db.Conn)

	app := fiber.New(fiber.Config{
		AppName: os.Getenv("APP_NAME"),
	})

	routes.Register(app, svc)

	return app.Listen(":" + os.Getenv("PORT"))
}
