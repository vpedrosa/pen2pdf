package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var version = "dev"

var rootCmd = &cobra.Command{
	Use:     "pen2pdf",
	Short:   "Generate PDF documents from .pen files",
	Long:    "pen2pdf parses .pen files (Pencil design format) and renders them into high-fidelity PDF output.",
	Version: version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
