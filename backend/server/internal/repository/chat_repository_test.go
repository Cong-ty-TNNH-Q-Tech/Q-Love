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
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestChatRepository_SaveMessage(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening gorm database", err)
	}

	repo := NewChatMessageRepository(gdb)

	msgID := uuid.New()
	matchID := uuid.New()
	senderID := uuid.New()
	now := time.Now()

	msg := &models.ChatMessage{
		ID:        msgID,
		MatchID:   matchID,
		SenderID:  senderID,
		Type:      "text",
		Content:   "hello",
		CreatedAt: now,
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO \"chat_messages\"").
		WithArgs(msg.ID, msg.MatchID, msg.SenderID, msg.Type, msg.Content, msg.BlurURL, msg.BlurLevel, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.Create(context.Background(), msg)
	if err != nil {
		t.Errorf("error was not expected while saving message: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestChatRepository_GetMessagesByMatchID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening gorm database", err)
	}

	repo := NewChatMessageRepository(gdb)

	matchID := uuid.New()
	
	rows := sqlmock.NewRows([]string{"id", "match_id", "sender_id", "type", "content"}).
		AddRow(uuid.New().String(), matchID.String(), uuid.New().String(), "text", "hi")

	mock.ExpectQuery("SELECT \\* FROM \"chat_messages\" WHERE match_id = \\$1 ORDER BY created_at DESC LIMIT \\$2").
		WithArgs(matchID, 50).
		WillReturnRows(rows)

	messages, err := repo.GetMessagesByMatchID(context.Background(), matchID, 50, nil)
	if err != nil {
		t.Errorf("error was not expected while getting messages: %s", err)
	}

	if len(messages) != 1 {
		t.Errorf("expected 1 message, got %d", len(messages))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
