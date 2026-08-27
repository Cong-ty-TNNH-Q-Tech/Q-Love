// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"context"
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/repository"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/pkg/logger"
	"go.uber.org/zap"
)

type IslandCronService interface {
	RunDailyGhostingCheck(ctx context.Context) error
}

type islandCronService struct {
	matchRepo repository.MatchRepository
}

func NewIslandCronService(matchRepo repository.MatchRepository) IslandCronService {
	return &islandCronService{
		matchRepo: matchRepo,
	}
}

func (s *islandCronService) RunDailyGhostingCheck(ctx context.Context) error {
	logger.Log.Info("Starting Daily Ghosting Check Cronjob (Island Streak & Level)")

	// 1. Reset streak for matches inactive > 24 hours
	if err := s.matchRepo.ResetStreakForInactiveMatches(ctx, 24*time.Hour); err != nil {
		logger.Log.Error("Failed to reset streak for inactive matches", zap.Error(err))
		return err
	}

	// 2. Destroy islands (reset level) for matches inactive > 7 days (168 hours)
	if err := s.matchRepo.ResetIslandLevelForInactiveMatches(ctx, 7*24*time.Hour); err != nil {
		logger.Log.Error("Failed to reset island level for inactive matches", zap.Error(err))
		return err
	}

	logger.Log.Info("Completed Daily Ghosting Check Cronjob")
	return nil
}
