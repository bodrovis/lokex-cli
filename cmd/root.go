package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"

	downloadcmd "github.com/bodrovis/lokex-cli/cmd/download"
	uploadcmd "github.com/bodrovis/lokex-cli/cmd/upload"
	"github.com/bodrovis/lokex-cli/internal/cli"
)

var version = "dev"

func RootCmd() *cobra.Command {
	cfg := cli.NewGlobalConfig()

	rootCmd := &cobra.Command{
		Use:              "lokex-cli",
		Short:            "CLI for uploading and downloading files with Lokalise",
		TraverseChildren: true,
		SilenceUsage:     true,
		SilenceErrors:    true,
	}

	cli.BindPersistentFlags(rootCmd.PersistentFlags(), cfg)

	rootCmd.AddCommand(newVersionCmd())
	rootCmd.AddCommand(newGenDocsCmd(rootCmd))
	rootCmd.AddCommand(downloadcmd.NewCommand(cfg))
	rootCmd.AddCommand(uploadcmd.NewCommand(cfg))

	return rootCmd
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version info",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("lokex-cli %s\n", version)
		},
	}
}

func newGenDocsCmd(rootCmd *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:    "gendocs",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			if err := generateDocs(rootCmd, "./docs"); err != nil {
				fmt.Fprintf(os.Stderr, "error generating docs: %v\n", err)
			}
		},
	}
}

func generateDocs(rootCmd *cobra.Command, dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	return doc.GenMarkdownTree(rootCmd, dir)
}
