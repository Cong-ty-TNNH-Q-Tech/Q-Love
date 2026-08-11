package main

import (
	"net/http"
	"testing"
	"time"
)

func TestMainServer(t *testing.T) {
	// Start the main function in a goroutine
	go main()

	// Give the server a moment to start
	time.Sleep(500 * time.Millisecond)

	// Test the health endpoint to ensure it's running
	resp, err := http.Get("http://localhost:3000/health")
	if err != nil {
		t.Fatalf("Failed to make request to health endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	// Gracefully shutdown the server so the test can finish
	if app != nil {
		if err := app.Shutdown(); err != nil {
			t.Errorf("Failed to shutdown server: %v", err)
		}
	}
}
