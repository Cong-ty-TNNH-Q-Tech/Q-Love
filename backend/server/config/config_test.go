// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("R2_ACCOUNT_ID", "test_id")
	os.Setenv("PORT", "4000")
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Clearenv()

	cfg := LoadConfig()
	if cfg.R2AccountID != "test_id" {
		t.Errorf("Expected test_id, got %s", cfg.R2AccountID)
	}
	if cfg.Port != "4000" {
		t.Errorf("Expected 4000, got %s", cfg.Port)
	}
}

func TestGetEnv(t *testing.T) {
	os.Clearenv()
	val := getEnv("NON_EXISTENT", "default_val")
	if val != "default_val" {
		t.Errorf("Expected default_val, got %s", val)
	}

	os.Setenv("EXISTS", "yes")
	val2 := getEnv("EXISTS", "default_val")
	if val2 != "yes" {
		t.Errorf("Expected yes, got %s", val2)
	}
}
