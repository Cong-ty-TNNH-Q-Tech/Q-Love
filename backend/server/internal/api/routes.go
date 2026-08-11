package api

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/api/handlers"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/services"
)

func RegisterRoutes(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api/v1")

	// Wingman routes
	wingmanService := services.NewWingmanService(db)
	wingmanHandler := handlers.NewWingmanHandler(wingmanService)
	wingmanGroup := api.Group("/wingmans")
	wingmanGroup.Post("/referral", wingmanHandler.CreateReferral)
	wingmanGroup.Post("/referral/:id/accept", wingmanHandler.AcceptReferral)
}
