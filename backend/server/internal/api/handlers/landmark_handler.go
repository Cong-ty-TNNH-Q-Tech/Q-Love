// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package handlers

import (
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type LandmarkHandler struct {
	landmarkService services.LandmarkService
}

func NewLandmarkHandler(landmarkService services.LandmarkService) *LandmarkHandler {
	return &LandmarkHandler{landmarkService: landmarkService}
}

func (h *LandmarkHandler) CheckIn(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	landmarkIDStr := c.Params("landmark_id")
	landmarkID, err := uuid.Parse(landmarkIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid landmark_id format"})
	}

	var req struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		IsMocked  bool    `json:"is_mocked"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	_, err = h.landmarkService.CheckIn(c.Context(), userID, landmarkID, req.Latitude, req.Longitude, req.IsMocked)
	if err != nil {
		if err == services.ErrFakeGPS {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"code":    403,
				"message": "ERR_FAKE_GPS_DETECTED",
			})
		}
		if err == services.ErrOutOfRange {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"code":    400,
				"message": "ERR_OUT_OF_RANGE",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Check-in successful. +10 points",
	})
}
