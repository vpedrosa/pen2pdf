package domain_test

import (
	"testing"

	resolver "github.com/vpedrosa/pen2pdf/internal/resolver/domain"
	shared "github.com/vpedrosa/pen2pdf/internal/shared/domain"
)

type stubResolver struct {
	err error
}

func (s *stubResolver) Resolve(_ *shared.Document) error {
	return s.err
}

func TestResolverInterfaceCompliance(t *testing.T) {
	var _ resolver.Resolver = &stubResolver{}
}

func TestStubResolverNoError(t *testing.T) {
	r := &stubResolver{}
	doc := &shared.Document{Version: "1.0"}
	if err := r.Resolve(doc); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
