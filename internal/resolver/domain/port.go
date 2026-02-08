package domain

import (
	shared "github.com/vpedrosa/pen2pdf/internal/shared/domain"
)

// Resolver replaces $variable references in a document with their concrete
// values from the document's variables map. It mutates the document in-place.
type Resolver interface {
	Resolve(doc *shared.Document) error
}
