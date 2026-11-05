## Why

Coverage and EP/BVA skills improve execution and design, but they do not prove tests assert enough to catch faults. This teaching fork already ships `gomutants`; a project skill is needed so the agent systematically adds tests aimed at **live** mutants (or a named mutant), finishes with JSON + HTML reports, and justifies likely **equivalent** survivors—without rewriting production code.

## What Changes

- Add Cursor skill `.cursor/skills/mutation-test-generation/SKILL.md`
- Workflow: run `gomutants` → inspect live mutants → add/adjust **tests only** to kill them (all live mutants or a user-specified mutant) → re-run → emit reports including `mutation-report.html` via `gomutants -html-output`
- Require stop-and-report (file, line/region, why) when a mutant cannot be killed without production changes; justify possible equivalent mutants that remain
- Align with black-box/white-box **test-only** guardrails (no production edits under this skill)

## Capabilities

### New Capabilities

- `mutation-test-generation`: Project skill for gomutants-guided test generation against live mutants, with HTML/JSON reports and equivalent-mutant justification

### Modified Capabilities

- (none)

## Impact

- New skill under `.cursor/skills/mutation-test-generation/`
- Generated artifacts typically include `mutation-report.json`, `.gomutants-cache.json`, and `mutation-report.html` (treat as generated unless the user asks to commit)
- Complements existing `black-box-test-generation` and `white-box-test-generation` skills; uses preinstalled `gomutants` (no new Docker/native dependency install as part of skill execution)
