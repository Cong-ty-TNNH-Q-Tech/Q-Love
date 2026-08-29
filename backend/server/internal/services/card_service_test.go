// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"context"
	"testing"
	"errors"

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

func TestGetProfile_Success(t *testing.T) {
	mockCardRepo := new(MockCardRepository)
	svc := NewCardService(mockCardRepo, nil, nil, nil, nil)
	userID := uuid.New()

	mockCardRepo.On("GetCardProfile", mock.Anything, userID).Return(&models.CardProfile{
		UserID: userID,
	}, nil)

	profile, err := svc.GetProfile(context.Background(), userID)
	assert.NoError(t, err)
	assert.NotNil(t, profile)
	assert.Equal(t, userID, profile.UserID)
}

func TestGetProfile_CreateOnFly(t *testing.T) {
	mockCardRepo := new(MockCardRepository)
	svc := NewCardService(mockCardRepo, nil, nil, nil, nil)
	userID := uuid.New()

	mockCardRepo.On("GetCardProfile", mock.Anything, userID).Return(nil, errors.New("not found"))
	mockCardRepo.On("CreateCardProfile", mock.Anything, mock.AnythingOfType("*models.CardProfile")).Return(nil)

	profile, err := svc.GetProfile(context.Background(), userID)
	assert.NoError(t, err)
	assert.NotNil(t, profile)
	assert.Equal(t, userID, profile.UserID)
}

func TestGetProfile_CreateOnFly_Error(t *testing.T) {
	mockCardRepo := new(MockCardRepository)
	svc := NewCardService(mockCardRepo, nil, nil, nil, nil)
	userID := uuid.New()

	mockCardRepo.On("GetCardProfile", mock.Anything, userID).Return(nil, errors.New("not found"))
	mockCardRepo.On("CreateCardProfile", mock.Anything, mock.AnythingOfType("*models.CardProfile")).Return(errors.New("create error"))

	profile, err := svc.GetProfile(context.Background(), userID)
	assert.Error(t, err)
	assert.Nil(t, profile)
}

func TestTradeCard_Success_Buy(t *testing.T) {
	mockCardRepo := new(MockCardRepository)
	mockWalletRepo := new(MockWalletRepository)
	mockUserRepo := new(MockUserRepository)
	
	svc := NewCardService(mockCardRepo, mockWalletRepo, mockUserRepo, &dummyTxManager{}, nil)

	collectorID := uuid.New()
	targetUserID := uuid.New()

	mockUserRepo.On("FindByID", mock.Anything, collectorID).Return(&models.User{Level: 5}, nil)

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
}

func TestTradeCard_Success_Sell(t *testing.T) {
	mockCardRepo := new(MockCardRepository)
	mockWalletRepo := new(MockWalletRepository)
	mockUserRepo := new(MockUserRepository)
	
	svc := NewCardService(mockCardRepo, mockWalletRepo, mockUserRepo, &dummyTxManager{}, nil)

	collectorID := uuid.New()
	targetUserID := uuid.New()

	mockUserRepo.On("FindByID", mock.Anything, collectorID).Return(&models.User{Level: 5}, nil)

	mockCardRepo.On("GetCardProfileForUpdate", mock.Anything, targetUserID).Return(&models.CardProfile{
		UserID: targetUserID,
		CurrentPrice: 100,
		AvailableCards: 900,
		TotalCards: 1000,
	}, nil)
	
	mockCardRepo.On("GetOwnedQuantity", mock.Anything, collectorID, targetUserID).Return(5, nil)

	mockWalletRepo.On("UpdateBalance", mock.Anything, collectorID, mock.AnythingOfType("float64")).Return(nil)
	mockCardRepo.On("UpdateCardProfile", mock.Anything, mock.AnythingOfType("*models.CardProfile")).Return(nil)
	mockCardRepo.On("CreateCardTransaction", mock.Anything, mock.AnythingOfType("*models.CardTransaction")).Return(nil)
	mockWalletRepo.On("CreateTransaction", mock.Anything, mock.AnythingOfType("*models.WalletTransaction")).Return(nil)

	err := svc.TradeCard(context.Background(), collectorID, targetUserID, "sell", 1)
	assert.NoError(t, err)
}

func TestTradeCard_Errors(t *testing.T) {
	mockCardRepo := new(MockCardRepository)
	mockWalletRepo := new(MockWalletRepository)
	mockUserRepo := new(MockUserRepository)
	
	svc := NewCardService(mockCardRepo, mockWalletRepo, mockUserRepo, &dummyTxManager{}, nil)

	collectorID := uuid.New()
	
	// Trade own card
	err := svc.TradeCard(context.Background(), collectorID, collectorID, "buy", 1)
	assert.Error(t, err)
	
	// Quantity <= 0
	err = svc.TradeCard(context.Background(), collectorID, uuid.New(), "buy", 0)
	assert.Error(t, err)
	
	// User not found
	mockUserRepo.On("FindByID", mock.Anything, collectorID).Return(nil, errors.New("not found")).Once()
	err = svc.TradeCard(context.Background(), collectorID, uuid.New(), "buy", 1)
	assert.Error(t, err)
	
	// Level < 5
	mockUserRepo.On("FindByID", mock.Anything, collectorID).Return(&models.User{Level: 4}, nil).Once()
	err = svc.TradeCard(context.Background(), collectorID, uuid.New(), "buy", 1)
	assert.Error(t, err)
}

func TestTradeCard_TxErrors(t *testing.T) {
	mockCardRepo := new(MockCardRepository)
	mockWalletRepo := new(MockWalletRepository)
	mockUserRepo := new(MockUserRepository)
	
	svc := NewCardService(mockCardRepo, mockWalletRepo, mockUserRepo, &dummyTxManager{}, nil)

	collectorID := uuid.New()
	targetUserID := uuid.New()

	mockUserRepo.On("FindByID", mock.Anything, collectorID).Return(&models.User{Level: 5}, nil)
	
	// Profile for update fail
	mockCardRepo.On("GetCardProfileForUpdate", mock.Anything, targetUserID).Return(nil, errors.New("fail")).Once()
	mockCardRepo.On("CreateCardProfile", mock.Anything, mock.Anything).Return(errors.New("fail")).Once()
	err := svc.TradeCard(context.Background(), collectorID, targetUserID, "buy", 1)
	assert.Error(t, err)
	
	// Buy but no cards available
	mockCardRepo.On("GetCardProfileForUpdate", mock.Anything, targetUserID).Return(&models.CardProfile{
		UserID: targetUserID, AvailableCards: 0,
	}, nil).Once()
	err = svc.TradeCard(context.Background(), collectorID, targetUserID, "buy", 1)
	assert.Error(t, err)
	
	// Sell but not enough owned
	mockCardRepo.On("GetCardProfileForUpdate", mock.Anything, targetUserID).Return(&models.CardProfile{
		UserID: targetUserID, AvailableCards: 1000,
	}, nil).Once()
	mockCardRepo.On("GetOwnedQuantity", mock.Anything, collectorID, targetUserID).Return(0, nil).Once()
	err = svc.TradeCard(context.Background(), collectorID, targetUserID, "sell", 1)
	assert.Error(t, err)
	
	// Invalid trade type
	mockCardRepo.On("GetCardProfileForUpdate", mock.Anything, targetUserID).Return(&models.CardProfile{
		UserID: targetUserID, AvailableCards: 1000,
	}, nil).Once()
	err = svc.TradeCard(context.Background(), collectorID, targetUserID, "invalid", 1)
	assert.Error(t, err)
}
