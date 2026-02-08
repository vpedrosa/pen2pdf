package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"
	parserApp "github.com/vpedrosa/pen2pdf/internal/parser/application"
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

	parseSvc := parserApp.NewParseService(parserInfra.NewJSONParser())

	doc, err := parseSvc.Parse(inputFile)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	cmd.Printf("File:    %s\n", inputPath)
	cmd.Printf("Version: %s\n", doc.Version)
	cmd.Println()

	cmd.Printf("Pages (%d):\n", len(doc.Children))
	for _, child := range doc.Children {
		if frame, ok := child.(*shared.Frame); ok {
			cmd.Printf("  - %s (%.0fx%.0f)\n", frame.Name, frame.Width.Value, frame.Height.Value)
		}
	}
	cmd.Println()

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

	fonts := shared.CollectFontFamilies(doc)
	if len(fonts) > 0 {
		cmd.Printf("Fonts (%d):\n", len(fonts))
		for _, f := range fonts {
			cmd.Printf("  - %s\n", f)
		}
	}

	return nil
}
