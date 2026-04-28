package main

import (
	"fmt"
	"io"
	"os"

	"github.com/bodrovis/lokex-cli/cmd"
	"github.com/spf13/cobra"
)

var (
	runFunc            = run
	exitFunc           = os.Exit
	stderr   io.Writer = os.Stderr
)

func main() {
	if err := runFunc(os.Args[1:]); err != nil {
		_, _ = fmt.Fprintf(stderr, "command failed: %+v\n", err)
		exitFunc(1)
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
