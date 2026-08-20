// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"context"
	"errors"
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/repository"
	"github.com/google/uuid"
)

type ExRatingService interface {
	SubmitRating(ctx context.Context, targetUserID, matchID uuid.UUID, score int, tags []string) error
	ViewRating(ctx context.Context, viewerID, targetUserID uuid.UUID) (float64, int64, map[string]int, error)
}

type exRatingService struct {
	exRatingRepo repository.ExRatingRepository
	walletRepo   repository.WalletRepository
	txManager    repository.TransactionManager
	chatRepo     repository.ChatRepository
	matchRepo    repository.MatchRepository
}

func NewExRatingService(
	exRatingRepo repository.ExRatingRepository,
	walletRepo repository.WalletRepository,
	txManager repository.TransactionManager,
	chatRepo repository.ChatRepository,
	matchRepo repository.MatchRepository,
) ExRatingService {
	return &exRatingService{
		exRatingRepo: exRatingRepo,
		walletRepo:   walletRepo,
		txManager:    txManager,
		chatRepo:     chatRepo,
		matchRepo:    matchRepo,
	}
}

func (s *exRatingService) SubmitRating(ctx context.Context, targetUserID, matchID uuid.UUID, score int, tags []string) error {
	// Verify match exists and is inactive/unmatched
	match, err := s.matchRepo.FindByID(ctx, matchID)
	if err != nil {
		return errors.New("match not found")
	}

	if match.Status != "unmatched" {
		return errors.New("chỉ được đánh giá sau khi đã unmatch")
	}

	// Verify target user is part of the match
	if match.User1ID != targetUserID && match.User2ID != targetUserID {
		return errors.New("target_user_id không thuộc về match này")
	}

	// Check if already rated
	hasRated, err := s.exRatingRepo.HasRated(ctx, matchID, targetUserID)
	if err != nil {
		return err
	}
	if hasRated {
		return errors.New("đã đánh giá người này trong match này rồi")
	}

	// Verify message count >= 50
	msgCount, err := s.chatRepo.CountMessagesByMatchID(ctx, matchID)
	if err != nil {
		return err
	}
	if msgCount < 50 {
		return errors.New("chưa đủ điều kiện đánh giá (phải có ít nhất 50 tin nhắn)")
	}

	rating := &models.ExRating{
		ID:           uuid.New(),
		TargetUserID: targetUserID,
		MatchID:      matchID,
		RatingScore:  score,
		CreatedAt:    time.Now(),
	}
	rating.SetTags(tags)

	return s.exRatingRepo.Create(ctx, rating)
}

func (s *exRatingService) ViewRating(ctx context.Context, viewerID, targetUserID uuid.UUID) (float64, int64, map[string]int, error) {
	// Transaction to deduct 50 Xu
	err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		wallet, err := s.walletRepo.GetBalance(txCtx, viewerID)
		if err != nil {
			return errors.New("lỗi truy xuất ví")
		}

		if wallet.Balance < 50 {
			return errors.New("không đủ xu để tra cứu")
		}

		// Deduct Xu
		err = s.walletRepo.UpdateBalance(txCtx, viewerID, -50, "deduct", "Phí tra cứu CV Tình trường")
		if err != nil {
			return errors.New("lỗi trừ xu")
		}

		return nil
	})

	if err != nil {
		return 0, 0, nil, err
	}

	// If deduction succeeds, get the rating summary
	return s.exRatingRepo.GetSummaryByUserID(ctx, targetUserID)
}
