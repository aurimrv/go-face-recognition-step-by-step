## ADDED Requirements

### Requirement: Project skill exists
The project SHALL provide a Cursor skill at `.cursor/skills/mutation-test-generation/SKILL.md` that instructs the agent to generate or improve Go tests guided by live mutants using `gomutants`.

#### Scenario: Skill discoverable
- **WHEN** a contributor inspects `.cursor/skills/mutation-test-generation/`
- **THEN** a `SKILL.md` with name `mutation-test-generation` is present

#### Scenario: Skill invoked for mutation-guided tests
- **WHEN** the user asks for mutation testing, killing live mutants, or gomutants-guided test generation in this repo
- **THEN** the agent follows the `mutation-test-generation` skill

### Requirement: Use gomutants for mutant generation and execution
While executing the skill, the agent MUST use `gomutants` as the mutant generation and execution tool (not alternate mutators such as gremlins).

#### Scenario: Default tool
- **WHEN** the agent measures or re-measures mutants under this skill
- **THEN** it invokes `gomutants` on the resolved package/path scope

### Requirement: Tests target live mutants
The skill SHALL direct the agent to inspect live (surviving) mutants and to create or modify tests that attempt to kill all live mutants in scope, or a user-specified mutant when one is named.

#### Scenario: Kill all live mutants in scope
- **WHEN** the user does not name a specific mutant
- **THEN** the agent works through live mutants from the gomutants report and adds or updates tests aimed at killing them

#### Scenario: Kill a specific mutant
- **WHEN** the user identifies a specific mutant (e.g. by file, line, and/or mutant id from the report)
- **THEN** the agent prioritizes tests that attempt to kill that mutant and re-runs `gomutants` to verify the outcome

### Requirement: Final JSON and HTML mutation reports
At the end of a skill run, the agent MUST produce the gomutants report artifacts, including an HTML mutation report named `mutation-report.html` written with `gomutants` `-html-output` (unless the user overrides the HTML path).

#### Scenario: HTML report written
- **WHEN** the agent finishes a mutation-guided pass without an alternate HTML path
- **THEN** it writes `mutation-report.html` via `gomutants … -html-output mutation-report.html` (alongside the JSON report)

#### Scenario: Report paths in summary
- **WHEN** the agent reports completion or stop-and-report
- **THEN** the summary includes paths to the JSON and HTML mutation reports used

### Requirement: Justify possible equivalent mutants
When mutants remain live after the kill loop (or cannot be killed under test-only constraints), the agent MUST justify possible equivalent mutants that it is not possible to kill, citing file, line/region, and reasoning (e.g. no observable behavioral change / oracle limitation). Optional use of `gomutants -detect-equivalent` MAY inform but MUST NOT replace written justification.

#### Scenario: Equivalent survivors documented
- **WHEN** one or more mutants still live at the end of the run
- **THEN** the agent lists each remaining live mutant with file, line/region, and a brief justification when it appears equivalent or otherwise unkillable with tests alone

### Requirement: Test-only edits and stop-and-report
While executing the mutation-test-generation skill, the agent MUST NOT create or modify production code. The agent MAY only create or modify test code and test fixtures (`*_test.go`, files under `testdata/`, and other clearly test-only fixtures), plus the mutation report/cache artifacts the skill documents. If a problem prevents killing a mutant without changing production code (including suspected bug, untestable design, missing fixtures, or would-require-production-change), the agent MUST stop, report the blocker with file path, line (or function/region), and reason, and MUST wait for the user to decide.

#### Scenario: Would need production change
- **WHEN** killing a live mutant would require changing production code
- **THEN** the agent stops and reports the file, line (or region), and why instead of editing production code

#### Scenario: Allowed test artifacts only
- **WHEN** the skill adds or updates tests to kill mutants
- **THEN** source edits are limited to `*_test.go` and test fixtures (e.g. under `testdata/`)
