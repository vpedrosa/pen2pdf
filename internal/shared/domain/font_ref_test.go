package domain_test

import (
	"testing"

	"github.com/vpedrosa/pen2pdf/internal/shared/domain"
)

func TestCollectFontRefsSingleText(t *testing.T) {
	doc := &domain.Document{
		Children: []domain.Node{
			&domain.Frame{
				ID: "page", Name: "page",
				Children: []domain.Node{
					&domain.Text{ID: "t1", FontFamily: "Inter", FontWeight: "400", FontStyle: "normal"},
				},
			},
		},
	}

	refs := domain.CollectFontRefs(doc)
	if len(refs) != 1 {
		t.Fatalf("expected 1 ref, got %d", len(refs))
	}
	if refs[0].Family != "Inter" || refs[0].Weight != "400" || refs[0].Style != "normal" {
		t.Errorf("unexpected ref: %+v", refs[0])
	}
}

func TestCollectFontRefsDeduplicate(t *testing.T) {
	doc := &domain.Document{
		Children: []domain.Node{
			&domain.Frame{
				ID: "page", Name: "page",
				Children: []domain.Node{
					&domain.Text{ID: "t1", FontFamily: "Inter", FontWeight: "400", FontStyle: "normal"},
					&domain.Text{ID: "t2", FontFamily: "Inter", FontWeight: "400", FontStyle: "normal"},
					&domain.Text{ID: "t3", FontFamily: "Inter", FontWeight: "700", FontStyle: "normal"},
				},
			},
		},
	}

	refs := domain.CollectFontRefs(doc)
	if len(refs) != 2 {
		t.Fatalf("expected 2 unique refs, got %d", len(refs))
	}
}

func TestCollectFontRefsNestedFrames(t *testing.T) {
	doc := &domain.Document{
		Children: []domain.Node{
			&domain.Frame{
				ID: "page", Name: "page",
				Children: []domain.Node{
					&domain.Frame{
						ID: "inner", Name: "inner",
						Children: []domain.Node{
							&domain.Frame{
								ID: "deep", Name: "deep",
								Children: []domain.Node{
									&domain.Text{ID: "t1", FontFamily: "Poppins", FontWeight: "600", FontStyle: "italic"},
								},
							},
						},
					},
				},
			},
		},
	}

	refs := domain.CollectFontRefs(doc)
	if len(refs) != 1 {
		t.Fatalf("expected 1 ref from nested frame, got %d", len(refs))
	}
	if refs[0].Family != "Poppins" {
		t.Errorf("expected family 'Poppins', got '%s'", refs[0].Family)
	}
}

func TestCollectFontRefsEmptyDocument(t *testing.T) {
	doc := &domain.Document{}
	refs := domain.CollectFontRefs(doc)
	if len(refs) != 0 {
		t.Errorf("expected 0 refs for empty doc, got %d", len(refs))
	}
}

func TestCollectFontRefsSkipsEmptyFamily(t *testing.T) {
	doc := &domain.Document{
		Children: []domain.Node{
			&domain.Frame{
				ID: "page", Name: "page",
				Children: []domain.Node{
					&domain.Text{ID: "t1", FontFamily: "", FontWeight: "400"},
					&domain.Text{ID: "t2", FontFamily: "Inter", FontWeight: "400"},
				},
			},
		},
	}

	refs := domain.CollectFontRefs(doc)
	if len(refs) != 1 {
		t.Fatalf("expected 1 ref (skip empty family), got %d", len(refs))
	}
	if refs[0].Family != "Inter" {
		t.Errorf("expected family 'Inter', got '%s'", refs[0].Family)
	}
}

func TestCollectFontRefsMultiplePages(t *testing.T) {
	doc := &domain.Document{
		Children: []domain.Node{
			&domain.Frame{
				ID: "p1", Name: "page1",
				Children: []domain.Node{
					&domain.Text{ID: "t1", FontFamily: "Inter", FontWeight: "400"},
				},
			},
			&domain.Frame{
				ID: "p2", Name: "page2",
				Children: []domain.Node{
					&domain.Text{ID: "t2", FontFamily: "Roboto", FontWeight: "700"},
				},
			},
		},
	}

	refs := domain.CollectFontRefs(doc)
	if len(refs) != 2 {
		t.Fatalf("expected 2 refs from 2 pages, got %d", len(refs))
	}
}
