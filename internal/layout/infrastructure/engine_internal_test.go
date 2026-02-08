package infrastructure

import (
	"testing"

	shared "github.com/vpedrosa/pen2pdf/internal/shared/domain"
)

type stubMeasurer struct{}

func (m *stubMeasurer) MeasureText(text, family string, size float64, weight string, maxWidth float64) (float64, float64) {
	// Simple estimation: 8px per char width, size for height
	w := float64(len(text)) * 8
	if maxWidth > 0 && w > maxWidth {
		w = maxWidth
	}
	return w, size
}

func TestIntrinsicSizeEmptyFrame(t *testing.T) {
	frame := &shared.Frame{ID: "f1"}
	w, h := intrinsicSize(frame, nil, 800)
	if w != 0 || h != 0 {
		t.Errorf("expected (0,0) for empty frame, got (%f,%f)", w, h)
	}
}

func TestIntrinsicSizeEmptyFrameWithPadding(t *testing.T) {
	frame := &shared.Frame{
		ID:      "f1",
		Padding: shared.UniformPadding(20),
	}
	w, h := intrinsicSize(frame, nil, 800)
	if w != 40 || h != 40 {
		t.Errorf("expected (40,40) for padded empty frame, got (%f,%f)", w, h)
	}
}

func TestIntrinsicSizeHorizontalFixedChildren(t *testing.T) {
	frame := &shared.Frame{
		ID:     "f1",
		Layout: "horizontal",
		Children: []shared.Node{
			&shared.Frame{ID: "a", Width: shared.FixedDimension(100), Height: shared.FixedDimension(50)},
			&shared.Frame{ID: "b", Width: shared.FixedDimension(200), Height: shared.FixedDimension(80)},
		},
	}
	w, h := intrinsicSize(frame, nil, 800)
	// horizontal: totalMain = 100+200 = 300, maxCross = 80
	if w != 300 {
		t.Errorf("expected width 300, got %f", w)
	}
	if h != 80 {
		t.Errorf("expected height 80, got %f", h)
	}
}

func TestIntrinsicSizeVerticalFixedChildren(t *testing.T) {
	frame := &shared.Frame{
		ID:     "f1",
		Layout: "vertical",
		Children: []shared.Node{
			&shared.Frame{ID: "a", Width: shared.FixedDimension(100), Height: shared.FixedDimension(50)},
			&shared.Frame{ID: "b", Width: shared.FixedDimension(200), Height: shared.FixedDimension(80)},
		},
	}
	w, h := intrinsicSize(frame, nil, 800)
	// vertical: totalMain = 50+80 = 130, maxCross = 200
	if w != 200 {
		t.Errorf("expected width 200, got %f", w)
	}
	if h != 130 {
		t.Errorf("expected height 130, got %f", h)
	}
}

func TestIntrinsicSizeWithGap(t *testing.T) {
	frame := &shared.Frame{
		ID:     "f1",
		Layout: "vertical",
		Gap:    10,
		Children: []shared.Node{
			&shared.Frame{ID: "a", Width: shared.FixedDimension(100), Height: shared.FixedDimension(50)},
			&shared.Frame{ID: "b", Width: shared.FixedDimension(100), Height: shared.FixedDimension(50)},
			&shared.Frame{ID: "c", Width: shared.FixedDimension(100), Height: shared.FixedDimension(50)},
		},
	}
	w, h := intrinsicSize(frame, nil, 800)
	// vertical: totalMain = 150, gaps = 2*10 = 20
	if w != 100 {
		t.Errorf("expected width 100, got %f", w)
	}
	if h != 170 { // 150 + 20
		t.Errorf("expected height 170, got %f", h)
	}
}

func TestIntrinsicSizeWithPaddingAndChildren(t *testing.T) {
	frame := &shared.Frame{
		ID:      "f1",
		Layout:  "horizontal",
		Padding: shared.Padding{Top: 10, Right: 20, Bottom: 10, Left: 20},
		Children: []shared.Node{
			&shared.Frame{ID: "a", Width: shared.FixedDimension(100), Height: shared.FixedDimension(50)},
		},
	}
	w, h := intrinsicSize(frame, nil, 800)
	// horizontal: totalMain=100 + padH=40, maxCross=50 + padV=20
	if w != 140 {
		t.Errorf("expected width 140, got %f", w)
	}
	if h != 70 {
		t.Errorf("expected height 70, got %f", h)
	}
}

func TestIntrinsicSizeNestedFrames(t *testing.T) {
	frame := &shared.Frame{
		ID:     "outer",
		Layout: "vertical",
		Children: []shared.Node{
			&shared.Frame{
				ID:     "inner",
				Layout: "horizontal",
				Children: []shared.Node{
					&shared.Frame{ID: "a", Width: shared.FixedDimension(60), Height: shared.FixedDimension(30)},
					&shared.Frame{ID: "b", Width: shared.FixedDimension(40), Height: shared.FixedDimension(30)},
				},
			},
		},
	}
	w, h := intrinsicSize(frame, nil, 800)
	// inner horizontal: w=100, h=30; outer vertical: w=100, h=30
	if w != 100 {
		t.Errorf("expected width 100, got %f", w)
	}
	if h != 30 {
		t.Errorf("expected height 30, got %f", h)
	}
}

func TestIntrinsicSizeWithTextNode(t *testing.T) {
	frame := &shared.Frame{
		ID:     "f1",
		Layout: "vertical",
		Children: []shared.Node{
			&shared.Text{ID: "t1", Content: "Hello", FontSize: 16},
		},
	}
	measurer := &stubMeasurer{}
	w, h := intrinsicSize(frame, measurer, 800)
	// "Hello" = 5 chars * 8 = 40 width, 16 height
	if w != 40 {
		t.Errorf("expected width 40, got %f", w)
	}
	if h != 16 {
		t.Errorf("expected height 16, got %f", h)
	}
}

func TestIntrinsicSizeTextWithFixedWidth(t *testing.T) {
	frame := &shared.Frame{
		ID:     "f1",
		Layout: "vertical",
		Children: []shared.Node{
			&shared.Text{ID: "t1", Content: "Hello", FontSize: 16, Width: shared.FixedDimension(200)},
		},
	}
	measurer := &stubMeasurer{}
	w, h := intrinsicSize(frame, measurer, 800)
	// Text has fixed width 200, measurer uses that as maxWidth
	if w != 200 {
		t.Errorf("expected width 200, got %f", w)
	}
	if h != 16 {
		t.Errorf("expected height 16, got %f", h)
	}
}

func TestCrossOffsetCenter(t *testing.T) {
	offset := crossOffset("center", 800, 200)
	if offset != 300 { // (800-200)/2
		t.Errorf("expected 300, got %f", offset)
	}
}

func TestCrossOffsetEnd(t *testing.T) {
	offset := crossOffset("end", 800, 200)
	if offset != 600 { // 800-200
		t.Errorf("expected 600, got %f", offset)
	}
}

func TestCrossOffsetStart(t *testing.T) {
	offset := crossOffset("start", 800, 200)
	if offset != 0 {
		t.Errorf("expected 0, got %f", offset)
	}
}

func TestCrossOffsetDefault(t *testing.T) {
	offset := crossOffset("", 800, 200)
	if offset != 0 {
		t.Errorf("expected 0 for empty alignItems, got %f", offset)
	}
}
