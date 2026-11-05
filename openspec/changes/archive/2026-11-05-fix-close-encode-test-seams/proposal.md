## Why

White-box coverage stopped at **98.7%** and reported three production statements that cannot be hit with tests alone: `drawer.loadImage` close-after-decode failure, `drawer.SaveImage` JPEG encode failure, and `load_persons` directory close-failure logging. The test-generation skills correctly refused to inject seams; this change is the user-authorized production fix so those error paths become testable and whole-module statement coverage can reach 100%.

## What Changes

- Add minimal package-level I/O seams in `internal/entity` (file close + JPEG encode) and `internal/usecases` (persons-dir close), defaulting to real `Close` / `jpeg.Encode`
- Wire `loadImage`, `SaveImage`, and `LoadPersonsUseCase.Execute` through those seams (behavior unchanged in production)
- Restore/add `*_test.go` cases that force seam failures and assert the existing error/log outcomes
- Re-measure with `go test ./... -coverprofile=coverage.out -coverpkg=./...` and confirm **100%** statement coverage
- Does **not** weaken the white-box/black-box test-only guardrails; those skills still must not invent production edits without user/OpenSpec authorization

## Capabilities

### New Capabilities

- `testable-io-error-paths`: Production packages expose narrow, default-real I/O seams so Close/Encode failure branches in drawer and load-persons can be exercised from tests

### Modified Capabilities

- (none)

## Impact

- `internal/entity/drawer.go` (+ matching `*_test.go`)
- `internal/usecases/load_persons.go` (+ matching `*_test.go`)
- Coverage artifacts regenerated during verify; no models/Docker/recognition-algorithm changes
