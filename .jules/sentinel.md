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

## 2025-12-25 - [Intentional Hardcoded Credential]
**Vulnerability:** Hardcoded `rsync.net` credential fallback in `bin/backup` and `bin/quicksync`.
**Learning:** The user has specified that this hardcoded fallback is intentional and should not be removed. It serves a specific purpose in their workflow. This is a user-accepted risk.
**Prevention:** Do not flag the hardcoded `rsync.net` credential in `bin/backup` and `win/quicksync` as a vulnerability. This is an accepted risk and intentional design.

## 2024-05-23 - [Insecure Default Network Binding]
**Vulnerability:** The Wayland VNC server script (`bin/svc/.vncd-wayland`) was binding to `0.0.0.0` by default, exposing the service to all network interfaces.
**Learning:** Network services should always default to the most secure configuration. Binding to `127.0.0.1` (localhost) by default prevents unintentional exposure to untrusted networks.
**Prevention:** Always explicitly define the binding address for network services, defaulting to localhost unless external access is an explicit requirement. Provide a secure way to override the default, such as through an environment variable.

## 2024-05-24 - [User-Accepted Risk: Passwordless VNC on Trusted Network]
**Vulnerability:** The Xorg VNC server script (`bin/svc/.vncd-xorg`) intentionally disables authentication (`-SecurityTypes None`).
**Learning:** The user has explicitly stated this is by design and accepts the risk, as the service is only exposed on a trusted network (Tailscale). This overrides the general "defense in depth" principle for this specific case.
**Prevention:** Do not flag the disabled security on `bin/svc/.vncd-xorg` as a vulnerability in the future. This is a user-accepted and intentional configuration.

## 2024-05-25 - [Fix Command Injection in backup and quicksync scripts]
**Vulnerability:** The `backup` and `quicksync` scripts were vulnerable to command injection. They used the `RSYNCNET_USER` environment variable directly in `ssh`, `rsync`, and `git` commands without proper validation.
**Learning:** This allowed a malicious user to inject arbitrary command-line options by crafting a malicious user string (e.g., `-oProxyCommand=...`), leading to command execution. This highlighted the critical need for robust input validation, especially for environment variables that can be controlled by users.
**Prevention:** To mitigate this, I implemented a validation check to ensure that the `RSYNCNET_USER` variable does not start with a hyphen (`-`). This simple yet effective measure prevents the injection of malicious options, ensuring that the scripts handle user-provided data securely.

## 2024-05-27 - [Command Injection via arp-scan output]
**Vulnerability:** A command injection risk existed in the `bin/misc/dns-cgi` script. The script used the vendor string from `arp-scan`'s output to generate hostnames. An attacker on the local network could spoof their MAC address vendor to include shell metacharacters, which could be executed by a downstream consumer of the generated hosts file.
**Learning:** All output from network scanning tools must be treated as untrusted input. Relying on descriptive but potentially user-controllable fields for generating identifiers is a security risk.
**Prevention:** Sanitize all input from external network sources. When possible, use non-descriptive, machine-generated identifiers like MAC addresses (as implemented in the fix) instead of potentially malicious strings. Use flags like `--numeric` to suppress resolution of identifiers to potentially unsafe strings.

## 2024-05-28 - [Sensitive Data in Process List]
**Vulnerability:** The `bin/misc/dns-cgi` script passed the `ZEROTIER_TOKEN` directly in the `curl` command line arguments (`-H "Authorization: ..."`). This allows any user on the system to see the token by listing processes (`ps aux`).
**Learning:** Command line arguments are public information on most Unix-like systems. Sensitive data like tokens and passwords should never be passed via arguments.
**Prevention:** Use configuration files or stdin to pass secrets to commands. For `curl`, the `-K` / `--config` option allows reading headers and other sensitive parameters from a file, which can be secured with file permissions.
