package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName: "Q-Love Backend v1.0",
	})

	// Middleware
	app.Use(logger.New())

	// Routes
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"message": "Q-Love API is running",
		})
	})

	log.Fatal(app.Listen(":3000"))
}
