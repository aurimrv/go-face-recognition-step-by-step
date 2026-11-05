package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/Kagami/go-face"
	"github.com/andre-ols/go-face-recognition/internal/entity"
	"github.com/andre-ols/go-face-recognition/internal/usecases"
)

func chdirRepoRoot(t *testing.T) {
	t.Helper()
	if err := os.Chdir(".."); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir("cmd")
	})
}

// reexecMain runs this test binary as a child with TEST_MAIN_CASE set so the
// child can call main() without killing the parent (log.Fatalf / os.Exit).
func reexecMain(t *testing.T, caseName string, extraEnv []string) *exec.Cmd {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run=^"+t.Name()+"$", "-test.v=false")
	cmd.Env = append(os.Environ(), append([]string{"TEST_MAIN_CASE=" + caseName}, extraEnv...)...)
	return cmd
}

func TestMain_ResolveError(t *testing.T) {
	// Covers: uncovered main ResolveAnalysisJobs error (log.Fatalf)
	if os.Getenv("TEST_MAIN_CASE") == "resolve" {
		chdirRepoRoot(t)
		os.Args = []string{"cmd", filepath.Join(t.TempDir(), "missing.jpg")}
		main()
		return
	}
	cmd := reexecMain(t, "resolve", nil)
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected non-zero exit")
	}
	if ee, ok := err.(*exec.ExitError); !ok || ee.ExitCode() == 0 {
		t.Fatalf("exit err %v", err)
	}
}

func TestMain_LoadPersonsError(t *testing.T) {
	// Covers: uncovered main load persons error
	if os.Getenv("TEST_MAIN_CASE") == "loadpersons" {
		chdirRepoRoot(t)
		personsDir = filepath.Join(t.TempDir(), "no-persons")
		os.Args = []string{"cmd", filepath.Join("images", "unknown.jpg")}
		main()
		return
	}
	cmd := reexecMain(t, "loadpersons", nil)
	if err := cmd.Run(); err == nil {
		t.Fatal("expected non-zero exit")
	}
}

func TestMain_BadModels(t *testing.T) {
	// Covers: uncovered main NewRecognizer error
	if os.Getenv("TEST_MAIN_CASE") == "models" {
		chdirRepoRoot(t)
		modelsDir = t.TempDir()
		os.Args = []string{"cmd", filepath.Join("images", "unknown.jpg")}
		main()
		return
	}
	cmd := reexecMain(t, "models", nil)
	if err := cmd.Run(); err == nil {
		t.Fatal("expected non-zero exit")
	}
}

func TestMain_RecognizeError(t *testing.T) {
	// Covers: uncovered main recognize persons error
	if os.Getenv("TEST_MAIN_CASE") == "recognize" {
		chdirRepoRoot(t)
		root := t.TempDir()
		pdir := filepath.Join(root, "NoFace")
		if err := os.MkdirAll(pdir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(pdir, "x.jpg"), []byte("not-a-jpeg"), 0o644); err != nil {
			t.Fatal(err)
		}
		personsDir = root
		os.Args = []string{"cmd", filepath.Join("images", "unknown.jpg")}
		main()
		return
	}
	cmd := reexecMain(t, "recognize", nil)
	if err := cmd.Run(); err == nil {
		t.Fatal("expected non-zero exit")
	}
}

func TestMain_FailFastClassifyError(t *testing.T) {
	// Covers: uncovered main failFast log.Fatalf when processJob fails
	if os.Getenv("TEST_MAIN_CASE") == "failfast" {
		chdirRepoRoot(t)
		fontPath = filepath.Join(t.TempDir(), "missing.ttf")
		os.Args = []string{"cmd"}
		main()
		return
	}
	cmd := reexecMain(t, "failfast", nil)
	if err := cmd.Run(); err == nil {
		t.Fatal("expected non-zero exit")
	}
}

func TestMain_BatchContinueOnError(t *testing.T) {
	// Covers: uncovered main non-failFast failed=true then os.Exit(1)
	if os.Getenv("TEST_MAIN_CASE") == "batch" {
		chdirRepoRoot(t)
		tmp := t.TempDir()
		goodSrc := filepath.Join("images", "Megan-Fox", "unknown-01.jpg")
		good := filepath.Join(tmp, "good.jpg")
		data, err := os.ReadFile(goodSrc)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(good, data, 0o644); err != nil {
			t.Fatal(err)
		}
		bad := filepath.Join(tmp, "bad.jpg")
		if err := os.WriteFile(bad, []byte("not-jpeg"), 0o644); err != nil {
			t.Fatal(err)
		}
		os.Args = []string{"cmd", bad, good}
		main()
		return
	}
	cmd := reexecMain(t, "batch", nil)
	if err := cmd.Run(); err == nil {
		t.Fatal("expected non-zero exit")
	}
}

func TestMain_SuccessWithExplicitImage(t *testing.T) {
	// Covers: uncovered main success path + defer Total-time print.
	// Safe in-process: success returns from main without os.Exit/log.Fatalf.
	chdirRepoRoot(t)
	outDir := t.TempDir()
	src := filepath.Join("images", "Megan-Fox", "unknown-01.jpg")
	tmpSrc := filepath.Join(outDir, "unknown-01.jpg")
	data, err := os.ReadFile(src)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(tmpSrc, data, 0o644); err != nil {
		t.Fatal(err)
	}
	oldArgs := os.Args
	os.Args = []string{"cmd", tmpSrc}
	t.Cleanup(func() { os.Args = oldArgs })
	main()
	if _, err := os.Stat(filepath.Join(outDir, "result_unknown-01.jpg")); err != nil {
		t.Fatalf("expected result image: %v", err)
	}
}

func TestProcessJob_Success(t *testing.T) {
	// Covers: processJob happy path
	chdirRepoRoot(t)

	rec, err := face.NewRecognizer(modelsDir)
	if err != nil {
		t.Fatalf("recognizer: %v", err)
	}
	defer rec.Close()

	load := usecases.NewLoadPersonsUseCase()
	persons, err := load.Execute(personsDir)
	if err != nil {
		t.Fatal(err)
	}
	if err := usecases.NewRecognizePersonsUseCase(rec).Execute(persons); err != nil {
		t.Fatal(err)
	}

	out := filepath.Join(t.TempDir(), "result_unknown.jpg")
	job := usecases.Job{
		Src: filepath.Join("images", "unknown.jpg"),
		Dst: out,
	}
	if err := processJob(job, persons, usecases.NewClassifyPersonsUseCase(rec)); err != nil {
		t.Fatalf("processJob: %v", err)
	}
	info, err := os.Stat(out)
	if err != nil || info.Size() == 0 {
		t.Fatalf("expected output image, err=%v", err)
	}
}

func TestProcessJob_ClassifyError(t *testing.T) {
	// Covers: processJob classify error
	chdirRepoRoot(t)

	rec, err := face.NewRecognizer(modelsDir)
	if err != nil {
		t.Fatal(err)
	}
	defer rec.Close()

	persons := map[int]*entity.Person{0: entity.NewPerson(0, "X", nil)}
	job := usecases.Job{Src: filepath.Join(t.TempDir(), "missing.jpg"), Dst: filepath.Join(t.TempDir(), "out.jpg")}
	if err := processJob(job, persons, usecases.NewClassifyPersonsUseCase(rec)); err == nil {
		t.Fatal("expected classify error")
	}
}

func TestProcessJob_DrawerError(t *testing.T) {
	// Covers: uncovered processJob NewDrawer error
	chdirRepoRoot(t)
	rec, err := face.NewRecognizer(modelsDir)
	if err != nil {
		t.Fatal(err)
	}
	defer rec.Close()
	persons, err := usecases.NewLoadPersonsUseCase().Execute(personsDir)
	if err != nil {
		t.Fatal(err)
	}
	if err := usecases.NewRecognizePersonsUseCase(rec).Execute(persons); err != nil {
		t.Fatal(err)
	}
	old := fontPath
	fontPath = filepath.Join(t.TempDir(), "no.ttf")
	t.Cleanup(func() { fontPath = old })
	job := usecases.Job{
		Src: filepath.Join("images", "unknown.jpg"),
		Dst: filepath.Join(t.TempDir(), "out.jpg"),
	}
	if err := processJob(job, persons, usecases.NewClassifyPersonsUseCase(rec)); err == nil {
		t.Fatal("expected drawer error")
	}
}

func TestProcessJob_SaveError(t *testing.T) {
	// Covers: uncovered processJob SaveImage error
	chdirRepoRoot(t)
	rec, err := face.NewRecognizer(modelsDir)
	if err != nil {
		t.Fatal(err)
	}
	defer rec.Close()
	persons, err := usecases.NewLoadPersonsUseCase().Execute(personsDir)
	if err != nil {
		t.Fatal(err)
	}
	if err := usecases.NewRecognizePersonsUseCase(rec).Execute(persons); err != nil {
		t.Fatal(err)
	}
	job := usecases.Job{
		Src: filepath.Join("images", "unknown.jpg"),
		Dst: filepath.Join(t.TempDir(), "no-dir", "out.jpg"),
	}
	if err := processJob(job, persons, usecases.NewClassifyPersonsUseCase(rec)); err == nil {
		t.Fatal("expected save error")
	}
}
