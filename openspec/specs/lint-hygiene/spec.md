# lint-hygiene Specification

## Purpose
TBD - created by archiving change fix-static-analysis-md-go. Update Purpose after archive.
## Requirements
### Requirement: Go packages pass configured golangci-lint
Go packages in this repository SHALL pass `go vet ./...` and `golangci-lint run` using the existing `.golangci.yml` without introducing blanket nolint directives.

#### Scenario: Clean Go analysis
- **WHEN** a developer runs `go vet ./...` and `golangci-lint run` with the repo `.golangci.yml`
- **THEN** there are no reported issues in project Go packages

### Requirement: Project Markdown passes markdownlint-cli2 scope
Project Markdown files outside `openspec/`, `.cursor/`, and `node_modules/` SHALL pass `markdownlint-cli2 "**/*.md" "#openspec" "#.cursor" "#node_modules"` without file-wide markdownlint-disable directives.

#### Scenario: Clean Markdown analysis
- **WHEN** a developer runs the project markdownlint-cli2 command
- **THEN** there are no reported issues in the linted project Markdown files

