package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
)

// generateStaticImages generates high-quality static PNG images of color models
func generateStaticImages() error {
	outputDir := "docs/models"
	os.MkdirAll(outputDir, 0755)

	// Generate static images for each model
	// Use frame 0 (or a good representative frame) for each
	models := []struct {
		name  string
		frame int // Which frame to use as the static image
		fn    func(frame int, totalFrames int) *image.RGBA
	}{
		{"rgb_cube", 0, generateRGBCubeFrame},
		{"hsl_cylinder", 0, generateHSLCylinderFrame},
		{"lab_space", 0, generateLABSpaceFrame},
		{"oklch_space", 0, generateOKLCHSpaceFrame},
	}

	for _, m := range models {
		// Generate high-quality image (we could increase resolution here if needed)
		img := m.fn(m.frame, numFrames)

		filename := filepath.Join(outputDir, fmt.Sprintf("model_%s_static.png", m.name))
		f, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("failed to create %s: %w", filename, err)
		}

		// Encode as PNG with best compression
		encoder := &png.Encoder{
			CompressionLevel: png.BestCompression,
		}

		if err := encoder.Encode(f, img); err != nil {
			f.Close()
			return fmt.Errorf("failed to encode %s: %w", filename, err)
		}

		f.Close()
		fmt.Printf("Generated high-quality static image: %s\n", filename)
	}

	return nil
}
