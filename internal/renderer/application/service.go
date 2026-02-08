package application

import (
	"io"

	layout "github.com/vpedrosa/pen2pdf/internal/layout/domain"
	renderer "github.com/vpedrosa/pen2pdf/internal/renderer/domain"
)

// RenderResult contains the outcome of a render operation.
type RenderResult struct {
	PageCount int
}

// RenderService orchestrates PDF rendering.
type RenderService struct {
	renderer renderer.Renderer
}

// NewRenderService creates a RenderService with the given Renderer port.
func NewRenderService(r renderer.Renderer) *RenderService {
	return &RenderService{renderer: r}
}

// Render produces PDF output from laid-out pages.
func (s *RenderService) Render(pages []layout.Page, output io.Writer) (*RenderResult, error) {
	if err := s.renderer.Render(pages, output); err != nil {
		return nil, err
	}
	return &RenderResult{PageCount: len(pages)}, nil
}
