package application

import (
	asset "github.com/vpedrosa/pen2pdf/internal/asset/domain"
	shared "github.com/vpedrosa/pen2pdf/internal/shared/domain"
)

// FontService handles font-related use cases.
type FontService struct {
	fontLoader asset.FontLoader
}

// NewFontService creates a FontService with the given FontLoader port.
func NewFontService(fl asset.FontLoader) *FontService {
	return &FontService{fontLoader: fl}
}

// DetectMissingFonts checks each font reference in the document against the
// configured font loader and returns those that cannot be loaded.
func (s *FontService) DetectMissingFonts(doc *shared.Document) []shared.FontRef {
	refs := shared.CollectFontRefs(doc)
	var missing []shared.FontRef
	for _, ref := range refs {
		_, err := s.fontLoader.LoadFont(ref.Family, ref.Weight, ref.Style)
		if err != nil {
			missing = append(missing, ref)
		}
	}
	return missing
}
