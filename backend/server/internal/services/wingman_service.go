// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/repository"
)

type WingmanService interface {
	CreateReferral(ctx context.Context, wingmanID, target1ID, target2ID uuid.UUID) (*models.WingmanReferral, error)
	AcceptReferral(ctx context.Context, referralID, acceptingUserID uuid.UUID) (*models.WingmanReferral, error)
	ProcessCommission(ctx context.Context, referralID uuid.UUID) error
}

type wingmanService struct {
	wingmanRepo repository.WingmanRepository
	walletRepo  repository.WalletRepository
	txManager   repository.TransactionManager
	matchRepo   repository.MatchRepository
}

func NewWingmanService(
	wingmanRepo repository.WingmanRepository, 
	walletRepo repository.WalletRepository, 
	txManager repository.TransactionManager,
	matchRepo repository.MatchRepository,
) WingmanService {
	return &wingmanService{
		wingmanRepo: wingmanRepo,
		walletRepo:  walletRepo,
		txManager:   txManager,
		matchRepo:   matchRepo,
	}
}

func (s *wingmanService) CreateReferral(ctx context.Context, wingmanID, target1ID, target2ID uuid.UUID) (*models.WingmanReferral, error) {
	// Generate a secure deep link token
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}
	token := hex.EncodeToString(bytes)
	deepLink := "qlove://match/" + token

	referral := &models.WingmanReferral{
		ID:        uuid.New(),
		WingmanID: wingmanID,
		Target1ID: target1ID,
		Target2ID: target2ID,
		Status:    "pending",
		DeepLink:  deepLink,
		ExpiresAt: time.Now().Add(48 * time.Hour), // Link expires in 48 hours
	}

	if err := s.wingmanRepo.CreateReferral(ctx, referral); err != nil {
		return nil, err
	}
	return referral, nil
}

func (s *wingmanService) AcceptReferral(ctx context.Context, referralID, acceptingUserID uuid.UUID) (*models.WingmanReferral, error) {
	var referral *models.WingmanReferral

	err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		referral, err = s.wingmanRepo.GetReferralByID(txCtx, referralID)
		if err != nil {
			return errors.New("referral not found")
		}

		if referral.Status != "pending" {
			return errors.New("referral is no longer pending")
		}

		if time.Now().After(referral.ExpiresAt) {
			return errors.New("referral link expired")
		}

		if acceptingUserID != referral.Target1ID && acceptingUserID != referral.Target2ID {
			return errors.New("user is not part of this referral")
		}

		if referral.Target1ID == referral.WingmanID || referral.Target2ID == referral.WingmanID {
			return errors.New("wingman cannot refer themselves")
		}

		// Tao Match thuc su
		match := &models.Match{
			ID:      uuid.New(),
			User1ID: referral.Target1ID,
			User2ID: referral.Target2ID,
		}
		if err := s.matchRepo.Create(txCtx, match); err != nil {
			return err
		}

		referral.Status = "matched"
		referral.MatchID = &match.ID
		if err := s.wingmanRepo.UpdateReferral(txCtx, referral); err != nil {
			return err
		}
		
		return nil
	})

	if err != nil {
		return nil, err
	}
	return referral, nil
}

func (s *wingmanService) ProcessCommission(ctx context.Context, referralID uuid.UUID) error {
	// Using SERIALIZABLE transaction to prevent race conditions on wallet
	return s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		referral, err := s.wingmanRepo.GetReferralByID(txCtx, referralID)
		if err != nil {
			return err
		}

		if referral.Status != "matched" && referral.Status != "dated" {
			return errors.New("invalid status for commission")
		}

		// Reward Wingman 10% (Assume 10 Xu for now)
		commissionAmount := 10.0

		// Update Wingman Wallet
		if err := s.walletRepo.AddCommission(txCtx, referral.WingmanID, commissionAmount); err != nil {
			return err
		}

		// Log Transaction
		txn := &models.WalletTransaction{
			ID:          uuid.New(),
			UserID:      referral.WingmanID,
			Amount:      commissionAmount,
			Type:        "wingman_commission",
			ReferenceID: referral.ID,
		}
		if err := s.walletRepo.CreateTransaction(txCtx, txn); err != nil {
			return err
		}

		// Mark referral as rewarded
		referral.Status = "rewarded"
		return s.wingmanRepo.UpdateReferral(txCtx, referral)
	}, &sql.TxOptions{Isolation: sql.LevelSerializable}) // Ensure strong consistency
}
