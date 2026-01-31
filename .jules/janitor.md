# Janitor's Journal - Critical Learnings

## 2026-01-31 - Robust Directory Changing & CI Fixes

**Issue:** `bin/shim/workspaced` used fragile `cd` patterns (SC2103). CI failed due to `shfmt` errors in `bin/prelude/999-atuin.sh`, impure derivation usage (`teste-impure`) in `flake.nix`, `mise.toml` `lint:sh` task failing to find `shfmt`, Go unused field linting errors, and nested `mise.toml` trust issues.
**Root Cause:**
1. `bin/shim/workspaced`: Linear directory state changes.
2. `flake.nix`: `teste-impure` required `impure-derivations` feature disabled in CI.
3. `bin/prelude/999-atuin.sh`: Formatting violation (extra space).
4. `mise.toml`: `lint:sh` used `$(shfmt -f .)` which failed because `shfmt` wasn't in the global PATH.
5. `pkg/common/cas.go`: Unused field `hash` in `CASWriter` struct.
6. CI: `mise-action` fails on untrusted nested `mise.toml` files in `infra/` and `nix/pkgs/workspaced/`.
**Solution:**
1. Refactored `bin/shim/workspaced` to use subshells `( cd ... )`.
2. Commented out `teste-impure` in `flake.nix` as it blocks CI.
3. Ran `shfmt -w` on `bin/prelude/999-atuin.sh`.
4. Updated `mise.toml` to use `$(mise exec shfmt -- shfmt -f .)` ensuring tool availability.
5. Removed unused `hash` field from `CASWriter` struct in `pkg/common/cas.go`.
6. Consolidated all tools into root `mise.toml` and removed nested `mise.toml` files.
**Pattern:**
- Always use subshells `( cd ... )` for temporary directory changes.
- Ensure strict compliance with CI linters (`shfmt`, `nix flake check`, `golangci-lint`) before submitting.
- Disable or fix experimental features in Flakes that conflict with CI security settings.
- When using `mise` tools in command substitutions inside `mise.toml`, explicitly wrap them with `mise exec` if they aren't guaranteed to be in PATH.
- Remove unused struct fields to satisfy `unused` linter.
- Avoid nested `mise.toml` files in CI environments that require explicit trust; consolidate dependencies in the root configuration.

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
