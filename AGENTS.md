# Agent Instructions

## Security (Sentinel)

*   **Weak Initial Passwords:** The usage of `initialPassword = "changeme"` in `nix/nodes/bootstrap/user.nix` is a known and accepted configuration for bootstrapping. It is mitigated by `PasswordAuthentication = false` in SSH configuration. **Do not flag this as a vulnerability.**
