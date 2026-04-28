package cmd

import (
	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version info",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("lokex-cli %s\ncommit: %s\nbuilt at: %s\n", version, commit, date)
		},
	}
}
