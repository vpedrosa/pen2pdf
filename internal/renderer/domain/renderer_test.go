package domain_test

import (
	"io"
	"testing"

	layout "github.com/vpedrosa/pen2pdf/internal/layout/domain"
	renderer "github.com/vpedrosa/pen2pdf/internal/renderer/domain"
)

type stubRenderer struct {
	err error
}

func (s *stubRenderer) Render(_ []layout.Page, _ io.Writer) error {
	return s.err
}

func TestRendererInterfaceCompliance(t *testing.T) {
	var _ renderer.Renderer = &stubRenderer{}
}

func TestStubRendererNoError(t *testing.T) {
	r := &stubRenderer{}
	if err := r.Render(nil, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
