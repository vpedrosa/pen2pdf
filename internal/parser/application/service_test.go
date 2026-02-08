package application_test

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/vpedrosa/pen2pdf/internal/parser/application"
	shared "github.com/vpedrosa/pen2pdf/internal/shared/domain"
)

type stubParser struct {
	doc *shared.Document
	err error
}

func (p *stubParser) Parse(_ io.Reader) (*shared.Document, error) {
	return p.doc, p.err
}

func TestParseServiceSuccess(t *testing.T) {
	doc := &shared.Document{Version: "1.0"}
	svc := application.NewParseService(&stubParser{doc: doc})

	result, err := svc.Parse(strings.NewReader("{}"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Version != "1.0" {
		t.Errorf("expected version '1.0', got '%s'", result.Version)
	}
}

func TestParseServiceError(t *testing.T) {
	svc := application.NewParseService(&stubParser{err: fmt.Errorf("parse failed")})

	_, err := svc.Parse(strings.NewReader("bad"))
	if err == nil {
		t.Fatal("expected error")
	}
}
