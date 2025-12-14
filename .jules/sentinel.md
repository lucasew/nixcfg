# Sentinel Journal

## 2025-02-24 - Hardcoded Rsync.net Credentials
**Vulnerability:** Hardcoded username and host `de3163@de3163.rsync.net` found in `bin/backup` and `bin/quicksync`.
**Learning:** Hardcoded credentials (even just usernames/hosts) make scripts less portable and expose infrastructure details in the codebase.
**Prevention:** Use environment variables for all credentials and service locations. Implemented `RSYNCNET_USER` with a deprecated fallback to avoid immediate breakage.
