package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
)

func setupUserRepoMock(t *testing.T) (UserRepository, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewUserRepository(gormDB)
	return repo, mock
}

func TestUserRepository_GetTopUsersByScore(t *testing.T) {
	repo, mock := setupUserRepoMock(t)

	mockID1 := uuid.New()
	mockID2 := uuid.New()
	mockRows := sqlmock.NewRows([]string{"user_id"}).
		AddRow(mockID1).
		AddRow(mockID2)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "user_id" FROM "card_profiles" ORDER BY (match_count_cached + locket_count_cached + clan_upvote_cached - court_penalty_cached) DESC LIMIT $1`)).
		WithArgs(5).
		WillReturnRows(mockRows)

	users, err := repo.GetTopUsersByScore(context.Background(), 5)
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, mockID1, users[0])
	assert.Equal(t, mockID2, users[1])
}

func TestUserRepository_GetTopUsersByScore_Error(t *testing.T) {
	repo, mock := setupUserRepoMock(t)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "user_id" FROM "card_profiles" ORDER BY (match_count_cached + locket_count_cached + clan_upvote_cached - court_penalty_cached) DESC LIMIT $1`)).
		WithArgs(5).
		WillReturnError(assert.AnError)

	users, err := repo.GetTopUsersByScore(context.Background(), 5)
	assert.Error(t, err)
	assert.Nil(t, users)
}
