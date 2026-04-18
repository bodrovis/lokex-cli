package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestRootCmd_HasExpectedCommands(t *testing.T) {
	root := RootCmd()

	if root == nil {
		t.Fatal("RootCmd() returned nil")
	}

	if root.Use != "lokex-cli" {
		t.Fatalf("unexpected root Use: got %q, want %q", root.Use, "lokex-cli")
	}

	expected := []string{"version", "gendocs", "download", "upload"}
	for _, name := range expected {
		if findSubcommand(root, name) == nil {
			t.Fatalf("expected subcommand %q to be registered", name)
		}
	}
}

func findSubcommand(root *cobra.Command, name string) *cobra.Command {
	for _, cmd := range root.Commands() {
		if cmd.Name() == name {
			return cmd
		}
	}
	return nil
}
