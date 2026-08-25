// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package handlers

import (
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type VoucherHandler struct {
	voucherService services.VoucherService
}

func NewVoucherHandler(s services.VoucherService) *VoucherHandler {
	return &VoucherHandler{voucherService: s}
}

type redeemRequest struct {
	Brand   string `json:"brand"`
	ValueXu int    `json:"value_xu"`
}

func (h *VoucherHandler) RedeemVoucher(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid user ID"})
	}

	var req redeemRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "dữ liệu không hợp lệ"})
	}

	if req.Brand == "" || req.ValueXu <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "thương hiệu và giá trị xu phải hợp lệ"})
	}

	voucher, err := h.voucherService.RedeemVoucher(c.Context(), userID, req.Brand, req.ValueXu)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "đổi voucher thành công",
		"voucher": fiber.Map{
			"brand":      voucher.Brand,
			"code":       voucher.Code,
			"value_xu":   voucher.ValueXu,
			"expires_at": voucher.ExpiresAt,
		},
	})
}

// GetAvailableGroups returns groups of available vouchers
// For simplicity in this PR, we'll just expose GetVouchers directly from Admin API instead of complex grouping here,
// or we can implement a basic count by brand and value.
func (h *VoucherHandler) GetAvailableVouchers(c *fiber.Ctx) error {
	// A simple workaround for the user endpoint without adding a new Repo method
	vouchers, err := h.voucherService.GetVouchers(c.Context(), 1000, 0)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	
	// Group available vouchers
	type group struct {
		Brand   string `json:"brand"`
		ValueXu int    `json:"value_xu"`
		Count   int    `json:"count"`
	}
	counts := make(map[string]*group)
	for _, v := range vouchers {
		if v.Status == "available" {
			key := v.Brand + string(rune(v.ValueXu))
			if _, ok := counts[key]; !ok {
				counts[key] = &group{Brand: v.Brand, ValueXu: v.ValueXu, Count: 0}
			}
			counts[key].Count++
		}
	}

	result := []group{}
	for _, g := range counts {
		result = append(result, *g)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}
