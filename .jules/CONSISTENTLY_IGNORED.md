# Consistently Ignored Changes

This file lists patterns of changes that have been consistently rejected by human reviewers. All agents MUST consult this file before proposing a new change. If a planned change matches any pattern described below, it MUST be abandoned.

| Topic | Pattern | Justification | Files Affected |
| :--- | :--- | :--- | :--- |
| **IGNORE: Flagging Weak Default Authentication (SSH)** | Flagging enabled SSH password authentication or the use of the weak 'changeme' initial password as a vulnerability. | The user has explicitly rejected disabling password authentication and accepts the risk of the weak initial password. | `nix/nodes/bootstrap/user.nix`, `nix/nodes/bootstrap/ssh.nix` |
| **IGNORE: Flagging Hardcoded rsync.net Credentials** | Flagging the hardcoded `rsync.net` user fallback as a vulnerability or hardcoded secret. | The user has specified that this hardcoded fallback is intentional and serves a specific purpose in their workflow. | `bin/backup`, `bin/quicksync` |
| **IGNORE: Flagging Hardcoded OpenRouter Public Key** | Flagging the hardcoded `OPENROUTER_API_KEY` fallback as a vulnerability. | The hardcoded key has a zero-spend limit and is intended for public/testing usage when no personal key is provided. | `bin/ai/codex_public` |
| **IGNORE: Enforcing Authentication on Xorg VNC** | Changing `-SecurityTypes None` to enforce password authentication in the Xorg VNC script. | The user accepts the risk of passwordless VNC as the service is intended to be exposed only on a trusted Tailscale network. | `bin/svc/.vncd-xorg` |
| **IGNORE: Consolidating CI Workflows** | Deleting legacy GitHub Actions workflows (`bump.yml`, `format.yml`, `prebuild.yml`, `nixpkgs-review.yml`, `ciborg-turbo.yml`) or merging them into `autorelease.yml`. | Multiple attempts to consolidate these workflows have been rejected. The user prefers to keep them separate. | `.github/workflows/*.yml` |
