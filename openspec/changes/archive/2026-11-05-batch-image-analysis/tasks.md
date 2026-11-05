## 1. Path resolution

- [x] 1.1 Add a small `Job` type (`Src`, `Dst`) and `ResolveAnalysisJobs(args []string) ([]Job, error)` helper (e.g. under `internal/usecases` or `internal/cli`)
- [x] 1.2 Implement no-arg default: `images/unknown.jpg` → `images/result.jpg`
- [x] 1.3 Implement explicit file args: JPEG only → beside-source `result_<stem>.jpg`
- [x] 1.4 Implement directory args: recursive walk, case-insensitive `.jpg`/`.jpeg`, same beside-source naming
- [x] 1.5 Skip inputs whose base name is `result.jpg` or matches `result_*.jpg`
- [x] 1.6 Add unit tests for default, files, nested directories, non-JPEG skip, result-file skip, and mixed args

## 2. Drawer robustness

- [x] 2.1 Surface image (and font) load errors from drawer construction so callers can handle failures
- [x] 2.2 Ensure a failed load does not panic on nil image when saving or drawing

## 3. CLI batch loop

- [x] 3.1 Wire `os.Args` through the resolver in `cmd/main.go`
- [x] 3.2 Keep persons load + recognizer init + `SetSamples` once; loop classify → draw → save per job
- [x] 3.3 For explicit-argument runs, continue after per-job failures, report each failure, exit non-zero if any failed
- [x] 3.4 Preserve fail-fast (`log.Fatalf`-style) behavior for the no-arg single default job

## 4. Docs and manual check

- [x] 4.1 Update README usage for: `go run ./cmd`, `go run ./cmd images/a.jpg images/b.jpg`, `go run ./cmd path/to/album/`
- [x] 4.2 Manually verify no-arg still writes `images/result.jpg`, and a small folder run writes beside-source `result_<stem>.jpg` without re-analyzing prior results
