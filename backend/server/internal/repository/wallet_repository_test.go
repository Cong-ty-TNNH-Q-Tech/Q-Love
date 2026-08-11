// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package repository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
)

func TestWalletRepository_AddCommission(t *testing.T) {
	db, mock, err := setupTestDB()
	if err != nil {
		t.Fatal(err)
	}

	repo := NewWalletRepository(db)
	userID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT .* FROM "user_wallets"`).
		WithArgs(userID, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "balance"}).AddRow(userID, 0))

	mock.ExpectExec(`UPDATE "user_wallets"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.AddCommission(context.Background(), userID, 10.0)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestWalletRepository_AddCommission_Error(t *testing.T) {
	db, mock, err := setupTestDB()
	if err != nil {
		t.Fatal(err)
	}

	repo := NewWalletRepository(db)
	userID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT .* FROM "user_wallets"`).
		WithArgs(userID, sqlmock.AnyArg()).
		WillReturnError(sqlmock.ErrCancelled)
	mock.ExpectRollback()

	err = repo.AddCommission(context.Background(), userID, 10.0)
	if err == nil {
		t.Error("Expected error")
	}
}

func TestWalletRepository_CreateTransaction(t *testing.T) {
	db, mock, err := setupTestDB()
	if err != nil {
		t.Fatal(err)
	}

	repo := NewWalletRepository(db)
	txn := &models.WalletTransaction{
		ID:          uuid.New(),
		UserID:      uuid.New(),
		Amount:      10.0,
		Type:        "wingman_commission",
		ReferenceID: uuid.New(),
	}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "wallet_transactions"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.CreateTransaction(context.Background(), txn)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
