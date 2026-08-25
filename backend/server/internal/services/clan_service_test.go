// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type mockClanRepo struct {
	findByNameErr  error
	findByNameClan *models.Clan
	createErr      error
	addMemberErr   error
}

func (m *mockClanRepo) FindByName(ctx context.Context, name string) (*models.Clan, error) {
	return m.findByNameClan, m.findByNameErr
}

func (m *mockClanRepo) CreateClan(ctx context.Context, clan *models.Clan) error {
	clan.ID = uuid.New()
	return m.createErr
}

func (m *mockClanRepo) AddClanMember(ctx context.Context, member *models.ClanMember) error {
	return m.addMemberErr
}

func (m *mockClanRepo) GetMembers(ctx context.Context, clanID uuid.UUID) ([]models.ClanMember, error) {
	return nil, nil
}
func (m *mockClanRepo) GetTopWeeklyClan(ctx context.Context) (*models.Clan, error) {
	return nil, nil
}

func (m *mockClanRepo) ResetWeeklyScores(ctx context.Context) error {
	return nil
}

type mockWalletRepoClan struct {
	wallet    *models.UserWallet
	walletErr error
	updateErr error
	createTxErr error
}

func (m *mockWalletRepoClan) AddCommission(ctx context.Context, userID uuid.UUID, amount float64) error {
	return nil
}

func (m *mockWalletRepoClan) CheckTransactionExists(ctx context.Context, txID uuid.UUID) (bool, error) {
	return false, nil
}

func (m *mockWalletRepoClan) GetWalletForUpdate(ctx context.Context, userID uuid.UUID) (*models.UserWallet, error) {
	return m.wallet, m.walletErr
}

func (m *mockWalletRepoClan) UpdateBalance(ctx context.Context, userID uuid.UUID, amount float64) error {
	return m.updateErr
}

func (m *mockWalletRepoClan) CreateTransaction(ctx context.Context, transaction *models.WalletTransaction) error {
	return m.createTxErr
}

type mockTxManagerClan struct{}

func (m *mockTxManagerClan) WithTransaction(ctx context.Context, fn func(ctx context.Context) error, opts ...*sql.TxOptions) error {
	return fn(ctx)
}

func TestCreateClan(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mClan := &mockClanRepo{findByNameErr: gorm.ErrRecordNotFound}
		mWallet := &mockWalletRepoClan{wallet: &models.UserWallet{Balance: 1000}}
		mTx := &mockTxManagerClan{}

		svc := NewClanService(mClan, mWallet, mTx)
		clan, err := svc.CreateClan(ctx, userID, "Test Clan", "Slogan", "url")
		
		assert.NoError(t, err)
		assert.NotNil(t, clan)
		assert.Equal(t, "Test Clan", clan.Name)
	})

	t.Run("Clan name taken", func(t *testing.T) {
		mClan := &mockClanRepo{findByNameClan: &models.Clan{}}
		mWallet := &mockWalletRepoClan{}
		mTx := &mockTxManagerClan{}

		svc := NewClanService(mClan, mWallet, mTx)
		clan, err := svc.CreateClan(ctx, userID, "Test Clan", "Slogan", "url")
		
		assert.Error(t, err)
		assert.Equal(t, "ERR_CLAN_NAME_TAKEN", err.Error())
		assert.Nil(t, clan)
	})

	t.Run("Insufficient balance", func(t *testing.T) {
		mClan := &mockClanRepo{findByNameErr: gorm.ErrRecordNotFound}
		mWallet := &mockWalletRepoClan{wallet: &models.UserWallet{Balance: 100}}
		mTx := &mockTxManagerClan{}

		svc := NewClanService(mClan, mWallet, mTx)
		clan, err := svc.CreateClan(ctx, userID, "Test Clan", "Slogan", "url")
		
		assert.Error(t, err)
		assert.Equal(t, "insufficient balance", err.Error())
		assert.Nil(t, clan)
	})

	t.Run("Wallet not found", func(t *testing.T) {
		mClan := &mockClanRepo{findByNameErr: gorm.ErrRecordNotFound}
		mWallet := &mockWalletRepoClan{walletErr: gorm.ErrRecordNotFound}
		mTx := &mockTxManagerClan{}

		svc := NewClanService(mClan, mWallet, mTx)
		clan, err := svc.CreateClan(ctx, userID, "Test Clan", "Slogan", "url")
		
		assert.Error(t, err)
		assert.Equal(t, "wallet not found", err.Error())
		assert.Nil(t, clan)
	})
	
	t.Run("DB error on create", func(t *testing.T) {
		mClan := &mockClanRepo{findByNameErr: gorm.ErrRecordNotFound, createErr: errors.New("db err")}
		mWallet := &mockWalletRepoClan{wallet: &models.UserWallet{Balance: 1000}}
		mTx := &mockTxManagerClan{}

		svc := NewClanService(mClan, mWallet, mTx)
		clan, err := svc.CreateClan(ctx, userID, "Test Clan", "Slogan", "url")
		
		assert.Error(t, err)
		assert.Equal(t, "db err", err.Error())
		assert.Nil(t, clan)
	})
}
