// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package repository

import (
	"context"
	"testing"
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

func TestUserViolationRepository(t *testing.T) {
	db, mock, err := setupTestDB()
	if err != nil {
		t.Fatal(err)
	}

	repo := NewUserViolationRepository(db)
	ctx := context.Background()
	userID := uuid.New()

	violation := &models.UserViolation{
		ID:        uuid.New(),
		UserID:    userID,
		Type:      "nsfw",
		Reason:    "detected",
		IsActive:  true,
		CreatedAt: time.Now(),
	}

	// Test Create
	mock.ExpectQuery(`INSERT INTO "user_violations"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(violation.ID))

	err = repo.Create(ctx, violation)
	if err != nil {
		t.Fatalf("Failed to create violation: %v", err)
	}

	// Test CountActiveViolationsByType
	mock.ExpectQuery(`SELECT count\(\*\) FROM "user_violations" WHERE.*`).
		WithArgs(userID, "nsfw").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	count, err := repo.CountActiveViolationsByType(ctx, userID, "nsfw")
	if err != nil || count != 2 {
		t.Fatalf("Failed to count violations: %v", err)
	}

	// Test BanUser
	mock.ExpectExec(`UPDATE "users" SET "is_shadowbanned"=\$1 WHERE id = \$2`).
		WithArgs(true, userID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.BanUser(ctx, userID)
	if err != nil {
		t.Fatalf("Failed to ban user: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserViolationRepository_GetViolations(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewUserViolationRepository(gormDB)

	mock.ExpectQuery(`SELECT count\(\*\) FROM "user_violations" WHERE is_active = true`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	mock.ExpectQuery(`SELECT \* FROM "user_violations" WHERE is_active = true ORDER BY created_at DESC LIMIT \$1`).
		WithArgs(10).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}).AddRow(uuid.New(), uuid.New()))

	violations, total, err := repo.GetViolations(context.Background(), 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, violations, 1)

	// Test Count Error
	mock.ExpectQuery(`SELECT count\(\*\) FROM "user_violations" WHERE is_active = true`).
		WillReturnError(assert.AnError)
	_, _, err = repo.GetViolations(context.Background(), 1, 10)
	assert.ErrorIs(t, err, assert.AnError)
}

func TestUserViolationRepository_DeleteViolation(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewUserViolationRepository(gormDB)
	id := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "user_violations" SET "is_active"=\$1 WHERE id = \$2`).
		WithArgs(false, id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.DeleteViolation(context.Background(), id)
	assert.NoError(t, err)
}

