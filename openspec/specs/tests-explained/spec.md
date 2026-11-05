# tests-explained Specification

## Purpose
Mainlined from change tests-explained-traceability: human-readable `TESTS_EXPLAINED.md` generated from `tests-explained.yaml`, with skill hooks so black/white/mutation test generation keeps the catalog current.
## Requirements
### Requirement: Human catalog TESTS_EXPLAINED.md
The project SHALL maintain a repository-root file `TESTS_EXPLAINED.md` that explains, for each documented test case, the test objective, the skill that generated or last updated it, and the part of the system it is responsible for testing, organized for human reading by package and test file.

#### Scenario: Document discoverable
- **WHEN** a contributor opens the repository root
- **THEN** `TESTS_EXPLAINED.md` is present and states that it is generated from `tests-explained.yaml`

#### Scenario: Organization by package and file
- **WHEN** a reader browses `TESTS_EXPLAINED.md`
- **THEN** entries are grouped under package headings and per-`*_test.go` subheadings with one subsection per test function

### Requirement: YAML registry as source of truth
The project SHALL maintain a repository-root file `tests-explained.yaml` as the machine-editable source of truth for test explanations. `TESTS_EXPLAINED.md` SHALL be regenerated from that registry so the Markdown stays synchronized with the registry.

#### Scenario: Registry fields
- **WHEN** a test case is recorded in `tests-explained.yaml`
- **THEN** the entry includes at least: stable `id`, `test` function name, `test_file`, `skill`, `objective`, `system_under_test`, and `production_refs`

#### Scenario: Markdown projection
- **WHEN** the registry is updated by a test-generation skill
- **THEN** `TESTS_EXPLAINED.md` is regenerated to reflect the current registry contents

### Requirement: Traceability to current suite
Documented cases SHALL remain consistent with the current test suite: each registry `id` SHALL refer to an existing test function in the listed `test_file`, and skills SHALL update or remove entries when they rename or delete tests they own in a pass.

#### Scenario: Remove stale entry on delete
- **WHEN** a test-generation skill deletes a test function it previously registered
- **THEN** it removes or updates the corresponding registry entry and regenerates `TESTS_EXPLAINED.md`

#### Scenario: Refresh production refs when touching a case
- **WHEN** a skill updates a test that exercises production code
- **THEN** it refreshes `system_under_test` and `production_refs` to match the current implementation paths or symbols under test
