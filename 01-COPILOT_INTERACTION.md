# Image analysis folder support
_Exported on 11/4/2026 at 19:35:31 PST from Cursor (3.28.7)_

---

**User**

/opsx-explore The current implementation analyzes a single image, called unknown.jpg, in the images folder and
produces the result.jpg image. I want to provide a specific image to analyze, or just a folder, and have it analyze
all images in that folder and its subfolders, naming each resulting image result_<original_image_name>.jpg. Can you
explore possible alternatives to do this? Run examples: go run ./cmd (images/unknown.jpg -> images/result.jpg); go
run ./cmd images/a.jpg images/b.jpg; go run ./cmd path/to/album/ (result_<stem>.jpg beside each source, including
subfolders).

---

**Cursor**

I'll explore how the current single-image flow works and what it would take to accept specific files or a folder tree. First I'm reading the explore-mode guidance and checking the project.

The no-argument run should stay a single special case. Everything else is one rule: each JPEG becomes `result_<stem>.jpg` in the same folder as the source. No OpenSpec change is open yet, so this is only the shape of the change.

Today `cmd/main.go` never reads arguments. It always classifies `images/unknown.jpg`, draws on a new `Drawer`, and writes `images/result.jpg`. The recognizer and the known-person samples are built once, before that single image. That one-time setup is the part worth keeping.

```
go run ./cmd                         images/unknown.jpg
                                              │
                                              ▼
 persons/ ──► samples ──► recognizer ──► classify ──► draw ──► images/result.jpg
 (once)        (once)       (once)

go run ./cmd images/a.jpg images/b.jpg
go run ./cmd path/to/album/          walk, subfolders included
                                              │
                                              ▼
                              result_<stem>.jpg beside each source
```

## What the three commands actually decide

| Invocation | Inputs | Output |
|---|---|---|
| `go run ./cmd` | only `images/unknown.jpg` | `images/result.jpg` |
| `go run ./cmd images/a.jpg images/b.jpg` | those files, in order | `images/result_a.jpg`, `images/result_b.jpg` |
| `go run ./cmd path/to/album/` | every JPEG under that tree | `path/to/album/.../result_<stem>.jpg` next to each file |

The default output name is the odd one. Passing `images/unknown.jpg` explicitly would write `images/result_unknown.jpg`, not `images/result.jpg`. Same picture, two names, depending on whether you passed it.

`images/` is already a tree. A walk of that folder is not the default run:

```
images/unknown.jpg                          → images/result_unknown.jpg
images/Denzel-Washington/unknown-01.jpg     → images/Denzel-Washington/result_unknown-01.jpg
images/Kate-Winslet/unknown-01.jpg          → images/Kate-Winslet/result_unknown-01.jpg
images/Tom-Cruise/unknown-01.avif           → skipped (decoder is JPEG-only)
images/Tom-Cruise/unknown-01.jpg            → images/Tom-Cruise/result_unknown-01.jpg
```

Putting every result in one flat folder would collide: two `unknown-01.jpg` files live in different person folders. Writing beside the source avoids that. The stem is the filename without its extension, so `unknown-01.jpg` becomes `result_unknown-01.jpg`, and `photo.JPEG` still becomes `result_photo.jpg`.

## Three ways to build it

**1. Loop in `main`.** Parse `os.Args`, turn them into paths, and call the existing classify and draw steps. Smallest diff, and it matches the script-shaped `main` you have now. The fragile part is the path rules: default versus files versus walk, JPEG filter, and "don't analyze the results of the last run." Those would sit untested inside `main`.

**2. A pure resolver, then the same loop.** Something like `Resolve(args) []Job{Src, Dst}` decides paths only. `main` still loads people, initializes the recognizer once, and loops. This is the cut I'd make. Discovery is easy to test without dlib, and classification stays a one-image use case.

```
args ──► Resolve ──► []Job{src, dst}
                         │
 persons + models ───────┤  (once)
                         ▼
                    for each job
                    classify → draw → save
```

**3. One new "analyze these paths" use case.** `main` becomes a thin CLI. The use case then owns the recognizer, the person map, the font path, walking, drawing, and saving. That hides the batch policy inside a type that is hard to test, and it mixes discovery with recognition. Heavier than this program needs. A CLI framework is the same kind of weight: three invocation forms don't need one.

All three share the same constraint. `face.NewRecognizer` and `SetSamples` stay outside the loop. One recognizer, used sequentially. The drawer reloads the font per image; for an album that cost is noise next to dlib.

## Edges the current code will hit

`Drawer.loadImage` only calls `jpeg.Decode`, and `NewDrawer` ignores that error. A bad file then panics on a nil image. With a hardcoded `unknown.jpg` that never comes up. With an album it will, on the first `.avif` or broken file. Whichever option you pick, decode has to return the error before drawing.

`Classify` treats "no face found" as a fatal error, but "faces found, none above 0.3" returns an empty slice and `main` still saves an unmarked copy. In a folder those are the same kind of photo. Fail-fast keeps today's behavior and throws away the rest of the album at the first landscape or `.avif`. Continue-and-report matches a folder better: log the path, move on, exit non-zero if any file failed.

A second walk will otherwise classify the previous outputs and write `result_result_<stem>.jpg`. Skip a file when its name is `result.jpg` or starts with `result_`. The precise version is: collect inputs, compute destinations, subtract those destinations from the inputs. A real photo named `result_party.jpg` would be skipped by the prefix rule and kept by the precise one.

`go run ./cmd dir/ a.jpg other/` falls out of the same resolver if every argument is either a file or a directory. The shell already expands `images/*/*.jpg` into the file form, so the walk exists for recursion the shell won't do.

## Open points

- A photo with no face: skip and continue, or abort the batch?
- `go run ./cmd images/unknown.jpg`: confirm `result_unknown.jpg`, leaving `result.jpg` only for the no-arg run.
- Mixed file and directory arguments: allow them, since the resolver is the same either way?
- An optional output directory would keep albums clean, and it brings back name collisions unless the relative path is kept.

I would lock option 2, beside-source naming, JPEG-only, and skip outputs. The error policy is the decision that changes the batch's character. Want me to turn that into a change proposal, or keep turning the error-policy question over?

---

**User**

/opsx-propose create the proposal from this exploration

---

**Cursor**

I'll turn the exploration into an OpenSpec change: create the change, then fill proposal, design, and tasks from what we locked in.

Creating the proposal from the exploration decisions.

Change **`batch-image-analysis`** is ready at `openspec/changes/batch-image-analysis/`.

| Artifact | What it captures |
|---|---|
| `proposal.md` | Why: hardcoded single image; what: CLI paths, naming, JPEG filter, skip results, batch errors |
| `design.md` | How: pure `Resolve` → `[]Job`, loop in `main`, beside-source `result_<stem>.jpg`, continue-on-error for explicit args |
| `specs/analysis-inputs/spec.md` | Requirements for default/files/dirs, JPEG-only, skip `result_*`, reuse recognizer, batch failures |
| `tasks.md` | Resolver + tests → drawer errors → CLI loop → README / manual check |

Locked from exploration: option 2 (resolver + loop), no-arg keeps `images/result.jpg`, explicit paths use stem naming, recursive dirs, JPEG-only, skip prior results. Open error-policy call was settled as continue-and-report for explicit-arg runs.

All artifacts created. Ready for implementation.

Run `/opsx:apply` or ask me to implement to start on the tasks.

---

**User**

/opsx-apply

---

**Cursor**

Using change: **batch-image-analysis** (override with `/opsx:apply <other>`). Reading apply instructions and context next.

## Implementing: batch-image-analysis (schema: spec-driven)

Progress: **0/14** tasks complete. Starting with path resolution (1.1–1.6).

While tests download deps, I'll update the drawer and wire the CLI loop.

Unit tests passed. Marking path-resolution and drawer/CLI tasks done, then updating the README and running a manual check.

## Implementation Complete

**Change:** batch-image-analysis  
**Schema:** spec-driven  
**Progress:** 14/14 tasks complete

### Completed This Session
- [x] 1.1–1.6 Path resolver (`ResolveAnalysisJobs`) + unit tests
- [x] 2.1–2.2 `NewDrawer` returns errors; no nil-image panic path
- [x] 3.1–3.4 CLI batch loop: reuse recognizer, continue on failure for args, fail-fast for no-arg
- [x] 4.1–4.2 README examples + manual runs verified

### Verified
- `go run ./cmd` → `images/result.jpg`
- `go run ./cmd images/Megan-Fox/` → `result_unknown-01.jpg` beside source; second run does not re-analyze prior results
- `go run ./cmd images/Megan-Fox/unknown-01.jpg images/Kate-Winslet/unknown-01.jpg` → both `result_<stem>.jpg` outputs

All tasks complete. You can archive this change with `/opsx:archive`.

---

**User**

/opsx-archive

---

**Cursor**

I'll list active changes and ask which one to archive.

Which change should I archive?

| Change | Status | Tasks |
|---|---|---|
| **batch-image-analysis** | complete | 14/14 |

Reply with the change name (e.g. `batch-image-analysis`) to continue.

---

**User**

/opsx-archive batch-image-analysis

---

**Cursor**

Archiving **batch-image-analysis**. Checking artifact/task completion and delta-spec sync state.

Artifacts and tasks are all complete (4/4 artifacts, 14/14 tasks).

**Delta spec sync assessment**

| Capability | Main spec | Delta action |
|---|---|---|
| `analysis-inputs` | *(missing)* | **ADD** full capability — 7 requirements (default target, explicit files, directories, JPEG-only, skip results, reuse recognizer, batch continue-on-failure) |

Would create `openspec/specs/analysis-inputs/spec.md` from the change delta.

How should I proceed?

1. **Sync now (recommended)** — copy requirements into main specs, then archive  
2. **Archive without syncing** — move the change to archive only  

Reply with `1` or `2`.

---

**User**

1

---

**Cursor**

Syncing the delta specs into main, then archiving.

Using `openspec archive` to sync specs and move the change.

## Archive Complete

**Change:** batch-image-analysis  
**Schema:** spec-driven  
**Archived to:** `openspec/changes/archive/2026-11-05-batch-image-analysis/`  
**Specs:** Synced to main specs (`analysis-inputs`: 7 requirements added)

All artifacts complete. All tasks complete.

