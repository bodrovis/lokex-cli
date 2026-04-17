package viper_helpers

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

func ApplyConfigValue[T any](v *viper.Viper, key string, set func(T)) {
	if !v.IsSet(key) {
		return
	}

	val, ok := v.Get(key).(T)
	if ok {
		set(val)
		return
	}

	var zero T
	switch any(zero).(type) {
	case string:
		set(any(v.GetString(key)).(T))
	case bool:
		set(any(v.GetBool(key)).(T))
	case int64:
		set(any(v.GetInt64(key)).(T))
	case time.Duration:
		set(any(v.GetDuration(key)).(T))
	}
}

func ApplyConfigStringSlice(v *viper.Viper, key string, set func([]string)) {
	if !v.IsSet(key) {
		return
	}

	raw := v.Get(key)

	switch val := raw.(type) {
	case []string:
		clean := make([]string, 0, len(val))
		for _, item := range val {
			item = strings.TrimSpace(item)
			if item != "" {
				clean = append(clean, item)
			}
		}
		if len(clean) > 0 {
			set(clean)
		}
	case []any:
		clean := make([]string, 0, len(val))
		for _, item := range val {
			s := strings.TrimSpace(fmt.Sprint(item))
			if s != "" {
				clean = append(clean, s)
			}
		}
		if len(clean) > 0 {
			set(clean)
		}
	case string:
		parts := strings.Split(val, ",")
		clean := make([]string, 0, len(parts))
		for _, item := range parts {
			item = strings.TrimSpace(item)
			if item != "" {
				clean = append(clean, item)
			}
		}
		if len(clean) > 0 {
			set(clean)
		}
	}
}
