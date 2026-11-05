---
name: white-box-test-generation
description: >-
  Improve Go statement coverage by measuring with go test -coverprofile,
  generating coverage.html, inspecting uncovered code, and adding tests until a
  threshold is met (default 100% project-wide). Use when the user asks for
  white-box tests, coverage improvement, or closing gaps after black-box EP/BVA.
---

# White-Box Test Generation (statement coverage)

Raise **statement coverage** with Go’s built-in cover tools. Measure → write
`coverage.html` → inspect uncovered statements → add tests → re-measure until
the target is met (or blockers are reported).

This complements `.cursor/skills/black-box-test-generation/` (EP/BVA design).
Prefer black-box first for systematic design; use this skill to close remaining
coverage gaps. Do not strip existing EP/BVA comments.

## When to use

- User asks to improve statement coverage or generate white-box tests
- User wants `coverage.html` / coverprofile-driven gap closing
- User asks to finish coverage after a black-box EP/BVA pass
- User names a file, package, threshold, or “whole project / 100%”

## Hard exclusions (test-only)

While executing this skill, **MUST NOT** create or modify **production code**
(non-test `.go` under `cmd/` / `internal/` or other app sources)—including
recognition pipeline edits, injecting test seams, changing control flow, or
deleting hard-to-cover statements—nor model weights under `models/`, nor Docker
native dependency installation.

**Allowed levers only:**

- `*_test.go` (external `package x_test` or white-box `package x` in test files)
- Test fixtures under `testdata/` or other clearly test-only fixture paths
- Generated coverage artifacts (`coverage.out`, `coverage.html`) the skill
  already documents

If a gap cannot be closed with tests/fixtures alone, **stop-and-report**—do not
edit production under this skill.

## Stop-and-report (wait for the user)

When blocked, **STOP** the coverage loop, report, and **wait for the user**.
Do not silently introduce production seams or other production workarounds.

Each blocker report **MUST** include:

- **File** path
- **Line** (or function/region)
- **Reason** category: `unreachable` | `suspected bug` | `untestable design` |
  `missing fixtures` | `would require production change`

Examples aligned with the skill spec:

- **WHEN** covering a statement would require editing non-test source to inject
  a seam or alternate path **THEN** stop and report file/line/why
  (`would require production change`).
- **WHEN** uncovered statements cannot be executed via tests and fixtures alone
  **THEN** report them as blocked (`unreachable` / `untestable design`).
- **WHEN** closing a gap appears to require fixing or redesigning production
  code **THEN** stop and report (`suspected bug`) rather than patching
  production under this skill.

## Scope and threshold

Resolve before measuring:

| User input | Scope | Threshold |
|---|---|---|
| Unspecified | Entire module (`./...`) | **100%** statement coverage |
| Package path (e.g. `./internal/entity`) | That package | User value or **100%** |
| File path (e.g. `internal/entity/drawer.go`) | Its package; judge that **file’s** statements | User value or **100%** |
| Explicit % (e.g. 80%) | As above | That percentage |

## Measurement recipes

Default artifact names (unless the user overrides): `coverage.out`,
**`coverage.html`**. Treat them as generated; do not commit unless asked.

### Whole project (default scope)

```bash
go test ./... -coverprofile=coverage.out -coverpkg=./...
go tool cover -html=coverage.out -o coverage.html
go tool cover -func=coverage.out
```

`-coverpkg=./...` is required for a meaningful module-wide profile when using
`go test ./... -coverprofile=…`.

### Single package

```bash
go test ./path/to/package -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
go tool cover -func=coverage.out
```

### Single file

Run the package recipe for that file’s package, then use `coverage.html` and/or
filter `go tool cover -func` output to functions in that file. Success = those
statements meet the threshold (aim for no uncovered lines in that file in the
HTML report).

## Iteration loop

Cap at **10** fix/measure rounds. After the 10th re-measure, stop and report
leftovers even if below threshold.

1. **Resolve** scope + threshold (defaults above).
2. **Measure** with the matching recipe; write `coverage.html`.
3. **Inspect** uncovered statements:
   - Open/read `coverage.html` (red/uncovered regions), and/or
   - Use `go tool cover -func=coverage.out` for functions below 100%.
   - Read the corresponding source.
4. **Add tests** — smallest cases that execute uncovered statements **using
   only test/fixture edits**:
   - Prefer same-package tests (external `package x_test` or white-box `package x`
     when unexported code must be reached)—still only in `*_test.go`.
   - Assert observable results/errors/side effects; avoid assertion-free “touch
     only” tests when behavior can be checked.
   - Annotate new cases, e.g. `// Covers: uncovered path/file.go:L12-L40`.
   - Preserve existing EP/BVA catalog comments.
   - If the next step would be a production edit (seam, control-flow change,
     delete-to-cover), do **not** do it—treat as blocked (see Stop-and-report).
5. **Re-measure** (same recipe).
6. **Stop** when:
   - coverage for the scope ≥ threshold → report success, final %, artifact paths; or
   - blocked gaps remain → list file:line (or function) + reason category, then
     **wait for the user**; or
   - round cap hit → list remaining uncovered areas + last %.

## Relationship to black-box skill

- Optional order: **black-box-test-generation** (EP/BVA) → **this skill** (coverage gaps).
- White-box may invent cases EP would not; keep intent comments (`uncovered …`).
- Mutation testing (`gomutants`) is out of scope here; mention only as a possible
  follow-on after coverage is adequate.

## Output expectations

Report:

- Scope and threshold used
- Rounds executed
- Final `go tool cover -func` summary (or overall %)
- Paths to `coverage.out` / `coverage.html`
- Tests/files added or updated (`*_test.go` / fixtures only)
- Any blocked uncovered statements (file, line/region, reason category) and that
  you are waiting on the user before any production change
