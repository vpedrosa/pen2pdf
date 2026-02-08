package domain

// ImageData holds a loaded image ready for use by the renderer.
type ImageData struct {
	Path   string
	Width  int
	Height int
	Data   []byte
}

// ImageLoader abstracts image loading from the filesystem.
type ImageLoader interface {
	LoadImage(path string) (*ImageData, error)
}
