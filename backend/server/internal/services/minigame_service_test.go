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
