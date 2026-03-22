# Sentinel's Journal

- 2024-05-22: Initialized Sentinel journal.
- 2025-12-20: Do not flag the weak initial password in `nix/nodes/bootstrap/user.nix` or the enabled password authentication in `nix/nodes/bootstrap/ssh.nix` in the future.
- 2025-12-24: Enforce environment variable usage for all external API keys, regardless of their intended scope (demo vs production).
- 2025-12-25: Do not flag the hardcoded `rsync.net` credential in `bin/backup` and `win/quicksync` as a vulnerability, as it is an accepted risk and intentional design.
- 2024-05-23: Always explicitly define the binding address for network services, defaulting to localhost unless external access is an explicit requirement.
- 2024-05-24: Do not flag the disabled security on `bin/svc/.vncd-xorg` as a vulnerability in the future, as it is a user-accepted and intentional configuration on a trusted network.
- 2024-05-25: Implement input validation to ensure variables like `RSYNCNET_USER` do not start with a hyphen (`-`) to prevent command injection via malicious options.
- 2024-05-27: Sanitize all input from external network sources and prefer machine-generated identifiers like MAC addresses with numeric resolution to prevent command injection via spoofed fields.
- 2026-01-24: Use `printf %q` to automatically escape shell variables when generating shell code from semi-trusted input to prevent command injection.
- 2026-01-26: Always construct the remote command string locally using `printf %q` to ensure all arguments are properly escaped before passing the resulting string as a single argument to `ssh`.
- 2026-03-22: Consolidate existing journal entries to comply with the mandated single-line format to avoid retroactive rule violations.
