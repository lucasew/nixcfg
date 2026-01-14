# Janitor's Journal - Critical Learnings
## 2026-01-14 - Simplify quicksync script by removing redundant call
**Issue:** The `bin/quicksync` script contained a duplicate `ssh` command at the end of the script.
**Root Cause:** The script was likely edited over time, and the redundant call was unintentionally left in. The initial call ensures that the `riverwood` machine is synced before the local machine syncs, making the final call unnecessary.
**Solution:** I removed the final `ssh lucasew@riverwood sdw quicksync` block from the script.
**Pattern:** Periodically review scripts for redundant or unnecessary operations, especially in logic that has been modified incrementally. A final check for redundant calls can simplify logic and improve performance.