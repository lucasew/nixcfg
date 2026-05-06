# AGENTS.md

## Mandates
- **Do NOT unify `mise.toml` files.** Keep them localized in `/`, `workspaced/`, and `infra/` to maintain context-specific tool versions and tasks.
- **Script usage**: Use `sdw` over `sd` to ensure the latest version from the dotfiles directory is used.

## Where To Find Things
- `nix/` -> Nix/NixOS configs
- `nix/nodes/` -> Machine-specific configs and non-NixOS device scripts like Android
- `nix/pkgs/custom/` -> Custom package sources (defined in `nix/overlay.nix`)
- `modules/` -> Workspaced modules for configs and base16 colors
- `bin/` -> CLI scripts (Env init: `source bin/source_me`)
- `config/` -> Raw dotfiles
- `infra/` -> Terraform definitions
- `flake.nix` -> Global settings (user, email, IPs, DE)
- `nix/nodes/common/sops.nix` -> Secrets (sops-nix managed)

## Machine Context
- **riverwood**: Laptop, Intel CPU/GPU, ext4, Sway/i3.
- **whiterun**: Desktop, Ryzen 5600G, ZFS, Monitoring/Containers.
- **ravenrock**: Cloud (GCP), Turbo VM (currently unused).

## Common Commands
- **NixOS Apply**: `sudo nixos-rebuild switch --flake .`
- **Deploy**: `sd nix rrun .#deploy`
- **Build**: `nix build .#nixosConfigurations.MACHINE.config.system.build.toplevel`
- **Workspaced**: `workspaced home apply` / `workspaced home plan` / `workspaced driver doctor`
