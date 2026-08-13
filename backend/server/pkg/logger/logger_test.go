// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package logger

import (
	"testing"
)

func TestInitLogger(t *testing.T) {
	// Test development env
	InitLogger("development", "")
	if Log == nil {
		t.Error("Expected Log to be initialized for development")
	}

	// Test production env without sentry
	InitLogger("production", "")
	if Log == nil {
		t.Error("Expected Log to be initialized for production")
	}

	// Test with dummy Sentry DSN (it should fail to init sentry but not panic)
	InitLogger("development", "http://public@localhost/1")
	if Log == nil {
		t.Error("Expected Log to be initialized even if sentry fails or succeeds")
	}
}
