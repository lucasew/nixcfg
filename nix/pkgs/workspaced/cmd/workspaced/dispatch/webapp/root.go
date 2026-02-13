package webapp

import (
	"workspaced/pkg/registry"

	"github.com/spf13/cobra"
)

var Registry registry.CommandRegistry

type WebappConfig struct {
	URL         string   `json:"url"`
	Profile     string   `json:"profile"`
	DesktopName string   `json:"desktop_name"`
	Icon        string   `json:"icon"`
	ExtraFlags  []string `json:"extra_flags"`
}

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "webapp",
		Short: "Manage webapps",
	}
	return Registry.GetCommand(cmd)
}
