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

func TestLandmarkRepository_FindByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewLandmarkRepository(gormDB)
	landmarkID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "landmarks" WHERE id = \$1 AND "landmarks"\."deleted_at" IS NULL ORDER BY "landmarks"\."id" LIMIT \$2`).
			WithArgs(landmarkID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(landmarkID))

		res, err := repo.FindByID(context.Background(), landmarkID)
		assert.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "landmarks" WHERE .*`).
			WithArgs(landmarkID, 1).
			WillReturnError(assert.AnError)

		res, err := repo.FindByID(context.Background(), landmarkID)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestLandmarkRepository_CheckDistance(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewLandmarkRepository(gormDB)
	landmarkID := uuid.New()
	lat := 10.0
	lng := 106.0

	t.Run("success true", func(t *testing.T) {
		mock.ExpectQuery(`SELECT ST_DWithin\(.*`).
			WithArgs(lng, lat, landmarkID).
			WillReturnRows(sqlmock.NewRows([]string{"is_within"}).AddRow(true))

		res, err := repo.CheckDistance(context.Background(), landmarkID, lat, lng)
		assert.NoError(t, err)
		assert.True(t, res)
	})

	t.Run("success false", func(t *testing.T) {
		mock.ExpectQuery(`SELECT ST_DWithin\(.*`).
			WithArgs(lng, lat, landmarkID).
			WillReturnRows(sqlmock.NewRows([]string{"is_within"}).AddRow(false))

		res, err := repo.CheckDistance(context.Background(), landmarkID, lat, lng)
		assert.NoError(t, err)
		assert.False(t, res)
	})

	t.Run("error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT ST_DWithin\(.*`).
			WithArgs(lng, lat, landmarkID).
			WillReturnError(assert.AnError)

		res, err := repo.CheckDistance(context.Background(), landmarkID, lat, lng)
		assert.Error(t, err)
		assert.False(t, res)
	})
}
