// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package handlers

import (
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AdminVoucherHandler struct {
	voucherService services.VoucherService
}

func NewAdminVoucherHandler(s services.VoucherService) *AdminVoucherHandler {
	return &AdminVoucherHandler{voucherService: s}
}

type createVoucherRequest struct {
	Brand     string    `json:"brand"`
	Code      string    `json:"code"`
	ValueXu   int       `json:"value_xu"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (h *AdminVoucherHandler) CreateVoucher(c *fiber.Ctx) error {
	var req createVoucherRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "dữ liệu không hợp lệ"})
	}

	if req.Brand == "" || req.Code == "" || req.ValueXu <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "dữ liệu không được để trống và giá trị xu phải lớn hơn 0"})
	}

	err := h.voucherService.CreateVoucher(c.Context(), services.CreateVoucherRequest{
		Brand:     req.Brand,
		Code:      req.Code,
		ValueXu:   req.ValueXu,
		ExpiresAt: req.ExpiresAt,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "thêm voucher thành công"})
}

func (h *AdminVoucherHandler) GetVouchers(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 50)
	offset := c.QueryInt("offset", 0)

	vouchers, err := h.voucherService.GetVouchers(c.Context(), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": vouchers})
}

func (h *AdminVoucherHandler) DeleteVoucher(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID không hợp lệ"})
	}

	err = h.voucherService.DeleteVoucher(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "xóa voucher thành công"})
}
