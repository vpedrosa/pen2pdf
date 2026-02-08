package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
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

	// Parse
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("open input: %w", err)
	}
	defer inputFile.Close()

	parser := parserInfra.NewJSONParser()
	doc, err := parser.Parse(inputFile)
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}

	// Resolve variables
	resolver := resolverInfra.NewVariableResolver()
	if err := resolver.Resolve(doc); err != nil {
		return fmt.Errorf("resolve: %w", err)
	}

	// Filter pages if --pages flag is set
	if pagesFlag != "" {
		doc.Children, err = filterPages(doc.Children, pagesFlag)
		if err != nil {
			return err
		}
	}

	// Set up asset loaders
	baseDir := filepath.Dir(inputPath)
	imageLoader := assetInfra.NewFSImageLoader(baseDir)
	fontsDir := filepath.Join(baseDir, "fonts")

	fontDirs := []string{fontsDir, "/usr/share/fonts", "/usr/local/share/fonts"}
	if home, err := os.UserHomeDir(); err == nil {
		fontDirs = append(fontDirs, filepath.Join(home, ".local", "share", "fonts"))
	}
	fontLoader := assetInfra.NewFSFontLoader(fontDirs...)

	// Check for missing fonts and offer to download
	if err := checkAndDownloadFonts(cmd, doc, fontLoader, fontsDir); err != nil {
		return err
	}

	// Layout
	measurer := layoutInfra.NewGopdfTextMeasurer(fontLoader)
	layoutEngine := layoutInfra.NewFlexboxEngine()
	pages, err := layoutEngine.Layout(doc, measurer)
	if err != nil {
		return fmt.Errorf("layout: %w", err)
	}

	// Render
	outputFile, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("create output: %w", err)
	}
	defer outputFile.Close()

	renderer := rendererInfra.NewPDFRenderer(imageLoader, fontLoader)
	if err := renderer.Render(pages, outputFile); err != nil {
		return fmt.Errorf("render: %w", err)
	}

	cmd.Printf("PDF written to %s (%d pages)\n", output, len(pages))
	return nil
}

func checkAndDownloadFonts(cmd *cobra.Command, doc *shared.Document, fontLoader *assetInfra.FSFontLoader, fontsDir string) error {
	refs := shared.CollectFontRefs(doc)

	var missing []shared.FontRef
	for _, ref := range refs {
		_, err := fontLoader.LoadFont(ref.Family, ref.Weight, ref.Style)
		if err != nil {
			missing = append(missing, ref)
		}
	}

	if len(missing) == 0 {
		return nil
	}

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

func filterPages(children []shared.Node, names string) ([]shared.Node, error) {
	nameList := strings.Split(names, ",")
	nameSet := make(map[string]bool, len(nameList))
	for _, n := range nameList {
		nameSet[strings.TrimSpace(n)] = true
	}

	var filtered []shared.Node
	for _, child := range children {
		if nameSet[child.GetName()] {
			filtered = append(filtered, child)
		}
	}

	if len(filtered) == 0 {
		return nil, fmt.Errorf("no pages match --pages %q", names)
	}
	return filtered, nil
}
