// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package cron

import (
	"context"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/services"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/pkg/logger"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type Scheduler struct {
	c          *cron.Cron
	cronService services.ClanCronService
	islandCronService services.IslandCronService
}

func NewScheduler(cronService services.ClanCronService, islandCronService services.IslandCronService) *Scheduler {
	return &Scheduler{
		c:          cron.New(cron.WithSeconds()),
		cronService: cronService,
		islandCronService: islandCronService,
	}
}

func (s *Scheduler) Start() {
	// 0 0 * * 1 (0:00 every Monday)
	// Because cron.WithSeconds is used, we need 6 fields: Sec Min Hour Day Month DayOfWeek
	// So: 0 0 0 * * 1
	_, err := s.c.AddFunc("0 0 0 * * 1", func() {
		ctx := context.Background()
		logger.Log.Info("Triggering weekly clan reset from cron")
		if err := s.cronService.RunWeeklyReset(ctx); err != nil {
			logger.Log.Error("Failed to run weekly reset cron", zap.Error(err))
		}
	})
	if err != nil {
		logger.Log.Fatal("Failed to register clan cron job", zap.Error(err))
	}

	// 0 0 0 * * * (0:00 every day)
	_, err = s.c.AddFunc("0 0 0 * * *", func() {
		ctx := context.Background()
		logger.Log.Info("Triggering daily island ghosting check from cron")
		if err := s.islandCronService.RunDailyGhostingCheck(ctx); err != nil {
			logger.Log.Error("Failed to run daily island ghosting cron", zap.Error(err))
		}
	})
	if err != nil {
		logger.Log.Fatal("Failed to register island cron job", zap.Error(err))
	}

	s.c.Start()
	logger.Log.Info("Cron scheduler started")
}

func (s *Scheduler) Stop() {
	if s.c != nil {
		s.c.Stop()
	}
}
