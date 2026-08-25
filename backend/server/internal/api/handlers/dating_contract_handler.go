// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package handlers

import (
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"time"
)

type DatingContractHandler struct {
	contractService services.DatingContractService
}

func NewDatingContractHandler(contractService services.DatingContractService) *DatingContractHandler {
	return &DatingContractHandler{contractService: contractService}
}

func (h *DatingContractHandler) CreateContract(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req struct {
		TargetUserID    string    `json:"target_user_id"`
		DepositAmount   float64   `json:"deposit_amount"`
		AppointmentTime time.Time `json:"appointment_time"`
		LocationNote    string    `json:"location_note"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	targetUserID, err := uuid.Parse(req.TargetUserID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid target_user_id format"})
	}

	contract, err := h.contractService.CreateContract(c.Context(), userID, targetUserID, req.DepositAmount, req.AppointmentTime)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(contract)
}

func (h *DatingContractHandler) AcceptContract(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	contractIDStr := c.Params("contract_id")
	contractID, err := uuid.Parse(contractIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid contract_id format"})
	}

	contract, err := h.contractService.AcceptContract(c.Context(), contractID, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(contract)
}

func (h *DatingContractHandler) CancelContract(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	contractIDStr := c.Params("contract_id")
	contractID, err := uuid.Parse(contractIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid contract_id format"})
	}

	var req struct {
		Reason string `json:"reason"`
	}

	if err := c.BodyParser(&req); err != nil {
		// reason is optional
	}

	err = h.contractService.CancelContract(c.Context(), contractID, userID, req.Reason)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"cancelled_by":    userID.String(),
		"penalty_applied": true,
		"message":         "Cancelled successfully",
	})
}

func (h *DatingContractHandler) ScanContract(c *fiber.Ctx) error {
	contractIDStr := c.Params("contract_id")
	contractID, err := uuid.Parse(contractIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid contract_id format"})
	}

	var req struct {
		QRToken string `json:"qr_token"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err = h.contractService.ScanContract(c.Context(), contractID, req.QRToken)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "completed"})
}
