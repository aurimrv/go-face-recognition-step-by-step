package entity

import (
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"os"
	"path/filepath"
	"strings"
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

	for _, x := range []int{rect.Min.X + 10, rect.Min.X + 30, rect.Max.X - 10} {
		got := impl.dst.RGBAAt(x, rect.Min.Y)
		if got != want {
			t.Fatalf("top-edge pixel (%d,%d)=%v, want drawer color %v", x, rect.Min.Y, got, want)
		}
	}
}

func TestDrawFace_AllEdgesAndName(t *testing.T) {
	// Kills: L44–L46 STATEMENT_REMOVE of remaining edges; L48 drawName remove;
	// L48 height -→+; L168 color assignment / channel tweaks; L166 missing draw.Draw baseline.
	imgPath := repoFile(t, "persons", "Tom-Cruise", "tom-01.jpg")
	font := repoFile(t, "fonts", "Roboto-Regular.ttf")
	drawer, err := NewDrawer(imgPath, font)
	if err != nil {
		t.Fatalf("NewDrawer: %v", err)
	}
	impl := drawer.(*DrawerImpl)
	if impl.color != (color.RGBA{0, 255, 0, 255}) {
		t.Fatalf("drawer color %v, want pure green", impl.color)
	}
	// dst must be filled from the JPEG (kills L166 STATEMENT_REMOVE of draw.Draw).
	filled := impl.dst.RGBAAt(5, 5)
	if filled.A == 0 {
		t.Fatal("dst not filled from source image")
	}

	rect := image.Rect(40, 50, 100, 110)
	// Clear below the box so photo pixels cannot mask drawName / height mutants.
	for y := rect.Max.Y; y < rect.Max.Y+80; y++ {
		for x := 0; x < 250; x++ {
			impl.dst.Set(x, y, color.RGBA{0, 0, 0, 255})
		}
	}
	drawer.DrawFace(rect, "WWWW")
	want := impl.color
	assertEdgeGreen(t, impl.dst, rect.Min.X+15, rect.Min.Y, want)
	assertEdgeGreen(t, impl.dst, rect.Max.X, rect.Min.Y+15, want)
	assertEdgeGreen(t, impl.dst, rect.Min.X+15, rect.Max.Y, want)
	assertEdgeGreen(t, impl.dst, rect.Min.X, rect.Min.Y+15, want)
	exact := countExactColor(impl.dst, want, rect.Min.X, rect.Max.Y+1, rect.Min.X+150, rect.Max.Y+80)
	if exact < 50 {
		t.Fatalf("drawName painted too few exact green pixels: %d", exact)
	}
	// Height Max-Min=60 → fontSize 20, max glyph y≈129; Max+Min mutant → huge font / lower baseline.
	maxY := 0
	for y := rect.Max.Y + 1; y < rect.Max.Y+80; y++ {
		for x := rect.Min.X; x < rect.Min.X+150; x++ {
			if impl.dst.RGBAAt(x, y) == want && y > maxY {
				maxY = y
			}
		}
	}
	if maxY > 145 {
		t.Fatalf("text extends to y=%d; face height should be Max-Min only", maxY)
	}
}

func TestDrawLine_HorizontalComplete(t *testing.T) {
	// Kills: Bresenham sx/boundary/e2 mutants that skip or stop early on axis-aligned lines
	img := repoFile(t, "persons", "Tom-Cruise", "tom-01.jpg")
	font := repoFile(t, "fonts", "Roboto-Regular.ttf")
	d, err := NewDrawer(img, font)
	if err != nil {
		t.Fatal(err)
	}
	impl := d.(*DrawerImpl)
	y, x0, x1 := 25, 10, 70
	impl.drawLine(impl.dst, x0, y, x1, y, impl.color)
	for x := x0; x <= x1; x++ {
		if impl.dst.RGBAAt(x, y) != impl.color {
			t.Fatalf("horizontal gap at (%d,%d)", x, y)
		}
	}
	x, y0, y1 := 25, 10, 70
	impl.drawLine(impl.dst, x, y0, x, y1, impl.color)
	for y := y0; y <= y1; y++ {
		if impl.dst.RGBAAt(x, y) != impl.color {
			t.Fatalf("vertical gap at (%d,%d)", x, y)
		}
	}
	// Leftward / upward (sx/sy negative) with full coverage
	impl.drawLine(impl.dst, 60, 30, 20, 30, impl.color)
	for x := 20; x <= 60; x++ {
		if impl.dst.RGBAAt(x, 30) != impl.color {
			t.Fatalf("leftward gap at (%d,30)", x)
		}
	}
}

func TestDrawLine_DiagonalComplete(t *testing.T) {
	// Kills: L68 e2 factor; L69/L73 Bresenham boundary mutants (pixel count / corridor)
	img := repoFile(t, "persons", "Tom-Cruise", "tom-01.jpg")
	font := repoFile(t, "fonts", "Roboto-Regular.ttf")
	d, err := NewDrawer(img, font)
	if err != nil {
		t.Fatal(err)
	}
	impl := d.(*DrawerImpl)
	// Clear region so only this line counts.
	for y := 0; y < 80; y++ {
		for x := 0; x < 80; x++ {
			impl.dst.Set(x, y, color.RGBA{0, 0, 0, 255})
		}
	}
	impl.drawLine(impl.dst, 5, 5, 45, 20, impl.color)
	wantPts := []image.Point{
		{5, 5}, {6, 5}, {7, 6}, {8, 6}, {9, 7}, {10, 7}, {11, 7}, {12, 8}, {13, 8}, {14, 8},
		{15, 9}, {16, 9}, {17, 10}, {18, 10}, {19, 10}, {20, 11}, {21, 11}, {22, 11}, {23, 12}, {24, 12},
		{25, 13}, {26, 13}, {27, 13}, {28, 14}, {29, 14}, {30, 14}, {31, 15}, {32, 15}, {33, 16}, {34, 16},
		{35, 16}, {36, 17}, {37, 17}, {38, 17}, {39, 18}, {40, 18}, {41, 19}, {42, 19}, {43, 19}, {44, 20},
		{45, 20},
	}
	for _, p := range wantPts {
		if impl.dst.RGBAAt(p.X, p.Y) != impl.color {
			t.Fatalf("missing Bresenham pixel %v", p)
		}
	}
	n := 0
	for y := 0; y < 80; y++ {
		for x := 0; x < 80; x++ {
			if impl.dst.RGBAAt(x, y) == impl.color {
				n++
			}
		}
	}
	if n != len(wantPts) {
		t.Fatalf("Bresenham painted %d pixels, want %d", n, len(wantPts))
	}

	// Steep short segment: e2==dy distinguishes >= vs > (same length, different pixels).
	for y := 0; y < 20; y++ {
		for x := 0; x < 20; x++ {
			impl.dst.Set(x, y, color.RGBA{0, 0, 0, 255})
		}
	}
	impl.drawLine(impl.dst, 0, 0, 1, 2, impl.color)
	for _, p := range []image.Point{{0, 0}, {1, 1}, {1, 2}} {
		if impl.dst.RGBAAt(p.X, p.Y) != impl.color {
			t.Fatalf("steep segment missing %v (L69 e2>=dy path)", p)
		}
	}
	if impl.dst.RGBAAt(0, 1) == impl.color {
		t.Fatal("unexpected (0,1): e2>dy mutant path")
	}
}

func TestDrawName_FontSizeAndPlacement(t *testing.T) {
	// Kills: L83 3→2 (larger glyphs); L88 DPI ±1 when raster differs; L91/L93/L96 placement
	img := repoFile(t, "persons", "Tom-Cruise", "tom-01.jpg")
	font := repoFile(t, "fonts", "Roboto-Regular.ttf")
	d, err := NewDrawer(img, font)
	if err != nil {
		t.Fatal(err)
	}
	impl := d.(*DrawerImpl)
	for y := 100; y < 180; y++ {
		for x := 0; x < 200; x++ {
			impl.dst.Set(x, y, color.RGBA{0, 0, 0, 255})
		}
	}
	impl.drawName("Hi", 20, 120, 60) // fontSize = 20
	greenExact := countExactColor(impl.dst, impl.color, 15, 115, 120, 175)
	if greenExact != 47 {
		t.Fatalf("exact glyph pixels %d, want 47 (fontSize h/3 @ DPI 72)", greenExact)
	}
	baselineGreen := countExactColor(impl.dst, impl.color, 15, 135, 100, 155)
	if baselineGreen < 5 {
		t.Fatalf("expected glyphs near baseline band, got %d", baselineGreen)
	}
}

func TestLoadFont_OpenErrorIsPathError(t *testing.T) {
	// Kills: L101 BRANCH_IF that swallows ReadFile error then fails in Parse
	d := &DrawerImpl{}
	err := d.loadFont(filepath.Join(t.TempDir(), "missing.ttf"))
	if err == nil {
		t.Fatal("expected error")
	}
	var pe *os.PathError
	if !errors.As(err, &pe) {
		t.Fatalf("want PathError from ReadFile, got %T %v", err, err)
	}
}

func TestLoadImage_OpenErrorIsPathError(t *testing.T) {
	// Kills: L118 BRANCH_IF swallow Open error
	d := &DrawerImpl{}
	err := d.loadImage(filepath.Join(t.TempDir(), "missing.jpg"))
	if err == nil {
		t.Fatal("expected error")
	}
	var pe *os.PathError
	if !errors.As(err, &pe) {
		t.Fatalf("want PathError from Open, got %T %v", err, err)
	}
}

func TestLoadImage_DecodeErrorNotMaskedByClose(t *testing.T) {
	// Kills: L122 err==nil→true / &&→|| masking decode errors when close also fails
	bad := filepath.Join(t.TempDir(), "bad.jpg")
	if err := os.WriteFile(bad, []byte("not-a-jpeg"), 0o644); err != nil {
		t.Fatal(err)
	}
	old := closeFile
	closeFile = func(f *os.File) error {
		_ = f.Close()
		return os.ErrClosed
	}
	t.Cleanup(func() { closeFile = old })
	d := &DrawerImpl{}
	err := d.loadImage(bad)
	if err == nil {
		t.Fatal("expected jpeg decode error")
	}
	if errors.Is(err, os.ErrClosed) {
		t.Fatalf("decode error masked by close: %v", err)
	}
	if strings.Contains(strings.ToLower(err.Error()), "invalid argument") {
		t.Fatalf("unexpected close-priority error: %v", err)
	}
}

func TestSaveImage_CreateErrorIsPathError(t *testing.T) {
	// Kills: L143 BRANCH_IF swallow Create error
	img := repoFile(t, "persons", "Tom-Cruise", "tom-01.jpg")
	font := repoFile(t, "fonts", "Roboto-Regular.ttf")
	drawer, err := NewDrawer(img, font)
	if err != nil {
		t.Fatal(err)
	}
	err = drawer.SaveImage(filepath.Join(t.TempDir(), "no-dir", "out.jpg"))
	if err == nil {
		t.Fatal("expected create error")
	}
	var pe *os.PathError
	if !errors.As(err, &pe) {
		t.Fatalf("want PathError from Create, got %T %v", err, err)
	}
}

func assertEdgeGreen(t *testing.T, dst *image.RGBA, x, y int, want color.RGBA) {
	t.Helper()
	got := dst.RGBAAt(x, y)
	if got != want {
		t.Fatalf("pixel (%d,%d)=%v, want %v", x, y, got, want)
	}
}

func countExactColor(img *image.RGBA, want color.RGBA, x0, y0, x1, y1 int) int {
	n := 0
	b := img.Bounds()
	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			if !image.Pt(x, y).In(b) {
				continue
			}
			if img.RGBAAt(x, y) == want {
				n++
			}
		}
	}
	return n
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
