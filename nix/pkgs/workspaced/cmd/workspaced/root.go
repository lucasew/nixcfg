package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"workspaced/cmd/workspaced/daemon"
	"workspaced/cmd/workspaced/dispatch"
	"workspaced/cmd/workspaced/dispatch/apply"
	"workspaced/cmd/workspaced/dispatch/config"
	"workspaced/cmd/workspaced/dispatch/history"
	"workspaced/cmd/workspaced/dispatch/sync"
	"workspaced/cmd/workspaced/is"
	"workspaced/cmd/workspaced/launch"
	"workspaced/cmd/workspaced/svc"
	"workspaced/pkg/prelude"

	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workspaced",
		Short: "Workspace daemon and client",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Ensure logs go to stderr so stdout stays clean for piping/capturing
			slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, nil)))

			path := os.Getenv("PATH")
			systemPath := "/run/current-system/sw/bin"
			if !strings.Contains(path, systemPath) {
				_ = os.Setenv("PATH", fmt.Sprintf("%s:%s", systemPath, path))
			}
		},
	}
	cmd.AddCommand(dispatch.NewCommand())
	cmd.AddCommand(daemon.Command)
	cmd.AddCommand(svc.NewCommand())
	cmd.AddCommand(is.GetCommand())
	cmd.AddCommand(launch.NewCommand())

	// Top-level aliases for common commands
	cmd.AddCommand(apply.GetCommand())
	cmd.AddCommand(sync.GetCommand())
	cmd.AddCommand(config.GetColorsCommand())

	// History search alias
	searchCmd := history.GetCommand()
	// Find the search subcommand and add it directly
	for _, sub := range searchCmd.Commands() {
		if sub.Name() == "search" {
			cmd.AddCommand(sub)
			break
		}
	}

	return cmd
}

func Execute() {
	rootCmd := NewRootCommand()

	// Set root command for prelude generators (e.g., completion)
	prelude.SetRootCommand(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
