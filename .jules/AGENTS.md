# Project Conventions and Memory

## Operational Memory
- `bin/` -> Categorized scripts and utilities (e.g. `bin/ai`, `bin/misc`, `bin/svc`).
- `nix/nodes/` -> NixOS machine configurations.
- `nix/pkgs/custom/` -> Custom package definitions and sources.
- `config/` -> User configurations and dotfiles managed by workspaced.
- `infra/` -> Infrastructure provisioning and deployments (e.g. GCP, GKE).
- `.jules/` -> Agent memory journals and rules.

## Error Handling
- The project MUST have a single, centralized error-reporting function. All code paths that handle unexpected errors MUST funnel through `src/utils/errorReporting.ts` (or equivalent).
- No silent failures: Every `catch` block MUST report errors centrally.
