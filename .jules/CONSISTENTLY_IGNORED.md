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
