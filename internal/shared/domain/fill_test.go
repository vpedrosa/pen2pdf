package domain_test

import (
	"testing"

	"github.com/vpedrosa/pen2pdf/internal/shared/domain"
)

func TestSolidFill(t *testing.T) {
	f := domain.SolidFill("#FF6B35")
	if f.Type != domain.FillSolid {
		t.Errorf("expected type solid, got %s", f.Type)
	}
	if f.Color != "#FF6B35" {
		t.Errorf("expected color '#FF6B35', got '%s'", f.Color)
	}
}

func TestSolidFillWithAlpha(t *testing.T) {
	f := domain.SolidFill("#000000BB")
	if f.Color != "#000000BB" {
		t.Errorf("expected color '#000000BB', got '%s'", f.Color)
	}
}

func TestSolidFillWithVariable(t *testing.T) {
	f := domain.SolidFill("$primary-color")
	if f.Color != "$primary-color" {
		t.Errorf("expected color '$primary-color', got '%s'", f.Color)
	}
}

func TestImageFill(t *testing.T) {
	f := domain.ImageFill("./images/bg.jpg", "fill", 0.3, true)
	if f.Type != domain.FillImage {
		t.Errorf("expected type image, got %s", f.Type)
	}
	if f.URL != "./images/bg.jpg" {
		t.Errorf("expected URL './images/bg.jpg', got '%s'", f.URL)
	}
	if f.Mode != "fill" {
		t.Errorf("expected mode 'fill', got '%s'", f.Mode)
	}
	if f.Opacity != 0.3 {
		t.Errorf("expected opacity 0.3, got %f", f.Opacity)
	}
	if !f.Enabled {
		t.Error("expected enabled to be true")
	}
}
