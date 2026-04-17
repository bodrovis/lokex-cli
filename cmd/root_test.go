package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
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

func TestVersionCmd_PrintsVersion(t *testing.T) {
	oldVersion := version
	version = "1.2.3"
	t.Cleanup(func() {
		version = oldVersion
	})

	cmd := newVersionCmd()

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)

	cmd.Run(cmd, nil)

	got := out.String()
	want := "lokex-cli 1.2.3\n"

	if got != want {
		t.Fatalf("unexpected output: got %q, want %q", got, want)
	}
}

func TestGenerateDocs_CreatesMarkdownFiles(t *testing.T) {
	root := RootCmd()
	dir := t.TempDir()

	if err := generateDocs(root, dir); err != nil {
		t.Fatalf("generateDocs() returned error: %v", err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("failed to read docs dir: %v", err)
	}

	if len(entries) == 0 {
		t.Fatal("expected generated docs, but directory is empty")
	}

	var foundRootDoc bool
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if strings.HasSuffix(entry.Name(), ".md") {
			foundRootDoc = true
			break
		}
	}

	if !foundRootDoc {
		t.Fatal("expected at least one markdown file to be generated")
	}
}

func TestGenDocsCmd_Execute_CreatesDocsDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir to temp dir: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(oldWD)
	})

	root := RootCmd()
	root.SetArgs([]string{"gendocs"})

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() returned error: %v", err)
	}

	docsDir := filepath.Join(tmpDir, "docs")
	info, err := os.Stat(docsDir)
	if err != nil {
		t.Fatalf("expected docs directory to exist: %v", err)
	}
	if !info.IsDir() {
		t.Fatalf("expected %q to be a directory", docsDir)
	}

	entries, err := os.ReadDir(docsDir)
	if err != nil {
		t.Fatalf("failed to read generated docs dir: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("expected generated docs files, but docs directory is empty")
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
