// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package handlers

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type mockLocketService struct {
	err error
}

func (m *mockLocketService) SendLocket(ctx context.Context, matchID, senderID uuid.UUID, file *multipart.FileHeader) (string, error) {
	return "mock-url", m.err
}

func TestLocketHandler_SendLocket_Unauthorized(t *testing.T) {
	app := fiber.New()
	handler := NewLocketHandler(&mockLocketService{})

	app.Post("/send", handler.SendLocket)

	req := httptest.NewRequest(http.MethodPost, "/send", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", fiber.StatusUnauthorized, resp.StatusCode)
	}
}

func TestLocketHandler_SendLocket_InvalidUserID(t *testing.T) {
	app := fiber.New()
	handler := NewLocketHandler(&mockLocketService{})

	app.Post("/send", func(c *fiber.Ctx) error {
		c.Locals("user_id", "invalid-uuid")
		return handler.SendLocket(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/send", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
	}
}

func TestLocketHandler_SendLocket_InvalidMatchID(t *testing.T) {
	app := fiber.New()
	handler := NewLocketHandler(&mockLocketService{})

	app.Post("/send", func(c *fiber.Ctx) error {
		c.Locals("user_id", uuid.New().String())
		return handler.SendLocket(c)
	})

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("match_id", "invalid")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/send", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, _ := app.Test(req)

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
	}
}

func TestLocketHandler_SendLocket_MissingImage(t *testing.T) {
	app := fiber.New()
	handler := NewLocketHandler(&mockLocketService{})

	app.Post("/send", func(c *fiber.Ctx) error {
		c.Locals("user_id", uuid.New().String())
		return handler.SendLocket(c)
	})

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("match_id", uuid.New().String())
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/send", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, _ := app.Test(req)

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
	}
}

func TestLocketHandler_SendLocket_ServiceError(t *testing.T) {
	app := fiber.New()
	handler := NewLocketHandler(&mockLocketService{err: errors.New("service error")})

	app.Post("/send", func(c *fiber.Ctx) error {
		c.Locals("user_id", uuid.New().String())
		return handler.SendLocket(c)
	})

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("match_id", uuid.New().String())
	part, _ := writer.CreateFormFile("image", "test.jpg")
	part.Write([]byte("test image data"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/send", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, _ := app.Test(req)

	if resp.StatusCode != fiber.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", fiber.StatusInternalServerError, resp.StatusCode)
	}
}

func TestLocketHandler_SendLocket_Success(t *testing.T) {
	app := fiber.New()
	handler := NewLocketHandler(&mockLocketService{})

	app.Post("/send", func(c *fiber.Ctx) error {
		c.Locals("user_id", uuid.New().String())
		return handler.SendLocket(c)
	})

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("match_id", uuid.New().String())
	part, _ := writer.CreateFormFile("image", "test.jpg")
	part.Write([]byte("test image data"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/send", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, _ := app.Test(req)

	if resp.StatusCode != fiber.StatusAccepted {
		t.Errorf("Expected status %d, got %d", fiber.StatusAccepted, resp.StatusCode)
	}
}
