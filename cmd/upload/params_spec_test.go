package upload

import (
	"strings"
	"testing"
)

func TestUploadParamSpecs_ConfigKeysAreNotBlank(t *testing.T) {
	for _, spec := range uploadParamSpecs {
		if strings.TrimSpace(spec.ConfigKey) == "" {
			t.Fatalf("spec for flag %q has blank ConfigKey", spec.FlagName)
		}
	}
}

func TestUploadParamSpecs_ConfigKeysAreUnique(t *testing.T) {
	seen := map[string]bool{}

	for _, spec := range uploadParamSpecs {
		key := strings.TrimSpace(spec.ConfigKey)
		if key == "" {
			continue
		}
		if seen[key] {
			t.Fatalf("duplicate config key %q", key)
		}
		seen[key] = true
	}
}
