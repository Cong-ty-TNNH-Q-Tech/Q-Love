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
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestUserPremiumRepository_IsUserPremium(t *testing.T) {
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

	repo := NewUserPremiumRepository(gormDB)
	userID := uuid.New()

	// Setup expectations for true
	mock.ExpectQuery(`SELECT count\(\*\) FROM "user_premiums" WHERE user_id = \$1 AND expires_at > \$2`).
		WithArgs(userID, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	isPremium, err := repo.IsUserPremium(context.Background(), userID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !isPremium {
		t.Errorf("Expected true, got false")
	}

	// Setup expectations for false
	mock.ExpectQuery(`SELECT count\(\*\) FROM "user_premiums" WHERE user_id = \$1 AND expires_at > \$2`).
		WithArgs(userID, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	isPremium, err = repo.IsUserPremium(context.Background(), userID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if isPremium {
		t.Errorf("Expected false, got true")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestUserPremiumRepository_IsUserPremium_Error(t *testing.T) {
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

	repo := NewUserPremiumRepository(gormDB)
	userID := uuid.New()

	mock.ExpectQuery(`SELECT count\(\*\) FROM "user_premiums" WHERE user_id = \$1 AND expires_at > \$2`).
		WithArgs(userID, sqlmock.AnyArg()).
		WillReturnError(gorm.ErrInvalidDB)

	_, err = repo.IsUserPremium(context.Background(), userID)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestUserPremiumRepository_ActivatePremium(t *testing.T) {
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

	repo := NewUserPremiumRepository(gormDB)
	userID := uuid.New()
	expiresAt := time.Now().AddDate(0, 1, 0)

	// Test success
	mock.ExpectExec(`INSERT INTO user_premiums`).
		WithArgs(userID, expiresAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.ActivatePremium(context.Background(), userID, expiresAt)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test error
	mock.ExpectExec(`INSERT INTO user_premiums`).
		WithArgs(userID, expiresAt).
		WillReturnError(gorm.ErrInvalidDB)

	err = repo.ActivatePremium(context.Background(), userID, expiresAt)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
