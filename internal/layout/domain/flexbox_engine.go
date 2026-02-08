package domain

import (
	"fmt"

	shared "github.com/vpedrosa/pen2pdf/internal/shared/domain"
)

// FlexboxEngine implements a subset of CSS Flexbox layout for .pen files.
type FlexboxEngine struct{}

func NewFlexboxEngine() *FlexboxEngine {
	return &FlexboxEngine{}
}

func (e *FlexboxEngine) Layout(doc *shared.Document, measurer TextMeasurer) ([]Page, error) {
	pages := make([]Page, 0, len(doc.Children))
	for _, child := range doc.Children {
		frame, ok := child.(*shared.Frame)
		if !ok {
			return nil, fmt.Errorf("top-level node %q must be a frame", child.GetID())
		}

		root := layoutFrame(frame, 0, 0, frame.Width.Value, frame.Height.Value, measurer)
		pages = append(pages, Page{
			Width:  frame.Width.Value,
			Height: frame.Height.Value,
			Root:   root,
		})
	}
	return pages, nil
}

func layoutFrame(frame *shared.Frame, x, y, w, h float64, measurer TextMeasurer) *LayoutBox {
	box := &LayoutBox{
		X:      x,
		Y:      y,
		Width:  w,
		Height: h,
		Node:   frame,
	}

	if len(frame.Children) == 0 {
		return box
	}

	contentX := x + frame.Padding.Left
	contentY := y + frame.Padding.Top
	contentW := w - frame.Padding.Left - frame.Padding.Right
	contentH := h - frame.Padding.Top - frame.Padding.Bottom

	isVertical := frame.Layout == "vertical"

	// Phase 1: Measure all children to determine fixed vs fill_container sizes
	type childInfo struct {
		node       shared.Node
		width      float64
		height     float64
		fillWidth  bool
		fillHeight bool
	}

	children := make([]childInfo, len(frame.Children))
	totalFixedMain := 0.0
	fillCount := 0
	gaps := 0.0
	if len(frame.Children) > 1 {
		gaps = float64(len(frame.Children)-1) * frame.Gap
	}

	for i, child := range frame.Children {
		info := childInfo{node: child}

		switch n := child.(type) {
		case *shared.Frame:
			info.fillWidth = n.Width.FillContainer
			info.fillHeight = n.Height.FillContainer
			if !info.fillWidth {
				info.width = n.Width.Value
			}
			if !info.fillHeight {
				info.height = n.Height.Value
			}
			// Auto-size: compute intrinsic size when dimension is missing
			if (info.width == 0 && !info.fillWidth) || (info.height == 0 && !info.fillHeight) {
				iw, ih := intrinsicSize(n, measurer, contentW)
				if info.width == 0 && !info.fillWidth {
					info.width = iw
				}
				if info.height == 0 && !info.fillHeight {
					info.height = ih
				}
			}
		case *shared.Text:
			info.fillWidth = n.Width.FillContainer
			if !info.fillWidth && n.Width.Value > 0 {
				info.width = n.Width.Value
			}
			// Measure text intrinsic size
			if measurer != nil {
				maxW := info.width
				if maxW == 0 {
					maxW = contentW
				}
				tw, th := measurer.MeasureText(n.Content, n.FontFamily, n.FontSize, n.FontWeight, maxW)
				if info.width == 0 && !info.fillWidth {
					info.width = tw
				}
				if info.height == 0 {
					info.height = th
				}
			}
		}

		if isVertical {
			if info.fillHeight {
				fillCount++
			} else {
				totalFixedMain += info.height
			}
		} else {
			if info.fillWidth {
				fillCount++
			} else {
				totalFixedMain += info.width
			}
		}

		children[i] = info
	}

	// Phase 2: Calculate fill_container sizes
	remainingMain := 0.0
	if isVertical {
		remainingMain = contentH - totalFixedMain - gaps
	} else {
		remainingMain = contentW - totalFixedMain - gaps
	}
	if remainingMain < 0 {
		remainingMain = 0
	}

	fillSize := 0.0
	if fillCount > 0 {
		fillSize = remainingMain / float64(fillCount)
	}

	for i := range children {
		if isVertical {
			if children[i].fillHeight {
				children[i].height = fillSize
			}
			if children[i].fillWidth {
				children[i].width = contentW
			}
		} else {
			if children[i].fillWidth {
				children[i].width = fillSize
			}
			if children[i].fillHeight {
				children[i].height = contentH
			}
		}
	}

	// Phase 3: Position children based on justifyContent and alignItems
	totalUsedMain := totalFixedMain + float64(fillCount)*fillSize + gaps

	var mainOffset float64
	var mainSpacing float64

	switch frame.JustifyContent {
	case "center":
		if isVertical {
			mainOffset = (contentH - totalUsedMain) / 2
		} else {
			mainOffset = (contentW - totalUsedMain) / 2
		}
	case "end":
		if isVertical {
			mainOffset = contentH - totalUsedMain
		} else {
			mainOffset = contentW - totalUsedMain
		}
	case "space-between":
		if len(children) > 1 {
			totalWithoutGaps := totalFixedMain + float64(fillCount)*fillSize
			if isVertical {
				mainSpacing = (contentH - totalWithoutGaps) / float64(len(children)-1)
			} else {
				mainSpacing = (contentW - totalWithoutGaps) / float64(len(children)-1)
			}
		}
	default: // "start" or empty
		mainOffset = 0
	}

	currentMain := mainOffset

	for i, info := range children {
		var childX, childY float64
		var childW, childH float64

		childW = info.width
		childH = info.height

		if isVertical {
			childY = contentY + currentMain
			childX = contentX + crossOffset(frame.AlignItems, contentW, childW)
			currentMain += childH
		} else {
			childX = contentX + currentMain
			childY = contentY + crossOffset(frame.AlignItems, contentH, childH)
			currentMain += childW
		}

		// Add gap or space-between spacing
		if i < len(children)-1 {
			if frame.JustifyContent == "space-between" {
				currentMain += mainSpacing
			} else {
				currentMain += frame.Gap
			}
		}

		var childBox *LayoutBox
		switch n := info.node.(type) {
		case *shared.Frame:
			childBox = layoutFrame(n, childX, childY, childW, childH, measurer)
		case *shared.Text:
			childBox = &LayoutBox{
				X:      childX,
				Y:      childY,
				Width:  childW,
				Height: childH,
				Node:   n,
			}
		}

		box.Children = append(box.Children, childBox)
	}

	return box
}

// intrinsicSize computes the natural size of a frame based on its children.
// Used when a frame has no explicit width/height and is not fill_container.
func intrinsicSize(frame *shared.Frame, measurer TextMeasurer, availableW float64) (float64, float64) {
	padH := frame.Padding.Left + frame.Padding.Right
	padV := frame.Padding.Top + frame.Padding.Bottom

	if len(frame.Children) == 0 {
		return padH, padV
	}

	contentW := availableW - padH
	if contentW < 0 {
		contentW = 0
	}

	isVertical := frame.Layout == "vertical"

	var totalMain, maxCross float64
	gaps := 0.0
	if len(frame.Children) > 1 {
		gaps = float64(len(frame.Children)-1) * frame.Gap
	}

	for _, child := range frame.Children {
		var cw, ch float64

		switch n := child.(type) {
		case *shared.Frame:
			if n.Width.Value > 0 {
				cw = n.Width.Value
			} else if !n.Width.FillContainer {
				cw, _ = intrinsicSize(n, measurer, contentW)
			}
			if n.Height.Value > 0 {
				ch = n.Height.Value
			} else if !n.Height.FillContainer {
				_, ch = intrinsicSize(n, measurer, contentW)
			}
		case *shared.Text:
			if n.Width.Value > 0 {
				cw = n.Width.Value
			}
			if measurer != nil {
				maxW := cw
				if maxW == 0 {
					maxW = contentW
				}
				tw, th := measurer.MeasureText(n.Content, n.FontFamily, n.FontSize, n.FontWeight, maxW)
				if cw == 0 {
					cw = tw
				}
				ch = th
			}
		}

		if isVertical {
			totalMain += ch
			if cw > maxCross {
				maxCross = cw
			}
		} else {
			totalMain += cw
			if ch > maxCross {
				maxCross = ch
			}
		}
	}

	if isVertical {
		return maxCross + padH, totalMain + gaps + padV
	}
	return totalMain + gaps + padH, maxCross + padV
}

func crossOffset(alignItems string, available, size float64) float64 {
	switch alignItems {
	case "center":
		return (available - size) / 2
	case "end":
		return available - size
	default: // "start" or empty
		return 0
	}
}
