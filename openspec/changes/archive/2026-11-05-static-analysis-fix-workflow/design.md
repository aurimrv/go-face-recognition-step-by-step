## Context

The repo already uses OpenSpec (`/opsx-propose`, `/opsx-apply`, `/opsx-archive`) and ships `.golangci.yml` (version `"2"`) enabling errcheck, govet, ineffassign, staticcheck, and unused. Sample markdownlint output lives in teaching notes (`01-TERMINAL_INTERACTION.md`). Agents need a single project skill plus a Cursor command that runs a bounded, policy-heavy fix loop without inventing new golangci config or touching recognition/Docker native stacks.

## Goals / Non-Goals

**Goals:**

- Ship `.cursor/skills/static-analysis-fix/SKILL.md` and `.cursor/commands/opsx-static-fix.md` mirroring existing opsx skill/command pairing.
- Encode analyzer invocation, triage → OpenSpec change → apply → re-run → archive/stop, and fix/deferral policies in the skill so agents execute them consistently.
- Reuse existing `.golangci.yml`; support Docker fallback for CGO/`go-face` analysis when the host cannot run the linters.

**Non-Goals:**

- Fixing current lint debt as part of *this* change (that is a later `/opsx-static-fix` run).
- Replacing OpenSpec; changing recognition pipeline, models, or Dockerfile native dependencies.
- Creating a new `.golangci.yml` or enabling extra linters beyond the existing file.
- Building a custom Go/Markdown parser UI or CI workflow YAML (unless already present—out of scope here).

## Decisions

### 1. Skill + thin command (same pattern as other opsx workflows)

**Choice:** Command file points agents at the skill; the skill holds the full procedure. Skill lives under `.cursor/skills/static-analysis-fix/` (project-scoped).

**Alternatives:** Command-only markdown (duplicates poorly); personal skill in `~/.cursor/skills` (not shared with the teaching fork).

**Rationale:** Matches `openspec-apply-change` / `opsx-apply` and keeps the workflow in-repo for students.

### 2. OpenSpec as the fix vehicle

**Choice:** After triage, the agent creates an OpenSpec change (proposal, design, specs, tasks), applies task-by-task, re-runs analyzers, and archives when clean—or stops after two fix+re-run rounds and reports leftovers.

**Alternatives:** Fix files directly without OpenSpec (faster but no audit trail); one mega-commit without tasks (harder to review).

**Rationale:** Aligns with the course’s OpenSpec teaching path and makes deferrals/tasks visible.

### 3. Analyzer matrix (locked)

| Tool | Invocation / notes |
|---|---|
| markdownlint-cli2 | `markdownlint-cli2 "**/*.md" "#openspec" "#.cursor" "#node_modules"` only |
| go vet | Standard package vet as part of the Go pass |
| golangci-lint | Use **existing** `.golangci.yml` (v2); includes staticcheck via that config—do not author a new config |
| Docker fallback | If CGO/`go-face` prevents host analysis, run the Go analyzers inside the project’s Docker/Compose app image |

**Ingest path:** Accept a pasted or file report instead of re-running when the user supplies findings.

### 4. Fix and suppression policy

- Prefer real fixes that preserve observable behavior.
- No blanket `//nolint` or `markdownlint-disable` (file- or project-wide).
- Narrow, line-scoped suppressions only when a real fix is impossible; treat as deferrals that MUST record finding id + rationale in the OpenSpec change (e.g. tasks or design).
- Do not change recognition pipeline code, models, or Docker native dependency installs to silence findings—defer those instead.

### 5. Round budget

**Choice:** At most two (fix → re-run) rounds. After the second re-run, stop, summarize remaining findings and recorded deferrals, and do not start a third fix pass unless the user explicitly asks.

**Rationale:** Prevents infinite lint thrash on style noise (e.g. MD013 on long teaching docs).

### 6. Change naming when the skill runs

**Choice:** Skill derives a kebab-case OpenSpec change name from the triage (e.g. `fix-static-analysis-md-go-<date>` or theme-based `fix-markdownlint-readme`). One change per `/opsx-static-fix` invocation unless the user names one.

### 7. Skill structure

**Choice:** Single `SKILL.md` with YAML frontmatter (`name: static-analysis-fix`, clear description). Optional short `reference.md` only if command tables overflow—prefer keeping analyzers and policies in `SKILL.md` for one-shot loads. Command mirrors other opsx commands: frontmatter name/id/category/description + steps that say “read and follow the skill”.

## Risks / Trade-offs

- [markdownlint MD013 floods README/DEVCONTAINER] → Group by rule; defer with id+rationale after two rounds rather than rewriting pedagogy docs wholesale.
- [CGO blocks host golangci] → Document Docker fallback using existing compose/service.
- [Agent invents new golangci config] → Skill MUST say use existing `.golangci.yml` only; tasks for this meta-change verify no new config file.
- [Skill scope creep into recognition “cleanup”] → Hard exclusion list in skill + specs.
- [OpenSpec overhead for one-line fixes] → Acceptable for teaching consistency; small finding sets still get a short change.

## Migration Plan

- Add skill and command only; no runtime migration.
- Rollback: delete `.cursor/skills/static-analysis-fix/` and `.cursor/commands/opsx-static-fix.md`.
- After merge, instructors run `/opsx-static-fix` (optionally with `01-TERMINAL_INTERACTION.md` as an ingest sample).

## Open Questions

- None blocking. Optional later: commit deferral log format under `openspec/changes/<name>/` vs. a standing `DEFERRED_LINT.md`—default to recording inside the active change artifacts.
