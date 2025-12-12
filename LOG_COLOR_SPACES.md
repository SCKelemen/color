# LOG Color Spaces Guide

Professional LOG color spaces for cinema camera workflows with HDR support.

## Overview

LOG (logarithmic) color spaces are used in professional cinema cameras to capture a wider dynamic range than standard video. They encode linear light values logarithmically, similar to how human perception works, allowing 12-14 stops of dynamic range to be captured in 10-bit or 12-bit recording formats.

This library provides full support for 6 major cinema camera LOG formats with accurate transfer functions and wide-gamut primaries.

## Supported LOG Formats

### Canon C-Log

**Cameras:** C300, C500, C700, C70, R5 C, R7 C
**Gamut:** Cinema Gamut (wider than Rec.709, similar to DCI-P3)
**Dynamic Range:** ~12-14 stops

```go
// Load C-Log footage
clog := color.NewSpaceColor(color.CLogSpace,
    []float64{0.45, 0.40, 0.35}, 1.0)

// Also accessible via registry
clog, _ := color.GetSpace("c-log")      // Primary name
clog, _ := color.GetSpace("clog")       // Alias
```

**Key Characteristics:**
- 18% gray encodes to approximately 0.34
- Wide Cinema Gamut primaries
- Optimized for Canon cinema cameras
- Good highlight rolloff

### Sony S-Log3

**Cameras:** FX6, FX9, FX3, Venice, BURANO, A7S III, A1
**Gamut:** S-Gamut3 (very wide, designed for Rec.2020 delivery)
**Dynamic Range:** ~14+ stops

```go
// Load S-Log3 footage
slog3 := color.NewSpaceColor(color.SLog3Space,
    []float64{0.41, 0.39, 0.35}, 1.0)

// Registry access
slog3, _ := color.GetSpace("s-log3")    // Primary name
slog3, _ := color.GetSpace("slog3")     // Alias
```

**Key Characteristics:**
- 18% gray encodes to approximately 0.41 (41 IRE)
- S-Gamut3 primaries (wider than Rec.2020)
- Latest Sony LOG curve (preferred over S-Log2)
- Excellent for HDR workflows

### Panasonic V-Log

**Cameras:** GH5, GH6, S1H, S5 II, EVA1, Varicam
**Gamut:** V-Gamut (wide gamut optimized for cinema)
**Dynamic Range:** ~14 stops

```go
// Load V-Log footage
vlog := color.NewSpaceColor(color.VLogSpace,
    []float64{0.48, 0.45, 0.40}, 1.0)

// Registry access
vlog, _ := color.GetSpace("v-log")      // Primary name
vlog, _ := color.GetSpace("vlog")       // Alias
```

**Key Characteristics:**
- 18% gray encodes to approximately 0.42
- V-Gamut primaries
- Clean shadows
- Popular in hybrid cameras

### Arri LogC

**Cameras:** Alexa Mini, Alexa LF, Alexa 35, Amira
**Gamut:** Arri Wide Gamut (industry standard for cinema)
**Dynamic Range:** ~14 stops

```go
// Load Arri LogC footage (V3, EI 800)
logc := color.NewSpaceColor(color.ArriLogCSpace,
    []float64{0.42, 0.38, 0.35}, 1.0)

// Registry access
logc, _ := color.GetSpace("arri-logc")  // Primary name
logc, _ := color.GetSpace("logc")       // Alias
```

**Key Characteristics:**
- 18% gray encodes to approximately 0.38-0.39
- Arri Wide Gamut primaries
- Industry standard for cinema
- Excellent highlight handling
- This implementation uses LogC V3 at EI 800

### Red Log3G10

**Cameras:** Komodo, V-Raptor, Ranger, DSMC2 lineup
**Gamut:** RedWideGamutRGB (extremely wide)
**Dynamic Range:** ~16+ stops

```go
// Load Red Log3G10 footage
redlog := color.NewSpaceColor(color.RedLog3G10Space,
    []float64{0.46, 0.42, 0.38}, 1.0)

// Registry access
redlog, _ := color.GetSpace("red-log3g10")  // Primary name
redlog, _ := color.GetSpace("log3g10")      // Alias
```

**Key Characteristics:**
- Log3G10 uses 10-bit encoding
- RedWideGamutRGB primaries
- Massive dynamic range
- Optimized for HDR delivery

### Blackmagic Film

**Cameras:** Pocket Cinema 4K/6K, URSA Mini Pro
**Gamut:** Wide gamut (similar to Rec.2020)
**Dynamic Range:** ~13 stops

```go
// Load Blackmagic Film footage
bmdfilm := color.NewSpaceColor(color.BMDFilmSpace,
    []float64{0.44, 0.40, 0.36}, 1.0)

// Registry access
bmdfilm, _ := color.GetSpace("bmd-film")    // Primary name
bmdfilm, _ := color.GetSpace("bmdfilm")     // Alias
```

**Key Characteristics:**
- Wide gamut similar to Rec.2020
- Popular in independent filmmaking
- Good balance of DR and gradability

## Common Workflows

### Basic Conversion

```go
// Load LOG footage
slog3 := color.NewSpaceColor(color.SLog3Space,
    []float64{0.41, 0.39, 0.35}, 1.0)

// Convert to display space
rec709 := slog3.ConvertTo(color.Rec709Space)    // HDTV
rec2020 := slog3.ConvertTo(color.Rec2020Space)  // HDR/UHD
srgb := slog3.ConvertTo(color.SRGBSpace)        // Web
```

### HDR Workflow

```go
// Cinema camera footage in LOG
logcColor := color.NewSpaceColor(color.ArriLogCSpace,
    []float64{0.65, 0.60, 0.55}, 1.0)  // Bright highlight

// Convert to linear for processing
linear := logcColor.ConvertTo(color.SRGBLinearSpace)

// Values > 1.0 represent HDR highlights
r, g, b, _ := linear.RGBA()
fmt.Printf("Linear: %.2f, %.2f, %.2f (HDR)\n", r, g, b)

// Master in Rec.2020 for HDR delivery
hdrMaster := logcColor.ConvertTo(color.Rec2020Space)
```

### Color Grading Pipeline

```go
// 1. Load camera LOG footage
cameraLog := color.NewSpaceColor(color.VLogSpace,
    []float64{0.48, 0.45, 0.40}, 1.0)

// 2. Convert to perceptual space for grading
oklch := color.ToOKLCH(cameraLog)

// 3. Apply grading
graded := color.Lighten(oklch, 0.05)
graded = color.Saturate(graded, 0.15)

// 4. Convert to delivery format
// HDR delivery (Rec.2020)
hdrOut := color.ConvertToRGBSpace(graded, "rec2020")

// SDR delivery (sRGB)
sdrOut := color.ConvertToRGBSpace(graded, "srgb")
```

### Batch Processing

```go
// Process multiple frames
for _, frame := range frames {
    // Load LOG values for each pixel
    logColor := color.NewSpaceColor(color.SLog3Space,
        []float64{frame.r, frame.g, frame.b}, 1.0)

    // Convert to output color space
    output := logColor.ConvertTo(color.Rec709Space)

    // Write to output
    r, g, b, _ := output.RGBA()
    writePixel(r, g, b)
}
```

## Metadata Access

```go
// Get metadata about LOG color spaces
meta := color.Metadata(color.SLog3Space)

fmt.Println(meta.Name)                      // "s-log3"
fmt.Println(meta.Family)                    // "RGB-LOG"
fmt.Println(meta.IsHDR)                     // true
fmt.Println(meta.GamutVolumeRelativeToSRGB) // 1.69 (69% wider)
```

## Technical Details

### Transfer Functions

All LOG transfer functions implement both encoding (linear → LOG) and decoding (LOG → linear) with accurate round-trip conversion:

```go
// Transfer functions preserve values through round-trip
linear := 0.18  // 18% gray
encoded := cLogTransfer(linear)     // Encode to LOG
decoded := cLogInverseTransfer(encoded)  // Decode back
// decoded ≈ 0.18 (within floating-point precision)
```

### HDR Support

All LOG color spaces support HDR values > 1.0:

```go
// LOG spaces can represent very bright values
bright := color.NewSpaceColor(color.ArriLogCSpace,
    []float64{0.80, 0.75, 0.70}, 1.0)

// Convert to linear
linear := bright.ConvertTo(color.SRGBLinearSpace)
r, g, b, _ := linear.RGBA()

// r, g, b can be > 1.0 (HDR highlights)
if r > 1.0 {
    fmt.Println("HDR highlight detected")
}
```

### Gamut Boundaries

LOG color spaces use wide-gamut primaries that extend beyond sRGB:

| LOG Space | Gamut Volume vs sRGB |
|-----------|---------------------|
| Canon C-Log | 1.56× (56% wider) |
| Sony S-Log3 | 1.69× (69% wider) |
| Panasonic V-Log | 1.58× (58% wider) |
| Arri LogC | 1.55× (55% wider) |
| Red Log3G10 | 1.68× (68% wider) |
| Blackmagic Film | 1.70× (70% wider) |

## Best Practices

### 1. Preserve LOG Until Final Output

```go
// ✅ Good: Keep in LOG as long as possible
logFootage := loadLogFootage()
// ... do non-color operations ...
finalOutput := logFootage.ConvertTo(deliverySpace)

// ❌ Bad: Converting too early loses information
rec709 := logFootage.ConvertTo(color.Rec709Space)
// ... operations lose LOG dynamic range ...
```

### 2. Use Appropriate Delivery Space

```go
// HDR delivery
hdr := logColor.ConvertTo(color.Rec2020Space)

// SDR delivery (web, broadcast)
sdr := logColor.ConvertTo(color.Rec709Space)

// Web/mobile
web := logColor.ConvertTo(color.SRGBSpace)
```

### 3. Validate Input Ranges

```go
// LOG values typically in [0, 1] range
if r < 0 || r > 1 || g < 0 || g > 1 || b < 0 || b > 1 {
    fmt.Println("Warning: LOG values outside normal range")
}
```

### 4. Consider Gamut Mapping

```go
// When converting to narrow gamuts, use gamut mapping
logColor := color.NewSpaceColor(color.SLog3Space, vals, 1.0)

// Convert with gamut mapping strategy
srgbColor := logColor.ConvertTo(color.SRGBSpace)
if !color.InGamut(srgbColor) {
    mapped := color.MapToGamut(srgbColor, color.GamutPreserveLightness)
}
```

## Testing & Validation

All LOG implementations include:

- **Round-trip validation** - Encode → Decode returns original value
- **Reference value tests** - 18% gray encodes to documented values
- **HDR value support** - Values > 1.0 handled correctly
- **Edge case handling** - Negative values, transparency, extremes

See `spaces_log_test.go` for comprehensive test suite.

## References

- Canon C-Log: Canon Cinema EOS documentation
- Sony S-Log3: Sony CineAlta technical documentation
- Panasonic V-Log: Varicam V-Log/V-Gamut white paper
- Arri LogC: Arri ALEXA LogC specification
- Red Log3G10: Red Digital Cinema white papers
- Blackmagic Film: Blackmagic Design camera manual

## See Also

- [COLOR_SPACE_ARCHITECTURE.md](COLOR_SPACE_ARCHITECTURE.md) - Color space system design
- [QUICKSTART.md](QUICKSTART.md) - Getting started guide
- [examples_log_test.go](examples_log_test.go) - Executable examples
