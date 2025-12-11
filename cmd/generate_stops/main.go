package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"

	col "github.com/SCKelemen/color"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

const (
	gradientWidth  = 830
	gradientHeight = 50
	cornerRadius   = 8
	padding        = 20
	scale          = 3 // Render at 3x for better quality
)

func main() {
	stopHeight := 33 // 2/3 of bar height (50 * 2/3 = ~33px)
	outputDir := "docs/gradients"
	os.MkdirAll(outputDir, 0755)

	start := col.RGB(1, 0, 0) // red
	end := col.RGB(0, 0, 1)   // blue

	// Create image at higher resolution
	scaledWidth := gradientWidth * scale
	scaledHeight := (stopHeight + padding*2) * scale
	img := image.NewRGBA(image.Rect(0, 0, scaledWidth, scaledHeight))

	// Fill with transparent background
	for y := 0; y < scaledHeight; y++ {
		for x := 0; x < scaledWidth; x++ {
			img.Set(x, y, color.RGBA{0, 0, 0, 0})
		}
	}

	// Draw start square at 0px
	drawRoundedSquareAA(img, 0, padding*scale, stopHeight*scale, start, scale)

	// Draw end square ending at 830px (so it starts at 830 - stopHeight)
	drawRoundedSquareAA(img, (gradientWidth-stopHeight)*scale, padding*scale, stopHeight*scale, end, scale)

	// Scale down to final size
	finalImg := scaleDown(img, gradientWidth, stopHeight+padding*2)

	// Save
	filename := filepath.Join(outputDir, "stops.png")
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	if err := png.Encode(f, finalImg); err != nil {
		panic(err)
	}
	f.Close()
	fmt.Printf("Generated %s\n", filename)
}

func drawRoundedSquareAA(img *image.RGBA, x, y, size int, c col.Color, scale int) {
	r, g, b, _ := c.RGBA()
	cr := uint8(clamp255(r * 255))
	cg := uint8(clamp255(g * 255))
	cb := uint8(clamp255(b * 255))
	colr := color.RGBA{cr, cg, cb, 255}

	// Draw rounded rectangle with anti-aliasing
	for yy := 0; yy < size; yy++ {
		for xx := 0; xx < size; xx++ {
			alpha := getRoundedRectAlpha(xx, yy, size, size, cornerRadius*scale)
			if alpha > 0 {
				// Blend with existing pixel
				existing := img.RGBAAt(x+xx, y+yy)
				blended := color.RGBA{
					R: uint8((int(existing.R)*(255-int(alpha)) + int(colr.R)*int(alpha)) / 255),
					G: uint8((int(existing.G)*(255-int(alpha)) + int(colr.G)*int(alpha)) / 255),
					B: uint8((int(existing.B)*(255-int(alpha)) + int(colr.B)*int(alpha)) / 255),
					A: 255,
				}
				img.Set(x+xx, y+yy, blended)
			}
		}
	}

	// Draw color code text (both black and white for visibility)
	hex := col.RGBToHex(c)
	drawTextScaled(img, x+size/2, y+size+15*scale, hex, "black", true, scale)             // center aligned, black
	drawTextScaled(img, x+size/2+scale, y+size+15*scale+scale, hex, "white", true, scale) // center aligned, white shadow
}

func getRoundedRectAlpha(x, y, w, h, r int) uint8 {
	// Check if point is in the main rectangle (excluding corners)
	if x >= r && x < w-r {
		return 255
	}
	if y >= r && y < h-r {
		return 255
	}

	// Check corners with anti-aliasing
	var dist float64
	var centerX, centerY int

	// Top-left
	if x < r && y < r {
		centerX, centerY = r, r
		dist = math.Sqrt(float64((x-centerX)*(x-centerX) + (y-centerY)*(y-centerY)))
	} else if x >= w-r && y < r {
		// Top-right
		centerX, centerY = w-r, r
		dist = math.Sqrt(float64((x-centerX)*(x-centerX) + (y-centerY)*(y-centerY)))
	} else if x < r && y >= h-r {
		// Bottom-left
		centerX, centerY = r, h-r
		dist = math.Sqrt(float64((x-centerX)*(x-centerX) + (y-centerY)*(y-centerY)))
	} else if x >= w-r && y >= h-r {
		// Bottom-right
		centerX, centerY = w-r, h-r
		dist = math.Sqrt(float64((x-centerX)*(x-centerX) + (y-centerY)*(y-centerY)))
	} else {
		return 0
	}

	// Anti-aliasing: smooth transition at edge
	if dist <= float64(r)-0.5 {
		return 255
	} else if dist >= float64(r)+0.5 {
		return 0
	} else {
		// Smooth edge
		alpha := 0.5 - (dist - float64(r))
		return uint8(alpha * 255)
	}
}

func drawTextScaled(img *image.RGBA, x, y int, text, textColor string, center bool, scale int) {
	face := basicfont.Face7x13
	textCol := color.RGBA{0, 0, 0, 255}
	if textColor == "white" {
		textCol = color.RGBA{255, 255, 255, 255}
	}

	// Calculate text width at base resolution
	textWidth := len(text) * 7
	if center {
		x = x - (textWidth*scale)/2
	}

	// Draw text scaled up by drawing at higher resolution
	// Create a temporary image for the text at base resolution
	tempImg := image.NewRGBA(image.Rect(0, 0, textWidth+2, 13+2))
	point := fixed.Point26_6{
		X: fixed.Int26_6(1 * 64),
		Y: fixed.Int26_6(13 * 64),
	}

	d := &font.Drawer{
		Dst:  tempImg,
		Src:  image.NewUniform(textCol),
		Face: face,
		Dot:  point,
	}
	d.DrawString(text)

	// Scale up the text by drawing each pixel scale*scale times
	for ty := 0; ty < tempImg.Bounds().Dy(); ty++ {
		for tx := 0; tx < tempImg.Bounds().Dx(); tx++ {
			c := tempImg.RGBAAt(tx, ty)
			if c.A > 0 {
				// Draw this pixel scaled up
				for sy := 0; sy < scale; sy++ {
					for sx := 0; sx < scale; sx++ {
						dstX := x + tx*scale + sx
						dstY := y + ty*scale + sy
						if dstX >= 0 && dstX < img.Bounds().Dx() && dstY >= 0 && dstY < img.Bounds().Dy() {
							img.Set(dstX, dstY, c)
						}
					}
				}
			}
		}
	}
}

func scaleDown(src *image.RGBA, dstWidth, dstHeight int) *image.RGBA {
	dst := image.NewRGBA(image.Rect(0, 0, dstWidth, dstHeight))
	srcWidth := src.Bounds().Dx()
	srcHeight := src.Bounds().Dy()

	for y := 0; y < dstHeight; y++ {
		for x := 0; x < dstWidth; x++ {
			// Use area sampling for better quality
			var r, g, b, a int64
			sampleSize := scale
			for sy := 0; sy < sampleSize; sy++ {
				for sx := 0; sx < sampleSize; sx++ {
					srcX := (x*srcWidth)/dstWidth + sx
					srcY := (y*srcHeight)/dstHeight + sy
					if srcX < srcWidth && srcY < srcHeight {
						c := src.RGBAAt(srcX, srcY)
						r += int64(c.R)
						g += int64(c.G)
						b += int64(c.B)
						a += int64(c.A)
					}
				}
			}
			samples := int64(sampleSize * sampleSize)
			dst.Set(x, y, color.RGBA{
				R: uint8(r / samples),
				G: uint8(g / samples),
				B: uint8(b / samples),
				A: uint8(a / samples),
			})
		}
	}
	return dst
}

func clamp255(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return v
}
