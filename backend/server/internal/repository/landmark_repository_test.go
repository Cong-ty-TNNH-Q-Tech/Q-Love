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

func TestLandmarkRepository_UpdateAllOwners(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewLandmarkRepository(gormDB)

	t.Run("success with top clan", func(t *testing.T) {
		topClan := &models.Clan{
			ID: uuid.New(),
		}
		
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "landmarks" SET "current_owner_clan_id"=$1 WHERE 1 = 1`)).
			WithArgs(topClan.ID.String()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.UpdateAllOwners(context.Background(), topClan)
		assert.NoError(t, err)
	})

	t.Run("success with nil clan", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "landmarks" SET "current_owner_clan_id"=$1 WHERE 1 = 1`)).
			WithArgs(nil).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.UpdateAllOwners(context.Background(), nil)
		assert.NoError(t, err)
	})
}
