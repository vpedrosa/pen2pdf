package infrastructure_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	asset "github.com/vpedrosa/pen2pdf/internal/asset/domain"
	"github.com/vpedrosa/pen2pdf/internal/asset/infrastructure"
)

func TestFSFontLoaderImplementsPort(t *testing.T) {
	var _ asset.FontLoader = infrastructure.NewFSFontLoader()
}

func TestLoadFontRegular(t *testing.T) {
	dir := t.TempDir()
	createTestFont(t, filepath.Join(dir, "Inter-Regular.ttf"))

	loader := infrastructure.NewFSFontLoader(dir)
	font, err := loader.LoadFont("Inter", "400", "normal")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if font.Family != "Inter" {
		t.Errorf("expected family 'Inter', got '%s'", font.Family)
	}
	if font.Weight != "400" {
		t.Errorf("expected weight '400', got '%s'", font.Weight)
	}
	if len(font.Data) == 0 {
		t.Error("expected non-empty data")
	}
}

func TestLoadFontBold(t *testing.T) {
	dir := t.TempDir()
	createTestFont(t, filepath.Join(dir, "Inter-Bold.ttf"))

	loader := infrastructure.NewFSFontLoader(dir)
	font, err := loader.LoadFont("Inter", "700", "normal")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if font.Weight != "700" {
		t.Errorf("expected weight '700', got '%s'", font.Weight)
	}
}

func TestLoadFontBoldItalic(t *testing.T) {
	dir := t.TempDir()
	createTestFont(t, filepath.Join(dir, "Inter-BoldItalic.ttf"))

	loader := infrastructure.NewFSFontLoader(dir)
	font, err := loader.LoadFont("Inter", "700", "italic")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if font.Style != "italic" {
		t.Errorf("expected style 'italic', got '%s'", font.Style)
	}
}

func TestLoadFontItalic(t *testing.T) {
	dir := t.TempDir()
	createTestFont(t, filepath.Join(dir, "Inter-Italic.ttf"))

	loader := infrastructure.NewFSFontLoader(dir)
	font, err := loader.LoadFont("Inter", "400", "italic")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if font.Style != "italic" {
		t.Errorf("expected style 'italic', got '%s'", font.Style)
	}
}

func TestLoadFontWithSpacesInName(t *testing.T) {
	dir := t.TempDir()
	createTestFont(t, filepath.Join(dir, "OpenSans-Regular.ttf"))

	loader := infrastructure.NewFSFontLoader(dir)
	font, err := loader.LoadFont("Open Sans", "400", "normal")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if font.Family != "Open Sans" {
		t.Errorf("expected family 'Open Sans', got '%s'", font.Family)
	}
}

func TestLoadFontOTF(t *testing.T) {
	dir := t.TempDir()
	createTestFont(t, filepath.Join(dir, "Inter-Regular.otf"))

	loader := infrastructure.NewFSFontLoader(dir)
	font, err := loader.LoadFont("Inter", "400", "normal")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.HasSuffix(font.Path, ".otf") {
		t.Errorf("expected .otf path, got '%s'", font.Path)
	}
}

func TestLoadFontMultipleDirs(t *testing.T) {
	dir1 := t.TempDir()
	dir2 := t.TempDir()
	createTestFont(t, filepath.Join(dir2, "Inter-Bold.ttf"))

	loader := infrastructure.NewFSFontLoader(dir1, dir2)
	font, err := loader.LoadFont("Inter", "700", "normal")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.HasPrefix(font.Path, dir2) {
		t.Errorf("expected font from dir2, got path '%s'", font.Path)
	}
}

func TestLoadFontNotFound(t *testing.T) {
	loader := infrastructure.NewFSFontLoader(t.TempDir())
	_, err := loader.LoadFont("NonExistent", "400", "normal")
	if err == nil {
		t.Fatal("expected error for missing font")
	}
	if !strings.Contains(err.Error(), "font not found") {
		t.Errorf("expected 'font not found' in error, got: %s", err)
	}
}

func TestLoadFontWeightVariants(t *testing.T) {
	tests := []struct {
		weight   string
		filename string
	}{
		{"100", "Inter-Thin.ttf"},
		{"200", "Inter-ExtraLight.ttf"},
		{"300", "Inter-Light.ttf"},
		{"500", "Inter-Medium.ttf"},
		{"600", "Inter-SemiBold.ttf"},
		{"800", "Inter-ExtraBold.ttf"},
		{"900", "Inter-Black.ttf"},
	}

	for _, tt := range tests {
		t.Run("weight_"+tt.weight, func(t *testing.T) {
			dir := t.TempDir()
			createTestFont(t, filepath.Join(dir, tt.filename))

			loader := infrastructure.NewFSFontLoader(dir)
			font, err := loader.LoadFont("Inter", tt.weight, "normal")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if font.Weight != tt.weight {
				t.Errorf("expected weight '%s', got '%s'", tt.weight, font.Weight)
			}
		})
	}
}

func createTestFont(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	// Write dummy font data (not a real font, just for testing file loading)
	if err := os.WriteFile(path, []byte("fake-font-data"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
}
