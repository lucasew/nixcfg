package brightness

import (
	"workspaced/pkg/common/registry"

	"github.com/spf13/cobra"
)

var Registry registry.CommandRegistry

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "brightness",
		Short: "Control screen brightness",
	}
	return Registry.GetCommand(cmd)
}
