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

	// test invalid sentry DSN url
	InitLogger("development", "https://invalid@sentry.io/123")
	if Log == nil {
		t.Error("Expected Log to not be nil")
	}

	// test invalid sentry DSN string
	InitLogger("development", "invalid-dsn")
	if Log == nil {
		t.Error("Expected Log to not be nil")
	}
}
