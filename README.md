# Dotfiles and Nix/NixOS settings

[![built with nix](https://builtwithnix.org/badge.svg)](https://builtwithnix.org)

See [CLAUDE.md](./CLAUDE.md) (or [AGENTS.md](./AGENTS.md) if migrated) for AI agent conventions and instructions.

## Project Structure

- `nix/`: Nix and NixOS configurations.
  - `nix/nodes/`: Machine-specific NixOS configurations.
  - `nix/pkgs/`: Custom package definitions.
- `modules/`: `workspaced` modules for applying configs and `base16` colors to applications.
- `bin/`: Categorized user scripts and CLI utilities.
- `config/`: Raw dotfile configurations.
- `infra/`: Infrastructure definitions (e.g., Terraform for cloud machines).

The repository is organized so that it is not necessary to place the Nix files in the default locations like `/etc/nixos/configuration.nix` or `~/.config/nixpkgs/home.nix`.

## Graphical Environments

- **i3/Sway**: Daily driver, works nicely, playback buttons work when locked.
- **GNOME, XFCE, and KDE**: Not currently used, might be removed in the future.

## Machines

- **riverwood**: Main laptop (Acer A315-51-51SL), Intel CPU/GPU, 12GB RAM, 1TB SSD, dual booted with Windows 10, running Sway/i3 on ext4.
- **whiterun**: Battlestation Desktop, Ryzen 5600G, 32GB RAM, 1TB SSD + 2x1TB DVR HDDs, uses ZFS, serves monitoring/containers.
- **ravenrock**: Cloud machine (GCP Turbo VM), provisioned using Terraform from `infra/turbo/gcp.tf` (currently unused).

## Common Commands

- **NixOS Apply**: `sudo nixos-rebuild switch --flake .`
- **Deploy via sdw**: `sd nix rrun .#deploy`
- **Build Top-level**: `nix build .#nixosConfigurations.MACHINE.config.system.build.toplevel`
- **Workspaced Home Apply**: `workspaced home apply`
- **Workspaced Home Plan**: `workspaced home plan`

- licence
    - nothing special
    - don't blame me
    - have fun

- suggestions?
    - open a issue
    - let's learn together :smile:

- NixOS > Arch
    - change my mind
    - (yes, I have used arch btw for around 1 year, it's a good distro but NixOS is better for my workflow)
    - `nix-shell` rocks
    - the possibility of rollback at any time in a simple way, even if the distro fails to boot, is like magic
    - you can also replicate very precisely your configuration on another machine, but only if that is defined in Nix, imperative settings are left behind
    - not perfect but it's really easy to feel physical pain using something else or packaging software that tries to download stuff at build time
