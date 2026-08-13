package services

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/repository"
	"github.com/google/uuid"
)

// MockTxManager
type MockTxManager struct{}

func (m *MockTxManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error, opts ...*sql.TxOptions) error {
	return fn(ctx)
}

// MockWalletRepo
type MockWalletRepoIAP struct {
	repository.WalletRepository
	CheckTxFunc func(ctx context.Context, txID uuid.UUID) (bool, error)
	AddBalFunc  func(ctx context.Context, userID uuid.UUID, amount float64) error
	LogTxFunc   func(ctx context.Context, txn *models.WalletTransaction) error
}

func (m *MockWalletRepoIAP) CheckTransactionExists(ctx context.Context, txID uuid.UUID) (bool, error) {
	if m.CheckTxFunc != nil {
		return m.CheckTxFunc(ctx, txID)
	}
	return false, nil
}

func (m *MockWalletRepoIAP) UpdateBalance(ctx context.Context, userID uuid.UUID, amount float64) error {
	if m.AddBalFunc != nil {
		return m.AddBalFunc(ctx, userID, amount)
	}
	return nil
}

func (m *MockWalletRepoIAP) CreateTransaction(ctx context.Context, txn *models.WalletTransaction) error {
	if m.LogTxFunc != nil {
		return m.LogTxFunc(ctx, txn)
	}
	return nil
}

// MockUserPremiumRepo
type MockUserPremiumRepoIAP struct {
	repository.UserPremiumRepository
	ActivateFunc func(ctx context.Context, userID uuid.UUID, expiresAt time.Time) error
}

func (m *MockUserPremiumRepoIAP) ActivatePremium(ctx context.Context, userID uuid.UUID, expiresAt time.Time) error {
	if m.ActivateFunc != nil {
		return m.ActivateFunc(ctx, userID, expiresAt)
	}
	return nil
}

func TestIAPService_ProcessRevenueCatWebhook_Deposit(t *testing.T) {
	walletRepo := &MockWalletRepoIAP{
		CheckTxFunc: func(ctx context.Context, txID uuid.UUID) (bool, error) {
			return false, nil
		},
		AddBalFunc: func(ctx context.Context, userID uuid.UUID, amount float64) error {
			if amount != 100.0 {
				t.Errorf("Expected 100 coins, got %v", amount)
			}
			return nil
		},
		LogTxFunc: func(ctx context.Context, txn *models.WalletTransaction) error {
			if txn.Type != "iap_deposit" {
				t.Errorf("Expected iap_deposit, got %s", txn.Type)
			}
			return nil
		},
	}
	userPremRepo := &MockUserPremiumRepoIAP{}
	txManager := &MockTxManager{}

	svc := NewIAPService(txManager, walletRepo, userPremRepo)

	event := RevenueCatEvent{
		Type:          "NON_RENEWING_PURCHASE",
		AppUserID:     uuid.New().String(),
		ProductID:     "coin_pack_100",
		TransactionID: "tx_12345",
	}

	err := svc.ProcessRevenueCatWebhook(context.Background(), event)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestIAPService_ProcessRevenueCatWebhook_Premium(t *testing.T) {
	walletRepo := &MockWalletRepoIAP{}
	userPremRepo := &MockUserPremiumRepoIAP{
		ActivateFunc: func(ctx context.Context, userID uuid.UUID, expiresAt time.Time) error {
			return nil
		},
	}
	txManager := &MockTxManager{}

	svc := NewIAPService(txManager, walletRepo, userPremRepo)

	event := RevenueCatEvent{
		Type:          "INITIAL_PURCHASE",
		AppUserID:     uuid.New().String(),
		ProductID:     "qlove_premium_1month",
		TransactionID: "tx_99999",
	}

	err := svc.ProcessRevenueCatWebhook(context.Background(), event)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestIAPService_ProcessRevenueCatWebhook_InvalidUser(t *testing.T) {
	svc := NewIAPService(&MockTxManager{}, &MockWalletRepoIAP{}, &MockUserPremiumRepoIAP{})

	event := RevenueCatEvent{
		Type:          "NON_RENEWING_PURCHASE",
		AppUserID:     "invalid-uuid",
		ProductID:     "coin_pack_100",
		TransactionID: "tx_12345",
	}

	err := svc.ProcessRevenueCatWebhook(context.Background(), event)
	if err != ErrInvalidWebhookPayload {
		t.Fatalf("Expected ErrInvalidWebhookPayload, got %v", err)
	}
}
