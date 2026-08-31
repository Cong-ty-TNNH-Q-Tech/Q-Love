// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"context"
	"sync"
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/pkg/logger"
	"go.uber.org/zap"
)

type PurgeWorker interface {
	Start(ctx context.Context, workers int)
}

type purgeWorkerImpl struct {
	purgeSvc PurgeService
}

func NewPurgeWorker(purgeSvc PurgeService) PurgeWorker {
	return &purgeWorkerImpl{purgeSvc: purgeSvc}
}

func (w *purgeWorkerImpl) Start(ctx context.Context, numWorkers int) {
	logger.Log.Info("Starting Purge Workers...", zap.Int("workers", numWorkers))

	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			w.runLoop(ctx, workerID)
		}(i)
	}
	
	go func() {
		wg.Wait()
		logger.Log.Info("All Purge Workers stopped.")
	}()
}

func (w *purgeWorkerImpl) runLoop(ctx context.Context, workerID int) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := w.purgeSvc.ProcessMatchmaking(ctx, 100) // Batch size 100
			if err != nil {
				logger.Log.Error("PurgeWorker error", zap.Int("workerID", workerID), zap.Error(err))
			}
		}
	}
}
