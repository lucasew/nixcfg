# base16-tmux Module

Generates tmux theme based on the base16 color palette.

## Features

- Complete tmux color scheme using base16 palette
- Modular design using `tmux.conf.d/` directory
- Doesn't interfere with keybindings or other tmux configs
- Colors for:
  - Pane borders (active and inactive)
  - Status bar
  - Window status (current and inactive)
  - Messages and command mode
  - Copy mode
  - Clock mode

## Generated Files

- `~/.config/tmux/tmux.conf` (concatenated from all `.d.tmpl` sources)

## Usage

The theme is automatically loaded by tmux. To reload after changes:

```bash
tmux source-file ~/.config/tmux/tmux.conf
```

Or use the keybinding:
```
Prefix + r
```

## Base16 Color Mapping

- **base00**: Status bar background
- **base01**: Inactive pane borders, window status background
- **base02**: Current window background, copy mode background
- **base04/base05**: Foreground colors
- **base0D**: Active borders, current window, messages (blue accent)

## Dependencies

Requires the `base16` module to be enabled with a valid color palette.

## Structure

This module uses the **workspaced `.d.tmpl` pattern**:
- `config/.config/tmux/tmux.conf.d.tmpl/10-base.conf` → bindings, settings
- `modules/base16-tmux/home/.config/tmux/tmux.conf.d.tmpl/50-base16-theme.conf.tmpl` → colors only
- Workspaced concatenates all `.d.tmpl/**/*.conf` files into final `~/.config/tmux/tmux.conf`

The numbered prefixes (10-, 50-) control the order of concatenation.
