# Janitor's Journal - Critical Learnings

## 2024-07-25 - Always Quote Shell Variables
**Issue:** A refactoring to make an SSH target configurable in a bash script introduced a command injection vulnerability.
**Root Cause:** The script used an unquoted variable in an `ssh` command (`ssh $quicksync_ssh_target`). If the variable contained spaces or shell metacharacters (e.g., `;`, `|`, `&`), it would be subject to word splitting and command execution.
**Solution:** The variable was double-quoted (`ssh "$quicksync_ssh_target"`). Double quotes prevent word splitting and glob expansion, ensuring the variable's value is treated as a single, literal string.
**Pattern:** In shell scripts, always double-quote variables that contain user-controllable, file path, or otherwise complex string data, especially when used as arguments to commands. This prevents security vulnerabilities and unexpected behavior.
