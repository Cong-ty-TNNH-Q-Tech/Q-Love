// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"context"
	"fmt"
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/repository"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"database/sql"
)

type CourtWorker struct {
	courtRepo          repository.CourtRepository
	userViolationRepo  repository.UserViolationRepository
	redisClient        *redis.Client
	logger             *zap.Logger
	walletRepo         repository.WalletRepository
	txManager          repository.TransactionManager
}

func NewCourtWorker(
	courtRepo repository.CourtRepository,
	userViolationRepo repository.UserViolationRepository,
	redisClient *redis.Client,
	logger *zap.Logger,
	walletRepo repository.WalletRepository,
	txManager repository.TransactionManager,
) *CourtWorker {
	return &CourtWorker{
		courtRepo:         courtRepo,
		userViolationRepo: userViolationRepo,
		redisClient:       redisClient,
		logger:            logger,
		walletRepo:        walletRepo,
		txManager:         txManager,
	}
}

func (w *CourtWorker) Start(ctx context.Context) {
	w.logger.Info("Starting CourtWorker...")
	
	// Start cron job to evaluate expired cases
	go w.evaluateExpiredCasesLoop(ctx)
	
	// Start consumer for new cases from redis
	if w.redisClient != nil {
		go w.consumeCourtCasesStream(ctx)
	}
}

func (w *CourtWorker) evaluateExpiredCasesLoop(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("evaluateExpiredCasesLoop stopped")
			return
		case <-ticker.C:
			w.evaluateExpiredCases(ctx)
		}
	}
}

func (w *CourtWorker) evaluateExpiredCases(ctx context.Context) {
	cases, err := w.courtRepo.GetExpiredVotingCases(ctx)
	if err != nil {
		w.logger.Error("Failed to get expired court cases", zap.Error(err))
		return
	}

	for _, c := range cases {
		total, guilty, err := w.courtRepo.CountVotesByCase(ctx, c.ID)
		if err != nil {
			w.logger.Error("Failed to count votes", zap.Error(err), zap.String("case_id", c.ID.String()))
			continue
		}

		var status models.CourtCaseStatus
		if total >= 50 && float64(guilty)/float64(total) > 0.65 {
			status = models.CourtCaseStatusGuilty
		} else {
			status = models.CourtCaseStatusNotGuilty
		}

		err = w.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
			if err := w.courtRepo.UpdateCaseStatus(txCtx, c.ID, status); err != nil {
				return err
			}

			courtFee := float64(50)
			if status == models.CourtCaseStatusGuilty {
				if err := w.walletRepo.UpdateBalance(txCtx, c.DefendantID, -courtFee); err != nil {
					return err
				}
				if err := w.walletRepo.UpdateBalance(txCtx, c.PlaintiffID, courtFee*2); err != nil {
					return err
				}
			} else {
				if err := w.walletRepo.UpdateBalance(txCtx, c.DefendantID, courtFee); err != nil {
					return err
				}
			}
			return nil
		}, &sql.TxOptions{Isolation: sql.LevelSerializable})

		if err != nil {
			w.logger.Error("Failed to resolve court case", zap.Error(err))
		}

		// Penalize defendant outside tx to avoid rolling back on external failure (if it fails, the case is already resolved)
		if status == models.CourtCaseStatusGuilty {
			violation := &models.UserViolation{
				UserID: c.DefendantID,
				Type:   "court_shadowban",
				Reason: "Found guilty by the Court of Love for ghosting",
			}
			if err := w.userViolationRepo.Create(ctx, violation); err != nil {
				w.logger.Error("Failed to create user violation", zap.Error(err))
			}
			if err := w.userViolationRepo.BanUser(ctx, c.DefendantID); err != nil {
				w.logger.Error("Failed to shadowban user", zap.Error(err))
			}
		}
	}
}

func (w *CourtWorker) consumeCourtCasesStream(ctx context.Context) {
	stream := "court_cases_stream"
	group := "court_worker_group"

	// Create consumer group
	err := w.redisClient.XGroupCreateMkStream(ctx, stream, group, "0").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		w.logger.Error("Failed to create redis consumer group", zap.Error(err))
		return
	}

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("consumeCourtCasesStream stopped")
			return
		default:
			streams, err := w.redisClient.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    group,
				Consumer: "worker-1",
				Streams:  []string{stream, ">"},
				Count:    10,
				Block:    5 * time.Second,
			}).Result()

			if err != nil {
				if err == redis.Nil {
					// No new messages
					continue
				}
				w.logger.Error("Failed to read from redis stream", zap.Error(err))
				time.Sleep(2 * time.Second)
				continue
			}

			for _, s := range streams {
				for _, msg := range s.Messages {
					caseIDStr, ok := msg.Values["case_id"].(string)
					if ok {
						w.logger.Info(fmt.Sprintf("Distributing case %s to 50 users (stubbed)", caseIDStr))
						// In real implementation, we would send FCM push notifications
						// to 50 random active users.
					}

					w.redisClient.XAck(ctx, stream, group, msg.ID)
				}
			}
		}
	}
}
