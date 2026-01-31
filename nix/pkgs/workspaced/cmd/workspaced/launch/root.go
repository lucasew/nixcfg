package launch

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"workspaced/pkg/common"

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
				if common.IsBinaryAvailable(ctx, term) {
					path, err := exec.LookPath(term)
					if err != nil {
						continue
					}

					// If local, just exec
					return syscall.Exec(path, append([]string{term}, args...), os.Environ())
				}
			}

			return fmt.Errorf("no terminal found")
		},
	})

	return cmd
}
