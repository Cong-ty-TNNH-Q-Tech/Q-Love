// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
)

type WalletRepository interface {
	AddCommission(ctx context.Context, userID uuid.UUID, amount float64) error
	CreateTransaction(ctx context.Context, txn *models.WalletTransaction) error
	GetWalletForUpdate(ctx context.Context, userID uuid.UUID) (*models.UserWallet, error)
	UpdateBalance(ctx context.Context, userID uuid.UUID, delta float64) error
	CheckTransactionExists(ctx context.Context, txID uuid.UUID) (bool, error)
}

type walletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) WalletRepository {
	return &walletRepository{db: db}
}

func (r *walletRepository) AddCommission(ctx context.Context, userID uuid.UUID, amount float64) error {
	var wallet models.UserWallet
	db := GetDB(ctx, r.db)
	
	if err := db.WithContext(ctx).FirstOrCreate(&wallet, models.UserWallet{UserID: userID}).Error; err != nil {
		return err
	}

	wallet.Balance += amount
	return db.WithContext(ctx).Save(&wallet).Error
}

func (r *walletRepository) GetWalletForUpdate(ctx context.Context, userID uuid.UUID) (*models.UserWallet, error) {
	var wallet models.UserWallet
	db := GetDB(ctx, r.db)
	if err := db.WithContext(ctx).Set("gorm:query_option", "FOR UPDATE").FirstOrCreate(&wallet, models.UserWallet{UserID: userID}).Error; err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *walletRepository) UpdateBalance(ctx context.Context, userID uuid.UUID, delta float64) error {
	db := GetDB(ctx, r.db)
	return db.WithContext(ctx).
		Model(&models.UserWallet{}).
		Where("user_id = ?", userID).
		UpdateColumn("balance", gorm.Expr("balance + ?", delta)).Error
}

func (r *walletRepository) CheckTransactionExists(ctx context.Context, txID uuid.UUID) (bool, error) {
	var count int64
	db := GetDB(ctx, r.db)
	err := db.WithContext(ctx).
		Model(&models.WalletTransaction{}).
		Where("id = ?", txID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *walletRepository) CreateTransaction(ctx context.Context, txn *models.WalletTransaction) error {
	return GetDB(ctx, r.db).WithContext(ctx).Create(txn).Error
}
