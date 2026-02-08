package cmd

import (
	"testing"

	shared "github.com/vpedrosa/pen2pdf/internal/shared/domain"
)

func TestInfoCommandRegistered(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Use == "info [input.pen]" {
			found = true
			break
		}
	}
	if !found {
		t.Error("info command not registered")
	}
}

func TestInfoCommandRequiresArg(t *testing.T) {
	err := infoCmd.Args(infoCmd, []string{})
	if err == nil {
		t.Error("expected error for missing argument")
	}
}

func TestCollectFonts(t *testing.T) {
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
	// Should be sorted
	expected := []string{"Inter", "Montserrat", "Playfair Display"}
	for i, f := range fonts {
		if f != expected[i] {
			t.Errorf("expected fonts[%d] '%s', got '%s'", i, expected[i], f)
		}
	}
}

func TestCollectFontsEmpty(t *testing.T) {
	nodes := []shared.Node{
		&shared.Frame{ID: "f1", Name: "empty"},
	}

	fonts := collectFonts(nodes)
	if len(fonts) != 0 {
		t.Errorf("expected 0 fonts, got %d", len(fonts))
	}
}

func TestCollectFontsSkipsEmptyFamily(t *testing.T) {
	nodes := []shared.Node{
		&shared.Text{ID: "t1", Name: "a", FontFamily: ""},
		&shared.Text{ID: "t2", Name: "b", FontFamily: "Inter"},
	}

	fonts := collectFonts(nodes)
	if len(fonts) != 1 {
		t.Fatalf("expected 1 font, got %d", len(fonts))
	}
	if fonts[0] != "Inter" {
		t.Errorf("expected 'Inter', got '%s'", fonts[0])
	}
}
