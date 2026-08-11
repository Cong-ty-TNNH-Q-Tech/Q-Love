package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/config"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/api"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/pkg/storage"
	"gorm.io/gorm"
)

var app *fiber.App

func setupApp() *fiber.App {
	// 1. Load config
	cfg := config.LoadConfig()

	// 2. Init R2 Storage Client
	r2Client, err := storage.NewR2Client(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize R2 Storage Client: %v", err)
	}

	a := fiber.New(fiber.Config{
		AppName: "Q-Love Backend v1.0",
	})

	// Middleware
	a.Use(logger.New())

	var db *gorm.DB // In a real setup, connect to PostgreSQL here
	
	// Health Check
	a.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Q-Love API is running",
		})
	})

	// Register API Routes
	api.RegisterRoutes(a, db, r2Client)

	return a
}

func main() {
	cfg := config.LoadConfig()
	app = setupApp()
	
	port := cfg.Port
	if port == "" {
		port = "3000"
	}

	log.Printf("Starting server on port %s...", port)
	if err := app.Listen(":" + port); err != nil {
		log.Printf("Server stopped: %v", err)
	}
}
