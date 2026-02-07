package shellgen

// GenerateDaemon generates the shell script to start the workspaced daemon.
//
// It checks if the `workspaced` binary is available in the path.
// If found, it attempts to start the daemon in the background (`--try` flag suggests idempotency).
// The process is detached to ensure the shell startup doesn't hang waiting for the daemon.
func GenerateDaemon() (string, error) {
	return `# Start workspaced daemon if available
if command -v workspaced >/dev/null 2>&1; then
	(workspaced daemon --try &) &>/dev/null
fi
`, nil
}

// GenerateFlags generates environment variables marking the shell state.
//
// It exports `WORKSPACED_SHELL_INIT=1`, which allows other tools
// or scripts to detect that the shell has been initialized by workspaced.
func GenerateFlags() (string, error) {
	return `# Flag to indicate workspaced shell init is being used
export WORKSPACED_SHELL_INIT=1
`, nil
}
