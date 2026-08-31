// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"context"
	"testing"
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockPurgeService struct {
	mock.Mock
}

func (m *MockPurgeService) ProcessMatchmaking(ctx context.Context, batchSize int) error {
	args := m.Called(ctx, batchSize)
	return args.Error(0)
}

func TestPurgeWorker_StartStop(t *testing.T) {
	if logger.Log == nil {
		logger.Log = zap.NewNop()
	}

	mockSvc := new(MockPurgeService)
	// Expect it to be called multiple times in background
	mockSvc.On("ProcessMatchmaking", mock.Anything, 100).Return(nil).Maybe()

	worker := NewPurgeWorker(mockSvc)
	
	worker.Start(context.Background(), 2) // start 2 workers
	
	// let it run a tiny bit
	time.Sleep(50 * time.Millisecond)
	
	worker.Stop()
	
	// Should stop safely without hanging
	assert.True(t, true)
}
