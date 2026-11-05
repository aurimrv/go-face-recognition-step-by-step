## ADDED Requirements

### Requirement: Update tests-explained for coverage cases
While executing the white-box-test-generation skill, after creating or materially changing a test case the agent MUST upsert a corresponding entry in `tests-explained.yaml` (skill: `white-box-test-generation`) with objective, system under test, and production refs (including uncovered loci when applicable), and MUST regenerate `TESTS_EXPLAINED.md` from the registry before finishing the pass.

#### Scenario: New coverage case documented
- **WHEN** the skill adds a test to cover previously uncovered statements
- **THEN** `tests-explained.yaml` gains or updates an entry for that test and `TESTS_EXPLAINED.md` includes it

#### Scenario: Objective and skill recorded
- **WHEN** the entry is written
- **THEN** it states why the case exists (coverage gap closed), names skill `white-box-test-generation`, and identifies the production code exercised
