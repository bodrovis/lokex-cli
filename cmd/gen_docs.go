package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func newGenDocsCmd(rootCmd *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:    "gendocs",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return generateDocs(rootCmd, "./docs")
		},
	}
}

func generateDocs(rootCmd *cobra.Command, dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	return doc.GenMarkdownTree(rootCmd, dir)
}
