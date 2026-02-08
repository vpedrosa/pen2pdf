package domain_test

import (
	"testing"

	"github.com/vpedrosa/pen2pdf/internal/shared/domain"
)

func TestFilterPagesByNameSingle(t *testing.T) {
	children := []domain.Node{
		&domain.Frame{ID: "p1", Name: "Front"},
		&domain.Frame{ID: "p2", Name: "Back"},
	}

	filtered, err := domain.FilterPagesByName(children, "Front")
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

func TestFilterPagesByNameMultiple(t *testing.T) {
	children := []domain.Node{
		&domain.Frame{ID: "p1", Name: "Front"},
		&domain.Frame{ID: "p2", Name: "Back"},
		&domain.Frame{ID: "p3", Name: "Extra"},
	}

	filtered, err := domain.FilterPagesByName(children, "Front,Back")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filtered) != 2 {
		t.Fatalf("expected 2 pages, got %d", len(filtered))
	}
}

func TestFilterPagesByNameNoMatch(t *testing.T) {
	children := []domain.Node{
		&domain.Frame{ID: "p1", Name: "Front"},
	}

	_, err := domain.FilterPagesByName(children, "NonExistent")
	if err == nil {
		t.Fatal("expected error for no matching pages")
	}
}
