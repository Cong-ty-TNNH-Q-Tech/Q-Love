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
	getWalletForUpdateFn func(ctx context.Context, userID uuid.UUID) (*models.UserWallet, error)
	updateBalanceFn      func(ctx context.Context, userID uuid.UUID, delta float64) error
}

func (m *mockWalletRepo) AddCommission(ctx context.Context, userID uuid.UUID, amount float64) error {
	return m.addCommissionFn(ctx, userID, amount)
}
func (m *mockWalletRepo) CreateTransaction(ctx context.Context, txn *models.WalletTransaction) error {
	return m.createTransactionFn(ctx, txn)
}
func (m *mockWalletRepo) GetWalletForUpdate(ctx context.Context, userID uuid.UUID) (*models.UserWallet, error) {
	if m.getWalletForUpdateFn != nil {
		return m.getWalletForUpdateFn(ctx, userID)
	}
	return nil, nil
}
func (m *mockWalletRepo) UpdateBalance(ctx context.Context, userID uuid.UUID, delta float64) error {
	if m.updateBalanceFn != nil {
		return m.updateBalanceFn(ctx, userID, delta)
	}
	return nil
}
func (m *mockWalletRepo) CheckTransactionExists(ctx context.Context, txID uuid.UUID) (bool, error) {
	return false, nil
}

type mockTxManager struct{}

func (m *mockTxManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error, opts ...*sql.TxOptions) error {
	// Simply execute the function directly, bypassing actual DB transactions
	return fn(ctx)
}

type mockWingmanMatchRepo struct {
	createFn func(ctx context.Context, match *models.Match) error
}

func (m *mockWingmanMatchRepo) Create(ctx context.Context, match *models.Match) error {
	if m.createFn != nil {
		return m.createFn(ctx, match)
	}
	return nil
}
func (m *mockWingmanMatchRepo) FindByID(ctx context.Context, id uuid.UUID) (*models.Match, error) { return nil, nil }
func (m *mockWingmanMatchRepo) UpdateLastInteraction(ctx context.Context, id uuid.UUID, t time.Time) error { return nil }

// --- Tests ---

func TestWingmanService_CreateReferral(t *testing.T) {
	wingmanRepo := &mockWingmanRepo{
		createReferralFn: func(ctx context.Context, referral *models.WingmanReferral) error {
			return nil
		},
	}
	service := NewWingmanService(wingmanRepo, nil, nil, &mockWingmanMatchRepo{})

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
	service := NewWingmanService(wingmanRepo, nil, &mockTxManager{}, &mockWingmanMatchRepo{})

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
	service := NewWingmanService(wingmanRepo, walletRepo, &mockTxManager{}, &mockWingmanMatchRepo{})

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
	service := NewWingmanService(wingmanRepo, nil, &mockTxManager{}, &mockWingmanMatchRepo{})

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

	// 5. Wingman refers themselves
	wingmanRepo.getReferralByIDFn = func(ctx context.Context, referralID uuid.UUID) (*models.WingmanReferral, error) {
		return &models.WingmanReferral{
			ID:        refID,
			Status:    "pending",
			ExpiresAt: time.Now().Add(1 * time.Hour),
			Target1ID: target1ID,
			Target2ID: uuid.New(),
			WingmanID: target1ID, // Wingman refers themselves
		}, nil
	}
	_, err = service.AcceptReferral(context.Background(), refID, target1ID)
	if err == nil || err.Error() != "wingman cannot refer themselves" {
		t.Errorf("Expected wingman cannot refer themselves error, got %v", err)
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
	service := NewWingmanService(wingmanRepo, nil, &mockTxManager{}, &mockWingmanMatchRepo{})

	// 1. Invalid status
	err := service.ProcessCommission(context.Background(), refID)
	if err == nil || err.Error() != "invalid status for commission" {
		t.Errorf("Expected invalid status for commission error, got %v", err)
	}
}

func TestWingmanService_CreateReferral_Errors(t *testing.T) {
	wingmanRepo := &mockWingmanRepo{
		createReferralFn: func(ctx context.Context, referral *models.WingmanReferral) error {
			return errors.New("db error")
		},
	}
	service := NewWingmanService(wingmanRepo, nil, nil, &mockWingmanMatchRepo{})
	
	_, err := service.CreateReferral(context.Background(), uuid.New(), uuid.New(), uuid.New())
	if err == nil || err.Error() != "db error" {
		t.Errorf("Expected db error, got %v", err)
	}
}

func TestWingmanService_AcceptReferral_UpdateError(t *testing.T) {
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
			return errors.New("update error")
		},
	}
	service := NewWingmanService(wingmanRepo, nil, &mockTxManager{}, &mockWingmanMatchRepo{})

	_, err := service.AcceptReferral(context.Background(), refID, target1ID)
	if err == nil || err.Error() != "update error" {
		t.Errorf("Expected update error, got %v", err)
	}
}

func TestWingmanService_ProcessCommission_DBErrors(t *testing.T) {
	refID := uuid.New()
	wingmanID := uuid.New()
	
	// Test AddCommission error
	wingmanRepo1 := &mockWingmanRepo{
		getReferralByIDFn: func(ctx context.Context, referralID uuid.UUID) (*models.WingmanReferral, error) {
			return &models.WingmanReferral{
				ID:        refID,
				Status:    "matched",
				WingmanID: wingmanID,
			}, nil
		},
	}
	walletRepo1 := &mockWalletRepo{
		addCommissionFn: func(ctx context.Context, userID uuid.UUID, amount float64) error {
			return errors.New("add commission error")
		},
	}
	service1 := NewWingmanService(wingmanRepo1, walletRepo1, &mockTxManager{}, &mockWingmanMatchRepo{})
	if err := service1.ProcessCommission(context.Background(), refID); err == nil || err.Error() != "add commission error" {
		t.Errorf("Expected add commission error, got %v", err)
	}

	// Test CreateTransaction error
	walletRepo2 := &mockWalletRepo{
		addCommissionFn: func(ctx context.Context, userID uuid.UUID, amount float64) error { return nil },
		createTransactionFn: func(ctx context.Context, txn *models.WalletTransaction) error {
			return errors.New("create txn error")
		},
	}
	service2 := NewWingmanService(wingmanRepo1, walletRepo2, &mockTxManager{}, &mockWingmanMatchRepo{})
	if err := service2.ProcessCommission(context.Background(), refID); err == nil || err.Error() != "create txn error" {
		t.Errorf("Expected create txn error, got %v", err)
	}

	// Test UpdateReferral error
	wingmanRepo3 := &mockWingmanRepo{
		getReferralByIDFn: func(ctx context.Context, referralID uuid.UUID) (*models.WingmanReferral, error) {
			return &models.WingmanReferral{ID: refID, Status: "matched", WingmanID: wingmanID}, nil
		},
		updateReferralFn: func(ctx context.Context, referral *models.WingmanReferral) error {
			return errors.New("update ref error")
		},
	}
	walletRepo3 := &mockWalletRepo{
		addCommissionFn: func(ctx context.Context, userID uuid.UUID, amount float64) error { return nil },
		createTransactionFn: func(ctx context.Context, txn *models.WalletTransaction) error { return nil },
	}
	service3 := NewWingmanService(wingmanRepo3, walletRepo3, &mockTxManager{}, &mockWingmanMatchRepo{})
	if err := service3.ProcessCommission(context.Background(), refID); err == nil || err.Error() != "update ref error" {
		t.Errorf("Expected update ref error, got %v", err)
	}
}
