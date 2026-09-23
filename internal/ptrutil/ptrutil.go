package ptrutil

import "strings"

func Value[T any](ptr *T) T {
	if ptr == nil {
		var zero T
		return zero
	}

	return *ptr
}

func TrimmedString(ptr *string) string {
	if ptr == nil {
		return ""
	}

	return strings.TrimSpace(*ptr)
}
