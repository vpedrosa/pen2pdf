package infrastructure

import (
	"testing"
)

func TestFallbackStyleKeyRegular(t *testing.T) {
	if got := fallbackStyleKey("400", ""); got != "regular" {
		t.Errorf("expected 'regular', got '%s'", got)
	}
}

func TestFallbackStyleKeyBold(t *testing.T) {
	for _, weight := range []string{"700", "800", "900"} {
		if got := fallbackStyleKey(weight, ""); got != "bold" {
			t.Errorf("weight %s: expected 'bold', got '%s'", weight, got)
		}
	}
}

func TestFallbackStyleKeyItalic(t *testing.T) {
	if got := fallbackStyleKey("400", "italic"); got != "italic" {
		t.Errorf("expected 'italic', got '%s'", got)
	}
}

func TestFallbackStyleKeyBoldItalic(t *testing.T) {
	if got := fallbackStyleKey("700", "italic"); got != "bolditalic" {
		t.Errorf("expected 'bolditalic', got '%s'", got)
	}
}

func TestFallbackStyleKeyNonBoldWeights(t *testing.T) {
	for _, weight := range []string{"100", "200", "300", "500", "600"} {
		if got := fallbackStyleKey(weight, ""); got != "regular" {
			t.Errorf("weight %s: expected 'regular', got '%s'", weight, got)
		}
	}
}

func TestFallbackStyleKeyEmptyWeight(t *testing.T) {
	if got := fallbackStyleKey("", ""); got != "regular" {
		t.Errorf("expected 'regular', got '%s'", got)
	}
}
