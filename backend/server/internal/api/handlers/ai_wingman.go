// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/services"
)

type AIWingmanHandler struct {
	service services.AIWingmanService
}

func NewAIWingmanHandler(service services.AIWingmanService) *AIWingmanHandler {
	return &AIWingmanHandler{service: service}
}

func (h *AIWingmanHandler) GetSuggestions(c *fiber.Ctx) error {
	matchIDStr := c.Params("id")
	matchID, err := uuid.Parse(matchIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid match ID"})
	}

	suggestions, err := h.service.GetSuggestions(c.Context(), matchID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate suggestions"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"suggestions": suggestions,
	})
}
