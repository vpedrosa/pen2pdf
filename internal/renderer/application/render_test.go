package application_test

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	layout "github.com/vpedrosa/pen2pdf/internal/layout/domain"
	"github.com/vpedrosa/pen2pdf/internal/renderer/application"
)

type stubRenderer struct {
	err error
}

func (r *stubRenderer) Render(_ []layout.Page, w io.Writer) error {
	if r.err != nil {
		return r.err
	}
	_, _ = w.Write([]byte("%PDF-stub"))
	return nil
}

func TestRenderServiceSuccess(t *testing.T) {
	svc := application.NewRenderService(&stubRenderer{})
	pages := []layout.Page{{Width: 800, Height: 1000}, {Width: 800, Height: 1000}}

	var buf bytes.Buffer
	result, err := svc.Render(pages, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.PageCount != 2 {
		t.Errorf("expected 2 pages, got %d", result.PageCount)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty output")
	}
}

func TestRenderServiceError(t *testing.T) {
	svc := application.NewRenderService(&stubRenderer{err: fmt.Errorf("render failed")})

	var buf bytes.Buffer
	_, err := svc.Render([]layout.Page{}, &buf)
	if err == nil {
		t.Fatal("expected error")
	}
}
