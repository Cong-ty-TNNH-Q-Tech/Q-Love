// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package repository

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestCardRepository_CreateCardProfile(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewCardRepository(gormDB)

	profile := &models.CardProfile{
		UserID:    uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "card_profiles"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(profile.UserID))
	mock.ExpectCommit()

	err = repo.CreateCardProfile(context.Background(), profile)
	assert.NoError(t, err)
}

func TestCardRepository_GetCardProfile(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewCardRepository(gormDB)
	userID := uuid.New()

	mock.ExpectQuery(`SELECT \* FROM "card_profiles" WHERE user_id = \$1`).
		WithArgs(userID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(userID))

	profile, err := repo.GetCardProfile(context.Background(), userID)
	assert.NoError(t, err)
	assert.NotNil(t, profile)
}
