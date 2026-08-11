// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package handlers

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Presigner interface {
	GeneratePresignedURL(ctx context.Context, objectKey string, contentType string, contentLength int64) (string, error)
}

type UploadHandler struct {
	presigner Presigner
}

func NewUploadHandler(presigner Presigner) *UploadHandler {
	return &UploadHandler{
		presigner: presigner,
	}
}

type PresignedURLRequest struct {
	Filename      string `json:"filename"`
	ContentType   string `json:"contentType"`
	ContentLength int64  `json:"contentLength"`
}

func (h *UploadHandler) GenerateUploadURL(c *fiber.Ctx) error {
	var req PresignedURLRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	// 1. Validate File Size (< 10MB)
	if req.ContentLength <= 0 || req.ContentLength > 10*1024*1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "File size must be between 1 byte and 10MB"})
	}

	// 2. Validate Content Type
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/jpg":  true,
	}
	if !allowedTypes[req.ContentType] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Only JPG and PNG formats are allowed"})
	}

	// 3. Generate unique object key (uuid + extension)
	ext := filepath.Ext(req.Filename)
	if ext == "" {
		ext = ".jpg" // default fallback
	}
	// normalize extension
	ext = strings.ToLower(ext)
	objectKey := fmt.Sprintf("avatars/%s%s", uuid.New().String(), ext)

	// 4. Generate Presigned URL via Presigner
	url, err := h.presigner.GeneratePresignedURL(c.Context(), objectKey, req.ContentType, req.ContentLength)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate upload URL"})
	}

	// TODO: Placeholder for NSFW Scan triggering after upload completes
	// This might be handled via a webhook from Cloudflare or a subsequent API call from the client once upload is done.

	return c.JSON(fiber.Map{
		"uploadUrl": url,
		"objectKey": objectKey,
		"expiresIn": 900, // 15 mins
	})
}
