# black-box-test-generation Specification

## Purpose
TBD - created by archiving change black-box-test-generation-skill. Update Purpose after archive.
## Requirements
### Requirement: Project skill exists
The project SHALL provide a Cursor skill at `.cursor/skills/black-box-test-generation/SKILL.md` that instructs the agent to generate Go tests using equivalence partitioning (EP) and boundary value analysis (BVA).

#### Scenario: Skill discoverable
- **WHEN** a contributor inspects `.cursor/skills/black-box-test-generation/`
- **THEN** a `SKILL.md` with name `black-box-test-generation` is present

#### Scenario: Skill invoked for test generation
- **WHEN** the user asks to generate black-box EP/BVA tests for a Go package or function in this repo
- **THEN** the agent follows the `black-box-test-generation` skill

### Requirement: Equivalence partitions documented in the test set
Generated Go test files SHALL include comments that define the equivalence partitions used for the unit under test (stable ids and brief meanings).

#### Scenario: Partition catalog present
- **WHEN** the skill generates or updates a `*_test.go` file for a unit
- **THEN** that file contains a comment block listing the EP ids and their definitions for that unit

### Requirement: Each test case cites partition or boundary coverage
Each generated test case SHALL state, in a comment and/or subtest name, which equivalence partition and/or boundary it covers.

#### Scenario: Per-case coverage annotation
- **WHEN** a generated test function or table row exercises a representative or boundary input
- **THEN** it is annotated with the corresponding partition id and/or boundary id

### Requirement: BVA included alongside EP
The skill SHALL derive boundary values (edges of partitions and notable numeric/string/path edges) and SHALL include test cases for selected boundaries, not only one interior representative per partition.

#### Scenario: Boundary cases present
- **WHEN** the unit under test has identifiable boundaries (e.g. empty vs non-empty args, threshold edges, JPEG vs non-JPEG)
- **THEN** the generated suite includes at least one annotated case per selected boundary the skill chooses to cover

### Requirement: Face and unknown image fixtures when match degree matters
When generating tests for face recognition or classification behavior where match strength matters, the skill SHALL create or reuse face images and unknown images (including similar-face unknowns) so the suite can exercise different degrees of matching, and SHALL document fixture intent in comments.

#### Scenario: Classification-oriented fixtures
- **WHEN** the unit under test classifies unknown faces against known samples with a similarity/threshold notion
- **THEN** the skill plans fixtures for strong match, weaker/near-boundary match, and non-match (or no-face) as needed for the chosen partitions/boundaries, with comments describing each fixture’s role

#### Scenario: Prefer hermetic fixture location
- **WHEN** new image fixtures are created for tests
- **THEN** they are placed under a test-oriented location such as `testdata/` (or copied into temp dirs at test time) rather than replacing model files

### Requirement: Black-box scope and hard exclusions
The skill SHALL target observable inputs and outputs of the unit under test. While executing the skill, the agent MUST NOT change production implementation (including recognition pipeline code), model weights, or Docker native dependency installation in order to make tests pass. Allowed levers are tests and fixtures only, subject to Requirement: Test-only edits and stop-and-report.

#### Scenario: Tests adapt to behavior, code under test stays
- **WHEN** a generated test fails because production behavior differs from an incorrect expectation
- **THEN** the agent corrects the test or fixtures under EP/BVA rationale and does not modify production code, models, or Docker native deps to force a pass

#### Scenario: Stop instead of production workaround
- **WHEN** correcting expectations or fixtures is insufficient and production code would need to change for the suite to meet its objective
- **THEN** the agent stops and reports the blocker rather than editing production code under this skill

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
