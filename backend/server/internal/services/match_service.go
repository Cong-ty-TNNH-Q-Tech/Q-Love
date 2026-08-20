// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.


package services

import (
	"context"
	"errors"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/repository"
	"github.com/google/uuid"
)

type MatchService interface {
	Unmatch(ctx context.Context, matchID, userID uuid.UUID) error
}

type matchService struct {
	matchRepo repository.MatchRepository
}

func NewMatchService(matchRepo repository.MatchRepository) MatchService {
	return &matchService{matchRepo: matchRepo}
}

func (s *matchService) Unmatch(ctx context.Context, matchID, userID uuid.UUID) error {
	match, err := s.matchRepo.FindByID(ctx, matchID)
	if err != nil {
		return err
	}

	if match == nil {
		return errors.New("match not found")
	}

	if match.User1ID != userID && match.User2ID != userID {
		return errors.New("forbidden")
	}

	err = s.matchRepo.SoftDelete(ctx, matchID)
	if err != nil {
		return err
	}

	// TODO: Trigger luồng cho phép đánh giá CV Tình trường (Ex-Rating)
	return nil
}

