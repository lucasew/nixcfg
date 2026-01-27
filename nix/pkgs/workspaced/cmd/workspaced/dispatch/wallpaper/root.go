package wallpaper

import (
	"workspaced/pkg/common/registry"

	"github.com/spf13/cobra"
)

var Registry registry.CommandRegistry

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wallpaper",
		Short: "Wallpaper management",
	}
	return Registry.GetCommand(cmd)
}
