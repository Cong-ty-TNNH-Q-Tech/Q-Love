// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type VoucherRepository interface {
	GetAvailableVoucher(ctx context.Context, brand string, valueXu int) (*models.Voucher, error)
	MarkAsClaimed(ctx context.Context, voucherID uuid.UUID, userID uuid.UUID) error
	Create(ctx context.Context, voucher *models.Voucher) error
	FindAll(ctx context.Context, limit, offset int) ([]models.Voucher, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type voucherRepository struct {
	db *gorm.DB
}

func NewVoucherRepository(db *gorm.DB) VoucherRepository {
	return &voucherRepository{db: db}
}

func (r *voucherRepository) GetAvailableVoucher(ctx context.Context, brand string, valueXu int) (*models.Voucher, error) {
	var voucher models.Voucher
	db := GetDB(ctx, r.db)
	query := db.WithContext(ctx).
		Where("brand = ? AND value_xu = ? AND status = 'available' AND expires_at > ?", brand, valueXu, time.Now()).
		Order("expires_at ASC") // Use earliest expiring voucher first

	if db.Dialector.Name() != "sqlite" {
		query = query.Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"})
	}

	err := query.First(&voucher).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("không có voucher khả dụng")
	}
	return &voucher, err
}

func (r *voucherRepository) MarkAsClaimed(ctx context.Context, voucherID uuid.UUID, userID uuid.UUID) error {
	db := GetDB(ctx, r.db)
	// Update voucher status
	err := db.WithContext(ctx).Model(&models.Voucher{}).
		Where("id = ?", voucherID).
		Update("status", "claimed").Error
	if err != nil {
		return err
	}

	// Create user_voucher mapping
	userVoucher := &models.UserVoucher{
		ID:        uuid.New(),
		UserID:    userID,
		VoucherID: voucherID,
		ClaimedAt: time.Now(),
	}
	return db.WithContext(ctx).Create(userVoucher).Error
}

func (r *voucherRepository) Create(ctx context.Context, voucher *models.Voucher) error {
	db := GetDB(ctx, r.db)
	return db.WithContext(ctx).Create(voucher).Error
}

func (r *voucherRepository) FindAll(ctx context.Context, limit, offset int) ([]models.Voucher, error) {
	var vouchers []models.Voucher
	db := GetDB(ctx, r.db)
	err := db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&vouchers).Error
	return vouchers, err
}

func (r *voucherRepository) Delete(ctx context.Context, id uuid.UUID) error {
	db := GetDB(ctx, r.db)
	// Only delete if it's still available (soft delete)
	result := db.WithContext(ctx).
		Where("id = ? AND status = 'available'", id).
		Delete(&models.Voucher{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("không tìm thấy voucher hoặc voucher đã được sử dụng")
	}
	return nil
}
