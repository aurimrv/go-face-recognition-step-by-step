## Why

Markdown and Go static analysis findings (markdownlint, go vet, staticcheck via golangci-lint) are fixed ad hoc today, with no shared agent workflow for ingesting reports, grouping issues, driving OpenSpec changes, re-verifying, or recording deferrals. A project skill and `/opsx-static-fix` command make that loop repeatable for this teaching fork without touching recognition, models, or Docker native deps.

## What Changes

- Add project skill `.cursor/skills/static-analysis-fix/SKILL.md` that defines a spec-driven fix loop: ingest or run analyzers → triage/group → OpenSpec change (proposal, design, specs, tasks) → apply fixes task-by-task → re-run → archive when clean (or stop after two fix/re-run rounds).
- Add Cursor command `.cursor/commands/opsx-static-fix.md` (`/opsx-static-fix`) that invokes that skill.
- Analyzer policy for this repo: `markdownlint-cli2` only on project Markdown with `"**/*.md" "#openspec" "#.cursor" "#node_modules"`; `go vet`; `staticcheck` / `golangci-lint` using the existing `.golangci.yml` (v2)—do not create a new golangci config; Docker fallback when CGO packages need the project image.
- Fix policy: preserve behavior; prefer real fixes over suppressions; no blanket `nolint` or `markdownlint-disable`; every deferral MUST record finding id and rationale.
- Out of scope for the skill’s edits: recognition pipeline, models, Docker native dependencies.

## Capabilities

### New Capabilities

- `static-analysis-fix`: Agent workflow and policies for resolving Markdown/Go static analysis findings via OpenSpec (ingest/run, triage, change artifacts, apply, re-run, archive/stop, fix and deferral rules, analyzer commands, and hard exclusions).

### Modified Capabilities

- (none)

## Impact

- New files under `.cursor/skills/static-analysis-fix/` and `.cursor/commands/opsx-static-fix.md`.
- No application runtime changes as part of *this* change; later runs of the skill may edit project Markdown/Go under the skill’s policies.
- Relies on existing OpenSpec CLI (`propose` / `apply` / `archive` patterns) and existing `.golangci.yml`.
- Does not modify `Dockerfile` native deps, `models/`, or face-recognition use cases as part of delivering the skill/command.
