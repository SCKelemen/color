package color

import (
	"testing"
)

// Benchmark color space conversions
func BenchmarkRGBToOKLCH(b *testing.B) {
	c := RGB(0.5, 0.6, 0.7)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ToOKLCH(c)
	}
}

func BenchmarkRGBToLAB(b *testing.B) {
	c := RGB(0.5, 0.6, 0.7)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ToLAB(c)
	}
}

func BenchmarkRGBToHSL(b *testing.B) {
	c := RGB(0.5, 0.6, 0.7)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ToHSL(c)
	}
}

func BenchmarkRGBToXYZ(b *testing.B) {
	c := RGB(0.5, 0.6, 0.7)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ToXYZ(c)
	}
}

// Benchmark color operations
func BenchmarkLighten(b *testing.B) {
	c := RGB(0.5, 0.3, 0.7)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Lighten(c, 0.2)
	}
}

func BenchmarkDarken(b *testing.B) {
	c := RGB(0.5, 0.3, 0.7)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Darken(c, 0.2)
	}
}

func BenchmarkSaturate(b *testing.B) {
	c := RGB(0.5, 0.3, 0.7)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Saturate(c, 0.3)
	}
}

func BenchmarkDesaturate(b *testing.B) {
	c := RGB(0.5, 0.3, 0.7)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Desaturate(c, 0.3)
	}
}

func BenchmarkMixOKLCH(b *testing.B) {
	c1 := RGB(1, 0, 0)
	c2 := RGB(0, 0, 1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MixOKLCH(c1, c2, 0.5)
	}
}

// Benchmark gradient generation
func BenchmarkGradient10Steps(b *testing.B) {
	start := RGB(1, 0, 0)
	end := RGB(0, 0, 1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Gradient(start, end, 10)
	}
}

func BenchmarkGradient100Steps(b *testing.B) {
	start := RGB(1, 0, 0)
	end := RGB(0, 0, 1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Gradient(start, end, 100)
	}
}

func BenchmarkGradientInSpaceOKLCH(b *testing.B) {
	start := RGB(1, 0, 0)
	end := RGB(0, 0, 1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GradientInSpace(start, end, 20, GradientOKLCH)
	}
}

func BenchmarkGradientInSpaceRGB(b *testing.B) {
	start := RGB(1, 0, 0)
	end := RGB(0, 0, 1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GradientInSpace(start, end, 20, GradientRGB)
	}
}

func BenchmarkGradientMultiStop(b *testing.B) {
	stops := []GradientStop{
		{Color: RGB(1, 0, 0), Position: 0.0},
		{Color: RGB(1, 1, 0), Position: 0.5},
		{Color: RGB(0, 0, 1), Position: 1.0},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GradientMultiStop(stops, 30, GradientOKLCH)
	}
}

// Benchmark wide-gamut operations
func BenchmarkConvertToRGBSpace(b *testing.B) {
	c := RGB(0.8, 0.2, 0.3)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ConvertToRGBSpace(c, "display-p3")
	}
}

func BenchmarkConvertFromRGBSpace(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ConvertFromRGBSpace(0.8, 0.2, 0.3, 1.0, "display-p3")
	}
}

func BenchmarkSpaceColorConvert(b *testing.B) {
	c := NewSpaceColor(DisplayP3Space, []float64{0.8, 0.2, 0.3}, 1.0)
	srgb, _ := GetSpace("srgb")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.ConvertTo(srgb)
	}
}

// Benchmark parsing
func BenchmarkParseColorHex(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ParseColor("#FF5733")
	}
}

func BenchmarkParseColorRGB(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ParseColor("rgb(255, 87, 51)")
	}
}

func BenchmarkParseColorOKLCH(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ParseColor("oklch(0.7 0.2 150)")
	}
}

func BenchmarkRGBToHex(b *testing.B) {
	c := RGB(0.8, 0.3, 0.2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RGBToHex(c)
	}
}

// Benchmark common workflows
func BenchmarkThemePaletteGeneration(b *testing.B) {
	base := RGB(0.23, 0.51, 0.96)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Lighten(base, 0.4)
		_ = Lighten(base, 0.3)
		_ = Lighten(base, 0.2)
		_ = Lighten(base, 0.1)
		_ = base
		_ = Darken(base, 0.1)
		_ = Darken(base, 0.2)
		_ = Darken(base, 0.3)
		_ = Darken(base, 0.4)
	}
}

func BenchmarkColorDifferenceWorkflow(b *testing.B) {
	c1 := RGB(1, 0, 0)
	c2 := RGB(0.95, 0.05, 0.05)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = DeltaE2000(c1, c2)
	}
}

func BenchmarkGamutMappingWorkflow(b *testing.B) {
	vivid := NewOKLCH(0.7, 0.35, 150, 1.0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MapToGamut(vivid, GamutPreserveLightness)
	}
}

// Benchmark memory allocations
func BenchmarkColorCreation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RGB(0.5, 0.6, 0.7)
	}
}

func BenchmarkNewOKLCH(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewOKLCH(0.7, 0.2, 150, 1.0)
	}
}

// Benchmark stdlib compatibility
func BenchmarkToStdColor(b *testing.B) {
	c := RGB(0.5, 0.6, 0.7)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ToStdColor(c)
	}
}

func BenchmarkFromStdColor(b *testing.B) {
	stdC := ToStdColor(RGB(0.5, 0.6, 0.7))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FromStdColor(stdC)
	}
}
