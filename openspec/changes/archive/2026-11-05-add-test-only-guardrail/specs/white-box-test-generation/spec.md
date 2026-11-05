## ADDED Requirements

### Requirement: Test-only edits and stop-and-report
While executing the white-box-test-generation skill, the agent MUST NOT create or modify production code. The agent MAY only create or modify test code and test fixtures (`*_test.go`, files under `testdata/`, and other clearly test-only fixtures). If a problem prevents reaching the coverage objective without changing production code (including unreachable code, a suspected bug, an untestable design, or missing fixtures the skill cannot supply under this rule), the agent MUST stop, report the blocker with file path, line (or function/region), and reason, and MUST wait for the user to decide—without applying production edits under this skill.

#### Scenario: Would need production change
- **WHEN** the skill would need to change production code (for example injecting seams, changing control flow, or deleting hard-to-cover statements) to raise coverage
- **THEN** it stops and reports the file, line (or region), and why instead of editing production code

#### Scenario: Unreachable or untestable statements
- **WHEN** uncovered statements cannot be executed via tests and fixtures alone
- **THEN** the agent stops and reports those statements as blocked (file, line/region, reason) for the user to decide

#### Scenario: Suspected bug blocking coverage
- **WHEN** closing a coverage gap appears to require fixing or redesigning production code
- **THEN** the agent stops and reports the suspected issue rather than patching production under this skill

#### Scenario: Allowed test artifacts only
- **WHEN** the skill adds or updates tests to execute uncovered statements
- **THEN** changes are limited to `*_test.go` and test fixtures (e.g. under `testdata/`)

## MODIFIED Requirements

### Requirement: Hard exclusions while raising coverage
While executing this skill, the agent MUST NOT change production implementation (including the recognition pipeline), model files under `models/`, or Docker native dependency installation solely to increase coverage. Tests and fixtures are the only allowed levers; if a gap cannot be closed under Requirement: Test-only edits and stop-and-report, the agent MUST report it as blocked.

#### Scenario: Coverage via tests not production rewrites
- **WHEN** uncovered code sits in recognition or model-loading paths
- **THEN** the agent adds tests/fixtures (or reports a blocker) rather than deleting, bypassing, or seam-injecting that production code to inflate coverage

#### Scenario: Stop instead of production seam
- **WHEN** covering a statement would require editing non-test source to introduce a test seam or alternate control path
- **THEN** the agent stops and reports the blocker (file, line/region, reason) instead of modifying production code under this skill
