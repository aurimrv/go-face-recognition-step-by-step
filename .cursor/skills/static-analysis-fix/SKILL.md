---
name: static-analysis-fix
description: >-
  Resolve Markdown and Go static analysis findings via a spec-driven OpenSpec
  workflow (ingest or run analyzers, triage, propose/apply/archive, re-run).
  Use when the user runs /opsx-static-fix, asks to fix markdownlint/go
  vet/staticcheck/golangci-lint findings, or supplies a lint report for this repo.
---

# Static Analysis Fix

Spec-driven loop for clearing Markdown and Go analyzer findings in this project.
OpenSpec is the fix vehicle; analyzers verify. Do not invent a parallel workflow.

## When to use

- User invokes `/opsx-static-fix`
- User pastes or points at a markdownlint / go vet / golangci-lint / staticcheck report
- User asks to fix static analysis findings in this repository

## Hard exclusions

While executing this skill, **MUST NOT** change any of the following to silence findings:

- Recognition pipeline (face load / recognize / classify use cases and related flow)
- Model files under `models/`
- Docker native dependency installation (Dockerfile `apt`/native packages for dlib/CGO)

If a finding appears to require those changes, **defer** it with finding id + rationale.

## Analyzer commands (locked)

### Markdown

Run **only**:

```bash
markdownlint-cli2 "**/*.md" "#openspec" "#.cursor" "#node_modules"
```

Do not lint `openspec/`, `.cursor/`, or `node_modules` as part of this workflow.

### Go

1. `go vet ./...` (or equivalent package scope for the repo modules under test)
2. `golangci-lint run` using the repository’s **existing** `.golangci.yml` (**v2**). That config already enables staticcheck (and errcheck, govet, ineffassign, unused).

**MUST NOT** create a new `.golangci.yml`, a parallel golangci config, or expand linter enablement beyond the existing file for this workflow.

### Docker fallback (CGO)

If host-side Go analysis fails because of CGO / go-face native libraries, re-run the Go analyzers inside the project Docker/Compose app environment (for example `docker compose` / `docker exec` on the `go-face-recognition` service), still using the same existing `.golangci.yml` in the repo.

### Ingest path

If the user supplies a report (paste or file), triage from that report first. A fresh analyzer run is not required before creating the OpenSpec change; always re-run analyzers after fix rounds to verify.

## Fix and deferral policy

- **Preserve behavior** — fixes must not change observable recognition or CLI behavior except where the finding itself is a true defect (e.g. ignored error). Prefer the smallest correct edit.
- **Prefer real fixes** over suppressions.
- **No blanket suppressions** — MUST NOT add file- or project-wide `//nolint`, `nolint:all`, or `markdownlint-disable` (including HTML comment disables that cover a whole file).
- Narrow, line-scoped suppressions only when a real fix is impossible; treat each as a **deferral**.
- **Every deferral MUST record** the finding id (rule + location as reported) and a written rationale in the OpenSpec change artifacts (tasks.md and/or design.md).

## Round limit

At most **two** (fix → re-run) rounds per `/opsx-static-fix` invocation.

After the **second** re-run, if findings remain: **stop**. Do not start a third fix round unless the user explicitly asks. Report what is left, including recorded deferrals.

## Workflow

Follow these steps in order.

### 1. Ingest or run

- If a report was provided, parse it into a finding list.
- Otherwise run Markdown and Go analyzers per **Analyzer commands** above (with Docker fallback if needed).

### 2. Triage and group

- Group by tool, rule id, and file/area (e.g. all MD013 in README; errcheck in one package).
- Mark candidates to fix vs defer (exclusions, impossible-without-behavior-change, or style noise after judging effort).
- Stable finding ids: prefer `TOOL:RULE:path:line` (or the tool’s native id if unique).

### 3. Create an OpenSpec change

- Derive a kebab-case name (e.g. `fix-markdownlint-readme`, `fix-static-analysis-md-go`) unless the user provided one.
- Create the change and produce **proposal**, **design**, **specs**, and **tasks** (same artifact completeness as `/opsx-propose`).
- Tasks must be small, checkbox-sized fix units. List deferred findings with id + rationale in design or tasks.
- Announce the change name.

Use existing project OpenSpec skills/commands (`openspec-propose` / `/opsx-propose` patterns) when helpful; the change must still be created and apply-ready before fixing.

### 4. Apply fixes task by task

- Implement only the fix tasks; mark each checkbox done when finished.
- Obey fix/deferral policy and hard exclusions.
- Do not “drive-by” refactor unrelated code.

### 5. Re-run analyzers

- Re-run Markdown and Go commands (Docker fallback if required).
- Compare to the triage list: fixed, remaining, new.

### 6. Second round or finish

- If clean (no remaining non-deferred findings): summarize, then **archive** the OpenSpec change (`/opsx-archive` / openspec-archive-change).
- If issues remain and fewer than two fix/re-run rounds used: update tasks if needed, apply another fix round, re-run (counts as round 2).
- If issues remain after the second re-run: **stop**, report leftovers and deferrals, leave the change open unless the user asks to archive with warnings.

## Output expectations

While running, briefly state: current step, change name, round number (1 or 2).

On completion or stop, report:

- What was fixed
- Deferred findings (id + rationale)
- Analyzer status after the last re-run
- Whether the change was archived or left active
