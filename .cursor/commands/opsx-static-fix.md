---
name: /opsx-static-fix
id: opsx-static-fix
category: Workflow
description: Resolve Markdown and Go static analysis findings via OpenSpec
---

Resolve static analysis findings using the project **static-analysis-fix** skill.

**Input**: Optionally a path to a report file, pasted analyzer output, or a focus area (e.g. Markdown only, a package path). If omitted, run the configured analyzers on the project.

**Steps**

1. **Read and follow the skill**

   Open and follow `.cursor/skills/static-analysis-fix/SKILL.md` for the full procedure, analyzer commands, fix/deferral policies, round limit, and hard exclusions.

2. **Execute the skill workflow**

   Ingest a report or run analyzers → triage/group → create an OpenSpec change (proposal, design, specs, tasks) → apply fixes task by task → re-run analyzers → archive when clean, or stop after two fix/re-run rounds and report what remains.

3. **Do not bypass skill policies**

   - Use existing `.golangci.yml` (v2) only — do not create a new golangci config
   - Markdown: `markdownlint-cli2 "**/*.md" "#openspec" "#.cursor" "#node_modules"`
   - Prefer real fixes; no blanket `nolint` / `markdownlint-disable`
   - Record every deferral with finding id and rationale
   - Do not change the recognition pipeline, models, or Docker native dependencies to silence findings
