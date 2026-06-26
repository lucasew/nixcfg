---
name: notify-me
description: >
  Sending messages "to me" (notifications) using sendmail on this dotfiles/NixOS setup.
  The sendmail command is backed by telegram-sendmail, which delivers to
  the user's Telegram. Use for "send X to me", notifications, results, tests.
  Triggers: sendmail, send to me, notify me, notify-me, message to self, banana.
---

# notify-me (send to self via sendmail)

This environment routes `sendmail` through a local `telegram-sendmail` service.
Messages arrive in Telegram instead of traditional email.

## The command

Always use the wrapped binary:

```bash
/run/wrappers/bin/sendmail << 'EOF'
Subject: Short descriptive title
From: grok@whiterun
To: lucas59356@gmail.com

Message body here.
Multi-line is fine.
Keep it mostly ASCII to avoid UTF-8 queue issues.
EOF
```

- `sendmail` in PATH is usually the wrapper too (`/run/wrappers/bin/sendmail`).
- The recipient (`To:` or argument) is typically ignored; everything goes to the pre-configured Telegram chat.
- On success the command prints `OK` (message enqueued for the daemon).

## Quick patterns

One-liner:

```bash
printf 'Subject: test\n\nbanana\n' | /run/wrappers/bin/sendmail
```

Heredoc (preferred for anything non-trivial).

## Configuration in this repo

- Service: `services.telegram-sendmail.enable = true` (riverwood + whiterun)
- Module: `nix/nodes/common/services/telegram_sendmail.nix`
- Credentials: `secrets/telegram_sendmail.env` (sops, provides MAIL_TELEGRAM_TOKEN + MAIL_TELEGRAM_CHAT)
- Other users of sendmail:
  - `nix/nodes/common/services/cloud-savegame.nix`
  - `bin/qrun`

Global identity (flake.nix):
- user: lucasew
- email: lucas59356@gmail.com

## When the user says "send ... to me" or asks to notify

1. Format a minimal valid message with a clear Subject.
2. Pipe or heredoc it to `/run/wrappers/bin/sendmail`.
3. Report success: "Sent via sendmail (delivered to Telegram)."

This skill is called "notify-me". Do not use `mail`, `mutt`, `msmtp` directly unless the user asks — sendmail (via the telegram-sendmail backend) is the established channel here.
