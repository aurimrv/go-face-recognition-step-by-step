## Why

Project Markdown (README, DEVCONTAINER_SETUP, teaching interaction dumps) and Go packages report 249 markdownlint issues and 6 golangci-lint issues (errcheck, ST1005). Clearing them improves teaching-repo hygiene under the static-analysis-fix workflow without changing recognition behavior, models, or Docker native deps.

## What Changes

- Fix Go errcheck on `Close` in `drawer.go` and `load_persons.go`; lowercase `processJob` error strings (ST1005).
- Fix Markdown structure and line-length issues in `README.md`, `DEVCONTAINER_SETUP.md`, `01-TERMINAL_INTERACTION.md`, and `01-COPILOT_INTERACTION.md` (real edits: wrap lines, blank lines around fences/headings, ordered-list style, fence languages, table pipe spacing, etc.).
- Defer any finding that cannot be fixed without harming transcript fidelity or violating hard exclusions, with id + rationale.

## Capabilities

### New Capabilities

- `lint-hygiene`: Requirements for keeping project Markdown and Go analyzer-clean under the repo’s configured markdownlint-cli2 scope and existing `.golangci.yml`.

### Modified Capabilities

- (none)

## Impact

- Touches project root Markdown and `cmd/main.go`, `internal/entity/drawer.go`, `internal/usecases/load_persons.go`.
- Does not modify recognition algorithms, models, or Dockerfile native packages.
- Uses existing `.golangci.yml` only.
