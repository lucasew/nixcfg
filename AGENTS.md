# AGENTS.md

## Mandates
- **Do NOT unify `mise.toml` files.** Keep them localized in `/`, `workspaced/`, and `infra/` to maintain context-specific tool versions and tasks.
- **Script usage**: Use `sdw` over `sd` to ensure the latest version from the dotfiles directory is used.

## Where To Find Things
- `nix/nodes/` -> NixOS machine configurations.
- `nix/overlay.nix` -> Defines `pkgs.custom.*`.
- `nix/pkgs/custom/` -> Sources for custom packages.
- `bin/` -> Categorized helper scripts and binaries.
- `nix/nodes/common/sops.nix` -> `sops-nix` managed secrets.
- `flake.nix` -> Global settings including user, email, IPs, and DE configuration.
- `workspaced/` -> User configs and dotfiles manager (see `workspaced/AGENTS.md` for development docs).

## Machine Context
- **riverwood**: Laptop, Intel CPU/GPU, ext4, Sway/i3.
- **whiterun**: Desktop, Ryzen 5600G, ZFS, Monitoring/Containers.
- **ravenrock**: Cloud (GCP), Turbo VM (currently unused).

## Common Commands
- **NixOS Apply**: `sudo nixos-rebuild switch --flake .`
- **Deploy**: `sd nix rrun .#deploy`
- **Build**: `nix build .#nixosConfigurations.MACHINE.config.system.build.toplevel`
- **Workspaced**: `workspaced home apply` / `workspaced home plan` / `workspaced driver doctor`
