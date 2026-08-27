// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type DeviceHandler struct {
	redisClient *redis.Client
}

func NewDeviceHandler(redisClient *redis.Client) *DeviceHandler {
	return &DeviceHandler{redisClient: redisClient}
}

// RegisterFCMToken allows mobile app to register their FCM push token
func (h *DeviceHandler) RegisterFCMToken(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	type TokenRequest struct {
		Token string `json:"token"`
	}

	var req TokenRequest
	if err := c.BodyParser(&req); err != nil || req.Token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid token payload"})
	}

	// Save token in Redis
	key := fmt.Sprintf("fcm_token:%s", userID.String())
	// Set expiration to 60 days
	err := h.redisClient.Set(context.Background(), key, req.Token, 60*24*time.Hour).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save token"})
	}

	return c.JSON(fiber.Map{"message": "Token registered successfully"})
}
