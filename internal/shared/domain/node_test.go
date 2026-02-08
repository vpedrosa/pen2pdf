package domain_test

import (
	"testing"

	"github.com/vpedrosa/pen2pdf/internal/shared/domain"
)

func TestFrameImplementsNode(t *testing.T) {
	var n domain.Node = &domain.Frame{ID: "f1", Name: "test"}
	if n.GetID() != "f1" {
		t.Errorf("expected ID 'f1', got '%s'", n.GetID())
	}
	if n.GetName() != "test" {
		t.Errorf("expected Name 'test', got '%s'", n.GetName())
	}
	if n.GetType() != domain.NodeTypeFrame {
		t.Errorf("expected type '%s', got '%s'", domain.NodeTypeFrame, n.GetType())
	}
}

func TestTextImplementsNode(t *testing.T) {
	var n domain.Node = &domain.Text{ID: "t1", Name: "label"}
	if n.GetID() != "t1" {
		t.Errorf("expected ID 't1', got '%s'", n.GetID())
	}
	if n.GetName() != "label" {
		t.Errorf("expected Name 'label', got '%s'", n.GetName())
	}
	if n.GetType() != domain.NodeTypeText {
		t.Errorf("expected type '%s', got '%s'", domain.NodeTypeText, n.GetType())
	}
}

func TestFixedDimension(t *testing.T) {
	d := domain.FixedDimension(800)
	if d.Value != 800 {
		t.Errorf("expected value 800, got %f", d.Value)
	}
	if d.FillContainer {
		t.Error("expected FillContainer to be false")
	}
}

func TestFillContainerDimension(t *testing.T) {
	d := domain.FillContainerDimension()
	if !d.FillContainer {
		t.Error("expected FillContainer to be true")
	}
}

func TestUniformPadding(t *testing.T) {
	p := domain.UniformPadding(20)
	if p.Top != 20 || p.Right != 20 || p.Bottom != 20 || p.Left != 20 {
		t.Errorf("expected all sides 20, got %+v", p)
	}
}

func TestPaddingAsArray(t *testing.T) {
	p := domain.Padding{Top: 10, Right: 20, Bottom: 30, Left: 40}
	if p.Top != 10 || p.Right != 20 || p.Bottom != 30 || p.Left != 40 {
		t.Errorf("expected [10,20,30,40], got %+v", p)
	}
}

func TestFrameChildren(t *testing.T) {
	frame := &domain.Frame{
		ID:   "parent",
		Name: "container",
		Children: []domain.Node{
			&domain.Frame{ID: "child1", Name: "inner"},
			&domain.Text{ID: "child2", Name: "label", Content: "hello"},
		},
	}
	if len(frame.Children) != 2 {
		t.Fatalf("expected 2 children, got %d", len(frame.Children))
	}
	if frame.Children[0].GetType() != domain.NodeTypeFrame {
		t.Error("first child should be frame")
	}
	if frame.Children[1].GetType() != domain.NodeTypeText {
		t.Error("second child should be text")
	}
}
