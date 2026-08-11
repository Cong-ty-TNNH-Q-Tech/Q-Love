package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/api/handlers"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/repository"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/services"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/pkg/storage"
	"gorm.io/gorm"
)

func RegisterRoutes(app *fiber.App, db *gorm.DB, r2Client *storage.R2Client) {
	api := app.Group("/api/v1")

	// Wingman routes
	wingmanRepo := repository.NewWingmanRepository(db)
	walletRepo := repository.NewWalletRepository(db)
	txManager := repository.NewTransactionManager(db)
	wingmanService := services.NewWingmanService(wingmanRepo, walletRepo, txManager)
	wingmanHandler := handlers.NewWingmanHandler(wingmanService)
	wingmanGroup := api.Group("/wingmans")
	wingmanGroup.Post("/referral", wingmanHandler.CreateReferral)
	wingmanGroup.Post("/referral/:id/accept", wingmanHandler.AcceptReferral)
	// Upload routes
	uploadHandler := handlers.NewUploadHandler(r2Client)
	uploadGroup := api.Group("/upload")
	uploadGroup.Post("/presigned-url", uploadHandler.GenerateUploadURL)
}
