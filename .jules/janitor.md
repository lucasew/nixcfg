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

## 2026-02-04 - Fix Broken Linting Tasks

**Issue:** `lint:sh` failed because `shfmt` wasn't found in command substitution, and `lint:go` failed due to version mismatch and existing code errors.
**Root Cause:** Command substitution `$(shfmt ...)` in `mise.toml` executed in the host shell instead of the mise environment. Additionally, `golangci-lint` 1.61.0 was incompatible with Go 1.25.6.
**Solution:** Wrapped the command substitution with `mise exec shfmt -- ...` and updated `golangci-lint` to 1.64.5. Fixed underlying ShellCheck and Go lint violations.
**Pattern:** When using command substitution in `mise` tasks, explicitly invoke tools via `mise exec` if they are not in the system PATH.

## 2026-02-04 - Fix Missing Tool Definitions

**Issue:** CI failed because `shfmt`, `shellcheck`, and `ruff` were not installed in the CI environment, causing `mise exec` to fail or implicitly attempt installation which failed.
**Root Cause:** These tools were missing from the `[tools]` section in the root `mise.toml`, so `mise run install` did not install them.
**Solution:** Explicitly added `shfmt`, `shellcheck`, and `ruff` to `mise.toml`.
**Pattern:** Always define tool dependencies explicitly in `mise.toml` to ensure deterministic builds in CI.

## 2026-02-04 - Robust Linting Pipeline

**Issue:** `lint:sh` failed in CI because command substitution hid errors and passed empty arguments to `shellcheck`. Additionally, tools in root `mise.toml` might not be installed explicitly by the `install` task.
**Root Cause:** Command substitution execution order and stderr handling in CI environment is fragile. `mise run install` does not implicitly run `mise install` for root tools.
**Solution:** Replaced command substitution with `shfmt -f=0 . | xargs -0 -r ...` to robustly pipe file lists. Added `install:tools` task to explicitly run `mise install`.
**Pattern:** Prefer pipelines with `xargs -0 -r` over command substitution for passing file lists to tools, as it handles empty lists and filenames with spaces correctly.
- 2026-06-03: [errcheck] Always check error returns for exec.RunCmd calls, reporting failures via logging.ReportError.
<!-- Trigger CI retry -->
<!-- Trigger CI retry 2 -->
