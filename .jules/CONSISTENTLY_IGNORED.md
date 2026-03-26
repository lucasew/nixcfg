# Consistently Ignored Changes

This file lists patterns of changes that have been consistently rejected by human reviewers. All agents MUST consult this file before proposing a new change. If a planned change matches any pattern described below, it MUST be abandoned.

| Pattern | Justification | Files Affected |
| :--- | :--- | :--- |
| Flagging enabled SSH password authentication or the use of the weak 'changeme' initial password as a vulnerability. | The user has explicitly rejected disabling password authentication and accepts the risk of the weak initial password. | `nix/nodes/bootstrap/user.nix`, `nix/nodes/bootstrap/ssh.nix` |
| Flagging the hardcoded `rsync.net` user fallback as a vulnerability or hardcoded secret. | The user has specified that this hardcoded fallback is intentional and serves a specific purpose in their workflow. | `bin/backup`, `bin/quicksync` |
| Flagging the hardcoded `OPENROUTER_API_KEY` fallback as a vulnerability. | The hardcoded key has a zero-spend limit and is intended for public/testing usage when no personal key is provided. | `bin/ai/codex_public` |
| Changing `-SecurityTypes None` to enforce password authentication in the Xorg VNC script. | The user accepts the risk of passwordless VNC as the service is intended to be exposed only on a trusted Tailscale network. | `bin/svc/.vncd-xorg` |
| Deleting legacy GitHub Actions workflows (`bump.yml`, `format.yml`, `prebuild.yml`, `nixpkgs-review.yml`, `ciborg-turbo.yml`) or merging them into `autorelease.yml`. | Multiple attempts to consolidate these workflows have been rejected. The user prefers to keep them separate. | `.github/workflows/*.yml` |
| Refactoring to flat package structure by removing `pkg/drivers/` or moving packages to `pkg/` or splitting `pkg/common`. | The project explicitly maintains its structure with `pkg/drivers/` containing implementation packages. Attempts to flatten this structure are consistently rejected. | `pkg/**` |
| Formatting massive amounts of Nix files inside unrelated fixes ("Trojan Horse" PRs). | Massive, unrequested formatting changes to Nix files (e.g., `flake.nix`, `nix/**/*.nix`) are consistently rejected. | `flake.nix`, `nix/**/*.nix` |
| Automated dependency updates bumping tool versions in `mise.toml` or modules in `go.mod` (e.g., node, ruff, sqlite, gopls). | These updates are consistently autoclosed/rejected and should be ignored. | `mise.toml`, `go.mod` |
