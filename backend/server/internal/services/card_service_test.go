// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"context"
	"testing"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"database/sql"
)

type MockCardRepository struct {
	mock.Mock
}

func (m *MockCardRepository) GetCardProfile(ctx context.Context, userID uuid.UUID) (*models.CardProfile, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CardProfile), args.Error(1)
}

func (m *MockCardRepository) GetCardProfileForUpdate(ctx context.Context, userID uuid.UUID) (*models.CardProfile, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CardProfile), args.Error(1)
}

func (m *MockCardRepository) CreateCardProfile(ctx context.Context, profile *models.CardProfile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

func (m *MockCardRepository) UpdateCardProfile(ctx context.Context, profile *models.CardProfile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

func (m *MockCardRepository) CreateCardTransaction(ctx context.Context, tx *models.CardTransaction) error {
	args := m.Called(ctx, tx)
	return args.Error(0)
}

func (m *MockCardRepository) GetOwnedQuantity(ctx context.Context, collectorID, targetUserID uuid.UUID) (int, error) {
	args := m.Called(ctx, collectorID, targetUserID)
	return args.Int(0), args.Error(1)
}

type MockWalletRepository struct {
	mock.Mock
}
func (m *MockWalletRepository) AddCommission(ctx context.Context, userID uuid.UUID, amount float64) error { return nil }
func (m *MockWalletRepository) UpdateBalance(ctx context.Context, userID uuid.UUID, amount float64) error {
	args := m.Called(ctx, userID, amount)
	return args.Error(0)
}
func (m *MockWalletRepository) CreateTransaction(ctx context.Context, transaction *models.WalletTransaction) error {
	args := m.Called(ctx, transaction)
	return args.Error(0)
}
func (m *MockWalletRepository) GetWalletForUpdate(ctx context.Context, userID uuid.UUID) (*models.UserWallet, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserWallet), args.Error(1)
}
func (m *MockWalletRepository) HoldBalance(ctx context.Context, userID uuid.UUID, amount float64) error { return nil }
func (m *MockWalletRepository) ReleaseHoldBalance(ctx context.Context, userID uuid.UUID, amount float64) error { return nil }
func (m *MockWalletRepository) CheckTransactionExists(ctx context.Context, txID uuid.UUID) (bool, error) { return false, nil }

type MockUserRepository struct {
	mock.Mock
}
func (m *MockUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}
func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error { return nil }
func (m *MockUserRepository) FindByPhone(ctx context.Context, phone string) (*models.User, error) { return nil, nil }
func (m *MockUserRepository) GetTopUsersByScore(ctx context.Context, limit int) ([]uuid.UUID, error) { return nil, nil }
func (m *MockUserRepository) GetFeed(ctx context.Context, userID uuid.UUID, radius int) ([]models.User, error) { return nil, nil }
func (m *MockUserRepository) GetSpiritualFeed(ctx context.Context, userID uuid.UUID, radius int) ([]models.User, error) { return nil, nil }

// Dummy TxManager
type dummyTxManager struct{}
func (d *dummyTxManager) WithTransaction(ctx context.Context, fn func(context.Context) error, opts ...*sql.TxOptions) error {
	return fn(ctx)
}

func TestTradeCard_Success_Buy(t *testing.T) {
	mockCardRepo := new(MockCardRepository)
	mockWalletRepo := new(MockWalletRepository)
	mockUserRepo := new(MockUserRepository)
	
	svc := NewCardService(mockCardRepo, mockWalletRepo, mockUserRepo, &dummyTxManager{}, nil)

	collectorID := uuid.New()
	targetUserID := uuid.New()

	// Level 5 check
	mockUserRepo.On("FindByID", mock.Anything, collectorID).Return(&models.User{Level: 5}, nil)

	// Profile fetching
	mockCardRepo.On("GetCardProfileForUpdate", mock.Anything, targetUserID).Return(&models.CardProfile{
		UserID: targetUserID,
		CurrentPrice: 100,
		AvailableCards: 1000,
		TotalCards: 1000,
	}, nil)

	mockWalletRepo.On("UpdateBalance", mock.Anything, collectorID, mock.AnythingOfType("float64")).Return(nil)
	mockCardRepo.On("UpdateCardProfile", mock.Anything, mock.AnythingOfType("*models.CardProfile")).Return(nil)
	mockCardRepo.On("CreateCardTransaction", mock.Anything, mock.AnythingOfType("*models.CardTransaction")).Return(nil)
	mockWalletRepo.On("CreateTransaction", mock.Anything, mock.AnythingOfType("*models.WalletTransaction")).Return(nil)

	err := svc.TradeCard(context.Background(), collectorID, targetUserID, "buy", 1)
	assert.NoError(t, err)

	mockCardRepo.AssertExpectations(t)
	mockWalletRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}
