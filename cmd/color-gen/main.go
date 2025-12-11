package main

import (
	"context"
	"fmt"
	"os"

	"github.com/SCKelemen/clix"
	col "github.com/SCKelemen/color"
)

func main() {
	app := clix.NewApp("color-gen")

	app.Root = clix.NewGroup("color-gen", "Generate color space visualizations and diagrams",
		func() *clix.Command {
			cmd := clix.NewCommand("gradients", clix.WithCommandShort("Generate gradient comparison images in different color spaces"))
			cmd.Run = func(ctx *clix.Context) error {
				var startColor, endColor string
				startArg := ctx.Arg(0)
				if startArg != "" {
					startColor = startArg
				} else {
					startColor = "rgb(255, 0, 0)" // default red
				}
				endArg := ctx.Arg(1)
				if endArg != "" {
					endColor = endArg
				} else {
					endColor = "rgb(0, 0, 255)" // default blue
				}

				start, err := col.ParseColor(startColor)
				if err != nil {
					return fmt.Errorf("failed to parse start color %q: %w", startColor, err)
				}

				end, err := col.ParseColor(endColor)
				if err != nil {
					return fmt.Errorf("failed to parse end color %q: %w", endColor, err)
				}

				return generateGradients(start, end)
			}
			return cmd
		}(),

		func() *clix.Command {
			cmd := clix.NewCommand("stops", clix.WithCommandShort("Generate color stop images"))
			cmd.Run = func(ctx *clix.Context) error {
				var startColor, endColor string
				start := ctx.Arg(0)
				if start != "" {
					startColor = start
				} else {
					startColor = "rgb(255, 0, 0)" // default red
				}
				end := ctx.Arg(1)
				if end != "" {
					endColor = end
				} else {
					endColor = "rgb(0, 0, 255)" // default blue
				}

				startCol, err := col.ParseColor(startColor)
				if err != nil {
					return fmt.Errorf("failed to parse start color %q: %w", startColor, err)
				}

				endCol, err := col.ParseColor(endColor)
				if err != nil {
					return fmt.Errorf("failed to parse end color %q: %w", endColor, err)
				}

				return generateStops(startCol, endCol)
			}
			return cmd
		}(),

		func() *clix.Command {
			cmd := clix.NewCommand("gamuts", clix.WithCommandShort("Generate gamut volume visualizations for different RGB color spaces"))
			cmd.Run = func(ctx *clix.Context) error {
				return generateGamuts()
			}
			return cmd
		}(),

		func() *clix.Command {
			cmd := clix.NewCommand("xyz-gamuts", clix.WithCommandShort("Generate comparison of all RGB gamuts in XYZ color space"))
			cmd.Run = func(ctx *clix.Context) error {
				return generateXYZGamuts()
			}
			return cmd
		}(),

		func() *clix.Command {
			cmd := clix.NewCommand("chromaticity", clix.WithCommandShort("Generate CIE xy chromaticity diagrams for different RGB color spaces"))
			cmd.Run = func(ctx *clix.Context) error {
				return generateChromaticity()
			}
			return cmd
		}(),

		func() *clix.Command {
			cmd := clix.NewCommand("models", clix.WithCommandShort("Generate color model diagrams (RGB cube, HSL cylinder, LAB, OKLCH)"))
			cmd.Run = func(ctx *clix.Context) error {
				return generateModels()
			}
			return cmd
		}(),

		func() *clix.Command {
			cmd := clix.NewCommand("animations", clix.WithCommandShort("Generate animation frames for rotating color model visualizations"))
			cmd.Run = func(ctx *clix.Context) error {
				return generateAnimations()
			}
			return cmd
		}(),

		func() *clix.Command {
			cmd := clix.NewCommand("gifs", clix.WithCommandShort("Convert animation frames to animated GIFs"))
			cmd.Run = func(ctx *clix.Context) error {
				return generateGIFs()
			}
			return cmd
		}(),

		func() *clix.Command {
			cmd := clix.NewCommand("all", clix.WithCommandShort("Generate all visualizations (gradients, stops, gamuts, chromaticity, models, animations, gifs)"))
			cmd.Run = func(ctx *clix.Context) error {
				var startColor, endColor string
				start := ctx.Arg(0)
				if start != "" {
					startColor = start
				} else {
					startColor = "rgb(255, 0, 0)" // default red
				}
				end := ctx.Arg(1)
				if end != "" {
					endColor = end
				} else {
					endColor = "rgb(0, 0, 255)" // default blue
				}

				startCol, err := col.ParseColor(startColor)
				if err != nil {
					return fmt.Errorf("failed to parse start color %q: %w", startColor, err)
				}

				endCol, err := col.ParseColor(endColor)
				if err != nil {
					return fmt.Errorf("failed to parse end color %q: %w", endColor, err)
				}

				// Generate everything
				if err := generateGradients(startCol, endCol); err != nil {
					return fmt.Errorf("gradients: %w", err)
				}
				if err := generateStops(startCol, endCol); err != nil {
					return fmt.Errorf("stops: %w", err)
				}
				if err := generateGamuts(); err != nil {
					return fmt.Errorf("gamuts: %w", err)
				}
				if err := generateChromaticity(); err != nil {
					return fmt.Errorf("chromaticity: %w", err)
				}
				if err := generateModels(); err != nil {
					return fmt.Errorf("models: %w", err)
				}
				if err := generateAnimations(); err != nil {
					return fmt.Errorf("animations: %w", err)
				}
				if err := generateGIFs(); err != nil {
					return fmt.Errorf("gifs: %w", err)
				}

				fmt.Println("✓ All visualizations generated successfully!")
				return nil
			}
			return cmd
		}(),
	)

	if err := app.Run(context.Background(), nil); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
