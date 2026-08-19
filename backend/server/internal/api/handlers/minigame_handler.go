// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package handlers

import (
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type MinigameHandler struct {
	minigameService services.MinigameService
}

func NewMinigameHandler(minigameService services.MinigameService) *MinigameHandler {
	return &MinigameHandler{minigameService: minigameService}
}

func (h *MinigameHandler) InitSteal(c *fiber.Ctx) error {
	attackerIDStr := c.Locals("user_id").(string)
	attackerID, err := uuid.Parse(attackerIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	var req struct {
		DefenderID   string `json:"defender_id"`
		TargetCardID string `json:"target_card_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	defID, err := uuid.Parse(req.DefenderID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid defender id"})
	}
	targetID, err := uuid.Parse(req.TargetCardID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid target card id"})
	}

	steal, err := h.minigameService.InitSteal(c.Context(), attackerID, defID, targetID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "steal session initiated",
		"data":    steal,
	})
}

func (h *MinigameHandler) SubmitStealResult(c *fiber.Ctx) error {
	attackerIDStr := c.Locals("user_id").(string)
	attackerID, err := uuid.Parse(attackerIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	var req struct {
		StealID string `json:"steal_id"`
		IsWin   bool   `json:"is_win"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	stealID, err := uuid.Parse(req.StealID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid steal id"})
	}

	err = h.minigameService.SubmitStealResult(c.Context(), stealID, attackerID, req.IsWin)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "steal result processed successfully",
	})
}
