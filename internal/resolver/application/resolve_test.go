package application_test

import (
	"fmt"
	"testing"

	"github.com/vpedrosa/pen2pdf/internal/resolver/application"
	shared "github.com/vpedrosa/pen2pdf/internal/shared/domain"
)

type stubResolver struct {
	err error
}

func (r *stubResolver) Resolve(_ *shared.Document) error {
	return r.err
}

func TestResolveServiceSuccess(t *testing.T) {
	svc := application.NewResolveService(&stubResolver{})
	err := svc.Resolve(&shared.Document{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestResolveServiceError(t *testing.T) {
	svc := application.NewResolveService(&stubResolver{err: fmt.Errorf("resolve failed")})
	err := svc.Resolve(&shared.Document{})
	if err == nil {
		t.Fatal("expected error")
	}
}
