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
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/pkg/logger"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/pkg/esms"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/pkg/storage"
	"github.com/gofiber/websocket/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"context"
)

func RegisterRoutes(ctx context.Context, app *fiber.App, db *gorm.DB, r2Client *storage.R2Client, redisClient *redis.Client, cfg *config.Config) {
	wingmanRepo := repository.NewWingmanRepository(db)
	walletRepo := repository.NewWalletRepository(db)
	shameRepo := repository.NewShameRepository(db)
	txManager := repository.NewTransactionManager(db)
	matchRepo := repository.NewMatchRepository(db)

	wingmanService := services.NewWingmanService(wingmanRepo, walletRepo, txManager, matchRepo)
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
	go hub.Run(ctx)
	chatHandler := handlers.NewChatHandler(chatService, hub)

	matchService := services.NewMatchService(matchRepo)
	matchHandler := handlers.NewMatchHandler(matchService)
	userPremRepo := repository.NewUserPremiumRepository(db)
	
	violationRepo := repository.NewUserViolationRepository(db)
	nsfwService := services.NewNSFWService(cfg)
	notificationRepo := repository.NewNotificationRepository(db)
	notificationService := services.NewNotificationService(notificationRepo, redisClient, cfg.FCMKey)
	locketService := services.NewLocketService(chatRepo, matchRepo, violationRepo, nsfwService, notificationService, r2Client)
	locketHandler := handlers.NewLocketHandler(locketService)

	iapService := services.NewIAPService(txManager, walletRepo, userPremRepo)
	webhookHandler := handlers.NewWebhookHandler(cfg, iapService)

	spotifyService := services.NewSpotifyService()
	vibeHandler := handlers.NewVibeHandler(spotifyService)
	stealRepo := repository.NewCardStealRepository(db)
	minigameService := services.NewMinigameService(stealRepo, walletRepo, txManager)
	minigameHandler := handlers.NewMinigameHandler(minigameService)

	// Auth routes setup
	esmsClient := esms.NewClient(cfg.ESMSAPIKey, cfg.ESMSSecretKey)
	userRepo := repository.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, esmsClient, redisClient, cfg.JWTSecret)
	authHandler := handlers.NewAuthHandler(authService)

	// API v1 group
	v1 := app.Group("/api/v1")

	// Auth routes
	authGroup := v1.Group("/auth")
	authGroup.Post("/send-otp", authHandler.SendOTP)
	authGroup.Post("/verify-otp", authHandler.VerifyOTP)
	authGroup.Post("/refresh", authHandler.RefreshToken)

	// Wingman routes
	wingmanGroup := v1.Group("/wingman", middleware.JWTMiddleware(cfg.JWTSecret))
	wingmanGroup.Post("/referrals", wingmanHandler.CreateReferral)
	wingmanGroup.Post("/referrals/:id/accept", wingmanHandler.AcceptReferral)

	// Shame (Wall of Shame) routes
	shameGroup := v1.Group("/shames", middleware.JWTMiddleware(cfg.JWTSecret))
	shameGroup.Get("/", shameHandler.GetActiveShames)
	shameGroup.Post("/:id/tomato", shameHandler.ThrowTomato)

	// AI Wingman routes
	aiService := services.NewAIWingmanService(chatRepo, cfg.OpenAIAPIKey)
	aiHandler := handlers.NewAIWingmanHandler(aiService)
	aiGroup := v1.Group("/ai", middleware.JWTMiddleware(cfg.JWTSecret))
	aiGroup.Post("/suggest", aiHandler.SuggestReplies)

	// Ex-Rating routes
	exRatingRepo := repository.NewExRatingRepository(db)
	exRatingService := services.NewExRatingService(exRatingRepo, walletRepo, txManager, chatRepo, matchRepo)
	exRatingHandler := handlers.NewExRatingHandler(exRatingService)
	v1.Post("/ex-ratings", middleware.JWTMiddleware(cfg.JWTSecret), exRatingHandler.SubmitRating)
	v1.Get("/users/:user_id/ex-rating", middleware.JWTMiddleware(cfg.JWTSecret), exRatingHandler.ViewRating)

	// Upload routes
	uploadHandler := handlers.NewUploadHandler(r2Client)
	uploadGroup := v1.Group("/upload", middleware.JWTMiddleware(cfg.JWTSecret))
	uploadGroup.Post("/presigned-url", uploadHandler.GenerateUploadURL)

	// Clan routes
	clanGroup := v1.Group("/clans", middleware.JWTMiddleware(cfg.JWTSecret))
	clanGroup.Post("/", clanHandler.CreateClan)

	// Chat routes
	chatGroup := v1.Group("/chat")
	chatGroup.Get("/ws", middleware.JWTMiddleware(cfg.JWTSecret), chatHandler.Upgrade, websocket.New(chatHandler.WSHandler))
	chatGroup.Post("/messages", middleware.JWTMiddleware(cfg.JWTSecret), chatHandler.SendMessage)
	chatGroup.Get("/messages/:match_id", middleware.JWTMiddleware(cfg.JWTSecret), chatHandler.GetMessages)

	// Locket routes
	locketGroup := v1.Group("/locket", middleware.JWTMiddleware(cfg.JWTSecret))
	locketGroup.Post("/send", middleware.LocketRateLimiter(), locketHandler.SendLocket)
	
	// Auction routes
	auctionRepo := repository.NewAuctionRepository(db)
	chatLockRepo := repository.NewChatLockRepository(db)
	auctionService := services.NewAuctionService(auctionRepo, walletRepo, txManager, userRepo, chatLockRepo)
	auctionHandler := handlers.NewAuctionHandler(auctionService)
	auctionGroup := v1.Group("/auctions", middleware.JWTMiddleware(cfg.JWTSecret))
	auctionGroup.Get("/active", auctionHandler.GetActiveAuctions)
	auctionGroup.Post("/:id/bid", auctionHandler.PlaceBid)

	// Admin routes
	courtCaseRepo := repository.NewCourtCaseRepository(db)
	adminService := services.NewAdminService(violationRepo, courtCaseRepo, r2Client, walletRepo, txManager)
	adminHandler := handlers.NewAdminHandler(adminService)
	adminGroup := app.Group("/admin/v1", middleware.AdminMiddleware(cfg.JWTSecret))
	adminGroup.Get("/violations", adminHandler.GetViolations)
	adminGroup.Post("/users/:id/ban", adminHandler.BanUser)
	adminGroup.Delete("/violations/:id/media", adminHandler.DeleteViolationMedia)
	adminGroup.Post("/court/:id/override", adminHandler.OverrideCourtCase)

	// Webhooks
	webhookGroup := v1.Group("/webhooks")
	webhookGroup.Post("/revenuecat", webhookHandler.HandleRevenueCat)

	// Device routes
	deviceHandler := handlers.NewDeviceHandler(redisClient)
	deviceGroup := v1.Group("/devices", middleware.JWTMiddleware(cfg.JWTSecret))
	deviceGroup.Post("/token", deviceHandler.RegisterFCMToken)

	// Vibe Check (Spotify)
	vibeGroup := v1.Group("/vibe", middleware.JWTMiddleware(cfg.JWTSecret))
	vibeGroup.Get("/status", vibeHandler.Status)
	vibeGroup.Get("/current-track", vibeHandler.CurrentTrack)
	vibeGroup.Post("/match", vibeHandler.Match)
	// Minigame Steal routes
	stealGroup := v1.Group("/minigame/steal", middleware.JWTMiddleware(cfg.JWTSecret))
	stealGroup.Post("/init", minigameHandler.InitSteal)
	stealGroup.Post("/submit", minigameHandler.SubmitStealResult)

	// Match API
	matchGroup := v1.Group("/matches", middleware.JWTMiddleware(cfg.JWTSecret))
	matchGroup.Delete("/:match_id", matchHandler.Unmatch)

	// Court System
	courtRepo := repository.NewCourtRepository(db)
	courtService := services.NewCourtService(courtRepo, matchRepo, redisClient, walletRepo, txManager)
	courtHandler := handlers.NewCourtHandler(courtService)

	courtGroup := v1.Group("/court", middleware.JWTMiddleware(cfg.JWTSecret))
	courtGroup.Post("/cases", courtHandler.FileLawsuit)
	courtGroup.Get("/feed", courtHandler.GetFeed)
	courtGroup.Post("/:case_id/vote", courtHandler.VoteCase)
	courtGroup.Post("/:case_id/withdraw", courtHandler.WithdrawCase)

	// Start Court Worker
	if redisClient != nil {
		courtWorker := services.NewCourtWorker(courtRepo, violationRepo, redisClient, logger.Log, walletRepo, txManager)
		courtWorker.Start(ctx)
	}

	// Vouchers
	voucherRepo := repository.NewVoucherRepository(db)
	voucherService := services.NewVoucherService(voucherRepo, walletRepo, txManager)
	voucherHandler := handlers.NewVoucherHandler(voucherService)
	adminVoucherHandler := handlers.NewAdminVoucherHandler(voucherService)

	voucherGroup := v1.Group("/vouchers", middleware.JWTMiddleware(cfg.JWTSecret))
	voucherGroup.Get("/", voucherHandler.GetAvailableVouchers)
	voucherGroup.Post("/redeem", voucherHandler.RedeemVoucher)

	// Admin
	// Admin Vouchers
	adminGroup.Get("/vouchers", adminVoucherHandler.GetVouchers)
	adminGroup.Post("/vouchers", adminVoucherHandler.CreateVoucher)
	adminGroup.Delete("/vouchers/:id", adminVoucherHandler.DeleteVoucher)
}
