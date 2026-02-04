# CLAUDE.md

## Mandates
- **Do NOT unify `mise.toml` files.** Keep them localized in `/`, `nix/pkgs/workspaced/`, and `infra/` to maintain context-specific tool versions and tasks.
- **Script usage**: Use `sdw` over `sd` to ensure the latest version from the dotfiles directory is used.

## Architecture
- **NixOS/Home**: `nix/nodes/` (machine configs), `nix/homes/` (user configs).
- **Packages**: `nix/overlay.nix` defines `pkgs.custom.*`. Sources in `nix/pkgs/custom/`.
- **Scripts**: Categorized in `bin/`. Env init: `source bin/source_me`.
- **Secrets**: `sops-nix` managed in `nix/nodes/common/sops.nix`.
- **Global Settings**: `flake.nix` contains `global` attr (user, email, IPs, DE).

## Machine Context
- **riverwood**: Laptop, Intel CPU/GPU, ext4, Sway/i3.
- **whiterun**: Desktop, Ryzen 5600G, ZFS, Monitoring/Containers.
- **ravenrock**: Cloud (GCP), Turbo VM (currently unused).

## Common Commands
- **Apply**: `sudo nixos-rebuild switch --flake .` or `home-manager switch --flake .#main`
- **Apply (Go)**: `workspaced apply` (for Go packages in `nix/pkgs/workspaced`)
- **Deploy**: `sd nix rrun .#deploy`
- **Build**: `nix build .#nixosConfigurations.MACHINE.config.system.build.toplevel`
