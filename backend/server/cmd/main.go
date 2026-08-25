// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/config"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/api"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/cron"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/repository"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/services"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/pkg/database"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/pkg/logger"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/pkg/redis"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/pkg/storage"
	"github.com/getsentry/sentry-go"
	"github.com/gofiber/contrib/fibersentry"
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func setupApp(ctx context.Context, cfg *config.Config) (*fiber.App, error) {
	// Init Logger & Sentry
	logger.InitLogger(cfg.Environment, cfg.SentryDSN)

	// Init PostgreSQL Database
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		logger.Log.Error("Failed to initialize PostgreSQL Database", zap.Error(err))
		return nil, err
	}

	// Init R2 Storage Client
	r2Client, err := storage.NewR2Client(cfg)
	if err != nil {
		logger.Log.Error("Failed to initialize R2 Storage Client", zap.Error(err))
		return nil, err
	}

	// Init Redis Client
	redisClient, err := redis.NewRedisClient(cfg)
	if err != nil {
		logger.Log.Error("Failed to initialize Redis Client", zap.Error(err))
		return nil, err
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

	// Cron Scheduler
	matchRepo := repository.NewMatchRepository(db)
	islandCronService := services.NewIslandCronService(matchRepo)
	
	clanRepo := repository.NewClanRepository(db)
	landmarkRepo := repository.NewLandmarkRepository(db)
	notifRepo := repository.NewNotificationRepository(db)
	pushService := services.NewPushService()
	walletRepo := repository.NewWalletRepository(db)
	txManager := repository.NewTransactionManager(db)
	clanCronService := services.NewClanCronService(clanRepo, landmarkRepo, notifRepo, pushService, walletRepo, txManager)
	
	scheduler := cron.NewScheduler(clanCronService, islandCronService)
	scheduler.Start()

	// Health Check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Q-Love API is running",
		})
	})

	// Register API Routes
	api.RegisterRoutes(ctx, app, db, r2Client, redisClient, cfg)

	return app, nil
}

func main() {
	cfg := config.LoadConfig()
	
	// Setup context that cancels on interrupt
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	app, err := setupApp(ctx, cfg)
	if err != nil {
		logger.Log.Fatal("Failed to setup app", zap.Error(err))
	}
	
	// Flush buffered events before the program terminates
	defer sentry.Flush(2 * time.Second)
	defer logger.Log.Sync()
	
	port := cfg.Port
	if port == "" {
		port = "3000"
	}

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		logger.Log.Info("Gracefully shutting down...")
		cancel()
		_ = app.Shutdown()
	}()

	logger.Log.Info("Starting server...", zap.String("port", port))
	if err := app.Listen(":" + port); err != nil {
		logger.Log.Fatal("Server stopped", zap.Error(err))
	}
}
