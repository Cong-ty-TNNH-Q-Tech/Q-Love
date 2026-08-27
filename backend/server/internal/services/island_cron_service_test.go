// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/stretchr/testify/assert"
)

type mockMatchRepoForIsland struct {
	resetStreakErr error
	resetIslandErr error
}

func (m *mockMatchRepoForIsland) Create(ctx context.Context, match *models.Match) error { return nil }
func (m *mockMatchRepoForIsland) FindByID(ctx context.Context, id uuid.UUID) (*models.Match, error) { return nil, nil }
func (m *mockMatchRepoForIsland) FindByIDUnscoped(ctx context.Context, id uuid.UUID) (*models.Match, error) { return nil, nil }
func (m *mockMatchRepoForIsland) UpdateLastInteraction(ctx context.Context, id uuid.UUID, t time.Time) error { return nil }
func (m *mockMatchRepoForIsland) SoftDelete(ctx context.Context, id uuid.UUID) error { return nil }
func (m *mockMatchRepoForIsland) FindByUsers(ctx context.Context, u1, u2 uuid.UUID) (*models.Match, error) { return nil, nil }

func (m *mockMatchRepoForIsland) ResetStreakForInactiveMatches(ctx context.Context, inactiveDuration time.Duration) error {
	return m.resetStreakErr
}

func (m *mockMatchRepoForIsland) ResetIslandLevelForInactiveMatches(ctx context.Context, inactiveDuration time.Duration) error {
	return m.resetIslandErr
}

func TestIslandCronService_RunDailyGhostingCheck(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &mockMatchRepoForIsland{}
		service := NewIslandCronService(repo)

		err := service.RunDailyGhostingCheck(ctx)
		assert.NoError(t, err)
	})

	t.Run("reset streak error", func(t *testing.T) {
		repo := &mockMatchRepoForIsland{
			resetStreakErr: errors.New("db error"),
		}
		service := NewIslandCronService(repo)

		err := service.RunDailyGhostingCheck(ctx)
		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
	})

	t.Run("reset island level error", func(t *testing.T) {
		repo := &mockMatchRepoForIsland{
			resetIslandErr: errors.New("db error 2"),
		}
		service := NewIslandCronService(repo)

		err := service.RunDailyGhostingCheck(ctx)
		assert.Error(t, err)
		assert.Equal(t, "db error 2", err.Error())
	})
}
