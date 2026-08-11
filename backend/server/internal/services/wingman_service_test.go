package services

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	
	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	
	gdb, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	
	if err != nil {
		return nil, nil, err
	}
	
	return gdb, mock, nil
}

func TestWingmanService_CreateReferral(t *testing.T) {
	db, mock, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to open mock sql db, got error: %v", err)
	}

	service := NewWingmanService(db)

	wingmanID := uuid.New()
	target1ID := uuid.New()
	target2ID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "wingman_referrals"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	referral, err := service.CreateReferral(context.Background(), wingmanID, target1ID, target2ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if referral == nil {
		t.Fatal("Expected referral to be returned")
	}

	if referral.WingmanID != wingmanID {
		t.Errorf("Expected wingmanID %s, got %s", wingmanID, referral.WingmanID)
	}
}

func TestWingmanService_AcceptReferral_Success(t *testing.T) {
	db, mock, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to open mock sql db, got error: %v", err)
	}

	service := NewWingmanService(db)
	refID := uuid.New()
	target1ID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT .* FROM "wingman_referrals"`).
		WithArgs(refID, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "status", "expires_at", "target1_id", "target2_id"}).
			AddRow(refID, "pending", time.Now().Add(1*time.Hour), target1ID, uuid.New()))
	
	mock.ExpectExec(`UPDATE "wingman_referrals"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	referral, err := service.AcceptReferral(context.Background(), refID, target1ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if referral.Status != "matched" {
		t.Errorf("Expected status matched, got %s", referral.Status)
	}
}

func TestWingmanService_ProcessCommission_Success(t *testing.T) {
	db, mock, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to open mock sql db, got error: %v", err)
	}

	service := NewWingmanService(db)
	refID := uuid.New()
	wingmanID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT .* FROM "wingman_referrals"`).
		WithArgs(refID, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "status", "wingman_id"}).
			AddRow(refID, "matched", wingmanID))
	
	// mock FirstOrCreate UserWallet
	mock.ExpectQuery(`SELECT .* FROM "user_wallets"`).
		WithArgs(wingmanID, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "balance"}).
			AddRow(wingmanID, 0))
	
	// save wallet
	mock.ExpectExec(`UPDATE "user_wallets"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	
	// Create transaction
	mock.ExpectExec(`INSERT INTO "wallet_transactions"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	
	// Save referral
	mock.ExpectExec(`UPDATE "wingman_referrals"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err = service.ProcessCommission(context.Background(), refID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestWingmanService_AcceptReferral_Errors(t *testing.T) {
	db, mock, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to open mock sql db, got error: %v", err)
	}

	service := NewWingmanService(db)
	refID := uuid.New()
	target1ID := uuid.New()
	invalidUserID := uuid.New()

	// 1. Referral not found
	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT .* FROM "wingman_referrals"`).
		WithArgs(refID, sqlmock.AnyArg()).
		WillReturnError(gorm.ErrRecordNotFound)
	mock.ExpectRollback()

	_, err = service.AcceptReferral(context.Background(), refID, target1ID)
	if err == nil || err.Error() != "referral not found" {
		t.Errorf("Expected referral not found error, got %v", err)
	}

	// 2. Referral not pending
	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT .* FROM "wingman_referrals"`).
		WithArgs(refID, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "status"}).
			AddRow(refID, "matched"))
	mock.ExpectRollback()

	_, err = service.AcceptReferral(context.Background(), refID, target1ID)
	if err == nil || err.Error() != "referral is no longer pending" {
		t.Errorf("Expected referral is no longer pending error, got %v", err)
	}

	// 3. Referral expired
	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT .* FROM "wingman_referrals"`).
		WithArgs(refID, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "status", "expires_at"}).
			AddRow(refID, "pending", time.Now().Add(-1*time.Hour)))
	mock.ExpectRollback()

	_, err = service.AcceptReferral(context.Background(), refID, target1ID)
	if err == nil || err.Error() != "referral link expired" {
		t.Errorf("Expected referral link expired error, got %v", err)
	}

	// 4. Invalid user
	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT .* FROM "wingman_referrals"`).
		WithArgs(refID, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "status", "expires_at", "target1_id", "target2_id"}).
			AddRow(refID, "pending", time.Now().Add(1*time.Hour), target1ID, uuid.New()))
	mock.ExpectRollback()

	_, err = service.AcceptReferral(context.Background(), refID, invalidUserID)
	if err == nil || err.Error() != "user is not part of this referral" {
		t.Errorf("Expected user is not part of this referral error, got %v", err)
	}
}

func TestWingmanService_ProcessCommission_Errors(t *testing.T) {
	db, mock, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to open mock sql db, got error: %v", err)
	}

	service := NewWingmanService(db)
	refID := uuid.New()
	wingmanID := uuid.New()

	// 1. Invalid status
	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT .* FROM "wingman_referrals"`).
		WithArgs(refID, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "status", "wingman_id"}).
			AddRow(refID, "pending", wingmanID))
	mock.ExpectRollback()

	err = service.ProcessCommission(context.Background(), refID)
	if err == nil || err.Error() != "invalid status for commission" {
		t.Errorf("Expected invalid status for commission error, got %v", err)
	}
}
