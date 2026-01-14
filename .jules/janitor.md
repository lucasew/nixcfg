# Janitor's Journal - Critical Learnings

## 2026-01-14 - Refactor Duplicated Logic in `quicksync` Script
**Issue:** The `bin/quicksync` shell script contained two identical `if ! sd is riverwood; then ... fi` blocks.
**Root Cause:** The script needs to ensure an external machine (`riverwood`) is synchronized both before and after the main script logic runs. This was implemented by copying and pasting the same conditional block.
**Solution:** I extracted the duplicated conditional logic into a dedicated function called `sync_riverwood`. This function is now called twice at the appropriate points in the script, making the code cleaner and the script's intent clearer.
**Pattern:** Duplicated blocks of code in shell scripts should be refactored into functions. This improves readability, reduces the chance of introducing errors during future modifications, and adheres to the Don't Repeat Yourself (DRY) principle.
