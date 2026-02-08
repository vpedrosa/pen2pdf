package domain

import (
	shared "github.com/vpedrosa/pen2pdf/internal/shared/domain"
)

// LayoutBox holds the computed absolute position and dimensions for a node.
type LayoutBox struct {
	X        float64
	Y        float64
	Width    float64
	Height   float64
	Node     shared.Node
	Children []*LayoutBox
}

// Page represents the computed layout tree for a root-level frame.
type Page struct {
	Width  float64
	Height float64
	Root   *LayoutBox
}

// TextMeasurer decouples text measurement from the layout engine,
// allowing tests to run without real font files.
type TextMeasurer interface {
	MeasureText(text string, fontFamily string, fontSize float64, fontWeight string, maxWidth float64) (width, height float64)
}

// LayoutEngine computes absolute positions for all nodes in a document.
// Each root-level frame produces a separate Page.
type LayoutEngine interface {
	Layout(doc *shared.Document, measurer TextMeasurer) ([]Page, error)
}
