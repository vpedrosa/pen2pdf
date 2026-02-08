package infrastructure_test

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	asset "github.com/vpedrosa/pen2pdf/internal/asset/domain"
	"github.com/vpedrosa/pen2pdf/internal/asset/infrastructure"
)

func TestFSImageLoaderImplementsPort(t *testing.T) {
	var _ asset.ImageLoader = infrastructure.NewFSImageLoader("")
}

func TestLoadImageRelativePath(t *testing.T) {
	dir := t.TempDir()
	createTestPNG(t, filepath.Join(dir, "images", "bg.png"), 100, 50)

	loader := infrastructure.NewFSImageLoader(dir)
	img, err := loader.LoadImage("images/bg.png")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if img.Width != 100 {
		t.Errorf("expected width 100, got %d", img.Width)
	}
	if img.Height != 50 {
		t.Errorf("expected height 50, got %d", img.Height)
	}
	if len(img.Data) == 0 {
		t.Error("expected non-empty data")
	}
	if img.Path != filepath.Join(dir, "images", "bg.png") {
		t.Errorf("expected absolute path, got '%s'", img.Path)
	}
}

func TestLoadImageAbsolutePath(t *testing.T) {
	dir := t.TempDir()
	absPath := filepath.Join(dir, "test.png")
	createTestPNG(t, absPath, 200, 300)

	loader := infrastructure.NewFSImageLoader("/other/dir")
	img, err := loader.LoadImage(absPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if img.Width != 200 || img.Height != 300 {
		t.Errorf("expected 200x300, got %dx%d", img.Width, img.Height)
	}
}

func TestLoadImageMissing(t *testing.T) {
	loader := infrastructure.NewFSImageLoader(t.TempDir())
	_, err := loader.LoadImage("nonexistent.png")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
	if !strings.Contains(err.Error(), "load image") {
		t.Errorf("expected 'load image' in error, got: %s", err)
	}
}

func createTestPNG(t *testing.T, path string, w, h int) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	defer f.Close()

	img := image.NewRGBA(image.Rect(0, 0, w, h))
	img.Set(0, 0, color.RGBA{R: 255, A: 255})
	if err := png.Encode(f, img); err != nil {
		t.Fatalf("encode png: %v", err)
	}
}
