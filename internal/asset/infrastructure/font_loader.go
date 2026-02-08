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
			// Search recursively within each font directory
			var found string
			filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
				if err != nil || found != "" || d.IsDir() {
					return nil
				}
				if strings.EqualFold(d.Name(), candidate) {
					found = path
				}
				return nil
			})
			if found != "" {
				data, err := os.ReadFile(found)
				if err != nil {
					continue
				}
				return &asset.FontData{
					Family: family,
					Weight: weight,
					Style:  style,
					Path:   found,
					Data:   data,
				}, nil
			}
		}
	}

	return nil, fmt.Errorf("font not found: %s %s %s (searched %v)", family, weight, style, l.fontDirs)
}

// fontFileCandidates generates potential filename patterns for a font.
// Includes static font names (e.g., "Inter-Bold.ttf") and variable font
// names (e.g., "Inter-VariableFont_wght.ttf").
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

	// Static font files (e.g., Inter-Bold.ttf, Poppins-SemiBold.ttf)
	for _, name := range []string{familyNoSpaces, family} {
		for _, ext := range []string{".ttf", ".otf"} {
			candidates = append(candidates, name+"-"+suffix+ext)
		}
	}

	// Variable font files — a single file contains all weights
	if style == "italic" {
		for _, name := range []string{familyNoSpaces, family} {
			candidates = append(candidates,
				name+"-Italic-VariableFont_wght.ttf",
				name+"-Italic-VariableFont_opsz,wght.ttf",
				name+"-Italic-VariableFont_wdth,wght.ttf",
			)
		}
	}
	// Non-italic variable font (also used as fallback for italic if italic variant doesn't exist)
	for _, name := range []string{familyNoSpaces, family} {
		candidates = append(candidates,
			name+"-VariableFont_wght.ttf",
			name+"-VariableFont_opsz,wght.ttf",
			name+"-VariableFont_wdth,wght.ttf",
		)
	}

	// Last resort: Regular static font (for single-weight fonts like Bebas Neue)
	if suffix != "Regular" {
		for _, name := range []string{familyNoSpaces, family} {
			for _, ext := range []string{".ttf", ".otf"} {
				candidates = append(candidates, name+"-Regular"+ext)
			}
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
