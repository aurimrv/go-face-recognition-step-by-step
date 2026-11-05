## ADDED Requirements

### Requirement: Update tests-explained for generated cases
While executing the black-box-test-generation skill, after creating or materially changing a test case the agent MUST upsert a corresponding entry in `tests-explained.yaml` (skill: `black-box-test-generation`) with objective, system under test, and production refs, and MUST regenerate `TESTS_EXPLAINED.md` from the registry before finishing the pass.

#### Scenario: New EP/BVA case documented
- **WHEN** the skill adds a new test function for a chosen partition or boundary
- **THEN** `tests-explained.yaml` gains or updates an entry for that test and `TESTS_EXPLAINED.md` includes it

#### Scenario: Objective and skill recorded
- **WHEN** the entry is written
- **THEN** it states why the case exists, names skill `black-box-test-generation`, and identifies the unit/behavior under test
