# Tests explained

> **Generated file** — source of truth is [`tests-explained.yaml`](tests-explained.yaml). Do not edit this Markdown by hand; upsert YAML then run `python3 scripts/render-tests-explained.py`.

## How to read this document

Each test case lists **why** it exists (objective), **which skill** generated or last updated it, and **what production behavior** it protects.

## Index by package

- [`internal/entity`](#package-internalentity)
- [`internal/usecases`](#package-internalusecases)

## Package `internal/entity`

### File `drawer_test.go`

Path: `internal/entity/drawer_test.go`

#### `TestDrawFace_TopEdgePixels`

- **Skill:** `mutation-test-generation`
- **Objective:** Assert top-edge rectangle pixels are drawer green so STATEMENT_REMOVE of DrawFace's first drawLine (drawer.go:43) is killed.
- **System under test:** internal/entity.DrawerImpl.DrawFace top edge
- **Production refs:** `internal/entity/drawer.go`
- **Kills mutant:** entity/drawer.go:43 STATEMENT_REMOVE top drawLine
- **Trace:** `id=github.com/andre-ols/go-face-recognition/internal/entity::TestDrawFace_TopEdgePixels`
- **Updated:** 2026-11-05

### File `person_test.go`

Path: `internal/entity/person_test.go`

#### `TestNewPerson`

- **Skill:** `white-box-test-generation`
- **Objective:** Cover NewPerson constructor fields (ID, Name, ImagesPath) previously uncovered under statement-coverage drive.
- **System under test:** internal/entity.NewPerson
- **Production refs:** `internal/entity/person.go`
- **Covers:** uncovered person.go NewPerson
- **Trace:** `id=github.com/andre-ols/go-face-recognition/internal/entity::TestNewPerson`
- **Updated:** 2026-11-05

## Package `internal/usecases`

### File `load_persons_test.go`

Path: `internal/usecases/load_persons_test.go`

#### `TestLoadPersons_CloseDirError`

- **Skill:** `white-box-test-generation`
- **Objective:** Force closePersonsDir seam failure after a successful load to cover the non-fatal directory close-error log branch.
- **System under test:** internal/usecases.LoadPersonsUseCaseImpl.Execute close path
- **Production refs:** `internal/usecases/load_persons.go`
- **Covers:** load_persons.go closePersonsDir error branch
- **Trace:** `id=github.com/andre-ols/go-face-recognition/internal/usecases::TestLoadPersons_CloseDirError`
- **Updated:** 2026-11-05

### File `recognize_classify_test.go`

Path: `internal/usecases/recognize_classify_test.go`

#### `TestClassify_CrossPersonNoMatch`

- **Skill:** `black-box-test-generation`
- **Objective:** Megan unknown vs Tom-only gallery at 0.3 returns zero hits (EP-cross-nomatch).
- **System under test:** ClassifyPersonsUseCase.Execute cross-identity non-match
- **Production refs:** `internal/usecases/classify_person.go`
- **Partitions:** `EP-cross-nomatch`
- **Trace:** `id=github.com/andre-ols/go-face-recognition/internal/usecases::TestClassify_CrossPersonNoMatch`
- **Updated:** 2026-11-05

#### `TestClassify_MissingFile`

- **Skill:** `black-box-test-generation`
- **Objective:** Missing unknown image path returns classify error (EP-missing-unk).
- **System under test:** ClassifyPersonsUseCase.Execute missing file
- **Production refs:** `internal/usecases/classify_person.go`
- **Partitions:** `EP-missing-unk`
- **Trace:** `id=github.com/andre-ols/go-face-recognition/internal/usecases::TestClassify_MissingFile`
- **Updated:** 2026-11-05

#### `TestClassify_NoFace`

- **Skill:** `black-box-test-generation`
- **Objective:** Valid JPEG with no face returns classify error (EP-noface-unk).
- **System under test:** ClassifyPersonsUseCase.Execute no-face unknown
- **Production refs:** `internal/usecases/classify_person.go`
- **Partitions:** `EP-noface-unk`
- **Trace:** `id=github.com/andre-ols/go-face-recognition/internal/usecases::TestClassify_NoFace`
- **Updated:** 2026-11-05

#### `TestClassify_StrictThresholdNoHit`

- **Skill:** `black-box-test-generation`
- **Objective:** Even a true Megan unknown yields zero hits at threshold 0.01 (B-threshold-strict).
- **System under test:** ClassifyPersonsUseCase.Execute strict threshold edge
- **Production refs:** `internal/usecases/classify_person.go`
- **Boundaries:** `B-threshold-strict`
- **Trace:** `id=github.com/andre-ols/go-face-recognition/internal/usecases::TestClassify_StrictThresholdNoHit`
- **Updated:** 2026-11-05

#### `TestClassify_StrongMatchLooseThreshold`

- **Skill:** `black-box-test-generation`
- **Objective:** Same Megan unknown still yields a Megan hit at threshold 0.6 (EP-strong-match, B-threshold-loose).
- **System under test:** ClassifyPersonsUseCase.Execute match strength @ 0.6
- **Production refs:** `internal/usecases/classify_person.go`
- **Partitions:** `EP-strong-match`
- **Boundaries:** `B-threshold-loose`
- **Trace:** `id=github.com/andre-ols/go-face-recognition/internal/usecases::TestClassify_StrongMatchLooseThreshold`
- **Updated:** 2026-11-05

#### `TestClassify_StrongMatchTypicalThreshold`

- **Skill:** `black-box-test-generation`
- **Objective:** Megan unknown against Tom+Megan gallery at threshold 0.3 returns Megan (EP-strong-match, B-threshold-typical).
- **System under test:** ClassifyPersonsUseCase.Execute match strength @ 0.3
- **Production refs:** `internal/usecases/classify_person.go`
- **Partitions:** `EP-strong-match`
- **Boundaries:** `B-threshold-typical`
- **Trace:** `id=github.com/andre-ols/go-face-recognition/internal/usecases::TestClassify_StrongMatchTypicalThreshold`
- **Updated:** 2026-11-05

#### `TestRecognize_BadImage`

- **Skill:** `black-box-test-generation`
- **Objective:** Corrupt gallery JPEG yields recognize error (EP-gallery-bad).
- **System under test:** RecognizePersonsUseCase.Execute invalid image path
- **Production refs:** `internal/usecases/recognize_persons.go`
- **Partitions:** `EP-gallery-bad`
- **Trace:** `id=github.com/andre-ols/go-face-recognition/internal/usecases::TestRecognize_BadImage`
- **Updated:** 2026-11-05

#### `TestRecognize_GalleryOK`

- **Skill:** `black-box-test-generation`
- **Objective:** Train Tom+Megan gallery from person JPEGs succeeds (EP-gallery-ok).
- **System under test:** RecognizePersonsUseCase.Execute with valid face samples
- **Production refs:** `internal/usecases/recognize_persons.go`
- **Partitions:** `EP-gallery-ok`
- **Trace:** `id=github.com/andre-ols/go-face-recognition/internal/usecases::TestRecognize_GalleryOK`
- **Updated:** 2026-11-05

#### `TestRecognize_NoFaceInImage`

- **Skill:** `black-box-test-generation`
- **Objective:** Gallery JPEG with no detectable face yields error (EP-gallery-noface).
- **System under test:** RecognizePersonsUseCase.Execute no-face JPEG
- **Production refs:** `internal/usecases/recognize_persons.go`
- **Partitions:** `EP-gallery-noface`
- **Trace:** `id=github.com/andre-ols/go-face-recognition/internal/usecases::TestRecognize_NoFaceInImage`
- **Updated:** 2026-11-05

### File `resolve_analysis_jobs_test.go`

Path: `internal/usecases/resolve_analysis_jobs_test.go`

#### `TestResolveAnalysisJobs_EmptyArgs`

- **Skill:** `black-box-test-generation`
- **Objective:** Verify empty CLI args resolve to the default job images/unknown.jpg → images/result.jpg (EP-empty).
- **System under test:** internal/usecases.ResolveAnalysisJobs — empty-args default job
- **Production refs:** `internal/usecases/resolve_analysis_jobs.go`
- **Partitions:** `EP-empty`
- **Trace:** `id=github.com/andre-ols/go-face-recognition/internal/usecases::TestResolveAnalysisJobs_EmptyArgs`
- **Updated:** 2026-11-05
