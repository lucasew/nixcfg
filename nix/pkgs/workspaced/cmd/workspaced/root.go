package main

import (
	"log/slog"
	"os"
	"workspaced/cmd/workspaced/daemon"
	"workspaced/cmd/workspaced/dispatch"
	"workspaced/cmd/workspaced/history"
	"workspaced/cmd/workspaced/input"
	"workspaced/cmd/workspaced/open"
	"workspaced/cmd/workspaced/state"
	"workspaced/cmd/workspaced/system"
	"workspaced/pkg/driver/media"
	"workspaced/pkg/shellgen"
	"workspaced/pkg/version"

	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:     "workspaced",
		Short:   "workspaced - declarative user environment manager",
		Version: version.BuildID,
	}

	// Main Command Groups
	cmd.AddCommand(input.NewCommand())
	cmd.AddCommand(open.NewCommand())
	cmd.AddCommand(system.NewCommand())
	cmd.AddCommand(state.NewCommand())
	cmd.AddCommand(history.NewCommand())

	// Top-level aliases for daily ergonomic
	stateCmd := state.NewCommand()
	for _, c := range stateCmd.Commands() {
		if c.Name() == "apply" || c.Name() == "plan" || c.Name() == "sync" || c.Name() == "doctor" {
			cmd.AddCommand(c)
		}
	}
	historyCmd := history.NewCommand()
	for _, c := range historyCmd.Commands() {
		if c.Name() == "search" {
			cmd.AddCommand(c)
		}
	}

	// Daemon and Internal
	cmd.AddCommand(daemon.Command)
	cmd.AddCommand(dispatch.NewCommand()) // Keep hidden or for internal use

	// Media shortcuts (still very common)
	mediaCmd := &cobra.Command{
		Use:   "media",
		Short: "Media player control",
	}
	mediaCmd.AddCommand(&cobra.Command{
		Use:   "next",
		Short: "Next track",
		RunE: func(c *cobra.Command, args []string) error {
			return media.RunAction(c.Context(), "next")
		},
	})
	mediaCmd.AddCommand(&cobra.Command{
		Use:   "previous",
		Short: "Previous track",
		RunE: func(c *cobra.Command, args []string) error {
			return media.RunAction(c.Context(), "previous")
		},
	})
	mediaCmd.AddCommand(&cobra.Command{
		Use:   "play-pause",
		Short: "Toggle play/pause",
		RunE: func(c *cobra.Command, args []string) error {
			return media.RunAction(c.Context(), "play-pause")
		},
	})
	cmd.AddCommand(mediaCmd)

	// Set root command for shell completion generation
	shellgen.SetRootCommand(cmd)

	if err := cmd.Execute(); err != nil {
		slog.Error("error", "err", err)
		os.Exit(1)
	}
}
