// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupExRatingMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	return db, mock
}

func TestExRatingRepository_Create(t *testing.T) {
	db, mock := setupExRatingMockDB(t)
	repo := NewExRatingRepository(db)

	rating := &models.ExRating{
		ID:           uuid.New(),
		TargetUserID: uuid.New(),
		MatchID:      uuid.New(),
		RatingScore:  5,
		TagsString:   `#tốt`,
		CreatedAt:    time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "ex_ratings"`)).
		WithArgs(rating.ID, rating.TargetUserID, rating.MatchID, rating.RatingScore, rating.TagsString, rating.CreatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Create(context.Background(), rating)
	assert.NoError(t, err)
}

func TestExRatingRepository_HasRated(t *testing.T) {
	db, mock := setupExRatingMockDB(t)
	repo := NewExRatingRepository(db)

	targetUserID := uuid.New()
	matchID := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "ex_ratings" WHERE match_id = $1 AND target_user_id = $2`)).
		WithArgs(matchID, targetUserID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	hasRated, err := repo.HasRated(context.Background(), matchID, targetUserID)
	assert.NoError(t, err)
	assert.True(t, hasRated)
}

func TestExRatingRepository_GetSummaryByUserID(t *testing.T) {
	db, mock := setupExRatingMockDB(t)
	repo := NewExRatingRepository(db)

	targetUserID := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ex_ratings" WHERE target_user_id = $1`)).
		WithArgs(targetUserID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "target_user_id", "match_id", "rating_score", "tags_string"}).
			AddRow(uuid.New(), targetUserID, uuid.New(), 5, `#good,#funny`).
			AddRow(uuid.New(), targetUserID, uuid.New(), 4, `#good`))

	avg, total, tags, err := repo.GetSummaryByUserID(context.Background(), targetUserID)
	assert.NoError(t, err)
	assert.Equal(t, float64(4.5), avg)
	assert.Equal(t, int64(2), total)
	assert.Equal(t, 2, tags["#good"])
	assert.Equal(t, 1, tags["#funny"])
}

func TestExRatingRepository_GetSummaryByUserID_Empty(t *testing.T) {
	db, mock := setupExRatingMockDB(t)
	repo := NewExRatingRepository(db)

	targetUserID := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ex_ratings" WHERE target_user_id = $1`)).
		WithArgs(targetUserID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "target_user_id", "match_id", "rating_score", "tags_string"}))

	avg, total, tags, err := repo.GetSummaryByUserID(context.Background(), targetUserID)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), avg)
	assert.Equal(t, int64(0), total)
	assert.Empty(t, tags)
}

func TestExRatingRepository_GetSummaryByUserID_Error(t *testing.T) {
	db, mock := setupExRatingMockDB(t)
	repo := NewExRatingRepository(db)

	targetUserID := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ex_ratings" WHERE target_user_id = $1`)).
		WithArgs(targetUserID).
		WillReturnError(errors.New("db error"))

	_, _, _, err := repo.GetSummaryByUserID(context.Background(), targetUserID)
	assert.Error(t, err)
}
