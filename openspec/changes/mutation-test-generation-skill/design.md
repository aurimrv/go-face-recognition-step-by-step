## Context

The repo already standardizes on **`gomutants`** (preinstalled in the teaching container) for mutation testing, alongside `go test -cover` and the black-box / white-box skills. High statement coverage can still leave live mutants when assertions are weak. This change adds a skill that drives test updates from **surviving (live) mutants**, then publishes machine- and human-readable reports.

## Goals / Non-Goals

**Goals:**

- Ship `.cursor/skills/mutation-test-generation/SKILL.md` with a clear gomutants-first loop
- Target **all live mutants** in scope, or a **user-specified** mutant/locus
- End with JSON report + **`mutation-report.html`** via `-html-output`
- Justify remaining survivors that appear **equivalent** (or otherwise unkillable under test-only rules)
- Enforce the same **test-only / stop-and-report** guardrail as other test-generation skills

**Non-Goals:**

- Replacing EP/BVA or coverprofile skills
- Installing alternate mutators (gremlins, etc.) or changing Docker native deps
- Guaranteeing 100% kill rate when equivalents remain
- Editing production code under the skill to “kill” mutants

## Decisions

1. **Tooling: `gomutants` only**  
   Skill documents `gomutants <packages…>` with `-output` / default JSON and **`-html-output mutation-report.html`**. Optional `-detect-equivalent` may be used as an aid when justifying survivors (Trivial Compiler Equivalence), but human justification is still required.  
   *Alternatives:* gremlins (rejected—class standardized on gomutants).

2. **Scope resolution**  
   Default to user-named package/path; if unspecified, prefer a focused package the user was discussing, else ask or use a small pure package (e.g. usecases path helpers) rather than blindly mutating the entire module (CGO/timeouts). Explicit user package/path always wins.

3. **Kill loop**  
   Run gomutants → parse live/survived mutants from console and/or `mutation-report.json` → write tests that assert observable differences → re-run. Cap practical rounds (e.g. 10) like white-box. For a **specific** mutant, prioritize tests for that file/line/operator until killed or blocked/equivalent.

4. **Test-only guardrail**  
   Allowed: `*_test.go`, `testdata/` fixtures, and gomutants report/cache artifacts the skill documents. Disallowed: production `.go`, models, Docker native install. Blocked mutants → stop with file, line/region, reason; wait for user.

5. **Equivalent mutants**  
   After the kill loop (or when blocked), the final report section MUST list remaining live mutants with a short justification when they appear equivalent (semantics unchanged / no observable oracle) or cite stop-and-report reasons. Do not claim “killed” for equivalents.

6. **Artifacts**  
   Default names: `mutation-report.json` (gomutants default/`-output`), `mutation-report.html` (`-html-output mutation-report.html`), cache `.gomutants-cache.json`. Treat as generated unless user asks to commit.

## Risks / Trade-offs

- [Risk] Long runtimes / timeouts on CGO packages → Mitigation: scope packages; document expected timeouts; prefer pure packages first  
- [Risk] False “equivalent” labels → Mitigation: prefer evidence (TCE flag + oracle reasoning); leave as live if unsure and justify uncertainty  
- [Risk] Agent tempted to patch production to kill mutant → Mitigation: hard stop-and-report in skill + spec  
- [Risk] Report file churn in git → Mitigation: skill says do not commit reports unless asked

## Migration Plan

1. Apply: write `SKILL.md` per tasks  
2. No runtime deploy  
3. Archive after apply; optional later: point white-box “follow-on mutation” note at this skill

## Open Questions

- None blocking; exact worker/timeout flags can follow `DEVCONTAINER_SETUP.md` defaults unless the user overrides.
