// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateTOTPSecret(t *testing.T) {
	secret, err := generateTOTPSecret()
	assert.NoError(t, err)
	assert.Len(t, secret, 40) // 20 bytes = 40 hex chars
}

func TestGenerateTOTPSecret_Uniqueness(t *testing.T) {
	s1, _ := generateTOTPSecret()
	s2, _ := generateTOTPSecret()
	assert.NotEqual(t, s1, s2)
}

func TestGenerateTOTP_ValidOutput(t *testing.T) {
	secret, _ := generateTOTPSecret()
	code, err := GenerateTOTP(secret, time.Now())
	assert.NoError(t, err)
	assert.Len(t, code, 6) // Always 6 digits
}

func TestGenerateTOTP_InvalidSecret(t *testing.T) {
	_, err := GenerateTOTP("not-hex", time.Now())
	assert.Error(t, err)
}

func TestGenerateTOTP_Deterministic(t *testing.T) {
	secret := "0123456789abcdef0123456789abcdef01234567"
	fixedTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	code1, _ := GenerateTOTP(secret, fixedTime)
	code2, _ := GenerateTOTP(secret, fixedTime)
	assert.Equal(t, code1, code2)
}

func TestGenerateTOTP_DifferentTimeSteps(t *testing.T) {
	secret := "0123456789abcdef0123456789abcdef01234567"
	t1 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 := t1.Add(31 * time.Second) // Next time step

	code1, _ := GenerateTOTP(secret, t1)
	code2, _ := GenerateTOTP(secret, t2)
	assert.NotEqual(t, code1, code2)
}

func TestValidateTOTP_CurrentCode(t *testing.T) {
	secret, _ := generateTOTPSecret()
	code, _ := GenerateTOTP(secret, time.Now())

	assert.True(t, validateTOTP(secret, code))
}

func TestValidateTOTP_InvalidCode(t *testing.T) {
	secret, _ := generateTOTPSecret()
	assert.False(t, validateTOTP(secret, "000000"))
}

func TestValidateTOTP_InvalidSecret(t *testing.T) {
	assert.False(t, validateTOTP("invalid-hex", "123456"))
}

func TestValidateTOTP_EmptyToken(t *testing.T) {
	secret, _ := generateTOTPSecret()
	assert.False(t, validateTOTP(secret, ""))
}
