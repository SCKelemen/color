# Color Space Architecture

## Design Goals

1. **Lossless conversions within a space**: Converting between color representations in the same space should be perfectly reversible (e.g., RGB ↔ HSL in sRGB, LAB ↔ LCH)
2. **Lossless operations within a space**: Operations like `Lighten()`, `Darken()`, `Saturate()` should work in the native color space without intermediate conversions
3. **Explicit conversion only**: Data loss only occurs when explicitly converting to a different color space or requesting a specific output format (e.g., `ToRGBA()`)

## Architecture

### Core Concepts

1. **Color Space**: A mathematical definition of how colors are represented
   - Primaries (for RGB spaces)
   - White point
   - Transfer function (gamma/linear)
   - Conversion matrices to/from XYZ

2. **Color Value**: A color with its space information
   - Stores channel values in the native space
   - Carries space metadata
   - Operations work in native space

3. **Reference Space**: XYZ (or OKLAB for perceptually uniform operations)
   - All inter-space conversions go through reference space
   - Provides lossless conversion hub

### Type System

```go
// Space defines a color space with conversion to/from XYZ
type Space interface {
    Name() string
    // Convert from this space to XYZ (linear, D65)
    ToXYZ(channels []float64) (x, y, z float64)
    // Convert from XYZ to this space
    FromXYZ(x, y, z float64) []float64
    // Number of channels (3 for RGB, 4 for CMYK, etc.)
    Channels() int
    // Get channel names
    ChannelNames() []string
}

// Color represents a color in a specific color space
type Color interface {
    // Space returns the color space this color is in
    Space() Space
    
    // Channels returns the color values in the native space
    Channels() []float64
    
    // Alpha returns the alpha channel [0, 1]
    Alpha() float64
    
    // WithAlpha returns a new color with the specified alpha
    WithAlpha(alpha float64) Color
    
    // Convert to a different space (explicit conversion)
    ConvertTo(space Space) Color
    
    // ToRGBA converts to sRGB RGBA (explicit conversion, may lose data)
    ToRGBA() *RGBA
}
```

### Implementation Strategy

1. **Native Space Operations**
   - `Lighten()`, `Darken()`, `Saturate()`, etc. work in the color's native space
   - For perceptually uniform spaces (OKLCH, OKLAB), operations are mathematically correct
   - For RGB spaces, operations convert to OKLCH, operate, convert back

2. **Lossless Conversions**
   - RGB ↔ HSL/HSV: Perfectly reversible in same RGB space
   - LAB ↔ LCH: Perfectly reversible
   - OKLAB ↔ OKLCH: Perfectly reversible
   - All conversions preserve full precision

3. **Inter-Space Conversions**
   - All go through XYZ (or OKLAB for perceptually uniform)
   - Explicit conversion: `color.ConvertTo(targetSpace)`
   - This is where gamut clipping may occur

4. **Backward Compatibility**
   - Keep existing `Color` interface for compatibility
   - New `SpaceColor` interface for explicit space-aware colors
   - Adapter functions to bridge old and new APIs

## Example Usage

```go
// Create a color in a specific space
p3Color := color.NewSpaceColor(
    color.DisplayP3Space,
    []float64{1.0, 0.5, 0.0}, // R, G, B in Display P3
    1.0, // alpha
)

// Operations work in native space (Display P3)
lightened := p3Color.Lighten(0.2) // Still in Display P3

// Explicit conversion to different space
srgbColor := lightened.ConvertTo(color.SRGBSpace) // Now in sRGB

// Convert to RGBA (explicit, may clip gamut)
rgba := srgbColor.ToRGBA() // sRGB RGBA for display

// Lossless conversion within same space family
hsl := srgbColor.ToHSL() // RGB → HSL in sRGB (lossless)
rgbBack := hsl.ToRGB()   // HSL → RGB in sRGB (lossless, perfect round-trip)
```

## Migration Path

1. **Phase 1**: Implement `Space` interface and basic spaces
2. **Phase 2**: Implement `SpaceColor` type
3. **Phase 3**: Update operations to work in native spaces
4. **Phase 4**: Add adapters for backward compatibility
5. **Phase 5**: Deprecate old API (optional)

## Benefits

- **No accidental data loss**: Operations don't convert through RGBA unless explicitly requested
- **Gamut preservation**: Wide-gamut colors stay in their space until explicitly converted
- **Mathematical correctness**: Operations use appropriate spaces (OKLCH for perceptual, native for others)
- **Explicit conversions**: Users know when they're converting and may lose data
- **Future-proof**: Easy to add new color spaces without breaking existing code

