// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.


package services

import (
	"context"
	"errors"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/repository"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"fmt"
	"time"
)

type MatchService interface {
	Unmatch(ctx context.Context, matchID, userID uuid.UUID) error
}

type matchService struct {
	matchRepo   repository.MatchRepository
	notifSvc    NotificationService
	redisClient *redis.Client
}

func NewMatchService(matchRepo repository.MatchRepository, notifSvc NotificationService, redisClient *redis.Client) MatchService {
	return &matchService{
		matchRepo:   matchRepo,
		notifSvc:    notifSvc,
		redisClient: redisClient,
	}
}

func (s *matchService) Unmatch(ctx context.Context, matchID, userID uuid.UUID) error {
	match, err := s.matchRepo.FindByID(ctx, matchID)
	if err != nil {
		return err
	}

	if match == nil {
		return errors.New("match not found")
	}

	if match.User1ID != userID && match.User2ID != userID {
		return errors.New("forbidden")
	}

	err = s.matchRepo.SoftDelete(ctx, matchID)
	if err != nil {
		return err
	}

	// Trigger luồng cho phép đánh giá CV Tình trường (Ex-Rating)
	if s.redisClient != nil {
		partnerID := match.User2ID
		if match.User2ID == userID {
			partnerID = match.User1ID
		}

		userKey := fmt.Sprintf("pending_ex_ratings:%s", userID.String())
		partnerKey := fmt.Sprintf("pending_ex_ratings:%s", partnerID.String())

		// Add to pending sets
		s.redisClient.SAdd(ctx, userKey, partnerID.String())
		s.redisClient.SAdd(ctx, partnerKey, userID.String())

		// Set 24h expiration
		s.redisClient.Expire(ctx, userKey, 24*time.Hour)
		s.redisClient.Expire(ctx, partnerKey, 24*time.Hour)

		if s.notifSvc != nil {
			s.notifSvc.SendPush(ctx, userID, "ex_rating", "Đánh giá tình cũ", "Bạn có 24h để đánh giá CV Tình trường của người cũ", nil)
			s.notifSvc.SendPush(ctx, partnerID, "ex_rating", "Đánh giá tình cũ", "Bạn có 24h để đánh giá CV Tình trường của người cũ", nil)
		}
	}

	return nil
}

