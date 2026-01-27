package dispatch

import (
	"os/exec"

	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:              "dispatch",
	Short:            "Dispatch workspace commands",
	TraverseChildren: true,
}

func runCmd(c *cobra.Command, name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	ctx := c.Context()
	if ctx != nil {
		if env, ok := ctx.Value(EnvKey).([]string); ok {
			cmd.Env = env
		}
	}
	return cmd
}

func GetFullCommandPath(c *cobra.Command) []string {
	var path []string
	curr := c
	for curr != nil && curr.Name() != "dispatch" && curr.Name() != "workspaced" {
		path = append([]string{curr.Name()}, path...)
		curr = curr.Parent()
	}
	return path
}

func FindCommand(name string, args []string) (*cobra.Command, []string, error) {
	return Command.Find(append([]string{name}, args...))
}
