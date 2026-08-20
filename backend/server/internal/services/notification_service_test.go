// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"context"
	"testing"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/google/uuid"
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
	// If redis is available but FCM key is empty, it should mock success
	// But we need a mini redis cluster to test it properly, so we will skip integration tests here
	// and rely on unit tests mocking.
}
