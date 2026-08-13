// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

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
	wingmanRepo := repository.NewWingmanRepository(db)
	walletRepo := repository.NewWalletRepository(db)
	shameRepo := repository.NewShameRepository(db)
	txManager := repository.NewTransactionManager(db)

	walletService := services.NewWalletService(walletRepo, txManager)
	wingmanService := services.NewWingmanService(wingmanRepo, walletRepo, txManager)
	shameService := services.NewShameService(shameRepo, walletRepo, txManager)

	walletHandler := handlers.NewWalletHandler(walletService)
	wingmanHandler := handlers.NewWingmanHandler(wingmanService)
	shameHandler := handlers.NewShameHandler(shameService)

	// API v1 group
	v1 := app.Group("/api/v1")

	// Wallet routes
	v1.Post("/wallets/deposit", walletHandler.Deposit)
	v1.Get("/wallets/:user_id", walletHandler.GetBalance)

	// Wingman routes
	v1.Post("/wingman/referrals", wingmanHandler.CreateReferral)
	v1.Post("/wingman/referrals/:id/accept", wingmanHandler.AcceptReferral)

	// Shame (Wall of Shame) routes
	v1.Get("/shames", shameHandler.GetActiveShames)
	v1.Post("/shames/:id/tomato", shameHandler.ThrowTomato)

	// Upload routes
	uploadHandler := handlers.NewUploadHandler(r2Client)
	uploadGroup := v1.Group("/upload")
	uploadGroup.Post("/presigned-url", uploadHandler.GenerateUploadURL)
}
