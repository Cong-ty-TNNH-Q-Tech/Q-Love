// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package repository

import (
	"context"
	"regexp"
	"testing"

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

func TestUserRepository_FindByPhone(t *testing.T) {
	repo, mock := setupUserRepoMock(t)
	phone := "0901234567"
	userID := uuid.New()

	rows := sqlmock.NewRows([]string{"id", "phone"}).
		AddRow(userID.String(), phone)

	mock.ExpectQuery("^SELECT \\* FROM \"users\" WHERE phone = \\$1").
		WithArgs(phone).
		WillReturnRows(rows)

	user, err := repo.FindByPhone(context.Background(), phone)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, phone, user.Phone)
}

func TestUserRepository_FindByPhone_NotFound(t *testing.T) {
	repo, mock := setupUserRepoMock(t)
	phone := "0901234567"

	mock.ExpectQuery("^SELECT \\* FROM \"users\" WHERE phone = \\$1").
		WithArgs(phone).
		WillReturnRows(sqlmock.NewRows([]string{"id", "phone"})) // Empty rows

	user, err := repo.FindByPhone(context.Background(), phone)
	assert.NoError(t, err)
	assert.Nil(t, user) // Should return nil if not found
}

func TestUserRepository_FindByID(t *testing.T) {
	repo, mock := setupUserRepoMock(t)
	userID := uuid.New()

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(userID.String())

	mock.ExpectQuery("^SELECT \\* FROM \"users\" WHERE id = \\$1").
		WithArgs(userID.String()).
		WillReturnRows(rows)

	user, err := repo.FindByID(context.Background(), userID)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, userID, user.ID)
}

func TestUserRepository_Create(t *testing.T) {
	repo, mock := setupUserRepoMock(t)
	userID := uuid.New()
	
	// Create doesn't use ExpectQuery, it uses ExpectExec unless RETURNING is used, but gorm often uses RETURNING id for postgres.
	// We'll just mock the BEGIN and COMMIT, and the INSERT
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO \"users\"").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userID.String()))
	mock.ExpectCommit()

	user := &models.User{
		ID:    userID,
		Phone: "0901234567",
	}

	err := repo.Create(context.Background(), user)
	assert.NoError(t, err)
}
