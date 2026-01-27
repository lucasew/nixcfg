package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"workspaced/cmd/workspaced/daemon"
	"workspaced/cmd/workspaced/dispatch"
)

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workspaced",
		Short: "Workspace daemon and client",
	}
	cmd.AddCommand(dispatch.NewCommand())
	cmd.AddCommand(daemon.Command) // daemon is a global command instance, but its Run doesn't use shared state in a racey way
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
