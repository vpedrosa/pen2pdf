package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"
	"github.com/vpedrosa/pen2pdf/internal/application"
	parserInfra "github.com/vpedrosa/pen2pdf/internal/parser/infrastructure"
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

	svc := application.NewInfoService(parserInfra.NewJSONParser())

	info, err := svc.GetInfo(inputFile)
	if err != nil {
		return err
	}

	cmd.Printf("File:    %s\n", inputPath)
	cmd.Printf("Version: %s\n", info.Version)
	cmd.Println()

	cmd.Printf("Pages (%d):\n", len(info.Pages))
	for _, page := range info.Pages {
		cmd.Printf("  - %s (%.0fx%.0f)\n", page.Name, page.Width, page.Height)
	}
	cmd.Println()

	if len(info.Variables) > 0 {
		cmd.Printf("Variables (%d):\n", len(info.Variables))
		names := make([]string, 0, len(info.Variables))
		for name := range info.Variables {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			v := info.Variables[name]
			cmd.Printf("  %-20s %s = %v\n", name, v.Type, v.Value)
		}
		cmd.Println()
	}

	if len(info.Fonts) > 0 {
		cmd.Printf("Fonts (%d):\n", len(info.Fonts))
		for _, f := range info.Fonts {
			cmd.Printf("  - %s\n", f)
		}
	}

	return nil
}
