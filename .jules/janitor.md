# Janitor's Journal - Critical Learnings

## 2024-07-25 - Simplify Shell Script Output

**Issue:** The `bin/backup` script was printing its final summary message to the console twice—once with a standard `echo` and again via a notification function that also writes to `stderr`.

**Root Cause:** This redundancy likely stemmed from an oversight during development. A command substitution block captured multi-line output into a variable, which was then explicitly printed before being passed to a helper function that *also* printed it, leading to duplicated output.

**Solution:** I removed the superfluous `echo "$variable"` command. The existing `update_status` helper function already handled the responsibility of printing the final block to the console, so the explicit echo was unnecessary.

**Pattern:** A common simplification pattern in this repository's shell scripts is to replace temporary file creation (`mktemp`) and I/O (`tee`, `cat`) for capturing command block output with direct capture into a variable via command substitution (`output=$(...)`). A single, dedicated function should then be responsible for displaying this captured output to avoid redundant commands and improve script clarity.
