package application

import (
	"io"

	parser "github.com/vpedrosa/pen2pdf/internal/parser/domain"
	shared "github.com/vpedrosa/pen2pdf/internal/shared/domain"
)

// ParseService orchestrates document parsing.
type ParseService struct {
	parser parser.Parser
}

// NewParseService creates a ParseService with the given Parser port.
func NewParseService(p parser.Parser) *ParseService {
	return &ParseService{parser: p}
}

// Parse reads a .pen file and returns the parsed document.
func (s *ParseService) Parse(input io.Reader) (*shared.Document, error) {
	return s.parser.Parse(input)
}
