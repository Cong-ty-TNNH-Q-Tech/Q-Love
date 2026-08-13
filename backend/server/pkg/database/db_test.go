// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package database

import (
	"testing"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/config"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/pkg/logger"
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
