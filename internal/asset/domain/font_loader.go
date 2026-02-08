package domain

// FontData holds a loaded font ready for use by the layout engine and renderer.
type FontData struct {
	Family string
	Weight string
	Style  string
	Path   string
	Data   []byte
}

// FontLoader abstracts font loading and resolution.
type FontLoader interface {
	LoadFont(family, weight, style string) (*FontData, error)
}
