// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"context"
	"errors"
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/repository"
	"github.com/google/uuid"
)

type VoucherService interface {
	RedeemVoucher(ctx context.Context, userID uuid.UUID, brand string, valueXu int) (*models.Voucher, error)
	CreateVoucher(ctx context.Context, req CreateVoucherRequest) error
	GetVouchers(ctx context.Context, limit, offset int) ([]models.Voucher, error)
	DeleteVoucher(ctx context.Context, id uuid.UUID) error
}

type voucherService struct {
	voucherRepo repository.VoucherRepository
	walletRepo  repository.WalletRepository
	txManager   repository.TransactionManager
}

type CreateVoucherRequest struct {
	Brand     string
	Code      string
	ValueXu   int
	ExpiresAt time.Time
}

func NewVoucherService(vRepo repository.VoucherRepository, wRepo repository.WalletRepository, txManager repository.TransactionManager) VoucherService {
	return &voucherService{
		voucherRepo: vRepo,
		walletRepo:  wRepo,
		txManager:   txManager,
	}
}

func (s *voucherService) RedeemVoucher(ctx context.Context, userID uuid.UUID, brand string, valueXu int) (*models.Voucher, error) {
	var claimedVoucher *models.Voucher

	err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		// 1. Get Wallet For Update
		wallet, err := s.walletRepo.GetWalletForUpdate(txCtx, userID)
		if err != nil {
			return errors.New("lỗi truy cập ví")
		}

		if wallet.Balance < float64(valueXu) {
			return errors.New("số dư xu không đủ")
		}

		// 2. Lock and Get Available Voucher
		voucher, err := s.voucherRepo.GetAvailableVoucher(txCtx, brand, valueXu)
		if err != nil {
			return err
		}

		// 3. Deduct Balance
		err = s.walletRepo.UpdateBalance(txCtx, userID, -float64(valueXu))
		if err != nil {
			return errors.New("lỗi trừ xu")
		}

		// 4. Record Transaction
		txn := &models.WalletTransaction{
			ID:          uuid.New(),
			UserID:      userID,
			Amount:      -float64(valueXu),
			Type:        "redeem_voucher",
			ReferenceID: voucher.ID,
			CreatedAt:   time.Now(),
		}
		if err := s.walletRepo.CreateTransaction(txCtx, txn); err != nil {
			return errors.New("lỗi lưu lịch sử giao dịch")
		}

		// 5. Mark as Claimed
		if err := s.voucherRepo.MarkAsClaimed(txCtx, voucher.ID, userID); err != nil {
			return errors.New("lỗi cập nhật trạng thái voucher")
		}

		claimedVoucher = voucher
		return nil
	})

	return claimedVoucher, err
}

func (s *voucherService) CreateVoucher(ctx context.Context, req CreateVoucherRequest) error {
	if req.ValueXu <= 0 {
		return errors.New("giá trị xu phải lớn hơn 0")
	}
	v := &models.Voucher{
		ID:        uuid.New(),
		Brand:     req.Brand,
		Code:      req.Code,
		ValueXu:   req.ValueXu,
		Status:    "available",
		ExpiresAt: req.ExpiresAt,
		CreatedAt: time.Now(),
	}
	return s.voucherRepo.Create(ctx, v)
}

func (s *voucherService) GetVouchers(ctx context.Context, limit, offset int) ([]models.Voucher, error) {
	return s.voucherRepo.FindAll(ctx, limit, offset)
}

func (s *voucherService) DeleteVoucher(ctx context.Context, id uuid.UUID) error {
	return s.voucherRepo.Delete(ctx, id)
}
