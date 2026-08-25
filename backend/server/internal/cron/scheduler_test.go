// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package cron

import (
	"context"
	"testing"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/pkg/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type mockCronService struct {
	called bool
	err    error
}

func (m *mockCronService) RunWeeklyReset(ctx context.Context) error {
	m.called = true
	return m.err
}

type mockIslandCronService struct {
	called bool
	err    error
}

func (m *mockIslandCronService) RunDailyGhostingCheck(ctx context.Context) error {
	m.called = true
	return m.err
}

func TestScheduler_StartStop(t *testing.T) {
	if logger.Log == nil {
		logger.Log = zap.NewNop()
	}

	mockService := &mockCronService{}
	mockIslandService := &mockIslandCronService{}
	scheduler := NewScheduler(mockService, mockIslandService)

	assert.NotNil(t, scheduler)
	
	// Start should register cron and start
	scheduler.Start()
	
	// Stop should stop it safely
	scheduler.Stop()
}

func TestScheduler_CronJob(t *testing.T) {
	if logger.Log == nil {
		logger.Log = zap.NewNop()
	}

	mockService := &mockCronService{}
	mockIslandService := &mockIslandCronService{}
	scheduler := NewScheduler(mockService, mockIslandService)

	scheduler.Start()
	defer scheduler.Stop()

	entries := scheduler.c.Entries()
	assert.Len(t, entries, 2)

	// Execute all jobs directly
	for _, entry := range entries {
		entry.Job.Run()
	}
	assert.True(t, mockService.called)
	assert.True(t, mockIslandService.called)

	// Test error case
	mockService.called = false
	mockService.err = assert.AnError
	mockIslandService.called = false
	mockIslandService.err = assert.AnError

	for _, entry := range entries {
		entry.Job.Run()
	}
	assert.True(t, mockService.called)
	assert.True(t, mockIslandService.called)
}
