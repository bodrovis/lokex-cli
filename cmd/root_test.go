package cmd

import (
	"errors"
	"testing"

	"github.com/spf13/cobra"

	downloadcmd "github.com/bodrovis/lokex-cli/cmd/download"
	uploadcmd "github.com/bodrovis/lokex-cli/cmd/upload"
	"github.com/bodrovis/lokex-cli/internal/global_config"
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

func TestNewPersistentPreRunE_LoadsUploadConfigForUploadCommand(t *testing.T) {
	cfg := &global_config.GlobalConfig{UserAgent: "test-agent"}
	uploadCfg := &uploadcmd.UploadConfig{}
	downloadCfg := &downloadcmd.DownloadConfig{}
	configFile := "test.yaml"

	globalCalled := false
	uploadCalled := false
	downloadCalled := false

	preRun := newPersistentPreRunE(
		cfg,
		uploadCfg,
		downloadCfg,
		&configFile,
		func(userAgent string, opts global_config.LoadOptions) (*global_config.GlobalConfigInput, error) {
			globalCalled = true

			if userAgent != "test-agent" {
				t.Fatalf("unexpected user agent: got %q, want %q", userAgent, "test-agent")
			}
			if opts.ConfigFile != "test.yaml" {
				t.Fatalf("unexpected config file: got %q, want %q", opts.ConfigFile, "test.yaml")
			}
			if opts.EnvPrefix != "LOKEX" {
				t.Fatalf("unexpected env prefix: got %q, want %q", opts.EnvPrefix, "LOKEX")
			}

			return &global_config.GlobalConfigInput{}, nil
		},
		func(cfg *uploadcmd.UploadConfig, configFile string, envPrefix string) error {
			uploadCalled = true

			if configFile != "test.yaml" {
				t.Fatalf("unexpected config file: got %q, want %q", configFile, "test.yaml")
			}
			if envPrefix != "LOKEX" {
				t.Fatalf("unexpected env prefix: got %q, want %q", envPrefix, "LOKEX")
			}

			return nil
		},
		func(cfg *downloadcmd.DownloadConfig, configFile string, envPrefix string) error {
			downloadCalled = true
			return nil
		},
	)

	cmd := &cobra.Command{Use: "upload"}

	if err := preRun(cmd, nil); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !globalCalled {
		t.Fatal("expected global loader to be called")
	}
	if !uploadCalled {
		t.Fatal("expected upload loader to be called")
	}
	if downloadCalled {
		t.Fatal("did not expect download loader to be called")
	}
}

func TestNewPersistentPreRunE_LoadsDownloadConfigForDownloadCommand(t *testing.T) {
	cfg := &global_config.GlobalConfig{UserAgent: "test-agent"}
	uploadCfg := &uploadcmd.UploadConfig{}
	downloadCfg := &downloadcmd.DownloadConfig{}
	configFile := "test.yaml"

	globalCalled := false
	uploadCalled := false
	downloadCalled := false

	preRun := newPersistentPreRunE(
		cfg,
		uploadCfg,
		downloadCfg,
		&configFile,
		func(userAgent string, opts global_config.LoadOptions) (*global_config.GlobalConfigInput, error) {
			globalCalled = true

			if userAgent != "test-agent" {
				t.Fatalf("unexpected user agent: got %q, want %q", userAgent, "test-agent")
			}
			if opts.ConfigFile != "test.yaml" {
				t.Fatalf("unexpected config file: got %q, want %q", opts.ConfigFile, "test.yaml")
			}
			if opts.EnvPrefix != "LOKEX" {
				t.Fatalf("unexpected env prefix: got %q, want %q", opts.EnvPrefix, "LOKEX")
			}

			return &global_config.GlobalConfigInput{}, nil
		},
		func(cfg *uploadcmd.UploadConfig, configFile string, envPrefix string) error {
			uploadCalled = true
			return nil
		},
		func(cfg *downloadcmd.DownloadConfig, configFile string, envPrefix string) error {
			downloadCalled = true

			if configFile != "test.yaml" {
				t.Fatalf("unexpected config file: got %q, want %q", configFile, "test.yaml")
			}
			if envPrefix != "LOKEX" {
				t.Fatalf("unexpected env prefix: got %q, want %q", envPrefix, "LOKEX")
			}

			return nil
		},
	)

	cmd := &cobra.Command{Use: "download"}

	if err := preRun(cmd, nil); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !globalCalled {
		t.Fatal("expected global loader to be called")
	}
	if uploadCalled {
		t.Fatal("did not expect upload loader to be called")
	}
	if !downloadCalled {
		t.Fatal("expected download loader to be called")
	}
}

func TestNewPersistentPreRunE_DoesNotLoadSubcommandConfigForOtherCommands(t *testing.T) {
	cfg := &global_config.GlobalConfig{UserAgent: "test-agent"}
	uploadCfg := &uploadcmd.UploadConfig{}
	downloadCfg := &downloadcmd.DownloadConfig{}
	configFile := "test.yaml"

	globalCalled := false
	uploadCalled := false
	downloadCalled := false

	preRun := newPersistentPreRunE(
		cfg,
		uploadCfg,
		downloadCfg,
		&configFile,
		func(userAgent string, opts global_config.LoadOptions) (*global_config.GlobalConfigInput, error) {
			globalCalled = true
			return &global_config.GlobalConfigInput{}, nil
		},
		func(cfg *uploadcmd.UploadConfig, configFile string, envPrefix string) error {
			uploadCalled = true
			return nil
		},
		func(cfg *downloadcmd.DownloadConfig, configFile string, envPrefix string) error {
			downloadCalled = true
			return nil
		},
	)

	cmd := &cobra.Command{Use: "version"}

	if err := preRun(cmd, nil); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !globalCalled {
		t.Fatal("expected global loader to be called")
	}
	if uploadCalled {
		t.Fatal("did not expect upload loader to be called")
	}
	if downloadCalled {
		t.Fatal("did not expect download loader to be called")
	}
}

func TestNewPersistentPreRunE_ReturnsGlobalLoadError(t *testing.T) {
	cfg := &global_config.GlobalConfig{UserAgent: "test-agent"}
	uploadCfg := &uploadcmd.UploadConfig{}
	downloadCfg := &downloadcmd.DownloadConfig{}
	configFile := "test.yaml"

	wantErr := errors.New("global load failed")

	preRun := newPersistentPreRunE(
		cfg,
		uploadCfg,
		downloadCfg,
		&configFile,
		func(userAgent string, opts global_config.LoadOptions) (*global_config.GlobalConfigInput, error) {
			return nil, wantErr
		},
		func(cfg *uploadcmd.UploadConfig, configFile string, envPrefix string) error {
			t.Fatal("upload loader should not be called when global loader fails")
			return nil
		},
		func(cfg *downloadcmd.DownloadConfig, configFile string, envPrefix string) error {
			t.Fatal("download loader should not be called when global loader fails")
			return nil
		},
	)

	cmd := &cobra.Command{Use: "upload"}

	err := preRun(cmd, nil)
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected error %v, got %v", wantErr, err)
	}
}

func TestNewPersistentPreRunE_ReturnsUploadLoadError(t *testing.T) {
	cfg := &global_config.GlobalConfig{UserAgent: "test-agent"}
	uploadCfg := &uploadcmd.UploadConfig{}
	downloadCfg := &downloadcmd.DownloadConfig{}
	configFile := "test.yaml"

	wantErr := errors.New("upload load failed")

	preRun := newPersistentPreRunE(
		cfg,
		uploadCfg,
		downloadCfg,
		&configFile,
		func(userAgent string, opts global_config.LoadOptions) (*global_config.GlobalConfigInput, error) {
			return &global_config.GlobalConfigInput{}, nil
		},
		func(cfg *uploadcmd.UploadConfig, configFile string, envPrefix string) error {
			return wantErr
		},
		func(cfg *downloadcmd.DownloadConfig, configFile string, envPrefix string) error {
			t.Fatal("download loader should not be called for upload command")
			return nil
		},
	)

	cmd := &cobra.Command{Use: "upload"}

	err := preRun(cmd, nil)
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected error %v, got %v", wantErr, err)
	}
}

func TestNewPersistentPreRunE_ReturnsDownloadLoadError(t *testing.T) {
	cfg := &global_config.GlobalConfig{UserAgent: "test-agent"}
	uploadCfg := &uploadcmd.UploadConfig{}
	downloadCfg := &downloadcmd.DownloadConfig{}
	configFile := "test.yaml"

	wantErr := errors.New("download load failed")

	preRun := newPersistentPreRunE(
		cfg,
		uploadCfg,
		downloadCfg,
		&configFile,
		func(userAgent string, opts global_config.LoadOptions) (*global_config.GlobalConfigInput, error) {
			return &global_config.GlobalConfigInput{}, nil
		},
		func(cfg *uploadcmd.UploadConfig, configFile string, envPrefix string) error {
			t.Fatal("upload loader should not be called for download command")
			return nil
		},
		func(cfg *downloadcmd.DownloadConfig, configFile string, envPrefix string) error {
			return wantErr
		},
	)

	cmd := &cobra.Command{Use: "download"}

	err := preRun(cmd, nil)
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected error %v, got %v", wantErr, err)
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
