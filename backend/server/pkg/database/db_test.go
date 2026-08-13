// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package database

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/config"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestNewPostgresDB_EmptyDSN(t *testing.T) {
	logger.InitLogger("development", "") // prevent nil pointer

	cfg := &config.Config{
		DatabaseDSN: "",
	}

	db, err := NewPostgresDB(cfg)
	if err == nil {
		t.Error("Expected error for empty DSN")
	}
	if db != nil {
		t.Error("Expected nil db for empty DSN")
	}
}

func TestNewPostgresDB_InvalidDSN(t *testing.T) {
	logger.InitLogger("development", "")

	cfg := &config.Config{
		DatabaseDSN: "invalid-dsn",
	}

	db, err := NewPostgresDB(cfg)
	if err == nil {
		t.Error("Expected error for invalid DSN connection attempt")
	}
	if db != nil {
		t.Error("Expected nil db for invalid DSN connection attempt")
	}
}

func TestNewPostgresDB_SkipDSN(t *testing.T) {
	logger.InitLogger("development", "")

	cfg := &config.Config{
		DatabaseDSN: "skip",
	}

	db, err := NewPostgresDB(cfg)
	if err != nil {
		t.Errorf("Expected no error for skip DSN, got %v", err)
	}
	if db != nil {
		t.Error("Expected nil db for skip DSN")
	}
}

func TestNewPostgresDB_Success(t *testing.T) {
	logger.InitLogger("development", "")

	// Create a mock DB
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	// Mock GormOpen and PostgresOpen
	oldGormOpen := GormOpen
	oldPostgresOpen := PostgresOpen
	defer func() {
		GormOpen = oldGormOpen
		PostgresOpen = oldPostgresOpen
	}()

	GormOpen = func(dialector gorm.Dialector, opts ...gorm.Option) (*gorm.DB, error) {
		return gorm.Open(postgres.New(postgres.Config{Conn: mockDB}), opts...)
	}
	PostgresOpen = func(dsn string) gorm.Dialector {
		return postgres.New(postgres.Config{Conn: mockDB})
	}

	cfg := &config.Config{
		DatabaseDSN: "mock-dsn",
	}

	db, err := NewPostgresDB(cfg)
	if err != nil {
		t.Errorf("Expected no error for success case, got %v", err)
	}
	if db == nil {
		t.Error("Expected valid db for success case")
	}
}
