# ❄️ nixcfg

[![Built with Nix](https://img.shields.io/badge/Built_with-Nix-5277C3?logo=nixos&logoColor=white)](https://builtwithnix.org)
[![CI Status](https://github.com/lucasew/nixcfg/workflows/autorelease/badge.svg)](https://github.com/lucasew/nixcfg/actions)
[![License: Unlicense](https://img.shields.io/badge/license-Unlicense-blue.svg)](http://unlicense.org/)

> My personal declarative infrastructure and dotfiles, powered by NixOS, Home Manager, and Flakes.

## 🌟 Highlights

- **Flake-based**: Fully reproducible configuration using Nix Flakes.
- **Multi-host**: Manages a diverse fleet of devices including laptops, desktops, cloud servers, and Android phones.
- **Secret Management**: Encrypted secrets stored safely in the repo using `sops-nix` and `age`.
- **Ephemeral Root**: Implements "erase your darlings" (impermanence) on supported hosts for a pristine state on every boot.
- **Styling**: Unified system-wide theming managed by `stylix` and `nix-colors`.
- **Dev Environments**: Project-specific environments managed with `direnv` and `mise`.

## 🖥️ Inventory

| Host | Type | Description |
| :--- | :--- | :--- |
| **Riverwood** | Laptop | Acer A315, 12GB RAM, Dual boot (Win10). Main mobile driver. |
| **Whiterun** | Desktop | Ryzen 5600G, 32GB RAM, 1TB SSD + 2TB HDD. Home Server & Battlestation. |
| **Ravenrock** | Cloud | GCP instance managed via Terraform. |
| **AtomicPi** | SBC | Low-power x86 board for lightweight services. |
| **Recovery** | Image | Custom NixOS recovery ISO. |
| **Phone** | Android | Managed environment via `nix-on-droid`. |

## 📂 Structure

- `flake.nix`: The entry point for the configuration.
- `nix/nodes/`: NixOS configurations for each machine.
- `nix/homes/`: Home Manager configurations (dotfiles).
- `nix/containers/`: OCI container definitions.
- `infra/`: Terraform configurations for cloud infrastructure.
- `secrets/`: Encrypted secrets.
- `assets/`: Wallpapers and other static assets.

## 🚀 Getting Started

### Development Environment

This repository uses `direnv` to automatically load the development environment.

```bash
direnv allow
```

This will set up `nix`, `sops`, and other necessary tools.

### Deployment

To deploy changes to the main hosts (riverwood/whiterun), use the custom deploy script:

```bash
nix run .#deploy
```

This script handles building the configurations, copying closures to the target machines, and switching to the new generation.

## 🛠️ Built With

- **[NixOS](https://nixos.org/)**: The operating system.
- **[Home Manager](https://github.com/nix-community/home-manager)**: User environment management.
- **[Sops-Nix](https://github.com/Mic92/sops-nix)**: Secrets management.
- **[Stylix](https://github.com/danth/stylix)**: System-wide styling.
- **[Nix-on-Droid](https://github.com/t184256/nix-on-droid)**: Nix on Android.

## 💭 Philosophy

> NixOS > Arch. Change my mind.

- **Reproducibility**: If it breaks, I can rollback. If I need a new machine, I can clone the config.
- **Declarative**: The state of the system is defined in code, not by running imperative commands.
- **Flexibility**: `nix-shell` makes trying new tools painless.
- **Magic**: The possibility of rollback at any time in a simple way, even if the distro fails to boot.
- **Sanity**: Not perfect, but it's really easy to feel physical pain using something else or packaging software that tries to download stuff at build time.

## 📄 License

> Nothing special. Don't blame me. Have fun.

This project is effectively unincensed. Feel free to copy, modify, and distribute.

---
*Maintained by [Lucasew](https://github.com/lucasew).*
