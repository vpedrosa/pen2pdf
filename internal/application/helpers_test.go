package application

import (
	asset "github.com/vpedrosa/pen2pdf/internal/asset/domain"
	layout "github.com/vpedrosa/pen2pdf/internal/layout/domain"
	parser "github.com/vpedrosa/pen2pdf/internal/parser/domain"
	renderer "github.com/vpedrosa/pen2pdf/internal/renderer/domain"
	resolver "github.com/vpedrosa/pen2pdf/internal/resolver/domain"

	assetInfra "github.com/vpedrosa/pen2pdf/internal/asset/infrastructure"
	layoutInfra "github.com/vpedrosa/pen2pdf/internal/layout/infrastructure"
	parserInfra "github.com/vpedrosa/pen2pdf/internal/parser/infrastructure"
	rendererInfra "github.com/vpedrosa/pen2pdf/internal/renderer/infrastructure"
	resolverInfra "github.com/vpedrosa/pen2pdf/internal/resolver/infrastructure"
)

func newJSONParser() parser.Parser {
	return parserInfra.NewJSONParser()
}

func newVariableResolver() resolver.Resolver {
	return resolverInfra.NewVariableResolver()
}

func newEmptyFontLoader() asset.FontLoader {
	return assetInfra.NewFSFontLoader() // no dirs = always returns "not found"
}

func newFlexboxEngine() layout.LayoutEngine {
	return layoutInfra.NewFlexboxEngine()
}

func newMeasurer(fl asset.FontLoader) layout.TextMeasurer {
	return layoutInfra.NewGopdfTextMeasurer(fl)
}

func newPDFRenderer(fl asset.FontLoader) renderer.Renderer {
	return rendererInfra.NewPDFRenderer(nil, fl)
}
