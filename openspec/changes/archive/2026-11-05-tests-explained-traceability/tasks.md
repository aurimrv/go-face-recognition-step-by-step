## 1. Registry and human document

- [x] 1.1 Add repo-root `tests-explained.yaml` with schema comment/header and empty or seeded `cases` list (required fields per design)
- [x] 1.2 Add initial `TESTS_EXPLAINED.md` generated from the YAML (package → file → test layout; “generated from YAML” banner)
- [x] 1.3 Document the upsert + regenerate procedure clearly enough that skills can follow it without a mandatory script (optional tiny helper only if needed)

## 2. Skill updates

- [x] 2.1 Update `.cursor/skills/black-box-test-generation/SKILL.md` to require registry upsert + `TESTS_EXPLAINED.md` regeneration for each new/changed case
- [x] 2.2 Update `.cursor/skills/white-box-test-generation/SKILL.md` likewise (include coverage locus in objective/refs when relevant)
- [x] 2.3 Update `.cursor/skills/mutation-test-generation/SKILL.md` likewise (include mutant identity when known)

## 3. Verify

- [x] 3.1 Spot-check skills against `specs/tests-explained` and the three skill delta requirements
- [x] 3.2 Optionally seed a few existing representative tests into the YAML/MD to demonstrate the format
