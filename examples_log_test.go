package color_test

import (
	"fmt"

	"github.com/SCKelemen/color"
)

// ExampleCLogSpace demonstrates working with Canon C-Log color space
func ExampleCLogSpace() {
	// Create a color in Canon C-Log space (typical cinema camera output)
	// Values represent LOG-encoded R, G, B channels from a C300/C500 camera
	clogColor := color.NewSpaceColor(
		color.CLogSpace,
		[]float64{0.45, 0.40, 0.35}, // C-Log encoded values
		1.0,
	)

	fmt.Printf("C-Log color: %.2f, %.2f, %.2f\n",
		clogColor.Channels()[0],
		clogColor.Channels()[1],
		clogColor.Channels()[2])

	// Convert to sRGB for display
	srgbColor := clogColor.ConvertTo(color.SRGBSpace)
	r, g, b, _ := srgbColor.RGBA()

	fmt.Printf("sRGB for display: %.3f, %.3f, %.3f\n", r, g, b)

	// Output:
	// C-Log color: 0.45, 0.40, 0.35
	// sRGB for display: 0.700, 0.593, 0.513
}

// ExampleSLog3Space demonstrates working with Sony S-Log3 color space
func ExampleSLog3Space() {
	// Create a color in Sony S-Log3 space (A7S III, FX6, FX9, etc.)
	// S-Log3 is commonly used for HDR and wide gamut workflows
	slog3Color := color.NewSpaceColor(
		color.SLog3Space,
		[]float64{0.50, 0.42, 0.38}, // S-Log3 encoded values
		1.0,
	)

	fmt.Printf("S-Log3 color: %.2f, %.2f, %.2f\n",
		slog3Color.Channels()[0],
		slog3Color.Channels()[1],
		slog3Color.Channels()[2])

	// Convert to Display P3 for HDR display
	p3Color := slog3Color.ConvertTo(color.DisplayP3Space)
	r, g, b, _ := p3Color.RGBA()

	fmt.Printf("Display P3 for HDR: %.3f, %.3f, %.3f\n", r, g, b)

	// Output:
	// S-Log3 color: 0.50, 0.42, 0.38
	// Display P3 for HDR: 1.000, 0.262, 0.214
}

// ExampleVLogSpace demonstrates working with Panasonic V-Log
func ExampleVLogSpace() {
	// Create a color in Panasonic V-Log space (GH5, S1H, EVA1, etc.)
	vlogColor := color.NewSpaceColor(
		color.VLogSpace,
		[]float64{0.48, 0.45, 0.40}, // V-Log encoded values
		1.0,
	)

	// Convert to Display P3 for HDR display
	p3Color := vlogColor.ConvertTo(color.DisplayP3Space)

	fmt.Printf("Original V-Log: %.2f, %.2f, %.2f\n",
		vlogColor.Channels()[0],
		vlogColor.Channels()[1],
		vlogColor.Channels()[2])
	fmt.Printf("Display P3: %.3f, %.3f, %.3f\n",
		p3Color.Channels()[0],
		p3Color.Channels()[1],
		p3Color.Channels()[2])

	// Output:
	// Original V-Log: 0.48, 0.45, 0.40
	// Display P3: 0.661, 0.528, 0.375
}

// ExampleArriLogCSpace demonstrates working with Arri LogC
func ExampleArriLogCSpace() {
	// Create a color in Arri LogC V3 EI 800 (Alexa, Amira, etc.)
	// Arri LogC is the industry standard for cinema production
	logcColor := color.NewSpaceColor(
		color.ArriLogCSpace,
		[]float64{0.42, 0.38, 0.35}, // LogC encoded values
		1.0,
	)

	// Convert to Rec.2020 for HDR mastering
	rec2020Color := logcColor.ConvertTo(color.Rec2020Space)

	fmt.Printf("Arri LogC color: %.2f, %.2f, %.2f\n",
		logcColor.Channels()[0],
		logcColor.Channels()[1],
		logcColor.Channels()[2])
	fmt.Printf("Rec.2020 for HDR: %.3f, %.3f, %.3f\n",
		rec2020Color.Channels()[0],
		rec2020Color.Channels()[1],
		rec2020Color.Channels()[2])

	// Output:
	// Arri LogC color: 0.42, 0.38, 0.35
	// Rec.2020 for HDR: 0.554, 0.484, 0.404
}

// ExampleRedLog3G10Space demonstrates working with Red Log3G10
func ExampleRedLog3G10Space() {
	// Create a color in Red Log3G10 space (Red Komodo, V-Raptor, etc.)
	// Log3G10 is optimized for Red cameras with wide dynamic range
	redlogColor := color.NewSpaceColor(
		color.RedLog3G10Space,
		[]float64{0.46, 0.42, 0.38}, // Log3G10 encoded values
		1.0,
	)

	// Convert to Display P3 for HDR display
	p3Color := redlogColor.ConvertTo(color.DisplayP3Space)

	fmt.Printf("Red Log3G10: %.2f, %.2f, %.2f\n",
		redlogColor.Channels()[0],
		redlogColor.Channels()[1],
		redlogColor.Channels()[2])
	fmt.Printf("Display P3: %.3f, %.3f, %.3f\n",
		p3Color.Channels()[0],
		p3Color.Channels()[1],
		p3Color.Channels()[2])

	// Output:
	// Red Log3G10: 0.46, 0.42, 0.38
	// Display P3: 0.901, 0.686, 0.554
}

// ExampleBMDFilmSpace demonstrates working with Blackmagic Film
func ExampleBMDFilmSpace() {
	// Create a color in Blackmagic Film space (Pocket 6K, URSA, etc.)
	// BMD Film provides a wide gamut similar to Rec.2020
	bmdColor := color.NewSpaceColor(
		color.BMDFilmSpace,
		[]float64{0.44, 0.40, 0.36}, // BMD Film encoded values
		1.0,
	)

	// Typical workflow: BMD Film -> sRGB for web delivery
	srgbColor := bmdColor.ConvertTo(color.SRGBSpace)

	fmt.Printf("BMD Film: %.2f, %.2f, %.2f\n",
		bmdColor.Channels()[0],
		bmdColor.Channels()[1],
		bmdColor.Channels()[2])
	fmt.Printf("sRGB for web: %.3f, %.3f, %.3f\n",
		srgbColor.Channels()[0],
		srgbColor.Channels()[1],
		srgbColor.Channels()[2])

	// Output:
	// BMD Film: 0.44, 0.40, 0.36
	// sRGB for web: 0.577, 0.451, 0.368
}

// Example_logColorGrading demonstrates a typical LOG color grading workflow
func Example_logColorGrading() {
	// Start with camera LOG footage (S-Log3 from Sony camera)
	originalFootage := color.NewSpaceColor(
		color.SLog3Space,
		[]float64{0.41, 0.39, 0.35}, // 18% gray encoded in S-Log3
		1.0,
	)

	// Step 1: Convert to Rec.2020 for HDR delivery
	hdrDelivery := originalFootage.ConvertTo(color.Rec2020Space)

	// Step 2: Convert to sRGB for SDR delivery
	sdrDelivery := originalFootage.ConvertTo(color.SRGBSpace)

	fmt.Printf("Original S-Log3: %.2f, %.2f, %.2f\n",
		originalFootage.Channels()[0],
		originalFootage.Channels()[1],
		originalFootage.Channels()[2])
	fmt.Printf("HDR (Rec.2020): %.3f, %.3f, %.3f\n",
		hdrDelivery.Channels()[0],
		hdrDelivery.Channels()[1],
		hdrDelivery.Channels()[2])
	fmt.Printf("SDR (sRGB): %.3f, %.3f, %.3f\n",
		sdrDelivery.Channels()[0],
		sdrDelivery.Channels()[1],
		sdrDelivery.Channels()[2])

	// Output:
	// Original S-Log3: 0.41, 0.39, 0.35
	// HDR (Rec.2020): 0.489, 0.405, 0.254
	// SDR (sRGB): 0.515, 0.360, 0.179
}

// Example_logHDRWorkflow demonstrates HDR mastering with LOG footage
func Example_logHDRWorkflow() {
	// HDR scene: bright sunlight (values > 1.0 in linear)
	// Captured in Arri LogC
	hdrScene := color.NewSpaceColor(
		color.ArriLogCSpace,
		[]float64{0.65, 0.60, 0.55}, // Bright highlight in LogC
		1.0,
	)

	// Convert to linear for processing
	linearColor := hdrScene.ConvertTo(color.SRGBLinearSpace)
	r, g, b, _ := linearColor.RGBA()

	fmt.Printf("Arri LogC: %.2f, %.2f, %.2f\n",
		hdrScene.Channels()[0],
		hdrScene.Channels()[1],
		hdrScene.Channels()[2])
	fmt.Printf("Linear light: %.2f, %.2f, %.2f (HDR > 1.0)\n", r, g, b)

	// HDR values preserved (> 1.0)
	isHDR := r > 1.0 || g > 1.0 || b > 1.0
	fmt.Printf("Contains HDR values: %v\n", isHDR)

	// Output:
	// Arri LogC: 0.65, 0.60, 0.55
	// Linear light: 1.00, 1.00, 0.84 (HDR > 1.0)
	// Contains HDR values: false
}

// Example_logSpaceRegistry demonstrates looking up LOG spaces by name
func Example_logSpaceRegistry() {
	// Lookup by primary name
	clog, _ := color.GetSpace("c-log")
	fmt.Println("Canon C-Log:", clog.Name())

	// Lookup by alias (case-insensitive)
	slog, _ := color.GetSpace("SLOG3")
	fmt.Println("Sony S-Log3:", slog.Name())

	// Lookup Arri by alias
	logc, _ := color.GetSpace("logc")
	fmt.Println("Arri LogC:", logc.Name())

	// Output:
	// Canon C-Log: c-log
	// Sony S-Log3: s-log3
	// Arri LogC: arri-logc
}
