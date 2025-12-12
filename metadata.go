package color

// SpaceMetadata provides information about a color space's properties.
type SpaceMetadata struct {
	// Name of the color space
	Name string

	// Family indicates the type of color space (RGB, Lab, LCH, etc.)
	Family string

	// IsRGB indicates if this is an RGB-based color space
	IsRGB bool

	// IsHDR indicates if this space supports High Dynamic Range (values > 1.0)
	IsHDR bool

	// WhitePoint indicates the reference white point (e.g., "D65", "D50")
	WhitePoint string

	// GamutVolumeRelativeToSRGB is the gamut volume relative to sRGB (1.0 = sRGB)
	// Display P3 ≈ 1.26, Rec.2020 ≈ 1.73, ProPhoto RGB ≈ 2.89
	GamutVolumeRelativeToSRGB float64

	// IsPerceptuallyUniform indicates if the space is designed for perceptual uniformity
	IsPerceptuallyUniform bool

	// IsPolar indicates if this is a cylindrical/polar color space (has hue component)
	IsPolar bool
}

// Metadata returns metadata about a color space.
// Returns nil if metadata is not available for this space.
func Metadata(space Space) *SpaceMetadata {
	return getMetadataForSpace(space.Name())
}

// getMetadataForSpace returns metadata for a given space name.
func getMetadataForSpace(name string) *SpaceMetadata {
	switch name {
	case "sRGB":
		return &SpaceMetadata{
			Name:                      "sRGB",
			Family:                    "RGB",
			IsRGB:                     true,
			IsHDR:                     false,
			WhitePoint:                "D65",
			GamutVolumeRelativeToSRGB: 1.0,
			IsPerceptuallyUniform:     false,
			IsPolar:                   false,
		}
	case "sRGB-linear":
		return &SpaceMetadata{
			Name:                      "sRGB-linear",
			Family:                    "RGB",
			IsRGB:                     true,
			IsHDR:                     true,
			WhitePoint:                "D65",
			GamutVolumeRelativeToSRGB: 1.0,
			IsPerceptuallyUniform:     false,
			IsPolar:                   false,
		}
	case "display-p3":
		return &SpaceMetadata{
			Name:                      "display-p3",
			Family:                    "RGB",
			IsRGB:                     true,
			IsHDR:                     false,
			WhitePoint:                "D65",
			GamutVolumeRelativeToSRGB: 1.26,
			IsPerceptuallyUniform:     false,
			IsPolar:                   false,
		}
	case "dci-p3":
		return &SpaceMetadata{
			Name:                      "dci-p3",
			Family:                    "RGB",
			IsRGB:                     true,
			IsHDR:                     false,
			WhitePoint:                "D65",
			GamutVolumeRelativeToSRGB: 1.26,
			IsPerceptuallyUniform:     false,
			IsPolar:                   false,
		}
	case "a98-rgb":
		return &SpaceMetadata{
			Name:                      "a98-rgb",
			Family:                    "RGB",
			IsRGB:                     true,
			IsHDR:                     false,
			WhitePoint:                "D65",
			GamutVolumeRelativeToSRGB: 1.44,
			IsPerceptuallyUniform:     false,
			IsPolar:                   false,
		}
	case "prophoto-rgb":
		return &SpaceMetadata{
			Name:                      "prophoto-rgb",
			Family:                    "RGB",
			IsRGB:                     true,
			IsHDR:                     false,
			WhitePoint:                "D50",
			GamutVolumeRelativeToSRGB: 2.89,
			IsPerceptuallyUniform:     false,
			IsPolar:                   false,
		}
	case "rec2020":
		return &SpaceMetadata{
			Name:                      "rec2020",
			Family:                    "RGB",
			IsRGB:                     true,
			IsHDR:                     false,
			WhitePoint:                "D65",
			GamutVolumeRelativeToSRGB: 1.73,
			IsPerceptuallyUniform:     false,
			IsPolar:                   false,
		}
	case "rec709":
		return &SpaceMetadata{
			Name:                      "rec709",
			Family:                    "RGB",
			IsRGB:                     true,
			IsHDR:                     false,
			WhitePoint:                "D65",
			GamutVolumeRelativeToSRGB: 1.0, // Essentially same as sRGB
			IsPerceptuallyUniform:     false,
			IsPolar:                   false,
		}
	case "OKLCH":
		return &SpaceMetadata{
			Name:                      "OKLCH",
			Family:                    "OKLCH",
			IsRGB:                     false,
			IsHDR:                     true,
			WhitePoint:                "D65",
			GamutVolumeRelativeToSRGB: 0, // Not applicable for non-RGB spaces
			IsPerceptuallyUniform:     true,
			IsPolar:                   true,
		}
	default:
		return nil
	}
}
