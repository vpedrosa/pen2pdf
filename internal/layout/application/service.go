package application

import (
	layout "github.com/vpedrosa/pen2pdf/internal/layout/domain"
	shared "github.com/vpedrosa/pen2pdf/internal/shared/domain"
)

// LayoutService orchestrates document layout computation.
type LayoutService struct {
	engine   layout.LayoutEngine
	measurer layout.TextMeasurer
}

// NewLayoutService creates a LayoutService with the given engine and measurer ports.
func NewLayoutService(e layout.LayoutEngine, m layout.TextMeasurer) *LayoutService {
	return &LayoutService{engine: e, measurer: m}
}

// Layout computes absolute positions for all nodes in the document.
func (s *LayoutService) Layout(doc *shared.Document) ([]layout.Page, error) {
	return s.engine.Layout(doc, s.measurer)
}
