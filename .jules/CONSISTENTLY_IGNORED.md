# Consistently Ignored Changes

This file lists patterns of changes that have been consistently rejected by human reviewers. All agents MUST consult this file before proposing a new change. If a planned change matches any pattern described below, it MUST be abandoned.

---

## IGNORE: Flagging Weak Default Authentication (SSH)

**- Pattern:** Flagging enabled SSH password authentication or the use of the weak 'changeme' initial password as a vulnerability.
**- Justification:** The user has explicitly rejected disabling password authentication and accepts the risk of the weak initial password.
**- Files Affected:** `nix/nodes/bootstrap/user.nix`, `nix/nodes/bootstrap/ssh.nix`

---

## IGNORE: Flagging Hardcoded rsync.net Credentials

**- Pattern:** Flagging the hardcoded `rsync.net` user fallback as a vulnerability or hardcoded secret.
**- Justification:** The user has specified that this hardcoded fallback is intentional and serves a specific purpose in their workflow.
**- Files Affected:** `bin/backup`, `bin/quicksync`

---

## IGNORE: Flagging Hardcoded OpenRouter Public Key

**- Pattern:** Flagging the hardcoded `OPENROUTER_API_KEY` fallback as a vulnerability.
**- Justification:** The hardcoded key has a zero-spend limit and is intended for public/testing usage when no personal key is provided.
**- Files Affected:** `bin/ai/codex_public`

---

## IGNORE: Enforcing Authentication on Xorg VNC

**- Pattern:** Changing `-SecurityTypes None` to enforce password authentication in the Xorg VNC script.
**- Justification:** The user accepts the risk of passwordless VNC as the service is intended to be exposed only on a trusted Tailscale network.
**- Files Affected:** `bin/svc/.vncd-xorg`

---

## IGNORE: Consolidating CI Workflows

**- Pattern:** Deleting legacy GitHub Actions workflows (`bump.yml`, `format.yml`, `prebuild.yml`, `nixpkgs-review.yml`, `ciborg-turbo.yml`) or merging them into `autorelease.yml`.
**- Justification:** Multiple attempts to consolidate these workflows have been rejected. The user prefers to keep them separate.
**- Files Affected:** `.github/workflows/*.yml`

---

## IGNORE: Security Hardenings in Specific Scripts

**- Pattern:** Proposing security fixes such as variable quoting in `bin/misc/clean-garbage-gcroots` to prevent arbitrary file deletion, slug validation in `pkg/sudo`, and removing shell execution in `pkg/wallpaper`.
**- Justification:** Specific security hardenings proposed by agents like Sentinel have been explicitly rejected and must not be reapplied.
**- Files Affected:** `bin/misc/clean-garbage-gcroots`, `nix/pkgs/workspaced/pkg/sudo/sudo.go`, `nix/pkgs/workspaced/pkg/wallpaper/wallpaper.go`

---

## IGNORE: Go Codebase Restructuring

**- Pattern:** Moving packages from `pkg/drivers/` to `pkg/` or splitting `pkg/common` into smaller packages (`pkg/config`, `pkg/host`, etc.) in `nix/pkgs/workspaced`.
**- Justification:** Extensive structural refactoring of the Go codebase is consistently rejected. The current structure (with `drivers` and `common`) should be preserved.
**- Files Affected:** `nix/pkgs/workspaced/pkg/**/*`

---

## IGNORE: Massive Automated Formatting

**- Pattern:** Submitting PRs that perform codebase-wide formatting (e.g., reformatting all Nix files) alongside unrelated functional changes.
**- Justification:** Large formatting diffs obscure the actual changes and are consistently rejected. Formatting should be done in dedicated PRs or only on touched lines.
**- Files Affected:** `flake.nix`, `nix/**/*.nix`, `mise.toml`

---

## IGNORE: Automated Dependency Updates

**- Pattern:** Automated PRs bumping dependencies such as versions in `mise.toml` or `go.mod`.
**- Justification:** Automated dependency bumps are consistently autoclosed/rejected. Dependency updates should be handled manually or are not desired in this repository.
**- Files Affected:** `mise.toml`, `nix/pkgs/workspaced/mise.toml`, `nix/pkgs/workspaced/go.mod`, `nix/pkgs/workspaced/go.sum`
