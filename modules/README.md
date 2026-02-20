# Workspaced Modules Guide

## Creating Base16 Modules

Base16 modules apply the color palette defined in `settings.toml` to various applications.

### Basic Structure

```
modules/base16-{app}/
├── module.toml         # Module metadata
├── defaults.toml       # Driver preferences (optional)
├── README.md          # Documentation
└── home/
    └── .config/{app}/
        └── config.ext.tmpl  # Template file
```

### module.toml

```toml
[module]
name = "base16-{app}"
requires = ["base16"]
```

### Accessing Base16 Colors in Templates

```go
{{- $base16 := index .Modules "base16" }}
color: #{{ $base16.base00 }}
```

## The `.d.tmpl` Pattern

For configs that need to be split across multiple sources (main config + theme), use the `.d.tmpl` pattern:

### How It Works

1. Create a directory with `.d.tmpl` suffix
2. Add multiple files (can be from different modules)
3. Workspaced concatenates them into a single output file

### Example: tmux

```
config/.config/tmux/tmux.conf.d.tmpl/
  └── 10-base.conf                    # Main config (bindings, etc)

modules/base16-tmux/home/.config/tmux/tmux.conf.d.tmpl/
  └── 50-base16-theme.conf.tmpl       # Colors from module

→ Generates: ~/.config/tmux/tmux.conf (single file)
```

### Ordering

Files are concatenated in **alphabetical order**. Use numeric prefixes to control order:
- `00-` - Critical early configs
- `10-` - Base configuration
- `50-` - Themes and styling
- `90-` - Late configs

### When to Use `.d.tmpl`

✅ **Use when:**
- Config has functional parts (bindings) + theme parts (colors)
- Multiple modules need to contribute to same file
- Want separation of concerns

❌ **Don't use when:**
- Single module owns entire config
- Config is theme-only (just use `.tmpl`)

## Base16 Color Reference

| Color | Purpose | Example Usage |
|-------|---------|---------------|
| base00 | Background | Main background |
| base01 | Lighter background | Line numbers, status bars |
| base02 | Selection background | Selected items |
| base03 | Comments, invisibles | Subtle elements |
| base04 | Dark foreground | Status bar text |
| base05 | Default foreground | Main text |
| base06 | Light foreground | Highlighted text |
| base07 | Light background | Unused in dark themes |
| base08 | Red | Error, destructive, urgent |
| base09 | Orange | Warnings, modified |
| base0A | Yellow | Warnings |
| base0B | Green | Success, additions |
| base0C | Cyan | Info, low urgency |
| base0D | Blue | Accent, links, active |
| base0E | Purple | Special, visited links |
| base0F | Brown | Deprecated |

## Examples

### Simple Theme Module (rofi, dunst)

Single template file, module owns entire config:
```
modules/base16-rofi/home/.config/rofi/theme.rasi.tmpl
```

### Complex Module with `.d.tmpl` (tmux, bashrc)

Multiple sources concatenated:
```
config/.config/tmux/tmux.conf.d.tmpl/10-base.conf
modules/base16-tmux/...tmux.conf.d.tmpl/50-theme.conf.tmpl
→ ~/.config/tmux/tmux.conf
```

### Multi-File Module (gtk)

Creates files in multiple locations:
```
modules/base16-gtk/home/
  ├── .local/share/themes/base16/gtk-2.0/gtkrc.tmpl
  ├── .local/share/themes/base16/gtk-3.0/gtk.css.tmpl
  └── .config/qt5ct/colors/base16.conf.tmpl
```

## Migration Checklist

When converting a legacy config to a module:

1. ✅ Create module structure
2. ✅ Move template from `config/` to `modules/{name}/home/`
3. ✅ Remove old template from `config/`
4. ✅ Add `[modules.{name}]` to `settings.toml`
5. ✅ Test with `workspaced home plan`
6. ✅ Apply with `workspaced home apply`

## Existing Base16 Modules

- `base16` - Core palette definition
- `base16-shell` - Terminal colors
- `base16-sway` - Sway WM + Waybar
- `base16-helix` - Helix editor
- `base16-vscode` - VS Code
- `base16-gtk` - GTK 2/3/4 + Qt 5/6
- `base16-rofi` - Rofi launcher
- `base16-dunst` - Dunst notifications
- `base16-tmux` - Tmux terminal multiplexer
