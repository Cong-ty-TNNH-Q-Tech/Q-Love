// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"context"
	"sort"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/repository"
	"github.com/google/uuid"
)

type FeedUserResponse struct {
	User                 models.User `json:"user"`
	SpiritualMatchScore  int         `json:"spiritual_match_score,omitempty"`
	Zodiac               string      `json:"zodiac,omitempty"`
	Numerology           int         `json:"numerology,omitempty"`
}

type FeedService interface {
	GetFeed(ctx context.Context, userID uuid.UUID, filter string, radius int) ([]FeedUserResponse, error)
}

type feedService struct {
	userRepo         repository.UserRepository
	spiritualService SpiritualService
}

func NewFeedService(userRepo repository.UserRepository, spiritualService SpiritualService) FeedService {
	return &feedService{
		userRepo:         userRepo,
		spiritualService: spiritualService,
	}
}

func (s *feedService) GetFeed(ctx context.Context, userID uuid.UUID, filter string, radius int) ([]FeedUserResponse, error) {
	requestor, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Lấy danh sách users trong bán kính (mặc định 50km)
	var users []models.User
	if filter == "spiritual" {
		users, err = s.userRepo.GetSpiritualFeed(ctx, userID, radius)
	} else {
		users, err = s.userRepo.GetFeed(ctx, userID, radius)
	}
	
	if err != nil {
		return nil, err
	}

	var results []FeedUserResponse
	for _, u := range users {
		score := s.spiritualService.CalculateSpiritualMatchScore(requestor.DOB, u.DOB)

		results = append(results, FeedUserResponse{
			User:                u,
			SpiritualMatchScore: score,
			Zodiac:              s.spiritualService.CalculateZodiac(u.DOB),
			Numerology:          s.spiritualService.CalculateNumerology(u.DOB),
		})
	}

	// Sort results by spiritual match score descending if filter is spiritual
	if filter == "spiritual" {
		sort.SliceStable(results, func(i, j int) bool {
			return results[i].SpiritualMatchScore > results[j].SpiritualMatchScore
		})
	}

	return results, nil
}
