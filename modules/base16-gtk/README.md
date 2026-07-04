# base16-gtk

GTK 2/3/4 theme plus qt5ct/qt6ct color schemes from `modules.base16.config`. Needs `base16` enabled. Palette notes: [../README.md](../README.md).

Writes:

- `~/.local/share/themes/base16/gtk-{2.0,3.0,3.20,4.0}/...`
- `~/.config/qt{5,6}ct/colors/base16.conf`

GTK: set theme name to `base16` (`gsettings set org.gnome.desktop.interface gtk-theme base16`, or `gtk-theme-name=base16` in the gtk settings file).

Qt: in `qt5ct` / `qt6ct` -> Appearance -> color scheme `base16`.
