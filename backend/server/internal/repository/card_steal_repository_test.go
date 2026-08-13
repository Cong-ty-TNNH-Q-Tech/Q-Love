// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package repository

import (
	"context"
	"testing"
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	db.AutoMigrate(&models.CardSteal{}, &models.CardTransaction{})
	return db
}

func TestCardStealRepository(t *testing.T) {
	db := setupTestDB()
	repo := NewCardStealRepository(db)

	steal := &models.CardSteal{
		ID:           uuid.New(),
		AttackerID:   uuid.New(),
		DefenderID:   uuid.New(),
		TargetCardID: uuid.New(),
		Result:       "pending",
		CreatedAt:    time.Now(),
	}

	err := repo.Create(context.Background(), steal)
	if err != nil {
		t.Fatalf("Failed to create steal: %v", err)
	}

	found, err := repo.FindByID(context.Background(), steal.ID)
	if err != nil || found.ID != steal.ID {
		t.Fatalf("Failed to find steal: %v", err)
	}

	err = repo.UpdateResult(context.Background(), steal.ID, "attacker_won")
	if err != nil {
		t.Fatalf("Failed to update steal: %v", err)
	}

	err = repo.TransferCardOwnership(context.Background(), uuid.New(), uuid.New())
	if err != nil {
		t.Fatalf("Failed to transfer ownership: %v", err)
	}
}
