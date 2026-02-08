package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vpedrosa/pen2pdf/internal/application"
	parserInfra "github.com/vpedrosa/pen2pdf/internal/parser/infrastructure"
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

	svc := application.NewValidateService(
		parserInfra.NewJSONParser(),
		resolverInfra.NewVariableResolver(),
	)

	result, err := svc.Validate(inputFile)
	if err != nil {
		return err
	}

	cmd.Printf("Valid: %s (%d pages, %d variables)\n", inputPath, result.PageCount, result.VariableCount)
	return nil
}
