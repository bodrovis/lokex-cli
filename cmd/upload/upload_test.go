package upload

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bodrovis/lokex-cli/internal/appstate"
	"github.com/bodrovis/lokex-cli/internal/global_config"
	"github.com/bodrovis/lokex-cli/internal/params"
	lokexupload "github.com/bodrovis/lokex/v2/client/upload"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockUploader struct {
	uploadCalled bool
	batchCalled  bool

	gotCtx     context.Context
	gotParams  lokexupload.UploadParams
	gotSrcPath string
	gotPoll    bool

	gotBatchCtx   context.Context
	gotBatchItems []lokexupload.BatchUploadItem
	gotBatchPoll  bool

	result string
	err    error

	batchResult lokexupload.BatchUploadResult
	batchErr    error
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

func (m *mockUploader) UploadBatch(
	ctx context.Context,
	items []lokexupload.BatchUploadItem,
	poll bool,
) (lokexupload.BatchUploadResult, error) {
	m.batchCalled = true
	m.gotBatchCtx = ctx
	m.gotBatchItems = items
	m.gotBatchPoll = poll

	return m.batchResult, m.batchErr
}

func TestNewUploader(t *testing.T) {
	t.Run("returns error when client config is invalid", func(t *testing.T) {
		cfg := &global_config.GlobalConfig{}

		got, err := newUploader(cfg)

		require.Error(t, err)
		require.Nil(t, got)
	})

	t.Run("returns uploader when client config is valid", func(t *testing.T) {
		cfg := &global_config.GlobalConfig{
			Token:     "token",
			ProjectID: "project-id",
		}

		got, err := newUploader(cfg)

		require.NoError(t, err)
		require.NotNil(t, got)
	})
}

func TestNewCommand(t *testing.T) {
	cfg := &global_config.GlobalConfig{}
	state := &appstate.State{}

	cmd := NewCommand(
		cfg,
		state,
	)

	require.NotNil(t, cmd)
	assert.Equal(t, "upload", cmd.Use)
	assert.Equal(t, "Upload translation files to Lokalise", cmd.Short)

	require.NotNil(t, cmd.PreRunE)
	require.NotNil(t, cmd.RunE)
}

func TestValidateCommand(t *testing.T) {
	tests := []struct {
		name      string
		cfg       *global_config.GlobalConfig
		uploadCfg *UploadConfig
		wantErr   string
	}{
		{
			name: "ok",
			cfg: &global_config.GlobalConfig{
				Token:     "token",
				ProjectID: "project-id",
			},
			uploadCfg: &UploadConfig{
				Filename: new("en.json"),
				LangISO:  new("en"),
			},
		},
		{
			name: "ok with manifest only",
			cfg: &global_config.GlobalConfig{
				Token:     "token",
				ProjectID: "project-id",
			},
			uploadCfg: &UploadConfig{
				Manifest: new("manifest.json"),
			},
		},
		{
			name: "ok with manifest and whitespace filename/lang iso",
			cfg: &global_config.GlobalConfig{
				Token:     "token",
				ProjectID: "project-id",
			},
			uploadCfg: &UploadConfig{
				Manifest: new("manifest.json"),
				Filename: new("   "),
				LangISO:  new("   "),
			},
		},
		{
			name: "missing config",
			cfg:  nil,
			uploadCfg: &UploadConfig{
				Filename: new("en.json"),
				LangISO:  new("en"),
			},
			wantErr: "global config is nil",
		},
		{
			name: "missing upload config",
			cfg: &global_config.GlobalConfig{
				Token:     "token",
				ProjectID: "project-id",
			},
			uploadCfg: nil,
			wantErr:   "upload config is nil",
		},
		{
			name: "missing token",
			cfg: &global_config.GlobalConfig{
				ProjectID: "project-id",
			},
			uploadCfg: &UploadConfig{
				Filename: new("en.json"),
				LangISO:  new("en"),
			},
			wantErr: "token is required",
		},
		{
			name: "missing token with manifest",
			cfg: &global_config.GlobalConfig{
				ProjectID: "project-id",
			},
			uploadCfg: &UploadConfig{
				Manifest: new("manifest.json"),
			},
			wantErr: "token is required",
		},
		{
			name: "missing project id",
			cfg: &global_config.GlobalConfig{
				Token: "token",
			},
			uploadCfg: &UploadConfig{
				Filename: new("en.json"),
				LangISO:  new("en"),
			},
			wantErr: "project-id is required",
		},
		{
			name: "missing project id with manifest",
			cfg: &global_config.GlobalConfig{
				Token: "token",
			},
			uploadCfg: &UploadConfig{
				Manifest: new("manifest.json"),
			},
			wantErr: "project-id is required",
		},
		{
			name: "missing filename",
			cfg: &global_config.GlobalConfig{
				Token:     "token",
				ProjectID: "project-id",
			},
			uploadCfg: &UploadConfig{
				LangISO: new("en"),
			},
			wantErr: "--filename is required",
		},
		{
			name: "whitespace filename",
			cfg: &global_config.GlobalConfig{
				Token:     "token",
				ProjectID: "project-id",
			},
			uploadCfg: &UploadConfig{
				Filename: new("   "),
				LangISO:  new("en"),
			},
			wantErr: "--filename is required",
		},
		{
			name: "missing lang iso",
			cfg: &global_config.GlobalConfig{
				Token:     "token",
				ProjectID: "project-id",
			},
			uploadCfg: &UploadConfig{
				Filename: new("en.json"),
			},
			wantErr: "--lang-iso is required",
		},
		{
			name: "whitespace lang iso",
			cfg: &global_config.GlobalConfig{
				Token:     "token",
				ProjectID: "project-id",
			},
			uploadCfg: &UploadConfig{
				Filename: new("en.json"),
				LangISO:  new("   "),
			},
			wantErr: "--lang-iso is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCommand(
				tt.cfg,
				tt.uploadCfg,
			)

			if tt.wantErr == "" {
				require.NoError(t, err)
				return
			}

			require.EqualError(t, err, tt.wantErr)
		})
	}
}

func TestNewCommand_BindsFlags(t *testing.T) {
	cmd := NewCommand(
		&global_config.GlobalConfig{},
		&appstate.State{},
	)

	for _, name := range []string{
		"filename",
		"src-path",
		"lang-iso",
		"poll",
		"manifest",
		"format",
		"tags",
		"filter-task-id",
	} {
		require.NotNilf(
			t,
			cmd.Flags().Lookup(name),
			"expected flag %q to be registered",
			name,
		)
	}
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

		got, err := performUpload(
			context.Background(),
			mu,
			params,
			"./locales/en.json",
			true,
		)

		require.NoError(t, err)
		require.Equal(t, "process-123", got)

		require.True(t, mu.uploadCalled)
		require.Equal(t, "./locales/en.json", mu.gotSrcPath)
		require.True(t, mu.gotPoll)
		require.Equal(t, "en.json", mu.gotParams["filename"])
	})

	t.Run("error", func(t *testing.T) {
		wantErr := errors.New("upload failed")

		mu := &mockUploader{
			err: wantErr,
		}

		_, err := performUpload(
			context.Background(),
			mu,
			params,
			"./locales/en.json",
			false,
		)

		require.ErrorIs(t, err, wantErr)

		require.True(t, mu.uploadCalled)
		require.Equal(t, "./locales/en.json", mu.gotSrcPath)
		require.False(t, mu.gotPoll)
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
			cmd := &cobra.Command{
				Use: "test",
			}

			var out bytes.Buffer
			cmd.SetOut(&out)
			cmd.SetErr(&out)

			printUploadResult(
				cmd,
				tt.result,
				tt.poll,
			)

			require.Equal(t, tt.want, out.String())
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

		newUploaderFunc = func(
			*global_config.GlobalConfig,
		) (uploader, error) {
			return mu, nil
		}

		cfg := &global_config.GlobalConfig{
			Token:     "token",
			ProjectID: "project-id",
		}

		uploadCfg := &UploadConfig{
			Filename: new("en.json"),
			LangISO:  new("en"),
			SrcPath:  new("./locales/en.json"),
			Poll:     new(false),
		}

		cmd := &cobra.Command{Use: "upload"}

		var out bytes.Buffer
		cmd.SetOut(&out)
		cmd.SetErr(&out)

		err := runCommand(
			cmd,
			cfg,
			uploadCfg,
		)
		require.NoError(t, err)

		require.True(t, mu.uploadCalled)
		require.Equal(t, "./locales/en.json", mu.gotSrcPath)
		require.False(t, mu.gotPoll)

		require.Equal(t, "en.json", mu.gotParams["filename"])
		require.Equal(t, "en", mu.gotParams["lang_iso"])

		assert.Contains(
			t,
			out.String(),
			"Upload started: process-123",
		)
	})

	t.Run("happy path with poll", func(t *testing.T) {
		old := newUploaderFunc

		t.Cleanup(func() {
			newUploaderFunc = old
		})

		mu := &mockUploader{
			result: "bundle-456",
		}

		newUploaderFunc = func(
			*global_config.GlobalConfig,
		) (uploader, error) {
			return mu, nil
		}

		cfg := &global_config.GlobalConfig{
			Token:     "token",
			ProjectID: "project-id",
		}

		uploadCfg := &UploadConfig{
			Filename: new("en.json"),
			LangISO:  new("en"),
			SrcPath:  new("./locales/en.json"),
			Poll:     new(true),
		}

		cmd := &cobra.Command{Use: "upload"}

		var out bytes.Buffer
		cmd.SetOut(&out)
		cmd.SetErr(&out)

		err := runCommand(
			cmd,
			cfg,
			uploadCfg,
		)
		require.NoError(t, err)

		require.True(t, mu.uploadCalled)
		require.True(t, mu.gotPoll)

		assert.Contains(
			t,
			out.String(),
			"Upload completed: bundle-456",
		)
	})

	t.Run("build params error", func(t *testing.T) {
		oldUploader := newUploaderFunc
		oldBuildParams := buildParamsFunc

		t.Cleanup(func() {
			newUploaderFunc = oldUploader
			buildParamsFunc = oldBuildParams
		})

		mu := &mockUploader{}

		newUploaderFunc = func(
			*global_config.GlobalConfig,
		) (uploader, error) {
			return mu, nil
		}

		wantErr := errors.New("build params failed")

		buildParamsFunc = func(
			*UploadConfig,
		) (lokexupload.UploadParams, error) {
			return nil, wantErr
		}

		cfg := &global_config.GlobalConfig{
			Token:     "token",
			ProjectID: "project-id",
		}

		uploadCfg := &UploadConfig{
			Filename: new("en.json"),
			LangISO:  new("en"),
			SrcPath:  new("./locales/en.json"),
		}

		cmd := &cobra.Command{Use: "upload"}

		err := runCommand(
			cmd,
			cfg,
			uploadCfg,
		)

		require.ErrorIs(t, err, wantErr)
		require.False(t, mu.uploadCalled)
		require.False(t, mu.batchCalled)
	})

	t.Run("uploader factory error", func(t *testing.T) {
		old := newUploaderFunc

		t.Cleanup(func() {
			newUploaderFunc = old
		})

		wantErr := errors.New("cannot create uploader")

		newUploaderFunc = func(
			*global_config.GlobalConfig,
		) (uploader, error) {
			return nil, wantErr
		}

		cfg := &global_config.GlobalConfig{
			Token:     "token",
			ProjectID: "project-id",
		}

		uploadCfg := &UploadConfig{
			Filename: new("en.json"),
			LangISO:  new("en"),
		}

		cmd := &cobra.Command{Use: "upload"}

		err := runCommand(
			cmd,
			cfg,
			uploadCfg,
		)

		require.ErrorIs(t, err, wantErr)
	})

	t.Run("upload error", func(t *testing.T) {
		old := newUploaderFunc

		t.Cleanup(func() {
			newUploaderFunc = old
		})

		wantErr := errors.New("upload failed")

		mu := &mockUploader{
			err: wantErr,
		}

		newUploaderFunc = func(
			*global_config.GlobalConfig,
		) (uploader, error) {
			return mu, nil
		}

		cfg := &global_config.GlobalConfig{
			Token:     "token",
			ProjectID: "project-id",
		}

		uploadCfg := &UploadConfig{
			Filename: new("en.json"),
			LangISO:  new("en"),
			SrcPath:  new("./locales/en.json"),
		}

		cmd := &cobra.Command{Use: "upload"}

		err := runCommand(
			cmd,
			cfg,
			uploadCfg,
		)

		require.ErrorIs(t, err, wantErr)
		require.True(t, mu.uploadCalled)
	})
}

func TestLoadUploadConfig_UsesConfigValues(t *testing.T) {
	v := viper.New()

	v.Set("upload.filename", "en.json")
	v.Set("upload.lang-iso", "en")

	cmd := &cobra.Command{
		Use: "upload",
	}

	params.BindFlags(
		cmd,
		uploadParamSpecs,
	)

	cfg := &UploadConfig{}

	err := LoadUploadConfig(
		v,
		cmd,
		cfg,
	)
	require.NoError(t, err)

	require.NotNil(t, cfg.Filename)
	require.Equal(t, "en.json", *cfg.Filename)

	require.NotNil(t, cfg.LangISO)
	require.Equal(t, "en", *cfg.LangISO)
}

func TestLoadUploadConfig_ExplicitFlagsOverrideConfig(t *testing.T) {
	v := viper.New()

	v.Set("upload.filename", "default.json")
	v.Set("upload.lang-iso", "fr")

	cmd := &cobra.Command{
		Use: "upload",
	}

	params.BindFlags(
		cmd,
		uploadParamSpecs,
	)

	require.NoError(
		t,
		cmd.ParseFlags([]string{
			"--filename=explicit.json",
			"--lang-iso=en",
		}),
	)

	cfg := &UploadConfig{}

	err := LoadUploadConfig(
		v,
		cmd,
		cfg,
	)
	require.NoError(t, err)

	require.NotNil(t, cfg.Filename)
	require.Equal(t, "explicit.json", *cfg.Filename)

	require.NotNil(t, cfg.LangISO)
	require.Equal(t, "en", *cfg.LangISO)
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

			newUploaderFunc = func(
				*global_config.GlobalConfig,
			) (uploader, error) {
				return mu, nil
			}

			cfg := &global_config.GlobalConfig{
				Token:          "token",
				ProjectID:      "project-id",
				ContextTimeout: tt.timeout,
			}

			uploadCfg := &UploadConfig{
				Filename: new("en.json"),
				LangISO:  new("en"),
			}

			cmd := &cobra.Command{
				Use: "upload",
			}

			var out bytes.Buffer
			cmd.SetOut(&out)
			cmd.SetErr(&out)

			err := runCommand(
				cmd,
				cfg,
				uploadCfg,
			)
			require.NoError(t, err)

			require.NotNil(t, mu.gotCtx)

			_, gotDeadline := mu.gotCtx.Deadline()
			require.Equal(t, tt.wantDeadline, gotDeadline)
		})
	}
}

func TestRunCommand_UsesUploadConfigInBuildParams(t *testing.T) {
	old := newUploaderFunc

	t.Cleanup(func() {
		newUploaderFunc = old
	})

	mu := &mockUploader{
		result: "process-123",
	}

	newUploaderFunc = func(
		*global_config.GlobalConfig,
	) (uploader, error) {
		return mu, nil
	}

	cfg := &global_config.GlobalConfig{
		Token:     "token",
		ProjectID: "project-id",
	}

	uploadCfg := &UploadConfig{
		Filename: new("en.json"),
		LangISO:  new("en"),
		ApplyTM:  new(true),
	}

	cmd := &cobra.Command{
		Use: "upload",
	}

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)

	err := runCommand(
		cmd,
		cfg,
		uploadCfg,
	)
	require.NoError(t, err)

	got, ok := mu.gotParams["apply_tm"]

	require.True(t, ok)
	require.Equal(t, true, got)
}

func TestRunCommand_UploadErrorDoesNotPrintSuccessOutput(t *testing.T) {
	old := newUploaderFunc

	t.Cleanup(func() {
		newUploaderFunc = old
	})

	wantErr := errors.New("upload failed")

	mu := &mockUploader{
		err: wantErr,
	}

	newUploaderFunc = func(
		*global_config.GlobalConfig,
	) (uploader, error) {
		return mu, nil
	}

	cfg := &global_config.GlobalConfig{
		Token:     "token",
		ProjectID: "project-id",
	}

	uploadCfg := &UploadConfig{
		Filename: new("en.json"),
		LangISO:  new("en"),
	}

	cmd := &cobra.Command{
		Use: "upload",
	}

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)

	err := runCommand(
		cmd,
		cfg,
		uploadCfg,
	)

	require.ErrorIs(t, err, wantErr)

	assert.NotContains(t, out.String(), "Upload started:")
	assert.NotContains(t, out.String(), "Upload completed:")
}

func TestNewCommand_Execute_UsesConfigValues(t *testing.T) {
	old := newUploaderFunc

	t.Cleanup(func() {
		newUploaderFunc = old
	})

	mu := &mockUploader{
		result: "bundle-456",
	}

	newUploaderFunc = func(
		*global_config.GlobalConfig,
	) (uploader, error) {
		return mu, nil
	}

	cfg := &global_config.GlobalConfig{
		Token:     "token",
		ProjectID: "project-id",
	}

	state := newUploadTestState(map[string]any{
		"filename": "en.json",
		"lang-iso": "en",
		"src-path": "./locales/from-config.json",
		"poll":     true,
	})

	cmd := NewCommand(
		cfg,
		state,
	)

	cmd.SetArgs(nil)

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)

	err := cmd.Execute()
	require.NoError(t, err)

	require.True(t, mu.uploadCalled)

	require.Equal(
		t,
		"./locales/from-config.json",
		mu.gotSrcPath,
	)

	require.True(t, mu.gotPoll)

	require.Equal(
		t,
		"en.json",
		mu.gotParams["filename"],
	)

	require.Equal(
		t,
		"en",
		mu.gotParams["lang_iso"],
	)

	assert.Contains(
		t,
		out.String(),
		"Upload completed: bundle-456",
	)
}

func TestNewCommand_Execute_ExplicitFlagsOverrideConfig(t *testing.T) {
	old := newUploaderFunc

	t.Cleanup(func() {
		newUploaderFunc = old
	})

	mu := &mockUploader{
		result: "process-123",
	}

	newUploaderFunc = func(
		*global_config.GlobalConfig,
	) (uploader, error) {
		return mu, nil
	}

	cfg := &global_config.GlobalConfig{
		Token:     "token",
		ProjectID: "project-id",
	}

	state := newUploadTestState(map[string]any{
		"filename": "default.json",
		"lang-iso": "fr",
		"src-path": "./locales/from-config.json",
		"poll":     true,
	})

	cmd := NewCommand(
		cfg,
		state,
	)

	cmd.SetArgs([]string{
		"--filename=explicit.json",
		"--lang-iso=en",
		"--src-path=./locales/explicit.json",
		"--poll=false",
	})

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)

	err := cmd.Execute()
	require.NoError(t, err)

	require.True(t, mu.uploadCalled)

	require.Equal(
		t,
		"./locales/explicit.json",
		mu.gotSrcPath,
	)

	require.False(t, mu.gotPoll)

	require.Equal(
		t,
		"explicit.json",
		mu.gotParams["filename"],
	)

	require.Equal(
		t,
		"en",
		mu.gotParams["lang_iso"],
	)

	assert.Contains(
		t,
		out.String(),
		"Upload started: process-123",
	)
}

func newUploadTestState(values map[string]any) *appstate.State {
	v := viper.New()

	for key, value := range values {
		v.Set("upload."+key, value)
	}

	return &appstate.State{
		Viper: v,
	}
}
