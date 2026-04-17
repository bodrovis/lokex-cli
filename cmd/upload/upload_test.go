package upload

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"

	"github.com/bodrovis/lokex-cli/internal/global_config"
	"github.com/bodrovis/lokex-cli/internal/upload_config"
	lokexupload "github.com/bodrovis/lokex/v2/client/upload"
)

type mockUploader struct {
	uploadCalled bool

	gotCtx     context.Context
	gotParams  lokexupload.UploadParams
	gotSrcPath string
	gotPoll    bool

	result string
	err    error
}

func (m *mockUploader) Upload(
	ctx context.Context,
	params lokexupload.UploadParams,
	srcPath string,
	poll bool,
) (string, error) {
	m.uploadCalled = true
	m.gotCtx = ctx
	m.gotParams = params
	m.gotSrcPath = srcPath
	m.gotPoll = poll
	return m.result, m.err
}

func TestNewCommand(t *testing.T) {
	cfg := &global_config.GlobalConfig{}
	uploadCfg := &upload_config.UploadConfig{}

	cmd := NewCommand(cfg, uploadCfg)
	if cmd == nil {
		t.Fatal("expected non-nil command")
	}

	if cmd.Use != "upload" {
		t.Fatalf("unexpected Use: got %q, want %q", cmd.Use, "upload")
	}
	if cmd.Short != "Upload translation files to Lokalise" {
		t.Fatalf("unexpected Short: got %q", cmd.Short)
	}
	if cmd.PreRunE == nil {
		t.Fatal("expected PreRunE to be set")
	}
	if cmd.RunE == nil {
		t.Fatal("expected RunE to be set")
	}
}

func TestValidateCommand(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *global_config.GlobalConfig
		flags   *Flags
		wantErr string
	}{
		{
			name: "ok",
			cfg: &global_config.GlobalConfig{
				Token:     "token",
				ProjectID: "project-id",
			},
			flags: &Flags{
				Filename: "en.json",
				LangISO:  "en",
			},
		},
		{
			name: "missing token",
			cfg: &global_config.GlobalConfig{
				ProjectID: "project-id",
			},
			flags: &Flags{
				Filename: "en.json",
				LangISO:  "en",
			},
			wantErr: "--token is required",
		},
		{
			name: "missing project id",
			cfg: &global_config.GlobalConfig{
				Token: "token",
			},
			flags: &Flags{
				Filename: "en.json",
				LangISO:  "en",
			},
			wantErr: "--project-id is required",
		},
		{
			name: "missing filename",
			cfg: &global_config.GlobalConfig{
				Token:     "token",
				ProjectID: "project-id",
			},
			flags: &Flags{
				LangISO: "en",
			},
			wantErr: "--filename is required",
		},
		{
			name: "whitespace filename",
			cfg: &global_config.GlobalConfig{
				Token:     "token",
				ProjectID: "project-id",
			},
			flags: &Flags{
				Filename: "   ",
				LangISO:  "en",
			},
			wantErr: "--filename is required",
		},
		{
			name: "missing lang iso",
			cfg: &global_config.GlobalConfig{
				Token:     "token",
				ProjectID: "project-id",
			},
			flags: &Flags{
				Filename: "en.json",
			},
			wantErr: "--lang-iso is required",
		},
		{
			name: "whitespace lang iso",
			cfg: &global_config.GlobalConfig{
				Token:     "token",
				ProjectID: "project-id",
			},
			flags: &Flags{
				Filename: "en.json",
				LangISO:  "   ",
			},
			wantErr: "--lang-iso is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCommand(tt.cfg, tt.flags)
			if tt.wantErr == "" && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error %q, got nil", tt.wantErr)
				}
				if err.Error() != tt.wantErr {
					t.Fatalf("unexpected error: got %q, want %q", err.Error(), tt.wantErr)
				}
			}
		})
	}
}

func TestNewCommandContext(t *testing.T) {
	t.Run("background context when timeout is zero", func(t *testing.T) {
		ctx, cancel := newCommandContext(0)
		defer cancel()

		if ctx == nil {
			t.Fatal("expected non-nil context")
		}
		if _, ok := ctx.Deadline(); ok {
			t.Fatal("expected no deadline for zero timeout")
		}
	})

	t.Run("background context when timeout is negative", func(t *testing.T) {
		ctx, cancel := newCommandContext(-1 * time.Second)
		defer cancel()

		if ctx == nil {
			t.Fatal("expected non-nil context")
		}
		if _, ok := ctx.Deadline(); ok {
			t.Fatal("expected no deadline for negative timeout")
		}
	})

	t.Run("timeout context when timeout is positive", func(t *testing.T) {
		ctx, cancel := newCommandContext(2 * time.Second)
		defer cancel()

		if ctx == nil {
			t.Fatal("expected non-nil context")
		}
		if _, ok := ctx.Deadline(); !ok {
			t.Fatal("expected deadline for positive timeout")
		}
	})
}

func TestPerformUpload(t *testing.T) {
	params := lokexupload.UploadParams{
		"filename": "en.json",
		"lang_iso": "en",
	}

	t.Run("success", func(t *testing.T) {
		mu := &mockUploader{
			result: "process-123",
		}
		flags := &Flags{
			SrcPath: "./locales/en.json",
			Poll:    true,
		}

		got, err := performUpload(context.Background(), mu, flags, params)
		if err != nil {
			t.Fatalf("performUpload() error = %v", err)
		}
		if got != "process-123" {
			t.Fatalf("unexpected result: got %q", got)
		}
		if !mu.uploadCalled {
			t.Fatal("expected Upload to be called")
		}
		if mu.gotSrcPath != "./locales/en.json" {
			t.Fatalf("unexpected src path: got %q", mu.gotSrcPath)
		}
		if !mu.gotPoll {
			t.Fatal("expected poll=true to be passed")
		}
		if mu.gotParams["filename"] != "en.json" {
			t.Fatalf("unexpected params: %#v", mu.gotParams)
		}
	})

	t.Run("error", func(t *testing.T) {
		mu := &mockUploader{
			err: errors.New("upload failed"),
		}
		flags := &Flags{
			SrcPath: "./locales/en.json",
			Poll:    false,
		}

		_, err := performUpload(context.Background(), mu, flags, params)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "upload failed" {
			t.Fatalf("unexpected error: %q", err.Error())
		}
	})
}

func TestPrintUploadResult(t *testing.T) {
	t.Run("poll false prints started", func(t *testing.T) {
		cmd := &cobra.Command{Use: "test"}
		var out bytes.Buffer
		cmd.SetOut(&out)
		cmd.SetErr(&out)

		printUploadResult(cmd, "process-123", false)

		got := out.String()
		want := "Upload started: process-123\n"

		if got != want {
			t.Fatalf("unexpected output: got %q, want %q", got, want)
		}
	})

	t.Run("poll true prints completed", func(t *testing.T) {
		cmd := &cobra.Command{Use: "test"}
		var out bytes.Buffer
		cmd.SetOut(&out)
		cmd.SetErr(&out)

		printUploadResult(cmd, "bundle-456", true)

		got := out.String()
		want := "Upload completed: bundle-456\n"

		if got != want {
			t.Fatalf("unexpected output: got %q, want %q", got, want)
		}
	})
}

func TestRunCommand(t *testing.T) {
	t.Run("happy path without poll", func(t *testing.T) {
		old := newUploaderFunc
		t.Cleanup(func() {
			newUploaderFunc = old
		})

		mu := &mockUploader{
			result: "process-123",
		}
		newUploaderFunc = func(cfg *global_config.GlobalConfig) (uploader, error) {
			return mu, nil
		}

		cfg := &global_config.GlobalConfig{
			Token:     "token",
			ProjectID: "project-id",
		}
		flags := newFlags()
		flags.SrcPath = "./locales/en.json"

		cmd := newBoundTestCommand(flags)
		if err := cmd.Flags().Parse([]string{
			"--filename=en.json",
			"--lang-iso=en",
		}); err != nil {
			t.Fatalf("parse flags: %v", err)
		}

		var out bytes.Buffer
		cmd.SetOut(&out)
		cmd.SetErr(&out)

		err := runCommand(cmd, cfg, flags, nil)
		if err != nil {
			t.Fatalf("runCommand() error = %v", err)
		}

		if !mu.uploadCalled {
			t.Fatal("expected Upload to be called")
		}
		if mu.gotSrcPath != "./locales/en.json" {
			t.Fatalf("unexpected src path: got %q", mu.gotSrcPath)
		}
		if mu.gotPoll {
			t.Fatal("expected poll=false to be passed")
		}
		if mu.gotParams["filename"] != "en.json" {
			t.Fatalf("unexpected params: %#v", mu.gotParams)
		}
		if mu.gotParams["lang_iso"] != "en" {
			t.Fatalf("unexpected params: %#v", mu.gotParams)
		}

		gotOutput := out.String()
		if !strings.Contains(gotOutput, "Upload started: process-123") {
			t.Fatalf("unexpected output: %q", gotOutput)
		}
	})

	t.Run("happy path with poll", func(t *testing.T) {
		old := newUploaderFunc
		t.Cleanup(func() {
			newUploaderFunc = old
		})

		mu := &mockUploader{
			result: "bundle-456",
		}
		newUploaderFunc = func(cfg *global_config.GlobalConfig) (uploader, error) {
			return mu, nil
		}

		cfg := &global_config.GlobalConfig{
			Token:     "token",
			ProjectID: "project-id",
		}
		flags := newFlags()
		flags.SrcPath = "./locales/en.json"

		cmd := newBoundTestCommand(flags)
		if err := cmd.Flags().Parse([]string{
			"--filename=en.json",
			"--lang-iso=en",
			"--poll",
		}); err != nil {
			t.Fatalf("parse flags: %v", err)
		}

		var out bytes.Buffer
		cmd.SetOut(&out)
		cmd.SetErr(&out)

		err := runCommand(cmd, cfg, flags, nil)
		if err != nil {
			t.Fatalf("runCommand() error = %v", err)
		}

		if !mu.uploadCalled {
			t.Fatal("expected Upload to be called")
		}
		if !mu.gotPoll {
			t.Fatal("expected poll=true to be passed")
		}

		gotOutput := out.String()
		if !strings.Contains(gotOutput, "Upload completed: bundle-456") {
			t.Fatalf("unexpected output: %q", gotOutput)
		}
	})

	t.Run("uploader factory error", func(t *testing.T) {
		old := newUploaderFunc
		t.Cleanup(func() {
			newUploaderFunc = old
		})

		newUploaderFunc = func(cfg *global_config.GlobalConfig) (uploader, error) {
			return nil, errors.New("cannot create uploader")
		}

		cfg := &global_config.GlobalConfig{
			Token:     "token",
			ProjectID: "project-id",
		}
		flags := newFlags()

		cmd := newBoundTestCommand(flags)
		if err := cmd.Flags().Parse([]string{
			"--filename=en.json",
			"--lang-iso=en",
		}); err != nil {
			t.Fatalf("parse flags: %v", err)
		}

		err := runCommand(cmd, cfg, flags, nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "cannot create uploader" {
			t.Fatalf("unexpected error: %q", err.Error())
		}
	})

	t.Run("upload error", func(t *testing.T) {
		old := newUploaderFunc
		t.Cleanup(func() {
			newUploaderFunc = old
		})

		mu := &mockUploader{
			err: errors.New("upload failed"),
		}
		newUploaderFunc = func(cfg *global_config.GlobalConfig) (uploader, error) {
			return mu, nil
		}

		cfg := &global_config.GlobalConfig{
			Token:     "token",
			ProjectID: "project-id",
		}
		flags := newFlags()
		flags.SrcPath = "./locales/en.json"

		cmd := newBoundTestCommand(flags)
		if err := cmd.Flags().Parse([]string{
			"--filename=en.json",
			"--lang-iso=en",
		}); err != nil {
			t.Fatalf("parse flags: %v", err)
		}

		err := runCommand(cmd, cfg, flags, nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "upload failed" {
			t.Fatalf("unexpected error: %q", err.Error())
		}
	})
}

func newBoundTestCommand(flags *Flags) *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	bindFlags(cmd, flags)
	return cmd
}
