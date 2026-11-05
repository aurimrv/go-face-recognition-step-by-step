## Why

Recent white-box coverage work introduced production seams (`exitFunc`, `jpegEncode`, `closeFile`, `closePersonsDir`) so tests could hit hard-to-reach branches. Test-generation skills should not rewrite production code to chase coverage or EP/BVA goals; they should stop and let the user decide when the design or code under test blocks the objective. Codifying that now prevents the same drift in future skill runs while leaving `static-analysis-fix` free to change production code.

## What Changes

- Strengthen **black-box-test-generation** skill + spec: tests/fixtures only; stop-and-report when production changes would be required
- Strengthen **white-box-test-generation** skill + spec: same test-only guardrail (narrower than today’s “don’t change recognition/models/Docker” exclusions)
- Add explicit scenarios: WHEN production edits would be needed THEN stop with file/line/reason and wait for the user
- **Do not** change `static-analysis-fix` skill or spec (production edits remain expected there)

## Capabilities

### New Capabilities

- (none)

### Modified Capabilities

- `black-box-test-generation`: Add a test-only edit scope requirement and stop-and-report behavior when objectives cannot be met without production changes
- `white-box-test-generation`: Same test-only guardrail and stop-and-report behavior (capability lives in the open `white-box-test-generation-skill` change / project skill; delta applies here so archive can merge after that change)

## Impact

- `.cursor/skills/black-box-test-generation/SKILL.md`
- `.cursor/skills/white-box-test-generation/SKILL.md`
- OpenSpec specs (delta) for both capabilities
- No application/runtime Go packages, models, Docker, or `static-analysis-fix` skill
