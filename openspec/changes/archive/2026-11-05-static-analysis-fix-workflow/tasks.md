## 1. Project skill

- [x] 1.1 Create `.cursor/skills/static-analysis-fix/SKILL.md` with YAML frontmatter (`name: static-analysis-fix`, description covering Markdown/Go static-analysis fix via OpenSpec)
- [x] 1.2 Document the workflow steps: ingest report or run analyzers → triage/group → OpenSpec change (proposal, design, specs, tasks) → apply task-by-task → re-run → archive when clean
- [x] 1.3 Document analyzer commands: markdownlint-cli2 `"**/*.md" "#openspec" "#.cursor" "#node_modules"`; go vet; golangci-lint with existing `.golangci.yml` (v2); Docker/Compose fallback for CGO
- [x] 1.4 Document fix policies: preserve behavior; prefer real fixes; no blanket nolint/markdownlint-disable; deferrals require finding id + rationale; max two fix/re-run rounds then report leftovers
- [x] 1.5 Document hard exclusions: do not change recognition pipeline, models, or Docker native dependencies to silence findings
- [x] 1.6 State explicitly that agents MUST NOT create a new `.golangci.yml`

## 2. Cursor command

- [x] 2.1 Create `.cursor/commands/opsx-static-fix.md` with frontmatter name `/opsx-static-fix`, id `opsx-static-fix`, category Workflow, and a short description
- [x] 2.2 Command body MUST instruct the agent to read and follow `.cursor/skills/static-analysis-fix/SKILL.md` (same thin-wrapper style as other opsx commands)

## 3. Sanity check

- [x] 3.1 Confirm no new golangci config file was added and existing `.golangci.yml` is unchanged by this change
- [x] 3.2 Spot-check that skill + command paths match the specs and that OpenSpec artifact set for this meta-change remains complete
