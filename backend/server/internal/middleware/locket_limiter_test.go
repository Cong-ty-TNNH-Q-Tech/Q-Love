package middleware

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type mockUserPremiumRepo struct {
	isPremium bool
	err       error
}

func (m *mockUserPremiumRepo) IsUserPremium(ctx context.Context, userID uuid.UUID) (bool, error) {
	return m.isPremium, m.err
}

func TestLocketRateLimiter_Bypass(t *testing.T) {
	app := fiber.New()
	
	repo := &mockUserPremiumRepo{isPremium: true}
	
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_id", uuid.New().String())
		return c.Next()
	})
	
	app.Post("/send", LocketRateLimiter(repo), func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusAccepted)
	})

	// Test 1: Premium user should bypass rate limit (can send more than 10)
	for i := 0; i < 15; i++ {
		req := httptest.NewRequest("POST", "/send", nil)
		resp, _ := app.Test(req)
		if resp.StatusCode != fiber.StatusAccepted {
			t.Errorf("Expected status 202 Accepted, got %d on request %d", resp.StatusCode, i)
		}
	}
}

func TestLocketRateLimiter_Limit(t *testing.T) {
	app := fiber.New()
	
	repo := &mockUserPremiumRepo{isPremium: false} // Not premium
	userID := uuid.New().String()
	
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_id", userID)
		return c.Next()
	})
	
	app.Post("/send", LocketRateLimiter(repo), func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusAccepted)
	})

	// Test 2: Normal user should hit limit after 10
	for i := 1; i <= 15; i++ {
		req := httptest.NewRequest("POST", "/send", nil)
		resp, _ := app.Test(req)
		
		if i <= 10 {
			if resp.StatusCode != fiber.StatusAccepted {
				t.Errorf("Expected status 202 Accepted, got %d on request %d", resp.StatusCode, i)
			}
		} else {
			if resp.StatusCode != fiber.StatusTooManyRequests {
				t.Errorf("Expected status 429 Too Many Requests, got %d on request %d", resp.StatusCode, i)
			}
		}
	}
}
