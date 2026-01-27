package common

import (
	"os"
	"os/exec"
	"strings"
	"workspaced/cmd/workspaced/dispatch/types"

	"github.com/spf13/cobra"
)

func RunCmd(c *cobra.Command, name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	ctx := c.Context()
	if ctx != nil {
		if env, ok := ctx.Value(types.EnvKey).([]string); ok {
			cmd.Env = env
		}
	}
	return cmd
}

func GetRPC(env []string) string {
	for _, e := range env {
		if strings.HasPrefix(e, "WAYLAND_DISPLAY=") {
			return "swaymsg"
		}
	}
	if os.Getenv("WAYLAND_DISPLAY") != "" {
		return "swaymsg"
	}
	return "i3-msg"
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
