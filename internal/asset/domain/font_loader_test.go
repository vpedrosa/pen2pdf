package domain_test

import (
	"testing"

	asset "github.com/vpedrosa/pen2pdf/internal/asset/domain"
)

type stubFontLoader struct {
	font *asset.FontData
	err  error
}

func (s *stubFontLoader) LoadFont(_, _, _ string) (*asset.FontData, error) {
	return s.font, s.err
}

func TestFontLoaderInterfaceCompliance(t *testing.T) {
	var _ asset.FontLoader = &stubFontLoader{}
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
