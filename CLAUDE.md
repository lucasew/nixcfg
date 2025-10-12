# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

This is a NixOS configuration repository (dotfiles) using Nix flakes and home-manager. It manages multiple machines (riverwood, whiterun, ravenrock, atomicpi, recovery) and provides declarative system and user configurations.

## Architecture

### Repository Structure

- `flake.nix` - Main flake definition with inputs, outputs, and machine configurations
- `nix/` - Core Nix configurations
  - `nix/nodes/` - NixOS system configurations (per-machine)
    - `nix/nodes/bootstrap/` - Base system configuration imported by all nodes
    - `nix/nodes/common/` - Shared services and modules
    - `nix/nodes/gui-common/` - GUI-related configurations (desktop environments, audio, steam)
    - `nix/nodes/{riverwood,whiterun,ravenrock,atomicpi,recovery}/` - Machine-specific configs
  - `nix/homes/` - home-manager configurations
    - `nix/homes/base/` - Base home configuration
    - `nix/homes/main/` - Primary user home configuration
  - `nix/nixOnDroid/` - Android/termux configurations via nix-on-droid
  - `nix/pkgs/` - Custom package definitions and overrides
  - `nix/lib/` - Custom Nix library functions
  - `nix/overlay.nix` - Package overlay defining custom packages and overrides
- `bin/` - Shell scripts and utilities organized by category
  - `bin/prelude/` - Shell initialization scripts sourced via `source_me`
  - `bin/source_me` - Script that sources prelude files to initialize environment
  - `bin/d/` - Deployment scripts
  - `bin/g/` - Git shortcuts
  - `bin/misc/` - Miscellaneous utilities
- `default.nix` - Flake compatibility layer for legacy nix commands

### Key Architectural Patterns

**Node-based configuration**: Each machine is a "node" with its own directory under `nix/nodes/`. Nodes import shared modules from `bootstrap/`, `common/`, and `gui-common/`.

**Overlay system**: The `nix/overlay.nix` file extends nixpkgs with:
- `custom.*` - Custom package builds (neovim, vscode, firefox, rofi, etc.)
- `unstable.*` - Packages from nixpkgs-unstable
- Helper functions and package wrappers

**Script directory system**: Scripts in `bin/` are organized by category and accessed via the `sd` command. The `sdw` (script-directory-wrapper) is a NixOS-installed wrapper that locates the dotfiles directory (either at `~/.dotfiles/bin` or `/etc/.dotfiles/bin`) to ensure the latest script versions are used, even when the environment is built from a Nix store path.

**Flake inputs**: External dependencies are tracked as flake inputs (home-manager, stylix, nixos-hardware, etc.) and referenced throughout the configuration.

## Common Commands

### Building Configurations

```bash
# Build NixOS system configuration
nix build .#nixosConfigurations.riverwood.config.system.build.toplevel
nix build .#nixosConfigurations.whiterun.config.system.build.toplevel

# Build home-manager configuration
nix build .#homeConfigurations.main.activationPackage

# Build everything (release)
nix build .#release
```

### Deployment

```bash
# Deploy to machines using the custom deploy script
sd nix rrun .#deploy
```

### Development Shell

```bash
# Enter development shell with tools
nix develop

# This provides custom environment via bin/source_me
```

### Applying Changes

```bash
# Switch NixOS configuration on the current machine
sudo nixos-rebuild switch --flake .

# Switch home-manager configuration
home-manager switch --flake .#main

# Test configuration without switching
sudo nixos-rebuild test --flake .
```

### Script Directory Commands

Scripts are accessed via `sd CATEGORY SCRIPT [ARGS]` or `sdw CATEGORY SCRIPT [ARGS]`. The `sdw` command is preferred as it finds the latest version in the dotfiles directory.

```bash
# Git shortcuts (bin/g/*)
sd g s       # git status
sd g a       # git add
sd g c       # git commit
sd g d       # git diff

# Deployment utilities (bin/d/*)
sd d deploy riverwood,whiterun  # Deploy to machines
sd d sw PACKAGE -- COMMAND      # Quick nix-shell wrapper
sd d root                       # Get dotfiles root path

# Check if on specific machine (bin/is/*)
sd is riverwood  # Returns 0 if on riverwood
sd is whiterun   # Returns 0 if on whiterun
```

## Important Configuration Details

### Global Settings

The `global` attribute in `flake.nix` defines shared settings:
- Username: `lucasew`
- Email: `lucas59356@gmail.com`
- Node IPs (tailscale and zerotier)
- Selected desktop environment: `i3`

### Machine Characteristics

- **riverwood**: Laptop (Acer A315-51-51SL), Intel CPU/GPU, 12GB RAM, ZFS, Sway WM enabled
- **whiterun**: Desktop (Ryzen 5600G), 32GB RAM, ZFS, monitoring stack, container services
- **ravenrock**: Cloud server (GCP), services like vaultwarden, nginx, postgres

### Desktop Environments

The repository supports multiple DEs (defined in `nix/nodes/gui-common/gui-variants/`):
- i3 (default, primary use)
- Sway (enabled on riverwood)
- GNOME, XFCE, Plasma5 (available but not actively used)

### Custom Packages

Custom package definitions are in `nix/pkgs/custom/`:
- neovim with custom configuration
- vscode with extensions organized by category (common, programming, kb)
- rofi with custom theme
- firefox with specific extensions
- Custom polybar, emacs, pidgin, retroarch, etc.

All custom packages are accessed via `pkgs.custom.PACKAGE_NAME`.

### Library Functions

Custom Nix library functions in `nix/lib/`:
- `importAllIn.nix` - Import all Nix files in a directory
- `listModules.nix` - List module files
- `patchNixpkgs.nix` - Apply patches to nixpkgs
- `buildDockerEnv.nix` - Build Docker environments
- `image2color.nix`, `jpg2png.nix` - Image utilities
- `unpack.nix`, `unpackRecursive.nix` - Archive extraction helpers

## Editing Patterns

### Adding a New NixOS Module

1. Create module file in appropriate location (`nix/nodes/common/`, `nix/nodes/bootstrap/`, or machine-specific directory)
2. Import it in the relevant `default.nix` or machine configuration
3. Use standard NixOS module structure with `{ config, pkgs, lib, ... }:`

### Adding a New Machine

1. Create directory under `nix/nodes/MACHINE_NAME/`
2. Create `default.nix` importing `../bootstrap` and other shared modules
3. Add hardware-configuration.nix (generate with `nixos-generate-config`)
4. Register in `flake.nix` under `nixosConfigurations`

### Modifying Packages

1. For package overrides: edit `nix/overlay.nix`
2. For custom packages: add/edit in `nix/pkgs/custom/`
3. Packages are accessed via `pkgs.custom.PACKAGE_NAME`

### Adding Scripts

1. Place script in appropriate `bin/CATEGORY/` directory
2. Scripts are automatically available via `sd CATEGORY SCRIPT`
3. The `sdw` wrapper ensures the latest script versions from the dotfiles directory are used

## Stylix Integration

The configuration uses stylix (home-manager and NixOS modules) for unified theming. The color scheme is defined in `flake.nix` outputs as `self.colors` (based on nix-colors "darkviolet" scheme).

## SOPS for Secrets

Secrets management uses sops-nix. Secret files are encrypted and referenced in configurations. The bootstrap includes sops configuration in `nix/nodes/common/sops.nix`.
