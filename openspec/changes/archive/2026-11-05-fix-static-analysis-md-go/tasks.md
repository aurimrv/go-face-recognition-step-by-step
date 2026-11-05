## 1. Go fixes

- [x] 1.1 Fix errcheck on `Close` in `internal/entity/drawer.go` (loadImage + SaveImage)
- [x] 1.2 Fix errcheck on `dir.Close` and failed-Open path in `internal/usecases/load_persons.go`
- [x] 1.3 Lowercase ST1005 error strings in `cmd/main.go` `processJob`

## 2. Markdown fixes

- [x] 2.1 Fix `README.md` (MD013, MD029, MD040, MD051, and related)
- [x] 2.2 Fix `DEVCONTAINER_SETUP.md` (MD013, MD031, MD060, and related)
- [x] 2.3 Fix `01-TERMINAL_INTERACTION.md` structure and line length
- [x] 2.4 Fix `01-COPILOT_INTERACTION.md` structure and line length

## 3. Verify

- [x] 3.1 Re-run markdownlint-cli2 and go vet / golangci-lint; apply a second fix round if needed; stop after two rounds

## Deferred (round limit — id + rationale)

- `markdownlint:MD013:01-TERMINAL_INTERACTION.md:*` (76 findings) — File is a verbatim fenced terminal transcript for the class; wrapping lines would falsify the captured lint dump. Added an H1 for MD041 only.
- `markdownlint:MD013:01-COPILOT_INTERACTION.md:57-58,172-175,271` (7 findings) — Remaining long table rows in the exported chat; converting them mid-round would churn the teaching export further after the two-round budget.
