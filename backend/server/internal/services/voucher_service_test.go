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

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockVoucherRepo struct {
	availableVoucher *models.Voucher
	getErr           error
	markErr          error
	createErr        error
	deleteErr        error
}

func (m *mockVoucherRepo) GetAvailableVoucher(ctx context.Context, brand string, valueXu int) (*models.Voucher, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	return m.availableVoucher, nil
}
func (m *mockVoucherRepo) MarkAsClaimed(ctx context.Context, voucherID uuid.UUID, userID uuid.UUID) error {
	return m.markErr
}
func (m *mockVoucherRepo) Create(ctx context.Context, voucher *models.Voucher) error {
	return m.createErr
}
func (m *mockVoucherRepo) FindAll(ctx context.Context, limit, offset int) ([]models.Voucher, error) {
	return nil, nil
}
func (m *mockVoucherRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return m.deleteErr
}

type mockWalletRepoForVoucher struct {
	balance float64
	updateErr error
	txnErr    error
}
func (m *mockWalletRepoForVoucher) GetWalletForUpdate(ctx context.Context, userID uuid.UUID) (*models.UserWallet, error) {
	return &models.UserWallet{Balance: m.balance}, nil
}
func (m *mockWalletRepoForVoucher) UpdateBalance(ctx context.Context, userID uuid.UUID, delta float64) error {
	return m.updateErr
}
func (m *mockWalletRepoForVoucher) AddCommission(ctx context.Context, userID uuid.UUID, amount float64) error {
	return nil
}
func (m *mockWalletRepoForVoucher) CreateTransaction(ctx context.Context, txn *models.WalletTransaction) error {
	return m.txnErr
}
func (m *mockWalletRepoForVoucher) CheckTransactionExists(ctx context.Context, txID uuid.UUID) (bool, error) { return false, nil }

type mockTxManagerForVoucher struct{}

func (m *mockTxManagerForVoucher) WithTransaction(ctx context.Context, fn func(txCtx context.Context) error, opts ...*sql.TxOptions) error {
	return fn(ctx)
}

func TestVoucherService_RedeemVoucher_Success(t *testing.T) {
	vRepo := &mockVoucherRepo{
		availableVoucher: &models.Voucher{ID: uuid.New(), Brand: "Highlands", ValueXu: 100},
	}
	wRepo := &mockWalletRepoForVoucher{balance: 150}
	tx := &mockTxManagerForVoucher{}

	svc := NewVoucherService(vRepo, wRepo, tx)
	v, err := svc.RedeemVoucher(context.Background(), uuid.New(), "Highlands", 100)
	assert.NoError(t, err)
	assert.NotNil(t, v)
}

func TestVoucherService_RedeemVoucher_InsufficientFunds(t *testing.T) {
	vRepo := &mockVoucherRepo{}
	wRepo := &mockWalletRepoForVoucher{balance: 50}
	tx := &mockTxManagerForVoucher{}

	svc := NewVoucherService(vRepo, wRepo, tx)
	_, err := svc.RedeemVoucher(context.Background(), uuid.New(), "Highlands", 100)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "không đủ")
}

func TestVoucherService_RedeemVoucher_NoVoucher(t *testing.T) {
	vRepo := &mockVoucherRepo{getErr: errors.New("không có voucher")}
	wRepo := &mockWalletRepoForVoucher{balance: 150}
	tx := &mockTxManagerForVoucher{}

	svc := NewVoucherService(vRepo, wRepo, tx)
	_, err := svc.RedeemVoucher(context.Background(), uuid.New(), "Highlands", 100)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "không có voucher")
}

func TestVoucherService_CreateVoucher(t *testing.T) {
	svc := NewVoucherService(&mockVoucherRepo{}, nil, nil)
	err := svc.CreateVoucher(context.Background(), CreateVoucherRequest{
		Brand:   "Highlands",
		Code:    "HL-12345",
		ValueXu: 50,
	})
	assert.NoError(t, err)

	err = svc.CreateVoucher(context.Background(), CreateVoucherRequest{
		Brand:   "Highlands",
		Code:    "HL-12345",
		ValueXu: -10,
	})
	assert.Error(t, err)
}

func TestVoucherService_GetVouchers(t *testing.T) {
	svc := NewVoucherService(&mockVoucherRepo{}, nil, nil)
	v, err := svc.GetVouchers(context.Background(), 10, 0)
	assert.NoError(t, err)
	assert.Nil(t, v)
}

func TestVoucherService_DeleteVoucher(t *testing.T) {
	svc := NewVoucherService(&mockVoucherRepo{}, nil, nil)
	err := svc.DeleteVoucher(context.Background(), uuid.New())
	assert.NoError(t, err)
}

func TestVoucherService_RedeemVoucher_UpdateWalletErr(t *testing.T) {
	vRepo := &mockVoucherRepo{
		availableVoucher: &models.Voucher{ID: uuid.New(), Brand: "Highlands", ValueXu: 100},
	}
	wRepo := &mockWalletRepoForVoucher{balance: 150, updateErr: errors.New("lỗi trừ xu")}
	tx := &mockTxManagerForVoucher{}

	svc := NewVoucherService(vRepo, wRepo, tx)
	_, err := svc.RedeemVoucher(context.Background(), uuid.New(), "Highlands", 100)
	assert.Error(t, err)
}

func TestVoucherService_RedeemVoucher_TxnErr(t *testing.T) {
	vRepo := &mockVoucherRepo{
		availableVoucher: &models.Voucher{ID: uuid.New(), Brand: "Highlands", ValueXu: 100},
	}
	wRepo := &mockWalletRepoForVoucher{balance: 150, txnErr: errors.New("lỗi lưu lịch sử")}
	tx := &mockTxManagerForVoucher{}

	svc := NewVoucherService(vRepo, wRepo, tx)
	_, err := svc.RedeemVoucher(context.Background(), uuid.New(), "Highlands", 100)
	assert.Error(t, err)
}

func TestVoucherService_RedeemVoucher_MarkErr(t *testing.T) {
	vRepo := &mockVoucherRepo{
		availableVoucher: &models.Voucher{ID: uuid.New(), Brand: "Highlands", ValueXu: 100},
		markErr:          errors.New("lỗi cập nhật"),
	}
	wRepo := &mockWalletRepoForVoucher{balance: 150}
	tx := &mockTxManagerForVoucher{}

	svc := NewVoucherService(vRepo, wRepo, tx)
	_, err := svc.RedeemVoucher(context.Background(), uuid.New(), "Highlands", 100)
	assert.Error(t, err)
}

type mockTxManagerWithErrorForVoucher struct{}
func (m *mockTxManagerWithErrorForVoucher) WithTransaction(ctx context.Context, fn func(txCtx context.Context) error, opts ...*sql.TxOptions) error {
	return errors.New("lỗi truy cập ví")
}

func TestVoucherService_RedeemVoucher_TxFail(t *testing.T) {
	svc := NewVoucherService(nil, nil, &mockTxManagerWithErrorForVoucher{})
	_, err := svc.RedeemVoucher(context.Background(), uuid.New(), "Highlands", 100)
	assert.Error(t, err)
}
