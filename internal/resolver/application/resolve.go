package application

import (
	resolver "github.com/vpedrosa/pen2pdf/internal/resolver/domain"
	shared "github.com/vpedrosa/pen2pdf/internal/shared/domain"
)

// ResolveService orchestrates variable resolution.
type ResolveService struct {
	resolver resolver.Resolver
}

// NewResolveService creates a ResolveService with the given Resolver port.
func NewResolveService(r resolver.Resolver) *ResolveService {
	return &ResolveService{resolver: r}
}

// Resolve replaces variable references in the document with their values.
func (s *ResolveService) Resolve(doc *shared.Document) error {
	return s.resolver.Resolve(doc)
}
