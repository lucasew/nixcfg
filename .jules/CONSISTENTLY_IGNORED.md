# Consistently Ignored Changes

This file lists patterns of changes that have been consistently rejected by human reviewers. All agents MUST consult this file before proposing a new change. If a planned change matches any pattern described below, it MUST be abandoned.

---

## IGNORE: Out-of-Scope Agent Modifications

**- Pattern:** Agents modifying files outside their strictly allowed paths (e.g., modifying `mise.toml` or `.github/workflows/` when limited to `src/`).
**- Justification:** Agents must adhere to their designated scope and leave out-of-scope fixes to the appropriate agent.
**- Files Affected:** `*`

---

## IGNORE: Adding Dummy/Fallback mise Tasks

**- Pattern:** Adding fallback tasks like `[tasks."test:dummy"]` or `[tasks."install:tools"]` to `mise.toml` to fix CI wildcard issues.
**- Justification:** Adding dummy/fallback tasks adds unnecessary noise and clutters the configuration.
**- Files Affected:** `mise.toml`, `*/mise.toml`

---

## IGNORE: Incomplete Utility Implementations

**- Pattern:** Adding centralized utilities (e.g., error reporting) without migrating existing scattered usages across the codebase to use them.
**- Justification:** Implementing new utilities without adopting them in existing code creates dead code and fails to reduce source complexity.
**- Files Affected:** `src/utils/*.sh`, `tests/utils/*.sh`

---

## IGNORE: Automated Dependency Updates in Lockfiles

**- Pattern:** Automated dependency version bumps or tool updates in `workspaced.lock.json`.
**- Justification:** These automated lockfile bumps are consistently autoclosed and should be ignored.
**- Files Affected:** `workspaced.lock.json`

---

## IGNORE: Unrequested Lockfile Modifications

**- Pattern:** Committing unrequested lockfile modifications (e.g., `workspaced.lock.json`) caused by implicit modifications during linting or testing.
**- Justification:** Lockfile modifications unrelated to the core task add noise and should be restored or unstaged before committing.
**- Files Affected:** `workspaced.lock.json`

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
