# white-box-test-generation Specification

## Purpose
Mainlined from change white-box-test-generation-skill, with test-only guardrail from add-test-only-guardrail.
## Requirements
### Requirement: Project skill exists
The project SHALL provide a Cursor skill at `.cursor/skills/white-box-test-generation/SKILL.md` that instructs the agent to improve Go statement coverage through measurement, HTML report inspection, and iterative test generation.

#### Scenario: Skill discoverable
- **WHEN** a contributor inspects `.cursor/skills/white-box-test-generation/`
- **THEN** a `SKILL.md` with name `white-box-test-generation` is present

#### Scenario: Skill invoked for coverage improvement
- **WHEN** the user asks to raise statement coverage or run white-box test generation for this repo
- **THEN** the agent follows the `white-box-test-generation` skill

### Requirement: Measure with coverprofile and HTML report
The skill SHALL direct the agent to measure statement coverage with `go test -coverprofile` and to generate an HTML coverage report named `coverage.html` using `go tool cover -html` (unless the user specifies alternate output names).

#### Scenario: Default report artifacts
- **WHEN** the agent runs a coverage measurement pass under this skill without alternate names
- **THEN** it produces a cover profile and writes `coverage.html` via `go tool cover -html=… -o coverage.html`

### Requirement: Default whole-project 100% threshold
When the user does not specify a coverage threshold, the skill SHALL target 100% statement coverage. When the user does not specify a file or package scope, the skill SHALL target the entire project module.

#### Scenario: Unspecified goal
- **WHEN** the user invokes white-box test generation without a threshold or scope
- **THEN** the agent works toward 100% statement coverage across the project

#### Scenario: Explicit threshold or scope
- **WHEN** the user names a threshold and/or a specific file or package
- **THEN** the agent uses that threshold and scopes measurement and test generation accordingly

### Requirement: Inspect uncovered code and add tests iteratively
The skill SHALL direct the agent to inspect uncovered statements (via the HTML report and/or `go tool cover -func`), add tests that execute those statements, re-measure, and repeat until the threshold is met or remaining gaps are reported as blocked.

#### Scenario: Gap-closing loop
- **WHEN** measured coverage is below the target for the chosen scope
- **THEN** the agent identifies uncovered statements, adds or extends tests to cover them, and re-runs the coverage measurement

#### Scenario: Stop when target met
- **WHEN** measured statement coverage for the chosen scope meets or exceeds the threshold
- **THEN** the agent stops the loop and reports the final coverage and artifacts

### Requirement: Hard exclusions while raising coverage
While executing this skill, the agent MUST NOT change production implementation (including the recognition pipeline), model files under `models/`, or Docker native dependency installation solely to increase coverage. Tests and fixtures are the only allowed levers; if a gap cannot be closed under Requirement: Test-only edits and stop-and-report, the agent MUST report it as blocked.

#### Scenario: Coverage via tests not production rewrites
- **WHEN** uncovered code sits in recognition or model-loading paths
- **THEN** the agent adds tests/fixtures (or reports a blocker) rather than deleting, bypassing, or seam-injecting that production code to inflate coverage

#### Scenario: Stop instead of production seam
- **WHEN** covering a statement would require editing non-test source to introduce a test seam or alternate control path
- **THEN** the agent stops and reports the blocker (file, line/region, reason) instead of modifying production code under this skill

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
