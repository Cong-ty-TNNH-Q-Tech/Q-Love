// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/repository"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/pkg/logger"
	"github.com/google/uuid"
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
	walletRepo   repository.WalletRepository
	txManager    repository.TransactionManager
}

func NewClanCronService(
	clanRepo repository.ClanRepository,
	landmarkRepo repository.LandmarkRepository,
	notifRepo repository.NotificationRepository,
	pushService PushService,
	walletRepo repository.WalletRepository,
	txManager repository.TransactionManager,
) ClanCronService {
	return &clanCronService{
		clanRepo:     clanRepo,
		landmarkRepo: landmarkRepo,
		notifRepo:    notifRepo,
		pushService:  pushService,
		walletRepo:   walletRepo,
		txManager:    txManager,
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

	// 2. Use Transaction to reward members and reset scores
	err = s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		// Reward members of top clan
		if topClan != nil {
			members, err := s.clanRepo.GetMembers(txCtx, topClan.ID)
			if err != nil {
				return err
			}
			for _, m := range members {
				if err := s.walletRepo.UpdateBalance(txCtx, m.UserID, 100); err != nil { // 100 Xu reward
					return err
				}
			}
		}

		// Update Landmarks owner
		if err := s.landmarkRepo.UpdateAllOwners(txCtx, topClan); err != nil {
			return err
		}

		// Reset scores
		if err := s.clanRepo.ResetWeeklyScores(txCtx); err != nil {
			return err
		}
		
		return nil
	}, &sql.TxOptions{Isolation: sql.LevelSerializable})

	if err != nil {
		logger.Log.Error("Failed to process weekly reset in transaction", zap.Error(err))
		return err
	}

	// 3. Broadcast notification
	if topClan != nil {
		body := fmt.Sprintf("Bang hội %s đã chiếm lĩnh toàn bộ bản đồ tuần này với %d điểm!", topClan.Name, topClan.WeeklyScore)
		
		notification := &models.Notification{
			ID:          uuid.New(),
			UserID:      uuid.Nil, // global notification
			Type:        "clan_weekly_result",
			Payload:     body,
			ReferenceID: &topClan.ID,
		}
		_ = s.notifRepo.Create(ctx, notification)

		_ = s.pushService.BroadcastToAll(ctx, "👑 Vị Vua Mới Của Tuần", body, map[string]string{
			"type": "clan_weekly_result",
			"clan_id": topClan.ID.String(),
		})
	}

	logger.Log.Info("Finished Clan Weekly Reset Cronjob successfully")
	return nil
}
