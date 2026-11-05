## Why

Statement coverage in this teaching fork is still low outside a few packages. Agents need a repeatable white-box skill that measures coverage with Go’s built-in tools, inspects uncovered statements via an HTML report, and adds tests until a stated threshold is met (default: 100% for the whole project).

## What Changes

- Add project skill `.cursor/skills/white-box-test-generation/SKILL.md` to improve Go statement coverage iteratively.
- Require measuring with `go test -coverprofile=…`, generating `coverage.html` via `go tool cover -html=… -o coverage.html`, inspecting uncovered code, and adding tests until the target is reached.
- Support a user-specified threshold and scope (specific file/package or whole project); if unspecified, target **100% statement coverage for the entire project**.
- Complement (not replace) the black-box EP/BVA skill: white-box may add cases driven by uncovered statements; it must not weaken hard exclusions around models/Docker native deps/recognition redesign solely to inflate coverage.

## Capabilities

### New Capabilities

- `white-box-test-generation`: Agent workflow for coverage-driven Go test generation using `go test -coverprofile`, HTML coverage reports, uncovered-statement analysis, and iterative test addition to a coverage threshold.

### Modified Capabilities

- (none)

## Impact

- New files under `.cursor/skills/white-box-test-generation/`.
- Later skill runs may add/update `*_test.go` and write `coverage.out` / `coverage.html` (and similar artifacts) in the workspace.
- No change to recognition algorithms, models, or Docker native dependencies as part of delivering the skill itself.
