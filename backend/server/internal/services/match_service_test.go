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
)

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

	service := NewMatchService(mockRepo)

	err := service.Unmatch(context.Background(), matchID, userID)
	assert.NoError(t, err)
	assert.Empty(t, mockRepo.matches)
}

func TestMatchService_Unmatch_NotFound(t *testing.T) {
	matchID := uuid.New()
	userID := uuid.New()

	mockRepo := &mockMatchServiceRepo{
		matches: map[uuid.UUID]*models.Match{},
	}

	service := NewMatchService(mockRepo)

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

	service := NewMatchService(mockRepo)

	err := service.Unmatch(context.Background(), matchID, randomUserID)
	assert.EqualError(t, err, "forbidden")
}

func TestMatchService_Unmatch_DBError(t *testing.T) {
	matchID := uuid.New()
	userID := uuid.New()

	mockRepo := &mockMatchServiceRepo{
		err: errors.New("db error"),
	}

	service := NewMatchService(mockRepo)

	err := service.Unmatch(context.Background(), matchID, userID)
	assert.EqualError(t, err, "db error")
}

