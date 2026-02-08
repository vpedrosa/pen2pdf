package infrastructure_test

import (
	"fmt"
	"testing"

	asset "github.com/vpedrosa/pen2pdf/internal/asset/domain"
	layout "github.com/vpedrosa/pen2pdf/internal/layout/domain"
	"github.com/vpedrosa/pen2pdf/internal/layout/infrastructure"
)

type errorFontLoader struct{}

func (l *errorFontLoader) LoadFont(_, _, _ string) (*asset.FontData, error) {
	return nil, fmt.Errorf("font not found")
}

func TestGopdfTextMeasurerImplementsPort(t *testing.T) {
	var _ layout.TextMeasurer = infrastructure.NewGopdfTextMeasurer(nil)
}

func TestMeasureTextFallbackSingleLine(t *testing.T) {
	measurer := infrastructure.NewGopdfTextMeasurer(&errorFontLoader{})

	w, h := measurer.MeasureText("Hello", "Inter", 16, "400", 0)
	if w <= 0 {
		t.Errorf("expected positive width, got %f", w)
	}
	if h <= 0 {
		t.Errorf("expected positive height, got %f", h)
	}
}

func TestMeasureTextFallbackMultiLine(t *testing.T) {
	measurer := infrastructure.NewGopdfTextMeasurer(&errorFontLoader{})

	_, h1 := measurer.MeasureText("Line 1", "Inter", 16, "400", 0)
	_, h2 := measurer.MeasureText("Line 1\nLine 2", "Inter", 16, "400", 0)

	if h2 <= h1 {
		t.Errorf("expected multi-line height (%f) > single-line height (%f)", h2, h1)
	}
}

func TestMeasureTextFallbackWrapping(t *testing.T) {
	measurer := infrastructure.NewGopdfTextMeasurer(&errorFontLoader{})

	wNoLimit, hNoLimit := measurer.MeasureText("This is a really long text that should wrap", "Inter", 16, "400", 0)
	wLimited, hLimited := measurer.MeasureText("This is a really long text that should wrap", "Inter", 16, "400", 100)

	if wLimited > 100 {
		t.Errorf("expected width <= 100 with maxWidth, got %f", wLimited)
	}
	if hLimited <= hNoLimit {
		if wNoLimit > 100 {
			t.Errorf("expected wrapped height (%f) > unwrapped height (%f)", hLimited, hNoLimit)
		}
	}
}

func TestMeasureTextFallbackDifferentSizes(t *testing.T) {
	measurer := infrastructure.NewGopdfTextMeasurer(&errorFontLoader{})

	_, h12 := measurer.MeasureText("Hello", "Inter", 12, "400", 0)
	_, h24 := measurer.MeasureText("Hello", "Inter", 24, "400", 0)

	if h24 <= h12 {
		t.Errorf("expected fontSize 24 height (%f) > fontSize 12 height (%f)", h24, h12)
	}
}
