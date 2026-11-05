## 1. Black-box skill

- [x] 1.1 Update `.cursor/skills/black-box-test-generation/SKILL.md` Hard exclusions / workflow to require test-only edits (`*_test.go`, `testdata/` / test fixtures) and forbid production code changes
- [x] 1.2 Add stop-and-report guidance (file, line/region, reason) for unreachable API behavior, suspected bugs, untestable design, or missing fixtures; include WHEN/THEN-style examples matching the delta scenarios

## 2. White-box skill

- [x] 2.1 Update `.cursor/skills/white-box-test-generation/SKILL.md` Hard exclusions and iteration “Stop” rules so coverage gaps that need production seams or other production edits are reported as blockers, not fixed by editing non-test source
- [x] 2.2 Align allowed levers and blocker report shape with the white-box delta (file, line/region, reason categories)

## 3. Verify scope

- [x] 3.1 Confirm `static-analysis-fix` skill and `openspec/specs/static-analysis-fix/spec.md` are unchanged
- [x] 3.2 Spot-check both skills mention waiting for the user when blocked (no silent production edits)
