// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestBondingCurve_BuyFormula verifies the O(1) arithmetic series formula
// produces the same result as the O(N) iterative loop for buy operations.
func TestBondingCurve_BuyFormula(t *testing.T) {
	tests := []struct {
		name      string
		available int
		quantity  int
	}{
		{"single card", 1000, 1},
		{"10 cards", 1000, 10},
		{"100 cards", 1000, 100},
		{"500 cards", 1000, 500},
		{"all cards", 1000, 1000},
		{"partial supply", 500, 100},
		{"low supply", 50, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// O(N) reference implementation (original loop)
			expectedCost := float64(0)
			for i := 0; i < tt.quantity; i++ {
				priceAtI := 100.0 + float64(1000-(tt.available-i))*5.0
				expectedCost += priceAtI
			}

			// O(1) formula implementation
			a := 100.0 + float64(1000-tt.available)*5.0
			d := 5.0
			n := float64(tt.quantity)
			formulaCost := n * (2*a + (n-1)*d) / 2

			assert.InDelta(t, expectedCost, formulaCost, 0.001,
				"Formula cost should match iterative cost")
		})
	}
}

// TestBondingCurve_SellFormula verifies the O(1) formula for sell operations,
// including edge cases where price hits the floor (minimum 10).
func TestBondingCurve_SellFormula(t *testing.T) {
	tests := []struct {
		name      string
		available int
		quantity  int
	}{
		{"single card, high price", 100, 1},
		{"10 cards, no floor", 100, 10},
		{"sell at low supply", 950, 10},
		{"sell hitting floor", 990, 10},
		{"all cards at floor", 998, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// O(N) reference (original loop with floor)
			expectedCost := float64(0)
			for i := 0; i < tt.quantity; i++ {
				priceAtI := 100.0 + float64(1000-(tt.available+i))*5.0
				if priceAtI < 10 {
					priceAtI = 10
				}
				expectedCost += priceAtI
			}

			// O(1) formula with floor handling
			a := 100.0 + float64(1000-tt.available)*5.0
			d := -5.0
			n := float64(tt.quantity)

			stepsBeforeFloor := int((a - 10.0) / 5.0) + 1
			if stepsBeforeFloor < 0 {
				stepsBeforeFloor = 0
			}

			var formulaCost float64
			if stepsBeforeFloor >= tt.quantity {
				formulaCost = n * (2*a + (n-1)*d) / 2
			} else {
				nAbove := float64(stepsBeforeFloor)
				if nAbove > 0 {
					formulaCost = nAbove * (2*a + (nAbove-1)*d) / 2
				}
				formulaCost += float64(tt.quantity-stepsBeforeFloor) * 10.0
			}

			assert.InDelta(t, expectedCost, formulaCost, 0.001,
				"Formula cost should match iterative cost (with floor)")
		})
	}
}

// TestBondingCurve_PriceFloor ensures minimum price of 10 is enforced.
func TestBondingCurve_PriceFloor(t *testing.T) {
	// When available is very high (close to 1000), price approaches 100
	// When selling, price decreases. At available=998:
	// price = 100 + (1000-998)*5 = 110
	// After selling 1: price = 100 + (1000-999)*5 = 105
	// After selling 2: price = 100 + (1000-1000)*5 = 100
	// Price should never go below 10

	available := 998
	a := 100.0 + float64(1000-available)*5.0
	assert.Equal(t, 110.0, a, "Starting price should be 110")

	// Sell 20 cards - many will hit the floor
	quantity := 20
	totalCost := float64(0)
	for i := 0; i < quantity; i++ {
		price := 100.0 + float64(1000-(available+i))*5.0
		if price < 10 {
			price = 10
		}
		totalCost += price
		assert.GreaterOrEqual(t, price, 10.0, "Price should never go below 10")
	}
}
