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
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	// Test GetViolations
	mock.ExpectQuery(`SELECT count\(\*\) FROM "user_violations" WHERE is_active = true AND "user_violations"\."deleted_at" IS NULL`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery(`SELECT \* FROM "user_violations" WHERE is_active = true AND "user_violations"\."deleted_at" IS NULL ORDER BY created_at DESC LIMIT \$1`).
		WithArgs(10).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(violation.ID))

	violations, total, err := repo.GetViolations(ctx, 1, 10)
	if err != nil || total != 1 || len(violations) != 1 {
		t.Fatalf("Failed to get violations: %v", err)
	}

	// Test DeleteViolation
	mock.ExpectExec(`UPDATE "user_violations" SET "is_active"=\$1 WHERE id = \$2 AND "user_violations"\."deleted_at" IS NULL`).
		WithArgs(false, violation.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.DeleteViolation(ctx, violation.ID)
	if err != nil {
		t.Fatalf("Failed to delete violation: %v", err)
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

	mock.ExpectQuery(`SELECT count\(\*\) FROM "user_violations" WHERE is_active = true AND "user_violations"\."deleted_at" IS NULL`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	mock.ExpectQuery(`SELECT \* FROM "user_violations" WHERE is_active = true AND "user_violations"\."deleted_at" IS NULL ORDER BY created_at DESC LIMIT \$1`).
		WithArgs(10).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}).AddRow(uuid.New(), uuid.New()))

	violations, total, err := repo.GetViolations(context.Background(), 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, violations, 1)

	// Test Count Error
	mock.ExpectQuery(`SELECT count\(\*\) FROM "user_violations" WHERE is_active = true AND "user_violations"\."deleted_at" IS NULL`).
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
	mock.ExpectExec(`UPDATE "user_violations" SET "is_active"=\$1 WHERE id = \$2 AND "user_violations"\."deleted_at" IS NULL`).
		WithArgs(false, id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.DeleteViolation(context.Background(), id)
	assert.NoError(t, err)
}

func TestUserViolationRepository_HasActiveFakeGPSBan(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewUserViolationRepository(gormDB)
	userID := uuid.New()

	t.Run("Has Ban", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "user_violations" WHERE .*`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(uuid.New(), time.Now().Add(-24*time.Hour)))

		hasBan, expiresAt, err := repo.HasActiveFakeGPSBan(context.Background(), userID)
		assert.NoError(t, err)
		assert.True(t, hasBan)
		assert.NotNil(t, expiresAt)
	})

	t.Run("No Ban", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "user_violations" WHERE .*`).
			WillReturnError(gorm.ErrRecordNotFound)

		hasBan, expiresAt, err := repo.HasActiveFakeGPSBan(context.Background(), userID)
		assert.NoError(t, err)
		assert.False(t, hasBan)
		assert.Nil(t, expiresAt)
	})
}
