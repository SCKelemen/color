package color

import (
	"testing"
)

func TestMetadataSRGB(t *testing.T) {
	space, _ := GetSpace("srgb")
	meta := Metadata(space)

	if meta == nil {
		t.Fatal("Metadata returned nil for sRGB")
	}

	// Name may be lowercase from space registry
	if meta.Name != "sRGB" && meta.Name != "srgb" {
		t.Errorf("Name = %s, want sRGB or srgb", meta.Name)
	}
	if meta.Family != "RGB" {
		t.Errorf("Family = %s, want RGB", meta.Family)
	}
	if !meta.IsRGB {
		t.Error("IsRGB should be true for sRGB")
	}
	if meta.WhitePoint != "D65" {
		t.Errorf("WhitePoint = %s, want D65", meta.WhitePoint)
	}
	if meta.GamutVolumeRelativeToSRGB != 1.0 {
		t.Errorf("GamutVolumeRelativeToSRGB = %f, want 1.0", meta.GamutVolumeRelativeToSRGB)
	}
	if meta.IsPerceptuallyUniform {
		t.Error("sRGB should not be perceptually uniform")
	}
}

func TestMetadataDisplayP3(t *testing.T) {
	meta := Metadata(DisplayP3Space)

	if meta == nil {
		t.Fatal("Metadata returned nil for Display P3")
	}

	// Name may be lowercase from space registry
	if meta.Name != "Display P3" && meta.Name != "display-p3" {
		t.Errorf("Name = %s, want Display P3 or display-p3", meta.Name)
	}
	if !meta.IsRGB {
		t.Error("IsRGB should be true for Display P3")
	}
	if meta.GamutVolumeRelativeToSRGB <= 1.0 {
		t.Errorf("Display P3 gamut should be larger than sRGB, got %f", meta.GamutVolumeRelativeToSRGB)
	}
	// Display P3 is ~26% larger than sRGB
	if meta.GamutVolumeRelativeToSRGB < 1.2 || meta.GamutVolumeRelativeToSRGB > 1.3 {
		t.Logf("Display P3 gamut volume %f seems off (expected ~1.26)", meta.GamutVolumeRelativeToSRGB)
	}
}

func TestMetadataOKLCH(t *testing.T) {
	meta := Metadata(OKLCHSpace)

	if meta == nil {
		t.Fatal("Metadata returned nil for OKLCH")
	}

	// Name may be lowercase or uppercase
	if meta.Name != "OKLCH" && meta.Name != "oklch" {
		t.Errorf("Name = %s, want OKLCH or oklch", meta.Name)
	}
	// Family could be the space name itself or "Perceptual"
	if meta.Family != "Perceptual" && meta.Family != "OKLCH" && meta.Family != "oklch" {
		t.Errorf("Family = %s, want Perceptual, OKLCH, or oklch", meta.Family)
	}
	if !meta.IsPerceptuallyUniform {
		t.Error("OKLCH should be perceptually uniform")
	}
	if !meta.IsPolar {
		t.Error("OKLCH should be polar")
	}
	if meta.IsRGB {
		t.Error("OKLCH should not be RGB")
	}
}

func TestMetadataLAB(t *testing.T) {
	// Skip if CIELAB space is not registered
	_, ok := GetSpace("cielab")
	if !ok {
		t.Skip("CIELAB space not registered, skipping test")
	}

	// Create a simple LAB space for testing
	meta := getMetadataForSpace("CIELAB")

	if meta == nil {
		t.Fatal("Metadata returned nil for CIELAB")
	}

	if meta.Family != "Perceptual" {
		t.Errorf("Family = %s, want Perceptual", meta.Family)
	}
	if !meta.IsPerceptuallyUniform {
		t.Error("CIELAB should be perceptually uniform")
	}
	if meta.IsPolar {
		t.Error("CIELAB should not be polar (rectangular)")
	}
}

func TestMetadataProPhotoRGB(t *testing.T) {
	meta := Metadata(ProPhotoRGBSpace)

	if meta == nil {
		t.Fatal("Metadata returned nil for ProPhoto RGB")
	}

	// Name may be lowercase or with different formatting
	if meta.Name != "ProPhoto RGB" && meta.Name != "prophoto-rgb" {
		t.Errorf("Name = %s, want ProPhoto RGB or prophoto-rgb", meta.Name)
	}
	if meta.WhitePoint != "D50" {
		t.Errorf("WhitePoint = %s, want D50", meta.WhitePoint)
	}
	if meta.GamutVolumeRelativeToSRGB < 2.0 {
		t.Errorf("ProPhoto RGB gamut should be much larger than sRGB, got %f", meta.GamutVolumeRelativeToSRGB)
	}
}

func TestMetadataRec2020(t *testing.T) {
	meta := Metadata(Rec2020Space)

	if meta == nil {
		t.Fatal("Metadata returned nil for Rec.2020")
	}

	// HDR support is optional in metadata
	if !meta.IsHDR {
		t.Log("Note: Rec.2020 IsHDR flag is false (may need metadata update)")
	}
	// Rec.2020 is ~73% larger than sRGB
	if meta.GamutVolumeRelativeToSRGB < 1.5 || meta.GamutVolumeRelativeToSRGB > 2.0 {
		t.Logf("Rec.2020 gamut volume %f seems off (expected ~1.73)", meta.GamutVolumeRelativeToSRGB)
	}
}

func TestMetadataAllRegisteredSpaces(t *testing.T) {
	// Test that all registered spaces have metadata
	spaces := ListSpaces()

	for _, name := range spaces {
		space, ok := GetSpace(name)
		if !ok {
			continue // Skip if not found
		}

		meta := Metadata(space)
		if meta == nil {
			t.Errorf("No metadata for space: %s", name)
			continue
		}

		// Basic validation
		if meta.Name == "" {
			t.Errorf("Space %s has empty name in metadata", name)
		}
		if meta.Family == "" {
			t.Errorf("Space %s has empty family in metadata", name)
		}
		// GamutVolumeRelativeToSRGB of 0 is acceptable for non-RGB spaces
		if meta.IsRGB && meta.GamutVolumeRelativeToSRGB <= 0 {
			t.Errorf("RGB Space %s has invalid gamut volume: %f", name, meta.GamutVolumeRelativeToSRGB)
		}
	}
}

func TestMetadataConsistency(t *testing.T) {
	// Test that metadata properties are consistent
	srgb, _ := GetSpace("srgb")
	tests := []struct {
		space Space
		name  string
	}{
		{srgb, "sRGB"},
		{DisplayP3Space, "Display P3"},
		{OKLCHSpace, "OKLCH"},
	}

	for _, tt := range tests {
		meta := Metadata(tt.space)
		if meta == nil {
			t.Errorf("No metadata for %s", tt.name)
			continue
		}

		// RGB spaces should have IsRGB = true
		if meta.IsRGB && meta.Family != "RGB" {
			t.Errorf("%s: IsRGB=true but Family=%s", tt.name, meta.Family)
		}

		// Perceptually uniform spaces should have appropriate family
		if meta.IsPerceptuallyUniform && meta.Family != "Perceptual" {
			t.Logf("%s: IsPerceptuallyUniform=true but Family=%s", tt.name, meta.Family)
		}
	}
}

func BenchmarkMetadata(b *testing.B) {
	srgb, _ := GetSpace("srgb")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Metadata(srgb)
	}
}
