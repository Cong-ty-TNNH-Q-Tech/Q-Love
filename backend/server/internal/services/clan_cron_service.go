// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"context"
	"fmt"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/repository"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/pkg/logger"
	"go.uber.org/zap"
)

type ClanCronService interface {
	RunWeeklyReset(ctx context.Context) error
}

type clanCronService struct {
	clanRepo     repository.ClanRepository
	landmarkRepo repository.LandmarkRepository
	notifRepo    repository.NotificationRepository
	pushService  PushService
}

func NewClanCronService(
	clanRepo repository.ClanRepository,
	landmarkRepo repository.LandmarkRepository,
	notifRepo repository.NotificationRepository,
	pushService PushService,
) ClanCronService {
	return &clanCronService{
		clanRepo:     clanRepo,
		landmarkRepo: landmarkRepo,
		notifRepo:    notifRepo,
		pushService:  pushService,
	}
}

func (s *clanCronService) RunWeeklyReset(ctx context.Context) error {
	logger.Log.Info("Starting Clan Weekly Reset Cronjob")
	
	// 1. Get Top Clan
	topClan, err := s.clanRepo.GetTopWeeklyClan(ctx)
	if err != nil {
		logger.Log.Error("Failed to get top weekly clan", zap.Error(err))
		return err
	}

	// 2. Update Landmarks owner
	if err := s.landmarkRepo.UpdateAllOwners(ctx, topClan); err != nil {
		logger.Log.Error("Failed to update landmarks owner", zap.Error(err))
		return err
	}

	// 3. Broadcast notification
	if topClan != nil {
		body := fmt.Sprintf("Bang hội %s đã chiếm lĩnh toàn bộ bản đồ tuần này với %d điểm!", topClan.Name, topClan.WeeklyScore)
		_ = s.pushService.BroadcastToAll(ctx, "👑 Vị Vua Mới Của Tuần", body, map[string]string{
			"type": "clan_weekly_result",
			"clan_id": topClan.ID.String(),
		})
	}

	// 4. Reset scores
	if err := s.clanRepo.ResetWeeklyScores(ctx); err != nil {
		logger.Log.Error("Failed to reset weekly scores", zap.Error(err))
		return err
	}

	logger.Log.Info("Finished Clan Weekly Reset Cronjob successfully")
	return nil
}
