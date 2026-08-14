// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

func TestCardStealRepository(t *testing.T) {
	db, mock, err := setupTestDB()
	if err != nil {
		t.Fatal(err)
	}

	repo := NewCardStealRepository(db)
	ctx := context.Background()
	testID := uuid.New()
	attackerID := uuid.New()
	defenderID := uuid.New()
	cardID := uuid.New()

	steal := &models.CardSteal{
		ID:           testID,
		AttackerID:   attackerID,
		DefenderID:   defenderID,
		TargetCardID: cardID,
		Result:       "pending",
		CreatedAt:    time.Now(),
	}

	// 1. Test Create
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "card_steals"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(testID))
	mock.ExpectCommit()

	err = repo.Create(ctx, steal)
	if err != nil {
		t.Fatalf("Failed to create steal: %v", err)
	}

	// 2. Test FindByID
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "card_steals" WHERE id = $1`)).
		WithArgs(testID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "attacker_id", "defender_id", "target_card_id", "result"}).
			AddRow(testID, attackerID, defenderID, cardID, "pending"))

	found, err := repo.FindByID(ctx, testID)
	if err != nil || found.ID != testID {
		t.Fatalf("Failed to find steal: %v", err)
	}

	// 3. Test UpdateResult
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "card_steals" SET "result"=$1`)).
		WithArgs("attacker_won", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.UpdateResult(ctx, testID, "attacker_won")
	if err != nil {
		t.Fatalf("Failed to update steal: %v", err)
	}

	// 4. Test TransferCardOwnership
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "card_transactions"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
	mock.ExpectCommit()

	err = repo.TransferCardOwnership(ctx, attackerID, cardID)
	if err != nil {
		t.Fatalf("Failed to transfer ownership: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
