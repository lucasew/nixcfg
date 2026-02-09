# CLAUDE.md

## Mandates
- **Do NOT unify `mise.toml` files.** Keep them localized in `/`, `nix/pkgs/workspaced/`, and `infra/` to maintain context-specific tool versions and tasks.
- **Script usage**: Use `sdw` over `sd` to ensure the latest version from the dotfiles directory is used.

## Architecture
- **NixOS/Home**: `nix/nodes/` (machine configs). Home-manager removed - configs migrated to workspaced.
- **Packages**: `nix/overlay.nix` defines `pkgs.custom.*`. Sources in `nix/pkgs/custom/`.
- **Scripts**: Categorized in `bin/`. Env init: `source bin/source_me`.
- **Secrets**: `sops-nix` managed in `nix/nodes/common/sops.nix`.
- **Global Settings**: `flake.nix` contains `global` attr (user, email, IPs, DE).
- **Workspaced**: User configs/dotfiles in `config/`. Settings in `settings.toml`. Templates use `{{ .Field }}` syntax.

## Machine Context
- **riverwood**: Laptop, Intel CPU/GPU, ext4, Sway/i3.
- **whiterun**: Desktop, Ryzen 5600G, ZFS, Monitoring/Containers.
- **ravenrock**: Cloud (GCP), Turbo VM (currently unused).

## Common Commands
- **Apply**: `sudo nixos-rebuild switch --flake .`
- **Apply (Go)**: `workspaced apply` (for Go packages in `nix/pkgs/workspaced`)
- **Deploy**: `sd nix rrun .#deploy`
- **Build**: `nix build .#nixosConfigurations.MACHINE.config.system.build.toplevel`

## Workspaced Development
When adding new config fields to `nix/pkgs/workspaced/pkg/config/config.go`:
1. Add field to `GlobalConfig` struct with `toml:"field_name"` tag
2. Create corresponding struct (e.g., `FooConfig`) with fields and tags
3. **CRITICAL**: Add `Merge()` method to new struct (see `PaletteConfig.Merge()` as example)
4. **CRITICAL**: Call merge in `GlobalConfig.Merge()`: `result.Foo = result.Foo.Merge(other.Foo)`
5. Add config section to `settings.toml`
6. Templates access via `{{ .Foo.Field }}`

**⚠️ IMPORTANTE - Merge Methods:**
- LoadConfig() cria defaults hardcoded, depois carrega settings.toml e faz merge
- Sem implementar `Merge()` e chamar no `GlobalConfig.Merge()`, o merge não acontece
- Resultado: valores do settings.toml são ignorados, templates geram campos vazios
- Sintoma: código compila OK, TOML é lido, mas `{{ .Field }}` retorna string vazia
- Sempre implementar Merge() para structs nested no GlobalConfig!
