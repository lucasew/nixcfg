# base16-sway Module

Generates a theme configuration for Sway window manager and Waybar based on the base16 color palette.

## Features

- Provides client configuration for Sway defining border and background colors.
- Generates a CSS variables file for Waybar, exporting base16 colors.
- Keeps Sway and Waybar styling consistent with the rest of the system.

## Generated Files

- `~/.config/sway/config.d/base16`
- `~/.config/waybar/base16.css`

## Usage

**Sway:** Ensure your Sway configuration (`~/.config/sway/config`) includes the `config.d` directory so that `base16` colors are loaded:

```swayconfig
include /home/$USER/.config/sway/config.d/*
```

**Waybar:** Import the base16 CSS variables in your `~/.config/waybar/style.css`:

```css
@import "base16.css";
```

Restart Sway (`swaymsg reload`) and Waybar to apply the changes.

## Base16 Color Mapping

- **base00-base02**: Backgrounds and inactive borders.
- **base0D**: Active window borders (blue accent).
- **base08**: Urgent borders (red).
- CSS colors map directly to the base16 variables for Waybar.

## Dependencies

Requires the `base16` module to be enabled with a valid color palette.
