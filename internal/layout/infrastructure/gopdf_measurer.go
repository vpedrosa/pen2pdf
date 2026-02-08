package infrastructure

import (
	"strings"

	"github.com/signintech/gopdf"
	asset "github.com/vpedrosa/pen2pdf/internal/asset/domain"
)

// GopdfTextMeasurer uses gopdf to measure text dimensions based on loaded fonts.
type GopdfTextMeasurer struct {
	pdf        *gopdf.GoPdf
	fontLoader asset.FontLoader
	loaded     map[string]bool
}

func NewGopdfTextMeasurer(fontLoader asset.FontLoader) *GopdfTextMeasurer {
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	return &GopdfTextMeasurer{
		pdf:        pdf,
		fontLoader: fontLoader,
		loaded:     make(map[string]bool),
	}
}

func (m *GopdfTextMeasurer) MeasureText(text, fontFamily string, fontSize float64, fontWeight string, maxWidth float64) (width, height float64) {
	fontKey := fontFamily + "-" + fontWeight
	if !m.loaded[fontKey] {
		fontData, err := m.fontLoader.LoadFont(fontFamily, fontWeight, "normal")
		if err != nil {
			// Fallback: estimate based on fontSize
			return estimateText(text, fontSize, maxWidth)
		}
		if err := m.pdf.AddTTFFontByReader(fontKey, strings.NewReader(string(fontData.Data))); err != nil {
			return estimateText(text, fontSize, maxWidth)
		}
		m.loaded[fontKey] = true
	}

	if err := m.pdf.SetFont(fontKey, "", fontSize); err != nil {
		return estimateText(text, fontSize, maxWidth)
	}

	lines := strings.Split(text, "\n")
	maxLineWidth := 0.0

	for _, line := range lines {
		lineWidth, _ := m.pdf.MeasureTextWidth(line)

		if maxWidth > 0 && lineWidth > maxWidth {
			// Wrap text: estimate wrapped lines
			wrappedLines := int(lineWidth/maxWidth) + 1
			height += float64(wrappedLines) * fontSize * 1.2
			if maxWidth > maxLineWidth {
				maxLineWidth = maxWidth
			}
		} else {
			height += fontSize * 1.2
			if lineWidth > maxLineWidth {
				maxLineWidth = lineWidth
			}
		}
	}

	return maxLineWidth, height
}

// estimateText provides a rough fallback when fonts aren't available.
func estimateText(text string, fontSize, maxWidth float64) (float64, float64) {
	charWidth := fontSize * 0.6
	lines := strings.Split(text, "\n")
	totalHeight := 0.0
	maxLineWidth := 0.0

	for _, line := range lines {
		lineWidth := float64(len(line)) * charWidth
		if maxWidth > 0 && lineWidth > maxWidth {
			wrappedLines := int(lineWidth/maxWidth) + 1
			totalHeight += float64(wrappedLines) * fontSize * 1.2
			if maxWidth > maxLineWidth {
				maxLineWidth = maxWidth
			}
		} else {
			totalHeight += fontSize * 1.2
			if lineWidth > maxLineWidth {
				maxLineWidth = lineWidth
			}
		}
	}
	return maxLineWidth, totalHeight
}
