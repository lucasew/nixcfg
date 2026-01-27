package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"workspaced/cmd/workspaced/dispatch"
)

var rootCmd = &cobra.Command{
	Use:   "workspaced",
	Short: "Workspace daemon and client",
}

func Execute() {
	rootCmd.AddCommand(dispatch.DaemonCmd)
	rootCmd.AddCommand(dispatch.ModnCmd)
	rootCmd.AddCommand(dispatch.MediaCmd)
	rootCmd.AddCommand(dispatch.RofiCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func GetRootCmd() *cobra.Command {
	// Re-add commands to ensure fresh state if needed, or just return existing
	// For daemon execution we might need a fresh clone or careful handling of args
	cmd := &cobra.Command{
		Use:   "workspaced",
		Short: "Workspace daemon and client",
	}
	cmd.AddCommand(dispatch.DaemonCmd)
	cmd.AddCommand(dispatch.ModnCmd)
	cmd.AddCommand(dispatch.MediaCmd)
	cmd.AddCommand(dispatch.RofiCmd)
	return cmd
}
