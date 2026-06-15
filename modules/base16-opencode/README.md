# base16-opencode Module

Generates a theme for Opencode based on the base16 color palette.

## Features

- Creates a local JSON theme (`workspaced`) for Opencode.
- Integrates Opencode syntax and UI colors with the system-wide base16 palette.

## Generated Files

- `~/.config/opencode/themes/workspaced.json`

## Usage

Opencode will load the generated JSON theme. Configure Opencode to use the `workspaced` theme through its configuration or settings interface to see the base16 colors applied.

## Base16 Color Mapping

- Maps standard base16 colors to Opencode UI and editor token settings.

## Dependencies

Requires the `base16` module to be enabled with a valid color palette.
