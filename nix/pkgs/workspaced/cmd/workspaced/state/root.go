package state

import (
	"workspaced/cmd/workspaced/dispatch/apply"
	"workspaced/cmd/workspaced/dispatch/config"
	"workspaced/cmd/workspaced/dispatch/doctor"
	"workspaced/cmd/workspaced/dispatch/plan"
	"workspaced/cmd/workspaced/dispatch/sync"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "state",
		Short: "Dotfiles and system state management",
	}

	cmd.AddCommand(apply.GetCommand())
	cmd.AddCommand(plan.GetCommand())
	cmd.AddCommand(sync.GetCommand())
	cmd.AddCommand(doctor.Command)
	cmd.AddCommand(config.GetCommand())

	return cmd
}
