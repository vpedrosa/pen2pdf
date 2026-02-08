package cmd

import (
	"testing"

	shared "github.com/vpedrosa/pen2pdf/internal/shared/domain"
)

func TestRenderCommandRegistered(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Use == "render [input.pen]" {
			found = true
			break
		}
	}
	if !found {
		t.Error("render command not registered")
	}
}

func TestRenderCommandHasOutputFlag(t *testing.T) {
	f := renderCmd.Flags().Lookup("output")
	if f == nil {
		t.Fatal("expected --output flag")
	}
	if f.Shorthand != "o" {
		t.Errorf("expected shorthand 'o', got '%s'", f.Shorthand)
	}
}

func TestRenderCommandHasPagesFlag(t *testing.T) {
	f := renderCmd.Flags().Lookup("pages")
	if f == nil {
		t.Error("expected --pages flag")
	}
}

func TestRenderCommandRequiresArg(t *testing.T) {
	err := renderCmd.Args(renderCmd, []string{})
	if err == nil {
		t.Error("expected error for missing argument")
	}
}

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
