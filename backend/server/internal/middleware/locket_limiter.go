package middleware

import (
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/google/uuid"
)

// LocketRateLimiter creates a middleware to limit Locket sends to 10 per hour.
// It bypasses the limit if the user has an active Premium subscription.
func LocketRateLimiter(userPremiumRepo repository.UserPremiumRepository) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        10,
		Expiration: 1 * time.Hour,
		KeyGenerator: func(c *fiber.Ctx) string {
			// Rate limit based on user ID and match ID.
			// Assuming User ID is in locals (from JWT middleware) and Match ID is in the form or body
			userID, _ := c.Locals("user_id").(string)
			matchID := c.FormValue("match_id")
			if matchID == "" {
				matchID = "unknown_match" // Fallback, though validation should catch it
			}
			return "locket_" + userID + "_" + matchID
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too Many Requests",
				"message": "Bạn đã gửi quá nhiều Locket trong giờ này. Hãy thư giãn và thử lại sau.",
			})
		},
		Next: func(c *fiber.Ctx) bool {
			userIDStr, ok := c.Locals("user_id").(string)
			if !ok {
				return false
			}
			userID, err := uuid.Parse(userIDStr)
			if err != nil {
				return false
			}

			// Check if user is premium to bypass the limit
			isPremium, err := userPremiumRepo.IsUserPremium(c.Context(), userID)
			if err != nil {
				// On error, we don't bypass to be safe, but ideally log it
				return false
			}
			return isPremium
		},
	})
}
