# Sentinel's Journal

## 2024-05-22 - [Sentinel Init]
**Vulnerability:** N/A
**Learning:** Initialized Sentinel journal.
**Prevention:** N/A

## 2024-05-22 - [Disabled SSH Password Authentication]
**Vulnerability:** SSH Password Authentication was enabled with a weak initial password "changeme".
**Learning:** Default configurations can inadvertently expose systems if not overridden, especially when weak defaults like "changeme" are used for convenience.
**Prevention:** Explicitly disabled `PasswordAuthentication` in `nix/nodes/bootstrap/ssh.nix` to enforce key-based authentication.
