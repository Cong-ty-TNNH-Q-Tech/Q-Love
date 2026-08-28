// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.


package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"fmt"
)

type mockNotificationService struct{}

func (m *mockNotificationService) SendPush(ctx context.Context, userID uuid.UUID, nType, title, body string, data map[string]string) error {
	return nil
}

func (m *mockNotificationService) SendSilentPush(ctx context.Context, userID uuid.UUID, data map[string]string) error {
	return nil
}

type mockMatchServiceRepo struct {
	matches map[uuid.UUID]*models.Match
	err     error
}

func (m *mockMatchServiceRepo) FindByID(ctx context.Context, id uuid.UUID) (*models.Match, error) {
	if m.err != nil {
		return nil, m.err
	}
	match, ok := m.matches[id]
	if !ok {
		return nil, nil // Return nil, nil for not found (or however FindByID behaves in GORM)
	}
	return match, nil
}

func (m *mockMatchServiceRepo) FindByIDUnscoped(ctx context.Context, id uuid.UUID) (*models.Match, error) {
	if m.err != nil {
		return nil, m.err
	}
	match, ok := m.matches[id]
	if !ok {
		return nil, nil
	}
	return match, nil
}

func (m *mockMatchServiceRepo) UpdateLastInteraction(ctx context.Context, id uuid.UUID, t time.Time) error {
	return nil
}

func (m *mockMatchServiceRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	if m.err != nil {
		return m.err
	}
	delete(m.matches, id)
	return nil
}

func (m *mockMatchServiceRepo) Create(ctx context.Context, match *models.Match) error {
	return nil
}
func (m *mockMatchServiceRepo) FindByUsers(ctx context.Context, u1, u2 uuid.UUID) (*models.Match, error) {
	return nil, nil
}
func (m *mockMatchServiceRepo) ResetStreakForInactiveMatches(ctx context.Context, inactiveDuration time.Duration) error { return nil }
func (m *mockMatchServiceRepo) ResetIslandLevelForInactiveMatches(ctx context.Context, inactiveDuration time.Duration) error { return nil }

func TestMatchService_Unmatch(t *testing.T) {
	matchID := uuid.New()
	userID := uuid.New()
	otherUserID := uuid.New()

	mockRepo := &mockMatchServiceRepo{
		matches: map[uuid.UUID]*models.Match{
			matchID: {
				ID:      matchID,
				User1ID: userID,
				User2ID: otherUserID,
			},
		},
	}

	mr, err := miniredis.Run()
	assert.NoError(t, err)
	defer mr.Close()
	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	service := NewMatchService(mockRepo, new(mockNotificationService), rdb)

	err = service.Unmatch(context.Background(), matchID, userID)
	assert.NoError(t, err)
	assert.Empty(t, mockRepo.matches)

	// Verify Redis state
	userKey := fmt.Sprintf("pending_ex_ratings:%s", userID.String())
	partnerKey := fmt.Sprintf("pending_ex_ratings:%s", otherUserID.String())

	isMemberUser, _ := rdb.SIsMember(context.Background(), userKey, otherUserID.String()).Result()
	assert.True(t, isMemberUser)

	isMemberPartner, _ := rdb.SIsMember(context.Background(), partnerKey, userID.String()).Result()
	assert.True(t, isMemberPartner)

	// Verify expiration is set
	ttl, _ := rdb.TTL(context.Background(), userKey).Result()
	assert.True(t, ttl > 0)
}

func TestMatchService_Unmatch_NotFound(t *testing.T) {
	matchID := uuid.New()
	userID := uuid.New()

	mockRepo := &mockMatchServiceRepo{
		matches: map[uuid.UUID]*models.Match{},
	}

	service := NewMatchService(mockRepo, new(mockNotificationService), nil)

	err := service.Unmatch(context.Background(), matchID, userID)
	assert.EqualError(t, err, "match not found")
}

func TestMatchService_Unmatch_Forbidden(t *testing.T) {
	matchID := uuid.New()
	userID := uuid.New()
	otherUserID := uuid.New()
	randomUserID := uuid.New()

	mockRepo := &mockMatchServiceRepo{
		matches: map[uuid.UUID]*models.Match{
			matchID: {
				ID:      matchID,
				User1ID: userID,
				User2ID: otherUserID,
			},
		},
	}

	service := NewMatchService(mockRepo, new(mockNotificationService), nil)

	err := service.Unmatch(context.Background(), matchID, randomUserID)
	assert.EqualError(t, err, "forbidden")
}

func TestMatchService_Unmatch_DBError(t *testing.T) {
	matchID := uuid.New()
	userID := uuid.New()

	mockRepo := &mockMatchServiceRepo{
		err: errors.New("db error"),
	}

	service := NewMatchService(mockRepo, new(mockNotificationService), nil)

	err := service.Unmatch(context.Background(), matchID, userID)
	assert.EqualError(t, err, "db error")
}

