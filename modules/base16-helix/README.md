# base16-helix Module

Generates a theme for the Helix editor based on the base16 color palette.

## Features

- Complete Helix color scheme using the base16 palette.
- Defines UI colors and syntax highlighting.

## Generated Files

- `~/.config/helix/themes/base16.toml`

## Usage

Helix needs to be configured to use the generated theme. Add or update your `~/.config/helix/config.toml`:

```toml
theme = "base16"
```

Then restart Helix or use the `:theme base16` command within the editor to apply it immediately.

## Base16 Color Mapping

- Follows the standard base16 mapping conventions for syntax highlighting and UI elements in Helix.

## Dependencies

Requires the `base16` module to be enabled with a valid color palette.
