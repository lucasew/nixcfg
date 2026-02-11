package palette

import "github.com/spf13/cobra"

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "palette",
		Short: "Color palette generation and management",
		Long:  "Generate base16/base24 color palettes from images using various extraction algorithms",
	}

	cmd.AddCommand(GetGenerateCommand())

	return cmd
}
