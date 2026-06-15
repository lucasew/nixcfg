# base16-shell Module

Generates shell prompt and UI colors based on the base16 color palette.

## Features

- Provides a `.bashrc.d` script that exports terminal colors.
- Integrates terminal colors matching the system-wide base16 palette.
- Loaded early to ensure subsequent bash prompt scripts or aliases have access to correct UI colors.

## Generated Files

- `~/.bashrc.d/00-ui-colors.sh`

## Usage

The module uses the workspaced `.d.tmpl` pattern to insert the script into your `.bashrc.d` directory. Make sure your `~/.bashrc` is configured to source files from `~/.bashrc.d/`.

When you open a new shell session, the base16 UI colors are automatically applied.

## Base16 Color Mapping

- Follows standard base16 terminal color mappings (e.g., base00 to base0F).

## Dependencies

Requires the `base16` module to be enabled with a valid color palette.
