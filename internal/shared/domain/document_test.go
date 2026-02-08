package domain_test

import (
	"testing"

	"github.com/vpedrosa/pen2pdf/internal/shared/domain"
)

func TestDocumentWithPagesAndVariables(t *testing.T) {
	doc := &domain.Document{
		Version: "2.7",
		Children: []domain.Node{
			&domain.Frame{ID: "page1", Name: "Front", Width: domain.FixedDimension(800), Height: domain.FixedDimension(1000)},
			&domain.Frame{ID: "page2", Name: "Back", Width: domain.FixedDimension(800), Height: domain.FixedDimension(1000)},
		},
		Variables: map[string]domain.Variable{
			"primary-color": {Type: domain.VariableColor, Value: "#FF6B35"},
			"font-body":     {Type: domain.VariableString, Value: "Open Sans"},
			"font-size-lg":  {Type: domain.VariableNumber, Value: 16.0},
		},
	}

	if doc.Version != "2.7" {
		t.Errorf("expected version '2.7', got '%s'", doc.Version)
	}
	if len(doc.Children) != 2 {
		t.Fatalf("expected 2 children, got %d", len(doc.Children))
	}
	if doc.Children[0].GetName() != "Front" {
		t.Errorf("expected first page 'Front', got '%s'", doc.Children[0].GetName())
	}
	if len(doc.Variables) != 3 {
		t.Fatalf("expected 3 variables, got %d", len(doc.Variables))
	}
	if doc.Variables["primary-color"].Type != domain.VariableColor {
		t.Error("expected primary-color to be a color variable")
	}
}
