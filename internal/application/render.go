package application

import (
	"fmt"
	"io"
	"strings"

	asset "github.com/vpedrosa/pen2pdf/internal/asset/domain"
	layout "github.com/vpedrosa/pen2pdf/internal/layout/domain"
	parser "github.com/vpedrosa/pen2pdf/internal/parser/domain"
	renderer "github.com/vpedrosa/pen2pdf/internal/renderer/domain"
	resolver "github.com/vpedrosa/pen2pdf/internal/resolver/domain"
	shared "github.com/vpedrosa/pen2pdf/internal/shared/domain"
)

// RenderService orchestrates the full .pen → PDF pipeline.
type RenderService struct {
	parser       parser.Parser
	resolver     resolver.Resolver
	fontLoader   asset.FontLoader
	imageLoader  asset.ImageLoader
	layoutEngine layout.LayoutEngine
	measurer     layout.TextMeasurer
	renderer     renderer.Renderer
}

func NewRenderService(
	p parser.Parser,
	r resolver.Resolver,
	fl asset.FontLoader,
	il asset.ImageLoader,
	le layout.LayoutEngine,
	m layout.TextMeasurer,
	rend renderer.Renderer,
) *RenderService {
	return &RenderService{
		parser:       p,
		resolver:     r,
		fontLoader:   fl,
		imageLoader:  il,
		layoutEngine: le,
		measurer:     m,
		renderer:     rend,
	}
}

// RenderInput holds the parameters for a render operation.
type RenderInput struct {
	Input  io.Reader
	Output io.Writer
	Pages  string // comma-separated page names, empty for all
}

// RenderResult contains the outcome of a render operation.
type RenderResult struct {
	PageCount int
}

// Render parses, resolves, lays out, and renders a .pen file to PDF.
func (s *RenderService) Render(input RenderInput) (*RenderResult, error) {
	doc, err := s.parser.Parse(input.Input)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}

	if err := s.resolver.Resolve(doc); err != nil {
		return nil, fmt.Errorf("resolve: %w", err)
	}

	if input.Pages != "" {
		doc.Children, err = filterPages(doc.Children, input.Pages)
		if err != nil {
			return nil, err
		}
	}

	pages, err := s.layoutEngine.Layout(doc, s.measurer)
	if err != nil {
		return nil, fmt.Errorf("layout: %w", err)
	}

	if err := s.renderer.Render(pages, input.Output); err != nil {
		return nil, fmt.Errorf("render: %w", err)
	}

	return &RenderResult{PageCount: len(pages)}, nil
}

// DetectMissingFonts parses a document and returns any font references
// that cannot be loaded by the configured font loader.
func (s *RenderService) DetectMissingFonts(input io.Reader) ([]shared.FontRef, *shared.Document, error) {
	doc, err := s.parser.Parse(input)
	if err != nil {
		return nil, nil, fmt.Errorf("parse: %w", err)
	}

	if err := s.resolver.Resolve(doc); err != nil {
		return nil, nil, fmt.Errorf("resolve: %w", err)
	}

	refs := shared.CollectFontRefs(doc)
	var missing []shared.FontRef
	for _, ref := range refs {
		_, err := s.fontLoader.LoadFont(ref.Family, ref.Weight, ref.Style)
		if err != nil {
			missing = append(missing, ref)
		}
	}

	return missing, doc, nil
}

// RenderDocument renders an already-parsed and resolved document.
func (s *RenderService) RenderDocument(doc *shared.Document, output io.Writer, pages string) (*RenderResult, error) {
	var err error
	if pages != "" {
		doc.Children, err = filterPages(doc.Children, pages)
		if err != nil {
			return nil, err
		}
	}

	laidOut, err := s.layoutEngine.Layout(doc, s.measurer)
	if err != nil {
		return nil, fmt.Errorf("layout: %w", err)
	}

	if err := s.renderer.Render(laidOut, output); err != nil {
		return nil, fmt.Errorf("render: %w", err)
	}

	return &RenderResult{PageCount: len(laidOut)}, nil
}

func filterPages(children []shared.Node, names string) ([]shared.Node, error) {
	nameList := strings.Split(names, ",")
	nameSet := make(map[string]bool, len(nameList))
	for _, n := range nameList {
		nameSet[strings.TrimSpace(n)] = true
	}

	var filtered []shared.Node
	for _, child := range children {
		if nameSet[child.GetName()] {
			filtered = append(filtered, child)
		}
	}

	if len(filtered) == 0 {
		return nil, fmt.Errorf("no pages match --pages %q", names)
	}
	return filtered, nil
}
