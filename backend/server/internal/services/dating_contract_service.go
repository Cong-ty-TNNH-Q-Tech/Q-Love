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

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/repository"
	"github.com/google/uuid"
)

type DatingContractService interface {
	CreateContract(ctx context.Context, userID uuid.UUID, targetUserID uuid.UUID, amount float64, appointmentTime time.Time) (*models.DatingContract, error)
	AcceptContract(ctx context.Context, contractID uuid.UUID, userID uuid.UUID) (*models.DatingContract, error)
	CancelContract(ctx context.Context, contractID uuid.UUID, userID uuid.UUID, reason string) error
	ScanContract(ctx context.Context, contractID uuid.UUID, qrToken string) error
}

type datingContractService struct {
	contractRepo repository.DatingContractRepository
	walletRepo   repository.WalletRepository
	matchRepo    repository.MatchRepository
	chatRepo     repository.ChatRepository
	premiumRepo  repository.UserPremiumRepository
	txManager    repository.TransactionManager
}

func NewDatingContractService(
	contractRepo repository.DatingContractRepository,
	walletRepo repository.WalletRepository,
	matchRepo repository.MatchRepository,
	chatRepo repository.ChatRepository,
	premiumRepo repository.UserPremiumRepository,
	txManager repository.TransactionManager,
) DatingContractService {
	return &datingContractService{
		contractRepo: contractRepo,
		walletRepo:   walletRepo,
		matchRepo:    matchRepo,
		chatRepo:     chatRepo,
		premiumRepo:  premiumRepo,
		txManager:    txManager,
	}
}

func generateTOTPSecret() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (s *datingContractService) CreateContract(ctx context.Context, userID uuid.UUID, targetUserID uuid.UUID, amount float64, appointmentTime time.Time) (*models.DatingContract, error) {
	if amount <= 0 {
		return nil, errors.New("deposit amount must be greater than 0")
	}

	// 1. Verify match and chat count
	match, err := s.matchRepo.FindByUsers(ctx, userID, targetUserID)
	if err != nil {
		return nil, errors.New("match not found")
	}

	chatCount, err := s.chatRepo.CountMessagesByMatchID(ctx, match.ID)
	if err != nil || chatCount < 20 {
		return nil, errors.New("users must chat at least 20 messages before creating a dating contract")
	}

	var newContract *models.DatingContract

	// 2. Transaction with SERIALIZABLE isolation to avoid race conditions
	err = s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		wallet, err := s.walletRepo.GetWalletForUpdate(txCtx, userID)
		if err != nil {
			return err
		}

		if wallet.Balance < amount {
			return errors.New("insufficient balance")
		}

		// Hold balance
		if err := s.walletRepo.HoldBalance(txCtx, userID, amount); err != nil {
			return err
		}

		// Log transaction
		tx := &models.WalletTransaction{
			ID:          uuid.New(),
			UserID:      userID,
			Amount:      -amount,
			Type:        "contract_hold",
		}
		if err := s.walletRepo.CreateTransaction(txCtx, tx); err != nil {
			return err
		}

		// Create Contract
		secret := generateTOTPSecret()
		contract := &models.DatingContract{
			UserAID:         userID,
			UserBID:         targetUserID,
			DepositAmount:   amount,
			Status:          "pending",
			TOTPSecret:      secret,
			AppointmentTime: &appointmentTime,
		}

		if err := s.contractRepo.Create(txCtx, contract); err != nil {
			return err
		}
		newContract = contract
		return nil
	}, &sql.TxOptions{Isolation: sql.LevelSerializable})

	if err != nil {
		return nil, err
	}
	return newContract, nil
}

func (s *datingContractService) AcceptContract(ctx context.Context, contractID uuid.UUID, userID uuid.UUID) (*models.DatingContract, error) {
	var acceptedContract *models.DatingContract

	err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		contract, err := s.contractRepo.GetByIDForUpdate(txCtx, contractID)
		if err != nil {
			return err
		}

		if contract.UserBID != userID {
			return errors.New("forbidden: you are not the receiver of this contract")
		}
		if contract.Status != "pending" {
			return errors.New("contract is not in pending state")
		}

		wallet, err := s.walletRepo.GetWalletForUpdate(txCtx, userID)
		if err != nil {
			return err
		}

		if wallet.Balance < contract.DepositAmount {
			return errors.New("insufficient balance to accept the contract")
		}

		// Hold balance for User B
		if err := s.walletRepo.HoldBalance(txCtx, userID, contract.DepositAmount); err != nil {
			return err
		}

		// Log transaction
		tx := &models.WalletTransaction{
			ID:          uuid.New(),
			UserID:      userID,
			Amount:      -contract.DepositAmount,
			Type:        "contract_hold",
			ReferenceID: contract.ID,
		}
		if err := s.walletRepo.CreateTransaction(txCtx, tx); err != nil {
			return err
		}

		contract.Status = "active"
		if err := s.contractRepo.Update(txCtx, contract); err != nil {
			return err
		}
		
		acceptedContract = contract
		return nil
	}, &sql.TxOptions{Isolation: sql.LevelSerializable})

	if err != nil {
		return nil, err
	}
	return acceptedContract, nil
}

func (s *datingContractService) CancelContract(ctx context.Context, contractID uuid.UUID, userID uuid.UUID, reason string) error {
	return s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		contract, err := s.contractRepo.GetByIDForUpdate(txCtx, contractID)
		if err != nil {
			return err
		}

		if contract.Status != "active" && contract.Status != "pending" {
			return errors.New("contract cannot be cancelled in its current state")
		}

		if contract.UserAID != userID && contract.UserBID != userID {
			return errors.New("forbidden")
		}

		// Save original status BEFORE overwriting
		originalStatus := contract.Status
		contract.Status = "cancelled"
		contract.CancelledByID = &userID
		
		// If was pending, just refund User A (since User B hasn't deposited yet)
		if originalStatus == "pending" {
			if err := s.refundDeposit(txCtx, contract.UserAID, contract.DepositAmount, contract.ID); err != nil {
				return err
			}
			return s.contractRepo.Update(txCtx, contract)
		}
		
		// If active, check for premium free cancellation
		premium, err := s.premiumRepo.FindByUserID(txCtx, userID)
		if err == nil && premium != nil && premium.ExpiresAt.After(time.Now()) && premium.FreeCancelLeft > 0 {
			// Free cancellation! Refund both users
			if err := s.refundDeposit(txCtx, contract.UserAID, contract.DepositAmount, contract.ID); err != nil {
				return err
			}
			if err := s.refundDeposit(txCtx, contract.UserBID, contract.DepositAmount, contract.ID); err != nil {
				return err
			}
			
			premium.FreeCancelLeft -= 1
			if err := s.premiumRepo.Update(txCtx, premium); err != nil {
				return err
			}
		} else {
			// Confiscate from the cancelling user
			victimID := contract.UserAID
			if contract.UserAID == userID {
				victimID = contract.UserBID
			}
			
			// Deduct hold from canceler, and add 0 to balance (confiscated)
			if err := s.walletRepo.ReleaseHoldBalance(txCtx, userID, contract.DepositAmount); err != nil {
				return err
			}
			if err := s.walletRepo.UpdateBalance(txCtx, userID, -contract.DepositAmount); err != nil {
				return err
			}

			// Add 90% to victim
			rewardAmount := contract.DepositAmount * 0.90
			if err := s.walletRepo.UpdateBalance(txCtx, victimID, rewardAmount); err != nil {
				return err
			}
			
			// Refund victim's own deposit
			if err := s.refundDeposit(txCtx, victimID, contract.DepositAmount, contract.ID); err != nil {
				return err
			}
		}

		return s.contractRepo.Update(txCtx, contract)
	}, &sql.TxOptions{Isolation: sql.LevelSerializable})
}

func (s *datingContractService) ScanContract(ctx context.Context, contractID uuid.UUID, qrToken string) error {
	// Simple TOTP logic check (In real implementation, we should use a proper TOTP package)
	// For this simulation, we assume any token equal to TOTPSecret is valid
	return s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		contract, err := s.contractRepo.GetByIDForUpdate(txCtx, contractID)
		if err != nil {
			return err
		}

		if contract.Status != "active" {
			return errors.New("contract must be active to scan")
		}

		// Mock TOTP validation
		if qrToken != contract.TOTPSecret {
			return errors.New("invalid QR token")
		}

		contract.Status = "completed"
		if err := s.contractRepo.Update(txCtx, contract); err != nil {
			return err
		}

		// Refund both
		if err := s.refundDeposit(txCtx, contract.UserAID, contract.DepositAmount, contract.ID); err != nil {
			return err
		}
		if err := s.refundDeposit(txCtx, contract.UserBID, contract.DepositAmount, contract.ID); err != nil {
			return err
		}
		
		return nil
	}, &sql.TxOptions{Isolation: sql.LevelSerializable})
}

func (s *datingContractService) refundDeposit(ctx context.Context, userID uuid.UUID, amount float64, refID uuid.UUID) error {
	if err := s.walletRepo.ReleaseHoldBalance(ctx, userID, amount); err != nil {
		return err
	}
	
	tx := &models.WalletTransaction{
		ID:          uuid.New(),
		UserID:      userID,
		Amount:      amount,
		Type:        "contract_refund",
		ReferenceID: refID,
	}
	return s.walletRepo.CreateTransaction(ctx, tx)
}
