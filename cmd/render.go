package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	assetApp "github.com/vpedrosa/pen2pdf/internal/asset/application"
	assetInfra "github.com/vpedrosa/pen2pdf/internal/asset/infrastructure"
	layoutApp "github.com/vpedrosa/pen2pdf/internal/layout/application"
	layoutInfra "github.com/vpedrosa/pen2pdf/internal/layout/infrastructure"
	parserApp "github.com/vpedrosa/pen2pdf/internal/parser/application"
	parserInfra "github.com/vpedrosa/pen2pdf/internal/parser/infrastructure"
	rendererApp "github.com/vpedrosa/pen2pdf/internal/renderer/application"
	rendererInfra "github.com/vpedrosa/pen2pdf/internal/renderer/infrastructure"
	resolverApp "github.com/vpedrosa/pen2pdf/internal/resolver/application"
	resolverInfra "github.com/vpedrosa/pen2pdf/internal/resolver/infrastructure"
	shared "github.com/vpedrosa/pen2pdf/internal/shared/domain"
)

var (
	outputPath string
	pagesFlag  string
	noPrompt   bool
)

var renderCmd = &cobra.Command{
	Use:   "render [input.pen]",
	Short: "Render a .pen file to PDF",
	Long:  "Parses the .pen file and renders it into a high-fidelity PDF document.",
	Args:  cobra.ExactArgs(1),
	RunE:  runRender,
}

func init() {
	renderCmd.Flags().StringVarP(&outputPath, "output", "o", "", "output PDF file path (default: input with .pdf extension)")
	renderCmd.Flags().StringVar(&pagesFlag, "pages", "", "comma-separated page names to render (default: all)")
	renderCmd.Flags().BoolVar(&noPrompt, "no-prompt", false, "skip interactive prompts (for CI/scripts)")
	rootCmd.AddCommand(renderCmd)
}

func runRender(cmd *cobra.Command, args []string) error {
	inputPath := args[0]

	// Determine output path
	output := outputPath
	if output == "" {
		ext := filepath.Ext(inputPath)
		output = strings.TrimSuffix(inputPath, ext) + ".pdf"
	}

	// Build infrastructure
	baseDir := filepath.Dir(inputPath)
	fontsDir := filepath.Join(baseDir, "fonts")

	fontDirs := []string{fontsDir, "/usr/share/fonts", "/usr/local/share/fonts"}
	if home, err := os.UserHomeDir(); err == nil {
		fontDirs = append(fontDirs, filepath.Join(home, ".local", "share", "fonts"))
	}

	fontLoader := assetInfra.NewFSFontLoader(fontDirs...)
	imageLoader := assetInfra.NewFSImageLoader(baseDir)
	measurer := layoutInfra.NewGopdfTextMeasurer(fontLoader)
	pdfRenderer := rendererInfra.NewPDFRenderer(imageLoader, fontLoader)

	// Build application services (inject ports via DI)
	parseSvc := parserApp.NewParseService(parserInfra.NewJSONParser())
	resolveSvc := resolverApp.NewResolveService(resolverInfra.NewVariableResolver())
	fontSvc := assetApp.NewFontService(fontLoader)
	layoutSvc := layoutApp.NewLayoutService(layoutInfra.NewFlexboxEngine(), measurer)
	renderSvc := rendererApp.NewRenderService(pdfRenderer)

	// 1. Parse
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("open input: %w", err)
	}
	doc, err := parseSvc.Parse(inputFile)
	inputFile.Close() //nolint:errcheck
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}

	// 2. Resolve variables
	if err := resolveSvc.Resolve(doc); err != nil {
		return fmt.Errorf("resolve: %w", err)
	}

	// 3. Detect and download missing fonts (interactive CLI concern)
	missing := fontSvc.DetectMissingFonts(doc)
	if len(missing) > 0 {
		if err := promptAndDownloadFonts(cmd, missing, fontsDir); err != nil {
			return err
		}
	}

	// 4. Filter pages
	if pagesFlag != "" {
		doc.Children, err = shared.FilterPagesByName(doc.Children, pagesFlag)
		if err != nil {
			return err
		}
	}

	// 5. Layout
	pages, err := layoutSvc.Layout(doc)
	if err != nil {
		return fmt.Errorf("layout: %w", err)
	}

	// 6. Render
	outputFile, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("create output: %w", err)
	}
	defer outputFile.Close() //nolint:errcheck

	result, err := renderSvc.Render(pages, outputFile)
	if err != nil {
		return fmt.Errorf("render: %w", err)
	}

	cmd.Printf("PDF written to %s (%d pages)\n", output, result.PageCount)
	return nil
}

func promptAndDownloadFonts(cmd *cobra.Command, missing []shared.FontRef, fontsDir string) error {
	cmd.Printf("Missing %d font(s):\n", len(missing))
	for _, ref := range missing {
		label := ref.Family + " " + ref.Weight
		if ref.Style != "" {
			label += " " + ref.Style
		}
		cmd.Printf("  - %s\n", label)
	}

	if noPrompt {
		cmd.Println("Skipping download (--no-prompt). Fallback fonts will be used.")
		return nil
	}

	cmd.Printf("\nDownload from Google Fonts to %s? [Y/n] ", fontsDir)
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))

	if answer != "" && answer != "y" && answer != "yes" && answer != "s" && answer != "si" && answer != "sí" {
		cmd.Println("Skipping download. Fallback fonts will be used.")
		return nil
	}

	downloader := assetInfra.NewGoogleFontDownloader()
	downloaded := 0
	for _, ref := range missing {
		path, err := downloader.Download(ref.Family, ref.Weight, ref.Style, fontsDir)
		if err != nil {
			cmd.PrintErrf("  warning: could not download %s %s: %v\n", ref.Family, ref.Weight, err)
			continue
		}
		cmd.Printf("  downloaded: %s\n", filepath.Base(path))
		downloaded++
	}

	if downloaded > 0 {
		cmd.Printf("%d font(s) downloaded to %s\n\n", downloaded, fontsDir)
	}

	return nil
}
