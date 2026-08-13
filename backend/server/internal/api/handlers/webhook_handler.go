package handlers

import (
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/config"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/services"
	"github.com/gofiber/fiber/v2"
)

type WebhookHandler struct {
	cfg        *config.Config
	iapService services.IAPService
}

func NewWebhookHandler(cfg *config.Config, iapService services.IAPService) *WebhookHandler {
	return &WebhookHandler{
		cfg:        cfg,
		iapService: iapService,
	}
}

func (h *WebhookHandler) HandleRevenueCat(c *fiber.Ctx) error {
	// Verify Authorization header
	authHeader := c.Get("Authorization")
	if authHeader != "Bearer "+h.cfg.RevenueCatWebhookSecret && authHeader != h.cfg.RevenueCatWebhookSecret {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var payload struct {
		Event services.RevenueCatEvent `json:"event"`
	}

	if err := c.BodyParser(&payload); err != nil {
		// RevenueCat expects 2xx so it doesn't retry on malformed payloads that we can't parse anyway.
		// Returning 200 even on parse error prevents endless retries.
		// But returning 400 is also acceptable if we want them to log it. Let's return 400.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	err := h.iapService.ProcessRevenueCatWebhook(c.Context(), payload.Event)
	if err != nil {
		if err == services.ErrInvalidWebhookPayload {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		// In case of internal errors (DB down), we return 500 so RevenueCat retries later.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	// Always return 200 OK so RevenueCat marks webhook as successful
	return c.SendStatus(fiber.StatusOK)
}
