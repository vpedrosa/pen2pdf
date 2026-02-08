package application

import (
	"fmt"
	"io"
	"sort"

	parser "github.com/vpedrosa/pen2pdf/internal/parser/domain"
	shared "github.com/vpedrosa/pen2pdf/internal/shared/domain"
)

// InfoService extracts metadata from .pen files.
type InfoService struct {
	parser parser.Parser
}

func NewInfoService(p parser.Parser) *InfoService {
	return &InfoService{parser: p}
}

// PageInfo describes a single page in the document.
type PageInfo struct {
	Name   string
	Width  float64
	Height float64
}

// DocumentInfo contains document metadata.
type DocumentInfo struct {
	Version   string
	Pages     []PageInfo
	Variables map[string]shared.Variable
	Fonts     []string
}

// GetInfo parses a .pen file and returns its metadata.
func (s *InfoService) GetInfo(input io.Reader) (*DocumentInfo, error) {
	doc, err := s.parser.Parse(input)
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	info := &DocumentInfo{
		Version:   doc.Version,
		Variables: doc.Variables,
	}

	for _, child := range doc.Children {
		if frame, ok := child.(*shared.Frame); ok {
			info.Pages = append(info.Pages, PageInfo{
				Name:   frame.Name,
				Width:  frame.Width.Value,
				Height: frame.Height.Value,
			})
		}
	}

	info.Fonts = collectFonts(doc.Children)
	return info, nil
}

func collectFonts(nodes []shared.Node) []string {
	fontSet := make(map[string]bool)
	walkFonts(nodes, fontSet)

	fonts := make([]string, 0, len(fontSet))
	for f := range fontSet {
		fonts = append(fonts, f)
	}
	sort.Strings(fonts)
	return fonts
}

func walkFonts(nodes []shared.Node, fonts map[string]bool) {
	for _, node := range nodes {
		switch n := node.(type) {
		case *shared.Frame:
			walkFonts(n.Children, fonts)
		case *shared.Text:
			if n.FontFamily != "" {
				fonts[n.FontFamily] = true
			}
		}
	}
}
