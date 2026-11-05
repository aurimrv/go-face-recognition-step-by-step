## Context

`cmd/main.go` hardcodes `images/unknown.jpg` → `images/result.jpg` after loading `persons/` and initializing one `face.Recognizer`. Classify and draw already work per image; only input discovery and output naming are missing. The teaching-fork layout already nests many JPEG test images under `images/<Person>/`, so a recursive walk is immediately useful. `Drawer.loadImage` only decodes JPEG and currently ignores errors, which will break as soon as a walk hits `.avif` or a corrupt file.

## Goals / Non-Goals

**Goals:**

- Support the three invocation forms from the exploration: no args; one or more files; one or more directories (recursive).
- Keep recognizer init and known-face samples outside the per-image loop.
- Isolate path rules in a pure resolver that can be unit-tested without dlib.
- Fail cleanly on non-decodable inputs; continue a multi-job run after logging per-image failures.

**Non-Goals:**

- Parallel recognition or a second recognizer instance.
- Non-JPEG formats (PNG, AVIF, etc.) or optional output directories.
- Changing person loading, thresholds, models, or drawing style.
- Introducing a CLI framework (cobra/urfave) for three argument shapes.

## Decisions

### 1. Pure resolver, then loop in `main` (not a mega use case)

**Choice:** `Resolve(args []string) ([]Job, error)` where `Job` has `Src` and `Dst`. `main` loads persons, builds the recognizer once, then iterates jobs.

**Alternatives considered:**
- Inline path logic only in `main` — smallest diff, but hard to test and easy to get wrong (result skipping, stem naming).
- One “analyze paths” use case owning recognizer + walk + draw — over-couples discovery to dlib and is harder to unit-test.

**Rationale:** Matches the script-shaped entrypoint while making the risky rules (default vs explicit naming, walk, skip results) testable in isolation.

### 2. Output naming

| Invocation | Destination |
|---|---|
| No args | `images/result.jpg` (unchanged) |
| Explicit file `…/name.ext` | `…/result_<stem>.jpg` beside source |
| File found under a walked directory | same as explicit file |

**Stem** = filename without extension. Extension check for JPEG is case-insensitive (`.jpg` / `.jpeg` / `.JPG` / …). Explicit `images/unknown.jpg` therefore yields `images/result_unknown.jpg`, not `result.jpg`.

**Rationale:** Beside-source naming avoids collisions when the same stem appears in different subfolders (e.g. multiple `unknown-01.jpg` under `images/`).

### 3. Argument expansion

- Zero args → one job for `images/unknown.jpg` → `images/result.jpg`.
- Each arg that is a regular file → one job if JPEG; otherwise error or skip per batch policy (non-JPEG explicit file is a job failure).
- Each arg that is a directory → `filepath.WalkDir`, collect JPEG files, map each to a beside-source `result_<stem>.jpg`.
- Mixed file and directory args are allowed; jobs are the concatenation in argument order (walk order within each directory is filesystem walk order).

### 4. Skipping previous outputs

**Choice:** After collecting candidate sources, drop any path whose base name is `result.jpg` or matches `result_*.jpg` (case-sensitive on the `result` prefix as produced by this tool). Also exclude any path that equals a computed destination of another job in the same run if that matters for ordering.

**Alternatives considered:** Precise “subtract all destinations from inputs” only — safer for a coincidental photo named `result_party.jpg`, slightly more code. Prefix skip is simpler and matches the tool’s own naming.

**Rationale:** Prevents `result_result_<stem>.jpg` on a second album run. Document that user photos should not use the `result_` prefix.

### 5. Batch error policy

**Choice:** For two or more jobs (or any run with explicit args), log the failing path and continue; exit non-zero if any job failed. “No faces” and decode errors are per-job failures. Empty recognition above threshold still saves an unmarked copy (today’s behavior) and counts as success for that job.

**Rationale:** Folder runs should not stop on the first landscape or `.avif`. The no-arg single-job path can still `log.Fatalf` for familiar teaching demos.

### 6. Drawer load errors

**Choice:** Make image (and preferably font) load errors return from construction or an exported init path so `main` can treat them as job failures instead of nil-image panics.

**Rationale:** Required once inputs are user-supplied.

### 7. Package placement

**Choice:** Prefer a small helper under something like `internal/usecases` or `internal/cli` (e.g. `ResolveAnalysisJobs`) with unit tests beside it. Avoid new top-level packages unless needed for clarity.

## Risks / Trade-offs

- [Second walk re-analyzes results] → Skip `result.jpg` / `result_*.jpg` basenames.
- [Name collision if someone stores a real photo as `result_foo.jpg`] → Document prefix convention; optional later: precise destination-set exclusion only.
- [Large albums are slow (sequential dlib)] → Acceptable; parallelization is a non-goal.
- [`.avif` under `images/Tom-Cruise/` fails decode] → JPEG-only filter on discovery; explicit non-JPEG args fail that job without aborting the whole batch when other jobs remain.
- [Teaching fork expects `go run cmd/main.go` with no args] → Preserve default job and `result.jpg` name.

## Migration Plan

- No data migration. Existing `images/result.jpg` remains the no-arg output.
- Update README / Docker exec examples to show file and directory forms.
- Rollback: revert CLI changes; hardcoded paths still work as today.

## Open Questions

- None blocking implementation. Optional later: `--out-dir` with relative-path preservation if albums must stay pristine.
