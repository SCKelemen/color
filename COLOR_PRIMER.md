# Color Primer

## What is color?
- **Physical**: Light with a spectrum of wavelengths.
- **Perceptual**: Human vision reduces spectra to three responses (cones), so colors are modeled as 3‑component vectors.

## What is a color space?
- A **coordinate system** for color.
- Defines primaries (or axes), white point, and a transfer/gamma function.
- Examples: sRGB, Display P3, Rec.2020, XYZ, Lab, OKLab.

## What is a color gamut?
- The **range of colors** a space or device can represent.
- Visualized as a volume in a reference space (e.g., XYZ or Lab/OKLab).
- Converting from a larger to a smaller gamut can cause clipping or compression.

## What is a reference space?
- A common hub for conversions between spaces.
- **In this library:** single reference = **CIE XYZ (D65)**. All inter-space conversions go through XYZ (with chromatic adaptation if needed).
- Working/perceptual space for operations: **OKLCH** (for lighten/darken/saturate), but conversions still route through XYZ.

## What are perceptual color spaces?
- Spaces designed for more uniform perceptual distances.
- **OKLab / OKLCH** (modern, recommended), **Lab / LCH** (classic).
- Benefits: smoother gradients, better “evenness” of lightness/saturation adjustments.

## Color spaces and gamuts we support
- **RGB spaces (with explicit primaries/transfer):**
  - sRGB, sRGB-linear
  - Display P3
  - Adobe RGB (a98-rgb)
  - ProPhoto RGB
  - Rec. 2020
- **CIE / perceptual:**
  - XYZ (D65 reference)
  - Lab / LCH
  - OKLab / OKLCH
- **Legacy parameterizations (sRGB-based):**
  - HSL, HSV

## Where conversions can be lossy
- **Gamut reduction:** converting wide-gamut (P3/2020/ProPhoto) to sRGB (e.g., `ToRGBA()`).
- **Explicit ConvertTo to a smaller gamut** (out-of-gamut values clip/compress).
- **Metadata loss** if using the legacy `Color` interface that assumes sRGB.
- **Quantization** if exporting to 8-bit or YCbCr video ranges (not enabled by default here).

## Gradients: why space matters
- Same stops, different results:
  - **RGB:** non-uniform, can “dip” or “bunch.”
  - **HSL:** hue-based; lightness may feel uneven.
  - **Lab/LCH:** more uniform than RGB/HSL.
  - **OKLab/OKLCH:** most uniform (recommended).
- Recommendation: Generate gradients in **OKLCH** (or OKLab) for perceptual smoothness.

### Example: gradients in different spaces
```go
red := color.RGB(1, 0, 0)
blue := color.RGB(0, 0, 1)

// Perceptually uniform (recommended)
gOKLCH := color.GradientInSpace(red, blue, 10, color.GradientOKLCH)

// HSL (hue-based)
gHSL := color.GradientInSpace(red, blue, 10, color.GradientHSL)

// RGB (fast, not uniform)
gRGB := color.GradientInSpace(red, blue, 10, color.GradientRGB)
```

## Planned visuals (add later)
- Gamut slices (2D) for sRGB vs P3 vs Rec.2020 in a reference space (Lab/OKLab).
- Gradient strips comparing RGB, HSL, Lab, OKLab/OKLCH for the same stops.

## Key takeaways
- Single reference hub: **XYZ (D65)**.
- Perceptual working space for operations: **OKLCH** (conversions still via XYZ).
- Explicit conversions only; data loss happens mainly on gamut reduction or legacy sRGB-only paths.

