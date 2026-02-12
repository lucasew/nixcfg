package shellgen

import (
	"fmt"
	"workspaced/pkg/env"
)

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
	root, _ := env.GetDotfilesRoot()
	return fmt.Sprintf(`# Flag to indicate workspaced shell init is being used
export WORKSPACED_SHELL_INIT=1
export SD_ROOT=%q
export DOTFILES=%q
export NIXCFG_ROOT_PATH=%q
`, root, root, root), nil
}
