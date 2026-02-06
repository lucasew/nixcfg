package launch

import (
	"fmt"
	"os"
	"syscall"
	"workspaced/pkg/exec"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "launch",
		Short: "Launch applications",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "terminal",
		Short: "Launch the preferred terminal",
		RunE: func(c *cobra.Command, args []string) error {
			ctx := c.Context()
			terminals := []string{"kitty", "alacritty", "foot", "st", "xterm"}

			for _, term := range terminals {
				termPath, err := exec.Which(ctx, term)
				if err != nil {
					continue
				}

				// If we are in the daemon, we want to detach it
				if os.Getenv("WORKSPACED_DAEMON") == "1" {
					cmd := exec.RunCmd(ctx, termPath, args...)
					cmd.Stdout = nil
					cmd.Stderr = nil
					return cmd.Start()
				}

				// If local, just exec
				return syscall.Exec(termPath, append([]string{term}, args...), os.Environ())
			}

			return fmt.Errorf("no terminal found")
		},
	})

	return cmd
}
