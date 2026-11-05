## Context

`black-box-test-generation` and `white-box-test-generation` already forbid changing models and Docker native deps, and discourage recognition-pipeline edits. In practice, white-box runs still introduced production seams so coverage could hit Close/Encode/`os.Exit` paths. `static-analysis-fix` remains a separate workflow that is supposed to edit production code. This change tightens only the two test-generation skills and their specs.

## Goals / Non-Goals

**Goals:**

- Require test-generation skills to edit only test artifacts (`*_test.go`, fixtures, `testdata/`)
- Require stop-and-report (file, line, why) when the objective cannot be met without production changes—including unreachable code, suspected bugs, untestable design, or missing fixtures the skill cannot create under the guardrail
- Keep `static-analysis-fix` unchanged

**Non-Goals:**

- Reverting prior production seams already in the tree
- Changing analyzer/lint workflows or recognition behavior
- Defining mutation-testing policy

## Decisions

1. **Guardrail as ADDED requirement + MODIFIED hard-exclusion text**  
   Add an explicit “test-only edits / stop-and-report” requirement with scenarios, and update existing hard-exclusion requirements so they state that *any* production source change is out of scope for these skills (not only recognition/models/Docker).  
   *Alternatives:* skill-only prose without spec updates (rejected—user asked for specs); shared third capability (rejected—overkill for two skills).

2. **Allowed edit surface**  
   Allowed: `*_test.go`, test helpers colocated as tests, fixtures under `testdata/` (or other clearly test-only fixture paths), and generated coverage artifacts the skill already documents. Disallowed: non-test `.go` under `cmd/` / `internal/`, `models/`, Dockerfile/native deps, and other production assets.  
   *Alternatives:* allow “minimal test seams” in production (rejected—that is what this change forbids).

3. **Stop-and-report shape**  
   When blocked, the agent MUST stop the skill loop, report at least file path, line (or function/region), and reason category (unreachable / suspected bug / untestable design / missing fixtures / would require production change), and MUST NOT proceed with production edits unless the user explicitly directs a different workflow.  
   *Alternatives:* auto-open a follow-on OpenSpec change (deferred—user decides).

4. **White-box capability location**  
   `white-box-test-generation` is not yet under `openspec/specs/` (still in the open skill change). This change still ships a delta under `specs/white-box-test-generation/` so archive can merge after that capability is mainlined. Apply updates the live `SKILL.md` immediately.

5. **static-analysis-fix untouched**  
   No delta, no skill edit—explicit non-goal so agents do not over-generalize the guardrail.

## Risks / Trade-offs

- [Risk] Coverage or EP goals stall more often without seams → Mitigation: clear blocker reports; user may authorize production changes outside these skills  
- [Risk] Ambiguity whether package-level `export_test.go` or build-tag test hooks count as tests → Mitigation: treat `*_test.go` and `testdata/` as allowed; anything else requires stop-and-report  
- [Risk] White-box mainline spec merge order with the open skill change → Mitigation: ADDED/MODIFIED deltas are additive; archive after white-box skill change when possible

## Migration Plan

1. Apply: update both `SKILL.md` files and land delta specs via tasks  
2. No runtime deploy; no rollback beyond reverting the skill/spec docs  
3. Optional later: archive this change (after or with white-box skill archive sync)

## Open Questions

- None blocking apply; user may later clarify whether `export_test.go` naming is preferred when they authorize a deliberate test seam outside these skills.
