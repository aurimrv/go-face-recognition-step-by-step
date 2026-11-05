---
name: black-box-test-generation
description: >-
  Generate Go tests for this repo using equivalence partitioning (EP) and
  boundary value analysis (BVA), with partition catalogs and per-case coverage
  comments. Use when the user asks for black-box, EP, or BVA tests, or for
  systematic tests of package APIs (including face-match fixtures when needed).
---

# Black-Box Test Generation (EP + BVA)

Produce **black-box** Go tests from **equivalence partitioning** and **boundary
value analysis**. Every suite must document partitions and label what each case
covers. Prefer observable inputs → outputs/errors; do not treat line coverage
or mutation operators as the design method here.

## When to use

- User asks for black-box, EP, and/or BVA tests for a Go package, function, or API
- User wants systematic tests of CLI/path helpers (e.g. `ResolveAnalysisJobs`) or
  classification/recognition behavior with match-strength fixtures
- User asks to improve the suite design before measuring `go test -cover` or `gomutants`

## Hard exclusions (test-only)

While executing this skill, **MUST NOT** create or modify **production code**
(non-test `.go` under `cmd/` / `internal/` or other app sources), including
recognition pipeline changes, test seams in production, model weights under
`models/`, or Docker native dependency installation.

**Allowed edits only:**

- `*_test.go` (and test helpers that live in test files)
- Test fixtures under `testdata/` or other clearly test-only fixture paths
- Temp dirs created at test time

If expectations were wrong, fix the **test** or **fixtures** under EP/BVA
rationale—never production code under this skill.

## Stop-and-report (wait for the user)

If a problem prevents meeting the EP/BVA objective **without** changing
production code, **STOP**. Do not silently edit production. Report and **wait
for the user to decide**.

Each blocker report **MUST** include:

- **File** path
- **Line** (or function/region)
- **Reason** category: `unreachable` | `suspected bug` | `untestable design` |
  `missing fixtures` | `would require production change`

Examples aligned with the skill spec:

- **WHEN** the skill would need to change production code to make an EP/BVA
  objective achievable **THEN** stop and report file/line/why instead of editing
  production code.
- **WHEN** observable behavior looks incorrect or the unit cannot be exercised
  meaningfully through its public API without production changes **THEN** stop
  and report (`suspected bug` / `untestable design`).
- **WHEN** required fixtures cannot be created or reused within test-only paths
  **THEN** stop and report what is missing (`missing fixtures`) rather than
  altering production assets or models.
- **WHEN** correcting expectations or fixtures is insufficient and production
  code would need to change **THEN** stop and report the blocker rather than
  applying a production workaround under this skill.

## Black-box stance

- Derive partitions from **observable** parameters, return values, and errors of
  exported APIs (and documented CLI argument shapes when testing those helpers).
- Reading source to list inputs is fine; tests **MUST NOT** assert on unexported
  internals or “covers line N” as the sole rationale.
- Table-driven tests and subtests are preferred when they keep one row ↔ one
  partition/boundary clear.

## Workflow

Follow these steps in order.

### 1. Identify the unit and observable I/O

- Name the package/function under test.
- List inputs (args, paths, thresholds, image content roles) and observable
  outcomes (return values, errors, written files, recognized IDs).

### 2. Derive equivalence partitions (EP)

- Split each interesting input domain into classes that should behave the same.
- Assign stable ids (`EP-empty`, `EP-files`, `EP-strong-match`, …).
- Include invalid/error classes when the API can fail observably.

### 3. Derive boundaries (BVA)

- Mark edges of partitions and notable edges (empty vs one element, `.jpg` vs
  `.jpeg` vs non-JPEG, threshold clearly below / near / clearly above, result
  basename skip patterns, etc.).
- Assign stable boundary ids (`B-ext`, `B-threshold-low`, …).

### 4. Write the catalog comments

In the generated `*_test.go` (file- or suite-level comment), **MUST** include:

1. Equivalence partitions: id + brief meaning
2. Boundaries: id + which edge

### 5. Emit Go tests

- At least one representative per selected partition.
- At least one annotated case per selected boundary you choose to cover (not
  only interior points).
- **Each** test function, subtest, or table row **MUST** cite the partition
  and/or boundary it covers (comment and/or subtest name), e.g. `// Covers: EP-empty`
  or `t.Run("EP-files/B-ext-jpeg", ...)`.

### 6. Fixtures (only if match degree matters)

See **Face and unknown image fixtures** below.

### 7. Run and adjust

```bash
go test ./path/to/package -count=1
```

Adjust tests/fixtures within EP/BVA scope only. If you would need production
edits, follow **Stop-and-report**. Optionally note coverage with `go test
-cover` afterward; do not redesign the suite as white-box to chase %.

## Comment contract (required)

### Suite-level catalog

```go
// Equivalence partitions for ResolveAnalysisJobs(args):
//   EP-empty:  len(args)==0 → default images/unknown.jpg → images/result.jpg
//   EP-files:  args are existing files → beside-source result_<stem>.jpg
//   EP-dirs:   args are directories → recursive JPEG discovery
// Boundaries:
//   B-ext:     under dirs, .jpg/.jpeg selected; non-JPEG omitted
//   B-result:  basename result.jpg / result_*.jpg skipped as inputs
```

### Per-case annotation

```go
func TestResolveAnalysisJobs_EmptyArgs(t *testing.T) {
	// Covers: EP-empty
	jobs, err := ResolveAnalysisJobs(nil)
	// ...
}

func TestResolveAnalysisJobs_Table(t *testing.T) {
	tests := []struct {
		name    string
		covers  string
		args    []string
		wantErr bool
	}{
		// Covers: EP-files, B-ext (explicit .JPG)
		{name: "explicit-jpeg", covers: "EP-files/B-ext", args: []string{"…"}},
	}
	for _, tt := range tests {
		t.Run(tt.name+"/"+tt.covers, func(t *testing.T) {
			// Covers: tt.covers
			// ...
		})
	}
}
```

## Face and unknown image fixtures

Use when the unit involves face descriptors, classification threshold, or
“recognized vs not” (e.g. `ClassifyPersonsUseCase.Execute`). Pure path logic
such as `ResolveAnalysisJobs` needs **no** face fixtures.

| Fixture role | Intent |
|---|---|
| Known-person gallery JPEGs | Samples for `SetSamples` / known identities |
| Unknown strong match | Same person, clear face → expect a classify hit |
| Unknown weak / near-threshold | Similar face → boundary around `threshold` (e.g. 0.3 in `cmd`) |
| Unknown non-match / no face | Different person or no detectable face → miss or error |

Rules:

- Prefer creating or copying fixtures under **`testdata/`**, or copy into
  **temp dirs** at test time. Prefer reusing/cropping teaching images from
  `persons/` / `images/` over inventing random pixels.
- Document fixture **intent** and provenance in comments next to paths.
- **MUST NOT** replace `models/`, alter Docker native deps, or change any
  production code to force a pass; if fixtures are insufficient, stop-and-report.
- Threshold BVA: aim for clearly above / near / clearly below with measured
  runs in the devcontainer—not theoretical pixel math alone.
- If CGO/`go-face`/models are unavailable, generate the EP/BVA structure and
  skip or tag recognition tests with a clear reason; still fully implement
  suites for pure packages.

## Project targets (guidance)

- **No CGO:** `internal/usecases.ResolveAnalysisJobs` — args, JPEG filter, result skip.
- **CGO + models (devcontainer):** classification/recognition with threshold and
  image-content partitions; slower; use fixtures above.

## Output expectations

When finishing a generation pass, briefly report:

- Unit under test
- EP and boundary ids chosen
- Files added/updated (`*_test.go`, `testdata/…` only)
- `go test` result for the package
- Any blockers (file, line/region, reason) — and that you are waiting on the user
