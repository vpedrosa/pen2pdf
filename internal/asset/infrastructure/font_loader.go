package infrastructure

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	asset "github.com/vpedrosa/pen2pdf/internal/asset/domain"
)

// FSFontLoader discovers and loads font files from configurable directories.
type FSFontLoader struct {
	fontDirs []string
}

func NewFSFontLoader(fontDirs ...string) *FSFontLoader {
	return &FSFontLoader{fontDirs: fontDirs}
}

func (l *FSFontLoader) LoadFont(family, weight, style string) (*asset.FontData, error) {
	candidates := fontFileCandidates(family, weight, style)

	for _, dir := range l.fontDirs {
		for _, candidate := range candidates {
			path := filepath.Join(dir, candidate)
			data, err := os.ReadFile(path)
			if err != nil {
				continue
			}
			return &asset.FontData{
				Family: family,
				Weight: weight,
				Style:  style,
				Path:   path,
				Data:   data,
			}, nil
		}
	}

	return nil, fmt.Errorf("font not found: %s %s %s (searched %v)", family, weight, style, l.fontDirs)
}

// fontFileCandidates generates potential filename patterns for a font.
// Example: family="Open Sans", weight="700", style="italic"
// → ["OpenSans-BoldItalic.ttf", "OpenSans-BoldItalic.otf",
//
//	"Open Sans-BoldItalic.ttf", ...]
func fontFileCandidates(family, weight, style string) []string {
	suffix := weightToSuffix(weight)
	if style == "italic" && suffix != "" {
		suffix += "Italic"
	} else if style == "italic" {
		suffix = "Italic"
	}
	if suffix == "" {
		suffix = "Regular"
	}

	familyNoSpaces := strings.ReplaceAll(family, " ", "")

	var candidates []string
	for _, name := range []string{familyNoSpaces, family} {
		for _, ext := range []string{".ttf", ".otf"} {
			candidates = append(candidates, name+"-"+suffix+ext)
		}
	}
	return candidates
}

func weightToSuffix(weight string) string {
	switch weight {
	case "100":
		return "Thin"
	case "200":
		return "ExtraLight"
	case "300":
		return "Light"
	case "400", "normal", "":
		return ""
	case "500":
		return "Medium"
	case "600":
		return "SemiBold"
	case "700":
		return "Bold"
	case "800":
		return "ExtraBold"
	case "900":
		return "Black"
	default:
		return weight
	}
}
