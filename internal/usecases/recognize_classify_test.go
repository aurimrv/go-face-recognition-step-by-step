package usecases

import (
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"
	"testing"

	"github.com/Kagami/go-face"
	"github.com/andre-ols/go-face-recognition/internal/entity"
)

func modelsDir(t *testing.T) string {
	t.Helper()
	// cwd = internal/usecases when running package tests
	p := filepath.Join("..", "..", "models")
	if _, err := os.Stat(p); err != nil {
		t.Skipf("models not available: %v", err)
	}
	return p
}

func personsFixture(t *testing.T, person, file string) string {
	t.Helper()
	p := filepath.Join("..", "..", "persons", person, file)
	if _, err := os.Stat(p); err != nil {
		t.Fatalf("fixture: %v", err)
	}
	return p
}

func imagesFixture(t *testing.T, parts ...string) string {
	t.Helper()
	p := filepath.Join(append([]string{"..", "..", "images"}, parts...)...)
	if _, err := os.Stat(p); err != nil {
		t.Fatalf("fixture: %v", err)
	}
	return p
}

func TestRecognizeAndClassify_HappyPath(t *testing.T) {
	// Covers: uncovered recognize_persons + classify_person Execute / constructors
	rec, err := face.NewRecognizer(modelsDir(t))
	if err != nil {
		t.Fatalf("NewRecognizer: %v", err)
	}
	defer rec.Close()

	persons := map[int]*entity.Person{
		0: entity.NewPerson(0, "Tom-Cruise", []string{
			personsFixture(t, "Tom-Cruise", "tom-01.jpg"),
			personsFixture(t, "Tom-Cruise", "tom-02.jpg"),
		}),
		1: entity.NewPerson(1, "Megan-Fox", []string{
			personsFixture(t, "Megan-Fox", "megan-01.jpg"),
		}),
	}

	recognize := NewRecognizePersonsUseCase(rec)
	if err := recognize.Execute(persons); err != nil {
		t.Fatalf("recognize: %v", err)
	}

	classify := NewClassifyPersonsUseCase(rec)
	unk := imagesFixture(t, "Megan-Fox", "unknown-01.jpg")
	got, err := classify.Execute(unk, 0.3)
	if err != nil {
		t.Fatalf("classify: %v", err)
	}
	if len(got) == 0 {
		t.Fatal("expected at least one recognized face")
	}
	if persons[got[0].ID].Name != "Megan-Fox" {
		t.Fatalf("got person %s, want Megan-Fox", persons[got[0].ID].Name)
	}
}

func TestRecognize_BadImage(t *testing.T) {
	// Covers: uncovered recognize error on bad file
	rec, err := face.NewRecognizer(modelsDir(t))
	if err != nil {
		t.Fatal(err)
	}
	defer rec.Close()

	bad := filepath.Join(t.TempDir(), "x.jpg")
	if err := os.WriteFile(bad, []byte("nope"), 0o644); err != nil {
		t.Fatal(err)
	}
	persons := map[int]*entity.Person{
		0: entity.NewPerson(0, "X", []string{bad}),
	}
	err = NewRecognizePersonsUseCase(rec).Execute(persons)
	if err == nil {
		t.Fatal("expected recognize error")
	}
}

func TestClassify_MissingFile(t *testing.T) {
	// Covers: uncovered classify RecognizeFile error
	rec, err := face.NewRecognizer(modelsDir(t))
	if err != nil {
		t.Fatal(err)
	}
	defer rec.Close()

	_, err = NewClassifyPersonsUseCase(rec).Execute(filepath.Join(t.TempDir(), "missing.jpg"), 0.3)
	if err == nil {
		t.Fatal("expected classify error")
	}
}

func TestClassify_NoFace(t *testing.T) {
	// Covers: uncovered classify len(unkFaces)==0 path when possible
	rec, err := face.NewRecognizer(modelsDir(t))
	if err != nil {
		t.Fatal(err)
	}
	defer rec.Close()

	// tiny valid JPEG with no face (1x1 pixel)
	noFace := filepath.Join(t.TempDir(), "noface.jpg")
	writeTinyJPEG(t, noFace)

	_, err = NewClassifyPersonsUseCase(rec).Execute(noFace, 0.3)
	if err == nil {
		t.Fatal("expected error when no faces found")
	}
}

func writeTinyJPEG(t *testing.T, path string) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			img.Set(x, y, color.RGBA{200, 200, 200, 255})
		}
	}
	if err := jpeg.Encode(f, img, nil); err != nil {
		t.Fatal(err)
	}
}
