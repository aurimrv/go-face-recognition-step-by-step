package entity

import "testing"

func TestNewPerson(t *testing.T) {
	// Covers: uncovered person.go NewPerson
	paths := []string{"a.jpg", "b.jpg"}
	p := NewPerson(3, "Ada", paths)
	if p.ID != 3 || p.Name != "Ada" || len(p.ImagesPath) != 2 {
		t.Fatalf("unexpected person: %+v", p)
	}
	if p.ImagesPath[0] != "a.jpg" || p.ImagesPath[1] != "b.jpg" {
		t.Fatalf("unexpected paths: %v", p.ImagesPath)
	}
}
