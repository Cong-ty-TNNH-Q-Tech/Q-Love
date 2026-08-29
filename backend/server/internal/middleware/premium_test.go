// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package middleware_test

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/middleware"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type mockPremiumRepo struct {
	premium *models.UserPremium
	err     error
}

func (m *mockPremiumRepo) Create(ctx context.Context, premium *models.UserPremium) error { return nil }
func (m *mockPremiumRepo) FindByUserID(ctx context.Context, userID uuid.UUID) (*models.UserPremium, error) {
	return m.premium, m.err
}
func (m *mockPremiumRepo) Update(ctx context.Context, premium *models.UserPremium) error { return nil }
func (m *mockPremiumRepo) IsUserPremium(ctx context.Context, userID uuid.UUID) (bool, error) {
	return m.premium != nil && m.premium.ExpiresAt.After(time.Now()), m.err
}
func (m *mockPremiumRepo) ActivatePremium(ctx context.Context, userID uuid.UUID, expiresAt time.Time) error {
	return nil
}

func TestPremiumMiddleware(t *testing.T) {
	app := fiber.New()
	userID := uuid.New()

	tests := []struct {
		name         string
		userID       interface{}
		repo         *mockPremiumRepo
		wantPremium  interface{}
	}{
		{
			name:   "No user_id in locals",
			userID: nil,
			repo:   &mockPremiumRepo{premium: nil, err: nil},
			wantPremium: nil,
		},
		{
			name:   "Active premium",
			userID: userID,
			repo: &mockPremiumRepo{
				premium: &models.UserPremium{
					ExpiresAt: time.Now().Add(24 * time.Hour),
				},
				err: nil,
			},
			wantPremium: true,
		},
		{
			name:   "Expired premium",
			userID: userID,
			repo: &mockPremiumRepo{
				premium: &models.UserPremium{
					ExpiresAt: time.Now().Add(-24 * time.Hour),
				},
				err: nil,
			},
			wantPremium: false,
		},
		{
			name:   "No premium found",
			userID: userID,
			repo:   &mockPremiumRepo{premium: nil, err: errors.New("not found")},
			wantPremium: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := middleware.PremiumMiddleware(tt.repo)

			app.Get("/test-"+tt.name, func(c *fiber.Ctx) error {
				if tt.userID != nil {
					c.Locals("user_id", tt.userID)
				}
				return handler(c)
			}, func(c *fiber.Ctx) error {
				val := c.Locals("is_premium")
				if val != tt.wantPremium {
					t.Errorf("expected is_premium %v, got %v", tt.wantPremium, val)
				}
				return c.SendStatus(fiber.StatusOK)
			})

			req := httptest.NewRequest("GET", "/test-"+tt.name, nil)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("failed to execute request: %v", err)
			}
			if resp.StatusCode != fiber.StatusOK {
				t.Errorf("expected status 200, got %d", resp.StatusCode)
			}
		})
	}
}
