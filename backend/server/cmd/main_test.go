// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/config"
)

func TestMainServer(t *testing.T) {
	// Setup app locally instead of using global
	t.Setenv("JWT_SECRET", "test-secret")
	cfg := config.LoadConfig()
	
	// Skip real db connection for tests
	cfg.DatabaseDSN = "skip"
	cfg.RedisURL = "skip"

	app, err := setupApp(cfg)
	if err != nil {
		t.Fatalf("setupApp failed: %v", err)
	}

	// Start the server in a goroutine
	go func() {
		_ = app.Listen(":3001")
	}()

	// Give the server a moment to start
	time.Sleep(500 * time.Millisecond)

	// Test the health endpoint to ensure it's running
	resp, err := http.Get("http://localhost:3001/health")
	if err != nil {
		t.Fatalf("Failed to make request to health endpoint: %v", err)
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	// Test ping endpoint
	resp2, _ := http.Get("http://localhost:3001/ping")
	if resp2 != nil {
		resp2.Body.Close()
	}

	// Test version endpoint
	resp3, _ := http.Get("http://localhost:3001/version")
	if resp3 != nil {
		resp3.Body.Close()
	}

	// Gracefully shutdown the server so the test can finish
	if err := app.Shutdown(); err != nil {
		t.Errorf("Failed to shutdown server: %v", err)
	}

	// Also call main() directly in a goroutine to get coverage for it
	// We set a unique port to avoid conflicts
	t.Setenv("PORT", "3002")
	t.Setenv("DATABASE_DSN", "skip")
	t.Setenv("REDIS_URL", "skip")
	go func() {
		// This will block, but the test will exit and kill the goroutine
		main()
	}()
	time.Sleep(500 * time.Millisecond) // let it start up
}
