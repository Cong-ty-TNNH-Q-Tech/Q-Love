// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package handlers

import (
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ExRatingHandler struct {
	exRatingService services.ExRatingService
}

func NewExRatingHandler(exRatingService services.ExRatingService) *ExRatingHandler {
	return &ExRatingHandler{exRatingService: exRatingService}
}

func (h *ExRatingHandler) SubmitRating(c *fiber.Ctx) error {
	type Request struct {
		TargetUserID string   `json:"target_user_id"`
		MatchID      string   `json:"match_id"`
		RatingScore  int      `json:"rating_score"`
		Tags         []string `json:"tags"`
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.RatingScore < 1 || req.RatingScore > 5 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Rating score must be between 1 and 5"})
	}

	targetUserID, err := uuid.Parse(req.TargetUserID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid target_user_id"})
	}

	matchID, err := uuid.Parse(req.MatchID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid match_id"})
	}

	err = h.exRatingService.SubmitRating(c.Context(), targetUserID, matchID, req.RatingScore, req.Tags)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Rating submitted anonymously",
	})
}

func (h *ExRatingHandler) ViewRating(c *fiber.Ctx) error {
	viewerIDStr, ok := c.Locals("user_id").(string)
	if !ok || viewerIDStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	viewerID, err := uuid.Parse(viewerIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user context"})
	}

	targetUserIDStr := c.Params("user_id")
	targetUserID, err := uuid.Parse(targetUserIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid target_user_id"})
	}

	avg, total, tags, err := h.exRatingService.ViewRating(c.Context(), viewerID, targetUserID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Format tags to API response structure
	type TagSummary struct {
		Tag   string `json:"tag"`
		Count int    `json:"count"`
	}
	var tagSummaries []TagSummary
	for t, count := range tags {
		tagSummaries = append(tagSummaries, TagSummary{Tag: t, Count: count})
	}
	if tagSummaries == nil {
		tagSummaries = make([]TagSummary, 0)
	}

	return c.JSON(fiber.Map{
		"target_user_id": targetUserIDStr,
		"average_rating": avg,
		"total_ratings":  total,
		"tag_summary":    tagSummaries,
	})
}
