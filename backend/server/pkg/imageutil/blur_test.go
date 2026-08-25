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

	// Test non-RGBA image
	grayImg := image.NewGray(image.Rect(0, 0, 10, 10))
	blurredGray := ApplyGaussianBlur(grayImg, 100)
	if blurredGray == grayImg {
		t.Errorf("Expected Gray image to be processed and return a new image")
	}

	// Test fast resize code path
	imgBig := image.NewRGBA(image.Rect(0, 0, 200, 200))
	for y := 0; y < 200; y++ {
		for x := 0; x < 200; x++ {
			imgBig.Set(x, y, color.White)
		}
	}
	blurredBig := ApplyGaussianBlur(imgBig, 50)
	if blurredBig.Bounds().Dx() > 150 {
		t.Errorf("Expected image to be resized to max width 150")
	}

	// Test small radius calculation resulting in <= 0
	smallImg := image.NewRGBA(image.Rect(0, 0, 5, 5))
	blurredSmall := ApplyGaussianBlur(smallImg, 1) // maxRadius = 0 -> radius = 0
	if blurredSmall != smallImg {
		t.Errorf("Expected radius <= 0 to return original image")
	}
}
