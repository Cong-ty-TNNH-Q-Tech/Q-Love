// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"context"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/pkg/logger"
	"go.uber.org/zap"
)

type PushService interface {
	SendPush(ctx context.Context, userID string, title string, body string, payload map[string]string) error
	BroadcastToAll(ctx context.Context, title string, body string, payload map[string]string) error
}

type pushService struct {
}

func NewPushService() PushService {
	return &pushService{}
}

func (s *pushService) SendPush(ctx context.Context, userID string, title string, body string, payload map[string]string) error {
	// MOCK: Integration with FCM or APNs goes here
	logger.Log.Info("Sending push notification", zap.String("userID", userID), zap.String("title", title))
	return nil
}

func (s *pushService) BroadcastToAll(ctx context.Context, title string, body string, payload map[string]string) error {
	// MOCK: Send to FCM topic "all_users"
	logger.Log.Info("Broadcasting push notification to all users", zap.String("title", title))
	return nil
}
