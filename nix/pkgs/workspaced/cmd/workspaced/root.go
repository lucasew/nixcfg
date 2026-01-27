package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"workspaced/cmd/workspaced/dispatch"
	"workspaced/cmd/workspaced/daemon"
)

var RootCmd = &cobra.Command{
	Use:   "workspaced",
	Short: "Workspace daemon and client",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.AddCommand(dispatch.Command)
	RootCmd.AddCommand(daemon.Command)
}

func main() {
	Execute()
}
