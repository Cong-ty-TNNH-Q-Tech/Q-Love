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

func TestChatRepository_GetRecentMessages(t *testing.T) {
	db, mock, err := setupTestDB()
	if err != nil {
		t.Fatal(err)
	}

	repo := NewChatRepository(db)
	matchID := uuid.New()

	mock.ExpectQuery(`SELECT .* FROM "chat_messages"`).
		WithArgs(matchID, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "content"}).AddRow(uuid.New(), "Hello"))

	messages, err := repo.GetRecentMessages(context.Background(), matchID, 10)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(messages))
	}
}

func TestChatRepository_GetRecentMessages_Error(t *testing.T) {
	db, mock, err := setupTestDB()
	if err != nil {
		t.Fatal(err)
	}

	repo := NewChatRepository(db)
	matchID := uuid.New()

	mock.ExpectQuery(`SELECT .* FROM "chat_messages"`).
		WithArgs(matchID, sqlmock.AnyArg()).
		WillReturnError(sqlmock.ErrCancelled)

	_, err = repo.GetRecentMessages(context.Background(), matchID, 10)
	if err == nil {
		t.Error("Expected error")
	}
}
