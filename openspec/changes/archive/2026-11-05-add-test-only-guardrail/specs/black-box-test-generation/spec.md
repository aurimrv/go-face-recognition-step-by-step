## ADDED Requirements

### Requirement: Test-only edits and stop-and-report
While executing the black-box-test-generation skill, the agent MUST NOT create or modify production code. The agent MAY only create or modify test code and test fixtures (`*_test.go`, files under `testdata/`, and other clearly test-only fixtures). If a problem prevents meeting the test-generation objective without changing production code (including unreachable code, a suspected bug, an untestable design, or missing fixtures the skill cannot supply under this rule), the agent MUST stop, report the blocker with file path, line (or function/region), and reason, and MUST wait for the user to decide—without applying production edits under this skill.

#### Scenario: Would need production change
- **WHEN** the skill would need to change production code to make a black-box EP/BVA objective achievable
- **THEN** it stops and reports the file, line (or region), and why instead of editing production code

#### Scenario: Suspected bug or untestable design
- **WHEN** observable behavior looks incorrect or the unit cannot be exercised meaningfully through its public API without production changes
- **THEN** the agent stops and reports the blocker (file, line/region, reason) for the user to decide

#### Scenario: Missing fixtures under test-only scope
- **WHEN** required fixtures cannot be created or reused within test-only paths
- **THEN** the agent stops and reports what is missing rather than altering production assets or models

#### Scenario: Allowed test artifacts only
- **WHEN** the skill adds or updates tests to cover chosen partitions and boundaries
- **THEN** changes are limited to `*_test.go` and test fixtures (e.g. under `testdata/`)

## MODIFIED Requirements

### Requirement: Black-box scope and hard exclusions
The skill SHALL target observable inputs and outputs of the unit under test. While executing the skill, the agent MUST NOT change production implementation (including recognition pipeline code), model weights, or Docker native dependency installation in order to make tests pass. Allowed levers are tests and fixtures only, subject to Requirement: Test-only edits and stop-and-report.

#### Scenario: Tests adapt to behavior, code under test stays
- **WHEN** a generated test fails because production behavior differs from an incorrect expectation
- **THEN** the agent corrects the test or fixtures under EP/BVA rationale and does not modify production code, models, or Docker native deps to force a pass

#### Scenario: Stop instead of production workaround
- **WHEN** correcting expectations or fixtures is insufficient and production code would need to change for the suite to meet its objective
- **THEN** the agent stops and reports the blocker rather than editing production code under this skill
