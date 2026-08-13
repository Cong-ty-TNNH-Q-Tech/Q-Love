// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package main

import (
	"log"
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/config"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/api"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/pkg/logger"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/pkg/storage"
	"github.com/getsentry/sentry-go"
	"github.com/gofiber/contrib/fibersentry"
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func setupApp(cfg *config.Config) *fiber.App {
	// Init Logger & Sentry
	logger.InitLogger(cfg.Environment, cfg.SentryDSN)

	// Init R2 Storage Client
	r2Client, err := storage.NewR2Client(cfg)
	if err != nil {
		logger.Log.Fatal("Failed to initialize R2 Storage Client", zap.Error(err))
	}

	app := fiber.New(fiber.Config{
		AppName: "Q-Love Backend v1.0",
	})

	// Middleware
	app.Use(fiberzap.New(fiberzap.Config{
		Logger: logger.Log,
	}))
	app.Use(fibersentry.New(fibersentry.Config{
		Repanic:         true,
		WaitForDelivery: true,
	}))

	var db *gorm.DB // In a real setup, connect to PostgreSQL here
	
	// Health Check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Q-Love API is running",
		})
	})

	// Register API Routes
	api.RegisterRoutes(app, db, r2Client)

	return app
}

func main() {
	cfg := config.LoadConfig()
	app := setupApp(cfg)
	
	// Flush buffered events before the program terminates
	defer sentry.Flush(2 * time.Second)
	defer logger.Log.Sync()
	
	port := cfg.Port
	if port == "" {
		port = "3000"
	}

	logger.Log.Info("Starting server...", zap.String("port", port))
	if err := app.Listen(":" + port); err != nil {
		logger.Log.Fatal("Server stopped", zap.Error(err))
	}
}
