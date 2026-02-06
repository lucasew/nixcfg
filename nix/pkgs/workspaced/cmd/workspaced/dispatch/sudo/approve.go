package sudo

import (
	"os"
	"os/exec"
	"workspaced/pkg/drivers/sudo"
	"workspaced/pkg/logging"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		cmd := &cobra.Command{
			Use:   "approve <slug>",
			Short: "Approve and execute a pending command",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				ctx := cmd.Context()
				logger := logging.GetLogger(ctx)
				slug := args[0]
				sc, err := sudo.Get(slug)
				if err != nil {
					return err
				}

				logger.Info("approving command", "command", sc.Command, "args", sc.Args, "slug", slug)

				// Always remove after attempting to run
				defer func() { _ = sudo.Remove(slug) }()

				ec := exec.Command("sudo", append([]string{"-E", sc.Command}, sc.Args...)...)
				ec.Stdout = os.Stdout
				ec.Stderr = os.Stderr
				ec.Stdin = os.Stdin
				ec.Dir = sc.Cwd
				ec.Env = sc.Env

				return ec.Run()
			},
		}
		parent.AddCommand(cmd)
	})
}
