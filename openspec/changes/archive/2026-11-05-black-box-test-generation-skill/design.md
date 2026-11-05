## Context

The repo is a teaching fork for GenAI-assisted testing/refactoring. Packages under test today include pure path logic (`ResolveAnalysisJobs`) and CGO-backed face recognition (`ClassifyPersons` with a float threshold, person loading, drawing). Only `resolve_analysis_jobs_test.go` exists so far. Students already use `go test -cover` and `gomutants`; they need a skill that produces *systematic black-box* suites (EP + BVA) with auditable partition comments—and realistic image fixtures when match strength matters.

## Goals / Non-Goals

**Goals:**

- Ship `.cursor/skills/black-box-test-generation/SKILL.md` that agents follow when asked to generate EP/BVA Go tests.
- Mandate partition catalog comments and per-case “covers …” annotations.
- Guide creation of face/unknown JPEG fixtures with controlled similarity for recognition/classification boundaries (e.g. around the 0.3 threshold used in `cmd`).

**Non-Goals:**

- Generating the full production test suite as part of *this* change (that is a later skill run).
- White-box / structural coverage design (branch tables, mutation operators)—EP/BVA only.
- Changing recognition algorithms, models, or Docker native deps.
- Requiring a Cursor command wrapper in this change (skill-only unless a thin command is added later).

## Decisions

### 1. Skill-only delivery (mirror static-analysis-fix pattern, no command required)

**Choice:** One project skill directory with `SKILL.md`. Optional command can wait.

**Rationale:** User asked for the skill path explicitly; keep scope tight.

### 2. Black-box stance

**Choice:** Derive partitions from observable inputs/outputs/errors of exported APIs (and CLI argument shapes where testing `main` helpers), not from internal branch structure. Reading code to *list* parameters is allowed; tests must not assert on unexported internals.

**Alternatives:** Pure white-box path covering—rejected for this skill’s charter.

### 3. Comment contract (locked)

Generated `_test.go` files SHALL include:

1. A file- or suite-level comment block listing **equivalence partitions** (id, domain, representative meaning).
2. A file- or suite-level note of **boundaries** (which edges of which partitions).
3. On **each** test function (or table-row via comment on the row / subtest name), which partition id and/or boundary id it covers.

Example shape (illustrative):

```go
// Equivalence partitions for ResolveAnalysisJobs(args):
//   EP-empty:   len(args)==0 → default unknown→result job
//   EP-files:   args are existing JPEG files → beside-source result_<stem>.jpg
//   EP-dirs:    args are directories → recursive JPEG discovery
// Boundaries:
//   B-ext:      .jpg vs .jpeg vs non-JPEG under dirs
//   B-result:   basename result.jpg / result_*.jpg skipped
func TestResolve_EmptyArgs_DefaultJob(t *testing.T) {
    // Covers: EP-empty
    ...
}
```

### 4. When image fixtures are required

**Choice:** If the unit under test involves face descriptors, classification threshold, or “recognized vs not,” the skill MUST plan fixtures:

| Fixture role | Intent |
|---|---|
| Known-person gallery JPEGs | Build samples (strong identity evidence) |
| Unknown strong match | Same person, clear face → expect classify hit |
| Unknown weak / near-threshold | Similar but harder match → boundary around `threshold` |
| Unknown non-match / no face | Different person or no detectable face → miss or error |

Prefer `testdata/` (or a dedicated test tree) over mutating class `persons/` / `images/` demos. Synthesize or copy/crop existing teaching images; document provenance in comments. Do not replace `models/`.

### 5. Workflow inside the skill

1. Identify unit under test and observable I/O.
2. Derive EP set (valid/invalid classes) and BVA edges.
3. Write partition/boundary catalog comments.
4. Emit table-driven Go tests where practical; one case per selected representative/boundary.
5. Create fixtures if needed; keep tests hermetic (temp dirs when filesystem is the API).
6. Run `go test` for the package; adjust only within EP/BVA scope.

### 6. Targets in this codebase (guidance, not exclusive)

- High value without CGO: `ResolveAnalysisJobs` (args shapes, JPEG filter, result skip).
- CGO / models: `ClassifyPersonsUseCase.Execute` threshold and image content partitions—needs recognizer + fixtures; slower; Docker/devcontainer assumed.

## Risks / Trade-offs

- [Synthetic faces may not hit intended match bands] → Prefer real photos from `persons/` / `images/` with documented similarity intent; treat threshold BVA as “near / clearly above / clearly below” with measured runs, not theoretical pixel math.
- [Agents invent white-box tests] → Skill insists on EP/BVA comments; reject “covers line 42” as the sole rationale.
- [Fixture bloat in git] → Prefer small JPEGs under `testdata/`; reuse existing assets via copy into temp dirs when possible.
- [CGO flakiness] → Skill should allow skipping or tagging recognition tests when models/CGO unavailable, while still generating EP/BVA structure for pure packages.

## Migration Plan

- Add skill only; rollback by deleting `.cursor/skills/black-box-test-generation/`.
- Later runs create tests/fixtures; those are separate student commits.

## Open Questions

- None blocking. Optional later: `/opsx-black-box-tests` command wrapper.
