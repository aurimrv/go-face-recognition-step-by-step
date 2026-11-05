package entity

import (
	"image"
	"image/jpeg"
	"io"
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

func TestDrawFace_TopEdgePixels(t *testing.T) {
	// Kills: entity/drawer.go:43 STATEMENT_REMOVE of top-edge drawLine
	// (LIVED in mutation-report.json — suite never asserted drawn pixels)
	img := repoFile(t, "persons", "Tom-Cruise", "tom-01.jpg")
	font := repoFile(t, "fonts", "Roboto-Regular.ttf")
	drawer, err := NewDrawer(img, font)
	if err != nil {
		t.Fatalf("NewDrawer: %v", err)
	}
	impl := drawer.(*DrawerImpl)
	rect := image.Rect(30, 40, 90, 100)
	want := impl.color
	drawer.DrawFace(rect, "X")

	// Sample the top edge (exclusive of corners shared with sides is fine).
	for _, x := range []int{rect.Min.X + 10, rect.Min.X + 30, rect.Max.X - 10} {
		got := impl.dst.RGBAAt(x, rect.Min.Y)
		if got != want {
			t.Fatalf("top-edge pixel (%d,%d)=%v, want drawer color %v", x, rect.Min.Y, got, want)
		}
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

func TestLoadImage_CloseError(t *testing.T) {
	// Covers: drawer.go loadImage closeFile error after successful decode
	img := repoFile(t, "persons", "Tom-Cruise", "tom-01.jpg")
	old := closeFile
	closeFile = func(f *os.File) error {
		_ = f.Close()
		return os.ErrClosed
	}
	t.Cleanup(func() { closeFile = old })
	d := &DrawerImpl{}
	if err := d.loadImage(img); err == nil {
		t.Fatal("expected close error")
	}
}

func TestSaveImage_EncodeError(t *testing.T) {
	// Covers: drawer.go SaveImage jpegEncode error branch
	img := repoFile(t, "persons", "Tom-Cruise", "tom-01.jpg")
	font := repoFile(t, "fonts", "Roboto-Regular.ttf")
	drawer, err := NewDrawer(img, font)
	if err != nil {
		t.Fatal(err)
	}
	old := jpegEncode
	jpegEncode = func(w io.Writer, m image.Image, o *jpeg.Options) error {
		return os.ErrInvalid
	}
	t.Cleanup(func() { jpegEncode = old })
	if err := drawer.SaveImage(filepath.Join(t.TempDir(), "out.jpg")); err == nil {
		t.Fatal("expected encode error")
	}
}

func TestSaveImage_CloseError(t *testing.T) {
	// Covers: drawer.go SaveImage closeFile error after encode
	img := repoFile(t, "persons", "Tom-Cruise", "tom-01.jpg")
	font := repoFile(t, "fonts", "Roboto-Regular.ttf")
	drawer, err := NewDrawer(img, font)
	if err != nil {
		t.Fatal(err)
	}
	old := closeFile
	closeFile = func(f *os.File) error {
		_ = f.Close()
		return os.ErrClosed
	}
	t.Cleanup(func() { closeFile = old })
	if err := drawer.SaveImage(filepath.Join(t.TempDir(), "out.jpg")); err == nil {
		t.Fatal("expected close error")
	}
}
