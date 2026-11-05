package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Kagami/go-face"
	"github.com/andre-ols/go-face-recognition/internal/entity"
	"github.com/andre-ols/go-face-recognition/internal/usecases"
)

func chdirRepoRoot(t *testing.T) {
	t.Helper()
	// package tests run with cwd = cmd/
	if err := os.Chdir(".."); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir("cmd")
	})
}

func TestProcessJob_Success(t *testing.T) {
	// Covers: uncovered cmd/main.go processJob happy path
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
	// Covers: uncovered processJob classify error return
	chdirRepoRoot(t)

	rec, err := face.NewRecognizer(modelsDir)
	if err != nil {
		t.Fatal(err)
	}
	defer rec.Close()

	persons := map[int]*entity.Person{0: entity.NewPerson(0, "X", nil)}
	job := usecases.Job{Src: filepath.Join(t.TempDir(), "missing.jpg"), Dst: filepath.Join(t.TempDir(), "out.jpg")}
	err = processJob(job, persons, usecases.NewClassifyPersonsUseCase(rec))
	if err == nil {
		t.Fatal("expected classify error")
	}
}
