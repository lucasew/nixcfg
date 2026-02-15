# base16-rofi Module

Generates rofi theme based on the base16 color palette.

## Features

- Complete rofi color scheme using base16 palette
- Supports normal, active, and urgent states
- Consistent look with other base16 modules

## Generated Files

- `~/.config/rofi/theme.rasi`

## Usage

Rofi will automatically use the generated theme. Launch rofi normally:

```bash
rofi -show drun
```

## Base16 Color Mapping

- **base00**: Background
- **base01**: Light background (selected items)
- **base04/base05**: Foreground colors
- **base08**: Red (urgent states)
- **base0D**: Blue (active states)

## Dependencies

Requires the `base16` module to be enabled with a valid color palette.
