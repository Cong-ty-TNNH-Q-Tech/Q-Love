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
	mock.ExpectQuery(`INSERT INTO "wingman_referrals"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
	mock.ExpectCommit()

	referral, err := service.CreateReferral(context.Background(), wingmanID, target1ID, target2ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
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
	mock.ExpectQuery(`SELECT \* FROM "wingman_referrals" WHERE id = \$1`).
		WithArgs(refID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "status", "expires_at", "target1_id", "target2_id"}).
			AddRow(refID, "pending", time.Now().Add(1*time.Hour), target1ID, uuid.New()))
	
	mock.ExpectExec(`UPDATE "wingman_referrals" SET`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	referral, err := service.AcceptReferral(context.Background(), refID, target1ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
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
	mock.ExpectQuery(`SELECT \* FROM "wingman_referrals" WHERE id = \$1`).
		WithArgs(refID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "status", "wingman_id"}).
			AddRow(refID, "matched", wingmanID))
	
	// mock FirstOrCreate UserWallet
	mock.ExpectQuery(`SELECT \* FROM "user_wallets" WHERE "user_wallets"."user_id" = \$1`).
		WithArgs(wingmanID).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "balance"}).
			AddRow(wingmanID, 0))
	
	// save wallet
	mock.ExpectExec(`UPDATE "user_wallets" SET`).WillReturnResult(sqlmock.NewResult(1, 1))
	
	// Create transaction
	mock.ExpectQuery(`INSERT INTO "wallet_transactions"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
	
	// Save referral
	mock.ExpectExec(`UPDATE "wingman_referrals" SET`).WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err = service.ProcessCommission(context.Background(), refID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
