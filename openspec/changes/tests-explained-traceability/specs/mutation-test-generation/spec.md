## ADDED Requirements

### Requirement: Update tests-explained for mutant-kill cases
While executing the mutation-test-generation skill, after creating or materially changing a test case the agent MUST upsert a corresponding entry in `tests-explained.yaml` (skill: `mutation-test-generation`) with objective, system under test, production refs, and mutant identity when known, and MUST regenerate `TESTS_EXPLAINED.md` from the registry before finishing the pass.

#### Scenario: New kill case documented
- **WHEN** the skill adds a test aimed at killing a live mutant
- **THEN** `tests-explained.yaml` gains or updates an entry for that test and `TESTS_EXPLAINED.md` includes it

#### Scenario: Objective and skill recorded
- **WHEN** the entry is written
- **THEN** it states why the case exists (which mutant/behavior), names skill `mutation-test-generation`, and identifies the production code under test
