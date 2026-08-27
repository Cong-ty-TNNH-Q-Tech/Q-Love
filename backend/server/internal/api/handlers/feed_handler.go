// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package handlers

import (
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type FeedHandler struct {
	feedService services.FeedService
}

func NewFeedHandler(feedService services.FeedService) *FeedHandler {
	return &FeedHandler{
		feedService: feedService,
	}
}

func (h *FeedHandler) GetFeed(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	filter := c.Query("filter", "default")
	radius := c.QueryInt("radius", 50000) // Default 50km

	feed, err := h.feedService.GetFeed(c.Context(), userID, filter, radius)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get feed"})
	}

	return c.JSON(fiber.Map{
		"message": "success",
		"data":    feed,
	})
}
