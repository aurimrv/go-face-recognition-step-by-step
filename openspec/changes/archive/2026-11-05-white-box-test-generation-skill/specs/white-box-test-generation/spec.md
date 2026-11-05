## ADDED Requirements

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
While executing this skill, the agent MUST NOT change the recognition pipeline implementation, model files under `models/`, or Docker native dependency installation solely to increase coverage; tests and fixtures are the allowed levers.

#### Scenario: Coverage via tests not production rewrites
- **WHEN** uncovered code sits in recognition or model-loading paths
- **THEN** the agent adds tests/fixtures (or reports a blocker) rather than deleting or bypassing that production code to inflate coverage
