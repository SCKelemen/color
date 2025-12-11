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
)

const (
	width     = 800
	height    = 800
	numFrames = 60 // 60 frames for smooth rotation
)

func main() {
	outputDir := "docs/animations"
	os.MkdirAll(outputDir, 0755)

	// Generate rotating animations for each model
	models := []struct {
		name string
		fn   func(frame int, totalFrames int) *image.RGBA
	}{
		{"rgb_cube", generateRGBCubeFrame},
		{"hsl_cylinder", generateHSLCylinderFrame},
		{"lab_space", generateLABSpaceFrame},
		{"oklch_space", generateOKLCHSpaceFrame},
	}

	for _, m := range models {
		frameDir := filepath.Join(outputDir, m.name)
		os.MkdirAll(frameDir, 0755)

		for frame := 0; frame < numFrames; frame++ {
			img := m.fn(frame, numFrames)
			filename := filepath.Join(frameDir, fmt.Sprintf("frame_%03d.png", frame))
			f, err := os.Create(filename)
			if err != nil {
				panic(err)
			}
			png.Encode(f, img)
			f.Close()
		}
		fmt.Printf("Generated %d frames for %s\n", numFrames, m.name)
	}
}

func generateRGBCubeFrame(frame, totalFrames int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill with transparent background
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{0, 0, 0, 0})
		}
	}

	centerX := float64(width) / 2
	centerY := float64(height) / 2
	scale := float64(width) * 0.4

	// Calculate rotation angle
	angle := float64(frame) * 2 * math.Pi / float64(totalFrames)

	// Sample RGB colors with rotation
	g := 0.5 // Fixed green value
	step := 0.01
	for r := 0.0; r <= 1.0; r += step {
		for b := 0.0; b <= 1.0; b += step {
			c := col.RGB(r, g, b)
			rgbR, rgbG, rgbB, _ := c.RGBA()

			// Rotate coordinates
			x := (r - 0.5)
			y := (b - 0.5)
			rotX := x*math.Cos(angle) - y*math.Sin(angle)
			rotY := x*math.Sin(angle) + y*math.Cos(angle)

			px := int(centerX + rotX*scale)
			py := int(centerY - rotY*scale)

			if px >= 0 && px < width && py >= 0 && py < height {
				img.Set(px, py, color.RGBA{
					uint8(clamp255(rgbR * 255)),
					uint8(clamp255(rgbG * 255)),
					uint8(clamp255(rgbB * 255)),
					255,
				})
			}
		}
	}

	// Draw axes
	drawLine(img, int(centerX-scale), int(centerY), int(centerX+scale), int(centerY), color.RGBA{255, 255, 255, 255}, 2)
	drawLine(img, int(centerX), int(centerY-scale), int(centerX), int(centerY+scale), color.RGBA{255, 255, 255, 255}, 2)

	return img
}

func generateHSLCylinderFrame(frame, totalFrames int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill with transparent background
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{0, 0, 0, 0})
		}
	}

	centerX := float64(width) / 2
	centerY := float64(height) / 2
	radius := float64(width) * 0.35

	// Calculate rotation angle
	angle := float64(frame) * 2 * math.Pi / float64(totalFrames)

	// Fill with HSL colors showing a circular slice (top-down view of cylinder)
	l := 0.5
	for s := 0.0; s <= 1.0; s += 0.005 {
		for h := 0.0; h < 1.0; h += 0.005 {
			c := col.NewHSL(h*360, s, l, 1.0)
			r, g, b, _ := c.RGBA()
			hueAngle := h*2*math.Pi + angle // Add rotation
			// Map to circular area centered in the image
			x := centerX + radius*s*math.Cos(hueAngle)
			y := centerY + radius*s*math.Sin(hueAngle)
			if x >= 0 && x < float64(width) && y >= 0 && y < float64(height) {
				// Check if point is within the circle
				dx := x - centerX
				dy := y - centerY
				dist := math.Sqrt(dx*dx + dy*dy)
				if dist <= radius {
					img.Set(int(x), int(y), color.RGBA{
						uint8(clamp255(r * 255)),
						uint8(clamp255(g * 255)),
						uint8(clamp255(b * 255)),
						255,
					})
				}
			}
		}
	}

	return img
}

func generateLABSpaceFrame(frame, totalFrames int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill with transparent background
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{0, 0, 0, 0})
		}
	}

	centerX := float64(width) / 2
	centerY := float64(height) / 2
	scale := float64(width) * 0.3

	// Calculate rotation angle
	angle := float64(frame) * 2 * math.Pi / float64(totalFrames)

	// Draw LAB space as a slice at L=50
	l := 50.0
	step := 0.3
	for a := -80.0; a <= 80.0; a += step {
		for b := -80.0; b <= 80.0; b += step {
			// Convert LAB to RGB
			lab := col.NewLAB(l, a, b, 1.0)
			r, g, b, _ := lab.RGBA()

			// Rotate coordinates
			x := a / 100.0
			y := b / 100.0
			rotX := x*math.Cos(angle) - y*math.Sin(angle)
			rotY := x*math.Sin(angle) + y*math.Cos(angle)

			px := int(centerX + rotX*scale)
			py := int(centerY - rotY*scale)

			if px >= 0 && px < width && py >= 0 && py < height {
				img.Set(px, py, color.RGBA{
					uint8(clamp255(r * 255)),
					uint8(clamp255(g * 255)),
					uint8(clamp255(b * 255)),
					255,
				})
			}
		}
	}

	// Draw axes
	drawLine(img, int(centerX-scale), int(centerY), int(centerX+scale), int(centerY), color.RGBA{255, 255, 255, 255}, 2)
	drawLine(img, int(centerX), int(centerY-scale), int(centerX), int(centerY+scale), color.RGBA{255, 255, 255, 255}, 2)

	return img
}

func generateOKLCHSpaceFrame(frame, totalFrames int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill with transparent background
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{0, 0, 0, 0})
		}
	}

	centerX := float64(width) / 2
	centerY := float64(height) / 2
	maxRadius := float64(width) * 0.35

	// Calculate rotation angle
	angle := float64(frame) * 2 * math.Pi / float64(totalFrames)

	// Draw OKLCH space as a slice at L=0.5
	l := 0.5
	for c := 0.0; c <= 0.3; c += 0.003 {
		for h := 0.0; h < 360.0; h += 0.5 {
			oklch := col.NewOKLCH(l, c, h, 1.0)
			r, g, b, _ := oklch.RGBA()

			// Only draw if color is valid and within sRGB gamut
			if r >= 0 && r <= 1 && g >= 0 && g <= 1 && b >= 0 && b <= 1 {
				// Convert polar to Cartesian with rotation
				hueAngle := h*math.Pi/180.0 + angle
				px := int(centerX + c*maxRadius*math.Cos(hueAngle))
				py := int(centerY - c*maxRadius*math.Sin(hueAngle))

				if px >= 0 && px < width && py >= 0 && py < height {
					img.Set(px, py, color.RGBA{
						uint8(clamp255(r * 255)),
						uint8(clamp255(g * 255)),
						uint8(clamp255(b * 255)),
						255,
					})
				}
			}
		}
	}

	// Draw center point
	for dy := -3; dy <= 3; dy++ {
		for dx := -3; dx <= 3; dx++ {
			if dx*dx+dy*dy <= 9 {
				img.Set(int(centerX)+dx, int(centerY)+dy, color.RGBA{255, 255, 255, 255})
			}
		}
	}

	return img
}

func drawLine(img *image.RGBA, x0, y0, x1, y1 int, col color.Color, thickness int) {
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	sx := 1
	if x0 > x1 {
		sx = -1
	}
	sy := 1
	if y0 > y1 {
		sy = -1
	}
	err := dx - dy

	x, y := x0, y0
	for {
		for ty := -thickness / 2; ty <= thickness/2; ty++ {
			for tx := -thickness / 2; tx <= thickness/2; tx++ {
				nx, ny := x+tx, y+ty
				if nx >= 0 && nx < img.Bounds().Dx() && ny >= 0 && ny < img.Bounds().Dy() {
					img.Set(nx, ny, col)
				}
			}
		}
		if x == x1 && y == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x += sx
		}
		if e2 < dx {
			err += dx
			y += sy
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func clamp255(v float64) int {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return int(v)
}

