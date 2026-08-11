package api

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func TestRegisterRoutes(t *testing.T) {
	app := fiber.New()
	
	// pass a dummy or nil db for now
	var db *gorm.DB
	RegisterRoutes(app, db)

	req := httptest.NewRequest("POST", "/api/v1/wingmans/referral", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("Failed to test route: %v", err)
	}

	// Because of invalid json, it returns 400 instead of 404
	if resp.StatusCode == 404 {
		t.Errorf("Expected route to be registered, got 404")
	}
}
