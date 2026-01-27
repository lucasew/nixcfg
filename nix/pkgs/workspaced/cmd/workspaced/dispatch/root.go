package dispatch

import (
	"os/exec"

	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "dispatch",
	Short: "Dispatch workspace commands",
}

func runCmd(c *cobra.Command, name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	ctx := c.Context()
	if ctx != nil {
		if env, ok := ctx.Value("env").([]string); ok {
			cmd.Env = env
		}
	}
	return cmd
}
