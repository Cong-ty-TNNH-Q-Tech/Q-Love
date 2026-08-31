// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package handlers

import (
	"strings"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CardHandler struct {
	cardService services.CardService
}

func NewCardHandler(cardService services.CardService) *CardHandler {
	return &CardHandler{
		cardService: cardService,
	}
}

func (h *CardHandler) GetCardProfile(c *fiber.Ctx) error {
	targetUserIDStr := c.Params("user_id")
	targetUserID, err := uuid.Parse(targetUserIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"code":    422,
			"message": "ERR_VALIDATION",
			"details": "invalid user_id format",
		})
	}

	profile, err := h.cardService.GetProfile(c.Context(), targetUserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    500,
			"message": "ERR_INTERNAL",
		})
	}

	return c.JSON(profile)
}

type TradeCardRequest struct {
	TargetUserID string `json:"target_user_id"`
	Type         string `json:"type"` // "buy" or "sell"
	Quantity     int    `json:"quantity"`
}

func (h *CardHandler) TradeCard(c *fiber.Ctx) error {
	collectorIDStr := c.Locals("user_id").(string)
	collectorID, err := uuid.Parse(collectorIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    401,
			"message": "ERR_UNAUTHORIZED",
		})
	}

	var req TradeCardRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"code":    422,
			"message": "ERR_VALIDATION",
			"details": err.Error(),
		})
	}

	targetUserID, err := uuid.Parse(req.TargetUserID)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"code":    422,
			"message": "ERR_VALIDATION",
			"details": "invalid target_user_id format",
		})
	}

	err = h.cardService.TradeCard(c.Context(), collectorID, targetUserID, req.Type, req.Quantity)
	if err != nil {
		if strings.Contains(err.Error(), "circuit breaker") {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"code":    503,
				"message": "ERR_CIRCUIT_BREAKER_ACTIVE",
				"details": err.Error(),
			})
		}
		if strings.Contains(err.Error(), "level 5") {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"code":    403,
				"message": "ERR_FORBIDDEN",
				"details": err.Error(),
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    400,
			"message": "ERR_TRADE_FAILED",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "trade successful",
	})
}
