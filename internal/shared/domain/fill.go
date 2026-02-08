package domain

type FillType string

const (
	FillSolid FillType = "solid"
	FillImage FillType = "image"
)

type Fill struct {
	Type    FillType
	Color   string
	URL     string
	Mode    string
	Opacity float64
	Enabled bool
}

func SolidFill(color string) *Fill {
	return &Fill{Type: FillSolid, Color: color}
}

func ImageFill(url, mode string, opacity float64, enabled bool) *Fill {
	return &Fill{
		Type:    FillImage,
		URL:     url,
		Mode:    mode,
		Opacity: opacity,
		Enabled: enabled,
	}
}
