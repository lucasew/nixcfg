## Issue: CI workflow task failures
**Issue:** `mise run ci` and CI jobs failed due to missing `install:*`, `test:*`, and `codegen:*` task configurations in `mise.toml`.
**Root Cause:** The `autorelease.yml` workflow and local CI task expected wildcard and specific task definitions in `mise.toml` that were not defined.
**Solution:** Added `install:tools`, `test:dummy`, and `codegen:dummy` tasks, and defined explicit wildcard `depends` for `test`, `install`, and `codegen` in `mise.toml`.

## Issue: CI workflow task failures
**Issue:** `mise run ci` and CI jobs failed due to missing `install:*`, `test:*`, and `codegen:*` task configurations in `mise.toml`.
**Root Cause:** The `autorelease.yml` workflow and local CI task expected wildcard and specific task definitions in `mise.toml` that were not defined.
**Solution:** Added `install:tools`, `test:dummy`, and `codegen:dummy` tasks, and defined explicit wildcard `depends` for `test`, `install`, and `codegen` in `mise.toml`.
