// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package handlers

import (
	"strconv"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AdminHandler struct {
	adminService services.AdminService
}

func NewAdminHandler(adminService services.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

// GetViolations returns paginated violations
func (h *AdminHandler) GetViolations(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	violations, total, err := h.adminService.GetViolations(c.Context(), page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data":  violations,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// BanUser bans a user
func (h *AdminHandler) BanUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	if err := h.adminService.BanUser(c.Context(), userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "User banned successfully"})
}

// DeleteViolationMedia removes the image from R2 and dismisses the violation
func (h *AdminHandler) DeleteViolationMedia(c *fiber.Ctx) error {
	idParam := c.Params("id")
	violationID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid violation ID"})
	}

	type DeleteReq struct {
		ObjectKey string `json:"object_key"`
	}
	var req DeleteReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.adminService.DeleteViolationMedia(c.Context(), violationID, req.ObjectKey); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Violation media deleted successfully"})
}

// OverrideCourtCase overrides the status of a court case
func (h *AdminHandler) OverrideCourtCase(c *fiber.Ctx) error {
	idParam := c.Params("id")
	caseID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid case ID"})
	}

	type OverrideReq struct {
		Status string `json:"status"` // e.g. "dismissed", "guilty", "innocent"
	}
	var req OverrideReq
	if err := c.BodyParser(&req); err != nil || req.Status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.adminService.OverrideCourtCase(c.Context(), caseID, req.Status); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Court case overridden successfully"})
}
