package application_test

import (
	"fmt"
	"testing"

	"github.com/vpedrosa/pen2pdf/internal/layout/application"
	layout "github.com/vpedrosa/pen2pdf/internal/layout/domain"
	shared "github.com/vpedrosa/pen2pdf/internal/shared/domain"
)

type stubLayoutEngine struct {
	pages []layout.Page
	err   error
}

func (e *stubLayoutEngine) Layout(_ *shared.Document, _ layout.TextMeasurer) ([]layout.Page, error) {
	return e.pages, e.err
}

func TestLayoutServiceSuccess(t *testing.T) {
	pages := []layout.Page{{Width: 800, Height: 1000}}
	svc := application.NewLayoutService(&stubLayoutEngine{pages: pages}, nil)

	result, err := svc.Layout(&shared.Document{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 page, got %d", len(result))
	}
	if result[0].Width != 800 {
		t.Errorf("expected width 800, got %f", result[0].Width)
	}
}

func TestLayoutServiceError(t *testing.T) {
	svc := application.NewLayoutService(&stubLayoutEngine{err: fmt.Errorf("layout failed")}, nil)

	_, err := svc.Layout(&shared.Document{})
	if err == nil {
		t.Fatal("expected error")
	}
}
