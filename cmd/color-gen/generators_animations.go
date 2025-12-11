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
	numFrames = 60 // 60 frames for smooth rotation
)

func generateAnimations() error {
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
	return nil
}

func generateRGBCubeFrame(frame, totalFrames int) *image.RGBA {
	width, height := 800, 800
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
				// Only draw if color is valid and within sRGB gamut
				// RGBA() already returns clamped [0, 1] values in sRGB space
				if rgbR >= 0 && rgbR <= 1 && rgbG >= 0 && rgbG <= 1 && rgbB >= 0 && rgbB <= 1 {
					// Convert normalized [0, 1] to uint8 [0, 255]
					img.Set(px, py, color.RGBA{
						uint8(rgbR * 255),
						uint8(rgbG * 255),
						uint8(rgbB * 255),
						255,
					})
				}
			}
		}
	}

	// Draw axes
	drawLine(img, int(centerX-scale), int(centerY), int(centerX+scale), int(centerY), color.RGBA{255, 255, 255, 255}, 2)
	drawLine(img, int(centerX), int(centerY-scale), int(centerX), int(centerY+scale), color.RGBA{255, 255, 255, 255}, 2)

	return img
}

func generateHSLCylinderFrame(frame, totalFrames int) *image.RGBA {
	width, height := 800, 800
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
					// Only draw if color is valid and within sRGB gamut
					// RGBA() already returns clamped [0, 1] values in sRGB space
					if r >= 0 && r <= 1 && g >= 0 && g <= 1 && b >= 0 && b <= 1 {
						// Convert normalized [0, 1] to uint8 [0, 255]
						img.Set(int(x), int(y), color.RGBA{
							uint8(r * 255),
							uint8(g * 255),
							uint8(b * 255),
							255,
						})
					}
				}
			}
		}
	}

	return img
}

func generateLABSpaceFrame(frame, totalFrames int) *image.RGBA {
	width, height := 800, 800
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
			// Convert to XYZ first to check if in gamut
			xyz := col.ToXYZ(lab)
			// Convert XYZ to linear RGB
			linearR := xyz.X*3.2404542 - xyz.Y*1.5371385 - xyz.Z*0.4985314
			linearG := -xyz.X*0.9692660 + xyz.Y*1.8760108 + xyz.Z*0.0415560
			linearB := xyz.X*0.0556434 - xyz.Y*0.2040259 + xyz.Z*1.0572252
			
			// Only draw if color is within sRGB gamut (linear RGB in [0, 1])
			if linearR >= 0 && linearR <= 1 && linearG >= 0 && linearG <= 1 && linearB >= 0 && linearB <= 1 {
				// Apply gamma correction (sRGB transfer function)
				var r, g, bVal float64
				if linearR <= 0.0031308 {
					r = 12.92 * linearR
				} else {
					r = 1.055*math.Pow(linearR, 1.0/2.4) - 0.055
				}
				if linearG <= 0.0031308 {
					g = 12.92 * linearG
				} else {
					g = 1.055*math.Pow(linearG, 1.0/2.4) - 0.055
				}
				if linearB <= 0.0031308 {
					bVal = 12.92 * linearB
				} else {
					bVal = 1.055*math.Pow(linearB, 1.0/2.4) - 0.055
				}

				// Rotate coordinates
				x := a / 100.0
				y := b / 100.0
				rotX := x*math.Cos(angle) - y*math.Sin(angle)
				rotY := x*math.Sin(angle) + y*math.Cos(angle)

				px := int(centerX + rotX*scale)
				py := int(centerY - rotY*scale)

				if px >= 0 && px < width && py >= 0 && py < height {
					// Convert normalized [0, 1] to uint8 [0, 255]
					img.Set(px, py, color.RGBA{
						uint8(r * 255),
						uint8(g * 255),
						uint8(b * 255),
						255,
					})
				}
			}
		}
	}

	// Draw axes
	drawLine(img, int(centerX-scale), int(centerY), int(centerX+scale), int(centerY), color.RGBA{255, 255, 255, 255}, 2)
	drawLine(img, int(centerX), int(centerY-scale), int(centerX), int(centerY+scale), color.RGBA{255, 255, 255, 255}, 2)

	return img
}

func generateOKLCHSpaceFrame(frame, totalFrames int) *image.RGBA {
	width, height := 800, 800
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
			// Convert OKLCH to OKLAB (polar to rectangular)
			rad := h * math.Pi / 180
			oka := c * math.Cos(rad)
			okb := c * math.Sin(rad)
			oklab := col.NewOKLAB(l, oka, okb, 1.0)
			// Convert OKLAB to linear RGB
			l_ := oklab.L + 0.3963377774*oklab.A + 0.2158037573*oklab.B
			m := oklab.L - 0.1055613458*oklab.A - 0.0638541728*oklab.B
			s := oklab.L - 0.0894841775*oklab.A - 1.2914855480*oklab.B

			l3 := l_ * l_ * l_
			m3 := m * m * m
			s3 := s * s * s

			linearR := +4.0767416621*l3 - 3.3077115913*m3 + 0.2309699292*s3
			linearG := -1.2684380046*l3 + 2.6097574011*m3 - 0.3413193965*s3
			linearB := -0.0041960863*l3 - 0.7034186147*m3 + 1.7076147010*s3

			// Only draw if color is within sRGB gamut (linear RGB in [0, 1])
			if linearR >= 0 && linearR <= 1 && linearG >= 0 && linearG <= 1 && linearB >= 0 && linearB <= 1 {
				// Apply gamma correction (sRGB transfer function)
				var r, g, b float64
				if linearR <= 0.0031308 {
					r = 12.92 * linearR
				} else {
					r = 1.055*math.Pow(linearR, 1.0/2.4) - 0.055
				}
				if linearG <= 0.0031308 {
					g = 12.92 * linearG
				} else {
					g = 1.055*math.Pow(linearG, 1.0/2.4) - 0.055
				}
				if linearB <= 0.0031308 {
					b = 12.92 * linearB
				} else {
					b = 1.055*math.Pow(linearB, 1.0/2.4) - 0.055
				}

				// Convert polar to Cartesian with rotation
				hueAngle := h*math.Pi/180.0 + angle
				px := int(centerX + c*maxRadius*math.Cos(hueAngle))
				py := int(centerY - c*maxRadius*math.Sin(hueAngle))

				if px >= 0 && px < width && py >= 0 && py < height {
					// Convert normalized [0, 1] to uint8 [0, 255]
					img.Set(px, py, color.RGBA{
						uint8(r * 255),
						uint8(g * 255),
						uint8(b * 255),
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

// clamp01 clamps a value to the range [0, 1]
func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}
