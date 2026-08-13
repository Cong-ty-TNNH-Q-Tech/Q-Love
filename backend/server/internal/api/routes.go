// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/config"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/api/handlers"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/middleware"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/repository"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/services"
	chatws "github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/websocket"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/pkg/storage"
	"github.com/gofiber/websocket/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"context"
)

func RegisterRoutes(app *fiber.App, db *gorm.DB, r2Client *storage.R2Client, redisClient *redis.Client, cfg *config.Config) {
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

	// Chat & Websocket
	chatRepo := repository.NewChatMessageRepository(db)
	chatService := services.NewChatService(chatRepo)
	hub := chatws.NewHub(redisClient)
	go hub.Run(context.Background())
	chatHandler := handlers.NewChatHandler(chatService, hub)

	matchRepo := repository.NewMatchRepository(db)
	userPremRepo := repository.NewUserPremiumRepository(db)
	locketService := services.NewLocketService(chatRepo, matchRepo, r2Client)
	locketHandler := handlers.NewLocketHandler(locketService)

	iapService := services.NewIAPService(txManager, walletRepo, userPremRepo)
	webhookHandler := handlers.NewWebhookHandler(cfg, iapService)

	stealRepo := repository.NewCardStealRepository(db)
	minigameService := services.NewMinigameService(stealRepo, walletRepo, txManager)
	minigameHandler := handlers.NewMinigameHandler(minigameService)

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

	// Chat routes
	chatGroup := v1.Group("/chat")
	chatGroup.Get("/ws", chatHandler.Upgrade, websocket.New(chatHandler.WSHandler))
	chatGroup.Post("/messages", middleware.JWTMiddleware(""), chatHandler.SendMessage)
	chatGroup.Get("/messages/:match_id", middleware.JWTMiddleware(""), chatHandler.GetMessages)

	// Locket routes
	locketGroup := v1.Group("/locket", middleware.JWTMiddleware(""))
	locketGroup.Post("/send", middleware.LocketRateLimiter(), locketHandler.SendLocket)
	// Webhooks
	webhookGroup := v1.Group("/webhooks")
	webhookGroup.Post("/revenuecat", webhookHandler.HandleRevenueCat)

	// Minigame Steal routes
	stealGroup := v1.Group("/minigame/steal", middleware.JWTMiddleware(""))
	stealGroup.Post("/init", minigameHandler.InitSteal)
	stealGroup.Post("/submit", minigameHandler.SubmitStealResult)
}
