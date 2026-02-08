package infrastructure_test

import (
	"testing"

	layout "github.com/vpedrosa/pen2pdf/internal/layout/domain"
	"github.com/vpedrosa/pen2pdf/internal/layout/infrastructure"
	shared "github.com/vpedrosa/pen2pdf/internal/shared/domain"
)

type fixedMeasurer struct {
	width  float64
	height float64
}

func (m *fixedMeasurer) MeasureText(_, _ string, _ float64, _ string, _ float64) (float64, float64) {
	return m.width, m.height
}

func TestFlexboxEngineImplementsPort(t *testing.T) {
	var _ layout.LayoutEngine = infrastructure.NewFlexboxEngine()
}

// --- Issue #19: Fixed dimensions and padding ---

func TestLayoutFixedDimensions(t *testing.T) {
	doc := &shared.Document{
		Children: []shared.Node{
			&shared.Frame{
				ID: "page", Name: "page",
				Width: shared.FixedDimension(800), Height: shared.FixedDimension(1000),
			},
		},
	}

	pages := mustLayout(t, doc)
	if len(pages) != 1 {
		t.Fatalf("expected 1 page, got %d", len(pages))
	}

	page := pages[0]
	if page.Width != 800 || page.Height != 1000 {
		t.Errorf("expected page 800x1000, got %fx%f", page.Width, page.Height)
	}
	if page.Root.X != 0 || page.Root.Y != 0 {
		t.Errorf("expected root at (0,0), got (%f,%f)", page.Root.X, page.Root.Y)
	}
	if page.Root.Width != 800 || page.Root.Height != 1000 {
		t.Errorf("expected root 800x1000, got %fx%f", page.Root.Width, page.Root.Height)
	}
}

func TestLayoutWithPosition(t *testing.T) {
	// Root frames always start at (0,0) — canvas X/Y is ignored for PDF pages
	doc := &shared.Document{
		Children: []shared.Node{
			&shared.Frame{
				ID: "page", Name: "page", X: 96, Y: 100,
				Width: shared.FixedDimension(800), Height: shared.FixedDimension(1000),
			},
		},
	}

	pages := mustLayout(t, doc)
	if pages[0].Root.X != 0 || pages[0].Root.Y != 0 {
		t.Errorf("expected root at (0,0), got (%f,%f)", pages[0].Root.X, pages[0].Root.Y)
	}
}

func TestLayoutUniformPadding(t *testing.T) {
	doc := &shared.Document{
		Children: []shared.Node{
			&shared.Frame{
				ID: "page", Name: "page",
				Width: shared.FixedDimension(800), Height: shared.FixedDimension(1000),
				Padding: shared.UniformPadding(40),
				Children: []shared.Node{
					&shared.Frame{
						ID: "child", Name: "child",
						Width: shared.FixedDimension(100), Height: shared.FixedDimension(50),
					},
				},
			},
		},
	}

	pages := mustLayout(t, doc)
	child := pages[0].Root.Children[0]

	if child.X != 40 {
		t.Errorf("expected child X 40, got %f", child.X)
	}
	if child.Y != 40 {
		t.Errorf("expected child Y 40, got %f", child.Y)
	}
}

func TestLayoutAsymmetricPadding(t *testing.T) {
	doc := &shared.Document{
		Children: []shared.Node{
			&shared.Frame{
				ID: "page", Name: "page",
				Width: shared.FixedDimension(800), Height: shared.FixedDimension(1000),
				Padding: shared.Padding{Top: 10, Right: 20, Bottom: 30, Left: 40},
				Children: []shared.Node{
					&shared.Frame{
						ID: "child", Name: "child",
						Width: shared.FixedDimension(100), Height: shared.FixedDimension(50),
					},
				},
			},
		},
	}

	pages := mustLayout(t, doc)
	child := pages[0].Root.Children[0]

	if child.X != 40 {
		t.Errorf("expected child X 40 (left padding), got %f", child.X)
	}
	if child.Y != 10 {
		t.Errorf("expected child Y 10 (top padding), got %f", child.Y)
	}
}

func TestLayoutZeroPadding(t *testing.T) {
	doc := &shared.Document{
		Children: []shared.Node{
			&shared.Frame{
				ID: "page", Name: "page",
				Width: shared.FixedDimension(800), Height: shared.FixedDimension(1000),
				Children: []shared.Node{
					&shared.Frame{
						ID: "child", Name: "child",
						Width: shared.FixedDimension(100), Height: shared.FixedDimension(50),
					},
				},
			},
		},
	}

	pages := mustLayout(t, doc)
	child := pages[0].Root.Children[0]

	if child.X != 0 || child.Y != 0 {
		t.Errorf("expected child at (0,0) with zero padding, got (%f,%f)", child.X, child.Y)
	}
}

func TestLayoutNestedFixed(t *testing.T) {
	doc := &shared.Document{
		Children: []shared.Node{
			&shared.Frame{
				ID: "page", Name: "page",
				Width: shared.FixedDimension(800), Height: shared.FixedDimension(1000),
				Padding: shared.UniformPadding(20),
				Children: []shared.Node{
					&shared.Frame{
						ID: "inner", Name: "inner",
						Width: shared.FixedDimension(400), Height: shared.FixedDimension(300),
						Padding: shared.UniformPadding(10),
						Children: []shared.Node{
							&shared.Frame{
								ID: "deep", Name: "deep",
								Width: shared.FixedDimension(100), Height: shared.FixedDimension(50),
							},
						},
					},
				},
			},
		},
	}

	pages := mustLayout(t, doc)
	inner := pages[0].Root.Children[0]
	if inner.X != 20 || inner.Y != 20 {
		t.Errorf("expected inner at (20,20), got (%f,%f)", inner.X, inner.Y)
	}

	deep := inner.Children[0]
	// deep.X = inner.X + inner.Padding.Left = 20 + 10 = 30
	if deep.X != 30 || deep.Y != 30 {
		t.Errorf("expected deep at (30,30), got (%f,%f)", deep.X, deep.Y)
	}
}

func TestLayoutMultiplePages(t *testing.T) {
	doc := &shared.Document{
		Children: []shared.Node{
			&shared.Frame{
				ID: "p1", Name: "Front",
				Width: shared.FixedDimension(800), Height: shared.FixedDimension(1000),
			},
			&shared.Frame{
				ID: "p2", Name: "Back",
				Width: shared.FixedDimension(800), Height: shared.FixedDimension(1000),
			},
		},
	}

	pages := mustLayout(t, doc)
	if len(pages) != 2 {
		t.Fatalf("expected 2 pages, got %d", len(pages))
	}
	if pages[0].Root.Node.GetName() != "Front" {
		t.Errorf("expected first page 'Front', got '%s'", pages[0].Root.Node.GetName())
	}
	if pages[1].Root.Node.GetName() != "Back" {
		t.Errorf("expected second page 'Back', got '%s'", pages[1].Root.Node.GetName())
	}
}

// --- Issue #20: fill_container and flexbox distribution ---

func TestLayoutVerticalStacking(t *testing.T) {
	doc := &shared.Document{
		Children: []shared.Node{
			&shared.Frame{
				ID: "page", Name: "page",
				Width: shared.FixedDimension(800), Height: shared.FixedDimension(1000),
				Layout: "vertical",
				Children: []shared.Node{
					&shared.Frame{ID: "a", Name: "a", Width: shared.FixedDimension(200), Height: shared.FixedDimension(100)},
					&shared.Frame{ID: "b", Name: "b", Width: shared.FixedDimension(200), Height: shared.FixedDimension(150)},
				},
			},
		},
	}

	pages := mustLayout(t, doc)
	a := pages[0].Root.Children[0]
	b := pages[0].Root.Children[1]

	if a.Y != 0 {
		t.Errorf("expected a.Y 0, got %f", a.Y)
	}
	if b.Y != 100 {
		t.Errorf("expected b.Y 100, got %f", b.Y)
	}
}

func TestLayoutHorizontalStacking(t *testing.T) {
	doc := &shared.Document{
		Children: []shared.Node{
			&shared.Frame{
				ID: "page", Name: "page",
				Width: shared.FixedDimension(800), Height: shared.FixedDimension(1000),
				Layout: "horizontal",
				Children: []shared.Node{
					&shared.Frame{ID: "a", Name: "a", Width: shared.FixedDimension(200), Height: shared.FixedDimension(100)},
					&shared.Frame{ID: "b", Name: "b", Width: shared.FixedDimension(300), Height: shared.FixedDimension(100)},
				},
			},
		},
	}

	pages := mustLayout(t, doc)
	a := pages[0].Root.Children[0]
	b := pages[0].Root.Children[1]

	if a.X != 0 {
		t.Errorf("expected a.X 0, got %f", a.X)
	}
	if b.X != 200 {
		t.Errorf("expected b.X 200, got %f", b.X)
	}
}

func TestLayoutVerticalGap(t *testing.T) {
	doc := &shared.Document{
		Children: []shared.Node{
			&shared.Frame{
				ID: "page", Name: "page",
				Width: shared.FixedDimension(800), Height: shared.FixedDimension(1000),
				Layout: "vertical", Gap: 20,
				Children: []shared.Node{
					&shared.Frame{ID: "a", Name: "a", Width: shared.FixedDimension(200), Height: shared.FixedDimension(100)},
					&shared.Frame{ID: "b", Name: "b", Width: shared.FixedDimension(200), Height: shared.FixedDimension(100)},
					&shared.Frame{ID: "c", Name: "c", Width: shared.FixedDimension(200), Height: shared.FixedDimension(100)},
				},
			},
		},
	}

	pages := mustLayout(t, doc)
	a := pages[0].Root.Children[0]
	b := pages[0].Root.Children[1]
	c := pages[0].Root.Children[2]

	if a.Y != 0 {
		t.Errorf("expected a.Y 0, got %f", a.Y)
	}
	if b.Y != 120 { // 100 + 20 gap
		t.Errorf("expected b.Y 120, got %f", b.Y)
	}
	if c.Y != 240 { // 120 + 100 + 20 gap
		t.Errorf("expected c.Y 240, got %f", c.Y)
	}
}

func TestLayoutFillContainerVertical(t *testing.T) {
	doc := &shared.Document{
		Children: []shared.Node{
			&shared.Frame{
				ID: "page", Name: "page",
				Width: shared.FixedDimension(800), Height: shared.FixedDimension(1000),
				Layout: "vertical",
				Children: []shared.Node{
					&shared.Frame{ID: "fixed", Name: "fixed", Width: shared.FixedDimension(800), Height: shared.FixedDimension(100)},
					&shared.Frame{ID: "fill", Name: "fill", Width: shared.FillContainerDimension(), Height: shared.FillContainerDimension()},
				},
			},
		},
	}

	pages := mustLayout(t, doc)
	fill := pages[0].Root.Children[1]

	if fill.Width != 800 {
		t.Errorf("expected fill width 800, got %f", fill.Width)
	}
	if fill.Height != 900 { // 1000 - 100 fixed
		t.Errorf("expected fill height 900, got %f", fill.Height)
	}
}

func TestLayoutFillContainerMultiple(t *testing.T) {
	doc := &shared.Document{
		Children: []shared.Node{
			&shared.Frame{
				ID: "page", Name: "page",
				Width: shared.FixedDimension(800), Height: shared.FixedDimension(1000),
				Layout: "vertical",
				Children: []shared.Node{
					&shared.Frame{ID: "fill1", Name: "fill1", Width: shared.FillContainerDimension(), Height: shared.FillContainerDimension()},
					&shared.Frame{ID: "fill2", Name: "fill2", Width: shared.FillContainerDimension(), Height: shared.FillContainerDimension()},
				},
			},
		},
	}

	pages := mustLayout(t, doc)
	fill1 := pages[0].Root.Children[0]
	fill2 := pages[0].Root.Children[1]

	if fill1.Height != 500 {
		t.Errorf("expected fill1 height 500, got %f", fill1.Height)
	}
	if fill2.Height != 500 {
		t.Errorf("expected fill2 height 500, got %f", fill2.Height)
	}
}

func TestLayoutJustifyContentCenter(t *testing.T) {
	doc := &shared.Document{
		Children: []shared.Node{
			&shared.Frame{
				ID: "page", Name: "page",
				Width: shared.FixedDimension(800), Height: shared.FixedDimension(1000),
				Layout: "vertical", JustifyContent: "center",
				Children: []shared.Node{
					&shared.Frame{ID: "a", Name: "a", Width: shared.FixedDimension(200), Height: shared.FixedDimension(100)},
				},
			},
		},
	}

	pages := mustLayout(t, doc)
	a := pages[0].Root.Children[0]

	// center: (1000 - 100) / 2 = 450
	if a.Y != 450 {
		t.Errorf("expected a.Y 450, got %f", a.Y)
	}
}

func TestLayoutJustifyContentEnd(t *testing.T) {
	doc := &shared.Document{
		Children: []shared.Node{
			&shared.Frame{
				ID: "page", Name: "page",
				Width: shared.FixedDimension(800), Height: shared.FixedDimension(1000),
				Layout: "vertical", JustifyContent: "end",
				Children: []shared.Node{
					&shared.Frame{ID: "a", Name: "a", Width: shared.FixedDimension(200), Height: shared.FixedDimension(100)},
				},
			},
		},
	}

	pages := mustLayout(t, doc)
	a := pages[0].Root.Children[0]

	if a.Y != 900 {
		t.Errorf("expected a.Y 900, got %f", a.Y)
	}
}

func TestLayoutJustifyContentSpaceBetween(t *testing.T) {
	doc := &shared.Document{
		Children: []shared.Node{
			&shared.Frame{
				ID: "page", Name: "page",
				Width: shared.FixedDimension(800), Height: shared.FixedDimension(1000),
				Layout: "vertical", JustifyContent: "space-between",
				Children: []shared.Node{
					&shared.Frame{ID: "a", Name: "a", Width: shared.FixedDimension(200), Height: shared.FixedDimension(100)},
					&shared.Frame{ID: "b", Name: "b", Width: shared.FixedDimension(200), Height: shared.FixedDimension(100)},
				},
			},
		},
	}

	pages := mustLayout(t, doc)
	a := pages[0].Root.Children[0]
	b := pages[0].Root.Children[1]

	if a.Y != 0 {
		t.Errorf("expected a.Y 0, got %f", a.Y)
	}
	// space-between: (1000 - 200) / 1 = 800 between items
	if b.Y != 900 { // 0 + 100 + 800
		t.Errorf("expected b.Y 900, got %f", b.Y)
	}
}

func TestLayoutAlignItemsCenter(t *testing.T) {
	doc := &shared.Document{
		Children: []shared.Node{
			&shared.Frame{
				ID: "page", Name: "page",
				Width: shared.FixedDimension(800), Height: shared.FixedDimension(1000),
				Layout: "vertical", AlignItems: "center",
				Children: []shared.Node{
					&shared.Frame{ID: "a", Name: "a", Width: shared.FixedDimension(200), Height: shared.FixedDimension(100)},
				},
			},
		},
	}

	pages := mustLayout(t, doc)
	a := pages[0].Root.Children[0]

	// center cross axis: (800 - 200) / 2 = 300
	if a.X != 300 {
		t.Errorf("expected a.X 300, got %f", a.X)
	}
}

func TestLayoutAlignItemsEnd(t *testing.T) {
	doc := &shared.Document{
		Children: []shared.Node{
			&shared.Frame{
				ID: "page", Name: "page",
				Width: shared.FixedDimension(800), Height: shared.FixedDimension(1000),
				Layout: "vertical", AlignItems: "end",
				Children: []shared.Node{
					&shared.Frame{ID: "a", Name: "a", Width: shared.FixedDimension(200), Height: shared.FixedDimension(100)},
				},
			},
		},
	}

	pages := mustLayout(t, doc)
	a := pages[0].Root.Children[0]

	if a.X != 600 { // 800 - 200
		t.Errorf("expected a.X 600, got %f", a.X)
	}
}

func TestLayoutCombinedPaddingGapAlignCenter(t *testing.T) {
	doc := &shared.Document{
		Children: []shared.Node{
			&shared.Frame{
				ID: "page", Name: "page",
				Width: shared.FixedDimension(800), Height: shared.FixedDimension(1000),
				Layout: "vertical", Gap: 20, AlignItems: "center",
				Padding: shared.UniformPadding(40),
				Children: []shared.Node{
					&shared.Frame{ID: "a", Name: "a", Width: shared.FixedDimension(200), Height: shared.FixedDimension(100)},
					&shared.Frame{ID: "b", Name: "b", Width: shared.FixedDimension(300), Height: shared.FixedDimension(100)},
				},
			},
		},
	}

	pages := mustLayout(t, doc)
	a := pages[0].Root.Children[0]
	b := pages[0].Root.Children[1]

	// Content area: 720 x 920 starting at (40, 40)
	// a centered: 40 + (720 - 200)/2 = 40 + 260 = 300
	if a.X != 300 {
		t.Errorf("expected a.X 300, got %f", a.X)
	}
	if a.Y != 40 {
		t.Errorf("expected a.Y 40, got %f", a.Y)
	}

	// b centered: 40 + (720 - 300)/2 = 40 + 210 = 250
	if b.X != 250 {
		t.Errorf("expected b.X 250, got %f", b.X)
	}
	if b.Y != 160 { // 40 + 100 + 20 gap
		t.Errorf("expected b.Y 160, got %f", b.Y)
	}
}

func TestLayoutTextNodeWithMeasurer(t *testing.T) {
	doc := &shared.Document{
		Children: []shared.Node{
			&shared.Frame{
				ID: "page", Name: "page",
				Width: shared.FixedDimension(800), Height: shared.FixedDimension(1000),
				Layout: "vertical",
				Children: []shared.Node{
					&shared.Text{ID: "t1", Name: "title", Content: "Hello", FontSize: 48},
				},
			},
		},
	}

	measurer := &fixedMeasurer{width: 200, height: 60}
	engine := infrastructure.NewFlexboxEngine()
	pages, err := engine.Layout(doc, measurer)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	text := pages[0].Root.Children[0]
	if text.Width != 200 {
		t.Errorf("expected text width 200, got %f", text.Width)
	}
	if text.Height != 60 {
		t.Errorf("expected text height 60, got %f", text.Height)
	}
}

func TestLayoutFillContainerWithGap(t *testing.T) {
	doc := &shared.Document{
		Children: []shared.Node{
			&shared.Frame{
				ID: "page", Name: "page",
				Width: shared.FixedDimension(800), Height: shared.FixedDimension(1000),
				Layout: "vertical", Gap: 20,
				Children: []shared.Node{
					&shared.Frame{ID: "fixed", Name: "fixed", Width: shared.FixedDimension(800), Height: shared.FixedDimension(100)},
					&shared.Frame{ID: "fill", Name: "fill", Width: shared.FillContainerDimension(), Height: shared.FillContainerDimension()},
				},
			},
		},
	}

	pages := mustLayout(t, doc)
	fill := pages[0].Root.Children[1]

	// remaining = 1000 - 100 (fixed) - 20 (gap) = 880
	if fill.Height != 880 {
		t.Errorf("expected fill height 880, got %f", fill.Height)
	}
}

func mustLayout(t *testing.T, doc *shared.Document) []layout.Page {
	t.Helper()
	engine := infrastructure.NewFlexboxEngine()
	pages, err := engine.Layout(doc, nil)
	if err != nil {
		t.Fatalf("unexpected layout error: %v", err)
	}
	return pages
}
