// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"database/sql"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
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

func generateTOTPSecret() (string, error) {
	bytes := make([]byte, 20)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate TOTP secret: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateTOTP generates a 6-digit TOTP code based on RFC 6238.
// It uses HMAC-SHA1 with a 30-second time step.
// Exported for testing purposes.
func GenerateTOTP(secret string, t time.Time) (string, error) {
	key, err := hex.DecodeString(secret)
	if err != nil {
		return "", fmt.Errorf("invalid TOTP secret: %w", err)
	}

	// Time step = 30 seconds (RFC 6238 default)
	counter := uint64(t.Unix()) / 30

	// Convert counter to big-endian bytes
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, counter)

	// HMAC-SHA1
	mac := hmac.New(sha1.New, key)
	mac.Write(buf)
	hash := mac.Sum(nil)

	// Dynamic truncation (RFC 4226 Section 5.4)
	offset := hash[len(hash)-1] & 0x0F
	code := binary.BigEndian.Uint32(hash[offset:offset+4]) & 0x7FFFFFFF
	otp := code % uint32(math.Pow10(6))

	return fmt.Sprintf("%06d", otp), nil
}

// validateTOTP checks the provided token against TOTP codes within a
// +/- 1 time step window to account for clock skew.
func validateTOTP(secret string, token string) bool {
	now := time.Now()
	for _, offset := range []time.Duration{-30 * time.Second, 0, 30 * time.Second} {
		expected, err := GenerateTOTP(secret, now.Add(offset))
		if err != nil {
			continue
		}
		if hmac.Equal([]byte(expected), []byte(token)) {
			return true
		}
	}
	return false
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
		secret, err := generateTOTPSecret()
		if err != nil {
			return err
		}
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

		contract.Status = "cancelled"
		contract.CancelledByID = &userID
		
		// If pending, just refund User A (since User B hasn't deposited yet)
		if contract.Status == "pending" {
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
	// RFC 6238 TOTP validation with +/- 30s window for clock skew tolerance
	return s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		contract, err := s.contractRepo.GetByIDForUpdate(txCtx, contractID)
		if err != nil {
			return err
		}

		if contract.Status != "active" {
			return errors.New("contract must be active to scan")
		}

		// Validate TOTP token (RFC 6238 with +/- 1 time step window)
		if !validateTOTP(contract.TOTPSecret, qrToken) {
			return errors.New("invalid or expired QR token")
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
