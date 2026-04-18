package main

import (
	"fmt"
	"os"

	"github.com/bodrovis/lokex-cli/cmd"
	"github.com/spf13/cobra"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "command failed: %v\n", err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	return cmd.RootCmd()
}

func run(args []string) error {
	root := newRootCmd()
	root.SetArgs(args)

	return root.Execute()
}
