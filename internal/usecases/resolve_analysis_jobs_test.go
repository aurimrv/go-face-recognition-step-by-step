package usecases

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveAnalysisJobs_Default(t *testing.T) {
	jobs, err := ResolveAnalysisJobs(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(jobs) != 1 {
		t.Fatalf("got %d jobs, want 1", len(jobs))
	}
	wantSrc := filepath.Join("images", "unknown.jpg")
	wantDst := filepath.Join("images", "result.jpg")
	if jobs[0].Src != wantSrc || jobs[0].Dst != wantDst {
		t.Fatalf("got %+v, want Src=%s Dst=%s", jobs[0], wantSrc, wantDst)
	}
}

func TestResolveAnalysisJobs_ExplicitFiles(t *testing.T) {
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

func TestResolveAnalysisJobs_ExplicitUnknownUsesStem(t *testing.T) {
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

func TestResolveAnalysisJobs_SkipResultFilesExplicit(t *testing.T) {
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

func TestResolveAnalysisJobs_CaseInsensitiveJPEG(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "Pic.JPG")
	mustTouch(t, src)

	jobs, err := ResolveAnalysisJobs([]string{dir})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(jobs) != 1 {
		t.Fatalf("got %d jobs, want 1", len(jobs))
	}
	assertJob(t, jobs[0], src, filepath.Join(dir, "result_Pic.jpg"))
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
