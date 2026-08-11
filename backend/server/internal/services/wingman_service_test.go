// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
)

// --- Mocks ---

type mockWingmanRepo struct {
	createReferralFn  func(ctx context.Context, referral *models.WingmanReferral) error
	getReferralByIDFn func(ctx context.Context, referralID uuid.UUID) (*models.WingmanReferral, error)
	updateReferralFn  func(ctx context.Context, referral *models.WingmanReferral) error
}

func (m *mockWingmanRepo) CreateReferral(ctx context.Context, referral *models.WingmanReferral) error {
	return m.createReferralFn(ctx, referral)
}
func (m *mockWingmanRepo) GetReferralByID(ctx context.Context, referralID uuid.UUID) (*models.WingmanReferral, error) {
	return m.getReferralByIDFn(ctx, referralID)
}
func (m *mockWingmanRepo) UpdateReferral(ctx context.Context, referral *models.WingmanReferral) error {
	return m.updateReferralFn(ctx, referral)
}

type mockWalletRepo struct {
	addCommissionFn     func(ctx context.Context, userID uuid.UUID, amount float64) error
	createTransactionFn func(ctx context.Context, txn *models.WalletTransaction) error
}

func (m *mockWalletRepo) AddCommission(ctx context.Context, userID uuid.UUID, amount float64) error {
	return m.addCommissionFn(ctx, userID, amount)
}
func (m *mockWalletRepo) CreateTransaction(ctx context.Context, txn *models.WalletTransaction) error {
	return m.createTransactionFn(ctx, txn)
}

type mockTxManager struct{}

func (m *mockTxManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error, opts ...*sql.TxOptions) error {
	// Simply execute the function directly, bypassing actual DB transactions
	return fn(ctx)
}

// --- Tests ---

func TestWingmanService_CreateReferral(t *testing.T) {
	wingmanRepo := &mockWingmanRepo{
		createReferralFn: func(ctx context.Context, referral *models.WingmanReferral) error {
			return nil
		},
	}
	service := NewWingmanService(wingmanRepo, nil, nil)

	wingmanID := uuid.New()
	target1ID := uuid.New()
	target2ID := uuid.New()

	referral, err := service.CreateReferral(context.Background(), wingmanID, target1ID, target2ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if referral == nil {
		t.Fatal("Expected referral to be returned")
	}
	if referral.WingmanID != wingmanID {
		t.Errorf("Expected wingmanID %s, got %s", wingmanID, referral.WingmanID)
	}
}

func TestWingmanService_AcceptReferral_Success(t *testing.T) {
	refID := uuid.New()
	target1ID := uuid.New()

	wingmanRepo := &mockWingmanRepo{
		getReferralByIDFn: func(ctx context.Context, referralID uuid.UUID) (*models.WingmanReferral, error) {
			return &models.WingmanReferral{
				ID:        refID,
				Status:    "pending",
				ExpiresAt: time.Now().Add(1 * time.Hour),
				Target1ID: target1ID,
				Target2ID: uuid.New(),
			}, nil
		},
		updateReferralFn: func(ctx context.Context, referral *models.WingmanReferral) error {
			return nil
		},
	}
	service := NewWingmanService(wingmanRepo, nil, &mockTxManager{})

	referral, err := service.AcceptReferral(context.Background(), refID, target1ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if referral.Status != "matched" {
		t.Errorf("Expected status matched, got %s", referral.Status)
	}
}

func TestWingmanService_ProcessCommission_Success(t *testing.T) {
	refID := uuid.New()
	wingmanID := uuid.New()

	wingmanRepo := &mockWingmanRepo{
		getReferralByIDFn: func(ctx context.Context, referralID uuid.UUID) (*models.WingmanReferral, error) {
			return &models.WingmanReferral{
				ID:        refID,
				Status:    "matched",
				WingmanID: wingmanID,
			}, nil
		},
		updateReferralFn: func(ctx context.Context, referral *models.WingmanReferral) error {
			return nil
		},
	}
	walletRepo := &mockWalletRepo{
		addCommissionFn: func(ctx context.Context, userID uuid.UUID, amount float64) error {
			return nil
		},
		createTransactionFn: func(ctx context.Context, txn *models.WalletTransaction) error {
			return nil
		},
	}
	service := NewWingmanService(wingmanRepo, walletRepo, &mockTxManager{})

	err := service.ProcessCommission(context.Background(), refID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestWingmanService_AcceptReferral_Errors(t *testing.T) {
	refID := uuid.New()
	target1ID := uuid.New()
	invalidUserID := uuid.New()

	wingmanRepo := &mockWingmanRepo{
		getReferralByIDFn: func(ctx context.Context, referralID uuid.UUID) (*models.WingmanReferral, error) {
			return nil, errors.New("not found")
		},
	}
	service := NewWingmanService(wingmanRepo, nil, &mockTxManager{})

	// 1. Referral not found
	_, err := service.AcceptReferral(context.Background(), refID, target1ID)
	if err == nil || err.Error() != "referral not found" {
		t.Errorf("Expected referral not found error, got %v", err)
	}

	// 2. Referral not pending
	wingmanRepo.getReferralByIDFn = func(ctx context.Context, referralID uuid.UUID) (*models.WingmanReferral, error) {
		return &models.WingmanReferral{
			ID:     refID,
			Status: "matched",
		}, nil
	}
	_, err = service.AcceptReferral(context.Background(), refID, target1ID)
	if err == nil || err.Error() != "referral is no longer pending" {
		t.Errorf("Expected referral is no longer pending error, got %v", err)
	}

	// 3. Referral expired
	wingmanRepo.getReferralByIDFn = func(ctx context.Context, referralID uuid.UUID) (*models.WingmanReferral, error) {
		return &models.WingmanReferral{
			ID:        refID,
			Status:    "pending",
			ExpiresAt: time.Now().Add(-1 * time.Hour),
		}, nil
	}
	_, err = service.AcceptReferral(context.Background(), refID, target1ID)
	if err == nil || err.Error() != "referral link expired" {
		t.Errorf("Expected referral link expired error, got %v", err)
	}

	// 4. Invalid user
	wingmanRepo.getReferralByIDFn = func(ctx context.Context, referralID uuid.UUID) (*models.WingmanReferral, error) {
		return &models.WingmanReferral{
			ID:        refID,
			Status:    "pending",
			ExpiresAt: time.Now().Add(1 * time.Hour),
			Target1ID: target1ID,
			Target2ID: uuid.New(),
		}, nil
	}
	_, err = service.AcceptReferral(context.Background(), refID, invalidUserID)
	if err == nil || err.Error() != "user is not part of this referral" {
		t.Errorf("Expected user is not part of this referral error, got %v", err)
	}
}

func TestWingmanService_ProcessCommission_Errors(t *testing.T) {
	refID := uuid.New()
	wingmanID := uuid.New()

	wingmanRepo := &mockWingmanRepo{
		getReferralByIDFn: func(ctx context.Context, referralID uuid.UUID) (*models.WingmanReferral, error) {
			return &models.WingmanReferral{
				ID:        refID,
				Status:    "pending", // invalid status for commission
				WingmanID: wingmanID,
			}, nil
		},
	}
	service := NewWingmanService(wingmanRepo, nil, &mockTxManager{})

	// 1. Invalid status
	err := service.ProcessCommission(context.Background(), refID)
	if err == nil || err.Error() != "invalid status for commission" {
		t.Errorf("Expected invalid status for commission error, got %v", err)
	}
}
