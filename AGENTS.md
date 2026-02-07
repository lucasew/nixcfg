# AGENTS.md

This file serves as the definitive guide for AI agents working on this repository. It consolidates project conventions, architectural decisions, and operational mandates.

## Mandates

- **Do NOT unify `mise.toml` files.** Keep them localized in `/`, `nix/pkgs/workspaced/`, and `infra/` to maintain context-specific tool versions and tasks.
- **Script usage**: Use `sdw` over `sd` to ensure the latest version from the dotfiles directory is used.
- **Error Handling**: All unexpected errors must be funneled through `logging.ReportError` (in Go) or equivalent centralized reporting. Silent failures are forbidden.
- **Naming Conventions**: Use self-documenting names. Avoid cryptic abbreviations.
- **Directory Structure**: Group by domain and responsibility. Sparse directories should be merged unless they represent a distinct domain concept.

## Architecture

### NixOS / Home Manager
- **Nodes**: `nix/nodes/` contains machine-specific configurations.
- **Homes**: `nix/homes/` contains user-specific configurations.
- **Packages**: `nix/overlay.nix` defines `pkgs.custom.*`. Sources reside in `nix/pkgs/custom/`.
- **Secrets**: Managed via `sops-nix` in `nix/nodes/common/sops.nix`.
- **Global Settings**: `flake.nix` contains the `global` attribute (user, email, IPs, DE).

### Go Codebase (`nix/pkgs/workspaced`)
- **Domain Packages**: Core business logic should reside in top-level packages under `pkg/` (e.g., `pkg/apply`, `pkg/config`).
- **Drivers**: `pkg/drivers/` contains adapters for external systems or hardware (e.g., `backup`, `git`, `wm`).
- **Commands**: `cmd/workspaced/` contains the CLI entry points, organized by subcommand.

### Scripts
- **Location**: Categorized in `bin/`.
- **Environment**: Initialize with `source bin/source_me`.

## Machine Context
- **riverwood**: Laptop, Intel CPU/GPU, ext4, Sway/i3.
- **whiterun**: Desktop, Ryzen 5600G, ZFS, Monitoring/Containers.
- **ravenrock**: Cloud (GCP), Turbo VM (currently unused).

## Common Commands
- **Apply (Nix)**: `sudo nixos-rebuild switch --flake .` or `home-manager switch --flake .#main`
- **Apply (Go)**: `workspaced apply` (executes configuration application logic).
- **Deploy**: `sd nix rrun .#deploy`
- **Build**: `nix build .#nixosConfigurations.MACHINE.config.system.build.toplevel`

## Refactoring Guidelines
- **Structure**: Prefer domain-driven packaging over layer-driven (e.g., `pkg/apply` instead of `pkg/drivers/apply` if `apply` is a core domain).
- **Testing**: Ensure critical logic (like `Plan` and `Execute` in `pkg/apply`) is testable by abstracting side-effects (e.g., FileSystem interface).
