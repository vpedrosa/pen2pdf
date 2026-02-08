package domain

import (
	"fmt"
	"strconv"
)

// RGBA holds a parsed color with alpha channel.
type RGBA struct {
	R, G, B uint8
	A        float64 // 0.0 to 1.0
}

// ParseHexColor parses #RRGGBB or #RRGGBBAA hex color strings.
func ParseHexColor(hex string) (RGBA, error) {
	if len(hex) == 0 || hex[0] != '#' {
		return RGBA{}, fmt.Errorf("invalid hex color: %q", hex)
	}

	hex = hex[1:]
	switch len(hex) {
	case 6:
		r, g, b, err := parseRGB(hex)
		if err != nil {
			return RGBA{}, err
		}
		return RGBA{R: r, G: g, B: b, A: 1.0}, nil
	case 8:
		r, g, b, err := parseRGB(hex[:6])
		if err != nil {
			return RGBA{}, err
		}
		a, err := strconv.ParseUint(hex[6:8], 16, 8)
		if err != nil {
			return RGBA{}, fmt.Errorf("invalid alpha: %w", err)
		}
		return RGBA{R: r, G: g, B: b, A: float64(a) / 255.0}, nil
	default:
		return RGBA{}, fmt.Errorf("invalid hex color length: %q", "#"+hex)
	}
}

func parseRGB(hex string) (uint8, uint8, uint8, error) {
	r, err := strconv.ParseUint(hex[0:2], 16, 8)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid red: %w", err)
	}
	g, err := strconv.ParseUint(hex[2:4], 16, 8)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid green: %w", err)
	}
	b, err := strconv.ParseUint(hex[4:6], 16, 8)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid blue: %w", err)
	}
	return uint8(r), uint8(g), uint8(b), nil
}
