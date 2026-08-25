// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
)

// --- Mocks for ShameRepo ---

type mockShameRepo struct {
	getActiveShamesFn      func(ctx context.Context, limit, offset int) ([]models.WallOfShameResponse, error)
	incrementTomatoCountFn func(ctx context.Context, shameID uuid.UUID, count int) error
}

func (m *mockShameRepo) GetActiveShames(ctx context.Context, limit, offset int) ([]models.WallOfShameResponse, error) {
	return m.getActiveShamesFn(ctx, limit, offset)
}

func (m *mockShameRepo) IncrementTomatoCount(ctx context.Context, shameID uuid.UUID, count int) error {
	return m.incrementTomatoCountFn(ctx, shameID, count)
}

// We will also use mockWalletRepo and mockTxManager from wingman_service_test.go
// Wait, actually mockWalletRepo in wingman_service_test.go does NOT have GetWalletForUpdate and UpdateBalance!
// I need to add them to a local struct or I'll just redefine it here to avoid conflicts.
// Wait, if I redefine `mockWalletRepo` in the same package, it will be a redeclaration error.
// So let's name it `shameMockWalletRepo`.

type shameMockWalletRepo struct {
	getWalletForUpdateFn func(ctx context.Context, userID uuid.UUID) (*models.UserWallet, error)
	updateBalanceFn      func(ctx context.Context, userID uuid.UUID, delta float64) error
	createTransactionFn  func(ctx context.Context, txn *models.WalletTransaction) error
}

// Ensure it implements repository.WalletRepository (it doesn't have AddCommission, but go interface duck typing handles it if we don't pass it as WalletRepository directly).
// Wait, I need to implement the full interface. The interface in WalletRepository probably has AddCommission, GetWalletForUpdate, UpdateBalance, CreateTransaction, etc.
// Since I can't be sure of the exact interface of WalletRepository, I'll use a mocked repository by just ignoring methods we don't use or implement them as panic.
// Actually let's just implement the ones I need.

func (m *shameMockWalletRepo) GetWallet(ctx context.Context, userID uuid.UUID) (*models.UserWallet, error) {
	return nil, nil
}
func (m *shameMockWalletRepo) CreateWallet(ctx context.Context, userID uuid.UUID) error { return nil }
func (m *shameMockWalletRepo) GetWalletForUpdate(ctx context.Context, userID uuid.UUID) (*models.UserWallet, error) {
	return m.getWalletForUpdateFn(ctx, userID)
}
func (m *shameMockWalletRepo) UpdateBalance(ctx context.Context, userID uuid.UUID, delta float64) error {
	return m.updateBalanceFn(ctx, userID, delta)
}
func (m *shameMockWalletRepo) CreateTransaction(ctx context.Context, txn *models.WalletTransaction) error {
	return m.createTransactionFn(ctx, txn)
}
func (m *shameMockWalletRepo) AddCommission(ctx context.Context, userID uuid.UUID, amount float64) error {
	return nil
}
func (m *shameMockWalletRepo) CheckTransactionExists(ctx context.Context, txID uuid.UUID) (bool, error) {
	return false, nil
}
func (m *shameMockWalletRepo) HoldBalance(ctx context.Context, userID uuid.UUID, amount float64) error {
	return nil
}
func (m *shameMockWalletRepo) ReleaseHoldBalance(ctx context.Context, userID uuid.UUID, amount float64) error {
	return nil
}

// --- Tests ---

func TestShameService_GetActiveShames(t *testing.T) {
	shameRepo := &mockShameRepo{
		getActiveShamesFn: func(ctx context.Context, limit, offset int) ([]models.WallOfShameResponse, error) {
			return []models.WallOfShameResponse{
				{WallOfShame: models.WallOfShame{ID: uuid.New(), UserID: uuid.New(), Reason: "Cheating", TomatoesThrown: 5, ExpiresAt: time.Now().Add(1 * time.Hour)}, UserName: "TestUser", AvatarURL: "test.jpg"},
			}, nil
		},
	}
	service := NewShameService(shameRepo, nil, nil)

	shames, err := service.GetActiveShames(context.Background(), 10, 0)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(shames) != 1 {
		t.Fatalf("Expected 1 shame, got %d", len(shames))
	}
}

func TestShameService_ThrowTomato_Success(t *testing.T) {
	throwerID := uuid.New()
	shameID := uuid.New()

	shameRepo := &mockShameRepo{
		incrementTomatoCountFn: func(ctx context.Context, shameID uuid.UUID, count int) error {
			return nil
		},
	}
	walletRepo := &shameMockWalletRepo{
		getWalletForUpdateFn: func(ctx context.Context, userID uuid.UUID) (*models.UserWallet, error) {
			return &models.UserWallet{UserID: throwerID, Balance: 10.0}, nil
		},
		updateBalanceFn: func(ctx context.Context, userID uuid.UUID, delta float64) error {
			return nil
		},
		createTransactionFn: func(ctx context.Context, txn *models.WalletTransaction) error {
			return nil
		},
	}

	service := NewShameService(shameRepo, walletRepo, &mockTxManager{})

	err := service.ThrowTomato(context.Background(), throwerID, shameID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestShameService_ThrowTomato_InsufficientBalance(t *testing.T) {
	throwerID := uuid.New()
	shameID := uuid.New()

	walletRepo := &shameMockWalletRepo{
		getWalletForUpdateFn: func(ctx context.Context, userID uuid.UUID) (*models.UserWallet, error) {
			return &models.UserWallet{UserID: throwerID, Balance: 0.5}, nil // Not enough
		},
	}

	service := NewShameService(nil, walletRepo, &mockTxManager{})

	err := service.ThrowTomato(context.Background(), throwerID, shameID)
	if err == nil || err.Error() != "insufficient balance to throw a tomato" {
		t.Fatalf("Expected insufficient balance error, got %v", err)
	}
}

func TestShameService_ThrowTomato_DBError(t *testing.T) {
	throwerID := uuid.New()
	shameID := uuid.New()

	walletRepo := &shameMockWalletRepo{
		getWalletForUpdateFn: func(ctx context.Context, userID uuid.UUID) (*models.UserWallet, error) {
			return nil, errors.New("db error")
		},
	}

	service := NewShameService(nil, walletRepo, &mockTxManager{})

	err := service.ThrowTomato(context.Background(), throwerID, shameID)
	if err == nil || err.Error() != "db error" {
		t.Fatalf("Expected db error, got %v", err)
	}
}
