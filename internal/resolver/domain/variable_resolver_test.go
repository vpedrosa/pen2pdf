package domain_test

import (
	"strings"
	"testing"

	resolver "github.com/vpedrosa/pen2pdf/internal/resolver/domain"
	shared "github.com/vpedrosa/pen2pdf/internal/shared/domain"
)

func TestVariableResolverImplementsPort(t *testing.T) {
	var _ resolver.Resolver = resolver.NewVariableResolver()
}

func TestResolveFrameFillVariable(t *testing.T) {
	doc := &shared.Document{
		Children: []shared.Node{
			&shared.Frame{
				ID:   "f1",
				Name: "box",
				Fill: shared.SolidFill("$primary-color"),
			},
		},
		Variables: map[string]shared.Variable{
			"primary-color": {Type: shared.VariableColor, Value: "#FF6B35"},
		},
	}

	r := resolver.NewVariableResolver()
	if err := r.Resolve(doc); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	frame := doc.Children[0].(*shared.Frame)
	if frame.Fill.Color != "#FF6B35" {
		t.Errorf("expected '#FF6B35', got '%s'", frame.Fill.Color)
	}
}

func TestResolveTextFillVariable(t *testing.T) {
	doc := &shared.Document{
		Children: []shared.Node{
			&shared.Text{
				ID:   "t1",
				Name: "label",
				Fill: "$text-primary",
			},
		},
		Variables: map[string]shared.Variable{
			"text-primary": {Type: shared.VariableColor, Value: "#2C3E50"},
		},
	}

	r := resolver.NewVariableResolver()
	if err := r.Resolve(doc); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	text := doc.Children[0].(*shared.Text)
	if text.Fill != "#2C3E50" {
		t.Errorf("expected '#2C3E50', got '%s'", text.Fill)
	}
}

func TestResolveNestedNodes(t *testing.T) {
	doc := &shared.Document{
		Children: []shared.Node{
			&shared.Frame{
				ID:   "parent",
				Name: "page",
				Fill: shared.SolidFill("$bg-color"),
				Children: []shared.Node{
					&shared.Frame{
						ID:   "child",
						Name: "inner",
						Fill: shared.SolidFill("$accent-color"),
						Children: []shared.Node{
							&shared.Text{ID: "t1", Name: "label", Fill: "$text-color"},
						},
					},
				},
			},
		},
		Variables: map[string]shared.Variable{
			"bg-color":     {Type: shared.VariableColor, Value: "#FFFFFF"},
			"accent-color": {Type: shared.VariableColor, Value: "#FF0000"},
			"text-color":   {Type: shared.VariableColor, Value: "#000000"},
		},
	}

	r := resolver.NewVariableResolver()
	if err := r.Resolve(doc); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	parent := doc.Children[0].(*shared.Frame)
	if parent.Fill.Color != "#FFFFFF" {
		t.Errorf("expected parent fill '#FFFFFF', got '%s'", parent.Fill.Color)
	}

	child := parent.Children[0].(*shared.Frame)
	if child.Fill.Color != "#FF0000" {
		t.Errorf("expected child fill '#FF0000', got '%s'", child.Fill.Color)
	}

	text := child.Children[0].(*shared.Text)
	if text.Fill != "#000000" {
		t.Errorf("expected text fill '#000000', got '%s'", text.Fill)
	}
}

func TestResolveNonVariableUnchanged(t *testing.T) {
	doc := &shared.Document{
		Children: []shared.Node{
			&shared.Frame{
				ID:   "f1",
				Name: "box",
				Fill: shared.SolidFill("#FF6B35"),
			},
			&shared.Text{ID: "t1", Name: "label", Fill: "#000000"},
		},
		Variables: map[string]shared.Variable{},
	}

	r := resolver.NewVariableResolver()
	if err := r.Resolve(doc); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	frame := doc.Children[0].(*shared.Frame)
	if frame.Fill.Color != "#FF6B35" {
		t.Errorf("expected '#FF6B35', got '%s'", frame.Fill.Color)
	}

	text := doc.Children[1].(*shared.Text)
	if text.Fill != "#000000" {
		t.Errorf("expected '#000000', got '%s'", text.Fill)
	}
}

func TestResolveNoVariablesMap(t *testing.T) {
	doc := &shared.Document{
		Children: []shared.Node{
			&shared.Frame{ID: "f1", Name: "box"},
		},
	}

	r := resolver.NewVariableResolver()
	if err := r.Resolve(doc); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestResolveUndefinedVariable(t *testing.T) {
	doc := &shared.Document{
		Children: []shared.Node{
			&shared.Text{ID: "t1", Name: "label", Fill: "$missing-var"},
		},
		Variables: map[string]shared.Variable{},
	}

	r := resolver.NewVariableResolver()
	err := r.Resolve(doc)
	if err == nil {
		t.Fatal("expected error for undefined variable")
	}
	if !strings.Contains(err.Error(), "undefined variable") {
		t.Errorf("expected 'undefined variable' in error, got: %s", err)
	}
}

func TestResolveNilFillSkipped(t *testing.T) {
	doc := &shared.Document{
		Children: []shared.Node{
			&shared.Frame{ID: "f1", Name: "nofill"},
		},
		Variables: map[string]shared.Variable{},
	}

	r := resolver.NewVariableResolver()
	if err := r.Resolve(doc); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestResolveImageFillUnchanged(t *testing.T) {
	doc := &shared.Document{
		Children: []shared.Node{
			&shared.Frame{
				ID:   "f1",
				Name: "bg",
				Fill: shared.ImageFill("./bg.jpg", "fill", 1.0, true),
			},
		},
		Variables: map[string]shared.Variable{},
	}

	r := resolver.NewVariableResolver()
	if err := r.Resolve(doc); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	frame := doc.Children[0].(*shared.Frame)
	if frame.Fill.Type != shared.FillImage {
		t.Errorf("expected FillImage, got '%s'", frame.Fill.Type)
	}
	if frame.Fill.URL != "./bg.jpg" {
		t.Errorf("expected URL unchanged, got '%s'", frame.Fill.URL)
	}
}

func TestResolveEmptyFillStringSkipped(t *testing.T) {
	doc := &shared.Document{
		Children: []shared.Node{
			&shared.Text{ID: "t1", Name: "label", Fill: ""},
		},
		Variables: map[string]shared.Variable{},
	}

	r := resolver.NewVariableResolver()
	if err := r.Resolve(doc); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	text := doc.Children[0].(*shared.Text)
	if text.Fill != "" {
		t.Errorf("expected empty fill, got '%s'", text.Fill)
	}
}

func TestResolveMultipleVariablesSameDocument(t *testing.T) {
	doc := &shared.Document{
		Children: []shared.Node{
			&shared.Frame{
				ID:   "f1",
				Name: "page",
				Fill: shared.SolidFill("$primary-color"),
				Children: []shared.Node{
					&shared.Frame{ID: "f2", Name: "div", Fill: shared.SolidFill("$secondary-color")},
					&shared.Text{ID: "t1", Name: "title", Fill: "$text-primary"},
					&shared.Text{ID: "t2", Name: "subtitle", Fill: "$text-muted"},
				},
			},
		},
		Variables: map[string]shared.Variable{
			"primary-color":   {Type: shared.VariableColor, Value: "#FF6B35"},
			"secondary-color": {Type: shared.VariableColor, Value: "#16A085"},
			"text-primary":    {Type: shared.VariableColor, Value: "#2C3E50"},
			"text-muted":      {Type: shared.VariableColor, Value: "#7F8C8D"},
		},
	}

	r := resolver.NewVariableResolver()
	if err := r.Resolve(doc); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	frame := doc.Children[0].(*shared.Frame)
	if frame.Fill.Color != "#FF6B35" {
		t.Errorf("expected '#FF6B35', got '%s'", frame.Fill.Color)
	}

	innerFrame := frame.Children[0].(*shared.Frame)
	if innerFrame.Fill.Color != "#16A085" {
		t.Errorf("expected '#16A085', got '%s'", innerFrame.Fill.Color)
	}

	title := frame.Children[1].(*shared.Text)
	if title.Fill != "#2C3E50" {
		t.Errorf("expected '#2C3E50', got '%s'", title.Fill)
	}

	subtitle := frame.Children[2].(*shared.Text)
	if subtitle.Fill != "#7F8C8D" {
		t.Errorf("expected '#7F8C8D', got '%s'", subtitle.Fill)
	}
}
