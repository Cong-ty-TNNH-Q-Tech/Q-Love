// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/repository"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type NotificationService interface {
	SendPush(ctx context.Context, userID uuid.UUID, nType, title, body string, data map[string]string) error
	SendSilentPush(ctx context.Context, userID uuid.UUID, data map[string]string) error
}

type notificationService struct {
	repo        repository.NotificationRepository
	redisClient *redis.Client
	fcmKey      string
	httpClient  *http.Client
}

func NewNotificationService(repo repository.NotificationRepository, redisClient *redis.Client, fcmKey string) NotificationService {
	return &notificationService{
		repo:        repo,
		redisClient: redisClient,
		fcmKey:      fcmKey,
		httpClient:  &http.Client{},
	}
}

func (s *notificationService) getFCMToken(ctx context.Context, userID uuid.UUID) (string, error) {
	if s.redisClient == nil {
		return "", fmt.Errorf("redis client is nil")
	}
	key := fmt.Sprintf("fcm_token:%s", userID.String())
	token, err := s.redisClient.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *notificationService) SendPush(ctx context.Context, userID uuid.UUID, nType, title, body string, data map[string]string) error {
	token, err := s.getFCMToken(ctx, userID)
	if err != nil {
		// Log failed attempt to DB if token is missing
		s.logNotification(ctx, userID, nType, data, "failed_no_token")
		return err
	}

	payload := map[string]interface{}{
		"to": token,
		"notification": map[string]string{
			"title": title,
			"body":  body,
		},
		"data": data,
	}

	err = s.sendToFCM(payload)
	status := "sent"
	if err != nil {
		status = "failed"
	}

	s.logNotification(ctx, userID, nType, data, status)
	return err
}

func (s *notificationService) SendSilentPush(ctx context.Context, userID uuid.UUID, data map[string]string) error {
	token, err := s.getFCMToken(ctx, userID)
	if err != nil {
		s.logNotification(ctx, userID, "silent_push", data, "failed_no_token")
		return err
	}

	// For silent push on iOS using FCM, content_available must be true and no notification block
	payload := map[string]interface{}{
		"to": token,
		"content_available": true,
		"priority":          "high",
		"data":              data,
	}

	err = s.sendToFCM(payload)
	status := "sent"
	if err != nil {
		status = "failed"
	}

	s.logNotification(ctx, userID, "silent_push", data, status)
	return err
}

func (s *notificationService) sendToFCM(payload map[string]interface{}) error {
	if s.fcmKey == "" {
		// Mock behavior when key is not provided (for dev/test)
		return nil
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://fcm.googleapis.com/fcm/send", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "key="+s.fcmKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("FCM API error: status %d", resp.StatusCode)
	}

	return nil
}

func (s *notificationService) logNotification(ctx context.Context, userID uuid.UUID, nType string, data map[string]string, status string) {
	payloadStr := "{}"
	if b, err := json.Marshal(data); err == nil {
		payloadStr = string(b)
	}
	notif := &models.Notification{
		UserID:  userID,
		Type:    nType,
		Payload: payloadStr,
		Status:  status,
	}
	_ = s.repo.Create(ctx, notif)
}
