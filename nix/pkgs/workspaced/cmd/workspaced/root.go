package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"workspaced/cmd/workspaced/daemon"
	"workspaced/cmd/workspaced/dispatch"
	"workspaced/cmd/workspaced/dispatch/apply"
	"workspaced/cmd/workspaced/dispatch/history"
	"workspaced/cmd/workspaced/dispatch/sync"
	"workspaced/cmd/workspaced/is"
	"workspaced/cmd/workspaced/launch"
	"workspaced/cmd/workspaced/svc"
	dapi "workspaced/pkg/api"
	"workspaced/pkg/shellgen"

	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workspaced",
		Short: "Workspace daemon and client",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Ensure logs go to stderr so stdout stays clean for piping/capturing
			slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, nil)))
			if os.Getenv("WORKSPACED_DEBUG") != "" {
				slog.SetLogLoggerLevel(slog.LevelDebug)
			}

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

	// Set root command for shell generators (e.g., completion)
	shellgen.SetRootCommand(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		if errors.Is(err, dapi.ErrCanceled) {
			os.Exit(0)
		}
		slog.Error("error", "err", err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
