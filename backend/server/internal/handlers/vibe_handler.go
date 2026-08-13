// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/services"
)

type VibeHandler struct {
	SpotifyService *services.SpotifyService
}

func NewVibeHandler(spotifyService *services.SpotifyService) *VibeHandler {
	return &VibeHandler{
		SpotifyService: spotifyService,
	}
}

func (h *VibeHandler) Status(c *fiber.Ctx) error {
	unlocked := h.SpotifyService.CheckUnlockTime()
	
	return c.JSON(fiber.Map{
		"unlocked": unlocked,
		"message":  "Vibe Check is open from 23:00 to 05:00",
	})
}

func (h *VibeHandler) CurrentTrack(c *fiber.Ctx) error {
	if !h.SpotifyService.CheckUnlockTime() {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Vibe Check is currently locked"})
	}

	userID := "mock-user" // c.Locals("userID").(string)
	track, err := h.SpotifyService.GetCurrentTrack(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch Spotify track"})
	}

	return c.JSON(fiber.Map{
		"data": track,
	})
}

func (h *VibeHandler) Match(c *fiber.Ctx) error {
	if !h.SpotifyService.CheckUnlockTime() {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Vibe Check is currently locked"})
	}

	type MatchRequest struct {
		TrackID string `json:"track_id"`
	}

	var req MatchRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Mock match logic: we just create a mock room with a random user
	roomID := uuid.New().String()
	
	// Create VibeMatch model
	match := models.VibeMatch{
		ID:      uuid.New(),
		UserA:   uuid.New(), // mock current user
		UserB:   uuid.New(), // mock matched user
		TrackID: req.TrackID,
		RoomID:  roomID,
		Status:  "active",
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Matched successfully",
		"data": match,
	})
}
