package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	parserApp "github.com/vpedrosa/pen2pdf/internal/parser/application"
	parserInfra "github.com/vpedrosa/pen2pdf/internal/parser/infrastructure"
	resolverApp "github.com/vpedrosa/pen2pdf/internal/resolver/application"
	resolverInfra "github.com/vpedrosa/pen2pdf/internal/resolver/infrastructure"
)

var validateCmd = &cobra.Command{
	Use:   "validate [input.pen]",
	Short: "Validate a .pen file without rendering",
	Long:  "Parses and validates a .pen file, checking for syntax errors and undefined variable references.",
	Args:  cobra.ExactArgs(1),
	RunE:  runValidate,
}

func init() {
	rootCmd.AddCommand(validateCmd)
}

func runValidate(cmd *cobra.Command, args []string) error {
	inputPath := args[0]

	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("open input: %w", err)
	}
	defer inputFile.Close() //nolint:errcheck

	parseSvc := parserApp.NewParseService(parserInfra.NewJSONParser())
	resolveSvc := resolverApp.NewResolveService(resolverInfra.NewVariableResolver())

	doc, err := parseSvc.Parse(inputFile)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	if err := resolveSvc.Resolve(doc); err != nil {
		return fmt.Errorf("resolve error: %w", err)
	}

	cmd.Printf("Valid: %s (%d pages, %d variables)\n", inputPath, len(doc.Children), len(doc.Variables))
	return nil
}
