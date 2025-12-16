## 2024-05-22 - [Sentinel Init]
**Vulnerability:** N/A
**Learning:** Initialized Sentinel journal.
**Prevention:** N/A

## 2025-02-14 - [Disable SSH Password Authentication]
**Vulnerability:** SSH password authentication was enabled (`PasswordAuthentication = true`), and the default `initialPassword` for the user/root is set to "changeme". This combination poses a high risk of unauthorized access if the initial password is not changed immediately or if the machine is exposed to the internet.
**Learning:** Declarative systems with "initial" secrets (like passwords) can be dangerous if they default to insecure values and rely on manual intervention to secure them. Hardening the access method (disabling password auth) acts as a strong compensating control.
**Prevention:** Enforce key-based authentication by setting `services.openssh.settings.PasswordAuthentication = false` in the base configuration (`nix/nodes/bootstrap/ssh.nix`). This ensures that even if the password is weak, it cannot be used for remote SSH access.
