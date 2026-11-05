## Context

Three test-generation skills already annotate `*_test.go` differently (EP/BVA catalogs, `// Covers:`, `// Kills:`). Reviewers and students still lack one place that answers: objective, generating skill, and production responsibility, kept current as code evolves. This change adds a dual human/machine documentation surface and wires all three skills to maintain it.

## Goals / Non-Goals

**Goals:**

- Define a durable layout for `TESTS_EXPLAINED.md` that humans can skim by package/file
- Define a YAML registry skills can upsert without fragile Markdown editing
- Require black-box, white-box, and mutation skills to update registry + regenerate Markdown whenever they add or materially change a test case
- Keep entries traceable to current test symbols and production units (paths / APIs)

**Non-Goals:**

- Auto-backfilling every historical test in the first apply (optional follow-on / best-effort seed)
- Replacing in-test comments (they remain; explanations are complementary)
- Changing production code or test-only guardrails
- Building a full CI doc site

## Decisions

1. **YAML as source of truth; Markdown as projection**  
   - Machine file: `tests-explained.yaml` at repo root  
   - Human file: `TESTS_EXPLAINED.md` regenerated from YAML after each skill pass  
   Skills **MUST NOT** hand-edit Markdown as the primary store (avoids merge conflicts and drift).  
   *Alternatives:* Markdown-only (hard for agents); JSON (less comment-friendly); XML (verbose). YAML chosen for readability + stable keys.

2. **Stable case identity**  
   Primary key: `id` = `{go_package_path}::{TestFunction}` (e.g. `github.com/andre-ols/go-face-recognition/internal/entity::TestDrawFace_TopEdgePixels`).  
   On rename/delete of a test, skills update or remove that id so the catalog matches the suite.

3. **Required fields per case**

   | Field | Purpose |
   |---|---|
   | `id` | Stable key |
   | `test` | `Test…` function name |
   | `test_file` | Path to `*_test.go` |
   | `skill` | `black-box-test-generation` \| `white-box-test-generation` \| `mutation-test-generation` \| `manual` |
   | `objective` | One short paragraph: why the case exists |
   | `system_under_test` | Package and symbol(s) / behavior under test |
   | `production_refs` | List of production file paths (optional line/region hints) |

   Optional: `partitions`, `boundaries`, `covers`, `kills_mutant`, `notes`, `updated_at` (ISO date).

4. **Markdown organization (human)**

   ```
   # Tests explained
   ## How to read this document
   ## Index by package
   ## Package `internal/entity`
   ### File `drawer_test.go`
   #### `TestDrawFace_TopEdgePixels`
   - Skill: …
   - Objective: …
   - System under test: …
   - Production refs: …
   - Trace: id=…
   ```

   Packages sorted alphabetically; files and tests sorted by name. Header notes that the file is generated from `tests-explained.yaml`.

5. **Skill workflow hook**  
   After writing/updating `*_test.go`, each test skill:  
   1) upsert YAML entries for touched tests;  
   2) remove entries whose `test_file`+`test` no longer exist if the skill deleted them;  
   3) regenerate `TESTS_EXPLAINED.md` from YAML (inline in the skill, or a tiny documented script).  
   Apply phase will document the exact regeneration procedure in skills (agent-driven is enough initially).

6. **Traceability with current implementation**  
   Skills MUST set `system_under_test` and `production_refs` from the code they actually exercise. When production paths move, the next skill pass that touches related tests SHOULD refresh refs. Periodic consistency check: every `id`’s `test_file` exists and contains `func {test}`.

## Risks / Trade-offs

- [Risk] YAML and MD diverge if someone edits MD only → Mitigation: skills overwrite MD from YAML; banner in MD says “do not edit by hand”  
- [Risk] Catalog grows stale for untouched old tests → Mitigation: best-effort seed on apply; skills refresh on touch; optional later audit task  
- [Risk] Verbose YAML → Mitigation: keep required fields minimal; optionals only when useful  
- [Risk] Merge conflicts on YAML → Mitigation: keyed by stable `id`; small per-case patches

## Migration Plan

1. Apply: add empty/seeded YAML + initial MD; update three skills  
2. Optionally seed entries for existing high-value tests  
3. Archive; future skill runs keep docs current  
4. Rollback: remove the two files and skill sections

## Open Questions

- None blocking; whether to add a `scripts/render-tests-explained.go` can be decided at apply (agent rewrite of MD is acceptable for v1).
