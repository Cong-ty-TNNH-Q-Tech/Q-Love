// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSpiritualService_CalculateZodiac(t *testing.T) {
	service := NewSpiritualService()

	// Aries: March 21 - April 19
	dobAries, _ := time.Parse("2006-01-02", "1998-03-25")
	assert.Equal(t, "Aries", service.CalculateZodiac(dobAries))

	// Taurus: April 20 - May 20
	dobTaurus, _ := time.Parse("2006-01-02", "1995-05-10")
	assert.Equal(t, "Taurus", service.CalculateZodiac(dobTaurus))

	// Scorpio: October 24 - November 21
	dobScorpio, _ := time.Parse("2006-01-02", "1999-11-05")
	assert.Equal(t, "Scorpio", service.CalculateZodiac(dobScorpio))
}

func TestSpiritualService_CalculateNumerology(t *testing.T) {
	service := NewSpiritualService()

	// 15/04/1998 = 1+5+0+4+1+9+9+8 = 37 => 3+7 = 10 => 1+0 = 1
	dob1, _ := time.Parse("2006-01-02", "1998-04-15")
	assert.Equal(t, 1, service.CalculateNumerology(dob1))

	// 22/11/1990 = 2+2+1+1+1+9+9+0 = 25 => 2+5 = 7
	dob2, _ := time.Parse("2006-01-02", "1990-11-22")
	assert.Equal(t, 7, service.CalculateNumerology(dob2))

	// Example that results in 11 (master number)
	// 29/08/1990 = 2+9+0+8+1+9+9+0 = 38 => 3+8 = 11
	dob3, _ := time.Parse("2006-01-02", "1990-08-29")
	assert.Equal(t, 11, service.CalculateNumerology(dob3))
}

func TestSpiritualService_CalculateSpiritualMatchScore(t *testing.T) {
	service := NewSpiritualService()

	// Same Zodiac (Aries), Same Numerology (Odd vs Odd)
	// Aries 1: 25/03/1998 -> Aries, Num = 2+5+0+3+1+9+9+8=37->1
	dobAries1, _ := time.Parse("2006-01-02", "1998-03-25")
	
	// Aries 2: 14/04/1998 -> Aries, Num = 1+4+0+4+1+9+9+8=36->9
	dobAries2, _ := time.Parse("2006-01-02", "1998-04-14")

	score := service.CalculateSpiritualMatchScore(dobAries1, dobAries2)
	// Base (40) + Same Zodiac (25) + Same parity Numerology (20) = 85
	assert.Equal(t, 85, score)

	// Fire (Aries) vs Water (Scorpio)
	dobScorpio, _ := time.Parse("2006-01-02", "1999-11-05") // 5+1+1+1+9+9+9 = 35->8 (Even)
	score2 := service.CalculateSpiritualMatchScore(dobAries1, dobScorpio)
	// Base (40) + Fire/Water (5) + Odd/Even (10) = 55
	assert.Equal(t, 55, score2)
}
