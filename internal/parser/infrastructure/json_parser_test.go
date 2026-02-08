package infrastructure_test

import (
	"strings"
	"testing"

	parser "github.com/vpedrosa/pen2pdf/internal/parser/domain"
	"github.com/vpedrosa/pen2pdf/internal/parser/infrastructure"
	shared "github.com/vpedrosa/pen2pdf/internal/shared/domain"
)

func TestJSONParserImplementsPort(t *testing.T) {
	var _ parser.Parser = infrastructure.NewJSONParser()
}

func TestParseMinimalDocument(t *testing.T) {
	input := `{"version": "2.7", "children": []}`
	doc := mustParse(t, input)

	if doc.Version != "2.7" {
		t.Errorf("expected version '2.7', got '%s'", doc.Version)
	}
	if len(doc.Children) != 0 {
		t.Errorf("expected 0 children, got %d", len(doc.Children))
	}
}

func TestParseFrameNode(t *testing.T) {
	input := `{
		"version": "1.0",
		"children": [{
			"type": "frame",
			"id": "f1",
			"name": "page",
			"x": 96,
			"y": 100,
			"width": 800,
			"height": 1000,
			"clip": true,
			"layout": "vertical",
			"gap": 20,
			"justifyContent": "center",
			"alignItems": "center",
			"cornerRadius": 16
		}]
	}`
	doc := mustParse(t, input)

	if len(doc.Children) != 1 {
		t.Fatalf("expected 1 child, got %d", len(doc.Children))
	}

	frame, ok := doc.Children[0].(*shared.Frame)
	if !ok {
		t.Fatalf("expected *Frame, got %T", doc.Children[0])
	}

	if frame.ID != "f1" {
		t.Errorf("expected ID 'f1', got '%s'", frame.ID)
	}
	if frame.Name != "page" {
		t.Errorf("expected Name 'page', got '%s'", frame.Name)
	}
	if frame.X != 96 {
		t.Errorf("expected X 96, got %f", frame.X)
	}
	if frame.Y != 100 {
		t.Errorf("expected Y 100, got %f", frame.Y)
	}
	if frame.Width.Value != 800 {
		t.Errorf("expected Width 800, got %f", frame.Width.Value)
	}
	if frame.Height.Value != 1000 {
		t.Errorf("expected Height 1000, got %f", frame.Height.Value)
	}
	if !frame.Clip {
		t.Error("expected Clip true")
	}
	if frame.Layout != "vertical" {
		t.Errorf("expected Layout 'vertical', got '%s'", frame.Layout)
	}
	if frame.Gap != 20 {
		t.Errorf("expected Gap 20, got %f", frame.Gap)
	}
	if frame.JustifyContent != "center" {
		t.Errorf("expected JustifyContent 'center', got '%s'", frame.JustifyContent)
	}
	if frame.AlignItems != "center" {
		t.Errorf("expected AlignItems 'center', got '%s'", frame.AlignItems)
	}
	if frame.CornerRadius != 16 {
		t.Errorf("expected CornerRadius 16, got %f", frame.CornerRadius)
	}
}

func TestParseTextNode(t *testing.T) {
	input := `{
		"version": "1.0",
		"children": [{
			"type": "text",
			"id": "t1",
			"name": "label",
			"content": "Hello World",
			"fill": "#FF0000",
			"fontFamily": "Inter",
			"fontSize": 16,
			"fontWeight": "700",
			"fontStyle": "italic",
			"letterSpacing": 2,
			"lineHeight": 1.5,
			"textAlign": "center"
		}]
	}`
	doc := mustParse(t, input)

	text, ok := doc.Children[0].(*shared.Text)
	if !ok {
		t.Fatalf("expected *Text, got %T", doc.Children[0])
	}

	if text.ID != "t1" {
		t.Errorf("expected ID 't1', got '%s'", text.ID)
	}
	if text.Content != "Hello World" {
		t.Errorf("expected Content 'Hello World', got '%s'", text.Content)
	}
	if text.Fill != "#FF0000" {
		t.Errorf("expected Fill '#FF0000', got '%s'", text.Fill)
	}
	if text.FontFamily != "Inter" {
		t.Errorf("expected FontFamily 'Inter', got '%s'", text.FontFamily)
	}
	if text.FontSize != 16 {
		t.Errorf("expected FontSize 16, got %f", text.FontSize)
	}
	if text.FontWeight != "700" {
		t.Errorf("expected FontWeight '700', got '%s'", text.FontWeight)
	}
	if text.FontStyle != "italic" {
		t.Errorf("expected FontStyle 'italic', got '%s'", text.FontStyle)
	}
	if text.LetterSpacing != 2 {
		t.Errorf("expected LetterSpacing 2, got %f", text.LetterSpacing)
	}
	if text.LineHeight != 1.5 {
		t.Errorf("expected LineHeight 1.5, got %f", text.LineHeight)
	}
	if text.TextAlign != "center" {
		t.Errorf("expected TextAlign 'center', got '%s'", text.TextAlign)
	}
}

func TestParseTextWithFixedWidth(t *testing.T) {
	input := `{
		"version": "1.0",
		"children": [{
			"type": "text",
			"id": "t1",
			"name": "body",
			"content": "long text",
			"textGrowth": "fixed-width",
			"width": 600
		}]
	}`
	doc := mustParse(t, input)

	text := doc.Children[0].(*shared.Text)
	if text.Width.Value != 600 {
		t.Errorf("expected Width 600, got %f", text.Width.Value)
	}
	if text.TextGrowth != "fixed-width" {
		t.Errorf("expected TextGrowth 'fixed-width', got '%s'", text.TextGrowth)
	}
}

func TestParseFillSolidString(t *testing.T) {
	input := `{
		"version": "1.0",
		"children": [{
			"type": "frame",
			"id": "f1",
			"name": "box",
			"fill": "#FF6B35"
		}]
	}`
	doc := mustParse(t, input)

	frame := doc.Children[0].(*shared.Frame)
	if frame.Fill == nil {
		t.Fatal("expected Fill to be set")
	}
	if frame.Fill.Type != shared.FillSolid {
		t.Errorf("expected FillSolid, got '%s'", frame.Fill.Type)
	}
	if frame.Fill.Color != "#FF6B35" {
		t.Errorf("expected color '#FF6B35', got '%s'", frame.Fill.Color)
	}
}

func TestParseFillSolidWithAlpha(t *testing.T) {
	input := `{
		"version": "1.0",
		"children": [{
			"type": "frame",
			"id": "f1",
			"name": "overlay",
			"fill": "#000000BB"
		}]
	}`
	doc := mustParse(t, input)

	frame := doc.Children[0].(*shared.Frame)
	if frame.Fill.Color != "#000000BB" {
		t.Errorf("expected color '#000000BB', got '%s'", frame.Fill.Color)
	}
}

func TestParseFillWithVariable(t *testing.T) {
	input := `{
		"version": "1.0",
		"children": [{
			"type": "frame",
			"id": "f1",
			"name": "div",
			"fill": "$primary-color"
		}]
	}`
	doc := mustParse(t, input)

	frame := doc.Children[0].(*shared.Frame)
	if frame.Fill.Color != "$primary-color" {
		t.Errorf("expected color '$primary-color', got '%s'", frame.Fill.Color)
	}
}

func TestParseFillImage(t *testing.T) {
	input := `{
		"version": "1.0",
		"children": [{
			"type": "frame",
			"id": "f1",
			"name": "bg",
			"fill": {
				"type": "image",
				"url": "./images/bg.jpg",
				"mode": "fill",
				"opacity": 0.3,
				"enabled": true
			}
		}]
	}`
	doc := mustParse(t, input)

	frame := doc.Children[0].(*shared.Frame)
	if frame.Fill == nil {
		t.Fatal("expected Fill to be set")
	}
	if frame.Fill.Type != shared.FillImage {
		t.Errorf("expected FillImage, got '%s'", frame.Fill.Type)
	}
	if frame.Fill.URL != "./images/bg.jpg" {
		t.Errorf("expected URL './images/bg.jpg', got '%s'", frame.Fill.URL)
	}
	if frame.Fill.Mode != "fill" {
		t.Errorf("expected mode 'fill', got '%s'", frame.Fill.Mode)
	}
	if frame.Fill.Opacity != 0.3 {
		t.Errorf("expected opacity 0.3, got %f", frame.Fill.Opacity)
	}
	if !frame.Fill.Enabled {
		t.Error("expected Enabled true")
	}
}

func TestParseFillAbsent(t *testing.T) {
	input := `{
		"version": "1.0",
		"children": [{
			"type": "frame",
			"id": "f1",
			"name": "nofill"
		}]
	}`
	doc := mustParse(t, input)

	frame := doc.Children[0].(*shared.Frame)
	if frame.Fill != nil {
		t.Errorf("expected nil Fill, got %+v", frame.Fill)
	}
}

func TestParseDimensionFixed(t *testing.T) {
	input := `{
		"version": "1.0",
		"children": [{
			"type": "frame",
			"id": "f1",
			"name": "fixed",
			"width": 800,
			"height": 1000
		}]
	}`
	doc := mustParse(t, input)

	frame := doc.Children[0].(*shared.Frame)
	if frame.Width.Value != 800 || frame.Width.FillContainer {
		t.Errorf("expected fixed width 800, got %+v", frame.Width)
	}
	if frame.Height.Value != 1000 || frame.Height.FillContainer {
		t.Errorf("expected fixed height 1000, got %+v", frame.Height)
	}
}

func TestParseDimensionFillContainer(t *testing.T) {
	input := `{
		"version": "1.0",
		"children": [{
			"type": "frame",
			"id": "f1",
			"name": "stretch",
			"width": "fill_container",
			"height": "fill_container"
		}]
	}`
	doc := mustParse(t, input)

	frame := doc.Children[0].(*shared.Frame)
	if !frame.Width.FillContainer {
		t.Error("expected Width.FillContainer true")
	}
	if !frame.Height.FillContainer {
		t.Error("expected Height.FillContainer true")
	}
}

func TestParseDimensionAbsent(t *testing.T) {
	input := `{
		"version": "1.0",
		"children": [{
			"type": "frame",
			"id": "f1",
			"name": "auto"
		}]
	}`
	doc := mustParse(t, input)

	frame := doc.Children[0].(*shared.Frame)
	if frame.Width.Value != 0 || frame.Width.FillContainer {
		t.Errorf("expected zero dimension, got %+v", frame.Width)
	}
}

func TestParsePaddingUniform(t *testing.T) {
	input := `{
		"version": "1.0",
		"children": [{
			"type": "frame",
			"id": "f1",
			"name": "box",
			"padding": 40
		}]
	}`
	doc := mustParse(t, input)

	frame := doc.Children[0].(*shared.Frame)
	p := frame.Padding
	if p.Top != 40 || p.Right != 40 || p.Bottom != 40 || p.Left != 40 {
		t.Errorf("expected uniform padding 40, got %+v", p)
	}
}

func TestParsePaddingTwoValues(t *testing.T) {
	input := `{
		"version": "1.0",
		"children": [{
			"type": "frame",
			"id": "f1",
			"name": "badge",
			"padding": [6, 16]
		}]
	}`
	doc := mustParse(t, input)

	frame := doc.Children[0].(*shared.Frame)
	p := frame.Padding
	if p.Top != 6 || p.Right != 16 || p.Bottom != 6 || p.Left != 16 {
		t.Errorf("expected padding [6,16,6,16], got %+v", p)
	}
}

func TestParsePaddingFourValues(t *testing.T) {
	input := `{
		"version": "1.0",
		"children": [{
			"type": "frame",
			"id": "f1",
			"name": "header",
			"padding": [0, 0, 10, 0]
		}]
	}`
	doc := mustParse(t, input)

	frame := doc.Children[0].(*shared.Frame)
	p := frame.Padding
	if p.Top != 0 || p.Right != 0 || p.Bottom != 10 || p.Left != 0 {
		t.Errorf("expected padding [0,0,10,0], got %+v", p)
	}
}

func TestParsePaddingAbsent(t *testing.T) {
	input := `{
		"version": "1.0",
		"children": [{
			"type": "frame",
			"id": "f1",
			"name": "nopad"
		}]
	}`
	doc := mustParse(t, input)

	frame := doc.Children[0].(*shared.Frame)
	p := frame.Padding
	if p.Top != 0 || p.Right != 0 || p.Bottom != 0 || p.Left != 0 {
		t.Errorf("expected zero padding, got %+v", p)
	}
}

func TestParseNestedChildren(t *testing.T) {
	input := `{
		"version": "1.0",
		"children": [{
			"type": "frame",
			"id": "parent",
			"name": "container",
			"children": [
				{
					"type": "frame",
					"id": "child1",
					"name": "inner",
					"children": [
						{"type": "text", "id": "t1", "name": "deep", "content": "hello"}
					]
				},
				{"type": "text", "id": "child2", "name": "label", "content": "world"}
			]
		}]
	}`
	doc := mustParse(t, input)

	parent := doc.Children[0].(*shared.Frame)
	if len(parent.Children) != 2 {
		t.Fatalf("expected 2 children, got %d", len(parent.Children))
	}

	inner := parent.Children[0].(*shared.Frame)
	if len(inner.Children) != 1 {
		t.Fatalf("expected 1 grandchild, got %d", len(inner.Children))
	}

	deep := inner.Children[0].(*shared.Text)
	if deep.Content != "hello" {
		t.Errorf("expected content 'hello', got '%s'", deep.Content)
	}

	label := parent.Children[1].(*shared.Text)
	if label.Content != "world" {
		t.Errorf("expected content 'world', got '%s'", label.Content)
	}
}

func TestParseVariables(t *testing.T) {
	input := `{
		"version": "1.0",
		"children": [],
		"variables": {
			"primary-color": {"type": "color", "value": "#FF6B35"},
			"font-body": {"type": "string", "value": "Open Sans"},
			"font-size-lg": {"type": "number", "value": 16}
		}
	}`
	doc := mustParse(t, input)

	if len(doc.Variables) != 3 {
		t.Fatalf("expected 3 variables, got %d", len(doc.Variables))
	}

	color := doc.Variables["primary-color"]
	if color.Type != shared.VariableColor || color.Value != "#FF6B35" {
		t.Errorf("unexpected primary-color: %+v", color)
	}

	font := doc.Variables["font-body"]
	if font.Type != shared.VariableString || font.Value != "Open Sans" {
		t.Errorf("unexpected font-body: %+v", font)
	}

	size := doc.Variables["font-size-lg"]
	if size.Type != shared.VariableNumber || size.Value != 16.0 {
		t.Errorf("unexpected font-size-lg: %+v", size)
	}
}

func TestParseNoVariables(t *testing.T) {
	input := `{"version": "1.0", "children": []}`
	doc := mustParse(t, input)

	if doc.Variables != nil {
		t.Errorf("expected nil variables, got %v", doc.Variables)
	}
}

func TestParseInvalidJSON(t *testing.T) {
	p := infrastructure.NewJSONParser()
	_, err := p.Parse(strings.NewReader("{invalid"))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
	if !strings.Contains(err.Error(), "invalid JSON") {
		t.Errorf("expected 'invalid JSON' in error, got: %s", err)
	}
}

func TestParseUnknownNodeType(t *testing.T) {
	input := `{
		"version": "1.0",
		"children": [{"type": "rectangle", "id": "r1"}]
	}`
	p := infrastructure.NewJSONParser()
	_, err := p.Parse(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected error for unknown node type")
	}
	if !strings.Contains(err.Error(), "unknown node type") {
		t.Errorf("expected 'unknown node type' in error, got: %s", err)
	}
}

func TestParseInvalidDimensionValue(t *testing.T) {
	input := `{
		"version": "1.0",
		"children": [{"type": "frame", "id": "f1", "name": "x", "width": "stretch"}]
	}`
	p := infrastructure.NewJSONParser()
	_, err := p.Parse(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected error for invalid dimension")
	}
	if !strings.Contains(err.Error(), "unknown dimension value") {
		t.Errorf("expected 'unknown dimension value' in error, got: %s", err)
	}
}

func TestParseInvalidPaddingArray(t *testing.T) {
	input := `{
		"version": "1.0",
		"children": [{"type": "frame", "id": "f1", "name": "x", "padding": [1, 2, 3]}]
	}`
	p := infrastructure.NewJSONParser()
	_, err := p.Parse(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected error for invalid padding array")
	}
	if !strings.Contains(err.Error(), "2 or 4 elements") {
		t.Errorf("expected '2 or 4 elements' in error, got: %s", err)
	}
}

func TestParseUnknownVariableType(t *testing.T) {
	input := `{
		"version": "1.0",
		"children": [],
		"variables": {"x": {"type": "boolean", "value": true}}
	}`
	p := infrastructure.NewJSONParser()
	_, err := p.Parse(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected error for unknown variable type")
	}
	if !strings.Contains(err.Error(), "unknown variable type") {
		t.Errorf("expected 'unknown variable type' in error, got: %s", err)
	}
}

func TestParseUnknownFillType(t *testing.T) {
	input := `{
		"version": "1.0",
		"children": [{"type": "frame", "id": "f1", "name": "x", "fill": {"type": "gradient"}}]
	}`
	p := infrastructure.NewJSONParser()
	_, err := p.Parse(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected error for unknown fill type")
	}
	if !strings.Contains(err.Error(), "unknown fill type") {
		t.Errorf("expected 'unknown fill type' in error, got: %s", err)
	}
}

func TestParseExampleFile(t *testing.T) {
	// Integration-style test using a realistic multi-page document
	input := `{
		"version": "2.7",
		"children": [
			{
				"type": "frame",
				"id": "page1",
				"name": "Front",
				"x": 96, "y": 100,
				"width": 800, "height": 1000,
				"clip": true,
				"fill": {"type": "image", "url": "./bg.jpg", "mode": "fill", "enabled": true},
				"children": [
					{
						"type": "frame",
						"id": "overlay",
						"name": "overlay",
						"width": "fill_container",
						"height": "fill_container",
						"fill": "#000000BB",
						"children": [
							{
								"type": "text",
								"id": "title",
								"name": "headline",
								"fill": "#FFFFFF",
								"content": "Hello World",
								"fontFamily": "Inter",
								"fontSize": 48,
								"fontWeight": "900"
							}
						]
					}
				]
			},
			{
				"type": "frame",
				"id": "page2",
				"name": "Back",
				"width": 800, "height": 1000,
				"fill": "#FFFFFF",
				"padding": [50, 40, 50, 40],
				"layout": "vertical",
				"gap": 20,
				"children": [
					{
						"type": "text",
						"id": "body",
						"name": "terms",
						"fill": "$text-primary",
						"content": "Terms and conditions apply.",
						"textGrowth": "fixed-width",
						"width": 600,
						"fontFamily": "Open Sans",
						"fontSize": 12,
						"lineHeight": 1.5
					}
				]
			}
		],
		"variables": {
			"primary-color": {"type": "color", "value": "#FF6B35"},
			"text-primary": {"type": "color", "value": "#2C3E50"},
			"font-body": {"type": "string", "value": "Open Sans"},
			"spacing-lg": {"type": "number", "value": 20}
		}
	}`
	doc := mustParse(t, input)

	if doc.Version != "2.7" {
		t.Errorf("expected version '2.7', got '%s'", doc.Version)
	}
	if len(doc.Children) != 2 {
		t.Fatalf("expected 2 pages, got %d", len(doc.Children))
	}

	// Front page
	front := doc.Children[0].(*shared.Frame)
	if front.Name != "Front" {
		t.Errorf("expected 'Front', got '%s'", front.Name)
	}
	if front.Fill.Type != shared.FillImage {
		t.Errorf("expected image fill, got '%s'", front.Fill.Type)
	}
	overlay := front.Children[0].(*shared.Frame)
	if overlay.Fill.Color != "#000000BB" {
		t.Errorf("expected overlay color '#000000BB', got '%s'", overlay.Fill.Color)
	}
	headline := overlay.Children[0].(*shared.Text)
	if headline.Content != "Hello World" {
		t.Errorf("expected 'Hello World', got '%s'", headline.Content)
	}

	// Back page
	back := doc.Children[1].(*shared.Frame)
	if back.Name != "Back" {
		t.Errorf("expected 'Back', got '%s'", back.Name)
	}
	if back.Padding.Top != 50 || back.Padding.Right != 40 {
		t.Errorf("expected padding [50,40,50,40], got %+v", back.Padding)
	}
	terms := back.Children[0].(*shared.Text)
	if terms.Fill != "$text-primary" {
		t.Errorf("expected fill '$text-primary', got '%s'", terms.Fill)
	}

	// Variables
	if len(doc.Variables) != 4 {
		t.Errorf("expected 4 variables, got %d", len(doc.Variables))
	}
}

func mustParse(t *testing.T, input string) *shared.Document {
	t.Helper()
	p := infrastructure.NewJSONParser()
	doc, err := p.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	return doc
}
