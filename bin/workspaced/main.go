package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"workspaced/pkg/cmd"
)

var rootCmd = &cobra.Command{
	Use:   "workspaced",
	Short: "Workspace daemon and client",
}

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Run the workspaced daemon",
	Run: func(c *cobra.Command, args []string) {
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

var modnCmd = &cobra.Command{
	Use:   "modn",
	Short: "Rotate workspaces across outputs",
	Run: func(c *cobra.Command, args []string) {
		runOrRoute("modn", args, cmd.RunModn)
	},
}

var mediaCmd = &cobra.Command{
	Use:   "media",
	Short: "Control media playback",
	Run: func(c *cobra.Command, args []string) {
		runOrRoute("media", args, func() (string, error) {
			return cmd.RunMedia(args)
		})
	},
}

var rofiCmd = &cobra.Command{
	Use:   "rofi",
	Short: "Rofi workspace switcher",
	Run: func(c *cobra.Command, args []string) {
		runOrRoute("rofi", args, func() (string, error) {
			return cmd.RunRofi(args, os.Environ())
		})
	},
}
