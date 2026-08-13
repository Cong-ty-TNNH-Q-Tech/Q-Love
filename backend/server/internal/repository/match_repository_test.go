package repository

import (
	"context"
	"testing"
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestMatchRepository_FindByID(t *testing.T) {
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

	repo := NewMatchRepository(gormDB)
	matchID := uuid.New()

	mock.ExpectQuery(`SELECT \* FROM "matches" WHERE id = \$1 ORDER BY "matches"\."id" LIMIT \$2`).
		WithArgs(matchID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(matchID))

	_, err = repo.FindByID(context.Background(), matchID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestMatchRepository_UpdateLastInteraction(t *testing.T) {
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

	repo := NewMatchRepository(gormDB)
	matchID := uuid.New()
	now := time.Now()

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "matches" SET "last_interaction_at"=\$1 WHERE id = \$2`).
		WithArgs(now, matchID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.UpdateLastInteraction(context.Background(), matchID, now)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
