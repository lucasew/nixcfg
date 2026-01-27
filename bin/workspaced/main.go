package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "workspaced",
	Short: "Workspace daemon and client",
}

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Run the workspaced daemon",
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunDaemon(); err != nil {
			fmt.Fprintf(os.Stderr, "Daemon error: %v\n", err)
			os.Exit(1)
		}
	},
}

func main() {
	rootCmd.AddCommand(daemonCmd)
	rootCmd.AddCommand(modnCmd)
	rootCmd.AddCommand(mediaCmd)
	rootCmd.AddCommand(rofiCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// Helper to run logic either remotely or locally
func runOrRoute(cmdName string, args []string, localFunc func() (string, error)) {
	output, ranRemote, err := TryRemote(cmdName, args)
	if ranRemote {
		if output != "" {
			fmt.Print(output)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Remote error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Local fallback
	out, err := localFunc()
	if out != "" {
		fmt.Print(out)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Local error: %v\n", err)
		os.Exit(1)
	}
}
