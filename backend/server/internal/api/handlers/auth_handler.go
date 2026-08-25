// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package handlers

import (
	"strings"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/services"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type SendOTPRequest struct {
	Phone string `json:"phone"`
}

func (h *AuthHandler) SendOTP(c *fiber.Ctx) error {
	var req SendOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Phone == "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": "Phone number is required",
		})
	}

	err := h.authService.SendOTP(c.Context(), req.Phone)
	if err != nil {
		if strings.Contains(err.Error(), "rate limit exceeded") {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send OTP",
		})
	}

	return c.JSON(fiber.Map{
		"message": "OTP sent",
		"retry_after_seconds": 60,
	})
}

type VerifyOTPRequest struct {
	Phone string `json:"phone"`
	OTP   string `json:"otp"`
}

func (h *AuthHandler) VerifyOTP(c *fiber.Ctx) error {
	var req VerifyOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Phone == "" || req.OTP == "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": "Phone and OTP are required",
		})
	}

	user, accessToken, refreshToken, isNewUser, err := h.authService.VerifyOTP(c.Context(), req.Phone, req.OTP)
	if err != nil {
		if err.Error() == "ERR_INVALID_OTP" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"code": 400,
				"message": "ERR_INVALID_OTP",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to verify OTP",
		})
	}

	return c.JSON(fiber.Map{
		"access_token": accessToken,
		"refresh_token": refreshToken,
		"user": user,
		"is_new_user": isNewUser,
	})
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.RefreshToken == "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": "Refresh token is required",
		})
	}

	newAccessToken, newRefreshToken, err := h.authService.RefreshToken(c.Context(), req.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"access_token": newAccessToken,
		"refresh_token": newRefreshToken,
	})
}
