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

func TestNotificationRepository_Create(t *testing.T) {
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

	repo := NewNotificationRepository(gdb)

	id := uuid.New()
	userID := uuid.New()
	now := time.Now()

	notif := &models.Notification{
		ID:        id,
		UserID:    userID,
		Type:      "type",
		Payload:   "payload",
		Status:    "sent",
		CreatedAt: now,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "notifications"`).
		WithArgs(notif.UserID, notif.Type, notif.Payload, notif.Status, sqlmock.AnyArg(), sqlmock.AnyArg(), notif.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(notif.ID))
	mock.ExpectCommit()

	err = repo.Create(context.Background(), notif)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestNotificationRepository_UpdateStatus(t *testing.T) {
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

	repo := NewNotificationRepository(gdb)

	id := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE \"notifications\"").
		WithArgs("sent", id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.UpdateStatus(context.Background(), id, "sent")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
