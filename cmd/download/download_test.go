package download

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"

	globalCfg "github.com/bodrovis/lokex-cli/internal/global_config"
	lokexdownload "github.com/bodrovis/lokex/v2/client/download"
)

type mockDownloader struct {
	downloadCalled      bool
	downloadAsyncCalled bool

	gotCtx    context.Context
	gotOut    string
	gotParams lokexdownload.DownloadParams

	downloadURL      string
	downloadAsyncURL string

	downloadErr      error
	downloadAsyncErr error
}

func (m *mockDownloader) Download(ctx context.Context, out string, params lokexdownload.DownloadParams) (string, error) {
	m.downloadCalled = true
	m.gotCtx = ctx
	m.gotOut = out
	m.gotParams = params
	return m.downloadURL, m.downloadErr
}

func (m *mockDownloader) DownloadAsync(ctx context.Context, out string, params lokexdownload.DownloadParams) (string, error) {
	m.downloadAsyncCalled = true
	m.gotCtx = ctx
	m.gotOut = out
	m.gotParams = params
	return m.downloadAsyncURL, m.downloadAsyncErr
}

func TestNewDownloader(t *testing.T) {
	t.Run("returns error when client config is invalid", func(t *testing.T) {
		cfg := &globalCfg.GlobalConfig{}

		got, err := newDownloader(cfg)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if got != nil {
			t.Fatalf("expected nil downloader, got %#v", got)
		}
	})

	t.Run("returns downloader when client config is valid", func(t *testing.T) {
		cfg := &globalCfg.GlobalConfig{
			Token:     "token",
			ProjectID: "project-id",
		}

		got, err := newDownloader(cfg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got == nil {
			t.Fatal("expected non-nil downloader")
		}
	})
}

func TestNewCommand(t *testing.T) {
	cfg := &globalCfg.GlobalConfig{}

	cmd := NewCommand(cfg, nil)
	if cmd == nil {
		t.Fatal("expected non-nil command")
	}

	if cmd.Use != "download" {
		t.Fatalf("unexpected Use: got %q, want %q", cmd.Use, "download")
	}
	if cmd.Short != "Download translation files from Lokalise" {
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
		cfg     *globalCfg.GlobalConfig
		flags   *Flags
		wantErr string
	}{
		{
			name: "ok",
			cfg: &globalCfg.GlobalConfig{
				Token:     "token",
				ProjectID: "project-id",
			},
			flags: &Flags{
				Format: "json",
			},
		},
		{
			name: "missing token",
			cfg: &globalCfg.GlobalConfig{
				ProjectID: "project-id",
			},
			flags: &Flags{
				Format: "json",
			},
			wantErr: "--token is required",
		},
		{
			name: "missing project id",
			cfg: &globalCfg.GlobalConfig{
				Token: "token",
			},
			flags: &Flags{
				Format: "json",
			},
			wantErr: "--project-id is required",
		},
		{
			name: "missing format",
			cfg: &globalCfg.GlobalConfig{
				Token:     "token",
				ProjectID: "project-id",
			},
			flags:   &Flags{},
			wantErr: "--format is required",
		},
		{
			name: "whitespace format",
			cfg: &globalCfg.GlobalConfig{
				Token:     "token",
				ProjectID: "project-id",
			},
			flags: &Flags{
				Format: "   ",
			},
			wantErr: "--format is required",
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

func TestPerformDownload(t *testing.T) {
	params := lokexdownload.DownloadParams{
		"format": "json",
	}

	t.Run("sync download", func(t *testing.T) {
		md := &mockDownloader{
			downloadURL: "https://example.com/sync.zip",
		}
		flags := &Flags{
			Out:   "./locales",
			Async: false,
		}

		got, err := performDownload(context.Background(), md, flags, params)
		if err != nil {
			t.Fatalf("performDownload() error = %v", err)
		}
		if got != "https://example.com/sync.zip" {
			t.Fatalf("unexpected url: got %q", got)
		}
		if !md.downloadCalled {
			t.Fatal("expected Download to be called")
		}
		if md.downloadAsyncCalled {
			t.Fatal("did not expect DownloadAsync to be called")
		}
		if md.gotOut != "./locales" {
			t.Fatalf("unexpected out: got %q", md.gotOut)
		}
		if md.gotParams["format"] != "json" {
			t.Fatalf("unexpected params: %#v", md.gotParams)
		}
	})

	t.Run("async download", func(t *testing.T) {
		md := &mockDownloader{
			downloadAsyncURL: "https://example.com/async.zip",
		}
		flags := &Flags{
			Out:   "./locales",
			Async: true,
		}

		got, err := performDownload(context.Background(), md, flags, params)
		if err != nil {
			t.Fatalf("performDownload() error = %v", err)
		}
		if got != "https://example.com/async.zip" {
			t.Fatalf("unexpected url: got %q", got)
		}
		if !md.downloadAsyncCalled {
			t.Fatal("expected DownloadAsync to be called")
		}
		if md.downloadCalled {
			t.Fatal("did not expect Download to be called")
		}
		if md.gotOut != "./locales" {
			t.Fatalf("unexpected out: got %q", md.gotOut)
		}
		if md.gotParams["format"] != "json" {
			t.Fatalf("unexpected params: %#v", md.gotParams)
		}
	})

	t.Run("sync error", func(t *testing.T) {
		md := &mockDownloader{
			downloadErr: errors.New("sync failed"),
		}
		flags := &Flags{
			Out: "./locales",
		}

		_, err := performDownload(context.Background(), md, flags, params)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "sync failed" {
			t.Fatalf("unexpected error: %q", err.Error())
		}
	})

	t.Run("async error", func(t *testing.T) {
		md := &mockDownloader{
			downloadAsyncErr: errors.New("async failed"),
		}
		flags := &Flags{
			Out:   "./locales",
			Async: true,
		}

		_, err := performDownload(context.Background(), md, flags, params)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "async failed" {
			t.Fatalf("unexpected error: %q", err.Error())
		}
	})
}

func TestPrintDownloadResult(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)

	printDownloadResult(cmd, "https://example.com/file.zip")

	got := out.String()
	want := "Bundle downloaded from: https://example.com/file.zip\n"

	if got != want {
		t.Fatalf("unexpected output: got %q, want %q", got, want)
	}
}

func TestTruncateURLForOutput(t *testing.T) {
	tests := []struct {
		name string
		url  string
		max  int
		want string
	}{
		{
			name: "max zero",
			url:  "https://example.com/file.zip",
			max:  0,
			want: "",
		},
		{
			name: "max negative",
			url:  "https://example.com/file.zip",
			max:  -1,
			want: "",
		},
		{
			name: "shorter than max",
			url:  "short",
			max:  10,
			want: "short",
		},
		{
			name: "equal to max",
			url:  "exact",
			max:  5,
			want: "exact",
		},
		{
			name: "longer than max",
			url:  "abcdefghij",
			max:  7,
			want: "abcd...",
		},
		{
			name: "max less than or equal to three",
			url:  "abcdefghij",
			max:  3,
			want: "abc",
		},
		{
			name: "max one",
			url:  "abcdefghij",
			max:  1,
			want: "a",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncateURLForOutput(tt.url, tt.max)
			if got != tt.want {
				t.Fatalf("unexpected result: got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNewCommand_Execute_RunE(t *testing.T) {
	old := newDownloaderFunc
	t.Cleanup(func() {
		newDownloaderFunc = old
	})

	md := &mockDownloader{
		downloadURL: "https://example.com/file.zip",
	}
	newDownloaderFunc = func(cfg *globalCfg.GlobalConfig) (downloader, error) {
		return md, nil
	}

	cfg := &globalCfg.GlobalConfig{
		Token:     "token",
		ProjectID: "project-id",
	}

	cmd := NewCommand(cfg, nil)

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{
		"--format=json",
		"--out=./locales",
	})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !md.downloadCalled {
		t.Fatal("expected Download to be called")
	}
	if md.gotOut != "./locales" {
		t.Fatalf("unexpected out: got %q", md.gotOut)
	}
	if md.gotParams["format"] != "json" {
		t.Fatalf("unexpected params: %#v", md.gotParams)
	}
}

func TestRunCommand(t *testing.T) {
	t.Run("sync happy path", func(t *testing.T) {
		old := newDownloaderFunc
		t.Cleanup(func() {
			newDownloaderFunc = old
		})

		md := &mockDownloader{
			downloadURL: "https://example.com/file.zip",
		}
		newDownloaderFunc = func(cfg *globalCfg.GlobalConfig) (downloader, error) {
			return md, nil
		}

		cfg := &globalCfg.GlobalConfig{
			Token:     "token",
			ProjectID: "project-id",
		}
		flags := newFlags()
		flags.Out = "./locales"

		cmd := newBoundTestCommand(flags)
		if err := cmd.Flags().Parse([]string{"--format=json"}); err != nil {
			t.Fatalf("parse flags: %v", err)
		}

		var out bytes.Buffer
		cmd.SetOut(&out)
		cmd.SetErr(&out)

		err := runCommand(cmd, cfg, flags, nil)
		if err != nil {
			t.Fatalf("runCommand() error = %v", err)
		}

		if !md.downloadCalled {
			t.Fatal("expected Download to be called")
		}
		if md.downloadAsyncCalled {
			t.Fatal("did not expect DownloadAsync to be called")
		}
		if md.gotOut != "./locales" {
			t.Fatalf("unexpected out: got %q", md.gotOut)
		}
		if md.gotParams["format"] != "json" {
			t.Fatalf("unexpected params: %#v", md.gotParams)
		}

		gotOutput := out.String()
		if !strings.Contains(gotOutput, "Bundle downloaded from: https://example.com/file.zip") {
			t.Fatalf("unexpected output: %q", gotOutput)
		}
	})

	t.Run("async happy path", func(t *testing.T) {
		old := newDownloaderFunc
		t.Cleanup(func() {
			newDownloaderFunc = old
		})

		md := &mockDownloader{
			downloadAsyncURL: "https://example.com/async.zip",
		}
		newDownloaderFunc = func(cfg *globalCfg.GlobalConfig) (downloader, error) {
			return md, nil
		}

		cfg := &globalCfg.GlobalConfig{
			Token:     "token",
			ProjectID: "project-id",
		}
		flags := newFlags()
		flags.Out = "./locales"

		cmd := newBoundTestCommand(flags)
		if err := cmd.Flags().Parse([]string{"--format=json", "--async"}); err != nil {
			t.Fatalf("parse flags: %v", err)
		}

		var out bytes.Buffer
		cmd.SetOut(&out)
		cmd.SetErr(&out)

		err := runCommand(cmd, cfg, flags, nil)
		if err != nil {
			t.Fatalf("runCommand() error = %v", err)
		}

		if !md.downloadAsyncCalled {
			t.Fatal("expected DownloadAsync to be called")
		}
		if md.downloadCalled {
			t.Fatal("did not expect Download to be called")
		}
	})

	t.Run("downloader factory error", func(t *testing.T) {
		old := newDownloaderFunc
		t.Cleanup(func() {
			newDownloaderFunc = old
		})

		newDownloaderFunc = func(cfg *globalCfg.GlobalConfig) (downloader, error) {
			return nil, errors.New("cannot create downloader")
		}

		cfg := &globalCfg.GlobalConfig{
			Token:     "token",
			ProjectID: "project-id",
		}
		flags := newFlags()
		flags.Format = "json"

		cmd := newBoundTestCommand(flags)

		err := runCommand(cmd, cfg, flags, nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "cannot create downloader" {
			t.Fatalf("unexpected error: %q", err.Error())
		}
	})

	t.Run("build params error", func(t *testing.T) {
		old := newDownloaderFunc
		t.Cleanup(func() {
			newDownloaderFunc = old
		})

		md := &mockDownloader{}
		newDownloaderFunc = func(cfg *globalCfg.GlobalConfig) (downloader, error) {
			return md, nil
		}

		cfg := &globalCfg.GlobalConfig{
			Token:     "token",
			ProjectID: "project-id",
		}
		flags := newFlags()

		cmd := newBoundTestCommand(flags)
		if err := cmd.Flags().Parse([]string{"--format=json", "--language-mapping={"}); err != nil {
			t.Fatalf("parse flags: %v", err)
		}

		err := runCommand(cmd, cfg, flags, nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "parse --language-mapping") {
			t.Fatalf("unexpected error: %q", err.Error())
		}
		if md.downloadCalled || md.downloadAsyncCalled {
			t.Fatal("did not expect downloader to be called when buildParams fails")
		}
	})

	t.Run("download error", func(t *testing.T) {
		old := newDownloaderFunc
		t.Cleanup(func() {
			newDownloaderFunc = old
		})

		md := &mockDownloader{
			downloadErr: errors.New("download failed"),
		}
		newDownloaderFunc = func(cfg *globalCfg.GlobalConfig) (downloader, error) {
			return md, nil
		}

		cfg := &globalCfg.GlobalConfig{
			Token:     "token",
			ProjectID: "project-id",
		}
		flags := newFlags()
		flags.Format = "json"

		cmd := newBoundTestCommand(flags)

		err := runCommand(cmd, cfg, flags, nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "download failed" {
			t.Fatalf("unexpected error: %q", err.Error())
		}
	})
}

func TestNewCommand_PreRunE_UsesDefaults(t *testing.T) {
	t.Parallel()

	cfg := &globalCfg.GlobalConfig{
		Token:     "token",
		ProjectID: "project-id",
	}
	defaults := &DownloadConfig{
		Format: new("json"),
	}

	cmd := NewCommand(cfg, defaults)
	if err := cmd.ParseFlags([]string{}); err != nil {
		t.Fatalf("parse flags: %v", err)
	}

	if err := cmd.PreRunE(cmd, nil); err != nil {
		t.Fatalf("PreRunE() error = %v", err)
	}

	got, err := cmd.Flags().GetString("format")
	if err != nil {
		t.Fatalf("GetString(format): %v", err)
	}
	if got != "json" {
		t.Fatalf("expected format from defaults to be %q, got %q", "json", got)
	}
}

func TestNewCommand_PreRunE_ExplicitFlagOverridesDefaults(t *testing.T) {
	t.Parallel()

	cfg := &globalCfg.GlobalConfig{
		Token:     "token",
		ProjectID: "project-id",
	}
	defaults := &DownloadConfig{
		Format: new("json"),
	}

	cmd := NewCommand(cfg, defaults)
	if err := cmd.ParseFlags([]string{"--format=xml"}); err != nil {
		t.Fatalf("parse flags: %v", err)
	}

	if err := cmd.PreRunE(cmd, nil); err != nil {
		t.Fatalf("PreRunE() error = %v", err)
	}

	got, err := cmd.Flags().GetString("format")
	if err != nil {
		t.Fatalf("GetString(format): %v", err)
	}
	if got != "xml" {
		t.Fatalf("expected explicit flag to win, got %q", got)
	}
}

func TestRunCommand_PassesContextToDownloader(t *testing.T) {
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
			old := newDownloaderFunc
			t.Cleanup(func() {
				newDownloaderFunc = old
			})

			md := &mockDownloader{
				downloadURL: "https://example.com/file.zip",
			}
			newDownloaderFunc = func(cfg *globalCfg.GlobalConfig) (downloader, error) {
				return md, nil
			}

			cfg := &globalCfg.GlobalConfig{
				Token:          "token",
				ProjectID:      "project-id",
				ContextTimeout: tt.timeout,
			}
			flags := newFlags()
			flags.Format = "json"

			cmd := newBoundTestCommand(flags)
			var out bytes.Buffer
			cmd.SetOut(&out)
			cmd.SetErr(&out)

			if err := runCommand(cmd, cfg, flags, nil); err != nil {
				t.Fatalf("runCommand() error = %v", err)
			}

			if md.gotCtx == nil {
				t.Fatal("expected context to be passed to downloader")
			}

			_, gotDeadline := md.gotCtx.Deadline()
			if gotDeadline != tt.wantDeadline {
				t.Fatalf("unexpected deadline presence: got %v, want %v", gotDeadline, tt.wantDeadline)
			}
		})
	}
}

func TestRunCommand_UsesDefaultsInBuildParams(t *testing.T) {
	old := newDownloaderFunc
	t.Cleanup(func() {
		newDownloaderFunc = old
	})

	md := &mockDownloader{
		downloadURL: "https://example.com/file.zip",
	}
	newDownloaderFunc = func(cfg *globalCfg.GlobalConfig) (downloader, error) {
		return md, nil
	}

	cfg := &globalCfg.GlobalConfig{
		Token:     "token",
		ProjectID: "project-id",
	}
	flags := newFlags()
	flags.Format = "json"

	defaults := &DownloadConfig{
		Compact: new(true),
	}

	cmd := newBoundTestCommand(flags)
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)

	if err := runCommand(cmd, cfg, flags, defaults); err != nil {
		t.Fatalf("runCommand() error = %v", err)
	}

	got, ok := md.gotParams["compact"]
	if !ok {
		t.Fatal("expected compact to be set from defaults")
	}
	if got != true {
		t.Fatalf("expected compact=true, got %#v", got)
	}
}

func TestRunCommand_AsyncDownloadError(t *testing.T) {
	old := newDownloaderFunc
	t.Cleanup(func() {
		newDownloaderFunc = old
	})

	md := &mockDownloader{
		downloadAsyncErr: errors.New("async download failed"),
	}
	newDownloaderFunc = func(cfg *globalCfg.GlobalConfig) (downloader, error) {
		return md, nil
	}

	cfg := &globalCfg.GlobalConfig{
		Token:     "token",
		ProjectID: "project-id",
	}
	flags := newFlags()
	flags.Format = "json"
	flags.Async = true

	cmd := newBoundTestCommand(flags)
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)

	err := runCommand(cmd, cfg, flags, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "async download failed" {
		t.Fatalf("unexpected error: %q", err.Error())
	}

	if !md.downloadAsyncCalled {
		t.Fatal("expected DownloadAsync to be called")
	}
	if md.downloadCalled {
		t.Fatal("did not expect Download to be called")
	}
	if strings.Contains(out.String(), "Bundle downloaded from:") {
		t.Fatalf("did not expect success output on error, got %q", out.String())
	}
}

func newBoundTestCommand(flags *Flags) *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	bindFlags(cmd, flags)
	return cmd
}
