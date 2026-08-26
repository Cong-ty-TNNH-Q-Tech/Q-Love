// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package handlers

import (
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AIWingmanHandler struct {
	aiService services.AIWingmanService
}

func NewAIWingmanHandler(aiService services.AIWingmanService) *AIWingmanHandler {
	return &AIWingmanHandler{aiService: aiService}
}

// SuggestReplies POST /api/v1/ai/suggest
func (h *AIWingmanHandler) SuggestReplies(c *fiber.Ctx) error {
	_, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	type Request struct {
		MatchID string `json:"match_id"`
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	matchID, err := uuid.Parse(req.MatchID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid match_id"})
	}

	replies, err := h.aiService.SuggestReplies(c.Context(), matchID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"replies": replies,
	})
}
