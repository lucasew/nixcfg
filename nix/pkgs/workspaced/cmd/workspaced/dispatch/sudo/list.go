package sudo

import (
	"sort"
	"time"
	"workspaced/pkg/sudo"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		cmd := &cobra.Command{
			Use:   "list",
			Short: "List pending commands",
			RunE: func(cmd *cobra.Command, args []string) error {
				cmds, err := sudo.List(cmd.Context())
				if err != nil {
					return err
				}

				sort.Slice(cmds, func(i, j int) bool {
					return cmds[i].Timestamp < cmds[j].Timestamp
				})

				if len(cmds) == 0 {
					cmd.Println("No pending commands.")
					return nil
				}

				cmd.Printf("%-15s %-10s %s\n", "SLUG", "TIME", "COMMAND")
				for _, c := range cmds {
					t := time.Unix(c.Timestamp, 0).Format("15:04:05")
					fullCmd := append([]string{c.Command}, c.Args...)
					cmd.Printf("%-15s %-10s %v\n", c.Slug, t, fullCmd)
				}
				return nil
			},
		}
		parent.AddCommand(cmd)
	})
}
