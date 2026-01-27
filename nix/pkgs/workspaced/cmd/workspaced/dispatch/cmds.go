package dispatch

import (
	"fmt"

	"github.com/spf13/cobra"
	"workspaced/pkg/cmd"
)

func init() {
	DaemonCmd.Run = func(c *cobra.Command, args []string) {
		// We can't import main (cycle), so we expect the main package to set this Run,
		// OR we handle "daemon" specially in root.go.
		// For now, let's print a placeholder or fail if executed directly without main's override.
		fmt.Println("Daemon must be run via the main entrypoint.")
	}
}

var DaemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Run the workspaced daemon",
}

var ModnCmd = &cobra.Command{
	Use:   "modn",
	Short: "Rotate workspaces across outputs",
	RunE: func(c *cobra.Command, args []string) error {
		_, err := cmd.RunModn()
		return err
	},
}

var RofiCmd = &cobra.Command{
	Use:   "rofi",
	Short: "Rofi workspace switcher",
	RunE: func(c *cobra.Command, args []string) error {
		ctx := c.Context()
		env, _ := ctx.Value("env").([]string)
		// If env is nil (e.g. CLI run), use os.Environ?
		// cmd.RunRofi defaults to using its own logic if we pass nil?
		// But in `cmd.RunRofi`, we check `if len(env) > 0`.
		// If we run from CLI, we want `os.Environ()`.
		// If we run from Daemon, we want `req.Env`.
		// So if context is missing env, we should pass nil (and let caller handle?)
		// Wait, `cmd.RunRofi` uses `cmd.Env = env` if len > 0.
		// If len == 0, `cmd.Env` inherits from parent process (daemon or shell).
		// If CLI, shell env is inherited.
		// If Daemon, daemon env is inherited (NOT what we want).
		// So Daemon MUST pass env.
		// CLI CAN pass env (os.Environ) or let it inherit.

		// If we are in CLI, `c.Context()` might not have "env".
		// In that case, we can rely on inheritance (so passing nil/empty is fine, as long as `RunRofi` doesn't clear env).
		// `cmd.RunRofi`: `if len(env) > 0 { cmd.Env = env }`.
		// If empty, `cmd.Env` is nil, causing `exec.Command` to use `os.Environ()`.
		// This is correct for CLI.
		// For Daemon, we MUST populate context.

		_, err := cmd.RunRofi(args, env)
		return err
	},
}
