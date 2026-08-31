// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/repository"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type CardService interface {
	GetProfile(ctx context.Context, userID uuid.UUID) (*models.CardProfile, error)
	TradeCard(ctx context.Context, collectorID, targetUserID uuid.UUID, tradeType string, quantity int) error
}

type cardService struct {
	cardRepo   repository.CardRepository
	walletRepo repository.WalletRepository
	userRepo   repository.UserRepository
	txManager  repository.TransactionManager
	redisClient *redis.Client
}

func NewCardService(
	cardRepo repository.CardRepository,
	walletRepo repository.WalletRepository,
	userRepo repository.UserRepository,
	txManager repository.TransactionManager,
	redisClient *redis.Client,
) CardService {
	return &cardService{
		cardRepo:   cardRepo,
		walletRepo: walletRepo,
		userRepo:   userRepo,
		txManager:  txManager,
		redisClient: redisClient,
	}
}

func (s *cardService) GetProfile(ctx context.Context, userID uuid.UUID) (*models.CardProfile, error) {
	profile, err := s.cardRepo.GetCardProfile(ctx, userID)
	if err != nil {
		// If not found, we create a default profile on the fly
		profile = &models.CardProfile{
			UserID:         userID,
			CurrentPrice:   100,
			TotalCards:     1000,
			AvailableCards: 1000,
		}
		if err := s.cardRepo.CreateCardProfile(ctx, profile); err != nil {
			return nil, err
		}
	}

	// Check Circuit Breaker status from Redis
	isHalted := false
	if s.redisClient != nil {
		val, _ := s.redisClient.Get(ctx, fmt.Sprintf("circuit_breaker:%s", userID.String())).Result()
		if val == "1" {
			isHalted = true
		}
	}
	profile.IsHalted = &isHalted

	// Generate Ticker (e.g. #UserID_Prefix)
	profile.Ticker = fmt.Sprintf("#%s", userID.String()[:6])

	return profile, nil
}

func (s *cardService) TradeCard(ctx context.Context, collectorID, targetUserID uuid.UUID, tradeType string, quantity int) error {
	if collectorID == targetUserID {
		return errors.New("cannot trade your own card")
	}

	if quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}

	// Check Level 5 Requirement
	collector, err := s.userRepo.FindByID(ctx, collectorID)
	if err != nil {
		return errors.New("user not found")
	}
	if collector.Level < 5 {
		return errors.New("level 5 required to access trading")
	}

	// Check Circuit Breaker
	if s.redisClient != nil {
		val, _ := s.redisClient.Get(ctx, fmt.Sprintf("circuit_breaker:%s", targetUserID.String())).Result()
		if val == "1" {
			return errors.New("trading is currently halted due to high volatility (circuit breaker)")
		}
	}

	return s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		profile, err := s.cardRepo.GetCardProfileForUpdate(txCtx, targetUserID)
		if err != nil {
			// Create if not exists
			profile = &models.CardProfile{
				UserID:         targetUserID,
				CurrentPrice:   100,
				TotalCards:     1000,
				AvailableCards: 1000,
			}
			if err := s.cardRepo.CreateCardProfile(txCtx, profile); err != nil {
				return err
			}
		}

		oldPrice := profile.CurrentPrice
		totalCost := float64(0)
		feePercent := 0.02

		// Calculate cost using arithmetic series formula: Sum = n*(2a + (n-1)*d) / 2
		// where a = starting price, d = price step (5.0), n = quantity
		// This is O(1) instead of O(N) for large quantities
		if tradeType == "buy" {
			if profile.AvailableCards < quantity {
				return errors.New("not enough available cards")
			}
			
			// Starting price for the first card purchased
			a := 100.0 + float64(1000-(profile.AvailableCards))*5.0
			d := 5.0
			n := float64(quantity)
			// Arithmetic series: Sum = n * (2a + (n-1)*d) / 2
			totalCost = n * (2*a + (n-1)*d) / 2

			// Add 2% fee
			totalCostWithFee := totalCost * (1.0 + feePercent)

			// Deduct from wallet
			if err := s.walletRepo.UpdateBalance(txCtx, collectorID, -totalCostWithFee); err != nil {
				return err
			}

			profile.AvailableCards -= quantity
			profile.CurrentPrice = 100.0 + float64(1000-profile.AvailableCards)*5.0

		} else if tradeType == "sell" {
			owned, err := s.cardRepo.GetOwnedQuantity(txCtx, collectorID, targetUserID)
			if err != nil {
				return err
			}
			if owned < quantity {
				return errors.New("not enough owned cards to sell")
			}

			// Starting price for the first card sold
			a := 100.0 + float64(1000-(profile.AvailableCards))*5.0
			d := -5.0
			n := float64(quantity)

			// Calculate how many cards hit the price floor
			// Price at step i: a + i*d >= 10 => i <= (a - 10) / 5
			stepsBeforeFloor := int((a - 10.0) / 5.0) + 1
			if stepsBeforeFloor < 0 {
				stepsBeforeFloor = 0
			}

			if stepsBeforeFloor >= quantity {
				// All cards are above floor price: use arithmetic series
				totalCost = n * (2*a + (n-1)*d) / 2
			} else {
				// Some cards hit the floor
				nAbove := float64(stepsBeforeFloor)
				if nAbove > 0 {
					totalCost = nAbove * (2*a + (nAbove-1)*d) / 2
				}
				// Remaining cards at floor price
				totalCost += float64(quantity-stepsBeforeFloor) * 10.0
			}

			// Deduct 2% fee
			totalRevenueAfterFee := totalCost * (1.0 - feePercent)

			// Add to wallet
			if err := s.walletRepo.UpdateBalance(txCtx, collectorID, totalRevenueAfterFee); err != nil {
				return err
			}

			profile.AvailableCards += quantity
			profile.CurrentPrice = 100.0 + float64(1000-profile.AvailableCards)*5.0
			if profile.CurrentPrice < 10 {
				profile.CurrentPrice = 10
			}

		} else {
			return errors.New("invalid trade type")
		}

		// Update profile
		if err := s.cardRepo.UpdateCardProfile(txCtx, profile); err != nil {
			return err
		}

		// Create transaction history
		txType := "card_buy"
		amount := -totalCost
		if tradeType == "sell" {
			txType = "card_sell"
			amount = totalCost
		}

		txHistory := &models.CardTransaction{
			CollectorID:        collectorID,
			TargetUserID:       targetUserID,
			Type:               tradeType,
			Quantity:           quantity,
			PriceAtTransaction: totalCost / float64(quantity), // Average price
		}
		if err := s.cardRepo.CreateCardTransaction(txCtx, txHistory); err != nil {
			return err
		}

		// Create wallet transaction
		walletTx := &models.WalletTransaction{
			UserID:      collectorID,
			Amount:      amount,
			Type:        txType,
			ReferenceID: txHistory.ID,
		}
		if err := s.walletRepo.CreateTransaction(txCtx, walletTx); err != nil {
			return err
		}

		// Circuit Breaker Check
		priceChangeRatio := (profile.CurrentPrice - oldPrice) / oldPrice
		if priceChangeRatio < 0 {
			priceChangeRatio = -priceChangeRatio
		}

		if priceChangeRatio > 0.30 && s.redisClient != nil {
			// Halt trading for 5 minutes
			s.redisClient.Set(ctx, fmt.Sprintf("circuit_breaker:%s", targetUserID.String()), "1", 5*time.Minute)
		}

		return nil
	}, &sql.TxOptions{Isolation: sql.LevelSerializable})
}
