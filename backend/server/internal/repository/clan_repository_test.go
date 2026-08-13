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
	"gorm.io/gorm/logger"
)

func setupClanRepoTest(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, ClanRepository) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	assert.NoError(t, err)

	repo := NewClanRepository(gormDB)
	return gormDB, mock, repo
}

func TestClanRepository_CreateClan(t *testing.T) {
	_, mock, repo := setupClanRepoTest(t)
	ctx := context.Background()
	clan := &models.Clan{
		ID:       uuid.New(),
		Name:     "Test Clan",
		LeaderID: uuid.New(),
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "clans"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(clan.ID))
	mock.ExpectCommit()

	err := repo.CreateClan(ctx, clan)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestClanRepository_AddClanMember(t *testing.T) {
	_, mock, repo := setupClanRepoTest(t)
	ctx := context.Background()
	member := &models.ClanMember{
		ClanID: uuid.New(),
		UserID: uuid.New(),
		Role:   "leader",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "clan_members"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"joined_at"}).AddRow(time.Now()))
	mock.ExpectCommit()

	err := repo.AddClanMember(ctx, member)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestClanRepository_FindByName(t *testing.T) {
	_, mock, repo := setupClanRepoTest(t)
	ctx := context.Background()
	clanID := uuid.New()
	leaderID := uuid.New()

	mock.ExpectQuery(`SELECT \* FROM "clans" WHERE name = \$1 ORDER BY "clans"\."id" LIMIT \$2`).
		WithArgs("Test Clan", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "leader_id", "created_at"}).
			AddRow(clanID, "Test Clan", leaderID, time.Now()))

	clan, err := repo.FindByName(ctx, "Test Clan")
	assert.NoError(t, err)
	assert.NotNil(t, clan)
	assert.Equal(t, clanID, clan.ID)
	assert.Equal(t, "Test Clan", clan.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}
