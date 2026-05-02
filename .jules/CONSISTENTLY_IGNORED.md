# Consistently Ignored Changes

This file lists patterns of changes that have been consistently rejected by human reviewers. All agents MUST consult this file before proposing a new change. If a planned change matches any pattern described below, it MUST be abandoned.

---

| Pattern | Justification | Files Affected |
| :--- | :--- | :--- |
| Flagging enabled SSH password authentication or the use of the weak 'changeme' initial password as a vulnerability. | The user has explicitly rejected disabling password authentication and accepts the risk of the weak initial password. | `nix/nodes/bootstrap/user.nix`, `nix/nodes/bootstrap/ssh.nix` |
| Flagging the hardcoded `rsync.net` user fallback as a vulnerability or hardcoded secret. | The user has specified that this hardcoded fallback is intentional and serves a specific purpose in their workflow. | `bin/backup`, `bin/quicksync` |
| Flagging the hardcoded `OPENROUTER_API_KEY` fallback as a vulnerability. | The hardcoded key has a zero-spend limit and is intended for public/testing usage when no personal key is provided. | `bin/ai/codex_public` |
| Changing `-SecurityTypes None` to enforce password authentication in the Xorg VNC script. | The user accepts the risk of passwordless VNC as the service is intended to be exposed only on a trusted Tailscale network. | `bin/svc/.vncd-xorg` |
| Deleting legacy GitHub Actions workflows (`bump.yml`, `format.yml`, `prebuild.yml`, `nixpkgs-review.yml`, `ciborg-turbo.yml`) or merging them into `autorelease.yml`. | Multiple attempts to consolidate these workflows have been rejected. The user prefers to keep them separate. | `.github/workflows/*.yml` |
| Security fixes in specific scripts (e.g., `clean-garbage-gcroots`, slug validation, removing shell execution) | Specific security hardenings proposed by Sentinel have been explicitly rejected and must not be reapplied. | `bin/misc/clean-garbage-gcroots`, `pkg/sudo`, `pkg/wallpaper` |
| Go codebase restructuring (`pkg/drivers` removal) | The workspaced codebase MUST maintain its structure with `pkg/drivers/` containing implementation packages. Attempts to flatten this structure by promoting packages to `pkg/` or splitting `pkg/common` are consistently rejected. | `pkg/drivers/**`, `pkg/common/**` |
| "Trojan Horse" PRs | PRs that hide massive, unrequested formatting changes to Nix files (e.g., `flake.nix`, `nix/**/*.nix`) inside unrelated fixes are consistently rejected. | `flake.nix`, `nix/**/*.nix` |
| Automated dependency updates | Automated dependency updates bumping tool versions in `mise.toml` or modules in `go.mod` (e.g., node, ruff, sqlite, gopls) are consistently autoclosed/rejected and should be ignored. | `mise.toml`, `go.mod` |
| Rejecting agent journal formatting | PRs or reviews rejecting the single-line format (`- YYYY-MM-DD: [Title] Description.`) for agent journals or requesting the retention of multiline structures with bold metadata fields (like `**Issue:**`) must be ignored. | `.jules/sentinel.md`, `.jules/janitor.md` |
