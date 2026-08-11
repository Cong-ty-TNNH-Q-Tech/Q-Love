package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/gorm"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/api"
)

var app *fiber.App

func setupApp() *fiber.App {
	a := fiber.New(fiber.Config{
		AppName: "Q-Love Backend v1.0",
	})

	// Middleware
	a.Use(logger.New())

	var db *gorm.DB // In a real setup, connect to PostgreSQL here
	api.RegisterRoutes(a, db)

	// Routes
	a.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Q-Love API is running",
		})
	})

	a.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	a.Get("/version", func(c *fiber.Ctx) error {
		return c.SendString("v1.0.0")
	})

	return a
}

func main() {
	app = setupApp()
	if err := app.Listen(":3000"); err != nil {
		log.Printf("Server stopped: %v", err)
	}
}
