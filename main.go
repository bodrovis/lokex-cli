package main

import (
	"fmt"
	"os"

	"github.com/bodrovis/lokex-cli/cmd"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "command failed: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	return cmd.RootCmd().Execute()
}
