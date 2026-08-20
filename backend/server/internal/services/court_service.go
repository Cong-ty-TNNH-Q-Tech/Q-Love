// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/repository"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type CourtService interface {
	FileLawsuit(ctx context.Context, plaintiffID, defendantID, matchID uuid.UUID, reason string) (*models.CourtCase, error)
	GetFeed(ctx context.Context, jurorID uuid.UUID, limit int) ([]models.CourtCase, error)
	VoteCase(ctx context.Context, caseID, jurorID uuid.UUID, voteType models.CourtVoteType) error
	WithdrawCase(ctx context.Context, caseID, plaintiffID uuid.UUID) error
}

type courtService struct {
	courtRepo    repository.CourtRepository
	matchRepo    repository.MatchRepository
	redisClient  *redis.Client
}

func NewCourtService(
	courtRepo repository.CourtRepository,
	matchRepo repository.MatchRepository,
	redisClient *redis.Client,
) CourtService {
	return &courtService{
		courtRepo:   courtRepo,
		matchRepo:   matchRepo,
		redisClient: redisClient,
	}
}

func (s *courtService) FileLawsuit(ctx context.Context, plaintiffID, defendantID, matchID uuid.UUID, reason string) (*models.CourtCase, error) {
	match, err := s.matchRepo.GetByID(ctx, matchID)
	if err != nil {
		return nil, errors.New("match not found")
	}

	if match.User1ID != plaintiffID && match.User2ID != plaintiffID {
		return nil, errors.New("unauthorized to file lawsuit for this match")
	}
	
	if match.User1ID != defendantID && match.User2ID != defendantID {
		return nil, errors.New("invalid defendant for this match")
	}

	if match.StreakScore <= 5 {
		return nil, errors.New("streak must be greater than 5 to file a lawsuit")
	}

	// check if last interaction is > 48h
	if time.Since(match.LastInteractionAt) < 48*time.Hour {
		return nil, errors.New("must wait at least 48 hours of silence before filing a lawsuit")
	}

	courtCase := &models.CourtCase{
		PlaintiffID: plaintiffID,
		DefendantID: defendantID,
		MatchID:     matchID,
		Reason:      reason,
		Status:      models.CourtCaseStatusVoting,
		ExpiresAt:   time.Now().Add(12 * time.Hour),
	}

	if err := s.courtRepo.CreateCase(ctx, courtCase); err != nil {
		return nil, err
	}

	// Publish to Redis Stream for CourtWorker
	if s.redisClient != nil {
		err := s.redisClient.XAdd(ctx, &redis.XAddArgs{
			Stream: "court_cases_stream",
			Values: map[string]interface{}{
				"case_id": courtCase.ID.String(),
			},
		}).Err()
		if err != nil {
			// Just log the error, don't fail the lawsuit creation
			fmt.Printf("Failed to push case to redis stream: %v\n", err)
		}
	}

	return courtCase, nil
}

func (s *courtService) GetFeed(ctx context.Context, jurorID uuid.UUID, limit int) ([]models.CourtCase, error) {
	// In a real app we might filter out cases the user already voted on.
	// For simplicity, we just return active cases.
	return s.courtRepo.GetActiveCases(ctx, limit)
}

func (s *courtService) VoteCase(ctx context.Context, caseID, jurorID uuid.UUID, voteType models.CourtVoteType) error {
	courtCase, err := s.courtRepo.GetCaseByID(ctx, caseID)
	if err != nil {
		return errors.New("case not found")
	}

	if courtCase.Status != models.CourtCaseStatusVoting {
		return errors.New("voting is closed for this case")
	}

	if courtCase.ExpiresAt.Before(time.Now()) {
		return errors.New("voting period has expired")
	}

	if courtCase.PlaintiffID == jurorID || courtCase.DefendantID == jurorID {
		return errors.New("plaintiff or defendant cannot vote in their own case")
	}

	hasVoted, err := s.courtRepo.HasUserVoted(ctx, caseID, jurorID)
	if err != nil {
		return err
	}
	if hasVoted {
		return errors.New("already voted")
	}

	vote := &models.CourtVote{
		CaseID:  caseID,
		JurorID: jurorID,
		Vote:    voteType,
	}

	return s.courtRepo.CreateVote(ctx, vote)
}

func (s *courtService) WithdrawCase(ctx context.Context, caseID, plaintiffID uuid.UUID) error {
	courtCase, err := s.courtRepo.GetCaseByID(ctx, caseID)
	if err != nil {
		return errors.New("case not found")
	}

	if courtCase.PlaintiffID != plaintiffID {
		return errors.New("unauthorized to withdraw this case")
	}

	if courtCase.Status != models.CourtCaseStatusVoting {
		return errors.New("can only withdraw cases that are currently in voting phase")
	}

	return s.courtRepo.UpdateCaseStatus(ctx, caseID, models.CourtCaseStatusSettled)
}
