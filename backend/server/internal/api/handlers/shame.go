// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"qlove/internal/services"
)

type ShameHandler struct {
	service services.ShameService
}

func NewShameHandler(service services.ShameService) *ShameHandler {
	return &ShameHandler{service: service}
}

func (h *ShameHandler) GetActiveShames(c *fiber.Ctx) error {
	limitStr := c.Query("limit", "10")
	offsetStr := c.Query("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	shames, err := h.service.GetActiveShames(c.Context(), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch shames"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": shames,
	})
}

func (h *ShameHandler) ThrowTomato(c *fiber.Ctx) error {
	idParam := c.Params("id")
	shameID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid shame ID"})
	}

	// Simulated authenticated user
	throwerID := uuid.New()
	if authHeader := c.Get("X-User-ID"); authHeader != "" {
		if parsed, err := uuid.Parse(authHeader); err == nil {
			throwerID = parsed
		}
	}

	err = h.service.ThrowTomato(c.Context(), throwerID, shameID)
	if err != nil {
		if err.Error() == "insufficient balance to throw a tomato" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to throw tomato"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Tomato thrown successfully",
	})
}
