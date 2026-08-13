package repository

import (
	"context"
	"testing"

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
