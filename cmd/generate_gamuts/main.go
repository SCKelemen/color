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
	scale    = 3  // Render at 3x for better quality
	fontSize = 16 // Font size in points
)

var defaultTT *opentype.Font

func init() {
	// Try to load Roboto font
	robotoFonts := []string{
		"fonts/Roboto/static/Roboto-Regular.ttf",
		"fonts/Roboto/static/Roboto-Medium.ttf",
		"fonts/Roboto/Roboto-VariableFont_wdth,wght.ttf",
		"fonts/Roboto/Roboto-Regular.ttf",
	}

	for _, fontPath := range robotoFonts {
		if _, err := os.Stat(fontPath); err == nil {
			fontData, err := os.ReadFile(fontPath)
			if err == nil {
				tt, err := opentype.Parse(fontData)
				if err == nil {
					defaultTT = tt
					return
				}
			}
		}
	}
}

func main() {
	outputDir := "docs/gamuts"
	os.MkdirAll(outputDir, 0755)

	spaces := []struct {
		name string
		fn   func(r, g, b float64) col.Color
	}{
		{"sRGB", func(r, g, b float64) col.Color { return col.RGB(r, g, b) }},
		{"DisplayP3", func(r, g, b float64) col.Color {
			c, _ := col.ConvertFromRGBSpace(r, g, b, 1.0, "display-p3")
			return c
		}},
		{"AdobeRGB", func(r, g, b float64) col.Color {
			c, _ := col.ConvertFromRGBSpace(r, g, b, 1.0, "a98-rgb")
			return c
		}},
		{"Rec2020", func(r, g, b float64) col.Color {
			c, _ := col.ConvertFromRGBSpace(r, g, b, 1.0, "rec2020")
			return c
		}},
	}

	for _, s := range spaces {
		// Generate both black and white text versions
		for _, textColor := range []string{"black", "white"} {
			img := generateGamutVolume(s.fn, s.name, textColor, 1000, 800)
			suffix := ""
			if textColor == "black" {
				suffix = "_black"
			} else {
				suffix = "_white"
			}
			filename := filepath.Join(outputDir, fmt.Sprintf("gamut_%s%s.png", s.name, suffix))
			f, err := os.Create(filename)
			if err != nil {
				panic(err)
			}
			png.Encode(f, img)
			f.Close()
			fmt.Printf("Generated %s\n", filename)
		}
	}
}

func generateGamutVolume(createColor func(r, g, b float64) col.Color, spaceName, textColor string, width, height int) *image.RGBA {
	// First, generate at high resolution
	scaledWidth := width * scale
	scaledHeight := height * scale
	img := image.NewRGBA(image.Rect(0, 0, scaledWidth, scaledHeight))

	// Fill with transparent background
	for y := 0; y < scaledHeight; y++ {
		for x := 0; x < scaledWidth; x++ {
			img.Set(x, y, color.RGBA{0, 0, 0, 0})
		}
	}

	// 3D rotation angles to show all axes clearly
	// Rotate around Y axis (45°) and X axis (30°) to get a good view
	angleY := 45.0 * math.Pi / 180.0 // Rotation around Y axis
	angleX := 30.0 * math.Pi / 180.0 // Rotation around X axis
	angleZ := 0.0                    // No Z rotation

	scaleFactor := float64(scaledWidth) * 0.35
	centerX := float64(scaledWidth) / 2
	centerY := float64(scaledHeight) / 2

	// Sample RGB cube and project to 3D view
	step := 0.03
	depthBuffer := make(map[int]float64) // For depth sorting

	type point struct {
		x, y int
		z    float64 // depth for sorting
		col  color.RGBA
	}
	points := make([]point, 0)

	for r := 0.0; r <= 1.0; r += step {
		for g := 0.0; g <= 1.0; g += step {
			for b := 0.0; b <= 1.0; b += step {
				// Create color in the target space
				c := createColor(r, g, b)

				// 3D rotation: rotate around Y, X, and Z axes
				// First rotate around Y axis
				x1 := r*math.Cos(angleY) - b*math.Sin(angleY)
				y1 := g
				z1 := r*math.Sin(angleY) + b*math.Cos(angleY)

				// Then rotate around X axis
				y2 := y1*math.Cos(angleX) - z1*math.Sin(angleX)
				z2 := y1*math.Sin(angleX) + z1*math.Cos(angleX)

				// Project to 2D (orthographic projection)
				xProj := x1 * scaleFactor
				yProj := y2 * scaleFactor

				// Map to image coordinates
				xCoord := int(centerX + xProj)
				yCoord := int(centerY - yProj)

				if xCoord >= 0 && xCoord < scaledWidth && yCoord >= 0 && yCoord < scaledHeight {
					// Get RGB for display (convert to sRGB)
					rgbR, rgbG, rgbB, _ := c.RGBA()
					key := yCoord*scaledWidth + xCoord

					// Depth sorting: only draw if this point is closer (larger z)
					if z2 > depthBuffer[key] {
						depthBuffer[key] = z2
						points = append(points, point{
							x: xCoord,
							y: yCoord,
							z: z2,
							col: color.RGBA{
								R: uint8(math.Max(0, math.Min(255, rgbR*255))),
								G: uint8(math.Max(0, math.Min(255, rgbG*255))),
								B: uint8(math.Max(0, math.Min(255, rgbB*255))),
								A: 255,
							},
						})
					}
				}
			}
		}
	}

	// Draw points (front to back for proper depth)
	for _, p := range points {
		img.Set(p.x, p.y, p.col)
		// Add slight glow for better visibility
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				if dx == 0 && dy == 0 {
					continue
				}
				nx, ny := p.x+dx, p.y+dy
				if nx >= 0 && nx < scaledWidth && ny >= 0 && ny < scaledHeight {
					existing := img.RGBAAt(nx, ny)
					if existing.A == 0 {
						// Only add glow to transparent areas
						blended := color.RGBA{
							R: uint8(int(p.col.R) / 3),
							G: uint8(int(p.col.G) / 3),
							B: uint8(int(p.col.B) / 3),
							A: 128,
						}
						img.Set(nx, ny, blended)
					}
				}
			}
		}
	}

	// Draw axes clearly showing R, G, B directions
	drawAxis3D(img, centerX, centerY, scaleFactor, angleY, angleX, angleZ, 1.0, 0.0, 0.0, color.RGBA{255, 50, 50, 255}, "R")
	drawAxis3D(img, centerX, centerY, scaleFactor, angleY, angleX, angleZ, 0.0, 1.0, 0.0, color.RGBA{50, 255, 50, 255}, "G")
	drawAxis3D(img, centerX, centerY, scaleFactor, angleY, angleX, angleZ, 0.0, 0.0, 1.0, color.RGBA{50, 50, 255, 255}, "B")

	// Find bounding box of non-transparent pixels
	minX, minY, maxX, maxY := findBoundingBox(img)
	
	// Add padding for label
	labelHeight := int(30 * float64(scale))
	padding := int(10 * float64(scale))
	
	// Create cropped image with label space
	croppedWidth := maxX - minX + padding*2
	croppedHeight := maxY - minY + padding*2 + labelHeight
	
	// Adjust bounds to include label
	if minY > labelHeight {
		minY -= labelHeight
	} else {
		croppedHeight += labelHeight
		minY = 0
	}
	
	croppedImg := image.NewRGBA(image.Rect(0, 0, croppedWidth, croppedHeight))
	
	// Fill with transparent background
	for y := 0; y < croppedHeight; y++ {
		for x := 0; x < croppedWidth; x++ {
			croppedImg.Set(x, y, color.RGBA{0, 0, 0, 0})
		}
	}
	
	// Copy gamut image to cropped image
	for y := minY; y < maxY; y++ {
		for x := minX; x < maxX; x++ {
			if x >= 0 && x < scaledWidth && y >= 0 && y < scaledHeight {
				c := img.RGBAAt(x, y)
				if c.A > 0 {
					croppedImg.Set(x-minX+padding, y-minY+padding+labelHeight, c)
				}
			}
		}
	}
	
	// Draw label
	drawGamutLabel(croppedImg, croppedWidth, croppedHeight, spaceName, textColor)
	
	// Scale down to final size
	return scaleDown(croppedImg, croppedWidth/scale, croppedHeight/scale)
}

func findBoundingBox(img *image.RGBA) (minX, minY, maxX, maxY int) {
	minX, minY = img.Bounds().Dx(), img.Bounds().Dy()
	maxX, maxY = 0, 0
	
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			if img.RGBAAt(x, y).A > 0 {
				if x < minX {
					minX = x
				}
				if x > maxX {
					maxX = x
				}
				if y < minY {
					minY = y
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}
	
	// Add small margin
	if minX > 0 {
		minX = max(0, minX-10*scale)
	}
	if minY > 0 {
		minY = max(0, minY-10*scale)
	}
	if maxX < img.Bounds().Dx() {
		maxX = min(img.Bounds().Dx(), maxX+10*scale)
	}
	if maxY < img.Bounds().Dy() {
		maxY = min(img.Bounds().Dy(), maxY+10*scale)
	}
	
	return
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func scaleDown(img *image.RGBA, targetWidth, targetHeight int) *image.RGBA {
	result := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))
	
	scaleX := float64(img.Bounds().Dx()) / float64(targetWidth)
	scaleY := float64(img.Bounds().Dy()) / float64(targetHeight)
	
	for y := 0; y < targetHeight; y++ {
		for x := 0; x < targetWidth; x++ {
			srcX := int(float64(x) * scaleX)
			srcY := int(float64(y) * scaleY)
			if srcX < img.Bounds().Dx() && srcY < img.Bounds().Dy() {
				result.Set(x, y, img.RGBAAt(srcX, srcY))
			}
		}
	}
	
	return result
}

func drawGamutLabel(img *image.RGBA, width, height int, spaceName, textColor string) {
	// Center the label text
	text := strings.ToUpper(spaceName)
	
	// Determine primary and shadow colors
	primaryColor := textColor
	var shadowColor string
	if textColor == "white" {
		shadowColor = "black"
	} else {
		shadowColor = ""
	}
	
	// Position text at the bottom, centered
	textY := height - int(10*float64(scale))
	drawTextScaled(img, width/2, textY, text, primaryColor, shadowColor, true, scale)
}

func drawTextScaled(img *image.RGBA, x, y int, text, primaryColor, shadowColor string, center bool, scale int) {
	textCol := color.RGBA{0, 0, 0, 255}
	if primaryColor == "white" {
		textCol = color.RGBA{255, 255, 255, 255}
	}
	
	shadowCol := color.RGBA{0, 0, 0, 255}
	if shadowColor == "white" {
		shadowCol = color.RGBA{255, 255, 255, 255}
	}
	
	if defaultTT != nil {
		scaledFace, err := opentype.NewFace(defaultTT, &opentype.FaceOptions{
			Size:    fontSize * float64(scale),
			DPI:     72,
			Hinting: font.HintingFull,
		})
		if err == nil {
			textWidth := font.MeasureString(scaledFace, text)
			if center {
				x = x - textWidth.Ceil()/2
			}
			
			point := fixed.Point26_6{
				X: fixed.Int26_6(x * 64),
				Y: fixed.Int26_6(y * 64),
			}
			
			// Draw shadow first if needed
			if shadowColor != "" {
				shadowPoint := fixed.Point26_6{
					X: fixed.Int26_6((x + scale) * 64),
					Y: fixed.Int26_6((y + scale) * 64),
				}
				d := &font.Drawer{
					Dst:  img,
					Src:  image.NewUniform(shadowCol),
					Face: scaledFace,
					Dot:  shadowPoint,
				}
				d.DrawString(text)
			}
			
			// Draw primary text
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
	
	// Fallback to basicfont
	face := basicfont.Face7x13
	fontScale := float64(fontSize) / 13.0
	textWidth := int(float64(len(text)*7) * fontScale)
	
	if center {
		x = x - (textWidth*scale)/2
	}
	
	tempImg := image.NewRGBA(image.Rect(0, 0, textWidth+2, int(13*fontScale)+2))
	point := fixed.Point26_6{
		X: fixed.Int26_6(1 * 64),
		Y: fixed.Int26_6(int(13*fontScale) * 64),
	}
	
	d := &font.Drawer{
		Dst:  tempImg,
		Src:  image.NewUniform(textCol),
		Face: face,
		Dot:  point,
	}
	d.DrawString(text)
	
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

func drawAxis3D(img *image.RGBA, centerX, centerY, scaleFactor, angleY, angleX, angleZ, r, g, b float64, axisColor color.RGBA, label string) {
	// 3D rotation same as main function
	x1 := r*math.Cos(angleY) - b*math.Sin(angleY)
	y1 := g
	z1 := r*math.Sin(angleY) + b*math.Cos(angleY)

	y2 := y1*math.Cos(angleX) - z1*math.Sin(angleX)

	// Project to 2D
	xProj := x1 * scaleFactor
	yProj := y2 * scaleFactor

	xEnd := int(centerX + xProj)
	yEnd := int(centerY - yProj)

	// Draw line from origin
	drawLine(img, int(centerX), int(centerY), xEnd, yEnd, axisColor, 3)

	// Draw label
	drawLabel(img, xEnd+8, yEnd, label, axisColor)
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
		// Draw thick line
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

func drawLabel(img *image.RGBA, x, y int, label string, col color.Color) {
	// Draw label as small squares
	for i := range label {
		px := x + i*10
		py := y
		// Draw a small square for visibility
		for dy := 0; dy < 4; dy++ {
			for dx := 0; dx < 8; dx++ {
				if px+dx >= 0 && px+dx < img.Bounds().Dx() && py+dy >= 0 && py+dy < img.Bounds().Dy() {
					img.Set(px+dx, py+dy, col)
				}
			}
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
