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

func TestImageLoaderInterfaceCompliance(t *testing.T) {
	var _ asset.ImageLoader = &stubImageLoader{}
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
