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

// Equivalence partitions for RecognizePersonsUseCase.Execute(persons) then
// ClassifyPersonsUseCase.Execute(unkImagePath, threshold):
//
//	EP-gallery-ok:     known-person JPEGs with detectable faces → SetSamples succeeds
//	EP-gallery-bad:    gallery path unreadable/corrupt → recognize error
//	EP-gallery-noface: gallery JPEG with no face → recognize error
//	EP-strong-match:   unknown of a trained identity at typical threshold → ≥1 hit, correct ID
//	EP-cross-nomatch:  unknown of identity absent from gallery at typical threshold → 0 hits, nil err
//	EP-noface-unk:     unknown JPEG with no face → classify error
//	EP-missing-unk:    unknown path missing → classify error
//
// Boundaries:
//
//	B-threshold-strict: threshold near 0 (e.g. 0.01) → even true identity yields 0 hits
//	B-threshold-typical: threshold 0.3 (CLI default) → strong match hits
//	B-threshold-loose:  threshold clearly above typical (e.g. 0.6) → strong match still hits

func modelsDir(t *testing.T) string {
	t.Helper()
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

func trainTomAndMegan(t *testing.T, rec *face.Recognizer) map[int]*entity.Person {
	t.Helper()
	persons := map[int]*entity.Person{
		0: entity.NewPerson(0, "Tom-Cruise", []string{
			personsFixture(t, "Tom-Cruise", "tom-01.jpg"),
			personsFixture(t, "Tom-Cruise", "tom-02.jpg"),
		}),
		1: entity.NewPerson(1, "Megan-Fox", []string{
			personsFixture(t, "Megan-Fox", "megan-01.jpg"),
			personsFixture(t, "Megan-Fox", "megan-02.jpg"),
		}),
	}
	if err := NewRecognizePersonsUseCase(rec).Execute(persons); err != nil {
		t.Fatalf("recognize gallery: %v", err)
	}
	return persons
}

func TestRecognize_GalleryOK(t *testing.T) {
	// Covers: EP-gallery-ok
	rec, err := face.NewRecognizer(modelsDir(t))
	if err != nil {
		t.Fatal(err)
	}
	defer rec.Close()
	_ = trainTomAndMegan(t, rec)
}

func TestRecognize_BadImage(t *testing.T) {
	// Covers: EP-gallery-bad
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
	if err := NewRecognizePersonsUseCase(rec).Execute(persons); err == nil {
		t.Fatal("expected recognize error")
	}
}

func TestRecognize_NoFaceInImage(t *testing.T) {
	// Covers: EP-gallery-noface
	rec, err := face.NewRecognizer(modelsDir(t))
	if err != nil {
		t.Fatal(err)
	}
	defer rec.Close()
	noFace := filepath.Join(t.TempDir(), "noface.jpg")
	writeTinyJPEG(t, noFace)
	persons := map[int]*entity.Person{
		0: entity.NewPerson(0, "X", []string{noFace}),
	}
	if err := NewRecognizePersonsUseCase(rec).Execute(persons); err == nil {
		t.Fatal("expected nil-face error")
	}
}

func TestClassify_StrongMatchTypicalThreshold(t *testing.T) {
	// Covers: EP-strong-match, B-threshold-typical
	rec, err := face.NewRecognizer(modelsDir(t))
	if err != nil {
		t.Fatal(err)
	}
	defer rec.Close()
	persons := trainTomAndMegan(t, rec)

	got, err := NewClassifyPersonsUseCase(rec).Execute(imagesFixture(t, "Megan-Fox", "unknown-01.jpg"), 0.3)
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

func TestClassify_StrongMatchLooseThreshold(t *testing.T) {
	// Covers: EP-strong-match, B-threshold-loose
	rec, err := face.NewRecognizer(modelsDir(t))
	if err != nil {
		t.Fatal(err)
	}
	defer rec.Close()
	persons := trainTomAndMegan(t, rec)

	got, err := NewClassifyPersonsUseCase(rec).Execute(imagesFixture(t, "Megan-Fox", "unknown-01.jpg"), 0.6)
	if err != nil {
		t.Fatalf("classify: %v", err)
	}
	foundMegan := false
	for _, rf := range got {
		if persons[rf.ID].Name == "Megan-Fox" {
			foundMegan = true
			break
		}
	}
	if !foundMegan {
		t.Fatalf("expected a Megan-Fox hit at loose threshold, got %+v", got)
	}
}

func TestClassify_CrossPersonNoMatch(t *testing.T) {
	// Covers: EP-cross-nomatch (Megan unknown vs Tom-only gallery @ 0.3)
	rec, err := face.NewRecognizer(modelsDir(t))
	if err != nil {
		t.Fatal(err)
	}
	defer rec.Close()
	persons := map[int]*entity.Person{
		0: entity.NewPerson(0, "Tom-Cruise", []string{
			personsFixture(t, "Tom-Cruise", "tom-01.jpg"),
			personsFixture(t, "Tom-Cruise", "tom-02.jpg"),
		}),
	}
	if err := NewRecognizePersonsUseCase(rec).Execute(persons); err != nil {
		t.Fatal(err)
	}
	got, err := NewClassifyPersonsUseCase(rec).Execute(imagesFixture(t, "Megan-Fox", "unknown-01.jpg"), 0.3)
	if err != nil {
		t.Fatalf("classify: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected no matches for cross-person @ 0.3, got %+v", got)
	}
}

func TestClassify_StrictThresholdNoHit(t *testing.T) {
	// Covers: B-threshold-strict (true identity still yields 0 hits)
	rec, err := face.NewRecognizer(modelsDir(t))
	if err != nil {
		t.Fatal(err)
	}
	defer rec.Close()
	_ = trainTomAndMegan(t, rec)
	got, err := NewClassifyPersonsUseCase(rec).Execute(imagesFixture(t, "Megan-Fox", "unknown-01.jpg"), 0.01)
	if err != nil {
		t.Fatalf("classify: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected no matches at strict threshold, got %+v", got)
	}
}

func TestClassify_MissingFile(t *testing.T) {
	// Covers: EP-missing-unk
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
	// Covers: EP-noface-unk
	rec, err := face.NewRecognizer(modelsDir(t))
	if err != nil {
		t.Fatal(err)
	}
	defer rec.Close()

	noFace := filepath.Join(t.TempDir(), "noface.jpg")
	writeTinyJPEG(t, noFace)

	_, err = NewClassifyPersonsUseCase(rec).Execute(noFace, 0.3)
	if err == nil {
		t.Fatal("expected error when no faces found")
	}
}
