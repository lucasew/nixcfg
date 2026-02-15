# base16-gtk Module

Generates GTK and Qt themes based on the base16 color palette.

## Features

- **GTK 2**: Full theme with all widget states
- **GTK 3/3.20**: Colorscheme overriding Adwaita-dark base
- **GTK 4**: Modern colorscheme for Adwaita with proper accent colors
- **Qt5/Qt6**: Color schemes for Qt5ct and Qt6ct

## Generated Files

### GTK
- `~/.local/share/themes/base16/gtk-2.0/gtkrc`
- `~/.local/share/themes/base16/gtk-3.0/gtk.css`
- `~/.local/share/themes/base16/gtk-3.20/gtk.css`
- `~/.local/share/themes/base16/gtk-4.0/gtk.css`

### Qt
- `~/.config/qt5ct/colors/base16.conf`
- `~/.config/qt6ct/colors/base16.conf`

## Usage

### GTK
The module is automatically applied when enabled. Set your GTK theme to "base16":

```bash
gsettings set org.gnome.desktop.interface gtk-theme base16
```

Or in your GTK settings files:
```ini
gtk-theme-name=base16
```

### Qt
For Qt applications, configure qt5ct/qt6ct to use the "base16" color scheme:
- Run `qt5ct` or `qt6ct`
- Go to "Appearance" tab
- Select "base16" in the "Color scheme" dropdown

## Base16 Color Mapping

- **base00**: Default Background
- **base01**: Lighter Background (status bars, line numbers)
- **base02**: Selection Background
- **base03**: Comments, Invisibles
- **base04**: Dark Foreground (status bars)
- **base05**: Default Foreground
- **base06**: Light Foreground
- **base07**: Light Background
- **base08**: Red (Error, Destructive)
- **base09**: Orange
- **base0A**: Yellow (Warning)
- **base0B**: Green (Success)
- **base0C**: Cyan
- **base0D**: Blue (Accent)
- **base0E**: Purple (Link visited)
- **base0F**: Brown

## Dependencies

Requires the `base16` module to be enabled with a valid color palette.
