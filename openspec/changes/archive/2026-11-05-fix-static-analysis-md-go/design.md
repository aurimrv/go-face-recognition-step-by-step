## Context

Baseline (host run): markdownlint-cli2 → 249 issues in 4 files; go vet clean; golangci-lint → 3× errcheck (`Close`), 3× staticcheck ST1005 in `cmd/main.go`. Interaction dumps (`01-*.md`) dominate MD013/MD036 noise.

## Goals / Non-Goals

**Goals:**

- Clear Go findings with behavior-preserving Close/error-string fixes.
- Clear Markdown findings with structural and wrap fixes; at most two fix/re-run rounds.

**Non-Goals:**

- Changing recognition pipeline, models, Docker native deps.
- New golangci or markdownlint config; blanket disable comments.

## Decisions

### Go Close handling

Use named return or explicit close-after-encode so `Close` errors are checked. On `load_persons`, return immediately if `Open` fails (avoids nil `Close`).

### Markdown

Prefer wrapping prose and fixing structure (MD022/031/029/040/060/032/047/012/049/051/041/036). For transcript dumps, convert emphasis-as-heading to real headings and wrap long lines; preserve command/output meaning.

### Deferrals

Record in tasks/design if a line cannot be wrapped without breaking a required single-line artifact (unlikely); otherwise fix.

## Risks / Trade-offs

- [Large MD013 volume in dumps] → Two rounds; stop and list leftovers if any remain.
- [load_persons Open error path] → Returning error changes control flow from “print and continue” to fail-fast — correct and safer; still not a recognition-algorithm change.

## Deferred after two rounds

| Finding id | Rationale |
|---|---|
| `markdownlint:MD013:01-TERMINAL_INTERACTION.md:*` (76) | Verbatim fenced terminal transcript; wrapping would falsify the lab dump. |
| `markdownlint:MD013:01-COPILOT_INTERACTION.md:57-58,172-175,271` (7) | Long table rows in chat export; stopped at round budget. |

## Open Questions

- None.
