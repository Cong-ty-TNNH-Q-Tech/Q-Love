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
	mock.ExpectExec(`INSERT INTO "card_profiles"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
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

func TestCardRepository_GetCardProfileForUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewCardRepository(gormDB)
	userID := uuid.New()

	mock.ExpectQuery(`SELECT \* FROM "card_profiles" WHERE user_id = \$1 .* FOR UPDATE`).
		WithArgs(userID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(userID))

	profile, err := repo.GetCardProfileForUpdate(context.Background(), userID)
	assert.NoError(t, err)
	assert.NotNil(t, profile)
}

func TestCardRepository_UpdateCardProfile(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewCardRepository(gormDB)

	profile := &models.CardProfile{
		UserID: uuid.New(),
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "card_profiles"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.UpdateCardProfile(context.Background(), profile)
	assert.NoError(t, err)
}

func TestCardRepository_CreateCardTransaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewCardRepository(gormDB)

	tx := &models.CardTransaction{
		ID: uuid.New(),
	}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "card_transactions"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.CreateCardTransaction(context.Background(), tx)
	assert.NoError(t, err)
}

func TestCardRepository_GetOwnedQuantity(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewCardRepository(gormDB)
	collectorID := uuid.New()
	targetUserID := uuid.New()

	mock.ExpectQuery(`SELECT COALESCE\(SUM\(quantity\), 0\) FROM "card_transactions"`).
		WithArgs(collectorID, targetUserID, "buy").
		WillReturnRows(sqlmock.NewRows([]string{"COALESCE"}).AddRow(10))

	mock.ExpectQuery(`SELECT COALESCE\(SUM\(quantity\), 0\) FROM "card_transactions"`).
		WithArgs(collectorID, targetUserID, "sell").
		WillReturnRows(sqlmock.NewRows([]string{"COALESCE"}).AddRow(3))

	qty, err := repo.GetOwnedQuantity(context.Background(), collectorID, targetUserID)
	assert.NoError(t, err)
	assert.Equal(t, 7, qty)
}
