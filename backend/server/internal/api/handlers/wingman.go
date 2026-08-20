// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/services"
)

type WingmanHandler struct {
	service services.WingmanService
}

func NewWingmanHandler(service services.WingmanService) *WingmanHandler {
	return &WingmanHandler{service: service}
}

type CreateReferralRequest struct {
	Target1ID uuid.UUID `json:"target1_id"`
	Target2ID uuid.UUID `json:"target2_id"`
}

func (h *WingmanHandler) CreateReferral(c *fiber.Ctx) error {
	var req CreateReferralRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	wingmanID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	referral, err := h.service.CreateReferral(c.Context(), wingmanID, req.Target1ID, req.Target2ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create referral link"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":         referral.ID,
		"deep_link":  referral.DeepLink,
		"expires_at": referral.ExpiresAt,
	})
}

func (h *WingmanHandler) AcceptReferral(c *fiber.Ctx) error {
	idParam := c.Params("id")
	referralID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid referral ID"})
	}

	acceptingUserID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	referral, err := h.service.AcceptReferral(c.Context(), referralID, acceptingUserID)
	if err != nil {
		if errors.Is(err, services.ErrReferralNotFound) ||
			errors.Is(err, services.ErrUserNotInReferral) ||
			errors.Is(err, services.ErrReferralNotPending) ||
			errors.Is(err, services.ErrReferralExpired) ||
			errors.Is(err, services.ErrWingmanCannotReferSelf) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to accept referral"})
	}

	// When accepted, we simulate that we trigger commission to wingman if they date. 
	// The process commission usually runs in a worker or another webhook, but we can expose it or run it synchronously for testing:
	// s.service.ProcessCommission(c.Context(), referral.ID)
	// For this endpoint, we just return matched.

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"match_id": uuid.New(), // dummy match ID generated
		"status":   referral.Status,
	})
}
