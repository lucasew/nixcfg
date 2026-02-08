package menu

import (
	"workspaced/pkg/exec"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "launch",
			Short: "Application launcher",
			RunE: func(c *cobra.Command, args []string) error {
				return exec.RunCmd(c.Context(), "rofi", "-show", "combi", "-combi-modi", "drun", "-show-icons").Run()
			},
		})
	})
}
