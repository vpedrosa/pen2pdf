package application

import (
	"fmt"
	"io"

	parser "github.com/vpedrosa/pen2pdf/internal/parser/domain"
	resolver "github.com/vpedrosa/pen2pdf/internal/resolver/domain"
)

// ValidateService validates .pen files without rendering.
type ValidateService struct {
	parser   parser.Parser
	resolver resolver.Resolver
}

func NewValidateService(p parser.Parser, r resolver.Resolver) *ValidateService {
	return &ValidateService{parser: p, resolver: r}
}

// ValidateResult contains the outcome of a validation.
type ValidateResult struct {
	PageCount     int
	VariableCount int
}

// Validate parses and resolves a .pen file, returning metadata.
func (s *ValidateService) Validate(input io.Reader) (*ValidateResult, error) {
	doc, err := s.parser.Parse(input)
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	if err := s.resolver.Resolve(doc); err != nil {
		return nil, fmt.Errorf("resolve error: %w", err)
	}

	return &ValidateResult{
		PageCount:     len(doc.Children),
		VariableCount: len(doc.Variables),
	}, nil
}
