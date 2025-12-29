# Janitor's Journal - Critical Learnings

## 2024-07-22 - Securely Handling Credentials in Shell Scripts
**Issue:** The `bin/quicksync` script contained a hardcoded fallback for the `RSYNCNET_USER` credential, posing a security risk and reducing maintainability.
**Root Cause:** A default value was likely included for development convenience and was not removed for production use.
**Solution:** I removed the hardcoded fallback and implemented the bash construct `${RSYNCNET_USER:?}`. This ensures the script exits immediately with an error if the required environment variable is not set.
**Pattern:** Shell scripts requiring credentials from environment variables should use the `"${VAR_NAME:?}"` pattern. This enforces a secure failure mode by preventing execution with default or missing credentials, avoiding hardcoded secrets.