// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package middleware

import (
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// PremiumMiddleware checks if the user has an active premium subscription
// and sets "is_premium" in locals.
func PremiumMiddleware(premiumRepo repository.UserPremiumRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("user_id").(uuid.UUID)
		if !ok {
			// If no user_id, just proceed (or return unauthorized, but let other middlewares handle it)
			return c.Next()
		}

		premium, err := premiumRepo.FindByUserID(c.Context(), userID)
		if err == nil && premium != nil && premium.ExpiresAt.After(time.Now()) {
			c.Locals("is_premium", true)
		} else {
			c.Locals("is_premium", false)
		}

		return c.Next()
	}
}
