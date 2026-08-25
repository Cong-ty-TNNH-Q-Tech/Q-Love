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
func TestSpiritualService_CalculateZodiac_AllMonths(t *testing.T) {
	service := NewSpiritualService()

	tests := []struct {
		date     string
		expected string
	}{
		{"2000-01-10", "Capricorn"},
		{"2000-01-25", "Aquarius"},
		{"2000-02-10", "Aquarius"},
		{"2000-02-20", "Pisces"},
		{"2000-03-10", "Pisces"},
		{"2000-03-22", "Aries"},
		{"2000-04-10", "Aries"},
		{"2000-04-22", "Taurus"},
		{"2000-05-10", "Taurus"},
		{"2000-05-22", "Gemini"},
		{"2000-06-10", "Gemini"},
		{"2000-06-25", "Cancer"},
		{"2000-07-10", "Cancer"},
		{"2000-07-25", "Leo"},
		{"2000-08-10", "Leo"},
		{"2000-08-25", "Virgo"},
		{"2000-09-10", "Virgo"},
		{"2000-09-25", "Libra"},
		{"2000-10-10", "Libra"},
		{"2000-10-25", "Scorpio"},
		{"2000-11-10", "Scorpio"},
		{"2000-11-25", "Sagittarius"},
		{"2000-12-10", "Sagittarius"},
		{"2000-12-25", "Capricorn"},
	}

	for _, tc := range tests {
		dob, _ := time.Parse("2006-01-02", tc.date)
		assert.Equal(t, tc.expected, service.CalculateZodiac(dob))
	}
}

func TestSpiritualService_CalculateNumerology_22(t *testing.T) {
	service := NewSpiritualService()
	// Let's use something that sums to 22. 29/09/2000 -> 2+9+0+9+2+0+0+0 = 22
	dob22, _ := time.Parse("2006-01-02", "2000-09-29")
	assert.Equal(t, 22, service.CalculateNumerology(dob22))
}

func TestSpiritualService_CalculateSpiritualMatchScore_Elements(t *testing.T) {
	service := NewSpiritualService()
	
	// Earth vs Water (Taurus vs Cancer)
	dobTaurus, _ := time.Parse("2006-01-02", "2000-05-10") // 10/05/2000 = 1+0+0+5+2+0+0+0 = 8
	dobCancer, _ := time.Parse("2006-01-02", "2000-07-10") // 10/07/2000 = 1+0+0+7+2+0+0+0 = 10 -> 1
	
	// score: Base(40) + Earth/Water(20) + Even/Odd(10) = 70
	score1 := service.CalculateSpiritualMatchScore(dobTaurus, dobCancer)
	assert.Equal(t, 70, score1)
	
	// Fire vs Air (Aries vs Gemini)
	dobAries, _ := time.Parse("2006-01-02", "2000-04-10") // 10/04/2000 = 1+0+0+4+2+0+0+0 = 7
	dobGemini, _ := time.Parse("2006-01-02", "2000-06-10") // 10/06/2000 = 1+0+0+6+2+0+0+0 = 9
	// score: Base(40) + Fire/Air(20) + Odd/Odd(20) = 80
	score2 := service.CalculateSpiritualMatchScore(dobAries, dobGemini)
	assert.Equal(t, 80, score2)
	
	// Same Element (Earth vs Earth: Taurus vs Virgo)
	dobVirgo, _ := time.Parse("2006-01-02", "2000-09-10") // 10/09/2000 = 1+0+0+9+2+0+0+0 = 12 -> 3
	// Taurus(8) vs Virgo(3)
	// score: Base(40) + Same Element(30) + Even/Odd(10) = 80
	score3 := service.CalculateSpiritualMatchScore(dobTaurus, dobVirgo)
	assert.Equal(t, 80, score3)
	
	// Air vs Fire (Gemini vs Aries)
	score4 := service.CalculateSpiritualMatchScore(dobGemini, dobAries)
	assert.Equal(t, 80, score4)

	// Water vs Earth (Cancer vs Taurus)
	score5 := service.CalculateSpiritualMatchScore(dobCancer, dobTaurus)
	assert.Equal(t, 70, score5)

	// Even vs Even
	// Taurus(8) vs Capricorn(4)
	dobCap, _ := time.Parse("2006-01-02", "2000-01-01") // 1+1+2 = 4
	score6 := service.CalculateSpiritualMatchScore(dobTaurus, dobCap)
	// Base(40) + Same Element(30) + Even/Even(20) = 90
	assert.Equal(t, 90, score6)

	// 100 max capping branch
	// Taurus (8) and Virgo (8) -> Same Element (+30) and Same Numerology (+30) -> 100
	dobVirgo8, _ := time.Parse("2006-01-02", "2000-09-05") // 5+9+2 = 16 -> 7 wait.
	// We want Virgo (August 23 - Sept 22) and Numerology 8.
	// Sept 14, 2000 -> 1+4+0+9+2+0+0+0 = 16 -> 7.
	// Sept 15, 2000 -> 1+5+0+9+2+0+0+0 = 17 -> 8!
	dobVirgo8_real, _ := time.Parse("2006-01-02", "2000-09-15")
	score7 := service.CalculateSpiritualMatchScore(dobTaurus, dobVirgo8_real)
	assert.Equal(t, 100, score7)
}
