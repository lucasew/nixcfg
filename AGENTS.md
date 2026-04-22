# AGENTS.md

## Mandates
- **Do NOT unify `mise.toml` files.** Keep them localized in `/`, `workspaced/`, and `infra/` to maintain context-specific tool versions and tasks.
- **Script usage**: Use `sdw` over `sd` to ensure the latest version from the dotfiles directory is used.

## Where To Find Things
- `nix/nodes/` -> machine configs (NixOS/Home configs, migrated to workspaced)
- `nix/overlay.nix` -> defines `pkgs.custom.*` packages
- `nix/pkgs/custom/` -> sources for custom packages
- `nix/nodes/common/sops.nix` -> `sops-nix` managed secrets
- `bin/` -> categorized scripts (run `source bin/source_me` for env init)
- `flake.nix` -> global settings (user, email, IPs, DE)
- `workspaced/` -> user configs/dotfiles manager (see `workspaced/AGENTS.md` for development docs)
- `modules/` -> base16 themes and templates for various applications
- `config/` -> raw config dotfiles templates
- `infra/` -> terraform definitions for remote instances

## Machine Context
- **riverwood**: Laptop, Intel CPU/GPU, ext4, Sway/i3.
- **whiterun**: Desktop, Ryzen 5600G, ZFS, Monitoring/Containers.
- **ravenrock**: Cloud (GCP), Turbo VM (currently unused).

## Common Commands
- **NixOS Apply**: `sudo nixos-rebuild switch --flake .`
- **Deploy**: `sd nix rrun .#deploy`
- **Build**: `nix build .#nixosConfigurations.MACHINE.config.system.build.toplevel`
- **Workspaced**: `workspaced home apply` / `workspaced home plan` / `workspaced driver doctor`
