package main

import (
	"bytes"
	"errors"
	"os"
	"testing"
)

func TestMain_Success(t *testing.T) {
	oldRun := runFunc
	oldExit := exitFunc
	oldStderr := stderr
	t.Cleanup(func() {
		runFunc = oldRun
		exitFunc = oldExit
		stderr = oldStderr
	})

	runCalled := false
	exitCalled := false

	runFunc = func(args []string) error {
		runCalled = true
		return nil
	}
	exitFunc = func(code int) {
		exitCalled = true
	}

	var errBuf bytes.Buffer
	stderr = &errBuf

	oldArgs := os.Args
	t.Cleanup(func() { os.Args = oldArgs })
	os.Args = []string{"lokex"}

	main()

	if !runCalled {
		t.Fatal("expected run to be called")
	}
	if exitCalled {
		t.Fatal("did not expect exit to be called")
	}
	if errBuf.Len() != 0 {
		t.Fatalf("unexpected stderr output: %q", errBuf.String())
	}
}

func TestMain_Error(t *testing.T) {
	oldRun := runFunc
	oldExit := exitFunc
	oldStderr := stderr
	t.Cleanup(func() {
		runFunc = oldRun
		exitFunc = oldExit
		stderr = oldStderr
	})

	runFunc = func(args []string) error {
		return errors.New("boom")
	}

	var gotExitCode int
	var exitCalled bool
	exitFunc = func(code int) {
		exitCalled = true
		gotExitCode = code
	}

	var errBuf bytes.Buffer
	stderr = &errBuf

	oldArgs := os.Args
	t.Cleanup(func() { os.Args = oldArgs })
	os.Args = []string{"lokex", "bad"}

	main()

	if !exitCalled {
		t.Fatal("expected exit to be called")
	}
	if gotExitCode != 1 {
		t.Fatalf("unexpected exit code: got %d, want 1", gotExitCode)
	}

	got := errBuf.String()
	want := "command failed: boom\n"
	if got != want {
		t.Fatalf("unexpected stderr: got %q, want %q", got, want)
	}
}

func TestRun_Smoke(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "no args",
			args:    nil,
			wantErr: false,
		},
		{
			name:    "help",
			args:    []string{"--help"},
			wantErr: false,
		},
		{
			name:    "invalid command",
			args:    []string{"definitely-not-a-command"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := run(tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("run(%v) error = %v, wantErr = %v", tt.args, err, tt.wantErr)
			}
		})
	}
}
