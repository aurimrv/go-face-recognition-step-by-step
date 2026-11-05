package usecases

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewLoadPersonsUseCase_Execute(t *testing.T) {
	// Covers: uncovered load_persons.go NewLoadPersonsUseCase + Execute happy path
	root := t.TempDir()
	alice := filepath.Join(root, "Alice")
	bob := filepath.Join(root, "Bob")
	if err := os.MkdirAll(alice, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(bob, 0o755); err != nil {
		t.Fatal(err)
	}
	mustTouch(t, filepath.Join(alice, "a1.jpg"))
	mustTouch(t, filepath.Join(bob, "b1.jpg"))
	// skip non-dir file at root
	mustTouch(t, filepath.Join(root, "readme.txt"))

	uc := NewLoadPersonsUseCase()
	persons, err := uc.Execute(root)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if len(persons) != 2 {
		t.Fatalf("got %d persons, want 2: %+v", len(persons), persons)
	}
	names := map[string]bool{}
	for _, p := range persons {
		names[p.Name] = true
		if len(p.ImagesPath) != 1 {
			t.Fatalf("person %s images: %v", p.Name, p.ImagesPath)
		}
	}
	if !names["Alice"] || !names["Bob"] {
		t.Fatalf("unexpected names: %v", names)
	}
}

func TestLoadPersons_MissingDir(t *testing.T) {
	// Covers: uncovered load_persons open error
	uc := NewLoadPersonsUseCase()
	_, err := uc.Execute(filepath.Join(t.TempDir(), "nope"))
	if err == nil {
		t.Fatal("expected error for missing persons dir")
	}
}

func TestLoadPersons_RootIsFile(t *testing.T) {
	// Covers: uncovered Readdir error when personsDir is a file
	root := filepath.Join(t.TempDir(), "not-a-dir")
	mustTouch(t, root)
	_, err := NewLoadPersonsUseCase().Execute(root)
	if err == nil {
		t.Fatal("expected Readdir error")
	}
}

func TestLoadPersons_UnreadablePersonDir(t *testing.T) {
	// Covers: uncovered ReadDir error on person subdirectory
	root := t.TempDir()
	p := filepath.Join(root, "Locked")
	if err := os.MkdirAll(p, 0o755); err != nil {
		t.Fatal(err)
	}
	mustTouch(t, filepath.Join(p, "a.jpg"))
	if err := os.Chmod(p, 0o000); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(p, 0o755) })
	_, err := NewLoadPersonsUseCase().Execute(root)
	if err == nil {
		t.Fatal("expected ReadDir error")
	}
}

func TestLoadPersons_CloseDirError(t *testing.T) {
	// Covers: load_persons.go closePersonsDir error branch (logged, non-fatal)
	root := t.TempDir()
	p := filepath.Join(root, "Alice")
	if err := os.MkdirAll(p, 0o755); err != nil {
		t.Fatal(err)
	}
	mustTouch(t, filepath.Join(p, "a.jpg"))
	old := closePersonsDir
	closePersonsDir = func(f *os.File) error {
		_ = f.Close()
		return os.ErrClosed
	}
	t.Cleanup(func() { closePersonsDir = old })
	persons, err := NewLoadPersonsUseCase().Execute(root)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if len(persons) != 1 {
		t.Fatalf("got %d persons, want 1", len(persons))
	}
}
