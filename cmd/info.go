package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"
	parserInfra "github.com/vpedrosa/pen2pdf/internal/parser/infrastructure"
	shared "github.com/vpedrosa/pen2pdf/internal/shared/domain"
)

var infoCmd = &cobra.Command{
	Use:   "info [input.pen]",
	Short: "Display document metadata",
	Long:  "Parses a .pen file and shows page count, dimensions, variables, and fonts used.",
	Args:  cobra.ExactArgs(1),
	RunE:  runInfo,
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

func runInfo(cmd *cobra.Command, args []string) error {
	inputPath := args[0]

	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("open input: %w", err)
	}
	defer inputFile.Close() //nolint:errcheck

	parser := parserInfra.NewJSONParser()
	doc, err := parser.Parse(inputFile)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	cmd.Printf("File:    %s\n", inputPath)
	cmd.Printf("Version: %s\n", doc.Version)
	cmd.Println()

	// Pages
	cmd.Printf("Pages (%d):\n", len(doc.Children))
	for _, child := range doc.Children {
		if frame, ok := child.(*shared.Frame); ok {
			cmd.Printf("  - %s (%.0fx%.0f)\n", frame.Name, frame.Width.Value, frame.Height.Value)
		}
	}
	cmd.Println()

	// Variables
	if len(doc.Variables) > 0 {
		cmd.Printf("Variables (%d):\n", len(doc.Variables))
		names := make([]string, 0, len(doc.Variables))
		for name := range doc.Variables {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			v := doc.Variables[name]
			cmd.Printf("  %-20s %s = %v\n", name, v.Type, v.Value)
		}
		cmd.Println()
	}

	// Fonts
	fonts := collectFonts(doc.Children)
	if len(fonts) > 0 {
		cmd.Printf("Fonts (%d):\n", len(fonts))
		for _, f := range fonts {
			cmd.Printf("  - %s\n", f)
		}
	}

	return nil
}

// collectFonts walks the node tree and returns sorted unique font families.
func collectFonts(nodes []shared.Node) []string {
	fontSet := make(map[string]bool)
	walkFonts(nodes, fontSet)

	fonts := make([]string, 0, len(fontSet))
	for f := range fontSet {
		fonts = append(fonts, f)
	}
	sort.Strings(fonts)
	return fonts
}

func walkFonts(nodes []shared.Node, fonts map[string]bool) {
	for _, node := range nodes {
		switch n := node.(type) {
		case *shared.Frame:
			walkFonts(n.Children, fonts)
		case *shared.Text:
			if n.FontFamily != "" {
				fonts[n.FontFamily] = true
			}
		}
	}
}
