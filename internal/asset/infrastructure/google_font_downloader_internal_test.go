package infrastructure

import (
	"testing"
)

func TestBuildFilenameRegular(t *testing.T) {
	got := buildFilename("Inter", "400", "")
	if got != "Inter-Regular.ttf" {
		t.Errorf("expected 'Inter-Regular.ttf', got '%s'", got)
	}
}

func TestBuildFilenameBold(t *testing.T) {
	got := buildFilename("Inter", "700", "")
	if got != "Inter-Bold.ttf" {
		t.Errorf("expected 'Inter-Bold.ttf', got '%s'", got)
	}
}

func TestBuildFilenameBoldItalic(t *testing.T) {
	got := buildFilename("Inter", "700", "italic")
	if got != "Inter-BoldItalic.ttf" {
		t.Errorf("expected 'Inter-BoldItalic.ttf', got '%s'", got)
	}
}

func TestBuildFilenameItalicOnly(t *testing.T) {
	got := buildFilename("Inter", "400", "italic")
	if got != "Inter-Italic.ttf" {
		t.Errorf("expected 'Inter-Italic.ttf', got '%s'", got)
	}
}

func TestBuildFilenameWithSpaces(t *testing.T) {
	got := buildFilename("Open Sans", "600", "")
	if got != "OpenSans-SemiBold.ttf" {
		t.Errorf("expected 'OpenSans-SemiBold.ttf', got '%s'", got)
	}
}

func TestBuildFilenameAllWeights(t *testing.T) {
	tests := []struct {
		weight   string
		expected string
	}{
		{"100", "Inter-Thin.ttf"},
		{"200", "Inter-ExtraLight.ttf"},
		{"300", "Inter-Light.ttf"},
		{"400", "Inter-Regular.ttf"},
		{"500", "Inter-Medium.ttf"},
		{"600", "Inter-SemiBold.ttf"},
		{"700", "Inter-Bold.ttf"},
		{"800", "Inter-ExtraBold.ttf"},
		{"900", "Inter-Black.ttf"},
	}

	for _, tt := range tests {
		t.Run("weight_"+tt.weight, func(t *testing.T) {
			got := buildFilename("Inter", tt.weight, "")
			if got != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, got)
			}
		})
	}
}

func TestBuildCSSURLNormal(t *testing.T) {
	got := buildCSSURL("Inter", "700", "")
	expected := "https://fonts.googleapis.com/css2?family=Inter:wght@700"
	if got != expected {
		t.Errorf("expected '%s', got '%s'", expected, got)
	}
}

func TestBuildCSSURLItalic(t *testing.T) {
	got := buildCSSURL("Inter", "400", "italic")
	expected := "https://fonts.googleapis.com/css2?family=Inter:ital,wght@1,400"
	if got != expected {
		t.Errorf("expected '%s', got '%s'", expected, got)
	}
}

func TestBuildCSSURLWithSpaces(t *testing.T) {
	got := buildCSSURL("Open Sans", "400", "")
	expected := "https://fonts.googleapis.com/css2?family=Open+Sans:wght@400"
	if got != expected {
		t.Errorf("expected '%s', got '%s'", expected, got)
	}
}
