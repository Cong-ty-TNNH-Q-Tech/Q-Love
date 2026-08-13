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

	mock.ExpectQuery(`SELECT .* FROM "user_wallets"`).
		WithArgs(userID, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "balance"}).AddRow(userID, 0))

	mock.ExpectExec(`UPDATE "user_wallets"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

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

	mock.ExpectQuery(`SELECT .* FROM "user_wallets"`).
		WithArgs(userID, sqlmock.AnyArg()).
		WillReturnError(sqlmock.ErrCancelled)

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

	mock.ExpectExec(`INSERT INTO "wallet_transactions"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.CreateTransaction(context.Background(), txn)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestWalletRepository_GetWalletForUpdate(t *testing.T) {
	db, mock, err := setupTestDB()
	if err != nil {
		t.Fatal(err)
	}
	repo := NewWalletRepository(db)
	userID := uuid.New()

	mock.ExpectQuery(`SELECT .* FROM "user_wallets"`).
		WithArgs(userID, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "balance"}).AddRow(userID, 100.0))

	wallet, err := repo.GetWalletForUpdate(context.Background(), userID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if wallet == nil || wallet.Balance != 100.0 {
		t.Errorf("Expected wallet with balance 100, got %v", wallet)
	}
}

func TestWalletRepository_GetWalletForUpdate_Error(t *testing.T) {
	db, mock, err := setupTestDB()
	if err != nil {
		t.Fatal(err)
	}
	repo := NewWalletRepository(db)
	userID := uuid.New()

	mock.ExpectQuery(`SELECT .* FROM "user_wallets"`).
		WithArgs(userID, sqlmock.AnyArg()).
		WillReturnError(sqlmock.ErrCancelled)

	_, err = repo.GetWalletForUpdate(context.Background(), userID)
	if err == nil {
		t.Error("Expected error")
	}
}

func TestWalletRepository_UpdateBalance(t *testing.T) {
	db, mock, err := setupTestDB()
	if err != nil {
		t.Fatal(err)
	}
	repo := NewWalletRepository(db)
	userID := uuid.New()

	mock.ExpectExec(`UPDATE "user_wallets"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateBalance(context.Background(), userID, -10.0)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestWalletRepository_UpdateBalance_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open gorm db: %v", err)
	}

	repo := NewWalletRepository(gormDB)
	userID := uuid.New()
	delta := 50.0

	// Expect UPDATE query and simulate an error
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "user_wallets" SET "balance"=balance \+ \$1 WHERE user_id = \$2`).
		WithArgs(delta, userID).
		WillReturnError(gorm.ErrInvalidDB)
	mock.ExpectRollback()

	err = repo.UpdateBalance(context.Background(), userID, delta)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestWalletRepository_CheckTransactionExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open gorm db: %v", err)
	}

	repo := NewWalletRepository(gormDB)
	txID := uuid.New()

	// 1. Transaction exists
	mock.ExpectQuery(`SELECT count\(\*\) FROM "wallet_transactions" WHERE id = \$1`).
		WithArgs(txID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	exists, err := repo.CheckTransactionExists(context.Background(), txID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !exists {
		t.Errorf("Expected exists to be true")
	}

	// 2. Transaction does not exist
	mock.ExpectQuery(`SELECT count\(\*\) FROM "wallet_transactions" WHERE id = \$1`).
		WithArgs(txID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	exists, err = repo.CheckTransactionExists(context.Background(), txID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if exists {
		t.Errorf("Expected exists to be false")
	}

	// 3. DB Error
	mock.ExpectQuery(`SELECT count\(\*\) FROM "wallet_transactions" WHERE id = \$1`).
		WithArgs(txID).
		WillReturnError(gorm.ErrInvalidDB)

	_, err = repo.CheckTransactionExists(context.Background(), txID)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
