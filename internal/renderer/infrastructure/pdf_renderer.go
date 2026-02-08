package infrastructure

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"strings"

	"github.com/signintech/gopdf"
	asset "github.com/vpedrosa/pen2pdf/internal/asset/domain"
	layout "github.com/vpedrosa/pen2pdf/internal/layout/domain"
	shared "github.com/vpedrosa/pen2pdf/internal/shared/domain"
)

// PDFRenderer renders layout pages to PDF using gopdf.
type PDFRenderer struct {
	imageLoader asset.ImageLoader
	fontLoader  asset.FontLoader
	loadedFonts map[string]bool
}

func NewPDFRenderer(imageLoader asset.ImageLoader, fontLoader asset.FontLoader) *PDFRenderer {
	return &PDFRenderer{
		imageLoader: imageLoader,
		fontLoader:  fontLoader,
		loadedFonts: make(map[string]bool),
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
	rgba, err := ParseHexColor(color)
	if err != nil {
		return err
	}

	pdf.SetFillColor(rgba.R, rgba.G, rgba.B)
	if rgba.A < 1.0 {
		pdf.SetTransparency(gopdf.Transparency{Alpha: rgba.A, BlendModeType: gopdf.NormalBlendMode})
	}

	if radius > 0 {
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
		return fmt.Errorf("load image %q: %w", fill.URL, err)
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

	// Load and set font
	fontKey := text.FontFamily + "-" + text.FontWeight
	if !r.loadedFonts[fontKey] && r.fontLoader != nil {
		style := text.FontStyle
		if style == "" {
			style = "normal"
		}
		fontData, err := r.fontLoader.LoadFont(text.FontFamily, text.FontWeight, style)
		if err != nil {
			return fmt.Errorf("load font %q: %w", fontKey, err)
		}
		if err := pdf.AddTTFFontByReader(fontKey, bytes.NewReader(fontData.Data)); err != nil {
			return fmt.Errorf("add font %q: %w", fontKey, err)
		}
		r.loadedFonts[fontKey] = true
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

	lines := strings.Split(text.Content, "\n")
	currentY := box.Y

	for _, line := range lines {
		lineWidth, _ := pdf.MeasureTextWidth(line)
		textX := box.X

		switch text.TextAlign {
		case "center":
			textX = box.X + (box.Width-lineWidth)/2
		case "right":
			textX = box.X + box.Width - lineWidth
		}

		pdf.SetX(textX)
		pdf.SetY(currentY)
		if err := pdf.Text(line); err != nil {
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
