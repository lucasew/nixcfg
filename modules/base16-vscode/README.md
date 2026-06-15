# base16-vscode Module

Generates a VS Code extension that applies the base16 color palette.

## Features

- Creates a local VS Code theme extension (`workspaced.base16-theme`).
- Automatically syncs the editor colors to the system-wide base16 palette.
- Supports flatpak installations of VS Code (`com.visualstudio.code`).

## Generated Files

- `~/.vscode/extensions/workspaced.base16-theme/package.json`
- `~/.var/app/com.visualstudio.code/data/vscode/extensions/workspaced.base16-theme/package.json`
- (And related theme JSON files configured via the module index)

## Usage

VS Code will automatically detect the local extension. Go to **Preferences: Color Theme** (Ctrl+K Ctrl+T) and select the base16 theme.

If VS Code is open when you apply changes, you may need to reload the window (Ctrl+Shift+P > "Developer: Reload Window") for the new colors to take effect.

## Dependencies

Requires the `base16` module to be enabled with a valid color palette.
