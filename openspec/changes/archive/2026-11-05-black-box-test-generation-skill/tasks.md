## 1. Skill scaffold

- [x] 1.1 Create `.cursor/skills/black-box-test-generation/SKILL.md` with YAML frontmatter (`name: black-box-test-generation`, description covering EP/BVA Go test generation for this repo)
- [x] 1.2 Document when to use the skill (user asks for black-box / EP / BVA tests for a package or API)

## 2. EP/BVA workflow and comment contract

- [x] 2.1 Document steps: identify observable I/O → derive EP → derive BVA → write catalog comments → emit Go tests → run `go test`
- [x] 2.2 Require suite-level comments defining partition ids and meanings
- [x] 2.3 Require each test case (function, subtest, or table row) to cite the partition and/or boundary it covers
- [x] 2.4 Include an illustrative comment/test shape in the skill (e.g. path-resolver style partitions)

## 3. Face-image fixtures and exclusions

- [x] 3.1 Document when to create/reuse known-face and unknown images (strong match, weak/near-threshold, non-match / no-face) with fixture-intent comments
- [x] 3.2 Prefer `testdata/` or temp copies; MUST NOT replace `models/` or change Docker native deps / recognition pipeline to force passes
- [x] 3.3 Note CGO/models availability (devcontainer) and that pure packages (e.g. `ResolveAnalysisJobs`) need no face fixtures

## 4. Sanity check

- [x] 4.1 Confirm skill path and required policies match `specs/black-box-test-generation/spec.md`
