// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package repository

import (
	"context"
	"testing"
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupVoucherTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	db.AutoMigrate(&models.Voucher{}, &models.UserVoucher{})
	return db
}

func TestVoucherRepository_Create(t *testing.T) {
	db := setupVoucherTestDB(t)
	repo := NewVoucherRepository(db)

	voucher := &models.Voucher{
		ID:        uuid.New(),
		Brand:     "Highlands",
		Code:      "HL-123",
		ValueXu:   100,
		Status:    "available",
		ExpiresAt: time.Now(),
		CreatedAt: time.Now(),
	}

	err := repo.Create(context.Background(), voucher)
	assert.NoError(t, err)

	// Check if created
	var count int64
	db.Model(&models.Voucher{}).Count(&count)
	assert.Equal(t, int64(1), count)
}

func TestVoucherRepository_FindAll(t *testing.T) {
	db := setupVoucherTestDB(t)
	repo := NewVoucherRepository(db)

	for i := 0; i < 3; i++ {
		db.Create(&models.Voucher{
			ID:        uuid.New(),
			Brand:     "Highlands",
			Code:      "HL-" + uuid.New().String(),
			ValueXu:   100,
			Status:    "available",
			ExpiresAt: time.Now(),
		})
	}

	vouchers, err := repo.FindAll(context.Background(), 10, 0)
	assert.NoError(t, err)
	assert.Len(t, vouchers, 3)
}

func TestVoucherRepository_FindAvailable(t *testing.T) {
	db := setupVoucherTestDB(t)
	repo := NewVoucherRepository(db)

	db.Create(&models.Voucher{
		ID:        uuid.New(),
		Brand:     "Highlands",
		Code:      "HL-1",
		ValueXu:   100,
		Status:    "available",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	})
	db.Create(&models.Voucher{
		ID:        uuid.New(),
		Brand:     "Highlands",
		Code:      "HL-2",
		ValueXu:   100,
		Status:    "claimed",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	})

	vouchers, err := repo.FindAvailable(context.Background(), 10, 0)
	assert.NoError(t, err)
	assert.Len(t, vouchers, 1)
	assert.Equal(t, "HL-1", vouchers[0].Code)
}

func TestVoucherRepository_FindByID(t *testing.T) {
	db := setupVoucherTestDB(t)
	repo := NewVoucherRepository(db)

	id := uuid.New()
	db.Create(&models.Voucher{
		ID:        id,
		Brand:     "Highlands",
		Code:      "HL-1",
		ValueXu:   100,
		Status:    "available",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	})

	voucher, err := repo.FindByID(context.Background(), id)
	assert.NoError(t, err)
	assert.Equal(t, id, voucher.ID)
}

func TestVoucherRepository_Claim(t *testing.T) {
	db := setupVoucherTestDB(t)
	repo := NewVoucherRepository(db)

	userID := uuid.New()
	voucherID := uuid.New()

	db.Create(&models.Voucher{
		ID:        voucherID,
		Brand:     "Highlands",
		Code:      "HL-1",
		ValueXu:   100,
		Status:    "available",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	})

	err := repo.Claim(context.Background(), userID, voucherID)
	assert.NoError(t, err)

	// Verify status updated
	var voucher models.Voucher
	db.First(&voucher, voucherID)
	assert.Equal(t, "claimed", voucher.Status)

	// Verify UserVoucher created
	var uv models.UserVoucher
	db.First(&uv, "voucher_id = ?", voucherID)
	assert.Equal(t, userID, uv.UserID)
}

func TestVoucherRepository_FindUserVouchers(t *testing.T) {
	db := setupVoucherTestDB(t)
	repo := NewVoucherRepository(db)

	userID := uuid.New()
	voucherID := uuid.New()

	db.Create(&models.Voucher{
		ID:        voucherID,
		Brand:     "Highlands",
		Code:      "HL-1",
		ValueXu:   100,
		Status:    "claimed",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	})

	db.Create(&models.UserVoucher{
		ID:        uuid.New(),
		UserID:    userID,
		VoucherID: voucherID,
		ClaimedAt: time.Now(),
	})

	vouchers, err := repo.FindUserVouchers(context.Background(), userID)
	assert.NoError(t, err)
	assert.Len(t, vouchers, 1)
	assert.Equal(t, voucherID, vouchers[0].VoucherID)
}

func TestVoucherRepository_Errors(t *testing.T) {
	db := setupVoucherTestDB(t)
	// Force close db to test errors
	sqlDB, _ := db.DB()
	sqlDB.Close()
	
	repo := NewVoucherRepository(db)
	
	err := repo.Create(context.Background(), &models.Voucher{})
	assert.Error(t, err)
	
	_, err = repo.FindAll(context.Background(), 10, 0)
	assert.Error(t, err)
	
	_, err = repo.FindAvailable(context.Background(), 10, 0)
	assert.Error(t, err)
	
	_, err = repo.FindByID(context.Background(), uuid.New())
	assert.Error(t, err)
	
	err = repo.Claim(context.Background(), uuid.New(), uuid.New())
	assert.Error(t, err)
	
	_, err = repo.FindUserVouchers(context.Background(), uuid.New())
	assert.Error(t, err)
}
