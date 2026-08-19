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
	mock.ExpectQuery(`SELECT count\(\*\) FROM "user_violations" WHERE user_id = \$1 AND type = \$2 AND is_active = true.*`).
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
