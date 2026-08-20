// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.
package services

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/repository"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrInvalidWebhookPayload = errors.New("invalid webhook payload")
	ErrTransactionExists     = errors.New("transaction already processed")
)

// Coin mapping (mocking a config/db dictionary)
var CoinPacks = map[string]float64{
	"coin_pack_100":  100.0,
	"coin_pack_500":  500.0,
	"coin_pack_1000": 1000.0,
}

type RevenueCatEvent struct {
	Type           string `json:"type"`
	AppUserID      string `json:"app_user_id"`
	ProductID      string `json:"product_id"`
	TransactionID  string `json:"transaction_id"`
	ExpirationAtMs int64  `json:"expiration_at_ms,omitempty"`
}

type IAPService interface {
	ProcessRevenueCatWebhook(ctx context.Context, event RevenueCatEvent) error
}

type iapService struct {
	txManager    repository.TransactionManager
	walletRepo   repository.WalletRepository
	userPremRepo repository.UserPremiumRepository
}

func NewIAPService(txManager repository.TransactionManager, walletRepo repository.WalletRepository, userPremRepo repository.UserPremiumRepository) IAPService {
	return &iapService{
		txManager:    txManager,
		walletRepo:   walletRepo,
		userPremRepo: userPremRepo,
	}
}

func (s *iapService) ProcessRevenueCatWebhook(ctx context.Context, event RevenueCatEvent) error {
	userID, err := uuid.Parse(event.AppUserID)
	if err != nil {
		logger.Log.Error("Invalid app_user_id in webhook", zap.String("app_user_id", event.AppUserID))
		return ErrInvalidWebhookPayload
	}

	txID, err := uuid.Parse(event.TransactionID)
	if err != nil {
		// RevenueCat transaction ID might not be a UUID, but our DB expects it as UUID?
		// Wait, RevenueCat sends string transaction_id like "123456789".
		// We can hash it to a UUID, or use string.
		// For now, assume it's UUID or we generate one based on it for idempotency.
		// Let's use uuid.NewSHA1 with a namespace to generate a deterministic UUID.
		namespace := uuid.MustParse("00000000-0000-0000-0000-000000000000")
		txID = uuid.NewSHA1(namespace, []byte(event.TransactionID))
	}

	// Xử lý nạp Xu (NON_RENEWING_PURCHASE)
	if event.Type == "NON_RENEWING_PURCHASE" {
		coins, ok := CoinPacks[event.ProductID]
		if !ok {
			logger.Log.Warn("Unknown product_id for NON_RENEWING_PURCHASE", zap.String("product_id", event.ProductID))
			return nil // Return 200 OK to stop retrying
		}

		err = s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
			// 1. Check idempotency: Ensure the transaction ID hasn't been processed
			exists, err := s.walletRepo.CheckTransactionExists(txCtx, txID)
			if err != nil {
				return err
			}
			if exists {
				return ErrTransactionExists
			}

			// 2. Process deposit
			err = s.walletRepo.UpdateBalance(txCtx, userID, coins)
			if err != nil {
				return err
			}

			// 3. Record transaction
			return s.walletRepo.CreateTransaction(txCtx, &models.WalletTransaction{
				ID:        txID,
				UserID:    userID,
				Amount:    coins,
				Type:      "iap_deposit",
				CreatedAt: time.Now(),
			})
		}, &sql.TxOptions{Isolation: sql.LevelSerializable})

		if err != nil {
			if errors.Is(err, ErrTransactionExists) {
				logger.Log.Info("IAP transaction already processed", zap.String("tx_id", event.TransactionID))
				return nil
			}
			logger.Log.Error("Failed to process IAP deposit", zap.Error(err))
			return err
		}
		
		logger.Log.Info("Successfully deposited coins", zap.String("user_id", userID.String()), zap.Float64("coins", coins))
		return nil
	}

	// Xử lý gói Premium (INITIAL_PURCHASE, RENEWAL)
	if event.Type == "INITIAL_PURCHASE" || event.Type == "RENEWAL" {
		if event.ProductID != "qlove_premium_1month" && event.ProductID != "qlove_premium_1year" {
			return nil
		}

		err = s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
			// 1. Check idempotency: Ensure the transaction ID hasn't been processed
			exists, err := s.walletRepo.CheckTransactionExists(txCtx, txID)
			if err != nil {
				return err
			}
			if exists {
				return ErrTransactionExists
			}

			// 2. Activate Premium
			var expiresAt time.Time
			if event.ExpirationAtMs > 0 {
				expiresAt = time.UnixMilli(event.ExpirationAtMs)
			} else {
				if event.ProductID == "qlove_premium_1year" {
					expiresAt = time.Now().AddDate(1, 0, 0) // Default 1 year
				} else {
					expiresAt = time.Now().AddDate(0, 1, 0) // Default 1 month
				}
			}
			err = s.userPremRepo.ActivatePremium(txCtx, userID, expiresAt)
			if err != nil {
				return err
			}

			// 3. Record transaction
			return s.walletRepo.CreateTransaction(txCtx, &models.WalletTransaction{
				ID:        txID,
				UserID:    userID,
				Amount:    0,
				Type:      "iap_premium",
				CreatedAt: time.Now(),
			})
		}, &sql.TxOptions{Isolation: sql.LevelSerializable})

		if err != nil {
			if errors.Is(err, ErrTransactionExists) {
				logger.Log.Info("IAP premium transaction already processed", zap.String("tx_id", event.TransactionID))
				return nil
			}
			logger.Log.Error("Failed to activate premium", zap.Error(err))
			return err
		}
		
		logger.Log.Info("Successfully activated premium", zap.String("user_id", userID.String()))
		return nil
	}

	// Ignore other types
	return nil
}

