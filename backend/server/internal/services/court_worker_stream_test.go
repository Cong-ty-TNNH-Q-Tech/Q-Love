// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	go_redis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func setupCourtWorkerTest() (*miniredis.Miniredis, *go_redis.Client) {
	mr, _ := miniredis.Run()
	client := go_redis.NewClient(&go_redis.Options{
		Addr: mr.Addr(),
	})
	return mr, client
}

func TestCourtWorker_ConsumeCourtCasesStream(t *testing.T) {
	mr, rdb := setupCourtWorkerTest()
	defer mr.Close()

	logger := zap.NewNop()
	worker := NewCourtWorker(nil, nil, rdb, logger)

	ctx, cancel := context.WithCancel(context.Background())
	
	// Start consuming
	go worker.consumeCourtCasesStream(ctx)
	
	// Give it time to create group
	time.Sleep(100 * time.Millisecond)
	
	// Add message to stream
	rdb.XAdd(context.Background(), &go_redis.XAddArgs{
		Stream: "court_cases_stream",
		Values: map[string]interface{}{
			"case_id": "test-case-id",
		},
	})
	
	time.Sleep(200 * time.Millisecond)
	cancel() // stop consumer
	time.Sleep(100 * time.Millisecond)
}

func TestCourtWorker_ConsumeCourtCasesStream_GroupCreateError(t *testing.T) {
	mr, rdb := setupCourtWorkerTest()
	defer mr.Close()

	logger := zap.NewNop()
	worker := NewCourtWorker(nil, nil, rdb, logger)

	// Inject error into miniredis
	mr.SetError("ERR_REDIS")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Should return early due to group create error
	worker.consumeCourtCasesStream(ctx)
}
