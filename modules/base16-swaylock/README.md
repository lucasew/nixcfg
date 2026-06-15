# base16-swaylock Module

Generates a lock screen script and systemd service tailored with the base16 color palette.

## Features

- Configures Swaylock with base16 colors for all visual states (typing, clearing, verifying, wrong password).
- Provides a `lock-screen` executable in `~/.local/bin`.
- Sets up an `xss-lock` user service to automatically trigger swaylock on system sleep or inactivity.

## Generated Files

- `~/.local/bin/lock-screen`
- `~/.config/systemd/user/xss-lock.service`

## Usage

You can trigger the lock screen manually by executing:

```bash
lock-screen
```

It is also automatically triggered by the system via `xss-lock` when locking is requested. Ensure the systemd user service `xss-lock.service` is enabled and started.

## Base16 Color Mapping

- **base00**: Inside/ring background colors
- **base01**: Line colors
- **base04**: Text colors
- **base08**: Error/wrong password indicators
- **base0A**: Warning indicators
- **base0B**: Typing indicators
- **base0D**: Verifying indicators
- **base0E**: Clearing indicators

## Dependencies

Requires the `base16` module to be enabled with a valid color palette.
