# Janitor's Journal - Critical Learnings

## 2026-01-18 - Fix unsafe shell expansion in toggle-monitor

**Issue:** `bin/misc/toggle-monitor` used unquoted command substitution `handler $(...)` which passes variable arguments depending on `xrandr` output. If output was missing or malformed, it caused syntax errors in `test`.
**Root Cause:** Implicit reliance on word splitting for argument parsing and lack of input validation.
**Solution:** Replaced with a `while read` loop to safely process output line by line. Quoted all variables and switched to `[[ ... ]]` for safer comparisons.
**Pattern:** Avoid `cmd $(generator)` where `generator` output is complex. Use `generator | while read ...` instead.
