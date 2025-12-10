# LOG and RAW Color Space Implementation Plan

This document outlines the plan for implementing LOG (logarithmic) and RAW (linear) color space support in the color library.

## Overview

LOG and RAW color spaces are essential for professional video/film production and raw image processing. They represent colors in a way that preserves more information in highlights and shadows, making them ideal for color grading and post-production workflows.

## LOG Color Spaces

### What are LOG Color Spaces?

LOG (logarithmic) color spaces apply a logarithmic transfer function to linear light values. This compresses the dynamic range, allowing more information to be stored in the same bit depth. LOG spaces are commonly used in:

- **Video Production**: C-Log (Canon), S-Log (Sony), V-Log (Panasonic), Arri LogC
- **Film Production**: Film scanning and digital intermediate workflows
- **Color Grading**: Professional color correction and grading pipelines

### Key Characteristics

1. **Logarithmic Transfer Function**: Compresses highlights and expands shadows
2. **Wide Dynamic Range**: Preserves more information than gamma-encoded spaces
3. **Camera-Specific**: Each camera manufacturer has their own LOG curve
4. **Linear-Like**: Appears flat/desaturated, requiring color grading

### Common LOG Color Spaces

| LOG Space | Manufacturer | Transfer Function | Use Case |
|-----------|--------------|-------------------|----------|
| **C-Log** | Canon | Canon-specific curve | Canon cinema cameras |
| **C-Log 2** | Canon | Updated C-Log curve | Canon cinema cameras (newer) |
| **C-Log 3** | Canon | Latest C-Log curve | Canon cinema cameras (latest) |
| **S-Log** | Sony | Sony-specific curve | Sony cinema cameras |
| **S-Log 2** | Sony | Updated S-Log curve | Sony cinema cameras (newer) |
| **S-Log 3** | Sony | Latest S-Log curve | Sony cinema cameras (latest) |
| **V-Log** | Panasonic | Panasonic-specific curve | Panasonic cinema cameras |
| **Arri LogC** | Arri | Arri-specific curve | Arri cinema cameras |
| **Red LogFilm** | Red | Red-specific curve | Red cinema cameras |
| **BMDFilm** | Blackmagic | Blackmagic-specific curve | Blackmagic cameras |

## RAW Color Spaces

### What are RAW Color Spaces?

RAW color spaces represent unprocessed sensor data from digital cameras. They are:

1. **Linear**: No gamma encoding applied
2. **Camera-Specific**: Each camera has its own color filter array (CFA) and sensor characteristics
3. **Wide Gamut**: Often exceed standard RGB gamuts
4. **High Bit Depth**: Typically 12-16 bits per channel

### Key Characteristics

1. **Linear Light**: Direct representation of sensor data
2. **No Color Space**: RAW data doesn't have a color space until demosaiced
3. **White Balance**: Applied during RAW processing
4. **Color Matrix**: Camera-specific conversion to RGB

### Common RAW Color Spaces

| RAW Format | Manufacturer | Color Space | Notes |
|-----------|--------------|-------------|-------|
| **Canon RAW** | Canon | Canon RGB | Canon-specific primaries |
| **Nikon RAW** | Nikon | Nikon RGB | Nikon-specific primaries |
| **Sony RAW** | Sony | Sony RGB | Sony-specific primaries |
| **Fuji RAW** | Fuji | Fuji RGB | Fuji-specific primaries |
| **Adobe DNG** | Adobe | Various | Standardized RAW format |

## Implementation Plan

### Phase 1: LOG Color Space Support

#### 1.1 Core LOG Infrastructure

```go
// LOGColorSpace represents a logarithmic color space
type LOGColorSpace struct {
    Name string
    // Transfer function: linear -> log
    LogTransfer func(float64) float64
    // Inverse transfer function: log -> linear
    LinearTransfer func(float64) float64
    // Base RGB color space (usually Rec. 2020 or camera-specific)
    BaseRGBSpace *RGBColorSpace
    // Black point (minimum value in log space)
    BlackPoint float64
    // White point (maximum value in log space)
    WhitePoint float64
}
```

#### 1.2 LOG Transfer Functions

Each LOG space needs its specific transfer function:

**C-Log Transfer Function:**
```go
func cLogTransfer(linear float64) float64 {
    // C-Log formula
    if linear < 0.0227218 {
        return 0.0928 + 0.4691 * math.Log10(linear + 0.0379)
    }
    return 0.0928 + 0.4691 * math.Log10(linear + 0.0379)
}
```

**S-Log 2 Transfer Function:**
```go
func sLog2Transfer(linear float64) float64 {
    // S-Log 2 formula
    if linear < 0 {
        return 0
    }
    return (0.432699 * math.Log10(linear * 0.616596 + 0.03) + 0.616596) / 0.432699
}
```

**S-Log 3 Transfer Function:**
```go
func sLog3Transfer(linear float64) float64 {
    // S-Log 3 formula
    if linear < 0.01125000 {
        return (linear * 171.2102946929 - 0.01125000) * 0.125
    }
    return (math.Log10((linear + 0.01) / (0.18 + 0.01)) * 0.432699 + 0.616596) / 0.432699
}
```

#### 1.3 LOG Color Type

```go
// LOG represents a color in a logarithmic color space
type LOG struct {
    R, G, B, A_ float64 // A_ to avoid conflict with Alpha() method
    Space *LOGColorSpace
}

// NewLOG creates a new LOG color
func NewLOG(r, g, b, alpha float64, space *LOGColorSpace) *LOG {
    return &LOG{
        R: r,
        G: g,
        B: b,
        A_: clamp01(alpha),
        Space: space,
    }
}

// RGBA converts LOG to RGBA via linear RGB
func (c *LOG) RGBA() (r, g, b, a float64) {
    // Convert LOG to linear RGB
    linearR := c.Space.LinearTransfer(c.R)
    linearG := c.Space.LinearTransfer(c.G)
    linearB := c.Space.LinearTransfer(c.B)
    
    // Convert linear RGB to sRGB (or base color space)
    // This would use the base RGB space conversion
    return convertLinearToSRGB(linearR, linearG, linearB, c.A_)
}

// ToLOG converts a color to LOG space
func ToLOG(c Color, space *LOGColorSpace) *LOG {
    // Convert to linear RGB first
    linearR, linearG, linearB, alpha := toLinearRGB(c)
    
    // Apply LOG transfer function
    return NewLOG(
        space.LogTransfer(linearR),
        space.LogTransfer(linearG),
        space.LogTransfer(linearB),
        alpha,
        space,
    )
}
```

### Phase 2: RAW Color Space Support

#### 2.1 RAW Color Space Definition

```go
// RAWColorSpace represents a camera RAW color space
type RAWColorSpace struct {
    Name string
    // Camera manufacturer
    Manufacturer string
    // Color primaries (RGB to XYZ matrix)
    RGBToXYZMatrix [9]float64
    XYZToRGBMatrix [9]float64
    // White point
    WhitePoint [3]float64
    // Linear transfer (RAW is always linear)
    TransferFunc func(float64) float64 // Identity function
    InverseTransferFunc func(float64) float64 // Identity function
}
```

#### 2.2 RAW Color Type

```go
// RAW represents a color in a RAW color space
type RAW struct {
    R, G, B, A_ float64
    Space *RAWColorSpace
}

// NewRAW creates a new RAW color
func NewRAW(r, g, b, alpha float64, space *RAWColorSpace) *RAW {
    return &RAW{
        R: r,
        G: g,
        B: b,
        A_: clamp01(alpha),
        Space: space,
    }
}

// RGBA converts RAW to RGBA
func (c *RAW) RGBA() (r, g, b, a float64) {
    // RAW is linear, so convert directly via XYZ
    xyz := c.Space.ConvertRGBToXYZ(c.R, c.G, c.B, c.A_)
    // Convert XYZ to sRGB
    return xyz.RGBA()
}

// ToRAW converts a color to RAW space
func ToRAW(c Color, space *RAWColorSpace) *RAW {
    // Convert to XYZ first
    xyz := ToXYZ(c)
    // Convert XYZ to RAW RGB
    rawR, rawG, rawB, alpha := space.ConvertXYZToRGB(xyz)
    return NewRAW(rawR, rawG, rawB, alpha, space)
}
```

### Phase 3: Predefined LOG Spaces

```go
var (
    // Canon LOG spaces
    CLogSpace = &LOGColorSpace{
        Name: "c-log",
        LogTransfer: cLogTransfer,
        LinearTransfer: cLogInverseTransfer,
        BaseRGBSpace: rec2020Space, // C-Log typically uses Rec. 2020
        BlackPoint: 0.0928,
        WhitePoint: 1.0,
    }
    
    CLog2Space = &LOGColorSpace{
        Name: "c-log-2",
        // C-Log 2 specific parameters
    }
    
    CLog3Space = &LOGColorSpace{
        Name: "c-log-3",
        // C-Log 3 specific parameters
    }
    
    // Sony LOG spaces
    SLogSpace = &LOGColorSpace{
        Name: "s-log",
        LogTransfer: sLogTransfer,
        LinearTransfer: sLogInverseTransfer,
        BaseRGBSpace: rec2020Space,
    }
    
    SLog2Space = &LOGColorSpace{
        Name: "s-log-2",
        LogTransfer: sLog2Transfer,
        LinearTransfer: sLog2InverseTransfer,
        BaseRGBSpace: rec2020Space,
    }
    
    SLog3Space = &LOGColorSpace{
        Name: "s-log-3",
        LogTransfer: sLog3Transfer,
        LinearTransfer: sLog3InverseTransfer,
        BaseRGBSpace: rec2020Space,
    }
    
    // Panasonic LOG
    VLogSpace = &LOGColorSpace{
        Name: "v-log",
        LogTransfer: vLogTransfer,
        LinearTransfer: vLogInverseTransfer,
        BaseRGBSpace: rec2020Space,
    }
    
    // Arri LOG
    ArriLogCSpace = &LOGColorSpace{
        Name: "arri-logc",
        LogTransfer: arriLogCTransfer,
        LinearTransfer: arriLogCInverseTransfer,
        BaseRGBSpace: rec2020Space,
    }
)
```

### Phase 4: Predefined RAW Spaces

```go
var (
    // Canon RAW
    CanonRAWSpace = &RAWColorSpace{
        Name: "canon-raw",
        Manufacturer: "Canon",
        RGBToXYZMatrix: [9]float64{
            // Canon-specific primaries
        },
        XYZToRGBMatrix: [9]float64{
            // Inverse matrix
        },
        WhitePoint: [3]float64{0.95047, 1.0, 1.08883}, // D65
        TransferFunc: linearTransfer,
        InverseTransferFunc: linearInverseTransfer,
    }
    
    // Nikon RAW
    NikonRAWSpace = &RAWColorSpace{
        Name: "nikon-raw",
        Manufacturer: "Nikon",
        // Nikon-specific primaries
    }
    
    // Sony RAW
    SonyRAWSpace = &RAWColorSpace{
        Name: "sony-raw",
        Manufacturer: "Sony",
        // Sony-specific primaries
    }
)
```

### Phase 5: Parsing Support

Add parsing for LOG and RAW colors:

```go
// Parse LOG colors
// Example: "c-log(0.5 0.3 0.2)"
// Example: "s-log-3(0.4 0.3 0.25)"

// Parse RAW colors
// Example: "canon-raw(0.8 0.6 0.4)"
// Example: "raw(0.8 0.6 0.4)" with space parameter
```

### Phase 6: Integration with Existing Systems

1. **Update Color Interface**: Ensure LOG and RAW types implement the Color interface
2. **Conversion Functions**: Add conversion functions between LOG/RAW and other spaces
3. **Gradient Support**: Support LOG/RAW in gradient generation
4. **Documentation**: Update README and examples

## Implementation Considerations

### 1. Transfer Function Accuracy

LOG transfer functions are complex and camera-specific. We need:
- Accurate mathematical formulas from camera manufacturers
- Proper handling of edge cases (black point, white point)
- Support for different bit depths (10-bit, 12-bit, 16-bit)

### 2. Color Space Conversion

LOG colors need to be converted through:
1. LOG → Linear RGB (via inverse transfer function)
2. Linear RGB → Base RGB space (usually Rec. 2020)
3. Base RGB space → XYZ
4. XYZ → Target color space

### 3. Bit Depth Handling

RAW and LOG spaces often use higher bit depths:
- Support for 10-bit, 12-bit, 16-bit values
- Normalization to [0, 1] range for internal representation
- Proper quantization when converting to 8-bit displays

### 4. White Balance

RAW colors may need white balance correction:
- Support for different white balance settings
- Conversion between color temperatures
- Integration with color grading workflows

### 5. Performance

LOG/RAW conversions are computationally intensive:
- Optimize transfer function calculations
- Cache conversion matrices
- Consider SIMD optimizations for batch processing

## Testing Strategy

1. **Unit Tests**: Test each LOG/RAW transfer function
2. **Round-Trip Tests**: Verify conversions preserve color information
3. **Reference Tests**: Compare against known reference values
4. **Integration Tests**: Test with real-world color grading workflows

## Future Enhancements

1. **LUT Support**: Load and apply Look-Up Tables (LUTs) for LOG spaces
2. **Color Grading Tools**: Helper functions for common color grading operations
3. **Camera Profiles**: Predefined profiles for specific camera models
4. **HDR Support**: Extended dynamic range for LOG/RAW workflows

## References

- [Canon C-Log White Paper](https://www.usa.canon.com/internet/portal/us/home/explore/product-showcases/cinema-eos/cinema-eos-features/c-log)
- [Sony S-Log Technical Paper](https://pro.sony/ue_US/products/broadcast-cameras/s-log-white-paper)
- [Arri LogC Specification](https://www.arri.com/en/learn-help/learn-help-camera-system/camera-workflow/working-with-arri-look-files)
- [Adobe DNG Specification](https://helpx.adobe.com/camera-raw/using/adobe-dng-converter.html)

## Example Usage

```go
// Create a LOG color
logColor := color.NewLOG(0.5, 0.3, 0.2, 1.0, color.CLog3Space)

// Convert to standard RGB for display
rgb := logColor.RGBA()

// Convert from RGB to LOG for color grading
rgbColor := color.RGB(1.0, 0.5, 0.0)
logColor2 := color.ToLOG(rgbColor, color.SLog3Space)

// Work with RAW colors
rawColor := color.NewRAW(0.8, 0.6, 0.4, 1.0, color.CanonRAWSpace)
oklch := color.ToOKLCH(rawColor) // Convert RAW to OKLCH for manipulation
```

## Conclusion

Implementing LOG and RAW color space support will make this library suitable for professional video/film production and raw image processing workflows. The implementation should prioritize accuracy, performance, and ease of use.

