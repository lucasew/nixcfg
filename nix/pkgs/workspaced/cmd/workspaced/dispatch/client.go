package dispatch

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	// The dispatch command itself acts as the entry point for local execution.
	// If executed locally (without subcommands?), we might want to allow forwarding?
	// But the PR says: "workspaced has only the subcommands dispatch and daemon".
	// "daemon server exposes the dispatch subcommand".
	// So `workspaced dispatch modn` is valid.

	// But what about the CLI client logic? "se o daemon tá rodando pede pelo socket... se não... faz diretamente".
	// This logic needs to be SOMEWHERE.
	// It's likely in `dispatch.Command.PersistentPreRun`?
	// Or we wrap the execution?

	// If `workspaced dispatch modn` is called:
	// We want to check if daemon is running.
	// If yes -> RPC call to daemon -> returns output.
	// If no -> continue execution (which runs `modnCmd.Run`).

	Command.PersistentPreRunE = func(c *cobra.Command, args []string) error {
		// Detect if we are running INSIDE the daemon to avoid recursion/loops.
		// The daemon calls `ExecuteContext` on this command.
		// We can check the context or an env var.
		// We pass "env" in context. We can also pass a "daemon_mode" flag.

		ctx := c.Context()
		if ctx.Value("daemon_mode") == true {
			return nil // We are the daemon, proceed to local logic
		}

		// We are the client. Try to connect to daemon.
		output, connected, err := TryRemote(c, args)
		if connected {
			if output != "" {
				fmt.Print(output)
			}
			// We handled it remotely. Stop local execution.
			// Returning an error stops execution, but prints usage?
			// SilenceUsage should be set.
			// Ideally we exit here? Or return a specific error that is handled gracefully?
			if err != nil {
				return err
			}
			os.Exit(0) // Clean exit after remote execution
		}

		// Daemon not connected, fall through to local execution
		return nil
	}
}

// Client logic (TryRemote) moved here or similar
// We need the `Request` struct which is defined in `daemon` package... circular dependency!
// `daemon` imports `dispatch`. `dispatch` cannot import `daemon`.
// So `Request` struct must be shared or defined in `dispatch` or `cmd/workspaced` (root)?
// OR we just define a local struct here for encoding.

type Request struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
	Env     []string `json:"env"`
}

type Response struct {
	Output string `json:"output"`
	Error  string `json:"error"`
}

func getSocketPath() string {
	runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
	if runtimeDir == "" {
		runtimeDir = fmt.Sprintf("/run/user/%d", os.Getuid())
	}
	return filepath.Join(runtimeDir, "workspaced.sock") // Duplicated logic, acceptable for decoupling
}

func TryRemote(c *cobra.Command, args []string) (string, bool, error) {
	socketPath := getSocketPath()
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return "", false, nil
	}
	defer conn.Close()

	// Reconstruct command: "modn" + args...
	// `c` is the command being run (e.g. modnCmd).
	// We need the name of the subcommand relative to `dispatch`.
	// If I run `workspaced dispatch modn`, `c` is `modnCmd`. `args` are flags?
	// Cobra passes non-flag args in `args`.

	// Wait, `PersistentPreRun` runs on `modnCmd` too.
	// `c.Name()` is "modn".
	// `args` are the arguments to modn.

	req := Request{
		Command: c.Name(),
		Args:    args,
		Env:     os.Environ(),
	}

	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(req); err != nil {
		return "", true, fmt.Errorf("failed to send request: %w", err)
	}

	var resp Response
	decoder := json.NewDecoder(conn)
	if err := decoder.Decode(&resp); err != nil {
		return "", true, fmt.Errorf("failed to decode response: %w", err)
	}

	if resp.Error != "" {
		return resp.Output, true, fmt.Errorf(resp.Error)
	}

	return resp.Output, true, nil
}
