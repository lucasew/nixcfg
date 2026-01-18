# Janitor's Journal - Critical Learnings

## 2026-01-18 - Refactor bin/misc/toggle-monitor

**Issue:** Vulnerable and brittle shell script using unquoted variables, non-POSIX `==` in tests, and command substitution for iteration.
**Root Cause:** Quick scripting without strict safety practices.
**Solution:** Replaced command substitution with a `while read` loop, quoted all variables, and switched to POSIX-compliant syntax (`=`). Removed unnecessary function wrapper.
**Pattern:** Always use `while read` loops to process command output safely instead of `$(...)` command substitution, and quote variables to prevent word splitting.
