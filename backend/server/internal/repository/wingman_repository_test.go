// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package repository

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
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
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, nil, err
	}
	return gdb, mock, nil
}

func TestWingmanRepository_CreateReferral(t *testing.T) {
	db, mock, err := setupTestDB()
	if err != nil {
		t.Fatal(err)
	}

	repo := NewWingmanRepository(db)
	referral := &models.WingmanReferral{
		ID:        uuid.New(),
		WingmanID: uuid.New(),
		Target1ID: uuid.New(),
		Target2ID: uuid.New(),
		Status:    "pending",
		DeepLink:  "link",
		ExpiresAt: time.Now(),
	}

	mock.ExpectExec(`INSERT INTO "wingman_referrals"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.CreateReferral(context.Background(), referral)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestWingmanRepository_GetReferralByID(t *testing.T) {
	db, mock, err := setupTestDB()
	if err != nil {
		t.Fatal(err)
	}

	repo := NewWingmanRepository(db)
	refID := uuid.New()

	mock.ExpectQuery(`SELECT .* FROM "wingman_referrals"`).
		WithArgs(refID, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "status"}).AddRow(refID, "pending"))

	referral, err := repo.GetReferralByID(context.Background(), refID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if referral.ID != refID {
		t.Errorf("Expected %s, got %s", refID, referral.ID)
	}
}

func TestWingmanRepository_GetReferralByID_Error(t *testing.T) {
	db, mock, err := setupTestDB()
	if err != nil {
		t.Fatal(err)
	}

	repo := NewWingmanRepository(db)
	refID := uuid.New()

	mock.ExpectQuery(`SELECT .* FROM "wingman_referrals"`).
		WithArgs(refID, sqlmock.AnyArg()).
		WillReturnError(gorm.ErrRecordNotFound)

	_, err = repo.GetReferralByID(context.Background(), refID)
	if err == nil {
		t.Error("Expected error")
	}
}

func TestWingmanRepository_UpdateReferral(t *testing.T) {
	db, mock, err := setupTestDB()
	if err != nil {
		t.Fatal(err)
	}

	repo := NewWingmanRepository(db)
	referral := &models.WingmanReferral{
		ID:        uuid.New(),
		Status:    "matched",
	}

	mock.ExpectExec(`UPDATE "wingman_referrals"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateReferral(context.Background(), referral)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
