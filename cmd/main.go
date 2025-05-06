package main

import (
	"log"

	"github.com/msdevbytes/go-microkit/internal/app"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
