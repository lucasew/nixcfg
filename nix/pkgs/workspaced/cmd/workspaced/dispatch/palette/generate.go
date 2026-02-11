package palette

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"

	"workspaced/pkg/config"
	"workspaced/pkg/palette"
	"workspaced/pkg/palette/api"
)

func GetGenerateCommand() *cobra.Command {
	var (
		driverName string
		polarity   string
		colorCount int
		outputJSON bool
	)

	cmd := &cobra.Command{
		Use:   "generate <image>",
		Short: "Generate color palette from an image",
		Long: `Generate a base16 or base24 color palette from an image using evolutionary algorithms.

The genetic driver uses an evolutionary algorithm to find colors that:
- Match the colors in the input image
- Create a harmonious color scheme
- Follow base16/base24 specifications for terminal and editor themes

Examples:
  # Generate dark theme from wallpaper
  workspaced palette generate ~/wallpaper.jpg --polarity dark

  # Generate light theme as JSON
  workspaced palette generate image.png --polarity light --json

  # Generate base24 palette (24 colors instead of 16)
  workspaced palette generate image.png --colors 24`,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			ctx := c.Context()
			imagePath := args[0]

			// Parse polarity flag
			pol, err := parsePolarityFlag(polarity)
			if err != nil {
				return err
			}

			opts := api.Options{
				Polarity:   pol,
				ColorCount: colorCount,
				MaxSamples: 10000,
			}

			// Extract palette
			pal, err := palette.ExtractFromFile(ctx, imagePath, driverName, opts)
			if err != nil {
				return fmt.Errorf("failed to extract palette: %w", err)
			}

			// Output format
			if outputJSON {
				encoder := json.NewEncoder(os.Stdout)
				encoder.SetIndent("", "  ")
				return encoder.Encode(pal)
			}

			return printPaletteTOML(pal)
		},
	}

	cmd.Flags().StringVar(&driverName, "driver", "genetic", "Extraction algorithm (genetic)")
	cmd.Flags().StringVar(&polarity, "polarity", "any", "Theme preference: dark, light, or any")
	cmd.Flags().IntVar(&colorCount, "colors", 16, "Number of colors (16 for base16, 24 for base24)")
	cmd.Flags().BoolVar(&outputJSON, "json", false, "Output as JSON instead of TOML")

	return cmd
}

// parsePolarityFlag converts string flag to Polarity enum
func parsePolarityFlag(s string) (api.Polarity, error) {
	switch strings.ToLower(s) {
	case "any":
		return api.PolarityAny, nil
	case "dark":
		return api.PolarityDark, nil
	case "light":
		return api.PolarityLight, nil
	default:
		return api.PolarityAny, fmt.Errorf("invalid polarity: %s (must be 'dark', 'light', or 'any')", s)
	}
}

// printPaletteTOML outputs palette in TOML format
func printPaletteTOML(pal *config.PaletteConfig) error {
	// Create a map for TOML encoding
	paletteMap := map[string]interface{}{
		"palette": pal,
	}

	encoder := toml.NewEncoder(os.Stdout)
	return encoder.Encode(paletteMap)
}
