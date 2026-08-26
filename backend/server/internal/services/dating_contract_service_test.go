// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocks
type mockDatingContractRepo struct{ mock.Mock }
func (m *mockDatingContractRepo) Create(ctx context.Context, contract *models.DatingContract) error {
	args := m.Called(ctx, contract)
	return args.Error(0)
}
func (m *mockDatingContractRepo) GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*models.DatingContract, error) {
	args := m.Called(ctx, id)
	var contract *models.DatingContract
	if args.Get(0) != nil {
		contract = args.Get(0).(*models.DatingContract)
	}
	return contract, args.Error(1)
}
func (m *mockDatingContractRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.DatingContract, error) {
	args := m.Called(ctx, id)
	var contract *models.DatingContract
	if args.Get(0) != nil {
		contract = args.Get(0).(*models.DatingContract)
	}
	return contract, args.Error(1)
}
func (m *mockDatingContractRepo) Update(ctx context.Context, contract *models.DatingContract) error {
	args := m.Called(ctx, contract)
	return args.Error(0)
}

type mockWalletRepo struct{ mock.Mock }
func (m *mockWalletRepo) AddCommission(ctx context.Context, userID uuid.UUID, amount float64) error {
	return m.Called(ctx, userID, amount).Error(0)
}
func (m *mockWalletRepo) CreateTransaction(ctx context.Context, txn *models.WalletTransaction) error {
	return m.Called(ctx, txn).Error(0)
}
func (m *mockWalletRepo) GetWalletForUpdate(ctx context.Context, userID uuid.UUID) (*models.UserWallet, error) {
	args := m.Called(ctx, userID)
	var w *models.UserWallet
	if args.Get(0) != nil {
		w = args.Get(0).(*models.UserWallet)
	}
	return w, args.Error(1)
}
func (m *mockWalletRepo) UpdateBalance(ctx context.Context, userID uuid.UUID, delta float64) error {
	return m.Called(ctx, userID, delta).Error(0)
}
func (m *mockWalletRepo) HoldBalance(ctx context.Context, userID uuid.UUID, amount float64) error {
	return m.Called(ctx, userID, amount).Error(0)
}
func (m *mockWalletRepo) ReleaseHoldBalance(ctx context.Context, userID uuid.UUID, amount float64) error {
	return m.Called(ctx, userID, amount).Error(0)
}
func (m *mockWalletRepo) CheckTransactionExists(ctx context.Context, txID uuid.UUID) (bool, error) {
	args := m.Called(ctx, txID)
	return args.Bool(0), args.Error(1)
}

type mockMatchRepo struct{ mock.Mock }
func (m *mockMatchRepo) Create(ctx context.Context, match *models.Match) error { return nil }
func (m *mockMatchRepo) FindByID(ctx context.Context, id uuid.UUID) (*models.Match, error) { return nil, nil }
func (m *mockMatchRepo) FindByIDUnscoped(ctx context.Context, id uuid.UUID) (*models.Match, error) { return nil, nil }
func (m *mockMatchRepo) UpdateLastInteraction(ctx context.Context, id uuid.UUID, t time.Time) error { return nil }
func (m *mockMatchRepo) SoftDelete(ctx context.Context, id uuid.UUID) error { return nil }
func (m *mockMatchRepo) ResetStreakForInactiveMatches(ctx context.Context, inactiveDuration time.Duration) error { return nil }
func (m *mockMatchRepo) ResetIslandLevelForInactiveMatches(ctx context.Context, inactiveDuration time.Duration) error { return nil }
func (m *mockMatchRepo) FindByUsers(ctx context.Context, u1, u2 uuid.UUID) (*models.Match, error) {
	args := m.Called(ctx, u1, u2)
	var match *models.Match
	if args.Get(0) != nil {
		match = args.Get(0).(*models.Match)
	}
	return match, args.Error(1)
}

type mockChatRepo struct{ mock.Mock }
func (m *mockChatRepo) Create(ctx context.Context, msg *models.ChatMessage) error { return nil }
func (m *mockChatRepo) GetMessagesByMatchID(ctx context.Context, matchID uuid.UUID, limit int, before *time.Time) ([]models.ChatMessage, error) { return nil, nil }
func (m *mockChatRepo) CountMessagesByMatchID(ctx context.Context, matchID uuid.UUID) (int64, error) {
	args := m.Called(ctx, matchID)
	return args.Get(0).(int64), args.Error(1)
}

type mockPremiumRepo struct{ mock.Mock }
func (m *mockPremiumRepo) Upsert(ctx context.Context, p *models.UserPremium) error { return nil }
func (m *mockPremiumRepo) ActivatePremium(ctx context.Context, userID uuid.UUID, expiresAt time.Time) error { return nil }
func (m *mockPremiumRepo) IsUserPremium(ctx context.Context, userID uuid.UUID) (bool, error) { return false, nil }
func (m *mockPremiumRepo) FindByUserID(ctx context.Context, id uuid.UUID) (*models.UserPremium, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.UserPremium), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockPremiumRepo) Update(ctx context.Context, p *models.UserPremium) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

type mockTxManager struct{}
func (m *mockTxManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error, opts ...*sql.TxOptions) error {
	return fn(ctx)
}

func TestDatingContractService_CreateContract(t *testing.T) {
	contractRepo := new(mockDatingContractRepo)
	walletRepo := new(mockWalletRepo)
	matchRepo := new(mockMatchRepo)
	chatRepo := new(mockChatRepo)
	premiumRepo := new(mockPremiumRepo)
	txManager := new(mockTxManager)

	service := services.NewDatingContractService(contractRepo, walletRepo, matchRepo, chatRepo, premiumRepo, txManager)
	ctx := context.Background()
	userA := uuid.New()
	userB := uuid.New()
	matchID := uuid.New()

	matchRepo.On("FindByUsers", ctx, userA, userB).Return(&models.Match{ID: matchID}, nil)
	chatRepo.On("CountMessagesByMatchID", ctx, matchID).Return(int64(25), nil)
	
	walletRepo.On("GetWalletForUpdate", ctx, userA).Return(&models.UserWallet{UserID: userA, Balance: 500}, nil)
	walletRepo.On("HoldBalance", ctx, userA, 100.0).Return(nil)
	walletRepo.On("CreateTransaction", ctx, mock.AnythingOfType("*models.WalletTransaction")).Return(nil)
	contractRepo.On("Create", ctx, mock.AnythingOfType("*models.DatingContract")).Return(nil)

	contract, err := service.CreateContract(ctx, userA, userB, 100.0, time.Now().Add(24*time.Hour))
	
	assert.NoError(t, err)
	assert.NotNil(t, contract)
	assert.Equal(t, userA, contract.UserAID)
	assert.Equal(t, 100.0, contract.DepositAmount)
	assert.Equal(t, "pending", contract.Status)
}

func TestDatingContractService_CancelContract_FreePremium(t *testing.T) {
	contractRepo := new(mockDatingContractRepo)
	walletRepo := new(mockWalletRepo)
	matchRepo := new(mockMatchRepo)
	chatRepo := new(mockChatRepo)
	premiumRepo := new(mockPremiumRepo)
	txManager := new(mockTxManager)

	service := services.NewDatingContractService(contractRepo, walletRepo, matchRepo, chatRepo, premiumRepo, txManager)
	ctx := context.Background()
	userA := uuid.New()
	userB := uuid.New()
	contractID := uuid.New()

	contract := &models.DatingContract{
		ID: contractID,
		UserAID: userA,
		UserBID: userB,
		DepositAmount: 100.0,
		Status: "active",
	}

	contractRepo.On("GetByIDForUpdate", ctx, contractID).Return(contract, nil)
	
	// UserA is premium
	premium := &models.UserPremium{
		UserID: userA,
		ExpiresAt: time.Now().Add(10*time.Hour),
		FreeCancelLeft: 1,
	}
	premiumRepo.On("FindByUserID", ctx, userA).Return(premium, nil)
	premiumRepo.On("Update", ctx, mock.AnythingOfType("*models.UserPremium")).Return(nil)
	
	walletRepo.On("ReleaseHoldBalance", ctx, userA, 100.0).Return(nil)
	walletRepo.On("CreateTransaction", ctx, mock.Anything).Return(nil)
	walletRepo.On("ReleaseHoldBalance", ctx, userB, 100.0).Return(nil)
	
	contractRepo.On("Update", ctx, mock.Anything).Return(nil)

	err := service.CancelContract(ctx, contractID, userA, "bận")
	assert.NoError(t, err)
	assert.Equal(t, "cancelled", contract.Status)
	assert.Equal(t, 0, premium.FreeCancelLeft)
}

func TestDatingContractService_AcceptContract(t *testing.T) {
	contractRepo := new(mockDatingContractRepo)
	walletRepo := new(mockWalletRepo)
	matchRepo := new(mockMatchRepo)
	chatRepo := new(mockChatRepo)
	premiumRepo := new(mockPremiumRepo)
	txManager := new(mockTxManager)

	service := services.NewDatingContractService(contractRepo, walletRepo, matchRepo, chatRepo, premiumRepo, txManager)
	ctx := context.Background()
	userB := uuid.New()
	contractID := uuid.New()

	contract := &models.DatingContract{
		ID:            contractID,
		UserBID:       userB,
		DepositAmount: 100.0,
		Status:        "pending",
	}

	contractRepo.On("GetByIDForUpdate", ctx, contractID).Return(contract, nil)
	walletRepo.On("GetWalletForUpdate", ctx, userB).Return(&models.UserWallet{UserID: userB, Balance: 500}, nil)
	walletRepo.On("HoldBalance", ctx, userB, 100.0).Return(nil)
	walletRepo.On("CreateTransaction", ctx, mock.AnythingOfType("*models.WalletTransaction")).Return(nil)
	contractRepo.On("Update", ctx, mock.AnythingOfType("*models.DatingContract")).Return(nil)

	acceptedContract, err := service.AcceptContract(ctx, contractID, userB)
	assert.NoError(t, err)
	assert.NotNil(t, acceptedContract)
	assert.Equal(t, "active", acceptedContract.Status)
}

func TestDatingContractService_ScanContract(t *testing.T) {
	contractRepo := new(mockDatingContractRepo)
	walletRepo := new(mockWalletRepo)
	matchRepo := new(mockMatchRepo)
	chatRepo := new(mockChatRepo)
	premiumRepo := new(mockPremiumRepo)
	txManager := new(mockTxManager)

	service := services.NewDatingContractService(contractRepo, walletRepo, matchRepo, chatRepo, premiumRepo, txManager)
	ctx := context.Background()
	userA := uuid.New()
	userB := uuid.New()
	contractID := uuid.New()

	contract := &models.DatingContract{
		ID: contractID,
		UserAID: userA,
		UserBID: userB,
		DepositAmount: 100.0,
		Status: "active",
		TOTPSecret: "mysecret",
	}

	contractRepo.On("GetByIDForUpdate", ctx, contractID).Return(contract, nil)
	walletRepo.On("ReleaseHoldBalance", ctx, userA, 100.0).Return(nil)
	walletRepo.On("ReleaseHoldBalance", ctx, userB, 100.0).Return(nil)
	walletRepo.On("CreateTransaction", ctx, mock.Anything).Return(nil)
	contractRepo.On("Update", ctx, mock.Anything).Return(nil)

	err := service.ScanContract(ctx, contractID, "mysecret")
	assert.NoError(t, err)
	assert.Equal(t, "completed", contract.Status)
}

func TestDatingContractService_AcceptContract_Errors(t *testing.T) {
	contractRepo := new(mockDatingContractRepo)
	walletRepo := new(mockWalletRepo)
	matchRepo := new(mockMatchRepo)
	chatRepo := new(mockChatRepo)
	premiumRepo := new(mockPremiumRepo)
	txManager := new(mockTxManager)

	service := services.NewDatingContractService(contractRepo, walletRepo, matchRepo, chatRepo, premiumRepo, txManager)
	ctx := context.Background()
	userB := uuid.New()
	contractID := uuid.New()

	t.Run("Forbidden", func(t *testing.T) {
		contractRepo.On("GetByIDForUpdate", ctx, contractID).Return(&models.DatingContract{UserBID: uuid.New()}, nil).Once()
		_, err := service.AcceptContract(ctx, contractID, userB)
		assert.ErrorContains(t, err, "forbidden")
	})

	t.Run("NotPending", func(t *testing.T) {
		contractRepo.On("GetByIDForUpdate", ctx, contractID).Return(&models.DatingContract{UserBID: userB, Status: "active"}, nil).Once()
		_, err := service.AcceptContract(ctx, contractID, userB)
		assert.ErrorContains(t, err, "not in pending state")
	})

	t.Run("InsufficientBalance", func(t *testing.T) {
		contractRepo.On("GetByIDForUpdate", ctx, contractID).Return(&models.DatingContract{UserBID: userB, Status: "pending", DepositAmount: 100}, nil).Once()
		walletRepo.On("GetWalletForUpdate", ctx, userB).Return(&models.UserWallet{UserID: userB, Balance: 50}, nil).Once()
		_, err := service.AcceptContract(ctx, contractID, userB)
		assert.ErrorContains(t, err, "insufficient balance")
	})
}

func TestDatingContractService_CancelContract_Errors(t *testing.T) {
	contractRepo := new(mockDatingContractRepo)
	walletRepo := new(mockWalletRepo)
	matchRepo := new(mockMatchRepo)
	chatRepo := new(mockChatRepo)
	premiumRepo := new(mockPremiumRepo)
	txManager := new(mockTxManager)

	service := services.NewDatingContractService(contractRepo, walletRepo, matchRepo, chatRepo, premiumRepo, txManager)
	ctx := context.Background()
	userID := uuid.New()
	contractID := uuid.New()

	t.Run("InvalidState", func(t *testing.T) {
		contractRepo.On("GetByIDForUpdate", ctx, contractID).Return(&models.DatingContract{Status: "completed"}, nil).Once()
		err := service.CancelContract(ctx, contractID, userID, "reason")
		assert.ErrorContains(t, err, "current state")
	})

	t.Run("Forbidden", func(t *testing.T) {
		contractRepo.On("GetByIDForUpdate", ctx, contractID).Return(&models.DatingContract{Status: "active", UserAID: uuid.New(), UserBID: uuid.New()}, nil).Once()
		err := service.CancelContract(ctx, contractID, userID, "reason")
		assert.ErrorContains(t, err, "forbidden")
	})
}
