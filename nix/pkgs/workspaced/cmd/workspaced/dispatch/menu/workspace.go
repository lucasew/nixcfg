package menu

import (
	"sort"
	"strings"
	"workspaced/pkg/exec"
	"workspaced/pkg/drivers/wm"

	"github.com/spf13/cobra"
	"workspaced/pkg/config"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		cmd := &cobra.Command{
			Use:   "workspace",
			Short: "Workspace switcher",
			RunE: func(c *cobra.Command, args []string) error {
				move, _ := c.Flags().GetBool("move")
				cfg, err := config.LoadConfig()
				if err != nil {
					return err
				}

				var keys []string
				for k := range cfg.Workspaces {
					keys = append(keys, k)
				}
				sort.Strings(keys)

				cmd := exec.RunCmd(c.Context(), "dmenu")
				cmd.Stdin = strings.NewReader(strings.Join(keys, "\n"))

				out, err := cmd.Output()
				if err != nil {
					return err
				}

				selected := strings.TrimSpace(string(out))
				if selected == "" {
					return nil
				}

				workspaceNum, ok := cfg.Workspaces[selected]
				if !ok {
					return nil
				}

				return wm.SwitchToWorkspace(c.Context(), workspaceNum, move)
			},
		}
		cmd.Flags().Bool("move", false, "Move container to workspace")
		parent.AddCommand(cmd)
	})
}
