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

func TestChatLockRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewChatLockRepository(gormDB)

	lock := &models.ChatLock{
		ID:        uuid.New(),
		UserID1:   uuid.New(),
		UserID2:   uuid.New(),
		ExpiresAt: time.Now().Add(time.Hour),
		CreatedAt: time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "chat_locks"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.Create(context.Background(), lock)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
