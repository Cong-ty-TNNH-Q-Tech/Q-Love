// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package imageutil

import (
	"image"
	"image/color"
	"image/draw"
)

// ApplyGaussianBlur approximates a Gaussian blur using 3 passes of a box blur.
// blurPercentage is a value from 0 to 100.
func ApplyGaussianBlur(img image.Image, blurPercentage int) image.Image {
	if blurPercentage <= 0 {
		return img
	}
	if blurPercentage > 100 {
		blurPercentage = 100
	}

	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	// Max radius is arbitrarily chosen as 1/10th of the smaller dimension.
	maxRadius := w
	if h < maxRadius {
		maxRadius = h
	}
	maxRadius = maxRadius / 10

	radius := int(float64(blurPercentage) / 100.0 * float64(maxRadius))
	if radius <= 0 {
		return img
	}

	// Convert image to RGBA for faster processing
	var rgba *image.RGBA
	if r, ok := img.(*image.RGBA); ok {
		rgba = r
	} else {
		rgba = image.NewRGBA(bounds)
		draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)
	}

	// 3 passes of box blur approximate a Gaussian blur
	b1 := boxBlur(rgba, radius)
	b2 := boxBlur(b1, radius)
	b3 := boxBlur(b2, radius)

	return b3
}

// boxBlur applies a horizontal then vertical box blur.
func boxBlur(img *image.RGBA, radius int) *image.RGBA {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	out := image.NewRGBA(bounds)
	tmp := image.NewRGBA(bounds)

	// Horizontal blur
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			var r, g, b, a uint32
			var count uint32

			for k := -radius; k <= radius; k++ {
				nx := x + k
				if nx >= 0 && nx < w {
					c := img.RGBAAt(nx, y)
					r += uint32(c.R)
					g += uint32(c.G)
					b += uint32(c.B)
					a += uint32(c.A)
					count++
				}
			}

			tmp.SetRGBA(x, y, color.RGBA{
				R: uint8(r / count),
				G: uint8(g / count),
				B: uint8(b / count),
				A: uint8(a / count),
			})
		}
	}

	// Vertical blur
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			var r, g, b, a uint32
			var count uint32

			for k := -radius; k <= radius; k++ {
				ny := y + k
				if ny >= 0 && ny < h {
					c := tmp.RGBAAt(x, ny)
					r += uint32(c.R)
					g += uint32(c.G)
					b += uint32(c.B)
					a += uint32(c.A)
					count++
				}
			}

			out.SetRGBA(x, y, color.RGBA{
				R: uint8(r / count),
				G: uint8(g / count),
				B: uint8(b / count),
				A: uint8(a / count),
			})
		}
	}

	return out
}
