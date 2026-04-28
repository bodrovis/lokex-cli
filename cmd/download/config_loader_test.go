package download

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestLoadDownloadConfig_IgnoresMissingImplicitConfig(t *testing.T) {
	cfg := &DownloadConfig{}

	err := LoadDownloadConfig(cfg, "", "LOKEX")
	if err != nil {
		t.Fatalf("expected no error for missing implicit config, got %v", err)
	}
}

func TestLoadDownloadConfig_ReturnsErrorForMissingExplicitConfig(t *testing.T) {
	cfg := &DownloadConfig{}

	err := LoadDownloadConfig(cfg, filepath.Join(t.TempDir(), "missing.yaml"), "LOKEX")
	if err == nil {
		t.Fatal("expected error for missing explicit config")
	}
}

func TestLoadDownloadConfig_ReturnsErrorWhenConfigIsNil(t *testing.T) {
	err := LoadDownloadConfig(nil, "", "LOKEX")
	if err == nil {
		t.Fatal("expected error for nil download config")
	}

	if !strings.Contains(err.Error(), "download config is nil") {
		t.Fatalf("expected nil config error, got %v", err)
	}
}

func TestLoadDownloadConfig_ReturnsErrorForInvalidExplicitConfig(t *testing.T) {
	configFile := filepath.Join(t.TempDir(), "broken.yaml")
	if err := os.WriteFile(configFile, []byte(":\nbad yaml"), 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	cfg := &DownloadConfig{}

	err := LoadDownloadConfig(cfg, configFile, "LOKEX")
	if err == nil {
		t.Fatal("expected error for invalid explicit config")
	}
}

func TestLoadDownloadConfig_LoadsValuesFromConfigFile(t *testing.T) {
	configFile := filepath.Join(t.TempDir(), "lokex.yaml")
	content := `
download:
  out: ./tmp/out
  format: json
  async: true
  original-filenames: true
  filter-langs:
    - en
    - fr
  filter-task-id: 123
`
	if err := os.WriteFile(configFile, []byte(strings.TrimSpace(content)), 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	cfg := &DownloadConfig{}

	err := LoadDownloadConfig(cfg, configFile, "LOKEX")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Out == nil || *cfg.Out != "./tmp/out" {
		t.Fatalf("expected Out to be loaded from config, got %#v", cfg.Out)
	}

	if cfg.Format == nil || *cfg.Format != "json" {
		t.Fatalf("expected Format to be loaded from config, got %#v", cfg.Format)
	}

	if cfg.Async == nil || *cfg.Async != true {
		t.Fatalf("expected Async to be true, got %#v", cfg.Async)
	}

	if cfg.OriginalFilenames == nil || *cfg.OriginalFilenames != true {
		t.Fatalf("expected OriginalFilenames to be true, got %#v", cfg.OriginalFilenames)
	}

	wantLangs := []string{"en", "fr"}
	if cfg.FilterLangs == nil || !reflect.DeepEqual(*cfg.FilterLangs, wantLangs) {
		t.Fatalf("expected FilterLangs %v, got %#v", wantLangs, cfg.FilterLangs)
	}

	if cfg.FilterTaskID == nil || *cfg.FilterTaskID != 123 {
		t.Fatalf("expected FilterTaskID to be 123, got %#v", cfg.FilterTaskID)
	}
}

func TestLoadDownloadConfig_LoadsValuesFromEnv(t *testing.T) {
	t.Setenv("LOKEX_DOWNLOAD_OUT", "./env-out")
	t.Setenv("LOKEX_DOWNLOAD_FORMAT", "xml")
	t.Setenv("LOKEX_DOWNLOAD_ASYNC", "true")
	t.Setenv("LOKEX_DOWNLOAD_ORIGINAL_FILENAMES", "true")
	t.Setenv("LOKEX_DOWNLOAD_FILTER_LANGS", "de, es")
	t.Setenv("LOKEX_DOWNLOAD_FILTER_TASK_ID", "987")

	cfg := &DownloadConfig{}

	err := LoadDownloadConfig(cfg, "", "LOKEX")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Out == nil || *cfg.Out != "./env-out" {
		t.Fatalf("expected Out to be loaded from env, got %#v", cfg.Out)
	}

	if cfg.Format == nil || *cfg.Format != "xml" {
		t.Fatalf("expected Format to be loaded from env, got %#v", cfg.Format)
	}

	if cfg.Async == nil || *cfg.Async != true {
		t.Fatalf("expected Async to be true, got %#v", cfg.Async)
	}

	if cfg.OriginalFilenames == nil || *cfg.OriginalFilenames != true {
		t.Fatalf("expected OriginalFilenames to be true, got %#v", cfg.OriginalFilenames)
	}

	wantLangs := []string{"de", "es"}
	if cfg.FilterLangs == nil || !reflect.DeepEqual(*cfg.FilterLangs, wantLangs) {
		t.Fatalf("expected FilterLangs %v, got %#v", wantLangs, cfg.FilterLangs)
	}

	if cfg.FilterTaskID == nil || *cfg.FilterTaskID != 987 {
		t.Fatalf("expected FilterTaskID to be 987, got %#v", cfg.FilterTaskID)
	}
}

func TestLoadDownloadConfig_EnvOverridesConfigFile(t *testing.T) {
	configFile := filepath.Join(t.TempDir(), "lokex.yaml")
	content := `
download:
  format: json
  async: false
  filter-langs:
    - en
    - fr
  filter-task-id: 111
`
	if err := os.WriteFile(configFile, []byte(strings.TrimSpace(content)), 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	t.Setenv("LOKEX_DOWNLOAD_FORMAT", "xml")
	t.Setenv("LOKEX_DOWNLOAD_ASYNC", "true")
	t.Setenv("LOKEX_DOWNLOAD_FILTER_LANGS", "de,es")
	t.Setenv("LOKEX_DOWNLOAD_FILTER_TASK_ID", "222")

	cfg := &DownloadConfig{}

	err := LoadDownloadConfig(cfg, configFile, "LOKEX")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Format == nil || *cfg.Format != "xml" {
		t.Fatalf("expected env to override config for Format, got %#v", cfg.Format)
	}

	if cfg.Async == nil || *cfg.Async != true {
		t.Fatalf("expected env to override config for Async, got %#v", cfg.Async)
	}

	wantLangs := []string{"de", "es"}
	if cfg.FilterLangs == nil || !reflect.DeepEqual(*cfg.FilterLangs, wantLangs) {
		t.Fatalf("expected env to override config for FilterLangs, got %#v", cfg.FilterLangs)
	}

	if cfg.FilterTaskID == nil || *cfg.FilterTaskID != 222 {
		t.Fatalf("expected env to override config for FilterTaskID, got %#v", cfg.FilterTaskID)
	}
}

func TestLoadDownloadConfig_EmptyStringSliceEnvDoesNotSetField(t *testing.T) {
	t.Setenv("LOKEX_DOWNLOAD_FILTER_LANGS", " , ,  , ")

	cfg := &DownloadConfig{}

	err := LoadDownloadConfig(cfg, "", "LOKEX")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.FilterLangs != nil {
		t.Fatalf("expected FilterLangs to stay nil, got %#v", cfg.FilterLangs)
	}
}
