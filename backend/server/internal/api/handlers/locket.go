package handlers

import (
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type LocketHandler struct {
	locketService services.LocketService
}

func NewLocketHandler(locketService services.LocketService) *LocketHandler {
	return &LocketHandler{locketService: locketService}
}

// SendLocket handles POST /api/v1/locket/send
func (h *LocketHandler) SendLocket(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	senderID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid User ID"})
	}

	matchIDStr := c.FormValue("match_id")
	matchID, err := uuid.Parse(matchIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Match ID"})
	}

	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Image file is required"})
	}

	if err := h.locketService.SendLocket(c.Context(), senderID, matchID, file); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "Yêu cầu gửi Locket đã được tiếp nhận",
	})
}
