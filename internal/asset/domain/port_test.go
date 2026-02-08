package domain_test

import (
	"io"
	"testing"

	asset "github.com/vpedrosa/pen2pdf/internal/asset/domain"
)

type stubImageLoader struct {
	img *asset.ImageData
	err error
}

func (s *stubImageLoader) LoadImage(_ string) (*asset.ImageData, error) {
	return s.img, s.err
}

type stubFontLoader struct {
	font *asset.FontData
	err  error
}

func (s *stubFontLoader) LoadFont(_, _, _ string) (*asset.FontData, error) {
	return s.font, s.err
}

func TestImageLoaderInterfaceCompliance(t *testing.T) {
	var _ asset.ImageLoader = &stubImageLoader{}
}

func TestFontLoaderInterfaceCompliance(t *testing.T) {
	var _ asset.FontLoader = &stubFontLoader{}
}

func TestStubImageLoaderReturnsData(t *testing.T) {
	expected := &asset.ImageData{Path: "bg.jpg", Width: 800, Height: 600, Data: []byte{0xFF}}
	loader := &stubImageLoader{img: expected}

	got, err := loader.LoadImage("bg.jpg")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Path != "bg.jpg" {
		t.Errorf("expected path 'bg.jpg', got '%s'", got.Path)
	}
	if got.Width != 800 || got.Height != 600 {
		t.Errorf("expected 800x600, got %dx%d", got.Width, got.Height)
	}
}

func TestStubImageLoaderReturnsError(t *testing.T) {
	loader := &stubImageLoader{err: io.ErrUnexpectedEOF}

	_, err := loader.LoadImage("missing.jpg")
	if err != io.ErrUnexpectedEOF {
		t.Errorf("expected ErrUnexpectedEOF, got %v", err)
	}
}

func TestStubFontLoaderReturnsData(t *testing.T) {
	expected := &asset.FontData{Family: "Inter", Weight: "700", Style: "normal", Path: "Inter-Bold.ttf"}
	loader := &stubFontLoader{font: expected}

	got, err := loader.LoadFont("Inter", "700", "normal")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Family != "Inter" {
		t.Errorf("expected family 'Inter', got '%s'", got.Family)
	}
	if got.Weight != "700" {
		t.Errorf("expected weight '700', got '%s'", got.Weight)
	}
}
