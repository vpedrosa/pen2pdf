package domain_test

import (
	"io"
	"strings"
	"testing"

	parser "github.com/vpedrosa/pen2pdf/internal/parser/domain"
	shared "github.com/vpedrosa/pen2pdf/internal/shared/domain"
)

type stubParser struct {
	doc *shared.Document
	err error
}

func (s *stubParser) Parse(_ io.Reader) (*shared.Document, error) {
	return s.doc, s.err
}

func TestParserInterfaceCompliance(t *testing.T) {
	var _ parser.Parser = &stubParser{}
}

func TestStubParserReturnsDocument(t *testing.T) {
	expected := &shared.Document{Version: "2.7"}
	p := &stubParser{doc: expected}

	got, err := p.Parse(strings.NewReader(""))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Version != "2.7" {
		t.Errorf("expected version '2.7', got '%s'", got.Version)
	}
}

func TestStubParserReturnsError(t *testing.T) {
	p := &stubParser{err: io.ErrUnexpectedEOF}

	_, err := p.Parse(strings.NewReader(""))
	if err != io.ErrUnexpectedEOF {
		t.Errorf("expected ErrUnexpectedEOF, got %v", err)
	}
}
