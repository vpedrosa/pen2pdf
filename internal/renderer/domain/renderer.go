package domain

import (
	"io"

	layout "github.com/vpedrosa/pen2pdf/internal/layout/domain"
)

// Renderer takes fully positioned layout trees and produces output.
// Each Page corresponds to a root-level frame in the document.
type Renderer interface {
	Render(pages []layout.Page, output io.Writer) error
}
