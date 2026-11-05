## Context

After the test-only guardrail, white-box generation reached **98.7%** and stopped. The three uncovered statements are error branches behind hard `os.File.Close` / `jpeg.Encode` calls that cannot be forced with real files under normal OS behavior. An earlier coverage push had used package-level seams; production was later restored without them while seam-based tests lingered (then removed). The user authorized a production fix via this OpenSpec change.

## Goals / Non-Goals

**Goals:**

- Make the three blocked branches reachable from same-package tests without changing recognition logic
- Keep default production behavior identical (seams call real Close/Encode)
- Restore tests and verify **100%** module statement coverage

**Non-Goals:**

- Changing white-box/black-box skills to allow unsolicited production edits
- Refactoring drawer/load-persons APIs for dependency injection frameworks
- Covering hypothetical bugs beyond these Close/Encode paths
- Reintroducing `exitFunc`/`run` in `cmd` (main is already covered via in-process success + subprocess Fatalf)

## Decisions

1. **Package-level function variables as seams**  
   Reintroduce unexported (or package-visible) vars such as `closeFile`, `jpegEncode` in `entity`, and `closePersonsDir` in `usecases`, assigned to real implementations at init. Tests swap them and restore via `t.Cleanup`.  
   *Alternatives:* interfaces wrapping every file op (heavier); build tags (more complex); leave uncovered (rejected—user wants fix).

2. **Preserve existing control flow**  
   Keep defer-style close-after-decode in `loadImage`, encode-then-close in `SaveImage`, and log-on-close-error in `Execute`—only replace the call sites with seam invocations so semantics stay the same.

3. **Tests stay in `*_test.go`**  
   White-box annotations (`// Covers: …`) for the three branches; no production asserts beyond error returns / successful Execute with forced close log path.

4. **Verify with full-module coverprofile**  
   Same recipe as the white-box skill: `go test ./... -coverprofile=coverage.out -coverpkg=./...` → `coverage.html` → `go tool cover -func` total 100%.

## Risks / Trade-offs

- [Risk] Seams are mutable globals → Mitigation: tests always restore; production never reassigns  
- [Risk] Future contributors treat seams as license to rewrite production under white-box skill → Mitigation: skills still require stop-and-report; this change is explicit OpenSpec authorization  
- [Risk] Over-broad seams → Mitigation: only Close/Encode for the three blocked sites

## Migration Plan

1. Apply seams + tests  
2. Confirm 100% coverage  
3. Archive change when done  
4. Rollback: revert seam vars and call sites (tests that force failures would fail to compile again)

## Open Questions

- None blocking apply.
