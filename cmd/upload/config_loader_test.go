package upload

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestLoadUploadConfig_IgnoresMissingImplicitConfig(t *testing.T) {
	cfg := &UploadConfig{}

	err := LoadUploadConfig(cfg, "", "LOKEX")
	if err != nil {
		t.Fatalf("expected no error for missing implicit config, got %v", err)
	}
}

func TestLoadUploadConfig_ReturnsErrorWhenConfigIsNil(t *testing.T) {
	err := LoadUploadConfig(nil, "", "LOKEX")
	if err == nil {
		t.Fatal("expected error for nil upload config")
	}

	if !strings.Contains(err.Error(), "upload config is nil") {
		t.Fatalf("expected nil config error, got %v", err)
	}
}

func TestLoadUploadConfig_ReturnsErrorForMissingExplicitConfig(t *testing.T) {
	cfg := &UploadConfig{}

	err := LoadUploadConfig(cfg, filepath.Join(t.TempDir(), "missing.yaml"), "LOKEX")
	if err == nil {
		t.Fatal("expected error for missing explicit config")
	}
}

func TestLoadUploadConfig_ReturnsErrorForInvalidExplicitConfig(t *testing.T) {
	configFile := filepath.Join(t.TempDir(), "broken.yaml")
	if err := os.WriteFile(configFile, []byte(":\nbad yaml"), 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	cfg := &UploadConfig{}

	err := LoadUploadConfig(cfg, configFile, "LOKEX")
	if err == nil {
		t.Fatal("expected error for invalid explicit config")
	}
}

func TestLoadUploadConfig_LoadsValuesFromConfigFile(t *testing.T) {
	configFile := filepath.Join(t.TempDir(), "lokex.yaml")
	content := `
upload:
  filename: messages.json
  src-path: ./locales/messages.json
  lang-iso: en
  poll: true
  format: json
  tags:
    - mobile
    - backend
  convert-placeholders: true
  filter-task-id: 123
`
	if err := os.WriteFile(configFile, []byte(strings.TrimSpace(content)), 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	cfg := &UploadConfig{}

	err := LoadUploadConfig(cfg, configFile, "LOKEX")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Filename == nil || *cfg.Filename != "messages.json" {
		t.Fatalf("expected Filename to be loaded from config, got %#v", cfg.Filename)
	}

	if cfg.SrcPath == nil || *cfg.SrcPath != "./locales/messages.json" {
		t.Fatalf("expected SrcPath to be loaded from config, got %#v", cfg.SrcPath)
	}

	if cfg.LangISO == nil || *cfg.LangISO != "en" {
		t.Fatalf("expected LangISO to be loaded from config, got %#v", cfg.LangISO)
	}

	if cfg.Poll == nil || *cfg.Poll != true {
		t.Fatalf("expected Poll to be true, got %#v", cfg.Poll)
	}

	if cfg.Format == nil || *cfg.Format != "json" {
		t.Fatalf("expected Format to be loaded from config, got %#v", cfg.Format)
	}

	wantTags := []string{"mobile", "backend"}
	if cfg.Tags == nil || !reflect.DeepEqual(*cfg.Tags, wantTags) {
		t.Fatalf("expected Tags %v, got %#v", wantTags, cfg.Tags)
	}

	if cfg.ConvertPlaceholders == nil || *cfg.ConvertPlaceholders != true {
		t.Fatalf("expected ConvertPlaceholders to be true, got %#v", cfg.ConvertPlaceholders)
	}

	if cfg.FilterTaskID == nil || *cfg.FilterTaskID != 123 {
		t.Fatalf("expected FilterTaskID to be 123, got %#v", cfg.FilterTaskID)
	}
}

func TestLoadUploadConfig_LoadsValuesFromEnv(t *testing.T) {
	t.Setenv("LOKEX_UPLOAD_FILENAME", "env-messages.json")
	t.Setenv("LOKEX_UPLOAD_SRC_PATH", "./env/messages.json")
	t.Setenv("LOKEX_UPLOAD_LANG_ISO", "fr")
	t.Setenv("LOKEX_UPLOAD_POLL", "true")
	t.Setenv("LOKEX_UPLOAD_CONTEXT_TIMEOUT", "2m")
	t.Setenv("LOKEX_UPLOAD_FORMAT", "xml")
	t.Setenv("LOKEX_UPLOAD_TAGS", "ios, backend")
	t.Setenv("LOKEX_UPLOAD_CONVERT_PLACEHOLDERS", "true")
	t.Setenv("LOKEX_UPLOAD_FILTER_TASK_ID", "987")

	cfg := &UploadConfig{}

	err := LoadUploadConfig(cfg, "", "LOKEX")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Filename == nil || *cfg.Filename != "env-messages.json" {
		t.Fatalf("expected Filename to be loaded from env, got %#v", cfg.Filename)
	}

	if cfg.SrcPath == nil || *cfg.SrcPath != "./env/messages.json" {
		t.Fatalf("expected SrcPath to be loaded from env, got %#v", cfg.SrcPath)
	}

	if cfg.LangISO == nil || *cfg.LangISO != "fr" {
		t.Fatalf("expected LangISO to be loaded from env, got %#v", cfg.LangISO)
	}

	if cfg.Poll == nil || *cfg.Poll != true {
		t.Fatalf("expected Poll to be true, got %#v", cfg.Poll)
	}

	if cfg.Format == nil || *cfg.Format != "xml" {
		t.Fatalf("expected Format to be loaded from env, got %#v", cfg.Format)
	}

	wantTags := []string{"ios", "backend"}
	if cfg.Tags == nil || !reflect.DeepEqual(*cfg.Tags, wantTags) {
		t.Fatalf("expected Tags %v, got %#v", wantTags, cfg.Tags)
	}

	if cfg.ConvertPlaceholders == nil || *cfg.ConvertPlaceholders != true {
		t.Fatalf("expected ConvertPlaceholders to be true, got %#v", cfg.ConvertPlaceholders)
	}

	if cfg.FilterTaskID == nil || *cfg.FilterTaskID != 987 {
		t.Fatalf("expected FilterTaskID to be 987, got %#v", cfg.FilterTaskID)
	}
}

func TestLoadUploadConfig_EnvOverridesConfigFile(t *testing.T) {
	configFile := filepath.Join(t.TempDir(), "lokex.yaml")
	content := `
upload:
  filename: config-messages.json
  lang-iso: en
  poll: false
  tags:
    - web
    - api
  convert-placeholders: false
  filter-task-id: 111
`
	if err := os.WriteFile(configFile, []byte(strings.TrimSpace(content)), 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	t.Setenv("LOKEX_UPLOAD_FILENAME", "env-messages.json")
	t.Setenv("LOKEX_UPLOAD_LANG_ISO", "de")
	t.Setenv("LOKEX_UPLOAD_POLL", "true")
	t.Setenv("LOKEX_UPLOAD_TAGS", "android,backend")
	t.Setenv("LOKEX_UPLOAD_CONVERT_PLACEHOLDERS", "true")
	t.Setenv("LOKEX_UPLOAD_FILTER_TASK_ID", "222")

	cfg := &UploadConfig{}

	err := LoadUploadConfig(cfg, configFile, "LOKEX")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Filename == nil || *cfg.Filename != "env-messages.json" {
		t.Fatalf("expected env to override config for Filename, got %#v", cfg.Filename)
	}

	if cfg.LangISO == nil || *cfg.LangISO != "de" {
		t.Fatalf("expected env to override config for LangISO, got %#v", cfg.LangISO)
	}

	if cfg.Poll == nil || *cfg.Poll != true {
		t.Fatalf("expected env to override config for Poll, got %#v", cfg.Poll)
	}

	wantTags := []string{"android", "backend"}
	if cfg.Tags == nil || !reflect.DeepEqual(*cfg.Tags, wantTags) {
		t.Fatalf("expected env to override config for Tags, got %#v", cfg.Tags)
	}

	if cfg.ConvertPlaceholders == nil || *cfg.ConvertPlaceholders != true {
		t.Fatalf("expected env to override config for ConvertPlaceholders, got %#v", cfg.ConvertPlaceholders)
	}

	if cfg.FilterTaskID == nil || *cfg.FilterTaskID != 222 {
		t.Fatalf("expected env to override config for FilterTaskID, got %#v", cfg.FilterTaskID)
	}
}

func TestLoadUploadConfig_EmptyStringSliceEnvDoesNotSetField(t *testing.T) {
	t.Setenv("LOKEX_UPLOAD_TAGS", " , ,  , ")

	cfg := &UploadConfig{}

	err := LoadUploadConfig(cfg, "", "LOKEX")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Tags != nil {
		t.Fatalf("expected Tags to stay nil, got %#v", cfg.Tags)
	}
}
