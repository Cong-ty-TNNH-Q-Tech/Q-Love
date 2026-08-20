// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package imageutil

import (
	"image"
	"image/color"
	"testing"
)

func TestApplyGaussianBlur(t *testing.T) {
	// Create a simple 10x10 image
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	
	// Fill with white
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			img.Set(x, y, color.White)
		}
	}

	// Test 0% blur (should return same image)
	blurred0 := ApplyGaussianBlur(img, 0)
	if blurred0 != img {
		t.Errorf("Expected 0%% blur to return original image")
	}

	// Test 100% blur (should process and return a new image)
	blurred100 := ApplyGaussianBlur(img, 100)
	if blurred100 == img {
		t.Errorf("Expected 100%% blur to return a new image")
	}
	
	// Test boundaries and invalid inputs
	blurredNeg := ApplyGaussianBlur(img, -10)
	if blurredNeg != img {
		t.Errorf("Expected negative blur to return original image")
	}
	
	blurredOver := ApplyGaussianBlur(img, 150)
	if blurredOver == img {
		t.Errorf("Expected >100%% blur to return a new image")
	}
}
