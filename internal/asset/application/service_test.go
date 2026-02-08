package application_test

import (
	"fmt"
	"testing"

	asset "github.com/vpedrosa/pen2pdf/internal/asset/domain"
	"github.com/vpedrosa/pen2pdf/internal/asset/application"
	shared "github.com/vpedrosa/pen2pdf/internal/shared/domain"
)

type stubFontLoader struct {
	available map[string]bool
}

func (l *stubFontLoader) LoadFont(family, weight, style string) (*asset.FontData, error) {
	key := family + "-" + weight + "-" + style
	if l.available[key] {
		return &asset.FontData{Family: family, Weight: weight, Style: style}, nil
	}
	return nil, fmt.Errorf("font not found: %s", key)
}

func TestDetectMissingFontsNone(t *testing.T) {
	loader := &stubFontLoader{available: map[string]bool{
		"Inter-400-normal": true,
	}}
	svc := application.NewFontService(loader)

	doc := &shared.Document{
		Children: []shared.Node{
			&shared.Frame{ID: "p", Children: []shared.Node{
				&shared.Text{ID: "t", FontFamily: "Inter", FontWeight: "400", FontStyle: "normal"},
			}},
		},
	}

	missing := svc.DetectMissingFonts(doc)
	if len(missing) != 0 {
		t.Errorf("expected 0 missing, got %d", len(missing))
	}
}

func TestDetectMissingFontsSome(t *testing.T) {
	loader := &stubFontLoader{available: map[string]bool{
		"Inter-400-normal": true,
	}}
	svc := application.NewFontService(loader)

	doc := &shared.Document{
		Children: []shared.Node{
			&shared.Frame{ID: "p", Children: []shared.Node{
				&shared.Text{ID: "t1", FontFamily: "Inter", FontWeight: "400", FontStyle: "normal"},
				&shared.Text{ID: "t2", FontFamily: "Missing", FontWeight: "700", FontStyle: "normal"},
			}},
		},
	}

	missing := svc.DetectMissingFonts(doc)
	if len(missing) != 1 {
		t.Fatalf("expected 1 missing, got %d", len(missing))
	}
	if missing[0].Family != "Missing" {
		t.Errorf("expected family 'Missing', got '%s'", missing[0].Family)
	}
}

func TestDetectMissingFontsEmptyDoc(t *testing.T) {
	svc := application.NewFontService(&stubFontLoader{})
	missing := svc.DetectMissingFonts(&shared.Document{})
	if len(missing) != 0 {
		t.Errorf("expected 0 missing for empty doc, got %d", len(missing))
	}
}
