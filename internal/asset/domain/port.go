package domain

// ImageData holds a loaded image ready for use by the renderer.
type ImageData struct {
	Path   string
	Width  int
	Height int
	Data   []byte
}

// FontData holds a loaded font ready for use by the layout engine and renderer.
type FontData struct {
	Family string
	Weight string
	Style  string
	Path   string
	Data   []byte
}

// ImageLoader abstracts image loading from the filesystem.
type ImageLoader interface {
	LoadImage(path string) (*ImageData, error)
}

// FontLoader abstracts font loading and resolution.
type FontLoader interface {
	LoadFont(family, weight, style string) (*FontData, error)
}
