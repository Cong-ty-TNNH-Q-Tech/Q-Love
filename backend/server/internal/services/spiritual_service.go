// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"fmt"
	"strconv"
	"time"
)

type SpiritualService interface {
	CalculateZodiac(dob time.Time) string
	CalculateNumerology(dob time.Time) int
	CalculateSpiritualMatchScore(dobA, dobB time.Time) int
}

type spiritualService struct{}

func NewSpiritualService() SpiritualService {
	return &spiritualService{}
}

// Cung Hoàng Đạo
func (s *spiritualService) CalculateZodiac(dob time.Time) string {
	day := dob.Day()
	month := dob.Month()

	switch month {
	case time.March:
		if day >= 21 {
			return "Aries"
		}
		return "Pisces"
	case time.April:
		if day >= 20 {
			return "Taurus"
		}
		return "Aries"
	case time.May:
		if day >= 21 {
			return "Gemini"
		}
		return "Taurus"
	case time.June:
		if day >= 22 {
			return "Cancer"
		}
		return "Gemini"
	case time.July:
		if day >= 23 {
			return "Leo"
		}
		return "Cancer"
	case time.August:
		if day >= 23 {
			return "Virgo"
		}
		return "Leo"
	case time.September:
		if day >= 23 {
			return "Libra"
		}
		return "Virgo"
	case time.October:
		if day >= 24 {
			return "Scorpio"
		}
		return "Libra"
	case time.November:
		if day >= 22 {
			return "Sagittarius"
		}
		return "Scorpio"
	case time.December:
		if day >= 22 {
			return "Capricorn"
		}
		return "Sagittarius"
	case time.January:
		if day >= 20 {
			return "Aquarius"
		}
		return "Capricorn"
	case time.February:
		if day >= 19 {
			return "Pisces"
		}
		return "Aquarius"
	}
	return "Unknown"
}

// Thần Số Học (Numerology)
func (s *spiritualService) CalculateNumerology(dob time.Time) int {
	dateStr := dob.Format("02012006")
	sum := 0
	for _, char := range dateStr {
		val, _ := strconv.Atoi(string(char))
		sum += val
	}

	for sum > 9 && sum != 11 && sum != 22 {
		temp := 0
		sumStr := fmt.Sprintf("%d", sum)
		for _, char := range sumStr {
			val, _ := strconv.Atoi(string(char))
			temp += val
		}
		sum = temp
	}

	return sum
}

func getZodiacElement(zodiac string) string {
	fire := []string{"Aries", "Leo", "Sagittarius"}
	earth := []string{"Taurus", "Virgo", "Capricorn"}
	air := []string{"Gemini", "Libra", "Aquarius"}
	water := []string{"Cancer", "Scorpio", "Pisces"}

	for _, z := range fire {
		if z == zodiac {
			return "Fire"
		}
	}
	for _, z := range earth {
		if z == zodiac {
			return "Earth"
		}
	}
	for _, z := range air {
		if z == zodiac {
			return "Air"
		}
	}
	for _, z := range water {
		if z == zodiac {
			return "Water"
		}
	}
	return "Unknown"
}

// Điểm tương thích tâm linh (0 - 100)
func (s *spiritualService) CalculateSpiritualMatchScore(dobA, dobB time.Time) int {
	score := 40 // Base score

	// 1. Zodiac Compatibility (Max 30)
	zA := s.CalculateZodiac(dobA)
	zB := s.CalculateZodiac(dobB)

	if zA == zB {
		score += 25
	} else {
		eleA := getZodiacElement(zA)
		eleB := getZodiacElement(zB)
		if eleA == eleB {
			score += 30 // Same element = great match
		} else if (eleA == "Fire" && eleB == "Air") || (eleA == "Air" && eleB == "Fire") {
			score += 20
		} else if (eleA == "Earth" && eleB == "Water") || (eleA == "Water" && eleB == "Earth") {
			score += 20
		} else {
			score += 5
		}
	}

	// 2. Numerology Compatibility (Max 30)
	nA := s.CalculateNumerology(dobA)
	nB := s.CalculateNumerology(dobB)

	// Odd and even compatibility
	if nA == nB {
		score += 30
	} else if (nA%2 == 0 && nB%2 == 0) || (nA%2 != 0 && nB%2 != 0) {
		score += 20 // Both even or both odd
	} else {
		score += 10
	}

	if score >= 100 {
		score = 100
	}

	return score
}
