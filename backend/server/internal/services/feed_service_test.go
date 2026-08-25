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

type mockUserRepository struct {
	user  models.User
	users []models.User
	feed  []models.User
	err   error
}

func (m *mockUserRepository) GetTopUsersByScore(ctx context.Context, limit int) ([]uuid.UUID, error) {
	return nil, nil
}

func (m *mockUserRepository) FindByPhone(ctx context.Context, phone string) (*models.User, error) {
	return nil, nil
}

func (m *mockUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &m.user, nil
}

func (m *mockUserRepository) GetFeed(ctx context.Context, userID uuid.UUID, radius int) ([]models.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.feed, nil
}

func (m *mockUserRepository) GetSpiritualFeed(ctx context.Context, userID uuid.UUID, radius int) ([]models.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.feed, nil
}

// Other mock methods returning nil
func (m *mockUserRepository) UpdateLocation(ctx context.Context, userID uuid.UUID, lat, lon float64) error { return nil }
func (m *mockUserRepository) Create(ctx context.Context, user *models.User) error { return nil }

func TestFeedService_GetFeed_Default(t *testing.T) {
	dob1 := time.Date(1995, 5, 5, 0, 0, 0, 0, time.UTC)
	dob2 := time.Date(1996, 6, 6, 0, 0, 0, 0, time.UTC)
	repo := &mockUserRepository{
		user: models.User{ID: uuid.New(), DOB: &dob1},
		feed: []models.User{
			{ID: uuid.New(), DOB: &dob2},
		},
	}
	spiritual := NewSpiritualService()
	svc := NewFeedService(repo, spiritual)

	res, err := svc.GetFeed(context.Background(), uuid.New(), "default", 50)
	assert.NoError(t, err)
	assert.Len(t, res, 1)
}

func TestFeedService_GetFeed_Spiritual(t *testing.T) {
	dob := time.Date(1995, 5, 5, 0, 0, 0, 0, time.UTC)
	repo := &mockUserRepository{
		user: models.User{ID: uuid.New(), DOB: &dob},
		feed: []models.User{
			{ID: uuid.New(), DOB: &dob},
			{ID: uuid.New(), DOB: &dob},
		},
	}
	spiritual := NewSpiritualService()
	svc := NewFeedService(repo, spiritual)

	res, err := svc.GetFeed(context.Background(), uuid.New(), "spiritual", 50)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
}

func TestFeedService_GetFeed_UserErr(t *testing.T) {
	repo := &mockUserRepository{
		err: errors.New("db error"),
	}
	spiritual := NewSpiritualService()
	svc := NewFeedService(repo, spiritual)

	res, err := svc.GetFeed(context.Background(), uuid.New(), "default", 50)
	assert.Error(t, err)
	assert.Nil(t, res)
}
