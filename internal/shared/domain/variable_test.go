package domain_test

import (
	"testing"

	"github.com/vpedrosa/pen2pdf/internal/shared/domain"
)

func TestVariableColor(t *testing.T) {
	v := domain.Variable{Type: domain.VariableColor, Value: "#FF6B35"}
	if v.Type != domain.VariableColor {
		t.Errorf("expected type color, got %s", v.Type)
	}
	if v.Value != "#FF6B35" {
		t.Errorf("expected value '#FF6B35', got '%v'", v.Value)
	}
}

func TestVariableString(t *testing.T) {
	v := domain.Variable{Type: domain.VariableString, Value: "Montserrat"}
	if v.Type != domain.VariableString {
		t.Errorf("expected type string, got %s", v.Type)
	}
	if v.Value != "Montserrat" {
		t.Errorf("expected value 'Montserrat', got '%v'", v.Value)
	}
}

func TestVariableNumber(t *testing.T) {
	v := domain.Variable{Type: domain.VariableNumber, Value: 16.0}
	if v.Type != domain.VariableNumber {
		t.Errorf("expected type number, got %s", v.Type)
	}
	if v.Value != 16.0 {
		t.Errorf("expected value 16.0, got '%v'", v.Value)
	}
}
