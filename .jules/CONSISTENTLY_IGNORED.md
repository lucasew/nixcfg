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

## IGNORE: Out-of-Scope Modifications by Scoped Agents

**- Pattern:** Agents (like Docs, Janitor, or Sentinel) modifying files outside their strictly defined allowed paths in order to fix unrelated issues (e.g., linters, CI).
**- Justification:** Agents must strictly adhere to their designated scope constraints (e.g., Docs agent must not change executable logic; Janitor must not modify `.github/workflows/**`).
**- Files Affected:** Any files outside an agent's allowed paths

## IGNORE: Installing workspaced via nix profile

**- Pattern:** Attempting to install the `workspaced` CLI tool natively using `nix profile install github:lucasew/workspaced`.
**- Justification:** Attempts to install `workspaced` via `nix profile install` in CI workflows are consistently rejected.
**- Files Affected:** `.github/workflows/*.yml`

## IGNORE: Adding id-token: write to GitHub Actions

**- Pattern:** Adding the `id-token: write` permission to GitHub Actions workflows like `autorelease.yml` and `ciborg-turbo.yml`.
**- Justification:** Attempts to add `id-token: write` permissions to workflow files are consistently rejected.
**- Files Affected:** `.github/workflows/*.yml`

## IGNORE: Formatting CONSISTENTLY_IGNORED.md as a Markdown Table

**- Pattern:** Restructuring `.jules/CONSISTENTLY_IGNORED.md` from its native list-based format (`## IGNORE:...`) into a Markdown table.
**- Justification:** The file must maintain its native list structure to ensure automated parsers and human reviewers can correctly extract the rules. Markdown tables are consistently rejected.
**- Files Affected:** `.jules/CONSISTENTLY_IGNORED.md`

## IGNORE: Automated Dependency Updates

**- Pattern:** Automated PRs bumping tool versions in `workspaced.lock.json` (e.g., zed-industries/zed).
**- Justification:** Automated dependency updates are consistently autoclosed/rejected as the project manages dependencies manually or through a different process.
**- Files Affected:** `workspaced.lock.json`

## IGNORE: Formatting Agent Journals with Empty Lines After Headers

**- Pattern:** Adding empty lines between Markdown headers and semantically structured metadata tags (e.g., `**Vulnerability:**`) in agent journals to comply with perceived Prettier formatting rules.
**- Justification:** Semantic metadata tags must remain contiguous and preserved structurally. Hallucinated rules about empty lines after headers are explicitly rejected.
**- Files Affected:** `.jules/*.md`

## IGNORE: Adding Dummy/Fallback Tasks and Wildcard Dependencies in mise.toml

**- Pattern:** Adding dummy or fallback tasks (e.g., `[tasks."test:fallback"]`) and changing task dependencies to use wildcards (e.g., `depends = ["test:*"]`) in `mise.toml` in an attempt to fix CI wildcard resolution failures.
**- Justification:** Attempts to fix CI wildcard issues by cluttering `mise.toml` with dummy tasks are consistently rejected. The configuration should be kept clean, and actual tasks should be resolved correctly without fallback clutter.
**- Files Affected:** `mise.toml`
