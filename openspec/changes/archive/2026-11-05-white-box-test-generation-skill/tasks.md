## 1. Skill scaffold

- [x] 1.1 Create `.cursor/skills/white-box-test-generation/SKILL.md` with YAML frontmatter (`name: white-box-test-generation`, description covering coverage-driven Go test generation)
- [x] 1.2 Document when to use (raise statement coverage, white-box tests, close gaps after black-box EP/BVA)

## 2. Measurement and loop

- [x] 2.1 Document `go test -coverprofile` measurement, including whole-project recipe with `-coverpkg=./...` when scope is the module
- [x] 2.2 Document generating `coverage.html` via `go tool cover -html=… -o coverage.html` and using `go tool cover -func` as an aid
- [x] 2.3 Document default threshold 100% and default scope entire project; document explicit threshold/file/package overrides
- [x] 2.4 Document the iterate loop: measure → inspect uncovered → add tests → re-measure → stop at threshold or report blockers; include a practical round cap

## 3. Policies

- [x] 3.1 Note relationship to `black-box-test-generation` (design first optional; white-box closes gaps; annotate new cases with uncovered locus)
- [x] 3.2 Document hard exclusions: no recognition-pipeline/models/Docker-native changes solely to raise coverage
- [x] 3.3 Spot-check skill path and policies against `specs/white-box-test-generation/spec.md`
