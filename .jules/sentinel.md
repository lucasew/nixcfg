# Sentinel Journal

## 2025-02-24 - Hardcoded Rsync.net Credentials
**Vulnerability:** Hardcoded username and host `de3163@de3163.rsync.net` found in `bin/backup` and `bin/quicksync`.
**Learning:** Hardcoded credentials are often intentional in personal/single-user repositories for convenience. Over-aggressive warning/blocking can disrupt workflows.
**Prevention:** Provide environment variable overrides (`RSYNCNET_USER`) for flexibility/security, but silently fall back to the "intentional" hardcoded values to preserve existing behavior and avoid noise.
