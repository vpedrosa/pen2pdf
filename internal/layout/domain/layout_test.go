package domain_test

import (
	"testing"

	layout "github.com/vpedrosa/pen2pdf/internal/layout/domain"
	shared "github.com/vpedrosa/pen2pdf/internal/shared/domain"
)

type stubTextMeasurer struct{}

func (s *stubTextMeasurer) MeasureText(_, _ string, _ float64, _ string, _ float64) (float64, float64) {
	return 100, 20
}

type stubLayoutEngine struct {
	pages []layout.Page
	err   error
}

func (s *stubLayoutEngine) Layout(_ *shared.Document, _ layout.TextMeasurer) ([]layout.Page, error) {
	return s.pages, s.err
}

func TestLayoutEngineInterfaceCompliance(t *testing.T) {
	var _ layout.LayoutEngine = &stubLayoutEngine{}
}

func TestTextMeasurerInterfaceCompliance(t *testing.T) {
	var _ layout.TextMeasurer = &stubTextMeasurer{}
}

func TestLayoutBoxTree(t *testing.T) {
	root := &layout.LayoutBox{
		X: 0, Y: 0, Width: 800, Height: 1000,
		Node: &shared.Frame{ID: "page", Name: "page"},
		Children: []*layout.LayoutBox{
			{
				X: 40, Y: 40, Width: 720, Height: 920,
				Node: &shared.Frame{ID: "content", Name: "content"},
				Children: []*layout.LayoutBox{
					{
						X: 40, Y: 40, Width: 720, Height: 48,
						Node: &shared.Text{ID: "title", Name: "title", Content: "Hello"},
					},
				},
			},
		},
	}

	if root.Width != 800 {
		t.Errorf("expected root width 800, got %f", root.Width)
	}
	if len(root.Children) != 1 {
		t.Fatalf("expected 1 child, got %d", len(root.Children))
	}
	content := root.Children[0]
	if content.X != 40 {
		t.Errorf("expected content X 40, got %f", content.X)
	}
	if len(content.Children) != 1 {
		t.Fatalf("expected 1 grandchild, got %d", len(content.Children))
	}
	title := content.Children[0]
	if title.Node.GetID() != "title" {
		t.Errorf("expected node ID 'title', got '%s'", title.Node.GetID())
	}
}

func TestPageStructure(t *testing.T) {
	page := layout.Page{
		Width:  800,
		Height: 1000,
		Root: &layout.LayoutBox{
			X: 0, Y: 0, Width: 800, Height: 1000,
			Node: &shared.Frame{ID: "p1", Name: "Front"},
		},
	}

	if page.Width != 800 {
		t.Errorf("expected width 800, got %f", page.Width)
	}
	if page.Height != 1000 {
		t.Errorf("expected height 1000, got %f", page.Height)
	}
	if page.Root.Node.GetName() != "Front" {
		t.Errorf("expected root node name 'Front', got '%s'", page.Root.Node.GetName())
	}
}

func TestStubLayoutEngineReturnsPages(t *testing.T) {
	engine := &stubLayoutEngine{
		pages: []layout.Page{
			{Width: 800, Height: 1000},
			{Width: 800, Height: 1000},
		},
	}

	pages, err := engine.Layout(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pages) != 2 {
		t.Errorf("expected 2 pages, got %d", len(pages))
	}
}
