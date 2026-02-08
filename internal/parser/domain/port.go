package domain

import (
	"io"

	shared "github.com/vpedrosa/pen2pdf/internal/shared/domain"
)

// Parser defines the contract for parsing design files into the domain model.
// Adapters implement this interface for specific formats (e.g., JSON).
type Parser interface {
	Parse(r io.Reader) (*shared.Document, error)
}
