// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/repository"
)

type ShameService interface {
	GetActiveShames(ctx context.Context, limit, offset int) ([]models.WallOfShame, error)
	ThrowTomato(ctx context.Context, throwerID uuid.UUID, shameID uuid.UUID) error
}

type shameService struct {
	shameRepo  repository.ShameRepository
	walletRepo repository.WalletRepository
	txManager  repository.TransactionManager
}

func NewShameService(
	shameRepo repository.ShameRepository,
	walletRepo repository.WalletRepository,
	txManager repository.TransactionManager,
) ShameService {
	return &shameService{
		shameRepo:  shameRepo,
		walletRepo: walletRepo,
		txManager:  txManager,
	}
}

func (s *shameService) GetActiveShames(ctx context.Context, limit, offset int) ([]models.WallOfShame, error) {
	return s.shameRepo.GetActiveShames(ctx, limit, offset)
}

func (s *shameService) ThrowTomato(ctx context.Context, throwerID uuid.UUID, shameID uuid.UUID) error {
	return s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		// 1. Get wallet and lock it for update
		wallet, err := s.walletRepo.GetWalletForUpdate(txCtx, throwerID)
		if err != nil {
			return err
		}

		// 2. Check balance (Tomato cost = 1 Xu)
		cost := 1.0
		if wallet.Balance < cost {
			return errors.New("insufficient balance to throw a tomato")
		}

		// 3. Deduct balance
		err = s.walletRepo.UpdateBalance(txCtx, throwerID, -cost)
		if err != nil {
			return err
		}

		// 4. Create wallet transaction record
		walletTx := &models.WalletTransaction{
			ID:          uuid.New(),
			UserID:      throwerID,
			Amount:      -cost,
			Type:        "THROW_TOMATO",
			ReferenceID: shameID,
		}
		err = s.walletRepo.CreateTransaction(txCtx, walletTx)
		if err != nil {
			return err
		}

		// 5. Increment tomato count on wall of shame
		err = s.shameRepo.IncrementTomatoCount(txCtx, shameID, 1)
		if err != nil {
			return err
		}

		return nil
	})
}
