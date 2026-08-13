// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/services"
)

type ClanHandler struct {
	service services.ClanService
}

func NewClanHandler(service services.ClanService) *ClanHandler {
	return &ClanHandler{service: service}
}

type CreateClanRequest struct {
	Name    string `json:"name"`
	Slogan  string `json:"slogan"`
	LogoURL string `json:"logo_url"`
}

func (h *ClanHandler) CreateClan(c *fiber.Ctx) error {
	var req CreateClanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Clan name is required"})
	}

	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	clan, err := h.service.CreateClan(c.Context(), userID, req.Name, req.Slogan, req.LogoURL)
	if err != nil {
		if err.Error() == "ERR_CLAN_NAME_TAKEN" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"code": 400, "message": "ERR_CLAN_NAME_TAKEN"})
		}
		if err.Error() == "insufficient balance" || err.Error() == "wallet not found" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create clan"})
	}

	return c.Status(fiber.StatusCreated).JSON(clan)
}
