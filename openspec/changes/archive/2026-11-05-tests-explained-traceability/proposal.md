## Why

Generated tests accumulate comments and skill-specific annotations, but humans still struggle to see **why** a case exists, **which skill** produced it, and **what production unit** it protects—especially as suites grow and code moves. A single, always-current explanation surface (with machine-friendly traceability) makes review and teaching coherent across black-box, white-box, and mutation skills.

## What Changes

- Introduce repo-root **`TESTS_EXPLAINED.md`** as the human-readable catalog of test objectives, generating skill, and system under test
- Introduce a machine source of truth (YAML) that skills upsert; regenerate `TESTS_EXPLAINED.md` from it so docs stay synchronized with the current suite/implementation
- Update **black-box-test-generation**, **white-box-test-generation**, and **mutation-test-generation** skills so every new/changed test case updates the registry and refreshed Markdown
- Define a coherent organization (by package → test file → test case) with stable IDs and optional links to production paths / EP-BVA / coverage / mutant refs
- Does **not** change production recognition logic or weaken test-only guardrails

## Capabilities

### New Capabilities

- `tests-explained`: Format, location, and maintenance rules for `TESTS_EXPLAINED.md` plus its YAML registry for human readability and automated traceability

### Modified Capabilities

- `black-box-test-generation`: Require updating the tests-explained registry/Markdown for each generated or materially changed test case
- `white-box-test-generation`: Same requirement for coverage-driven cases
- `mutation-test-generation`: Same requirement for mutant-kill cases

## Impact

- New files: `tests-explained.yaml` (registry), `TESTS_EXPLAINED.md` (generated/synced human doc); optional small helper script later if needed
- Skill updates under `.cursor/skills/{black-box,white-box,mutation}-test-generation/SKILL.md`
- OpenSpec deltas for the four capabilities above
- Existing `*_test.go` files unchanged by this proposal except when a later apply pass backfills explanations
