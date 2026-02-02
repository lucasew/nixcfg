package prelude

// GenerateDaemon generates daemon startup code
func GenerateDaemon() (string, error) {
	return `# Start workspaced daemon if available
if command -v workspaced >/dev/null 2>&1; then
	(workspaced daemon --try &) &>/dev/null
fi
`, nil
}

// GenerateFlags generates shell init flags
func GenerateFlags() (string, error) {
	return `# Flag to indicate workspaced shell init is being used
export WORKSPACED_SHELL_INIT=1
`, nil
}
