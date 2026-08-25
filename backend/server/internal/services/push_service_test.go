// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"context"
	"testing"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/pkg/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestPushService_SendPush(t *testing.T) {
	if logger.Log == nil {
		logger.Log = zap.NewNop()
	}

	service := NewPushService()

	err := service.SendPush(context.Background(), "user123", "Hello", "World", nil)
	assert.NoError(t, err)
}

func TestPushService_BroadcastToAll(t *testing.T) {
	if logger.Log == nil {
		logger.Log = zap.NewNop()
	}

	service := NewPushService()

	err := service.BroadcastToAll(context.Background(), "Hello All", "World", nil)
	assert.NoError(t, err)
}
