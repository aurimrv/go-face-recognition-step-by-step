## Context

The class already teaches `go test -coverprofile` and `go tool cover -html` (see terminal notes). `black-box-test-generation` designs suites via EP/BVA; this skill closes remaining gaps by chasing **statement coverage** with Go’s built-in profiler. Default ambition is aggressive: **100% project-wide** unless the user scopes a file/package or lowers the threshold.

## Goals / Non-Goals

**Goals:**

- Ship `.cursor/skills/white-box-test-generation/SKILL.md` with a measure → HTML → inspect → add tests → re-measure loop.
- Default threshold **100%** statement coverage for `./...` when unspecified.
- Allow scope: one file, one package, or whole module; allow explicit threshold (e.g. 80%).
- Produce `coverage.html` (and a coverprofile) for human/agent inspection of red/uncovered regions.

**Non-Goals:**

- Replacing EP/BVA design (black-box skill remains the systematic design path).
- Mutation testing (`gomutants`)—orthogonal; may be mentioned as a follow-on only.
- Changing production recognition logic, models, or Docker native deps to make coverage easier.
- Guaranteeing 100% on CGO/`main` without tests that need fixtures/devcontainer—skill must still try, and report blockers honestly.

## Decisions

### 1. Tooling (locked)

| Step | Command |
|---|---|
| Measure | `go test -coverprofile=coverage.out <packages>` (and `-covermode=set` default is fine) |
| HTML | `go tool cover -html=coverage.out -o coverage.html` |
| Text aid | `go tool cover -func=coverage.out` for per-function % |

Profile and HTML default names: `coverage.out` and **`coverage.html`** at repo root unless the user names otherwise. Prefer adding both to `.gitignore` only if the repo already ignores them—do not expand gitignore scope in the skill unless already project policy; skill may note they are generated artifacts.

### 2. Scope and threshold resolution

1. If user names a **file** → cover its package; judge success by statements in that file (via HTML/`go tool cover -func` filtered to the file) meeting threshold.
2. If user names a **package** → `go test` that package; package statement % ≥ threshold.
3. If unspecified → `go test ./...` with a merged profile (or per-package profiles merged carefully); **entire project** statement coverage ≥ threshold.
4. If threshold unspecified → **100%**.

Implementation note for agents: `go test ./... -coverprofile=coverage.out` writes one profile for the last package in older Go; for whole-module profiles use:

```bash
go test ./... -coverprofile=coverage.out -coverpkg=./...
```

or package-by-package merge. The skill **MUST** document a reliable whole-project measurement recipe for this module (recommend `-coverpkg=./...` with tests under packages that exercise code, or iterative per-package until each hits the bar when project-wide merge is awkward).

**Choice:** Skill documents:

```bash
go test ./... -coverprofile=coverage.out -coverpkg=./...
go tool cover -html=coverage.out -o coverage.html
go tool cover -func=coverage.out
```

When scoped to one package: omit `-coverpkg` or set `-coverpkg` to that package path.

### 3. Iteration loop

1. Resolve scope + threshold.
2. Run coverprofile + write `coverage.html`.
3. Identify uncovered statements (HTML and/or `go tool cover -func` showing <100% functions; read source around uncovered blocks).
4. Add the **smallest** tests that execute those statements (table cases, error paths, helpers). Prefer tests in the same package (`_test` or white-box same-package tests for unexported if needed).
5. Re-measure; stop when threshold met **or** remaining gaps are blocked (document why: e.g. requires live models and user forbade fixture work—but skill may create test fixtures similar to black-box for recognition paths).
6. Cap runaway loops: suggest a practical max rounds (e.g. 10) then report leftover uncovered statements—mirrors static-analysis-fix stop discipline without being the same policy.

### 4. Relationship to black-box skill

- White-box may add cases that black-box would not have invented from EP alone.
- Do not strip EP/BVA comments from existing suites; new white-box cases should note `// Covers: uncovered <file:line>…` or similar so intent stays auditable.
- If both skills apply, black-box first for design, white-box to close gaps—skill text should say so briefly.

### 5. Hard exclusions

Same family as other project skills: do not rewrite recognition pipeline, models, or Docker native deps merely to raise coverage. Tests and fixtures are the lever.

## Risks / Trade-offs

- [100% on `cmd/main` / CGO is expensive] → Allow scoped runs; report honest blockers; use testdata/fixtures in devcontainer.
- [Merged `./...` coverprofile pitfalls] → Skill mandates `-coverpkg=./...` (or equivalent) for project-wide runs.
- [coverage.html noise / git dirty tree] → Treat as generated; optional restore; don’t commit unless user asks.
- [Tests that only “touch” lines without asserting] → Skill requires meaningful assertions where behavior is observable; void smoke-only when no return value exists (e.g. some draw paths)—prefer asserting side effects.

## Migration Plan

- Add skill only. Rollback: delete `.cursor/skills/white-box-test-generation/`.

## Open Questions

- None blocking. Optional later: Cursor command `/opsx-white-box-tests`.
