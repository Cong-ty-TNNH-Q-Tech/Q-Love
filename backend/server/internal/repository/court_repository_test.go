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
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupCourtTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	db.AutoMigrate(&models.CourtCase{}, &models.CourtVote{})
	return db
}

func TestCourtRepository_CreateCase(t *testing.T) {
	db := setupCourtTestDB(t)
	repo := NewCourtRepository(db)

	courtCase := &models.CourtCase{
		ID:        uuid.New(),
		MatchID:   uuid.New(),
		AccuserID: uuid.New(),
		AccusedID: uuid.New(),
		Reason:    "Toxic behavior",
		Status:    "Voting",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := repo.CreateCase(context.Background(), courtCase)
	assert.NoError(t, err)
}

func TestCourtRepository_FindCaseByID(t *testing.T) {
	db := setupCourtTestDB(t)
	repo := NewCourtRepository(db)

	caseID := uuid.New()
	courtCase := &models.CourtCase{
		ID:        caseID,
		MatchID:   uuid.New(),
		AccuserID: uuid.New(),
		AccusedID: uuid.New(),
		Reason:    "Toxic behavior",
		Status:    "Voting",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_ = repo.CreateCase(context.Background(), courtCase)

	found, err := repo.FindCaseByID(context.Background(), caseID)
	assert.NoError(t, err)
	assert.Equal(t, caseID, found.ID)
}

func TestCourtRepository_HasUserVoted(t *testing.T) {
	db := setupCourtTestDB(t)
	repo := NewCourtRepository(db)

	caseID := uuid.New()
	voterID := uuid.New()

	vote := &models.CourtVote{
		ID:      uuid.New(),
		CaseID:  caseID,
		VoterID: voterID,
		Vote:    "Guilty",
	}
	db.Create(vote)

	hasVoted, err := repo.HasUserVoted(context.Background(), caseID, voterID)
	assert.NoError(t, err)
	assert.True(t, hasVoted)
}

func TestCourtRepository_CreateVote(t *testing.T) {
	db := setupCourtTestDB(t)
	repo := NewCourtRepository(db)

	caseID := uuid.New()
	courtCase := &models.CourtCase{
		ID:        caseID,
		MatchID:   uuid.New(),
		AccuserID: uuid.New(),
		AccusedID: uuid.New(),
		Reason:    "Toxic behavior",
		Status:    "Voting",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_ = repo.CreateCase(context.Background(), courtCase)

	vote := &models.CourtVote{
		ID:      uuid.New(),
		CaseID:  caseID,
		VoterID: uuid.New(),
		Vote:    "Guilty",
	}

	err := repo.CreateVote(context.Background(), vote)
	assert.NoError(t, err)

	// Check if case counts were updated
	found, _ := repo.FindCaseByID(context.Background(), caseID)
	assert.Equal(t, 1, found.GuiltyVotes)
}

func TestCourtRepository_UpdateCaseStatus(t *testing.T) {
	db := setupCourtTestDB(t)
	repo := NewCourtRepository(db)

	caseID := uuid.New()
	courtCase := &models.CourtCase{
		ID:        caseID,
		MatchID:   uuid.New(),
		AccuserID: uuid.New(),
		AccusedID: uuid.New(),
		Reason:    "Toxic behavior",
		Status:    "Voting",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_ = repo.CreateCase(context.Background(), courtCase)

	err := repo.UpdateCaseStatus(context.Background(), caseID, "Guilty")
	assert.NoError(t, err)

	found, _ := repo.FindCaseByID(context.Background(), caseID)
	assert.Equal(t, "Guilty", found.Status)
}
