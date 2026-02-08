//go:build integration

package integration_test

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	assetInfra "github.com/vpedrosa/pen2pdf/internal/asset/infrastructure"
	layoutInfra "github.com/vpedrosa/pen2pdf/internal/layout/infrastructure"
	parserInfra "github.com/vpedrosa/pen2pdf/internal/parser/infrastructure"
	rendererInfra "github.com/vpedrosa/pen2pdf/internal/renderer/infrastructure"
	resolverInfra "github.com/vpedrosa/pen2pdf/internal/resolver/infrastructure"
)

func projectRoot() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "..")
}

func TestFullPipelineWithExampleFile(t *testing.T) {
	root := projectRoot()
	inputPath := filepath.Join(root, "example", "test.pen")

	// Parse
	inputFile, err := os.Open(inputPath)
	if err != nil {
		t.Fatalf("open input: %v", err)
	}
	defer inputFile.Close()

	parser := parserInfra.NewJSONParser()
	doc, err := parser.Parse(inputFile)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	if doc.Version != "2.7" {
		t.Errorf("expected version '2.7', got '%s'", doc.Version)
	}
	if len(doc.Children) != 2 {
		t.Fatalf("expected 2 pages, got %d", len(doc.Children))
	}
	if doc.Children[0].GetName() != "Travel Flyer" {
		t.Errorf("expected first page 'Travel Flyer', got '%s'", doc.Children[0].GetName())
	}
	if doc.Children[1].GetName() != "Travel Flyer - Back" {
		t.Errorf("expected second page 'Travel Flyer - Back', got '%s'", doc.Children[1].GetName())
	}

	// Resolve
	resolver := resolverInfra.NewVariableResolver()
	if err := resolver.Resolve(doc); err != nil {
		t.Fatalf("resolve: %v", err)
	}

	// Layout (without real fonts, using fallback measurement)
	baseDir := filepath.Dir(inputPath)
	fontLoader := assetInfra.NewFSFontLoader(filepath.Join(baseDir, "fonts"))
	measurer := layoutInfra.NewGopdfTextMeasurer(fontLoader)
	layoutEngine := layoutInfra.NewFlexboxEngine()
	pages, err := layoutEngine.Layout(doc, measurer)
	if err != nil {
		t.Fatalf("layout: %v", err)
	}

	if len(pages) != 2 {
		t.Fatalf("expected 2 layout pages, got %d", len(pages))
	}
	if pages[0].Width != 800 || pages[0].Height != 1000 {
		t.Errorf("expected page 1 800x1000, got %fx%f", pages[0].Width, pages[0].Height)
	}

	// Render (without images, since they may not exist in CI)
	renderer := rendererInfra.NewPDFRenderer(nil, nil)
	var buf bytes.Buffer
	if err := renderer.Render(pages, &buf); err != nil {
		t.Fatalf("render: %v", err)
	}

	if buf.Len() == 0 {
		t.Error("expected non-empty PDF output")
	}
	if !bytes.HasPrefix(buf.Bytes(), []byte("%PDF")) {
		t.Error("output does not start with %PDF header")
	}
}

func TestParseAndValidateExampleFile(t *testing.T) {
	root := projectRoot()
	inputPath := filepath.Join(root, "example", "test.pen")

	inputFile, err := os.Open(inputPath)
	if err != nil {
		t.Fatalf("open input: %v", err)
	}
	defer inputFile.Close()

	parser := parserInfra.NewJSONParser()
	doc, err := parser.Parse(inputFile)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	// Verify all variables are defined
	if len(doc.Variables) == 0 {
		t.Error("expected variables in example file")
	}

	// Verify variable resolution succeeds
	resolver := resolverInfra.NewVariableResolver()
	if err := resolver.Resolve(doc); err != nil {
		t.Fatalf("resolve: %v", err)
	}
}
