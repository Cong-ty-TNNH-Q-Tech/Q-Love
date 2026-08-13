// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/api/handlers"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/middleware"
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

	wingmanService := services.NewWingmanService(wingmanRepo, walletRepo, txManager)
	shameService := services.NewShameService(shameRepo, walletRepo, txManager)
	clanRepo := repository.NewClanRepository(db)
	clanService := services.NewClanService(clanRepo, walletRepo, txManager)

	wingmanHandler := handlers.NewWingmanHandler(wingmanService)
	shameHandler := handlers.NewShameHandler(shameService)
	clanHandler := handlers.NewClanHandler(clanService)

	// API v1 group
	v1 := app.Group("/api/v1")

	// Wingman routes
	wingmanGroup := v1.Group("/wingman", middleware.JWTMiddleware(""))
	wingmanGroup.Post("/referrals", wingmanHandler.CreateReferral)
	wingmanGroup.Post("/referrals/:id/accept", wingmanHandler.AcceptReferral)

	// Shame (Wall of Shame) routes
	shameGroup := v1.Group("/shames", middleware.JWTMiddleware(""))
	shameGroup.Get("/", shameHandler.GetActiveShames)
	shameGroup.Post("/:id/tomato", shameHandler.ThrowTomato)

	// Upload routes
	uploadHandler := handlers.NewUploadHandler(r2Client)
	uploadGroup := v1.Group("/upload", middleware.JWTMiddleware(""))
	uploadGroup.Post("/presigned-url", uploadHandler.GenerateUploadURL)

	// Clan routes
	clanGroup := v1.Group("/clans", middleware.JWTMiddleware(""))
	clanGroup.Post("/", clanHandler.CreateClan)

}
