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
