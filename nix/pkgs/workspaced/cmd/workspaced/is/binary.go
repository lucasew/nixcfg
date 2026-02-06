package is

import (
	"fmt"
	"workspaced/pkg/exec"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "binary <name>",
			Short: "Check if binary is available",
			Args:  cobra.ExactArgs(1),
			RunE: func(c *cobra.Command, args []string) error {
				if !exec.IsBinaryAvailable(c.Context(), args[0]) {
					return fmt.Errorf("binary %s not available", args[0])
				}
				return nil
			},
		})
	})
}
