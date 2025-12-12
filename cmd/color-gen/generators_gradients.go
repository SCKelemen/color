package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"strings"

	col "github.com/SCKelemen/color"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

const (
	labelHeight = 30
)

var defaultFace font.Face

func generateGradients(start, end col.Color) error {
	outputDir := "docs/gradients"
	os.MkdirAll(outputDir, 0755)

	spaces := []struct {
		name string
		s    col.GradientSpace
	}{
		{"rgb", col.GradientRGB},
		{"hsl", col.GradientHSL},
		{"lab", col.GradientLAB},
		{"oklab", col.GradientOKLAB},
		{"lch", col.GradientLCH},
		{"oklch", col.GradientOKLCH},
	}

	for _, sp := range spaces {
		// Generate gradient
		gradient := col.GradientInSpace(start, end, gradientWidth, sp.s)

		// Generate images with transparent background and both text colors
		for _, textColor := range []string{"black", "white"} {
			// Create image at higher resolution
			// Reduced height: text is closer, so we need less space
			scaledWidth := gradientWidth * scale
			scaledHeight := (gradientHeight + padding + 20) * scale // Reduced from padding + labelHeight to padding + 20
			img := image.NewRGBA(image.Rect(0, 0, scaledWidth, scaledHeight))

			// Fill with transparent background
			for y := 0; y < scaledHeight; y++ {
				for x := 0; x < scaledWidth; x++ {
					img.Set(x, y, color.RGBA{0, 0, 0, 0})
				}
			}

			// Draw gradient (with anti-aliased rounded corners)
			drawRoundedGradientAA(img, 0, 0, scaledWidth, gradientHeight*scale, gradient, scale)

			// Draw labels with both black and white text
			// Center text vertically between bottom of gradient bar and bottom of image
			gradientBottom := gradientHeight * scale
			imageBottom := scaledHeight
			spaceBelowBar := imageBottom - gradientBottom
			textCenterY := gradientBottom + spaceBelowBar/2
			drawLabels(img, textCenterY, start, end, sp.name, textColor, scale)

			// Scale down to final size with better filtering
			finalImg := scaleDownGradients(img, gradientWidth, gradientHeight+padding+20)

			// Save
			filename := filepath.Join(outputDir, fmt.Sprintf("gradient_%s_%s.png", sp.name, textColor))
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
	}
	return nil
}

func drawRoundedGradientAA(img *image.RGBA, x, y, w, h int, gradient []col.Color, scale int) {
	for xx := 0; xx < w; xx++ {
		gradientX := xx / scale
		if gradientX >= len(gradient) {
			gradientX = len(gradient) - 1
		}
		c := gradient[gradientX]
		r, g, b, _ := c.RGBA()
		cr := uint8(clamp255(r * 255))
		cg := uint8(clamp255(g * 255))
		cb := uint8(clamp255(b * 255))
		colr := color.RGBA{cr, cg, cb, 255}

		for yy := 0; yy < h; yy++ {
			// Anti-aliased rounded rectangle check
			alpha := getRoundedRectAlpha(xx, yy, w, h, cornerRadius*scale)
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

func drawLabels(img *image.RGBA, centerY int, start, end col.Color, spaceName, textColor string, scale int) {
	// CenterY is the vertical center point where text should be centered
	// Text height is 16 pixels (scaled up)
	// Adjust Y position to account for font baseline - TrueType fonts position from baseline
	textHeight := 16 * scale
	baselineOffset := 12 * scale                     // Approximate baseline offset for 16pt font
	textY := centerY - textHeight/2 + baselineOffset // Move text down to center it properly

	// Determine primary and shadow colors based on textColor parameter
	// textColor="black" means all black text (no shadow)
	// textColor="white" means white text with black shadow for visibility
	primaryColor := textColor
	var shadowColor string
	if textColor == "white" {
		shadowColor = "black"
	} else {
		shadowColor = "" // No shadow for black text
	}

	// Left: start color in native format
	startStr := formatColorInSpace(start, spaceName)
	if shadowColor != "" {
		drawTextScaledGradients(img, padding*scale+scale, textY+scale, startStr, shadowColor, false, scale) // shadow first
	}
	drawTextScaledGradients(img, padding*scale, textY, startStr, primaryColor, false, scale) // primary text on top

	// Center: space name (uppercase)
	if shadowColor != "" {
		drawTextScaledGradients(img, (gradientWidth*scale)/2+scale, textY+scale, strings.ToUpper(spaceName), shadowColor, true, scale) // shadow first
	}
	drawTextScaledGradients(img, (gradientWidth*scale)/2, textY, strings.ToUpper(spaceName), primaryColor, true, scale) // primary text on top

	// Right: end color in native format - align to end of gradient bar
	endStr := formatColorInSpace(end, spaceName)
	// Measure text width properly for TrueType font
	var textWidth fixed.Int26_6
	if defaultTT != nil {
		scaledFace, _ := opentype.NewFace(defaultTT, &opentype.FaceOptions{
			Size:    fontSize * float64(scale),
			DPI:     72,
			Hinting: font.HintingFull,
		})
		if scaledFace != nil {
			textWidth = font.MeasureString(scaledFace, endStr)
		}
	}
	if textWidth == 0 {
		// Fallback calculation
		textWidth = fixed.Int26_6(len(endStr) * 7 * scale * 64)
	}
	rightX := (gradientWidth-padding)*scale - textWidth.Ceil() // Align right edge to padding from end
	if shadowColor != "" {
		drawTextScaledGradients(img, rightX+scale, textY+scale, endStr, shadowColor, false, scale) // shadow first
	}
	drawTextScaledGradients(img, rightX, textY, endStr, primaryColor, false, scale) // primary text on top
}

func formatColorInSpace(c col.Color, spaceName string) string {
	switch spaceName {
	case "rgb":
		r, g, b, _ := c.RGBA()
		return fmt.Sprintf("rgb(%.0f, %.0f, %.0f)", r*255, g*255, b*255)
	case "hsl":
		hsl := col.ToHSL(c)
		return fmt.Sprintf("hsl(%.0f, %.0f%%, %.0f%%)", hsl.H, hsl.S*100, hsl.L*100)
	case "lab":
		lab := col.ToLAB(c)
		return fmt.Sprintf("lab(%.0f, %.1f, %.1f)", lab.L, lab.A, lab.B)
	case "oklab":
		oklab := col.ToOKLAB(c)
		return fmt.Sprintf("oklab(%.2f, %.2f, %.2f)", oklab.L, oklab.A, oklab.B)
	case "lch":
		lch := col.ToLCH(c)
		return fmt.Sprintf("lch(%.0f, %.1f, %.0f)", lch.L, lch.C, lch.H)
	case "oklch":
		oklch := col.ToOKLCH(c)
		return fmt.Sprintf("oklch(%.2f, %.2f, %.0f)", oklch.L, oklch.C, oklch.H)
	default:
		return col.RGBToHex(c)
	}
}

func drawTextScaledGradients(img *image.RGBA, x, y int, text, textColor string, center bool, scale int) {
	textCol := color.RGBA{0, 0, 0, 255}
	if textColor == "white" {
		textCol = color.RGBA{255, 255, 255, 255}
	}

	if defaultTT != nil {
		// Create a scaled font face for high-res rendering
		// The font is 16pt, but we're rendering at 3x, so we need 16*3 = 48pt
		scaledFace, err := opentype.NewFace(defaultTT, &opentype.FaceOptions{
			Size:    fontSize * float64(scale), // Scale font size for high-res rendering
			DPI:     72,
			Hinting: font.HintingFull,
		})
		if err == nil {
			// Measure text width for centering
			textWidth := font.MeasureString(scaledFace, text)
			if center {
				x = x - textWidth.Ceil()/2
			}

			// Draw text directly at high resolution
			point := fixed.Point26_6{
				X: fixed.Int26_6(x * 64),
				Y: fixed.Int26_6(y * 64),
			}

			d := &font.Drawer{
				Dst:  img,
				Src:  image.NewUniform(textCol),
				Face: scaledFace,
				Dot:  point,
			}
			d.DrawString(text)
			return
		}
	}

	if defaultFace != nil {
		// Fallback: use defaultFace but measure correctly
		textWidth := font.MeasureString(defaultFace, text)
		if center {
			x = x - textWidth.Ceil()*scale/2
		}

		// Draw text directly at high resolution
		point := fixed.Point26_6{
			X: fixed.Int26_6(x * 64),
			Y: fixed.Int26_6(y * 64),
		}

		d := &font.Drawer{
			Dst:  img,
			Src:  image.NewUniform(textCol),
			Face: defaultFace,
			Dot:  point,
		}
		d.DrawString(text)
		return
	} else {
		// Fallback to basicfont (scaled)
		face := basicfont.Face7x13
		fontScale := float64(fontSize) / 13.0
		textWidth := int(float64(len(text)*7) * fontScale)

		if center {
			x = x - (textWidth*scale)/2
		}

		tempImg := image.NewRGBA(image.Rect(0, 0, int(float64(len(text)*7)*fontScale)+2, fontSize+2))
		point := fixed.Point26_6{
			X: fixed.Int26_6(1 * 64),
			Y: fixed.Int26_6(fontSize * 64),
		}

		d := &font.Drawer{
			Dst:  tempImg,
			Src:  image.NewUniform(textCol),
			Face: face,
			Dot:  point,
		}
		d.DrawString(text)

		// Scale up the text by the render scale factor
		for ty := 0; ty < tempImg.Bounds().Dy(); ty++ {
			for tx := 0; tx < tempImg.Bounds().Dx(); tx++ {
				c := tempImg.RGBAAt(tx, ty)
				if c.A > 0 {
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
}

func scaleDownGradients(src *image.RGBA, dstWidth, dstHeight int) *image.RGBA {
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

