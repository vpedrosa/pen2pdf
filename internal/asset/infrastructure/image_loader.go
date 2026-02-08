package infrastructure

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"

	asset "github.com/vpedrosa/pen2pdf/internal/asset/domain"
)

// FSImageLoader loads images from the filesystem relative to a base directory.
type FSImageLoader struct {
	baseDir string
}

func NewFSImageLoader(baseDir string) *FSImageLoader {
	return &FSImageLoader{baseDir: baseDir}
}

func (l *FSImageLoader) LoadImage(path string) (*asset.ImageData, error) {
	absPath := path
	if !filepath.IsAbs(path) {
		absPath = filepath.Join(l.baseDir, path)
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("load image %q: %w", path, err)
	}

	f, err := os.Open(absPath)
	if err != nil {
		return nil, fmt.Errorf("open image %q: %w", path, err)
	}
	defer f.Close() //nolint:errcheck

	cfg, _, err := image.DecodeConfig(f)
	if err != nil {
		return nil, fmt.Errorf("decode image config %q: %w", path, err)
	}

	return &asset.ImageData{
		Path:   absPath,
		Width:  cfg.Width,
		Height: cfg.Height,
		Data:   data,
	}, nil
}
