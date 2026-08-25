// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package api

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/config"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/pkg/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func TestRegisterRoutes(t *testing.T) {
	app := fiber.New()
	// setup dummy configs
	cfg := &config.Config{
		R2AccountID:       "dummy",
		R2AccessKeyID:     "dummy",
		R2SecretAccessKey: "dummy",
		R2BucketName:      "dummy",
		JWTSecret:         "dummy",
	}
	r2Client, _ := storage.NewR2Client(cfg)
	var db *gorm.DB
	var redisClient *redis.Client

	RegisterRoutes(context.Background(), app, db, r2Client, redisClient, cfg)

	// test wingmans route
	req := httptest.NewRequest("POST", "/api/v1/wingman/referrals", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("Failed to test route: %v", err)
	}
	if resp.StatusCode == 404 {
		t.Errorf("Expected route to be registered, got 404")
	}

	// Verify upload route was added
	var routeExists bool
	for _, route := range app.GetRoutes(true) {
		if route.Path == "/api/v1/upload/presigned-url" && route.Method == "POST" {
			routeExists = true
			break
		}
	}

	if !routeExists {
		t.Errorf("Expected route /api/v1/upload/presigned-url to be registered")
	}
}
