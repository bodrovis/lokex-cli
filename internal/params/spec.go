package params

type Kind uint8

const (
	KindString Kind = iota
	KindBool
	KindStringSlice
	KindInt
	KindInt64
	KindDuration
)

type Spec struct {
	Name    string
	Kind    Kind
	Usage   string
	Default any
}

func configKey(prefix, name string) string {
	if prefix == "" {
		return name
	}

	return prefix + "." + name
}
