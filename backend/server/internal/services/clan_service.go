// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"context"
	"errors"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClanService interface {
	CreateClan(ctx context.Context, userID uuid.UUID, name, slogan, logoURL string) (*models.Clan, error)
}

type clanService struct {
	clanRepo   repository.ClanRepository
	walletRepo repository.WalletRepository
	txManager  repository.TransactionManager
}

func NewClanService(clanRepo repository.ClanRepository, walletRepo repository.WalletRepository, txManager repository.TransactionManager) ClanService {
	return &clanService{
		clanRepo:   clanRepo,
		walletRepo: walletRepo,
		txManager:  txManager,
	}
}

func (s *clanService) CreateClan(ctx context.Context, userID uuid.UUID, name, slogan, logoURL string) (*models.Clan, error) {
	// 1. Validate clan name (must not exist)
	existingClan, err := s.clanRepo.FindByName(ctx, name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existingClan != nil {
		return nil, errors.New("ERR_CLAN_NAME_TAKEN")
	}

	var createdClan *models.Clan

	// 2. Start Transaction
	err = s.txManager.ExecuteInTx(ctx, func(tx *gorm.DB) error {
		// 3. Deduct 500 Xu from user wallet
		wallet, err := s.walletRepo.GetByUserID(ctx, tx, userID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("wallet not found")
			}
			return err
		}

		if wallet.Balance < 500 {
			return errors.New("insufficient balance")
		}

		err = s.walletRepo.UpdateBalance(ctx, tx, wallet.ID, -500)
		if err != nil {
			return err
		}

		// Create transaction record
		txRecord := &models.WalletTransaction{
			WalletID: wallet.ID,
			Amount:   -500,
			Type:     "clan_create",
		}
		err = s.walletRepo.CreateTransaction(ctx, tx, txRecord)
		if err != nil {
			return err
		}

		// 4. Create Clan
		clan := &models.Clan{
			Name:     name,
			Slogan:   slogan,
			LogoURL:  logoURL,
			LeaderID: userID,
		}
		if err := s.clanRepo.CreateClan(ctx, tx, clan); err != nil {
			return err
		}
		createdClan = clan

		// 5. Add user to ClanMembers as leader
		member := &models.ClanMember{
			ClanID: clan.ID,
			UserID: userID,
			Role:   "leader",
		}
		if err := s.clanRepo.AddClanMember(ctx, tx, member); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return createdClan, nil
}
