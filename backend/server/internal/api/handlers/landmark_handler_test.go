// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockLandmarkService struct {
	mock.Mock
}

func (m *mockLandmarkService) CheckIn(ctx context.Context, userID, landmarkID uuid.UUID, lat, lng float64, isMocked bool) (*models.Landmark, error) {
	args := m.Called(ctx, userID, landmarkID, lat, lng, isMocked)
	var v *models.Landmark
	if args.Get(0) != nil {
		v = args.Get(0).(*models.Landmark)
	}
	return v, args.Error(1)
}

func TestLandmarkHandler_CheckIn(t *testing.T) {
	app := fiber.New()
	mockService := new(mockLandmarkService)
	handler := NewLandmarkHandler(mockService)

	app.Post("/checkin/:landmark_id", func(c *fiber.Ctx) error {
		c.Locals("user_id", uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"))
		return handler.CheckIn(c)
	})

	userID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	landmarkID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockService.On("CheckIn", mock.Anything, userID, landmarkID, 10.0, 106.0, false).Return(nil, nil).Once()

		reqBody := map[string]interface{}{
			"latitude":  10.0,
			"longitude": 106.0,
			"is_mocked": false,
		}
		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/checkin/"+landmarkID.String(), bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Fake GPS Error", func(t *testing.T) {
		mockService.On("CheckIn", mock.Anything, userID, landmarkID, 10.0, 106.0, true).Return(nil, services.ErrFakeGPS).Once()

		reqBody := map[string]interface{}{
			"latitude":  10.0,
			"longitude": 106.0,
			"is_mocked": true,
		}
		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/checkin/"+landmarkID.String(), bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})
	
	t.Run("Out of Range Error", func(t *testing.T) {
		mockService.On("CheckIn", mock.Anything, userID, landmarkID, 10.0, 106.0, false).Return(nil, services.ErrOutOfRange).Once()

		reqBody := map[string]interface{}{
			"latitude":  10.0,
			"longitude": 106.0,
			"is_mocked": false,
		}
		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/checkin/"+landmarkID.String(), bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
