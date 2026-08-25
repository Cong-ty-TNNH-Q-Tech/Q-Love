// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package repository

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func TestCourtCaseRepository(t *testing.T) {
	db, mock, err := setupTestDB()
	if err != nil {
		t.Fatal(err)
	}

	repo := NewCourtCaseRepository(db)
	ctx := context.Background()
	caseID := uuid.New()
	userID := uuid.New()
	targetID := uuid.New()

	// Test UpdateStatus
	mock.ExpectExec(`UPDATE "court_cases" SET "status"=\$1 WHERE id = \$2`).
		WithArgs("rejected", caseID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateStatus(ctx, caseID, "rejected")
	if err != nil {
		t.Fatalf("Failed to update status: %v", err)
	}

	// Test FindByID Success
	createdAt := time.Now()
	mock.ExpectQuery(`SELECT \* FROM "court_cases" WHERE id = \$1 ORDER BY "court_cases"."id" LIMIT \$2`).
		WithArgs(caseID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "plaintiff_id", "defendant_id", "status", "created_at"}).
			AddRow(caseID, userID, targetID, "pending", createdAt))

	foundCase, err := repo.FindByID(ctx, caseID)
	if err != nil {
		t.Fatalf("Failed to find court case: %v", err)
	}
	if foundCase == nil || foundCase.ID != caseID {
		t.Fatalf("Expected court case ID %v, got %v", caseID, foundCase.ID)
	}

	// Test FindByID Error
	mock.ExpectQuery(`SELECT \* FROM "court_cases" WHERE id = \$1 ORDER BY "court_cases"."id" LIMIT \$2`).
		WithArgs(caseID, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	_, err = repo.FindByID(ctx, caseID)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
