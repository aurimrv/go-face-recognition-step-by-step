# static-analysis-fix Specification

## Purpose
TBD - created by archiving change static-analysis-fix-workflow. Update Purpose after archive.
## Requirements
### Requirement: Project skill and command exist
The project SHALL provide a Cursor skill at `.cursor/skills/static-analysis-fix/SKILL.md` and a command at `.cursor/commands/opsx-static-fix.md` that instructs the agent to follow that skill for resolving static analysis findings.

#### Scenario: Command entry point
- **WHEN** the user invokes `/opsx-static-fix`
- **THEN** the agent reads and follows the `static-analysis-fix` project skill

#### Scenario: Skill discoverable in repo
- **WHEN** a contributor inspects `.cursor/skills/` and `.cursor/commands/`
- **THEN** both `static-analysis-fix` and `opsx-static-fix` are present in the project tree

### Requirement: Spec-driven fix workflow
The skill SHALL direct the agent to: (1) ingest a report or run the configured analyzers, (2) triage and group findings, (3) create an OpenSpec change with proposal, design, specs, and tasks, (4) apply fixes task by task, (5) re-run the analyzers, and (6) archive the change when findings are clean, subject to the round limit.

#### Scenario: End-to-end when clean after fixes
- **WHEN** a fix round clears all non-deferred findings on re-run
- **THEN** the agent archives the OpenSpec change after reporting success

#### Scenario: Ingest supplied report
- **WHEN** the user supplies a static analysis report instead of asking to run tools
- **THEN** the agent triages from that report without requiring a fresh analyzer run before the first OpenSpec change is created

### Requirement: Markdown analyzer scope
When running Markdown analysis, the agent SHALL invoke markdownlint-cli2 with exactly the project scope `"**/*.md" "#openspec" "#.cursor" "#node_modules"` and MUST NOT lint excluded trees as part of this workflow.

#### Scenario: Markdown lint command
- **WHEN** the agent runs Markdown static analysis under this skill
- **THEN** it uses `markdownlint-cli2 "**/*.md" "#openspec" "#.cursor" "#node_modules"`

### Requirement: Go analyzers and existing golangci config
When running Go analysis, the agent SHALL run `go vet` and golangci-lint (including staticcheck as configured) using the repository’s existing `.golangci.yml` (v2). The agent MUST NOT create a new golangci configuration file.

#### Scenario: Use existing config only
- **WHEN** the agent runs Go static analysis under this skill
- **THEN** it uses the existing `.golangci.yml` and does not add a replacement or parallel golangci config

#### Scenario: Docker fallback for CGO
- **WHEN** host-side Go analysis fails because of CGO or go-face native dependencies
- **THEN** the agent re-runs the Go analyzers via the project’s Docker/Compose app environment

### Requirement: Fix quality and suppression limits
Fixes SHALL preserve observable behavior. The agent SHALL prefer real fixes over suppressions. The agent MUST NOT introduce blanket `nolint` or blanket `markdownlint-disable` directives. Every deferred finding MUST be recorded with its finding id and a rationale in the OpenSpec change artifacts.

#### Scenario: Prefer real fix
- **WHEN** a finding can be resolved by editing code or Markdown without changing behavior
- **THEN** the agent applies that fix rather than adding a suppression

#### Scenario: Deferral recorded
- **WHEN** a finding is deferred
- **THEN** the OpenSpec change records the finding id and a written rationale

#### Scenario: No blanket disable
- **WHEN** the agent cannot fix a cluster of findings in one file
- **THEN** it MUST NOT add a file- or project-wide `nolint` or `markdownlint-disable` to silence them

### Requirement: Two-round stop
The agent SHALL perform at most two fix-and-re-run rounds per `/opsx-static-fix` invocation. After the second re-run, if findings remain, the agent MUST stop further automatic fixes and report what is left (including deferrals).

#### Scenario: Stop after second re-run
- **WHEN** findings remain after the second re-run
- **THEN** the agent stops the fix loop and reports remaining findings without starting a third fix round

### Requirement: Hard exclusions
While executing this skill, the agent MUST NOT change the recognition pipeline, model files, or Docker native dependency installation in order to clear findings; such items SHALL be deferred with id and rationale if they surface as findings.

#### Scenario: Recognition code out of scope for silencing
- **WHEN** a finding appears to require changing face-recognition use cases, models, or Dockerfile native packages solely to silence the analyzer
- **THEN** the agent defers that finding with id and rationale instead of making those changes

