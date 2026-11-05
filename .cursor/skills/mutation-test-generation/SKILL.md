---
name: mutation-test-generation
description: >-
  Generate or improve Go tests guided by live mutants using gomutants: kill all
  live mutants in scope or a specific mutant, write mutation-report.html, and
  justify possible equivalent survivors. Use when the user asks for mutation
  testing, killing survivors, or a HTML/JSON mutation report.
---

# Mutation Test Generation (gomutants / live mutants)

Improve the suite by **killing live mutants**. Measure with **`gomutants`** →
inspect survivors → add or tighten **tests only** → re-run → finish with JSON +
**`mutation-report.html`**, and justify possible **equivalent** mutants that
remain.

This complements `.cursor/skills/black-box-test-generation/` (EP/BVA design) and
`.cursor/skills/white-box-test-generation/` (statement coverage). Prefer those
first when the suite is thin or coverage is low; use this skill when mutants
still live despite coverage.

**MUST** use **`gomutants`** only (not gremlins or other mutators).

## When to use

- User asks for mutation testing, killing live/surviving mutants, or
  gomutants-guided test generation
- User names a **specific mutant** (file, line, and/or id from a report)
- User wants `mutation-report.html` / JSON mutation reports
- Follow-on after black-box or white-box when assertions look weak

## Hard exclusions (test-only)

While executing this skill, **MUST NOT** create or modify **production code**
(non-test `.go` under `cmd/` / `internal/` or other app sources)—including
recognition pipeline edits, injecting test seams, or “fixing” production to
make a mutant die—nor model weights under `models/`, nor Docker native
dependency installation.

**Allowed levers only:**

- `*_test.go` (and test helpers in test files)
- Test fixtures under `testdata/` or other clearly test-only fixture paths
- Mutation artifacts this skill documents (`mutation-report.json`,
  `mutation-report.html`, `.gomutants-cache.json`)

If a mutant cannot be killed with tests/fixtures alone, **stop-and-report**—do
not edit production under this skill.

## Stop-and-report (wait for the user)

When blocked, **STOP** the kill loop, report, and **wait for the user**. Do not
silently change production to kill a mutant.

Each blocker report **MUST** include:

- **File** path
- **Line** (or function/region)
- **Reason** category: `unreachable` | `suspected bug` | `untestable design` |
  `missing fixtures` | `would require production change` | `possible equivalent`

Examples:

- **WHEN** killing a live mutant would require changing production code **THEN**
  stop and report file/line/why (`would require production change`).
- **WHEN** no observable oracle can distinguish the mutant from the original
  **THEN** treat as possible equivalent, justify it, and do not claim killed.

## Scope resolution

Resolve before the first `gomutants` run:

| User input | Scope |
|---|---|
| Package/path (e.g. `./internal/usecases`) | That package (or listed packages) |
| Specific mutant (file/line/id) | Package containing that locus; prioritize that mutant |
| Unspecified | Prefer the package under discussion; otherwise ask or start with a small **pure** package (path helpers). Avoid blindly mutating the whole module (CGO timeouts). |

Explicit user package/path always wins.

## Measurement recipe

Default artifact names (unless the user overrides):

- JSON: `mutation-report.json` (`gomutants` default / `-output`)
- HTML: **`mutation-report.html`** via **`-html-output mutation-report.html`**
- Cache: `.gomutants-cache.json` (optional; leave enabled unless user disables)

Example:

```bash
gomutants ./path/to/package \
  -output mutation-report.json \
  -html-output mutation-report.html
```

Optional aid when justifying survivors (does **not** replace written
justification):

```bash
gomutants ./path/to/package \
  -detect-equivalent \
  -output mutation-report.json \
  -html-output mutation-report.html
```

Inspect live mutants from the console summary and/or `mutation-report.json`.
Treat report/cache files as **generated**; do not commit unless asked.

Timeouts on CGO-heavy packages are expected in this teaching fork; prefer
narrower scopes and note timeouts in the summary when they dominate.

## Kill loop

Cap at **10** fix/measure rounds. After the 10th re-run, stop and report
leftovers even if live mutants remain.

1. **Resolve** scope (and optional specific mutant).
2. **Measure** with `gomutants` (always include HTML output for the final pass;
   include it on intermediate passes when useful).
3. **Inspect** live/surviving mutants (file, line, operator/description).
4. **Add or tighten tests** — assert observable outcomes that the mutant should
   break:
   - Prefer focused cases; annotate intent, e.g. `// Kills: mutant at file:line …`
   - Preserve existing EP/BVA catalog comments when present.
   - If the next step would be a production edit, **stop-and-report**.
5. **Re-run** `gomutants` on the same scope.
6. **Stop** when:
   - target met (no live mutants in scope, or the specific mutant is killed); or
   - blocked / possible equivalents remain → justify and wait if production
     change would be needed; or
   - round cap hit → list remaining live mutants + justifications.

### All live mutants vs specific mutant

- **All live (default):** Work through survivors in the report; kill as many as
  possible with test/fixture edits; re-measure after batches of related tests.
- **Specific mutant:** Prioritize tests for that file/line/operator until it is
  **killed**, justified as **possible equivalent**, or **blocked**
  (stop-and-report). Still emit final JSON + HTML reports for the scope.

## Equivalent and unkillable survivors

At the end of the run, for each remaining **live** mutant the agent **MUST**
document:

- File and line (or region)
- Brief justification: e.g. possible **equivalent** (no observable behavior
  change / no oracle), timeout-only noise, or stop-and-report reason

Optional `-detect-equivalent` (Trivial Compiler Equivalence) may **inform** the
justification; it **MUST NOT** replace written reasoning. Do not mark a mutant
killed if it only looks equivalent.

## Relationship to other skills

- Optional order: **black-box** (design) → **white-box** (coverage) → **this
  skill** (assertion strength / mutants).
- Do not strip EP/BVA comments while adding kill-focused cases.
- White-box statement coverage is out of scope here except as context.

## Tests explained registry (required)

After creating or **materially changing** any test case, keep human/machine
traceability current:

1. **Upsert** `tests-explained.yaml` (repo root) for each touched test:
   - `id`: `{module_import_path}::{TestFunction}`  
     (module = `github.com/andre-ols/go-face-recognition`)
   - `test`, `test_file`, `skill: mutation-test-generation`
   - `objective` — why the case exists (mutant / assertion intent)
   - `system_under_test` — package/symbol/behavior
   - `production_refs` — production paths under test
   - Optional: `kills_mutant` (file:line / type), `updated_at`
2. If you **delete** a test you registered, remove its YAML entry.
3. **Regenerate** Markdown:

```bash
python3 scripts/render-tests-explained.py
```

Do not hand-edit `TESTS_EXPLAINED.md` as the primary store.

## Output expectations

Report:

- Scope (and specific mutant, if any)
- Rounds executed
- Final gomutants summary (killed / lived / timed out / efficacy if shown)
- Paths to `mutation-report.json` and `mutation-report.html`
- Tests/files added or updated (`*_test.go` / fixtures only)
- Confirmation that `tests-explained.yaml` + `TESTS_EXPLAINED.md` were updated
- Remaining live mutants with file, line/region, and equivalent/unkillable
  justification—and that you are waiting on the user before any production
  change when blocked
