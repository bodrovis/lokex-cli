package viper_helpers

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestNewConfigViper_UsesExplicitConfigFile(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "custom.yaml")

	if err := os.WriteFile(configFile, []byte("token: explicit-token\n"), 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	v := NewConfigViper(configFile, "LOKEX")

	if err := v.ReadInConfig(); err != nil {
		t.Fatalf("expected explicit config file to be readable, got error: %v", err)
	}

	if got := v.GetString("token"); got != "explicit-token" {
		t.Fatalf("expected token from explicit config file, got %q", got)
	}
}

func TestNewConfigViper_ReadsDefaultConfigFromCurrentDirectory(t *testing.T) {
	origWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer func() {
		_ = os.Chdir(origWD)
	}()

	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change working directory: %v", err)
	}

	configFile := filepath.Join(tmpDir, "lokex.yaml")
	if err := os.WriteFile(configFile, []byte("token: from-current-dir\n"), 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	v := NewConfigViper("", "LOKEX")

	if err := v.ReadInConfig(); err != nil {
		t.Fatalf("expected config from current directory to be readable, got error: %v", err)
	}

	if got := v.GetString("token"); got != "from-current-dir" {
		t.Fatalf("expected token from current directory config, got %q", got)
	}
}

func TestNewConfigViper_ReadsDefaultConfigFromHomeConfigDir(t *testing.T) {
	origWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer func() {
		_ = os.Chdir(origWD)
	}()

	emptyWD := t.TempDir()
	if err := os.Chdir(emptyWD); err != nil {
		t.Fatalf("failed to change working directory: %v", err)
	}

	homeDir := t.TempDir()

	origUserHomeDir := userHomeDir
	userHomeDir = func() (string, error) {
		return homeDir, nil
	}
	defer func() {
		userHomeDir = origUserHomeDir
	}()

	configDir := filepath.Join(homeDir, ".config", "lokex-cli")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("failed to create config directory: %v", err)
	}

	configFile := filepath.Join(configDir, "lokex.yaml")
	if err := os.WriteFile(configFile, []byte("token: from-home-config\n"), 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	v := NewConfigViper("", "LOKEX")

	if err := v.ReadInConfig(); err != nil {
		t.Fatalf("expected config from home config dir to be readable, got error: %v", err)
	}

	if got := v.GetString("token"); got != "from-home-config" {
		t.Fatalf("expected token from home config dir, got %q", got)
	}
}

func TestNewConfigViper_CurrentDirectoryTakesPrecedenceOverHomeConfigDir(t *testing.T) {
	origWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer func() {
		_ = os.Chdir(origWD)
	}()

	wd := t.TempDir()
	if err := os.Chdir(wd); err != nil {
		t.Fatalf("failed to change working directory: %v", err)
	}

	homeDir := t.TempDir()

	origUserHomeDir := userHomeDir
	userHomeDir = func() (string, error) {
		return homeDir, nil
	}
	defer func() {
		userHomeDir = origUserHomeDir
	}()

	if err := os.WriteFile(filepath.Join(wd, "lokex.yaml"), []byte("token: from-current-dir\n"), 0o644); err != nil {
		t.Fatalf("failed to write current directory config: %v", err)
	}

	homeConfigDir := filepath.Join(homeDir, ".config", "lokex-cli")
	if err := os.MkdirAll(homeConfigDir, 0o755); err != nil {
		t.Fatalf("failed to create home config directory: %v", err)
	}

	if err := os.WriteFile(filepath.Join(homeConfigDir, "lokex.yaml"), []byte("token: from-home-config\n"), 0o644); err != nil {
		t.Fatalf("failed to write home config: %v", err)
	}

	v := NewConfigViper("", "LOKEX")

	if err := v.ReadInConfig(); err != nil {
		t.Fatalf("expected config to be readable, got error: %v", err)
	}

	if got := v.GetString("token"); got != "from-current-dir" {
		t.Fatalf("expected current directory config to take precedence, got %q", got)
	}
}

func TestNewConfigViper_ReadsPrefixedEnvVar(t *testing.T) {
	t.Setenv("LOKEX_TOKEN", "env-token")

	v := NewConfigViper("", "LOKEX")

	if got := v.GetString("token"); got != "env-token" {
		t.Fatalf("expected token from prefixed env var, got %q", got)
	}
}

func TestNewConfigViper_ReadsEnvVarWithoutPrefix(t *testing.T) {
	t.Setenv("TOKEN", "plain-env-token")

	v := NewConfigViper("", "")

	if got := v.GetString("token"); got != "plain-env-token" {
		t.Fatalf("expected token from unprefixed env var, got %q", got)
	}
}

func TestNewConfigViper_ReplacesDashAndDotInEnvKeys(t *testing.T) {
	t.Setenv("LOKEX_PROJECT_ID", "project-from-env")
	t.Setenv("LOKEX_PARENT_CHILD", "nested-from-env")

	v := NewConfigViper("", "LOKEX")

	if got := v.GetString("project-id"); got != "project-from-env" {
		t.Fatalf("expected project-id from env var, got %q", got)
	}

	if got := v.GetString("parent.child"); got != "nested-from-env" {
		t.Fatalf("expected parent.child from env var, got %q", got)
	}
}

func TestNewConfigViper_EnvOverridesConfigValues(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "lokex.yaml")

	if err := os.WriteFile(configFile, []byte("token: config-token\n"), 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	t.Setenv("LOKEX_TOKEN", "env-token")

	v := NewConfigViper(configFile, "LOKEX")

	if err := v.ReadInConfig(); err != nil {
		t.Fatalf("expected config file to be readable, got error: %v", err)
	}

	if got := v.GetString("token"); got != "env-token" {
		t.Fatalf("expected env to override config value, got %q", got)
	}
}

func TestNewConfigViper_IgnoresHomeDirWhenUserHomeDirFails(t *testing.T) {
	origWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer func() {
		_ = os.Chdir(origWD)
	}()

	emptyWD := t.TempDir()
	if err := os.Chdir(emptyWD); err != nil {
		t.Fatalf("failed to change working directory: %v", err)
	}

	origUserHomeDir := userHomeDir
	userHomeDir = func() (string, error) {
		return "", errors.New("boom")
	}
	defer func() {
		userHomeDir = origUserHomeDir
	}()

	v := NewConfigViper("", "LOKEX")

	err = v.ReadInConfig()
	if err == nil {
		t.Fatal("expected error because no config file should be found")
	}
}
