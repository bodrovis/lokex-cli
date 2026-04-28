package cmd

import (
	"bytes"
	"testing"
)

func TestVersionCmd_PrintsVersion(t *testing.T) {
	oldVersion := version
	version = "1.2.3"
	oldCommit := commit
	commit = "42"

	oldDate := date
	date = "some date"
	t.Cleanup(func() {
		version = oldVersion
		commit = oldCommit
		date = oldDate
	})

	cmd := newVersionCmd()

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)

	cmd.Run(cmd, nil)

	got := out.String()
	want := "lokex-cli 1.2.3\ncommit: 42\nbuilt at: some date\n"

	if got != want {
		t.Fatalf("unexpected output: got %q, want %q", got, want)
	}
}
