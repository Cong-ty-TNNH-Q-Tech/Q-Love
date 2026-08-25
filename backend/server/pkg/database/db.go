// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package database

import (
	"fmt"
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/config"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	GormOpen     = gorm.Open
	PostgresOpen = postgres.Open
)

func NewPostgresDB(cfg *config.Config) (*gorm.DB, error) {
	if cfg.DatabaseDSN == "" {
		return nil, fmt.Errorf("DatabaseDSN is required in configuration")
	}
	if cfg.DatabaseDSN == "skip" {
		logger.Log.Info("Skipping database connection for testing")
		return nil, nil
	}

	db, err := GormOpen(PostgresOpen(cfg.DatabaseDSN), &gorm.Config{})
	if err != nil {
		logger.Log.Error("Failed to connect to database", zap.Error(err))
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Log.Error("Failed to get database connection pool", zap.Error(err))
		return nil, err
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	logger.Log.Info("Successfully connected to PostgreSQL database")



	return db, nil
}
