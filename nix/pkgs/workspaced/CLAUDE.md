# Workspaced Development

## Overview
User configs/dotfiles in `config/`. Settings in `settings.toml`. Templates use `{{ .Field }}` syntax.
- **⚠️ CRITICAL**: See `TEMPLATES.md` for complete template system documentation (5 types: static, simple, multi-file, index, .d.tmpl)

## Common Commands
- **Apply**: `workspaced apply` (builds Go code and applies configs)
- **Plan**: `workspaced plan` (dry-run, shows what would change)
- **Doctor**: `workspaced doctor` (check driver status)
  - Use `--verbose` flag to see full interface/provider paths

## Adding New Config Fields
When adding new config fields to `pkg/config/config.go`:
1. Add field to `GlobalConfig` struct with `toml:"field_name"` tag
2. Create corresponding struct (e.g., `FooConfig`) with fields and tags
3. **CRITICAL**: Add `Merge()` method to new struct (see `PaletteConfig.Merge()` as example)
4. **CRITICAL**: Call merge in `GlobalConfig.Merge()`: `result.Foo = result.Foo.Merge(other.Foo)`
5. Add config section to `settings.toml`
6. Templates access via `{{ .Foo.Field }}`

**⚠️ IMPORTANT - Merge Methods:**
- LoadConfig() creates hardcoded defaults, then loads settings.toml and merges
- Without implementing `Merge()` and calling it in `GlobalConfig.Merge()`, the merge doesn't happen
- Result: values from settings.toml are ignored, templates generate empty fields
- Symptom: code compiles OK, TOML is read, but `{{ .Field }}` returns empty string
- Always implement Merge() for structs nested in GlobalConfig!

## CLI & Architecture
- **Intention-based Structure**: Commands are grouped by user intent:
  - `input`: User interaction (`text`, `confirm`, `choose`, `menu`).
  - `open`: Resource launching (`webapp`, `terminal`, generic URLs/files).
  - `system`: Hardware and session state (`audio`, `brightness`, `power`, `screen`).
  - `state`: Dotfiles lifecycle (`apply`, `plan`, `sync`, `doctor`).
- **Local-First**: CLI binary executes hardware/system logic locally whenever possible. Daemon handles shared state, tray, watchers, and cross-client coordination (OSD IDs).
- **Module System**:
  - Located in `modules/`. Atomic, parametric, and strictly unique (no claim collisions).
  - Uses `module.toml` (deps), `defaults.toml` (base config), and `schema.json` (validation).
  - **Zero-Intermediate**: Files are processed in-memory and streamed directly to targets.
- **Lazy Processing**: `source.File` interface delays content reading/rendering until strictly needed.
- **Strict Config**: No lists in module configs. Deep merge with zero substitution policy between different modules.
- **Top-level Aliases**: `sync`, `apply`, `plan`, and `open` are mirrored at root for ergonomics.

## Driver System
- Drivers provide platform-specific implementations for various features (audio, clipboard, notifications, etc.)
- Each driver implements a provider interface with:
  - `ID()`: Unique slug (e.g., "audio_pulse")
  - `Name()`: User-friendly name (e.g., "PulseAudio")
  - `DefaultWeight()`: Priority (0-100)
  - `CheckCompatibility()`: Verify if driver can run
  - `New()`: Create instance
- Use `workspaced doctor` to see all drivers and their status
- Configure driver weights in `settings.toml` under `[driver.weights]`
