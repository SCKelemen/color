# Release v1.0.2

## Summary

This patch release improves parser correctness, color-space fidelity, and documentation consistency.

## Highlights

- Parser hardening:
  - Consistent argument-count validation across color functions.
  - Consistent alpha-range validation (`0..1` / `0..100%`) across supported formats.
  - RGB numeric channel behavior aligned with CSS byte semantics (`0..255`).
- Better color-space preservation:
  - `ParseColor("color(<wide-gamut> ... )")` preserves source space as `SpaceColor`.
  - `ConvertFromRGBSpace(...)` now preserves source-space channels instead of immediately clipping through sRGB.
  - `ToXYZ(...)` now uses native `SpaceColor` channels when available.
- CSS compatibility improvements:
  - Hue angle units support: `deg`, `rad`, `grad`, `turn`.
  - Hue wrapping support for over/under-range values (for HSL/HSV/HWB/LCH/OKLCH parsing paths).
  - `color(xyz-d50 ...)` now performs D50→D65 chromatic adaptation.
- Named color coverage:
  - Expanded to full CSS named-color set, including `rebeccapurple`, plus `transparent`.
- Gradients:
  - Fixed multistop edge cases (`steps <= 0`, `steps == 1`) and aligned behavior across variants.

## Quality

- `go test ./...` passes
- `go build ./...` passes
- `go vet ./...` passes

## Notes

- This release adds a dependency on `github.com/SCKelemen/units` to standardize angle-unit handling semantics across libraries.
