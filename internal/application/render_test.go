package application

import (
	"bytes"
	"strings"
	"testing"

	shared "github.com/vpedrosa/pen2pdf/internal/shared/domain"
)

func TestFilterPagesMatchSingle(t *testing.T) {
	children := []shared.Node{
		&shared.Frame{ID: "p1", Name: "Front"},
		&shared.Frame{ID: "p2", Name: "Back"},
	}

	filtered, err := filterPages(children, "Front")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filtered) != 1 {
		t.Fatalf("expected 1 page, got %d", len(filtered))
	}
	if filtered[0].GetName() != "Front" {
		t.Errorf("expected 'Front', got '%s'", filtered[0].GetName())
	}
}

func TestFilterPagesMatchMultiple(t *testing.T) {
	children := []shared.Node{
		&shared.Frame{ID: "p1", Name: "Front"},
		&shared.Frame{ID: "p2", Name: "Back"},
		&shared.Frame{ID: "p3", Name: "Extra"},
	}

	filtered, err := filterPages(children, "Front,Back")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filtered) != 2 {
		t.Fatalf("expected 2 pages, got %d", len(filtered))
	}
}

func TestFilterPagesNoMatch(t *testing.T) {
	children := []shared.Node{
		&shared.Frame{ID: "p1", Name: "Front"},
	}

	_, err := filterPages(children, "NonExistent")
	if err == nil {
		t.Fatal("expected error for no matching pages")
	}
}

func TestValidateService(t *testing.T) {
	input := `{
		"version": "1.0",
		"variables": {"$primary": {"type": "color", "value": "#FF0000"}},
		"children": [
			{"type": "frame", "id": "p1", "name": "page", "width": 800, "height": 1000, "children": []}
		]
	}`

	svc := newTestValidateService()
	result, err := svc.Validate(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.PageCount != 1 {
		t.Errorf("expected 1 page, got %d", result.PageCount)
	}
	if result.VariableCount != 1 {
		t.Errorf("expected 1 variable, got %d", result.VariableCount)
	}
}

func TestValidateServiceParseError(t *testing.T) {
	svc := newTestValidateService()
	_, err := svc.Validate(strings.NewReader("invalid json"))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestInfoService(t *testing.T) {
	input := `{
		"version": "1.0",
		"variables": {},
		"children": [
			{
				"type": "frame", "id": "p1", "name": "Front", "width": 800, "height": 1000,
				"children": [
					{"type": "text", "id": "t1", "name": "title", "content": "Hello", "fontFamily": "Inter", "fontSize": 16, "fontWeight": "400"}
				]
			}
		]
	}`

	svc := newTestInfoService()
	info, err := svc.GetInfo(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Version != "1.0" {
		t.Errorf("expected version '1.0', got '%s'", info.Version)
	}
	if len(info.Pages) != 1 {
		t.Fatalf("expected 1 page, got %d", len(info.Pages))
	}
	if info.Pages[0].Name != "Front" {
		t.Errorf("expected page name 'Front', got '%s'", info.Pages[0].Name)
	}
	if len(info.Fonts) != 1 || info.Fonts[0] != "Inter" {
		t.Errorf("expected fonts [Inter], got %v", info.Fonts)
	}
}

func TestInfoServiceCollectFonts(t *testing.T) {
	nodes := []shared.Node{
		&shared.Frame{
			ID: "f1", Name: "page",
			Children: []shared.Node{
				&shared.Text{ID: "t1", Name: "a", FontFamily: "Inter"},
				&shared.Text{ID: "t2", Name: "b", FontFamily: "Montserrat"},
				&shared.Text{ID: "t3", Name: "c", FontFamily: "Inter"}, // duplicate
				&shared.Frame{
					ID: "f2", Name: "inner",
					Children: []shared.Node{
						&shared.Text{ID: "t4", Name: "d", FontFamily: "Playfair Display"},
					},
				},
			},
		},
	}

	fonts := collectFonts(nodes)
	if len(fonts) != 3 {
		t.Fatalf("expected 3 unique fonts, got %d: %v", len(fonts), fonts)
	}
	expected := []string{"Inter", "Montserrat", "Playfair Display"}
	for i, f := range fonts {
		if f != expected[i] {
			t.Errorf("expected fonts[%d] '%s', got '%s'", i, expected[i], f)
		}
	}
}

func TestInfoServiceCollectFontsEmpty(t *testing.T) {
	nodes := []shared.Node{
		&shared.Frame{ID: "f1", Name: "empty"},
	}
	fonts := collectFonts(nodes)
	if len(fonts) != 0 {
		t.Errorf("expected 0 fonts, got %d", len(fonts))
	}
}

func TestInfoServiceCollectFontsSkipsEmptyFamily(t *testing.T) {
	nodes := []shared.Node{
		&shared.Text{ID: "t1", Name: "a", FontFamily: ""},
		&shared.Text{ID: "t2", Name: "b", FontFamily: "Inter"},
	}
	fonts := collectFonts(nodes)
	if len(fonts) != 1 {
		t.Fatalf("expected 1 font, got %d", len(fonts))
	}
}

func TestDetectMissingFonts(t *testing.T) {
	input := `{
		"version": "1.0",
		"variables": {},
		"children": [
			{
				"type": "frame", "id": "p1", "name": "page", "width": 800, "height": 1000,
				"children": [
					{"type": "text", "id": "t1", "name": "title", "content": "Hello", "fontFamily": "NonExistent", "fontSize": 16, "fontWeight": "400"}
				]
			}
		]
	}`

	svc := newTestRenderService()
	missing, doc, err := svc.DetectMissingFonts(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if doc == nil {
		t.Fatal("expected non-nil document")
	}
	if len(missing) != 1 {
		t.Fatalf("expected 1 missing font, got %d", len(missing))
	}
	if missing[0].Family != "NonExistent" {
		t.Errorf("expected family 'NonExistent', got '%s'", missing[0].Family)
	}
}

func TestRenderDocumentEmptyPage(t *testing.T) {
	doc := &shared.Document{
		Children: []shared.Node{
			&shared.Frame{
				ID: "p1", Name: "page",
				Width: shared.FixedDimension(800), Height: shared.FixedDimension(1000),
			},
		},
	}

	svc := newTestRenderService()
	var buf bytes.Buffer
	result, err := svc.RenderDocument(doc, &buf, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.PageCount != 1 {
		t.Errorf("expected 1 page, got %d", result.PageCount)
	}
	if !bytes.HasPrefix(buf.Bytes(), []byte("%PDF")) {
		t.Error("output does not start with %PDF header")
	}
}

// Test helpers — use real infrastructure (lightweight, no I/O)

func newTestValidateService() *ValidateService {
	return NewValidateService(
		newJSONParser(),
		newVariableResolver(),
	)
}

func newTestInfoService() *InfoService {
	return NewInfoService(newJSONParser())
}

func newTestRenderService() *RenderService {
	fl := newEmptyFontLoader()
	return NewRenderService(
		newJSONParser(),
		newVariableResolver(),
		fl,
		nil, // imageLoader not needed for basic tests
		newFlexboxEngine(),
		newMeasurer(fl),
		newPDFRenderer(fl),
	)
}
