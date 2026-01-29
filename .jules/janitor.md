# Janitor's Journal - Critical Learnings

## 2026-01-24 - Robust Shell Argument Parsing

**Issue:** Argument parsing in `bin/prelude/020-notification.sh` was fragile, using double-shift which could consume flags as values if arguments were missing, and using unsafe `[ ... ]` syntax.
**Root Cause:** Manual `while` loop with `case` and `shift` without checking if the next argument existed or was a valid value.
**Solution:** Refactored to use `local` variables, `[[ ... ]]`, and verified `${2:-}` existence before shifting.
**Pattern:** When parsing arguments manually in Bash, always validate `${2:-}` before `shift`ing to assign a value, and use `[[ -n ... ]]` for robustness.

## 2026-01-18 - Simplify Prelude Sourcing

**Issue:** `bin/source_me` used a fragile `ls | sort` subshell to iterate over scripts, and `bin/prelude/999-path-legacy.sh` was dead code.
**Root Cause:** Legacy iteration pattern and unused function left over from previous refactors.
**Solution:** Replaced `ls` loop with a safer glob pattern `"$SD_ROOT/prelude/"*` and deleted the unused file.
**Pattern:** Prefer shell globs over parsing `ls` output for file iteration.
