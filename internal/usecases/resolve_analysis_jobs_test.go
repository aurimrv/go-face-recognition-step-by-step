package usecases

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Equivalence partitions for ResolveAnalysisJobs(args):
//   EP-empty:   len(args)==0 (nil or empty) → one default job
//               images/unknown.jpg → images/result.jpg
//   EP-files:   args are existing regular files (non-result) → one job each with
//               beside-source result_<stem>.jpg
//   EP-dirs:    args are directories → recursive JPEG discovery, same naming
//   EP-mixed:   file and directory args combined in order
//   EP-missing: arg path does not exist → error
//   EP-nosel:   args resolve to zero jobs (only skipped results / empty dir) → error
// Boundaries:
//   B-stem:     explicit …/unknown.jpg → result_unknown.jpg (not images/result.jpg)
//   B-ext-dir:  under dirs, .jpg/.jpeg/.JPG selected; non-JPEG (e.g. .avif) omitted
//   B-result:   basename result.jpg or result_*.jpg skipped as inputs
//   B-one-file: single-file edge of EP-files
//   B-empty-dir: directory with no selectable JPEGs contributes no jobs

func TestResolveAnalysisJobs_EmptyArgs(t *testing.T) {
	// Covers: EP-empty
	jobs, err := ResolveAnalysisJobs(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(jobs) != 1 {
		t.Fatalf("got %d jobs, want 1", len(jobs))
	}
	wantSrc := filepath.Join("images", "unknown.jpg")
	wantDst := filepath.Join("images", "result.jpg")
	assertJob(t, jobs[0], wantSrc, wantDst)
}

func TestResolveAnalysisJobs_EmptySlice(t *testing.T) {
	// Covers: EP-empty (empty slice same as nil: len==0)
	jobs, err := ResolveAnalysisJobs([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(jobs) != 1 {
		t.Fatalf("got %d jobs, want 1", len(jobs))
	}
	assertJob(t, jobs[0], filepath.Join("images", "unknown.jpg"), filepath.Join("images", "result.jpg"))
}

func TestResolveAnalysisJobs_ExplicitFiles(t *testing.T) {
	// Covers: EP-files
	dir := t.TempDir()
	a := filepath.Join(dir, "a.jpg")
	b := filepath.Join(dir, "b.jpg")
	mustTouch(t, a)
	mustTouch(t, b)

	jobs, err := ResolveAnalysisJobs([]string{a, b})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(jobs) != 2 {
		t.Fatalf("got %d jobs, want 2", len(jobs))
	}
	assertJob(t, jobs[0], a, filepath.Join(dir, "result_a.jpg"))
	assertJob(t, jobs[1], b, filepath.Join(dir, "result_b.jpg"))
}

func TestResolveAnalysisJobs_OneFile(t *testing.T) {
	// Covers: EP-files, B-one-file
	dir := t.TempDir()
	src := filepath.Join(dir, "solo.jpg")
	mustTouch(t, src)

	jobs, err := ResolveAnalysisJobs([]string{src})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(jobs) != 1 {
		t.Fatalf("got %d jobs, want 1", len(jobs))
	}
	assertJob(t, jobs[0], src, filepath.Join(dir, "result_solo.jpg"))
}

func TestResolveAnalysisJobs_ExplicitUnknownUsesStem(t *testing.T) {
	// Covers: EP-files, B-stem
	dir := t.TempDir()
	src := filepath.Join(dir, "unknown.jpg")
	mustTouch(t, src)

	jobs, err := ResolveAnalysisJobs([]string{src})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertJob(t, jobs[0], src, filepath.Join(dir, "result_unknown.jpg"))
}

func TestResolveAnalysisJobs_NestedDirectory(t *testing.T) {
	// Covers: EP-dirs, B-ext-dir, B-result
	root := t.TempDir()
	nested := filepath.Join(root, "album", "person")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	photo := filepath.Join(nested, "photo.jpg")
	avif := filepath.Join(nested, "photo.avif")
	prior := filepath.Join(nested, "result_photo.jpg")
	mustTouch(t, photo)
	mustTouch(t, avif)
	mustTouch(t, prior)

	jobs, err := ResolveAnalysisJobs([]string{root})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(jobs) != 1 {
		t.Fatalf("got %d jobs, want 1 (skip avif and prior result): %+v", len(jobs), jobs)
	}
	assertJob(t, jobs[0], photo, filepath.Join(nested, "result_photo.jpg"))
}

func TestResolveAnalysisJobs_DirJPEGCaseAndExt(t *testing.T) {
	// Covers: EP-dirs, B-ext-dir (.JPG and .jpeg)
	dir := t.TempDir()
	jpgUpper := filepath.Join(dir, "Pic.JPG")
	jpegLower := filepath.Join(dir, "other.jpeg")
	mustTouch(t, jpgUpper)
	mustTouch(t, jpegLower)

	jobs, err := ResolveAnalysisJobs([]string{dir})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(jobs) != 2 {
		t.Fatalf("got %d jobs, want 2: %+v", len(jobs), jobs)
	}
	bySrc := map[string]Job{}
	for _, j := range jobs {
		bySrc[j.Src] = j
	}
	assertJob(t, bySrc[jpgUpper], jpgUpper, filepath.Join(dir, "result_Pic.jpg"))
	assertJob(t, bySrc[jpegLower], jpegLower, filepath.Join(dir, "result_other.jpg"))
}

func TestResolveAnalysisJobs_SkipResultFilesExplicit(t *testing.T) {
	// Covers: EP-files, B-result
	dir := t.TempDir()
	photo := filepath.Join(dir, "photo.jpg")
	result := filepath.Join(dir, "result_photo.jpg")
	legacy := filepath.Join(dir, "result.jpg")
	mustTouch(t, photo)
	mustTouch(t, result)
	mustTouch(t, legacy)

	jobs, err := ResolveAnalysisJobs([]string{photo, result, legacy})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(jobs) != 1 {
		t.Fatalf("got %d jobs, want 1: %+v", len(jobs), jobs)
	}
	assertJob(t, jobs[0], photo, filepath.Join(dir, "result_photo.jpg"))
}

func TestResolveAnalysisJobs_MixedFileAndDirectory(t *testing.T) {
	// Covers: EP-mixed, B-ext-dir
	root := t.TempDir()
	file := filepath.Join(root, "solo.jpg")
	album := filepath.Join(root, "album")
	if err := os.MkdirAll(album, 0o755); err != nil {
		t.Fatal(err)
	}
	nested := filepath.Join(album, "shot.JPEG")
	mustTouch(t, file)
	mustTouch(t, nested)

	jobs, err := ResolveAnalysisJobs([]string{file, album})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(jobs) != 2 {
		t.Fatalf("got %d jobs, want 2: %+v", len(jobs), jobs)
	}
	assertJob(t, jobs[0], file, filepath.Join(root, "result_solo.jpg"))
	assertJob(t, jobs[1], nested, filepath.Join(album, "result_shot.jpg"))
}

func TestResolveAnalysisJobs_MissingPath(t *testing.T) {
	// Covers: EP-missing
	_, err := ResolveAnalysisJobs([]string{filepath.Join(t.TempDir(), "no-such.jpg")})
	if err == nil {
		t.Fatal("expected error for missing path")
	}
	if !strings.Contains(err.Error(), "stat") {
		t.Fatalf("error %q should mention stat", err)
	}
}

func TestResolveAnalysisJobs_OnlySkippedResults(t *testing.T) {
	// Covers: EP-nosel, B-result
	dir := t.TempDir()
	mustTouch(t, filepath.Join(dir, "result.jpg"))
	mustTouch(t, filepath.Join(dir, "result_old.jpg"))

	_, err := ResolveAnalysisJobs([]string{
		filepath.Join(dir, "result.jpg"),
		filepath.Join(dir, "result_old.jpg"),
	})
	if err == nil {
		t.Fatal("expected error when all inputs are skipped results")
	}
	if !strings.Contains(err.Error(), "no analysis jobs") {
		t.Fatalf("error %q should mention no analysis jobs", err)
	}
}

func TestResolveAnalysisJobs_EmptyDirectory(t *testing.T) {
	// Covers: EP-nosel, B-empty-dir
	dir := t.TempDir()
	_, err := ResolveAnalysisJobs([]string{dir})
	if err == nil {
		t.Fatal("expected error for empty directory")
	}
	if !strings.Contains(err.Error(), "no analysis jobs") {
		t.Fatalf("error %q should mention no analysis jobs", err)
	}
}

func mustTouch(t *testing.T, path string) {
	t.Helper()
	if err := os.WriteFile(path, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
}

func assertJob(t *testing.T, job Job, src, dst string) {
	t.Helper()
	if job.Src != src || job.Dst != dst {
		t.Fatalf("got Src=%q Dst=%q, want Src=%q Dst=%q", job.Src, job.Dst, src, dst)
	}
}
