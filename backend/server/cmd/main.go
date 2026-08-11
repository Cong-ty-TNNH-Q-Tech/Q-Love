package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var app *fiber.App

func setupApp() *fiber.App {
	a := fiber.New(fiber.Config{
		AppName: "Q-Love Backend v1.0",
	})

	// Middleware
	a.Use(logger.New())

	// Routes
	a.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Q-Love API is running",
		})
	})

	return a
}

func main() {
	app = setupApp()
	if err := app.Listen(":3000"); err != nil {
		log.Printf("Server stopped: %v", err)
	}
}
