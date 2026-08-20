// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.
package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/google/uuid"
)

// LocketRateLimiter creates a middleware to limit Locket sends to 10 per hour.
// It bypasses the limit if the user has an active Premium subscription.
func LocketRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        10,
		Expiration: 1 * time.Hour,
		KeyGenerator: func(c *fiber.Ctx) string {
			// Rate limit based on user ID and match ID.
			// Assuming User ID is in locals (from JWT middleware) and Match ID is in the form or body
			userID, ok := c.Locals("user_id").(uuid.UUID)
			if !ok {
				return ""
			}
			matchID := c.FormValue("match_id")
			if matchID == "" {
				matchID = "unknown_match" // Fallback, though validation should catch it
			}
			return "locket_" + userID.String() + "_" + matchID
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too Many Requests",
				"message": "Bạn đã gửi quá nhiều Locket trong giờ này. Hãy thư giãn và thử lại sau.",
			})
		},
		Next: func(c *fiber.Ctx) bool {
			isPremium, _ := c.Locals("is_premium").(bool)
			return isPremium
		},
	})
}

