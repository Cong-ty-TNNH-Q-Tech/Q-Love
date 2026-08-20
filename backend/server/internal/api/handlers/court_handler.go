// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package handlers

import (
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CourtHandler struct {
	courtService services.CourtService
}

func NewCourtHandler(courtService services.CourtService) *CourtHandler {
	return &CourtHandler{
		courtService: courtService,
	}
}

type FileLawsuitRequest struct {
	DefendantID uuid.UUID `json:"defendant_id" validate:"required"`
	MatchID     uuid.UUID `json:"match_id" validate:"required"`
	Reason      string    `json:"reason" validate:"required"`
}

type VoteRequest struct {
	Vote string `json:"vote" validate:"required"`
}

func (h *CourtHandler) FileLawsuit(c *fiber.Ctx) error {
	userIDStr := c.Locals("userID").(string)
	plaintiffID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid user id"})
	}

	var req FileLawsuitRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	courtCase, err := h.courtService.FileLawsuit(c.Context(), plaintiffID, req.DefendantID, req.MatchID, req.Reason)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "lawsuit filed successfully",
		"data":    courtCase,
	})
}

func (h *CourtHandler) GetFeed(c *fiber.Ctx) error {
	userIDStr := c.Locals("userID").(string)
	jurorID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid user id"})
	}

	limit := c.QueryInt("limit", 20)
	cases, err := h.courtService.GetFeed(c.Context(), jurorID, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch court feed"})
	}

	return c.JSON(fiber.Map{
		"data": cases,
	})
}

func (h *CourtHandler) VoteCase(c *fiber.Ctx) error {
	userIDStr := c.Locals("userID").(string)
	jurorID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid user id"})
	}

	caseIDStr := c.Params("case_id")
	caseID, err := uuid.Parse(caseIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid case id"})
	}

	var req VoteRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	voteType := models.CourtVoteType(req.Vote)
	if voteType != models.CourtVoteGuilty && voteType != models.CourtVoteNotGuilty {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid vote type"})
	}

	if err := h.courtService.VoteCase(c.Context(), caseID, jurorID, voteType); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "vote submitted successfully",
	})
}

func (h *CourtHandler) WithdrawCase(c *fiber.Ctx) error {
	userIDStr := c.Locals("userID").(string)
	plaintiffID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid user id"})
	}

	caseIDStr := c.Params("case_id")
	caseID, err := uuid.Parse(caseIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid case id"})
	}

	if err := h.courtService.WithdrawCase(c.Context(), caseID, plaintiffID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "case withdrawn successfully",
	})
}
