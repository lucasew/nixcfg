package screen

import (
	"fmt"
	"os"
	"strings"
	"workspaced/cmd/workspaced/dispatch/common"
	"workspaced/cmd/workspaced/dispatch/types"

	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "screen",
	Short: "Screen and power management",
}

var lockCmd = &cobra.Command{
	Use:   "lock",
	Short: "Lock the screen and turn it off",
	RunE: func(c *cobra.Command, args []string) error {
		if err := common.RunCmd(c, "loginctl", "lock-session").Run(); err != nil {
			return fmt.Errorf("failed to lock session: %w", err)
		}
		return setDPMS(c, false)
	},
}

var offCmd = &cobra.Command{
	Use:   "off",
	Short: "Turn off the screen (DPMS)",
	RunE: func(c *cobra.Command, args []string) error {
		return setDPMS(c, false)
	},
}

var onCmd = &cobra.Command{
	Use:   "on",
	Short: "Turn on the screen (DPMS)",
	RunE: func(c *cobra.Command, args []string) error {
		return setDPMS(c, true)
	},
}

var toggleCmd = &cobra.Command{
	Use:   "toggle",
	Short: "Toggle screen state (DPMS)",
	RunE: func(c *cobra.Command, args []string) error {
		isOn, err := isDPMSOn(c)
		if err != nil {
			return err
		}
		return setDPMS(c, !isOn)
	},
}

func isDPMSOn(c *cobra.Command) (bool, error) {
	ctx := c.Context()
	var env []string
	if ctx != nil {
		if val, ok := ctx.Value(types.EnvKey).([]string); ok {
			env = val
		}
	}

	isWayland := false
	for _, e := range env {
		if len(e) > 16 && e[:16] == "WAYLAND_DISPLAY=" {
			isWayland = true
			break
		}
	}
	if !isWayland && os.Getenv("WAYLAND_DISPLAY") != "" {
		isWayland = true
	}

	if isWayland {
		// Sway: check if any output has DPMS on
		out, err := common.RunCmd(c, "swaymsg", "-t", "get_outputs").Output()
		if err != nil {
			return false, err
		}
		// Simple check for "dpms": true
		return strings.Contains(string(out), `"dpms": true`), nil
	}

	// X11: xset q
	out, err := common.RunCmd(c, "xset", "q").Output()
	if err != nil {
		return false, err
	}
	return strings.Contains(string(out), "Monitor is On"), nil
}

func setDPMS(c *cobra.Command, on bool) error {
	ctx := c.Context()
	var env []string
	if ctx != nil {
		if val, ok := ctx.Value(types.EnvKey).([]string); ok {
			env = val
		}
	}

	state := "off"
	if on {
		state = "on"
	}

	isWayland := false
	for _, e := range env {
		if len(e) > 16 && e[:16] == "WAYLAND_DISPLAY=" {
			isWayland = true
			break
		}
	}
	if !isWayland && os.Getenv("WAYLAND_DISPLAY") != "" {
		isWayland = true
	}

	if isWayland {
		return common.RunCmd(c, "swaymsg", "output * dpms "+state).Run()
	}

	xsetArg := "off"
	if on {
		xsetArg = "on"
	}
	return common.RunCmd(c, "xset", "dpms", "force", xsetArg).Run()
}

func init() {
	Command.AddCommand(lockCmd)
	Command.AddCommand(offCmd)
	Command.AddCommand(onCmd)
	Command.AddCommand(toggleCmd)
}
