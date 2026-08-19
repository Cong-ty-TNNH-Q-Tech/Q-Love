// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package repository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

func TestShameRepository_GetActiveShames(t *testing.T) {
	db, mock, err := setupTestDB()
	if err != nil {
		t.Fatal(err)
	}
	repo := NewShameRepository(db)

	mock.ExpectQuery(`SELECT (.*) FROM "wall_of_shames"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}).AddRow(uuid.New(), uuid.New()))

	shames, err := repo.GetActiveShames(context.Background(), 10, 0)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(shames) != 1 {
		t.Errorf("Expected 1 shame, got %d", len(shames))
	}
}

func TestShameRepository_GetActiveShames_Error(t *testing.T) {
	db, mock, err := setupTestDB()
	if err != nil {
		t.Fatal(err)
	}
	repo := NewShameRepository(db)

	mock.ExpectQuery(`SELECT (.*) FROM "wall_of_shames"`).
		WithArgs(sqlmock.AnyArg()).
		WillReturnError(sqlmock.ErrCancelled)

	_, err = repo.GetActiveShames(context.Background(), 10, 0)
	if err == nil {
		t.Error("Expected error")
	}
}

func TestShameRepository_IncrementTomatoCount(t *testing.T) {
	db, mock, err := setupTestDB()
	if err != nil {
		t.Fatal(err)
	}
	repo := NewShameRepository(db)
	shameID := uuid.New()

	mock.ExpectExec(`UPDATE "wall_of_shames"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.IncrementTomatoCount(context.Background(), shameID, 1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestShameRepository_IncrementTomatoCount_Error(t *testing.T) {
	db, mock, err := setupTestDB()
	if err != nil {
		t.Fatal(err)
	}
	repo := NewShameRepository(db)
	shameID := uuid.New()

	mock.ExpectExec(`UPDATE "wall_of_shames"`).
		WithArgs(1, shameID).
		WillReturnError(sqlmock.ErrCancelled)

	err = repo.IncrementTomatoCount(context.Background(), shameID, 1)
	if err == nil {
		t.Error("Expected error")
	}
}
