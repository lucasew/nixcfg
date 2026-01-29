package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"workspaced/cmd/workspaced/daemon"
	"workspaced/cmd/workspaced/dispatch"
	"workspaced/cmd/workspaced/svc"
)

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workspaced",
		Short: "Workspace daemon and client",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			path := os.Getenv("PATH")
			systemPath := "/run/current-system/sw/bin"
			if !strings.Contains(path, systemPath) {
				_ = os.Setenv("PATH", fmt.Sprintf("%s:%s", systemPath, path))
			}
		},
	}
	cmd.AddCommand(dispatch.NewCommand())
	cmd.AddCommand(daemon.Command) // daemon is a global command instance, but its Run doesn't use shared state in a racey way
	cmd.AddCommand(svc.NewCommand())
	return cmd
}

func Execute() {
	if err := NewRootCommand().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
