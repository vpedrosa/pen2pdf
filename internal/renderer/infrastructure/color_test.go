package infrastructure_test

import (
	"math"
	"testing"

	"github.com/vpedrosa/pen2pdf/internal/renderer/infrastructure"
)

func TestParseHexColorRGB(t *testing.T) {
	c, err := infrastructure.ParseHexColor("#FF6B35")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.R != 0xFF || c.G != 0x6B || c.B != 0x35 {
		t.Errorf("expected (255,107,53), got (%d,%d,%d)", c.R, c.G, c.B)
	}
	if c.A != 1.0 {
		t.Errorf("expected alpha 1.0, got %f", c.A)
	}
}

func TestParseHexColorRGBA(t *testing.T) {
	c, err := infrastructure.ParseHexColor("#000000BB")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.R != 0 || c.G != 0 || c.B != 0 {
		t.Errorf("expected (0,0,0), got (%d,%d,%d)", c.R, c.G, c.B)
	}
	expected := float64(0xBB) / 255.0
	if math.Abs(c.A-expected) > 0.01 {
		t.Errorf("expected alpha ~%f, got %f", expected, c.A)
	}
}

func TestParseHexColorWhite(t *testing.T) {
	c, err := infrastructure.ParseHexColor("#FFFFFF")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.R != 255 || c.G != 255 || c.B != 255 {
		t.Errorf("expected (255,255,255), got (%d,%d,%d)", c.R, c.G, c.B)
	}
}

func TestParseHexColorBlack(t *testing.T) {
	c, err := infrastructure.ParseHexColor("#000000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.R != 0 || c.G != 0 || c.B != 0 {
		t.Errorf("expected (0,0,0), got (%d,%d,%d)", c.R, c.G, c.B)
	}
}

func TestParseHexColorFullyTransparent(t *testing.T) {
	c, err := infrastructure.ParseHexColor("#FF000000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.A != 0.0 {
		t.Errorf("expected alpha 0.0, got %f", c.A)
	}
}

func TestParseHexColorFullyOpaque(t *testing.T) {
	c, err := infrastructure.ParseHexColor("#FF0000FF")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.A != 1.0 {
		t.Errorf("expected alpha 1.0, got %f", c.A)
	}
}

func TestParseHexColorInvalidNoHash(t *testing.T) {
	_, err := infrastructure.ParseHexColor("FF6B35")
	if err == nil {
		t.Fatal("expected error for missing #")
	}
}

func TestParseHexColorInvalidLength(t *testing.T) {
	_, err := infrastructure.ParseHexColor("#FFF")
	if err == nil {
		t.Fatal("expected error for wrong length")
	}
}

func TestParseHexColorInvalidChars(t *testing.T) {
	_, err := infrastructure.ParseHexColor("#GGGGGG")
	if err == nil {
		t.Fatal("expected error for invalid hex chars")
	}
}

func TestParseHexColorEmpty(t *testing.T) {
	_, err := infrastructure.ParseHexColor("")
	if err == nil {
		t.Fatal("expected error for empty string")
	}
}
