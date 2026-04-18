package cmd

import (
	"bytes"
	"testing"
)

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
