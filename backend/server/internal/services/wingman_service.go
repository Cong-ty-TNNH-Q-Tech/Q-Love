package services

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
)

type WingmanService interface {
	CreateReferral(ctx context.Context, wingmanID, target1ID, target2ID uuid.UUID) (*models.WingmanReferral, error)
	AcceptReferral(ctx context.Context, referralID, acceptingUserID uuid.UUID) (*models.WingmanReferral, error)
	ProcessCommission(ctx context.Context, referralID uuid.UUID) error
}

type wingmanService struct {
	db *gorm.DB
}

func NewWingmanService(db *gorm.DB) WingmanService {
	return &wingmanService{db: db}
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

	if err := s.db.WithContext(ctx).Create(referral).Error; err != nil {
		return nil, err
	}
	return referral, nil
}

func (s *wingmanService) AcceptReferral(ctx context.Context, referralID, acceptingUserID uuid.UUID) (*models.WingmanReferral, error) {
	var referral models.WingmanReferral

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&referral, "id = ?", referralID).Error; err != nil {
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

		// For simplicity, we assume one person clicking the link accepts it and creates a match.
		// In a real scenario, we might need both to accept. Let's mark it as matched.
		referral.Status = "matched"
		if err := tx.Save(&referral).Error; err != nil {
			return err
		}
		
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &referral, nil
}

func (s *wingmanService) ProcessCommission(ctx context.Context, referralID uuid.UUID) error {
	// Using SERIALIZABLE transaction to prevent race conditions on wallet
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var referral models.WingmanReferral
		if err := tx.First(&referral, "id = ?", referralID).Error; err != nil {
			return err
		}

		if referral.Status != "matched" && referral.Status != "dated" {
			return errors.New("invalid status for commission")
		}

		// Reward Wingman 10% (Assume 10 Xu for now)
		commissionAmount := 10.0

		// Update Wingman Wallet
		var wallet models.UserWallet
		if err := tx.FirstOrCreate(&wallet, models.UserWallet{UserID: referral.WingmanID}).Error; err != nil {
			return err
		}

		wallet.Balance += commissionAmount
		if err := tx.Save(&wallet).Error; err != nil {
			return err
		}

		// Log Transaction
		txn := models.WalletTransaction{
			ID:          uuid.New(),
			UserID:      referral.WingmanID,
			Amount:      commissionAmount,
			Type:        "wingman_commission",
			ReferenceID: referral.ID,
		}
		if err := tx.Create(&txn).Error; err != nil {
			return err
		}

		// Mark referral as rewarded
		referral.Status = "rewarded"
		return tx.Save(&referral).Error
	}, &sql.TxOptions{Isolation: sql.LevelSerializable}) // Ensure strong consistency
}
