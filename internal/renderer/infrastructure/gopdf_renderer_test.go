package infrastructure_test

import (
	"bytes"
	"testing"

	layout "github.com/vpedrosa/pen2pdf/internal/layout/domain"
	renderer "github.com/vpedrosa/pen2pdf/internal/renderer/domain"
	"github.com/vpedrosa/pen2pdf/internal/renderer/infrastructure"
	shared "github.com/vpedrosa/pen2pdf/internal/shared/domain"
)

func TestPDFRendererImplementsPort(t *testing.T) {
	var _ renderer.Renderer = infrastructure.NewPDFRenderer(nil, nil)
}

func TestRenderEmptyDocument(t *testing.T) {
	r := infrastructure.NewPDFRenderer(nil, nil)
	var buf bytes.Buffer
	err := r.Render([]layout.Page{}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// gopdf should produce valid PDF even with no pages
	if buf.Len() == 0 {
		t.Error("expected non-empty output")
	}
}

func TestRenderSinglePage(t *testing.T) {
	r := infrastructure.NewPDFRenderer(nil, nil)
	pages := []layout.Page{
		{
			Width:  800,
			Height: 1000,
			Root: &layout.LayoutBox{
				X: 0, Y: 0, Width: 800, Height: 1000,
				Node: &shared.Frame{ID: "page", Name: "page"},
			},
		},
	}

	var buf bytes.Buffer
	err := r.Render(pages, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// PDF magic bytes
	if !bytes.HasPrefix(buf.Bytes(), []byte("%PDF")) {
		t.Error("output does not start with %PDF header")
	}
}

func TestRenderMultiplePages(t *testing.T) {
	r := infrastructure.NewPDFRenderer(nil, nil)
	pages := []layout.Page{
		{
			Width: 800, Height: 1000,
			Root: &layout.LayoutBox{
				X: 0, Y: 0, Width: 800, Height: 1000,
				Node: &shared.Frame{ID: "p1", Name: "Front"},
			},
		},
		{
			Width: 800, Height: 1000,
			Root: &layout.LayoutBox{
				X: 0, Y: 0, Width: 800, Height: 1000,
				Node: &shared.Frame{ID: "p2", Name: "Back"},
			},
		},
	}

	var buf bytes.Buffer
	err := r.Render(pages, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty output")
	}
}

func TestRenderSolidFillRect(t *testing.T) {
	r := infrastructure.NewPDFRenderer(nil, nil)
	pages := []layout.Page{
		{
			Width: 800, Height: 1000,
			Root: &layout.LayoutBox{
				X: 0, Y: 0, Width: 800, Height: 1000,
				Node: &shared.Frame{
					ID: "page", Name: "page",
					Fill: shared.SolidFill("#FF6B35"),
				},
			},
		},
	}

	var buf bytes.Buffer
	err := r.Render(pages, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty output")
	}
}

func TestRenderSolidFillWithAlpha(t *testing.T) {
	r := infrastructure.NewPDFRenderer(nil, nil)
	pages := []layout.Page{
		{
			Width: 800, Height: 1000,
			Root: &layout.LayoutBox{
				X: 0, Y: 0, Width: 800, Height: 1000,
				Node: &shared.Frame{
					ID: "page", Name: "page",
					Fill: shared.SolidFill("#000000BB"),
				},
			},
		},
	}

	var buf bytes.Buffer
	err := r.Render(pages, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRenderRoundedRect(t *testing.T) {
	r := infrastructure.NewPDFRenderer(nil, nil)
	pages := []layout.Page{
		{
			Width: 800, Height: 1000,
			Root: &layout.LayoutBox{
				X: 0, Y: 0, Width: 800, Height: 1000,
				Node: &shared.Frame{
					ID: "page", Name: "page",
					Fill:         shared.SolidFill("#FF6B35"),
					CornerRadius: 20,
				},
			},
		},
	}

	var buf bytes.Buffer
	err := r.Render(pages, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRenderNestedFrames(t *testing.T) {
	r := infrastructure.NewPDFRenderer(nil, nil)
	pages := []layout.Page{
		{
			Width: 800, Height: 1000,
			Root: &layout.LayoutBox{
				X: 0, Y: 0, Width: 800, Height: 1000,
				Node: &shared.Frame{ID: "page", Name: "page", Fill: shared.SolidFill("#FFFFFF")},
				Children: []*layout.LayoutBox{
					{
						X: 40, Y: 40, Width: 720, Height: 920,
						Node: &shared.Frame{ID: "inner", Name: "inner", Fill: shared.SolidFill("#FF0000")},
					},
				},
			},
		},
	}

	var buf bytes.Buffer
	err := r.Render(pages, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRenderNoFillFrame(t *testing.T) {
	r := infrastructure.NewPDFRenderer(nil, nil)
	pages := []layout.Page{
		{
			Width: 800, Height: 1000,
			Root: &layout.LayoutBox{
				X: 0, Y: 0, Width: 800, Height: 1000,
				Node: &shared.Frame{ID: "page", Name: "page"},
			},
		},
	}

	var buf bytes.Buffer
	err := r.Render(pages, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRenderEmptyTextSkipped(t *testing.T) {
	r := infrastructure.NewPDFRenderer(nil, nil)
	pages := []layout.Page{
		{
			Width: 800, Height: 1000,
			Root: &layout.LayoutBox{
				X: 0, Y: 0, Width: 800, Height: 1000,
				Node: &shared.Frame{ID: "page", Name: "page"},
				Children: []*layout.LayoutBox{
					{
						X: 40, Y: 40, Width: 200, Height: 30,
						Node: &shared.Text{ID: "t1", Name: "empty", Content: ""},
					},
				},
			},
		},
	}

	var buf bytes.Buffer
	err := r.Render(pages, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRenderTextWithFallbackFont(t *testing.T) {
	// No font loader — should use embedded Go font fallback
	r := infrastructure.NewPDFRenderer(nil, nil)
	pages := []layout.Page{
		{
			Width: 800, Height: 1000,
			Root: &layout.LayoutBox{
				X: 0, Y: 0, Width: 800, Height: 1000,
				Node: &shared.Frame{ID: "page", Name: "page"},
				Children: []*layout.LayoutBox{
					{
						X: 40, Y: 40, Width: 200, Height: 30,
						Node: &shared.Text{
							ID: "t1", Name: "label",
							Content:    "Hello World",
							Fill:       "#FF0000",
							FontFamily: "NonExistentFont",
							FontSize:   16,
							FontWeight: "400",
						},
					},
				},
			},
		},
	}

	var buf bytes.Buffer
	err := r.Render(pages, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.HasPrefix(buf.Bytes(), []byte("%PDF")) {
		t.Error("output does not start with %PDF header")
	}
}

func TestRenderTextBoldFallback(t *testing.T) {
	r := infrastructure.NewPDFRenderer(nil, nil)
	pages := []layout.Page{
		{
			Width: 800, Height: 1000,
			Root: &layout.LayoutBox{
				X: 0, Y: 0, Width: 800, Height: 1000,
				Node: &shared.Frame{ID: "page", Name: "page"},
				Children: []*layout.LayoutBox{
					{
						X: 40, Y: 40, Width: 200, Height: 30,
						Node: &shared.Text{
							ID: "t1", Name: "bold",
							Content:    "Bold Text",
							FontFamily: "MissingFont",
							FontSize:   24,
							FontWeight: "700",
						},
					},
				},
			},
		},
	}

	var buf bytes.Buffer
	err := r.Render(pages, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRenderTextItalicFallback(t *testing.T) {
	r := infrastructure.NewPDFRenderer(nil, nil)
	pages := []layout.Page{
		{
			Width: 800, Height: 1000,
			Root: &layout.LayoutBox{
				X: 0, Y: 0, Width: 800, Height: 1000,
				Node: &shared.Frame{ID: "page", Name: "page"},
				Children: []*layout.LayoutBox{
					{
						X: 40, Y: 40, Width: 200, Height: 30,
						Node: &shared.Text{
							ID: "t1", Name: "italic",
							Content:    "Italic Text",
							FontFamily: "MissingFont",
							FontSize:   14,
							FontWeight: "400",
							FontStyle:  "italic",
						},
					},
				},
			},
		},
	}

	var buf bytes.Buffer
	err := r.Render(pages, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRenderMultipleTextNodesFallback(t *testing.T) {
	r := infrastructure.NewPDFRenderer(nil, nil)
	pages := []layout.Page{
		{
			Width: 800, Height: 1000,
			Root: &layout.LayoutBox{
				X: 0, Y: 0, Width: 800, Height: 1000,
				Node: &shared.Frame{ID: "page", Name: "page"},
				Children: []*layout.LayoutBox{
					{
						X: 40, Y: 40, Width: 200, Height: 30,
						Node: &shared.Text{
							ID: "t1", Name: "a",
							Content: "First", FontFamily: "FontA", FontSize: 16, FontWeight: "400",
						},
					},
					{
						X: 40, Y: 80, Width: 200, Height: 30,
						Node: &shared.Text{
							ID: "t2", Name: "b",
							Content: "Second", FontFamily: "FontB", FontSize: 16, FontWeight: "700",
						},
					},
				},
			},
		},
	}

	var buf bytes.Buffer
	err := r.Render(pages, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
