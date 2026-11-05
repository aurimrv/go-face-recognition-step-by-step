## Why

The CLI always analyzes `images/unknown.jpg` and writes `images/result.jpg`, so album folders and extra test images under `images/` must be handled one file at a time by renaming. Accepting file and directory arguments lets the same trained recognizer annotate many images in one run.

## What Changes

- Accept zero or more path arguments: no args keeps today’s default; files and directories (including subfolders) become analysis jobs.
- Name each non-default output `result_<stem>.jpg` beside its source; the no-arg run still writes `images/result.jpg`.
- Restrict discovered inputs to JPEG files; skip previously generated result images so a second walk does not re-analyze them.
- Keep a single recognizer and known-person samples for the whole run; process jobs sequentially.
- On per-image failure in a batch (decode, no faces, etc.), log and continue; exit non-zero if any job failed. The no-arg path may keep fail-fast behavior for parity with today.
- Update usage docs to show the three invocation forms.

## Capabilities

### New Capabilities

- `analysis-inputs`: Resolve CLI path arguments into source/destination jobs (default image, explicit files, recursive directories) and define naming, JPEG filtering, result skipping, and batch error behavior for face analysis runs.

### Modified Capabilities

- (none — no existing specs)

## Impact

- `cmd/main.go`: read `os.Args`, resolve jobs, loop classify/draw/save instead of a single hardcoded path.
- New small package or helpers for path resolution (pure, testable without dlib).
- `internal/entity/drawer.go`: surface image-load errors so non-JPEG or corrupt files fail cleanly instead of panicking.
- Possibly `classify_person` / `main` error handling so batch runs can continue after a bad image.
- `README.md` (and any Docker/devcontainer run examples) for the new CLI forms.
- No new external dependencies; no change to person loading or recognition models.
