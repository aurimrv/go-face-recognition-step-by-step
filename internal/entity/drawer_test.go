package entity

import (
	"image"
	"os"
	"path/filepath"
	"testing"
)

func repoFile(t *testing.T, parts ...string) string {
	t.Helper()
	// tests run with cwd = package dir (internal/entity)
	p := filepath.Join(append([]string{"..", ".."}, parts...)...)
	if _, err := os.Stat(p); err != nil {
		t.Fatalf("fixture %s: %v", p, err)
	}
	return p
}

func TestNewDrawer_DrawFace_SaveImage(t *testing.T) {
	// Covers: uncovered drawer.go NewDrawer, loadImage, loadFont, DrawFace, drawLine, drawName, SaveImage
	img := repoFile(t, "persons", "Tom-Cruise", "tom-01.jpg")
	font := repoFile(t, "fonts", "Roboto-Regular.ttf")
	out := filepath.Join(t.TempDir(), "out.jpg")

	drawer, err := NewDrawer(img, font)
	if err != nil {
		t.Fatalf("NewDrawer: %v", err)
	}

	drawer.DrawFace(image.Rect(20, 20, 120, 140), "Tom")
	// reverse-direction edges also exercised by DrawFace (bottom/left)

	if err := drawer.SaveImage(out); err != nil {
		t.Fatalf("SaveImage: %v", err)
	}
	info, err := os.Stat(out)
	if err != nil || info.Size() == 0 {
		t.Fatalf("expected non-empty output, err=%v size=%v", err, info)
	}
}

func TestNewDrawer_BadImage(t *testing.T) {
	// Covers: uncovered NewDrawer loadImage error path
	font := repoFile(t, "fonts", "Roboto-Regular.ttf")
	_, err := NewDrawer(filepath.Join(t.TempDir(), "missing.jpg"), font)
	if err == nil {
		t.Fatal("expected error for missing image")
	}
}

func TestNewDrawer_BadFont(t *testing.T) {
	// Covers: uncovered NewDrawer loadFont error path
	img := repoFile(t, "persons", "Tom-Cruise", "tom-01.jpg")
	_, err := NewDrawer(img, filepath.Join(t.TempDir(), "missing.ttf"))
	if err == nil {
		t.Fatal("expected error for missing font")
	}
}

func TestNewDrawer_InvalidJPEG(t *testing.T) {
	// Covers: uncovered loadImage jpeg.Decode error
	font := repoFile(t, "fonts", "Roboto-Regular.ttf")
	bad := filepath.Join(t.TempDir(), "bad.jpg")
	if err := os.WriteFile(bad, []byte("not-a-jpeg"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := NewDrawer(bad, font)
	if err == nil {
		t.Fatal("expected error for invalid jpeg")
	}
}

func TestNewDrawer_InvalidFontBytes(t *testing.T) {
	// Covers: uncovered loadFont parse error
	img := repoFile(t, "persons", "Tom-Cruise", "tom-01.jpg")
	bad := filepath.Join(t.TempDir(), "bad.ttf")
	if err := os.WriteFile(bad, []byte("not-a-font"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := NewDrawer(img, bad)
	if err == nil {
		t.Fatal("expected error for invalid font")
	}
}

func TestSaveImage_NoDestination(t *testing.T) {
	// Covers: uncovered SaveImage nil dst
	d := &DrawerImpl{}
	if err := d.SaveImage(filepath.Join(t.TempDir(), "x.jpg")); err == nil {
		t.Fatal("expected error when dst is nil")
	}
}

func TestSaveImage_BadPath(t *testing.T) {
	// Covers: uncovered SaveImage create error
	img := repoFile(t, "persons", "Tom-Cruise", "tom-01.jpg")
	font := repoFile(t, "fonts", "Roboto-Regular.ttf")
	drawer, err := NewDrawer(img, font)
	if err != nil {
		t.Fatal(err)
	}
	badOut := filepath.Join(t.TempDir(), "no-such-dir", "out.jpg")
	if err := drawer.SaveImage(badOut); err == nil {
		t.Fatal("expected error for invalid output path")
	}
}

func TestDrawLine_Diagonals(t *testing.T) {
	// Covers: uncovered drawLine branches (sx/sy negative and steep lines)
	img := repoFile(t, "persons", "Tom-Cruise", "tom-01.jpg")
	font := repoFile(t, "fonts", "Roboto-Regular.ttf")
	d, err := NewDrawer(img, font)
	if err != nil {
		t.Fatal(err)
	}
	impl := d.(*DrawerImpl)
	impl.drawLine(impl.dst, 80, 80, 20, 20, impl.color) // up-left
	impl.drawLine(impl.dst, 20, 80, 80, 20, impl.color) // up-right
	impl.drawLine(impl.dst, 50, 10, 50, 90, impl.color) // vertical
}
