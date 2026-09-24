package download

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bodrovis/lokex-cli/internal/appstate"
	globalCfg "github.com/bodrovis/lokex-cli/internal/global_config"
	"github.com/bodrovis/lokex-cli/internal/params"
	lokexdownload "github.com/bodrovis/lokex/v2/client/download"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func (m *mockDownloader) Download(
	ctx context.Context,
	out string,
	params lokexdownload.DownloadParams,
) (string, error) {
	m.downloadCalled = true
	m.gotCtx = ctx
	m.gotOut = out
	m.gotParams = params

	return m.downloadURL, m.downloadErr
}

func (m *mockDownloader) DownloadAsync(
	ctx context.Context,
	out string,
	params lokexdownload.DownloadParams,
) (string, error) {
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

		require.Error(t, err)
		require.Nil(t, got)
	})

	t.Run("returns downloader when client config is valid", func(t *testing.T) {
		cfg := &globalCfg.GlobalConfig{
			Token:     "token",
			ProjectID: "project-id",
		}

		got, err := newDownloader(cfg)

		require.NoError(t, err)
		require.NotNil(t, got)
	})
}

func TestNewCommand(t *testing.T) {
	cfg := &globalCfg.GlobalConfig{}
	state := &appstate.State{}

	cmd := NewCommand(
		cfg,
		state,
	)

	require.NotNil(t, cmd)
	assert.Equal(t, "download", cmd.Use)
	assert.Equal(
		t,
		"Download translation files from Lokalise",
		cmd.Short,
	)

	require.NotNil(t, cmd.PreRunE)
	require.NotNil(t, cmd.RunE)
}

func TestNewCommand_BindsFlags(t *testing.T) {
	cmd := NewCommand(
		&globalCfg.GlobalConfig{},
		&appstate.State{},
	)

	for _, name := range []string{
		"out",
		"format",
		"async",
		"filter-langs",
		"bundle-structure",
		"filter-task-id",
		"language-mapping",
	} {
		require.NotNilf(
			t,
			cmd.Flags().Lookup(name),
			"expected flag %q to be registered",
			name,
		)
	}
}

func TestValidateCommand(t *testing.T) {
	tests := []struct {
		name        string
		cfg         *globalCfg.GlobalConfig
		downloadCfg *DownloadConfig
		wantErr     string
	}{
		{
			name: "ok",
			cfg: &globalCfg.GlobalConfig{
				Token:     "token",
				ProjectID: "project-id",
			},
			downloadCfg: &DownloadConfig{
				Format: new("json"),
			},
		},
		{
			name: "missing cfg",
			cfg:  nil,
			downloadCfg: &DownloadConfig{
				Format: new("json"),
			},
			wantErr: "global config is nil",
		},
		{
			name: "missing download config",
			cfg: &globalCfg.GlobalConfig{
				Token:     "token",
				ProjectID: "project-id",
			},
			downloadCfg: nil,
			wantErr:     "download config is nil",
		},
		{
			name: "missing token",
			cfg: &globalCfg.GlobalConfig{
				ProjectID: "project-id",
			},
			downloadCfg: &DownloadConfig{
				Format: new("json"),
			},
			wantErr: "token is required",
		},
		{
			name: "missing project id",
			cfg: &globalCfg.GlobalConfig{
				Token: "token",
			},
			downloadCfg: &DownloadConfig{
				Format: new("json"),
			},
			wantErr: "project-id is required",
		},
		{
			name: "missing format",
			cfg: &globalCfg.GlobalConfig{
				Token:     "token",
				ProjectID: "project-id",
			},
			downloadCfg: &DownloadConfig{},
			wantErr:     "format is required",
		},
		{
			name: "whitespace format",
			cfg: &globalCfg.GlobalConfig{
				Token:     "token",
				ProjectID: "project-id",
			},
			downloadCfg: &DownloadConfig{
				Format: new("   "),
			},
			wantErr: "format is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCommand(
				tt.cfg,
				tt.downloadCfg,
			)

			if tt.wantErr == "" {
				require.NoError(t, err)
				return
			}

			require.EqualError(t, err, tt.wantErr)
		})
	}
}

func TestPerformDownload(t *testing.T) {
	params := lokexdownload.DownloadParams{
		"format": "json",
	}

	t.Run("sync download", func(t *testing.T) {
		md := &mockDownloader{
			downloadURL: "https://example.com/sync.zip",
		}

		got, err := performDownload(
			context.Background(),
			md,
			params,
			"./locales",
			false,
		)

		require.NoError(t, err)
		require.Equal(t, "https://example.com/sync.zip", got)

		require.True(t, md.downloadCalled)
		require.False(t, md.downloadAsyncCalled)
		require.Equal(t, "./locales", md.gotOut)
		require.Equal(t, "json", md.gotParams["format"])
	})

	t.Run("async download", func(t *testing.T) {
		md := &mockDownloader{
			downloadAsyncURL: "https://example.com/async.zip",
		}

		got, err := performDownload(
			context.Background(),
			md,
			params,
			"./locales",
			true,
		)

		require.NoError(t, err)
		require.Equal(t, "https://example.com/async.zip", got)

		require.True(t, md.downloadAsyncCalled)
		require.False(t, md.downloadCalled)
		require.Equal(t, "./locales", md.gotOut)
		require.Equal(t, "json", md.gotParams["format"])
	})

	t.Run("sync error", func(t *testing.T) {
		wantErr := errors.New("sync failed")

		md := &mockDownloader{
			downloadErr: wantErr,
		}

		_, err := performDownload(
			context.Background(),
			md,
			params,
			"./locales",
			false,
		)

		require.ErrorIs(t, err, wantErr)
		require.True(t, md.downloadCalled)
		require.False(t, md.downloadAsyncCalled)
	})

	t.Run("async error", func(t *testing.T) {
		wantErr := errors.New("async failed")

		md := &mockDownloader{
			downloadAsyncErr: wantErr,
		}

		_, err := performDownload(
			context.Background(),
			md,
			params,
			"./locales",
			true,
		)

		require.ErrorIs(t, err, wantErr)
		require.True(t, md.downloadAsyncCalled)
		require.False(t, md.downloadCalled)
	})
}

func TestPrintDownloadResult(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)

	printDownloadResult(
		cmd,
		"https://example.com/file.zip",
	)

	require.Equal(
		t,
		"Bundle downloaded from: https://example.com/file.zip\n",
		out.String(),
	)
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
			require.Equal(
				t,
				tt.want,
				truncateURLForOutput(tt.url, tt.max),
			)
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

	newDownloaderFunc = func(
		*globalCfg.GlobalConfig,
	) (downloader, error) {
		return md, nil
	}

	cfg := &globalCfg.GlobalConfig{
		Token:     "token",
		ProjectID: "project-id",
	}

	state := &appstate.State{
		Viper: viper.New(),
	}

	cmd := NewCommand(
		cfg,
		state,
	)

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)

	cmd.SetArgs([]string{
		"--format=json",
		"--out=./locales",
	})

	err := cmd.Execute()
	require.NoError(t, err)

	require.True(t, md.downloadCalled)
	require.False(t, md.downloadAsyncCalled)

	require.Equal(
		t,
		"./locales",
		md.gotOut,
	)

	require.Equal(
		t,
		"json",
		md.gotParams["format"],
	)
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

		newDownloaderFunc = func(
			*globalCfg.GlobalConfig,
		) (downloader, error) {
			return md, nil
		}

		cfg := &globalCfg.GlobalConfig{
			Token:     "token",
			ProjectID: "project-id",
		}

		downloadCfg := &DownloadConfig{
			Out:    new("./locales"),
			Format: new("json"),
			Async:  new(false),
		}

		cmd := &cobra.Command{
			Use: "download",
		}

		var out bytes.Buffer
		cmd.SetOut(&out)
		cmd.SetErr(&out)

		err := runCommand(
			cmd,
			cfg,
			downloadCfg,
		)
		require.NoError(t, err)

		require.True(t, md.downloadCalled)
		require.False(t, md.downloadAsyncCalled)

		require.Equal(
			t,
			"./locales",
			md.gotOut,
		)

		require.Equal(
			t,
			"json",
			md.gotParams["format"],
		)

		assert.Contains(
			t,
			out.String(),
			"Bundle downloaded from: https://example.com/file.zip",
		)
	})

	t.Run("async happy path", func(t *testing.T) {
		old := newDownloaderFunc

		t.Cleanup(func() {
			newDownloaderFunc = old
		})

		md := &mockDownloader{
			downloadAsyncURL: "https://example.com/async.zip",
		}

		newDownloaderFunc = func(
			*globalCfg.GlobalConfig,
		) (downloader, error) {
			return md, nil
		}

		cfg := &globalCfg.GlobalConfig{
			Token:     "token",
			ProjectID: "project-id",
		}

		downloadCfg := &DownloadConfig{
			Out:    new("./locales"),
			Format: new("json"),
			Async:  new(true),
		}

		cmd := &cobra.Command{
			Use: "download",
		}

		var out bytes.Buffer
		cmd.SetOut(&out)
		cmd.SetErr(&out)

		err := runCommand(
			cmd,
			cfg,
			downloadCfg,
		)
		require.NoError(t, err)

		require.True(t, md.downloadAsyncCalled)
		require.False(t, md.downloadCalled)

		require.Equal(
			t,
			"./locales",
			md.gotOut,
		)

		require.Equal(
			t,
			"json",
			md.gotParams["format"],
		)
	})

	t.Run("downloader factory error", func(t *testing.T) {
		old := newDownloaderFunc

		t.Cleanup(func() {
			newDownloaderFunc = old
		})

		wantErr := errors.New("cannot create downloader")

		newDownloaderFunc = func(
			*globalCfg.GlobalConfig,
		) (downloader, error) {
			return nil, wantErr
		}

		cfg := &globalCfg.GlobalConfig{
			Token:     "token",
			ProjectID: "project-id",
		}

		downloadCfg := &DownloadConfig{
			Format: new("json"),
		}

		cmd := &cobra.Command{
			Use: "download",
		}

		err := runCommand(
			cmd,
			cfg,
			downloadCfg,
		)

		require.ErrorIs(t, err, wantErr)
	})

	t.Run("build params error", func(t *testing.T) {
		old := newDownloaderFunc

		t.Cleanup(func() {
			newDownloaderFunc = old
		})

		md := &mockDownloader{}

		newDownloaderFunc = func(
			*globalCfg.GlobalConfig,
		) (downloader, error) {
			return md, nil
		}

		cfg := &globalCfg.GlobalConfig{
			Token:     "token",
			ProjectID: "project-id",
		}

		downloadCfg := &DownloadConfig{
			Format:              new("json"),
			LanguageMappingJSON: new("{"),
		}

		cmd := &cobra.Command{
			Use: "download",
		}

		err := runCommand(
			cmd,
			cfg,
			downloadCfg,
		)

		require.Error(t, err)
		assert.Contains(
			t,
			err.Error(),
			"parse language-mapping",
		)

		require.False(t, md.downloadCalled)
		require.False(t, md.downloadAsyncCalled)
	})

	t.Run("download error", func(t *testing.T) {
		old := newDownloaderFunc

		t.Cleanup(func() {
			newDownloaderFunc = old
		})

		wantErr := errors.New("download failed")

		md := &mockDownloader{
			downloadErr: wantErr,
		}

		newDownloaderFunc = func(
			*globalCfg.GlobalConfig,
		) (downloader, error) {
			return md, nil
		}

		cfg := &globalCfg.GlobalConfig{
			Token:     "token",
			ProjectID: "project-id",
		}

		downloadCfg := &DownloadConfig{
			Out:    new("./locales"),
			Format: new("json"),
		}

		cmd := &cobra.Command{
			Use: "download",
		}

		err := runCommand(
			cmd,
			cfg,
			downloadCfg,
		)

		require.ErrorIs(t, err, wantErr)
		require.True(t, md.downloadCalled)
		require.False(t, md.downloadAsyncCalled)
	})
}

func TestLoadDownloadConfig_UsesConfigValues(t *testing.T) {
	v := viper.New()
	v.Set("download.format", "json")

	cmd := &cobra.Command{Use: "download"}
	params.BindFlags(cmd, downloadParamSpecs)

	cfg := &DownloadConfig{}

	err := LoadDownloadConfig(
		v,
		cmd,
		cfg,
	)
	require.NoError(t, err)

	require.NotNil(t, cfg.Format)
	require.Equal(t, "json", *cfg.Format)
}

func TestLoadDownloadConfig_ExplicitFlagOverridesConfig(t *testing.T) {
	v := viper.New()
	v.Set("download.format", "json")

	cmd := &cobra.Command{Use: "download"}
	params.BindFlags(cmd, downloadParamSpecs)

	require.NoError(
		t,
		cmd.ParseFlags([]string{
			"--format=xml",
		}),
	)

	cfg := &DownloadConfig{}

	err := LoadDownloadConfig(
		v,
		cmd,
		cfg,
	)
	require.NoError(t, err)

	require.NotNil(t, cfg.Format)
	require.Equal(t, "xml", *cfg.Format)
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

			newDownloaderFunc = func(
				*globalCfg.GlobalConfig,
			) (downloader, error) {
				return md, nil
			}

			cfg := &globalCfg.GlobalConfig{
				Token:          "token",
				ProjectID:      "project-id",
				ContextTimeout: tt.timeout,
			}

			downloadCfg := &DownloadConfig{
				Format: new("json"),
			}

			cmd := &cobra.Command{
				Use: "download",
			}

			var out bytes.Buffer
			cmd.SetOut(&out)
			cmd.SetErr(&out)

			err := runCommand(
				cmd,
				cfg,
				downloadCfg,
			)
			require.NoError(t, err)

			require.NotNil(t, md.gotCtx)

			_, gotDeadline := md.gotCtx.Deadline()

			require.Equal(
				t,
				tt.wantDeadline,
				gotDeadline,
			)
		})
	}
}

func TestRunCommand_UsesDownloadConfigInBuildParams(t *testing.T) {
	old := newDownloaderFunc

	t.Cleanup(func() {
		newDownloaderFunc = old
	})

	md := &mockDownloader{
		downloadURL: "https://example.com/file.zip",
	}

	newDownloaderFunc = func(
		*globalCfg.GlobalConfig,
	) (downloader, error) {
		return md, nil
	}

	cfg := &globalCfg.GlobalConfig{
		Token:     "token",
		ProjectID: "project-id",
	}

	downloadCfg := &DownloadConfig{
		Out:     new("./locales"),
		Format:  new("json"),
		Compact: new(true),
	}

	cmd := &cobra.Command{
		Use: "download",
	}

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)

	err := runCommand(
		cmd,
		cfg,
		downloadCfg,
	)
	require.NoError(t, err)

	got, ok := md.gotParams["compact"]

	require.True(t, ok)
	require.Equal(t, true, got)
}

func TestRunCommand_AsyncDownloadError(t *testing.T) {
	old := newDownloaderFunc

	t.Cleanup(func() {
		newDownloaderFunc = old
	})

	wantErr := errors.New("async download failed")

	md := &mockDownloader{
		downloadAsyncErr: wantErr,
	}

	newDownloaderFunc = func(
		*globalCfg.GlobalConfig,
	) (downloader, error) {
		return md, nil
	}

	cfg := &globalCfg.GlobalConfig{
		Token:     "token",
		ProjectID: "project-id",
	}

	downloadCfg := &DownloadConfig{
		Out:    new("./locales"),
		Format: new("json"),
		Async:  new(true),
	}

	cmd := &cobra.Command{
		Use: "download",
	}

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)

	err := runCommand(
		cmd,
		cfg,
		downloadCfg,
	)

	require.ErrorIs(t, err, wantErr)

	require.True(t, md.downloadAsyncCalled)
	require.False(t, md.downloadCalled)

	assert.NotContains(
		t,
		out.String(),
		"Bundle downloaded from:",
	)
}
