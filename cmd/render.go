package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vpedrosa/pen2pdf/internal/application"
	assetInfra "github.com/vpedrosa/pen2pdf/internal/asset/infrastructure"
	layoutInfra "github.com/vpedrosa/pen2pdf/internal/layout/infrastructure"
	parserInfra "github.com/vpedrosa/pen2pdf/internal/parser/infrastructure"
	rendererInfra "github.com/vpedrosa/pen2pdf/internal/renderer/infrastructure"
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

	parser := parserInfra.NewJSONParser()
	resolver := resolverInfra.NewVariableResolver()
	fontLoader := assetInfra.NewFSFontLoader(fontDirs...)
	imageLoader := assetInfra.NewFSImageLoader(baseDir)
	measurer := layoutInfra.NewGopdfTextMeasurer(fontLoader)
	layoutEngine := layoutInfra.NewFlexboxEngine()
	renderer := rendererInfra.NewPDFRenderer(imageLoader, fontLoader)

	svc := application.NewRenderService(parser, resolver, fontLoader, imageLoader, layoutEngine, measurer, renderer)

	// Check for missing fonts (interactive CLI concern)
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("open input: %w", err)
	}

	missing, doc, err := svc.DetectMissingFonts(inputFile)
	inputFile.Close() //nolint:errcheck
	if err != nil {
		return err
	}

	if len(missing) > 0 {
		if err := promptAndDownloadFonts(cmd, missing, fontsDir); err != nil {
			return err
		}
	}

	// Render using the already-parsed document
	outputFile, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("create output: %w", err)
	}
	defer outputFile.Close() //nolint:errcheck

	result, err := svc.RenderDocument(doc, outputFile, pagesFlag)
	if err != nil {
		return err
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
