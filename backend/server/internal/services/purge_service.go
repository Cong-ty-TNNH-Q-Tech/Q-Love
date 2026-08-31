// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/repository"
	"go.uber.org/zap"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/pkg/logger"
)

type PurgeService interface {
	ProcessMatchmaking(ctx context.Context, batchSize int) error
}

type purgeServiceImpl struct {
	queueRepo  repository.PurgeQueueRepository
	matchRepo  repository.MatchRepository
	pushSvc    PushService
	spiritualSvc SpiritualService
}

func NewPurgeService(queueRepo repository.PurgeQueueRepository, matchRepo repository.MatchRepository, pushSvc PushService, spiritualSvc SpiritualService) PurgeService {
	return &purgeServiceImpl{
		queueRepo: queueRepo,
		matchRepo: matchRepo,
		pushSvc:   pushSvc,
		spiritualSvc: spiritualSvc,
	}
}

func (s *purgeServiceImpl) ProcessMatchmaking(ctx context.Context, batchSize int) error {
	// Dequeue Normal users
	normalUsers, err := s.queueRepo.DequeueUsers(ctx, int64(batchSize))
	if err != nil {
		logger.Log.Error("Failed to dequeue normal users", zap.Error(err))
		return err
	}

	// Simple random matching for normal users
	s.matchUsersInPairs(ctx, normalUsers, false)

	// Dequeue VIP users
	vipUsers, err := s.queueRepo.DequeueVIPUsers(ctx, int64(batchSize))
	if err != nil {
		logger.Log.Error("Failed to dequeue VIP users", zap.Error(err))
		return err
	}

	// We can use Spiritual matching for VIP, for now just simple pairs
	s.matchUsersInPairs(ctx, vipUsers, true)

	return nil
}

func (s *purgeServiceImpl) matchUsersInPairs(ctx context.Context, users []string, isVIP bool) {
	// Pair adjacent users
	for i := 0; i < len(users)-1; i += 2 {
		user1 := users[i]
		user2 := users[i+1]
		uid1, err1 := uuid.Parse(user1)
		uid2, err2 := uuid.Parse(user2)
		if err1 != nil || err2 != nil {
			logger.Log.Error("Invalid UUID for purge match", zap.String("user1", user1), zap.String("user2", user2))
			continue
		}
		match := &models.Match{
			User1ID: uid1,
			User2ID: uid2,
		}
		err := s.matchRepo.Create(ctx, match)
		if err != nil {
			logger.Log.Error("Failed to create purge match", zap.String("user1", user1), zap.String("user2", user2), zap.Error(err))
			continue
		}
		
		// Send notification
		_ = s.pushSvc.SendPush(ctx, user1, "The Purge", "Bạn đã được ghép đôi trong Đêm Săn Mồi!", nil)
		_ = s.pushSvc.SendPush(ctx, user2, "The Purge", "Bạn đã được ghép đôi trong Đêm Săn Mồi!", nil)
	}
	
	// If one user left over, requeue them
	if len(users) > 0 && len(users)%2 != 0 {
		leftoverUser := users[len(users)-1]
		if _, err := uuid.Parse(leftoverUser); err == nil {
			_ = s.queueRepo.EnqueueUser(ctx, leftoverUser, isVIP)
		} else {
			logger.Log.Error("Invalid UUID for leftover user", zap.String("user", leftoverUser))
		}
	}
}
