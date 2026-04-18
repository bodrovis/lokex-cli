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
	uploadCfg := &UploadConfig{}

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
	tests := []struct {
		name   string
		result string
		poll   bool
		want   string
	}{
		{
			name:   "poll false prints started with process id",
			result: "process-123",
			poll:   false,
			want:   "Upload started: process-123\n",
		},
		{
			name:   "poll true prints completed with process id",
			result: "bundle-456",
			poll:   true,
			want:   "Upload completed: bundle-456\n",
		},
		{
			name:   "poll false prints unknown when process id is empty",
			result: "",
			poll:   false,
			want:   "Upload started (process ID unknown)\n",
		},
		{
			name:   "poll true prints unknown when process id is empty",
			result: "",
			poll:   true,
			want:   "Upload completed (process ID unknown)\n",
		},
		{
			name:   "poll false trims process id",
			result: "  process-789  ",
			poll:   false,
			want:   "Upload started: process-789\n",
		},
		{
			name:   "poll true treats whitespace-only process id as unknown",
			result: "   \t   ",
			poll:   true,
			want:   "Upload completed (process ID unknown)\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{Use: "test"}
			var out bytes.Buffer
			cmd.SetOut(&out)
			cmd.SetErr(&out)

			printUploadResult(cmd, tt.result, tt.poll)

			if got := out.String(); got != tt.want {
				t.Fatalf("unexpected output: got %q, want %q", got, tt.want)
			}
		})
	}
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

func TestNewCommand_PreRunE_UsesDefaults(t *testing.T) {
	t.Parallel()

	cfg := &global_config.GlobalConfig{
		Token:     "token",
		ProjectID: "project-id",
	}
	defaults := &UploadConfig{
		Filename: new("en.json"),
		LangISO:  new("en"),
	}

	cmd := NewCommand(cfg, defaults)
	if err := cmd.ParseFlags([]string{}); err != nil {
		t.Fatalf("parse flags: %v", err)
	}

	if err := cmd.PreRunE(cmd, nil); err != nil {
		t.Fatalf("PreRunE() error = %v", err)
	}

	gotFilename, err := cmd.Flags().GetString("filename")
	if err != nil {
		t.Fatalf("GetString(filename): %v", err)
	}
	if gotFilename != "en.json" {
		t.Fatalf("expected filename from defaults to be %q, got %q", "en.json", gotFilename)
	}

	gotLangISO, err := cmd.Flags().GetString("lang-iso")
	if err != nil {
		t.Fatalf("GetString(lang-iso): %v", err)
	}
	if gotLangISO != "en" {
		t.Fatalf("expected lang-iso from defaults to be %q, got %q", "en", gotLangISO)
	}
}

func TestNewCommand_PreRunE_ExplicitFlagsOverrideDefaults(t *testing.T) {
	t.Parallel()

	cfg := &global_config.GlobalConfig{
		Token:     "token",
		ProjectID: "project-id",
	}
	defaults := &UploadConfig{
		Filename: new("default.json"),
		LangISO:  new("fr"),
	}

	cmd := NewCommand(cfg, defaults)
	if err := cmd.ParseFlags([]string{"--filename=explicit.json", "--lang-iso=en"}); err != nil {
		t.Fatalf("parse flags: %v", err)
	}

	if err := cmd.PreRunE(cmd, nil); err != nil {
		t.Fatalf("PreRunE() error = %v", err)
	}

	gotFilename, err := cmd.Flags().GetString("filename")
	if err != nil {
		t.Fatalf("GetString(filename): %v", err)
	}
	if gotFilename != "explicit.json" {
		t.Fatalf("expected explicit filename to win, got %q", gotFilename)
	}

	gotLangISO, err := cmd.Flags().GetString("lang-iso")
	if err != nil {
		t.Fatalf("GetString(lang-iso): %v", err)
	}
	if gotLangISO != "en" {
		t.Fatalf("expected explicit lang-iso to win, got %q", gotLangISO)
	}
}

func TestRunCommand_PassesContextToUploader(t *testing.T) {
	tests := []struct {
		name         string
		timeout      time.Duration
		wantDeadline bool
	}{
		{
			name:         "positive timeout adds deadline",
			timeout:      2 * time.Second,
			wantDeadline: true,
		},
		{
			name:         "zero timeout uses background context",
			timeout:      0,
			wantDeadline: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
				Token:          "token",
				ProjectID:      "project-id",
				ContextTimeout: tt.timeout,
			}
			flags := newFlags()

			cmd := newBoundTestCommand(flags)
			if err := cmd.Flags().Parse([]string{"--filename=en.json", "--lang-iso=en"}); err != nil {
				t.Fatalf("parse flags: %v", err)
			}

			var out bytes.Buffer
			cmd.SetOut(&out)
			cmd.SetErr(&out)

			if err := runCommand(cmd, cfg, flags, nil); err != nil {
				t.Fatalf("runCommand() error = %v", err)
			}

			if mu.gotCtx == nil {
				t.Fatal("expected context to be passed to uploader")
			}

			_, gotDeadline := mu.gotCtx.Deadline()
			if gotDeadline != tt.wantDeadline {
				t.Fatalf("unexpected deadline presence: got %v, want %v", gotDeadline, tt.wantDeadline)
			}
		})
	}
}

func TestRunCommand_UsesDefaultsInBuildParams(t *testing.T) {
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

	defaults := &UploadConfig{
		ApplyTM: new(true),
	}

	cmd := newBoundTestCommand(flags)
	if err := cmd.Flags().Parse([]string{"--filename=en.json", "--lang-iso=en"}); err != nil {
		t.Fatalf("parse flags: %v", err)
	}

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)

	if err := runCommand(cmd, cfg, flags, defaults); err != nil {
		t.Fatalf("runCommand() error = %v", err)
	}

	got, ok := mu.gotParams["apply_tm"]
	if !ok {
		t.Fatal("expected apply_tm to be set from defaults")
	}
	if got != true {
		t.Fatalf("expected apply_tm=true, got %#v", got)
	}
}

func TestRunCommand_UploadErrorDoesNotPrintSuccessOutput(t *testing.T) {
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

	cmd := newBoundTestCommand(flags)
	if err := cmd.Flags().Parse([]string{"--filename=en.json", "--lang-iso=en"}); err != nil {
		t.Fatalf("parse flags: %v", err)
	}

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)

	err := runCommand(cmd, cfg, flags, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "upload failed" {
		t.Fatalf("unexpected error: %q", err.Error())
	}

	if strings.Contains(out.String(), "Upload started:") || strings.Contains(out.String(), "Upload completed:") {
		t.Fatalf("did not expect success output on error, got %q", out.String())
	}
}

func TestNewCommand_PreRunE_UsesDefaultsToPassValidation(t *testing.T) {
	t.Parallel()

	cfg := &global_config.GlobalConfig{
		Token:     "token",
		ProjectID: "project-id",
	}
	defaults := &UploadConfig{
		Filename: new("en.json"),
		LangISO:  new("en"),
	}

	cmd := NewCommand(cfg, defaults)
	if err := cmd.ParseFlags([]string{}); err != nil {
		t.Fatalf("parse flags: %v", err)
	}

	if err := cmd.PreRunE(cmd, nil); err != nil {
		t.Fatalf("PreRunE() error = %v", err)
	}

	gotFilename, err := cmd.Flags().GetString("filename")
	if err != nil {
		t.Fatalf("GetString(filename): %v", err)
	}
	if gotFilename != "en.json" {
		t.Fatalf("expected filename from defaults to be %q, got %q", "en.json", gotFilename)
	}

	gotLangISO, err := cmd.Flags().GetString("lang-iso")
	if err != nil {
		t.Fatalf("GetString(lang-iso): %v", err)
	}
	if gotLangISO != "en" {
		t.Fatalf("expected lang-iso from defaults to be %q, got %q", "en", gotLangISO)
	}
}

func TestNewCommand_Execute_UsesDefaultLocalOptions(t *testing.T) {
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
	defaults := &UploadConfig{
		Filename: new("en.json"),
		LangISO:  new("en"),
		SrcPath:  new("./locales/from-default.json"),
		Poll:     new(true),
	}

	cmd := NewCommand(cfg, defaults)
	cmd.SetArgs([]string{})

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !mu.uploadCalled {
		t.Fatal("expected Upload to be called")
	}
	if mu.gotSrcPath != "./locales/from-default.json" {
		t.Fatalf("expected default src path to be used, got %q", mu.gotSrcPath)
	}
	if !mu.gotPoll {
		t.Fatal("expected default poll=true to be used")
	}
	if got := mu.gotParams["filename"]; got != "en.json" {
		t.Fatalf("unexpected filename param: got %#v", got)
	}
	if got := mu.gotParams["lang_iso"]; got != "en" {
		t.Fatalf("unexpected lang_iso param: got %#v", got)
	}

	gotOutput := out.String()
	if !strings.Contains(gotOutput, "Upload completed: bundle-456") {
		t.Fatalf("unexpected output: %q", gotOutput)
	}
}

func TestNewCommand_Execute_ExplicitFlagsOverrideDefaultLocalOptions(t *testing.T) {
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
	defaults := &UploadConfig{
		Filename: new("default.json"),
		LangISO:  new("fr"),
		SrcPath:  new("./locales/from-default.json"),
		Poll:     new(true),
	}

	cmd := NewCommand(cfg, defaults)
	cmd.SetArgs([]string{
		"--filename=explicit.json",
		"--lang-iso=en",
		"--src-path=./locales/explicit.json",
		"--poll=false",
	})

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !mu.uploadCalled {
		t.Fatal("expected Upload to be called")
	}
	if mu.gotSrcPath != "./locales/explicit.json" {
		t.Fatalf("expected explicit src path to win, got %q", mu.gotSrcPath)
	}
	if mu.gotPoll {
		t.Fatal("expected explicit poll=false to win")
	}
	if got := mu.gotParams["filename"]; got != "explicit.json" {
		t.Fatalf("unexpected filename param: got %#v", got)
	}
	if got := mu.gotParams["lang_iso"]; got != "en" {
		t.Fatalf("unexpected lang_iso param: got %#v", got)
	}

	gotOutput := out.String()
	if !strings.Contains(gotOutput, "Upload started: process-123") {
		t.Fatalf("unexpected output: %q", gotOutput)
	}
}

func newBoundTestCommand(flags *Flags) *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	bindFlags(cmd, flags)
	return cmd
}
