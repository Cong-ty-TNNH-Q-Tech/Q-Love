package handlers

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type mockLocketService struct {
	err error
}

func (m *mockLocketService) SendLocket(ctx context.Context, senderID uuid.UUID, matchID uuid.UUID, file *multipart.FileHeader) error {
	return m.err
}

func TestLocketHandler_SendLocket(t *testing.T) {
	app := fiber.New()
	service := &mockLocketService{}
	handler := NewLocketHandler(service)

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_id", uuid.New().String())
		return c.Next()
	})

	app.Post("/send", handler.SendLocket)

	// Create multipart form body
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("match_id", uuid.New().String())
	
	part, _ := writer.CreateFormFile("image", "test.jpg")
	part.Write([]byte("fake image data"))
	writer.Close()

	req := httptest.NewRequest("POST", "/send", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}

	if resp.StatusCode != fiber.StatusAccepted {
		t.Errorf("Expected status 202 Accepted, got %d", resp.StatusCode)
	}
}
