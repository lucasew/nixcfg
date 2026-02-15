# base16-dunst Module

Generates dunst notification daemon theme based on the base16 color palette.

## Features

- Notification colors based on base16 palette
- Different colors for urgency levels:
  - **Critical**: Red (base08) background
  - **Normal**: Blue (base0D) background
  - **Low**: Cyan (base0C) background

## Generated Files

- `~/.config/dunst/dunstrc`

## Usage

Restart dunst to apply changes:

```bash
killall dunst
# Dunst will auto-restart on next notification
```

Or manually:
```bash
dunst &
```

## Base16 Color Mapping

- **base06**: Foreground text
- **base08**: Critical urgency background (red)
- **base0C**: Low urgency background (cyan)
- **base0D**: Normal urgency background (blue)

## Dependencies

Requires the `base16` module to be enabled with a valid color palette.
