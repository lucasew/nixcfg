package audio

import (
	"workspaced/pkg/registry"

	"github.com/spf13/cobra"
)

var Registry registry.CommandRegistry

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "audio",
		Short: "Control audio volume",
	}
	return Registry.GetCommand(cmd)
}
