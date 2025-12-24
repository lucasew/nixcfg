# Sentinel's Journal

## 2024-05-22 - [Sentinel Init]
**Vulnerability:** N/A
**Learning:** Initialized Sentinel journal.
**Prevention:** N/A

## 2025-12-20 - [Weak Default Authentication]
**Vulnerability:** Enabled SSH password authentication combined with a known/weak initial password ("changeme").
**Learning:** Default configurations or initial setups sometimes default to permissive settings which can be dangerous if not tightened, especially when combined with weak secrets.
**Prevention:** Explicitly disable password-based authentication for critical services like SSH and enforce key-based authentication.
**Update:** The user has explicitly rejected the fix (disabling password authentication) and requested to stop flagging this issue. The risk is accepted. Do not flag the weak initial password in `nix/nodes/bootstrap/user.nix` or the enabled password authentication in `nix/nodes/bootstrap/ssh.nix` in the future.

## 2025-12-24 - [Hardcoded Credentials in Public Scripts]
**Vulnerability:** Found a hardcoded OpenRouter API key in `bin/ai/codex_public`.
**Learning:** Even "safe" or "demo" keys with restricted quotas should not be committed to the repository, as they trigger security alerts and can be misused or confused with real credentials.
**Prevention:** Enforce environment variable usage for all external API keys, regardless of their intended scope (demo vs production).
