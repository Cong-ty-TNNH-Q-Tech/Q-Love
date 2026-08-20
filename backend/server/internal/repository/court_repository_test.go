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
	dbName := "file:" + uuid.New().String() + "?mode=memory&cache=shared"
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
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
		ID:          uuid.New(),
		PlaintiffID: uuid.New(),
		DefendantID: uuid.New(),
		MatchID:     uuid.New(),
		Reason:      "Ghosting",
		Status:      models.CourtCaseStatusVoting,
		ExpiresAt:   time.Now().Add(24 * time.Hour),
		CreatedAt:   time.Now(),
	}

	err := repo.CreateCase(context.Background(), courtCase)
	assert.NoError(t, err)
}

func TestCourtRepository_FindCaseByID(t *testing.T) {
	db := setupCourtTestDB(t)
	repo := NewCourtRepository(db)

	caseID := uuid.New()
	courtCase := &models.CourtCase{
		ID:          caseID,
		MatchID:     uuid.New(),
		PlaintiffID: uuid.New(),
		DefendantID: uuid.New(),
		Reason:      "Toxic behavior",
		Status:      "Voting",
		CreatedAt:   time.Now(),
	}
	_ = repo.CreateCase(context.Background(), courtCase)

	found, err := repo.GetCaseByID(context.Background(), caseID)
	assert.NoError(t, err)
	assert.Equal(t, caseID, found.ID)
}

func TestCourtRepository_HasUserVoted(t *testing.T) {
	db := setupCourtTestDB(t)
	repo := NewCourtRepository(db)

	caseID := uuid.New()
	jurorID := uuid.New()

	vote := &models.CourtVote{
		CaseID:    caseID,
		JurorID:   jurorID,
		Vote:      models.CourtVoteGuilty,
		CreatedAt: time.Now(),
	}
	db.Create(vote)

	hasVoted, err := repo.HasUserVoted(context.Background(), caseID, jurorID)
	assert.NoError(t, err)
	assert.True(t, hasVoted)
}

func TestCourtRepository_CreateVote(t *testing.T) {
	db := setupCourtTestDB(t)
	repo := NewCourtRepository(db)

	caseID := uuid.New()
	courtCase := &models.CourtCase{
		ID:          caseID,
		MatchID:     uuid.New(),
		PlaintiffID: uuid.New(),
		DefendantID: uuid.New(),
		Reason:      "Toxic behavior",
		Status:      "Voting",
		CreatedAt:   time.Now(),
	}
	_ = repo.CreateCase(context.Background(), courtCase)

	vote := &models.CourtVote{
		CaseID:    caseID,
		JurorID:   uuid.New(),
		Vote:      models.CourtVoteGuilty,
		CreatedAt: time.Now(),
	}

	err := repo.CreateVote(context.Background(), vote)
	assert.NoError(t, err)

	// We removed GuiltyVotes assertion because GuiltyVotes is not in models.CourtCase
}

func TestCourtRepository_UpdateCaseStatus(t *testing.T) {
	db := setupCourtTestDB(t)
	repo := NewCourtRepository(db)

	caseID := uuid.New()
	courtCase := &models.CourtCase{
		ID:          caseID,
		MatchID:     uuid.New(),
		PlaintiffID: uuid.New(),
		DefendantID: uuid.New(),
		Reason:      "Toxic behavior",
		Status:      models.CourtCaseStatusVoting,
		ExpiresAt:   time.Now().Add(24 * time.Hour),
		CreatedAt:   time.Now(),
	}
	_ = repo.CreateCase(context.Background(), courtCase)

	err := repo.UpdateCaseStatus(context.Background(), caseID, models.CourtCaseStatusGuilty)
	assert.NoError(t, err)

	found, _ := repo.GetCaseByID(context.Background(), caseID)
	assert.Equal(t, models.CourtCaseStatusGuilty, found.Status)
}

func TestCourtRepository_GetActiveCases(t *testing.T) {
	db := setupCourtTestDB(t)
	repo := NewCourtRepository(db)

	db.Create(&models.CourtCase{
		ID:          uuid.New(),
		PlaintiffID: uuid.New(),
		DefendantID: uuid.New(),
		MatchID:     uuid.New(),
		Reason:      "Active",
		Status:      models.CourtCaseStatusVoting,
		ExpiresAt:   time.Now().Add(48 * time.Hour),
	})
	
	cases, err := repo.GetActiveCases(context.Background(), 10)
	assert.NoError(t, err)
	assert.Len(t, cases, 1)
}

func TestCourtRepository_GetExpiredVotingCases(t *testing.T) {
	db := setupCourtTestDB(t)
	repo := NewCourtRepository(db)

	db.Create(&models.CourtCase{
		ID:          uuid.New(),
		PlaintiffID: uuid.New(),
		DefendantID: uuid.New(),
		MatchID:     uuid.New(),
		Reason:      "Expired",
		Status:      models.CourtCaseStatusVoting,
		ExpiresAt:   time.Now().Add(-1 * time.Hour),
	})
	
	cases, err := repo.GetExpiredVotingCases(context.Background())
	assert.NoError(t, err)
	assert.Len(t, cases, 1)
}

func TestCourtRepository_CountVotesByCase(t *testing.T) {
	db := setupCourtTestDB(t)
	repo := NewCourtRepository(db)

	caseID := uuid.New()
	db.Create(&models.CourtVote{CaseID: caseID, JurorID: uuid.New(), Vote: models.CourtVoteGuilty})
	db.Create(&models.CourtVote{CaseID: caseID, JurorID: uuid.New(), Vote: models.CourtVoteGuilty})
	db.Create(&models.CourtVote{CaseID: caseID, JurorID: uuid.New(), Vote: models.CourtVoteNotGuilty})
	
	total, guilty, err := repo.CountVotesByCase(context.Background(), caseID)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Equal(t, int64(2), guilty)
}
