// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"context"
	"testing"
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/pkg/logger"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// Mocks

type MockPurgeQueueRepository struct {
	mock.Mock
}

func (m *MockPurgeQueueRepository) EnqueueUser(ctx context.Context, userID string, isVIP bool) error {
	args := m.Called(ctx, userID, isVIP)
	return args.Error(0)
}

func (m *MockPurgeQueueRepository) DequeueUsers(ctx context.Context, count int64) ([]string, error) {
	args := m.Called(ctx, count)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockPurgeQueueRepository) DequeueVIPUsers(ctx context.Context, count int64) ([]string, error) {
	args := m.Called(ctx, count)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockPurgeQueueRepository) RemoveUser(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockPurgeQueueRepository) ClearQueue(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

type MockMatchRepository struct {
	mock.Mock
}

func (m *MockMatchRepository) Create(ctx context.Context, match *models.Match) error {
	args := m.Called(ctx, match)
	return args.Error(0)
}
func (m *MockMatchRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Match, error) { return nil, nil }
func (m *MockMatchRepository) FindByIDUnscoped(ctx context.Context, id uuid.UUID) (*models.Match, error) { return nil, nil }
func (m *MockMatchRepository) FindByUsers(ctx context.Context, user1ID, user2ID uuid.UUID) (*models.Match, error) { return nil, nil }
func (m *MockMatchRepository) UpdateLastInteraction(ctx context.Context, id uuid.UUID, t time.Time) error { return nil }
func (m *MockMatchRepository) SoftDelete(ctx context.Context, id uuid.UUID) error { return nil }
func (m *MockMatchRepository) ResetStreakForInactiveMatches(ctx context.Context, d time.Duration) error { return nil }
func (m *MockMatchRepository) ResetIslandLevelForInactiveMatches(ctx context.Context, d time.Duration) error { return nil }

type MockPushService struct {
	mock.Mock
}

func (m *MockPushService) SendPush(ctx context.Context, userID string, title string, body string, payload map[string]string) error {
	args := m.Called(ctx, userID, title, body, payload)
	return args.Error(0)
}

func (m *MockPushService) BroadcastToAll(ctx context.Context, title string, body string, payload map[string]string) error {
	args := m.Called(ctx, title, body, payload)
	return args.Error(0)
}

type MockSpiritualService struct {
	mock.Mock
}

func (m *MockSpiritualService) CalculateZodiac(dob time.Time) string { return "" }
func (m *MockSpiritualService) CalculateNumerology(dob time.Time) int { return 0 }
func (m *MockSpiritualService) CalculateSpiritualMatchScore(dobA, dobB time.Time) int { return 0 }

// Tests

func TestPurgeService_ProcessMatchmaking(t *testing.T) {
	if logger.Log == nil {
		logger.Log = zap.NewNop()
	}

	queueRepo := new(MockPurgeQueueRepository)
	matchRepo := new(MockMatchRepository)
	pushSvc := new(MockPushService)
	spiritualSvc := new(MockSpiritualService)

	u1 := uuid.New().String()
	u2 := uuid.New().String()
	v1 := uuid.New().String()
	v2 := uuid.New().String()

	queueRepo.On("DequeueUsers", mock.Anything, int64(4)).Return([]string{u1, u2}, nil)
	queueRepo.On("DequeueVIPUsers", mock.Anything, int64(4)).Return([]string{v1, v2}, nil)

	matchRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Match")).Return(nil).Times(2)

	pushSvc.On("SendPush", mock.Anything, mock.AnythingOfType("string"), "The Purge", "Bạn đã được ghép đôi trong Đêm Săn Mồi!", mock.Anything).Return(nil).Times(4)

	svc := NewPurgeService(queueRepo, matchRepo, pushSvc, spiritualSvc)
	err := svc.ProcessMatchmaking(context.Background(), 4)

	assert.NoError(t, err)
	queueRepo.AssertExpectations(t)
	matchRepo.AssertExpectations(t)
	pushSvc.AssertExpectations(t)
}
