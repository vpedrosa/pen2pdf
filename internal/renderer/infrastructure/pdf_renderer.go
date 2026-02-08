package infrastructure

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"strings"

	"github.com/signintech/gopdf"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/gobolditalic"
	"golang.org/x/image/font/gofont/goitalic"
	"golang.org/x/image/font/gofont/goregular"

	asset "github.com/vpedrosa/pen2pdf/internal/asset/domain"
	layout "github.com/vpedrosa/pen2pdf/internal/layout/domain"
	shared "github.com/vpedrosa/pen2pdf/internal/shared/domain"
)

const fallbackFontFamily = "GoFont"

// fallbackFonts maps style keys to embedded Go font TTF data.
var fallbackFonts = map[string][]byte{
	"regular":    goregular.TTF,
	"bold":       gobold.TTF,
	"italic":     goitalic.TTF,
	"bolditalic": gobolditalic.TTF,
}

// PDFRenderer renders layout pages to PDF using gopdf.
type PDFRenderer struct {
	imageLoader   asset.ImageLoader
	fontLoader    asset.FontLoader
	loadedFonts   map[string]bool
	fallbackReady map[string]bool
	warned        map[string]bool
}

func NewPDFRenderer(imageLoader asset.ImageLoader, fontLoader asset.FontLoader) *PDFRenderer {
	return &PDFRenderer{
		imageLoader:   imageLoader,
		fontLoader:    fontLoader,
		loadedFonts:   make(map[string]bool),
		fallbackReady: make(map[string]bool),
		warned:        make(map[string]bool),
	}
}

func (r *PDFRenderer) Render(pages []layout.Page, output io.Writer) error {
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	for i, page := range pages {
		pdf.AddPageWithOption(gopdf.PageOption{
			PageSize: &gopdf.Rect{W: page.Width, H: page.Height},
		})
		if err := r.renderBox(pdf, page.Root); err != nil {
			return fmt.Errorf("page %d: %w", i, err)
		}
	}

	return pdf.Write(output)
}

func (r *PDFRenderer) renderBox(pdf *gopdf.GoPdf, box *layout.LayoutBox) error {
	switch node := box.Node.(type) {
	case *shared.Frame:
		if err := r.renderFrame(pdf, box, node); err != nil {
			return err
		}
	case *shared.Text:
		if err := r.renderText(pdf, box, node); err != nil {
			return err
		}
	}

	for _, child := range box.Children {
		if err := r.renderBox(pdf, child); err != nil {
			return err
		}
	}
	return nil
}

func (r *PDFRenderer) renderFrame(pdf *gopdf.GoPdf, box *layout.LayoutBox, frame *shared.Frame) error {
	if frame.Fill == nil {
		return nil
	}

	switch frame.Fill.Type {
	case shared.FillSolid:
		return r.drawSolidRect(pdf, box.X, box.Y, box.Width, box.Height, frame.Fill.Color, frame.CornerRadius, frame.Clip)
	case shared.FillImage:
		return r.drawImage(pdf, box.X, box.Y, box.Width, box.Height, frame.Fill, frame.Clip)
	}
	return nil
}

func (r *PDFRenderer) drawSolidRect(pdf *gopdf.GoPdf, x, y, w, h float64, color string, radius float64, _ bool) error {
	if w <= 0 || h <= 0 {
		return nil
	}

	rgba, err := ParseHexColor(color)
	if err != nil {
		return err
	}

	pdf.SetFillColor(rgba.R, rgba.G, rgba.B)
	if rgba.A < 1.0 {
		pdf.SetTransparency(gopdf.Transparency{Alpha: rgba.A, BlendModeType: gopdf.NormalBlendMode})
	}

	if radius > 0 {
		// Clamp radius to half the smallest dimension
		maxRadius := math.Min(w, h) / 2
		if radius > maxRadius {
			radius = maxRadius
		}
		if err := pdf.Rectangle(x, y, x+w, y+h, "F", radius, 20); err != nil {
			return err
		}
	} else {
		pdf.RectFromUpperLeftWithStyle(x, y, w, h, "F")
	}

	if rgba.A < 1.0 {
		pdf.SetTransparency(gopdf.Transparency{Alpha: 1.0, BlendModeType: gopdf.NormalBlendMode})
	}
	return nil
}

func (r *PDFRenderer) drawImage(pdf *gopdf.GoPdf, x, y, w, h float64, fill *shared.Fill, clip bool) error {
	if r.imageLoader == nil {
		return nil
	}

	imgData, err := r.imageLoader.LoadImage(fill.URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: image %q not found, skipping\n", fill.URL)
		return nil
	}

	if fill.Opacity > 0 && fill.Opacity < 1.0 {
		pdf.SetTransparency(gopdf.Transparency{Alpha: fill.Opacity, BlendModeType: gopdf.NormalBlendMode})
	}

	// Calculate fill mode dimensions (cover)
	imgW := float64(imgData.Width)
	imgH := float64(imgData.Height)
	scaleW := w / imgW
	scaleH := h / imgH
	scale := math.Max(scaleW, scaleH)

	drawW := imgW * scale
	drawH := imgH * scale
	drawX := x - (drawW-w)/2
	drawY := y - (drawH-h)/2

	if clip {
		pdf.SaveGraphicsState()
		pdf.ClipPolygon([]gopdf.Point{
			{X: x, Y: y},
			{X: x + w, Y: y},
			{X: x + w, Y: y + h},
			{X: x, Y: y + h},
		})
	}

	imgHolder, err := gopdf.ImageHolderByBytes(imgData.Data)
	if err != nil {
		return fmt.Errorf("create image holder: %w", err)
	}

	if err := pdf.ImageByHolderWithOptions(imgHolder, gopdf.ImageOptions{
		X:    drawX,
		Y:    drawY,
		Rect: &gopdf.Rect{W: drawW, H: drawH},
	}); err != nil {
		return fmt.Errorf("draw image: %w", err)
	}

	if clip {
		pdf.RestoreGraphicsState()
	}

	if fill.Opacity > 0 && fill.Opacity < 1.0 {
		pdf.SetTransparency(gopdf.Transparency{Alpha: 1.0, BlendModeType: gopdf.NormalBlendMode})
	}

	return nil
}

func (r *PDFRenderer) renderText(pdf *gopdf.GoPdf, box *layout.LayoutBox, text *shared.Text) error {
	if text.Content == "" {
		return nil
	}

	fontKey := text.FontFamily + "-" + text.FontWeight
	if !r.loadedFonts[fontKey] {
		if err := r.loadFont(pdf, fontKey, text.FontFamily, text.FontWeight, text.FontStyle); err != nil {
			return err
		}
	}

	if err := pdf.SetFont(fontKey, "", text.FontSize); err != nil {
		return fmt.Errorf("set font %q: %w", fontKey, err)
	}

	// Set text color
	if text.Fill != "" {
		rgba, err := ParseHexColor(text.Fill)
		if err != nil {
			return err
		}
		pdf.SetTextColor(rgba.R, rgba.G, rgba.B)
		if rgba.A < 1.0 {
			pdf.SetTransparency(gopdf.Transparency{Alpha: rgba.A, BlendModeType: gopdf.NormalBlendMode})
		}
	}

	// Set letter spacing
	if text.LetterSpacing != 0 {
		pdf.SetCharSpacing(text.LetterSpacing)
	}

	// Render lines
	lineHeight := text.FontSize * 1.2
	if text.LineHeight > 0 {
		lineHeight = text.FontSize * text.LineHeight
	}

	// Determine horizontal alignment for CellWithOption
	hAlign := gopdf.Left
	switch text.TextAlign {
	case "center":
		hAlign = gopdf.Center
	case "right":
		hAlign = gopdf.Right
	}

	lines := strings.Split(text.Content, "\n")
	currentY := box.Y

	for _, line := range lines {
		pdf.SetX(box.X)
		pdf.SetY(currentY)
		if err := pdf.CellWithOption(&gopdf.Rect{W: box.Width, H: lineHeight}, line, gopdf.CellOption{
			Align: hAlign | gopdf.Top,
		}); err != nil {
			return fmt.Errorf("render text: %w", err)
		}
		currentY += lineHeight
	}

	// Reset letter spacing
	if text.LetterSpacing != 0 {
		pdf.SetCharSpacing(0)
	}
	// Reset transparency
	pdf.SetTransparency(gopdf.Transparency{Alpha: 1.0, BlendModeType: gopdf.NormalBlendMode})

	return nil
}

// loadFont tries the font loader first, then falls back to embedded Go fonts.
func (r *PDFRenderer) loadFont(pdf *gopdf.GoPdf, fontKey, family, weight, style string) error {
	if style == "" {
		style = "normal"
	}

	// Try the real font loader first
	if r.fontLoader != nil {
		fontData, err := r.fontLoader.LoadFont(family, weight, style)
		if err == nil {
			if err := pdf.AddTTFFontByReader(fontKey, bytes.NewReader(fontData.Data)); err != nil {
				return fmt.Errorf("add font %q: %w", fontKey, err)
			}
			r.loadedFonts[fontKey] = true
			return nil
		}
		// Font not found — warn and fall back
		if !r.warned[fontKey] {
			fmt.Fprintf(os.Stderr, "warning: font %q not found, using fallback (Go font)\n", fontKey)
			r.warned[fontKey] = true
		}
	}

	// Fallback to embedded Go fonts
	fbKey := fallbackStyleKey(weight, style)
	fbFontKey := fallbackFontFamily + "-" + fbKey

	if !r.fallbackReady[fbFontKey] {
		data, ok := fallbackFonts[fbKey]
		if !ok {
			data = fallbackFonts["regular"]
			fbFontKey = fallbackFontFamily + "-regular"
		}
		if !r.fallbackReady[fbFontKey] {
			if err := pdf.AddTTFFontData(fbFontKey, data); err != nil {
				return fmt.Errorf("add fallback font: %w", err)
			}
			r.fallbackReady[fbFontKey] = true
		}
	}

	// Map the requested fontKey to the fallback so SetFont works
	if err := pdf.AddTTFFontData(fontKey, fallbackFonts[fbKey]); err != nil {
		// Already loaded under another key, try regular
		if err2 := pdf.AddTTFFontData(fontKey, fallbackFonts["regular"]); err2 != nil {
			return fmt.Errorf("add fallback font for %q: %w", fontKey, err)
		}
	}
	r.loadedFonts[fontKey] = true
	return nil
}

// fallbackStyleKey maps weight/style to a Go font variant key.
func fallbackStyleKey(weight, style string) string {
	isBold := weight == "700" || weight == "800" || weight == "900"
	isItalic := style == "italic"

	switch {
	case isBold && isItalic:
		return "bolditalic"
	case isBold:
		return "bold"
	case isItalic:
		return "italic"
	default:
		return "regular"
	}
}
