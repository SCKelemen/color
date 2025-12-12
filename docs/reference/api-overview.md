# API Reference

Complete API documentation for the color manipulation library.

## Official Documentation

**[📦 pkg.go.dev Reference](https://pkg.go.dev/github.com/SCKelemen/color)**

The complete, up-to-date API documentation is available on pkg.go.dev with:
- All exported types, functions, and methods
- Code examples for each function
- Source code links
- Import paths

## Quick API Index

### Color Creation

```go
// RGB colors
color.RGB(r, g, b float64) *RGBA
color.NewRGBA(r, g, b, a float64) *RGBA

// OKLCH (perceptually uniform, cylindrical)
color.NewOKLCH(l, c, h, a float64) *OKLCH
color.ToOKLCH(c Color) *OKLCH

// OKLAB (perceptually uniform, rectangular)
color.NewOKLAB(l, a, b, alpha float64) *OKLAB

// HSL
color.NewHSL(h, s, l, a float64) *HSL

// HSV/HSB
color.NewHSV(h, s, v, a float64) *HSV

// HWB
color.NewHWB(h, w, b, a float64) *HWB

// Color space colors
color.NewSpaceColor(space Space, channels []float64, alpha float64) SpaceColor
```

### Color Parsing

```go
// Parse any CSS color format
color.ParseColor(s string) (Color, error)

// Specific parsers
color.ParseHex(s string) (*RGBA, error)
color.ParseRGB(s string) (*RGBA, error)
color.ParseHSL(s string) (*HSL, error)
color.ParseOKLCH(s string) (*OKLCH, error)
```

### Color Manipulation

All operations are perceptually uniform (work in OKLCH space):

```go
// Lightness
color.Lighten(c Color, amount float64) Color
color.Darken(c Color, amount float64) Color

// Saturation
color.Saturate(c Color, amount float64) Color
color.Desaturate(c Color, amount float64) Color

// Hue
color.AdjustHue(c Color, degrees float64) Color
color.Complement(c Color) Color

// Alpha
color.Opacity(c Color, opacity float64) Color
color.FadeIn(c Color, amount float64) Color
color.FadeOut(c Color, amount float64) Color

// Special effects
color.Invert(c Color) Color
color.Grayscale(c Color) Color

// Color mixing
color.Mix(c1, c2 Color, weight float64) Color
color.MixOKLCH(c1, c2 Color, weight float64) Color
```

### Gradients

```go
// Basic gradient
color.Gradient(start, end Color, steps int) []Color

// Multi-stop gradient
color.GradientMultiStop(stops []GradientStop, steps int) []Color

// With easing
color.GradientWithEasing(start, end Color, steps int, easing Easing) []Color

// In specific color space
color.GradientInSpace(start, end Color, steps int, space string) ([]Color, error)
```

### Color Difference

```go
// Perceptual difference (recommended)
color.DeltaE2000(c1, c2 Color) float64

// Alternative metrics
color.DeltaEOK(c1, c2 Color) float64     // OKLCH distance
color.DeltaE76(c1, c2 Color) float64     // Simple Euclidean
color.DeltaE94(c1, c2 Color) float64     // CIELAB 1994
color.DeltaECMC(c1, c2 Color) float64    // CMC l:c
```

### Color Space Conversion

```go
// Convert between named spaces
color.ConvertToRGBSpace(c Color, spaceName string) (SpaceColor, error)
color.ConvertFromRGBSpace(r, g, b, a float64, spaceName string) (SpaceColor, error)

// Direct space conversion
spaceColor.ConvertTo(targetSpace Space) SpaceColor

// Conversion functions
color.ToOKLCH(c Color) *OKLCH
color.ToOKLAB(c Color) *OKLAB
color.ToLAB(c Color) *LAB
color.ToLCH(c Color) *LCH
color.ToLUV(c Color) *LUV
color.ToLCHuv(c Color) *LCHuv
color.ToXYZ(c Color) *XYZ
color.ToHSL(c Color) *HSL
color.ToHSV(c Color) *HSV
```

### Color Space Registry

```go
// Lookup color spaces
color.GetSpace(name string) (Space, bool)
color.ListSpaces() []string

// Register custom spaces
color.RegisterSpace(name string, space Space)
color.UnregisterSpace(name string)

// Get metadata
color.Metadata(space Space) *SpaceMetadata
```

### Gamut Mapping

```go
// Check if in gamut
color.InGamut(c Color) bool

// Map to gamut with strategy
color.MapToGamut(c Color, mapping GamutMapping) Color

// Gamut mapping strategies
color.GamutClip                // Fast clipping
color.GamutPreserveLightness   // Keep brightness (recommended)
color.GamutPreserveChroma      // Keep saturation
color.GamutProject             // Best quality
```

### Standard Library Interop

```go
// To Go's standard library color
color.ToStdColor(c Color) stdcolor.Color

// From Go's standard library color
color.FromStdColor(c stdcolor.Color) Color
```

### String Conversion

```go
// To hex string
color.RGBToHex(c Color) string

// To CSS color() function
color.ToColorFunction(c Color, space string) string
```

## Type Interfaces

### Color Interface

```go
type Color interface {
    RGBA() (r, g, b, a float64)
    Alpha() float64
    WithAlpha(alpha float64) Color
}
```

All color types implement this interface.

### SpaceColor Interface

```go
type SpaceColor interface {
    Color  // Embeds Color interface
    Space() Space
    Channels() []float64
    ConvertTo(space Space) SpaceColor
    ToRGBA() *RGBA
}
```

### Space Interface

```go
type Space interface {
    Name() string
    Channels() int
    ChannelNames() []string
    ToXYZ(channels []float64) (x, y, z float64)
    FromXYZ(x, y, z float64) []float64
}
```

## Constants and Enums

### Gamut Mapping

```go
const (
    GamutClip              GamutMapping = iota
    GamutPreserveChroma
    GamutPreserveLightness
    GamutProject
)
```

### Easing Functions

```go
const (
    EasingLinear        Easing = iota
    EasingEaseIn
    EasingEaseOut
    EasingEaseInOut
    EasingSmoothstep
)
```

### Hue Interpolation

```go
const (
    HueShorter         HueInterpolation = iota
    HueLonger
    HueIncreasing
    HueDecreasing
)
```

## Predefined Color Spaces

### RGB Spaces

```go
color.SRGBSpace         // Standard RGB (web)
color.SRGBLinearSpace   // Linear sRGB
color.DisplayP3Space    // Apple displays (26% wider)
color.DCIP3Space        // Digital cinema
color.A98RGBSpace       // Adobe RGB 1998
color.ProPhotoRGBSpace  // ProPhoto RGB (189% wider)
color.Rec2020Space      // UHDTV (73% wider)
color.Rec709Space       // HDTV
```

### LOG Spaces (Cinema Cameras)

```go
color.CLogSpace         // Canon C-Log
color.SLog3Space        // Sony S-Log3
color.VLogSpace         // Panasonic V-Log
color.ArriLogCSpace     // Arri LogC
color.RedLog3G10Space   // Red Log3G10
color.BMDFilmSpace      // Blackmagic Film
```

### Perceptual Spaces

```go
color.OKLCHSpace        // OKLCH (recommended)
color.OKLABSpace        // OKLAB
```

## Usage Patterns

### Basic Color Manipulation

```go
// Parse a color
blue, _ := color.ParseColor("#0000FF")

// Manipulate
lighter := color.Lighten(blue, 0.2)
vivid := color.Saturate(lighter, 0.3)

// Output
hex := color.RGBToHex(vivid)
```

### Wide-Gamut Workflow

```go
// Create in Display P3
p3Color := color.NewSpaceColor(color.DisplayP3Space,
    []float64{1.0, 0.2, 0.3}, 1.0)

// Manipulate (preserves wide gamut)
result := color.Lighten(p3Color, 0.2)

// Convert back
final, _ := color.ConvertToRGBSpace(result, "display-p3")
```

### Gradient Generation

```go
// Simple gradient
start := color.RGB(1, 0, 0)
end := color.RGB(0, 0, 1)
gradient := color.Gradient(start, end, 20)

// Multi-stop with easing
stops := []color.GradientStop{
    {Pos: 0.0, Color: color.RGB(1, 0, 0)},
    {Pos: 0.5, Color: color.RGB(0, 1, 0)},
    {Pos: 1.0, Color: color.RGB(0, 0, 1)},
}
gradient = color.GradientMultiStop(stops, 100)
```

### Color Comparison

```go
color1, _ := color.ParseColor("#FF0000")
color2, _ := color.ParseColor("#FE0000")

diff := color.DeltaE2000(color1, color2)
if diff < 1.0 {
    fmt.Println("Imperceptible to humans")
}
```

## Error Handling

Most functions that can fail return `(result, error)`:

```go
parsed, err := color.ParseColor("#invalid")
if err != nil {
    // Handle error
}

spaceColor, err := color.ConvertToRGBSpace(c, "unknown-space")
if err != nil {
    // Handle error
}
```

## Performance Tips

1. **Reuse color objects** - They're immutable, safe to share
2. **Batch conversions** - Convert once, use many times
3. **Cache gradients** - Generate once, reuse
4. **Use appropriate precision** - float64 is usually overkill for 8-bit displays

## Thread Safety

All operations are thread-safe:
- Color objects are immutable
- Color space registry uses sync.RWMutex
- No global mutable state

## See Also

- **[Color Space List](color-space-list.md)** - All supported color spaces
- **[Format Support](format-support.md)** - Parsing capabilities
- **[Examples on pkg.go.dev](https://pkg.go.dev/github.com/SCKelemen/color#pkg-examples)** - Executable examples
