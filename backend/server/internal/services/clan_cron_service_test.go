// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"context"
	"errors"
	"testing"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// Mock ClanRepository
type mockCronClanRepo struct {
	topClan       *models.Clan
	topClanErr    error
	resetScoreErr error
}

func (m *mockCronClanRepo) GetTopWeeklyClan(ctx context.Context) (*models.Clan, error) {
	return m.topClan, m.topClanErr
}

func (m *mockCronClanRepo) ResetWeeklyScores(ctx context.Context) error {
	return m.resetScoreErr
}

func (m *mockCronClanRepo) CreateClan(ctx context.Context, clan *models.Clan) error {
	return nil
}

func (m *mockCronClanRepo) AddClanMember(ctx context.Context, member *models.ClanMember) error {
	return nil
}

func (m *mockCronClanRepo) FindByName(ctx context.Context, name string) (*models.Clan, error) {
	return nil, nil
}

// Mock LandmarkRepository
type mockCronLandmarkRepo struct {
	updateErr error
}

func (m *mockCronLandmarkRepo) UpdateAllOwners(ctx context.Context, ownerClan *models.Clan) error {
	return m.updateErr
}

// Mock NotificationRepository
type mockCronNotifRepo struct {
}

func (m *mockCronNotifRepo) Create(ctx context.Context, notification *models.Notification) error {
	return nil
}

// Mock PushService
type mockCronPushService struct {
	broadcastCalled bool
}

func (m *mockCronPushService) SendPush(ctx context.Context, userID string, title string, body string, payload map[string]string) error {
	return nil
}

func (m *mockCronPushService) BroadcastToAll(ctx context.Context, title string, body string, payload map[string]string) error {
	m.broadcastCalled = true
	return nil
}

func TestClanCronService_RunWeeklyReset(t *testing.T) {
	t.Run("success with top clan", func(t *testing.T) {
		topClan := &models.Clan{
			ID:          uuid.New(),
			Name:        "Top Clan",
			WeeklyScore: 1000,
		}

		mockClanRepo := &mockCronClanRepo{topClan: topClan}
		mockLandmarkRepo := &mockCronLandmarkRepo{}
		mockNotifRepo := &mockCronNotifRepo{}
		mockPushService := &mockCronPushService{}

		service := NewClanCronService(mockClanRepo, mockLandmarkRepo, mockNotifRepo, mockPushService)

		err := service.RunWeeklyReset(context.Background())
		assert.NoError(t, err)
		assert.True(t, mockPushService.broadcastCalled)
	})

	t.Run("success without top clan", func(t *testing.T) {
		mockClanRepo := &mockCronClanRepo{topClan: nil}
		mockLandmarkRepo := &mockCronLandmarkRepo{}
		mockNotifRepo := &mockCronNotifRepo{}
		mockPushService := &mockCronPushService{}

		service := NewClanCronService(mockClanRepo, mockLandmarkRepo, mockNotifRepo, mockPushService)

		err := service.RunWeeklyReset(context.Background())
		assert.NoError(t, err)
		assert.False(t, mockPushService.broadcastCalled)
	})

	t.Run("error getting top clan", func(t *testing.T) {
		mockClanRepo := &mockCronClanRepo{topClanErr: errors.New("db error")}
		mockLandmarkRepo := &mockCronLandmarkRepo{}
		mockNotifRepo := &mockCronNotifRepo{}
		mockPushService := &mockCronPushService{}

		service := NewClanCronService(mockClanRepo, mockLandmarkRepo, mockNotifRepo, mockPushService)

		err := service.RunWeeklyReset(context.Background())
		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
	})
}
