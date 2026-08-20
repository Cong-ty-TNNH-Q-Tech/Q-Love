// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"context"
	"net/http"
	"testing"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/alicebob/miniredis/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

type mockNotificationRepo struct {
	CreateFn       func(ctx context.Context, notif *models.Notification) error
	UpdateStatusFn func(ctx context.Context, id uuid.UUID, status string) error
}

func (m *mockNotificationRepo) Create(ctx context.Context, notif *models.Notification) error {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, notif)
	}
	return nil
}
func (m *mockNotificationRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	if m.UpdateStatusFn != nil {
		return m.UpdateStatusFn(ctx, id, status)
	}
	return nil
}

func TestNotificationService_SendPush_NoRedis(t *testing.T) {
	mockRepo := &mockNotificationRepo{}
	// Nil redis client to force error in getting token
	svc := NewNotificationService(mockRepo, nil, "fake-key")

	err := svc.SendPush(context.Background(), uuid.New(), "alert", "Hello", "World", nil)
	assert.Error(t, err)
}

func TestNotificationService_SendSilentPush_NoRedis(t *testing.T) {
	mockRepo := &mockNotificationRepo{}
	// Nil redis client to force error in getting token
	svc := NewNotificationService(mockRepo, nil, "fake-key")

	err := svc.SendSilentPush(context.Background(), uuid.New(), nil)
	assert.Error(t, err)
}

func TestNotificationService_SendPush_MockFCMKeyEmpty(t *testing.T) {
	// Setup miniredis
	mr, err := miniredis.Run()
	assert.NoError(t, err)
	defer mr.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	userID := uuid.New()
	key := "fcm_token:" + userID.String()
	redisClient.Set(context.Background(), key, "mock-token", 0)

	mockRepo := &mockNotificationRepo{
		CreateFn: func(ctx context.Context, notif *models.Notification) error {
			return nil
		},
		UpdateStatusFn: func(ctx context.Context, id uuid.UUID, status string) error {
			return nil
		},
	}

	svc := NewNotificationService(mockRepo, redisClient, "")
	err = svc.SendPush(context.Background(), userID, "title", "body", "type", nil)
	assert.NoError(t, err)
}

type mockTransport struct {
	roundTripFn func(req *http.Request) (*http.Response, error)
}

func (t *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.roundTripFn != nil {
		return t.roundTripFn(req)
	}
	return nil, nil
}

func TestNotificationService_SendPush_FCMError(t *testing.T) {
	mr, _ := miniredis.Run()
	defer mr.Close()

	redisClient := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	userID := uuid.New()
	redisClient.Set(context.Background(), "fcm_token:"+userID.String(), "token", 0)

	mockRepo := &mockNotificationRepo{}
	
	svc := NewNotificationService(mockRepo, redisClient, "real-key").(*notificationService)
	svc.httpClient = &http.Client{
		Transport: &mockTransport{
			roundTripFn: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 400,
					Body:       http.NoBody,
				}, nil
			},
		},
	}

	err := svc.SendPush(context.Background(), userID, "alert", "Hello", "World", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "FCM API error: status 400")
}

func TestNotificationService_SendSilentPush_FCMSuccess(t *testing.T) {
	mr, _ := miniredis.Run()
	defer mr.Close()

	redisClient := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	userID := uuid.New()
	redisClient.Set(context.Background(), "fcm_token:"+userID.String(), "token", 0)

	mockRepo := &mockNotificationRepo{}
	
	svc := NewNotificationService(mockRepo, redisClient, "real-key").(*notificationService)
	svc.httpClient = &http.Client{
		Transport: &mockTransport{
			roundTripFn: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       http.NoBody,
				}, nil
			},
		},
	}

	err := svc.SendSilentPush(context.Background(), userID, nil)
	assert.NoError(t, err)
}
