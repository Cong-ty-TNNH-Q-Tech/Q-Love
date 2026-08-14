// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"context"
	"testing"
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/google/uuid"
)

type mockStealRepo struct {
	steal *models.CardSteal
	err   error
}

func (m *mockStealRepo) Create(ctx context.Context, steal *models.CardSteal) error {
	return m.err
}

func (m *mockStealRepo) FindByID(ctx context.Context, id uuid.UUID) (*models.CardSteal, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.steal, nil
}

func (m *mockStealRepo) UpdateResult(ctx context.Context, id uuid.UUID, result string) error {
	return m.err
}

func (m *mockStealRepo) TransferCardOwnership(ctx context.Context, collectorID uuid.UUID, targetCardID uuid.UUID) error {
	return m.err
}

type mockMinigameWalletRepo struct {
	wallet *models.UserWallet
	err    error
}

func (m *mockMinigameWalletRepo) AddCommission(ctx context.Context, userID uuid.UUID, amount float64) error { return nil }
func (m *mockMinigameWalletRepo) CreateTransaction(ctx context.Context, txn *models.WalletTransaction) error { return nil }
func (m *mockMinigameWalletRepo) GetWalletForUpdate(ctx context.Context, userID uuid.UUID) (*models.UserWallet, error) {
	return m.wallet, m.err
}
func (m *mockMinigameWalletRepo) UpdateBalance(ctx context.Context, userID uuid.UUID, delta float64) error { return nil }
func (m *mockMinigameWalletRepo) CheckTransactionExists(ctx context.Context, txID uuid.UUID) (bool, error) { return false, nil }

func TestMinigameService_InitSteal_InsufficientBalance(t *testing.T) {
	walletRepo := &mockMinigameWalletRepo{wallet: &models.UserWallet{Balance: 500}} // Less than 1000
	stealRepo := &mockStealRepo{}
	txManager := &mockTxManager{}

	service := NewMinigameService(stealRepo, walletRepo, txManager)

	_, err := service.InitSteal(context.Background(), uuid.New(), uuid.New(), uuid.New())
	if err == nil || err.Error() != "insufficient balance to buy Thẻ Đạo Tặc" {
		t.Errorf("Expected insufficient balance error, got %v", err)
	}
}

func TestMinigameService_SubmitStealResult_Win(t *testing.T) {
	attackerID := uuid.New()
	steal := &models.CardSteal{
		ID:         uuid.New(),
		AttackerID: attackerID,
		Result:     "pending",
		CreatedAt:  time.Now().Add(-12 * time.Second), // Valid time for anti-cheat
	}

	walletRepo := &mockMinigameWalletRepo{wallet: &models.UserWallet{Balance: 2000}}
	stealRepo := &mockStealRepo{steal: steal}
	txManager := &mockTxManager{}

	service := NewMinigameService(stealRepo, walletRepo, txManager)

	err := service.SubmitStealResult(context.Background(), steal.ID, attackerID, true)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestMinigameService_SubmitStealResult_Lose(t *testing.T) {
	attackerID := uuid.New()
	steal := &models.CardSteal{
		ID:         uuid.New(),
		AttackerID: attackerID,
		Result:     "pending",
		CreatedAt:  time.Now().Add(-12 * time.Second),
	}

	walletRepo := &mockMinigameWalletRepo{wallet: &models.UserWallet{Balance: 2000}}
	stealRepo := &mockStealRepo{steal: steal}
	txManager := &mockTxManager{}

	service := NewMinigameService(stealRepo, walletRepo, txManager)

	err := service.SubmitStealResult(context.Background(), steal.ID, attackerID, false)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestMinigameService_InitSteal_Success(t *testing.T) {
	walletRepo := &mockMinigameWalletRepo{wallet: &models.UserWallet{Balance: 5000}}
	stealRepo := &mockStealRepo{}
	txManager := &mockTxManager{}

	service := NewMinigameService(stealRepo, walletRepo, txManager)

	steal, err := service.InitSteal(context.Background(), uuid.New(), uuid.New(), uuid.New())
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if steal == nil {
		t.Errorf("Expected steal object, got nil")
	}
}

func TestMinigameService_InitSteal_WalletError(t *testing.T) {
	walletRepo := &mockMinigameWalletRepo{err: context.DeadlineExceeded}
	stealRepo := &mockStealRepo{}
	txManager := &mockTxManager{}

	service := NewMinigameService(stealRepo, walletRepo, txManager)

	_, err := service.InitSteal(context.Background(), uuid.New(), uuid.New(), uuid.New())
	if err != context.DeadlineExceeded {
		t.Errorf("Expected wallet error, got %v", err)
	}
}

func TestMinigameService_SubmitStealResult_AntiCheat(t *testing.T) {
	attackerID := uuid.New()
	steal := &models.CardSteal{
		ID:         uuid.New(),
		AttackerID: attackerID,
		Result:     "pending",
		CreatedAt:  time.Now().Add(-2 * time.Second), // Too fast
	}

	walletRepo := &mockMinigameWalletRepo{wallet: &models.UserWallet{Balance: 2000}}
	stealRepo := &mockStealRepo{steal: steal}
	txManager := &mockTxManager{}

	service := NewMinigameService(stealRepo, walletRepo, txManager)

	err := service.SubmitStealResult(context.Background(), steal.ID, attackerID, true)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestMinigameService_SubmitStealResult_Unauthorized(t *testing.T) {
	attackerID := uuid.New()
	steal := &models.CardSteal{
		ID:         uuid.New(),
		AttackerID: uuid.New(), // Mismatch
		Result:     "pending",
		CreatedAt:  time.Now().Add(-15 * time.Second),
	}

	walletRepo := &mockMinigameWalletRepo{wallet: &models.UserWallet{Balance: 2000}}
	stealRepo := &mockStealRepo{steal: steal}
	txManager := &mockTxManager{}

	service := NewMinigameService(stealRepo, walletRepo, txManager)

	err := service.SubmitStealResult(context.Background(), steal.ID, attackerID, true)
	if err == nil || err.Error() != "unauthorized to submit result for this steal session" {
		t.Errorf("Expected unauthorized error, got %v", err)
	}
}

func TestMinigameService_SubmitStealResult_AlreadyDone(t *testing.T) {
	attackerID := uuid.New()
	steal := &models.CardSteal{
		ID:         uuid.New(),
		AttackerID: attackerID,
		Result:     "attacker_won", // Not pending
		CreatedAt:  time.Now().Add(-15 * time.Second),
	}

	walletRepo := &mockMinigameWalletRepo{wallet: &models.UserWallet{Balance: 2000}}
	stealRepo := &mockStealRepo{steal: steal}
	txManager := &mockTxManager{}

	service := NewMinigameService(stealRepo, walletRepo, txManager)

	err := service.SubmitStealResult(context.Background(), steal.ID, attackerID, true)
	if err == nil || err.Error() != "steal session already completed" {
		t.Errorf("Expected already concluded error, got %v", err)
	}
}

func TestMinigameService_SubmitStealResult_NotFound(t *testing.T) {
	walletRepo := &mockMinigameWalletRepo{wallet: &models.UserWallet{Balance: 2000}}
	stealRepo := &mockStealRepo{err: context.DeadlineExceeded} // simulate db error
	txManager := &mockTxManager{}

	service := NewMinigameService(stealRepo, walletRepo, txManager)

	err := service.SubmitStealResult(context.Background(), uuid.New(), uuid.New(), true)
	if err == nil || err.Error() != "steal session not found" {
		t.Errorf("Expected not found error, got %v", err)
	}
}
