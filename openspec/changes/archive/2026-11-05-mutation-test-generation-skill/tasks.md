## 1. Skill scaffold

- [x] 1.1 Create `.cursor/skills/mutation-test-generation/SKILL.md` with YAML frontmatter (`name: mutation-test-generation`, description covering gomutants-guided killing of live mutants)
- [x] 1.2 Document when to use (mutation testing, kill live mutants, specific mutant, HTML mutation report)

## 2. Workflow content

- [x] 2.1 Document scope resolution and the measure → inspect live mutants → add tests → re-run loop using `gomutants` only
- [x] 2.2 Document targeting all live mutants vs a user-specified mutant
- [x] 2.3 Document final reports: JSON (default/`-output`) and `mutation-report.html` via `-html-output mutation-report.html`
- [x] 2.4 Document justifying possible equivalent / unkillable survivors (file, line/region, reason; optional `-detect-equivalent` as aid only)

## 3. Guardrails and verify

- [x] 3.1 Document test-only edits and stop-and-report (file, line/region, why; wait for user)—aligned with other test-generation skills
- [x] 3.2 Document hard exclusions (no production code, models, or Docker native install under this skill)
- [x] 3.3 Spot-check skill path and policies against `specs/mutation-test-generation/spec.md`
