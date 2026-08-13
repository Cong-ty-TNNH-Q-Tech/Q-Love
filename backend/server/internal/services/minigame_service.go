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
	"github.com/google/uuid"
)

type MinigameService interface {
	InitSteal(ctx context.Context, attackerID uuid.UUID, defenderID uuid.UUID, targetCardID uuid.UUID) (*models.CardSteal, error)
	SubmitStealResult(ctx context.Context, stealID uuid.UUID, attackerID uuid.UUID, isWin bool) error
}

type minigameService struct {
	stealRepo  repository.CardStealRepository
	walletRepo repository.WalletRepository
	txManager  repository.TransactionManager
}

func NewMinigameService(
	stealRepo repository.CardStealRepository,
	walletRepo repository.WalletRepository,
	txManager repository.TransactionManager,
) MinigameService {
	return &minigameService{
		stealRepo:  stealRepo,
		walletRepo: walletRepo,
		txManager:  txManager,
	}
}

func (s *minigameService) InitSteal(ctx context.Context, attackerID uuid.UUID, defenderID uuid.UUID, targetCardID uuid.UUID) (*models.CardSteal, error) {
	var steal *models.CardSteal
	err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		// 1. Lock wallet
		wallet, err := s.walletRepo.GetWalletForUpdate(txCtx, attackerID)
		if err != nil {
			return err
		}

		// 2. Check balance (Ticket price = 1000 Xu)
		cost := 1000.0
		if wallet.Balance < cost {
			return errors.New("insufficient balance to buy Thẻ Đạo Tặc")
		}

		// 3. Deduct balance
		if err := s.walletRepo.UpdateBalance(txCtx, attackerID, -cost); err != nil {
			return err
		}

		// 4. Create wallet transaction record
		walletTx := &models.WalletTransaction{
			ID:          uuid.New(),
			UserID:      attackerID,
			Amount:      -cost,
			Type:        "BUY_STEAL_TICKET",
		}
		if err := s.walletRepo.CreateTransaction(txCtx, walletTx); err != nil {
			return err
		}

		// 5. Create CardSteal record
		steal = &models.CardSteal{
			ID:           uuid.New(),
			AttackerID:   attackerID,
			DefenderID:   defenderID,
			TargetCardID: targetCardID,
			Result:       "pending",
		}
		if err := s.stealRepo.Create(txCtx, steal); err != nil {
			return err
		}

		return nil
	}, &sql.TxOptions{Isolation: sql.LevelSerializable})

	return steal, err
}

func (s *minigameService) SubmitStealResult(ctx context.Context, stealID uuid.UUID, attackerID uuid.UUID, isWin bool) error {
	return s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		// 1. Find steal record
		steal, err := s.stealRepo.FindByID(txCtx, stealID)
		if err != nil {
			return errors.New("steal session not found")
		}
		if steal.AttackerID != attackerID {
			return errors.New("unauthorized to submit result for this steal session")
		}
		if steal.Result != "pending" {
			return errors.New("steal session already completed")
		}

		// 2. Validate Time (Anti-cheat: must be around 10 seconds, allow 8-20 seconds for network latency)
		elapsed := time.Since(steal.CreatedAt)
		if elapsed < 8*time.Second || elapsed > 30*time.Second {
			// If outside bounds, auto-fail it to prevent cheating
			isWin = false
		}

		// 3. Process Result
		if isWin {
			// Attacker Won -> Transfer Card Ownership
			if err := s.stealRepo.UpdateResult(txCtx, stealID, "attacker_won"); err != nil {
				return err
			}
			if err := s.stealRepo.TransferCardOwnership(txCtx, attackerID, steal.TargetCardID); err != nil {
				return err
			}
		} else {
			// Attacker Lost -> Compensate 500 Xu to Attacker
			if err := s.stealRepo.UpdateResult(txCtx, stealID, "defender_won"); err != nil {
				return err
			}
			
			// Refund/compensate 500 Xu
			compensation := 500.0
			if err := s.walletRepo.UpdateBalance(txCtx, attackerID, compensation); err != nil {
				return err
			}
			walletTx := &models.WalletTransaction{
				ID:          uuid.New(),
				UserID:      attackerID,
				Amount:      compensation,
				Type:        "STEAL_COMPENSATION",
				ReferenceID: stealID,
			}
			if err := s.walletRepo.CreateTransaction(txCtx, walletTx); err != nil {
				return err
			}
		}

		return nil
	}, &sql.TxOptions{Isolation: sql.LevelSerializable})
}
