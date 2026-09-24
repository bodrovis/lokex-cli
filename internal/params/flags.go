package params

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func BindFlags(
	cmd *cobra.Command,
	specs []Spec,
) {
	BindFlagSet(cmd.Flags(), specs)
}

func BindFlagSet(
	fs *pflag.FlagSet,
	specs []Spec,
) {
	fs.SortFlags = false

	for _, spec := range specs {
		switch spec.Kind {
		case KindString:
			value, _ := spec.Default.(string)
			fs.String(spec.Name, value, spec.Usage)

		case KindBool:
			value, _ := spec.Default.(bool)
			fs.Bool(spec.Name, value, spec.Usage)

		case KindStringSlice:
			value, _ := spec.Default.([]string)
			fs.StringSlice(spec.Name, value, spec.Usage)

		case KindInt:
			value, _ := spec.Default.(int)
			fs.Int(spec.Name, value, spec.Usage)

		case KindInt64:
			value, _ := spec.Default.(int64)
			fs.Int64(spec.Name, value, spec.Usage)

		case KindDuration:
			value, _ := spec.Default.(time.Duration)
			fs.Duration(spec.Name, value, spec.Usage)

		default:
			panic(fmt.Sprintf(
				"unsupported parameter kind %d for %q",
				spec.Kind,
				spec.Name,
			))
		}
	}
}

func flagValue(
	cmd *cobra.Command,
	spec Spec,
) (any, error) {
	switch spec.Kind {
	case KindString:
		return cmd.Flags().GetString(spec.Name)

	case KindBool:
		return cmd.Flags().GetBool(spec.Name)

	case KindStringSlice:
		return cmd.Flags().GetStringSlice(spec.Name)

	case KindInt:
		return cmd.Flags().GetInt(spec.Name)

	case KindInt64:
		return cmd.Flags().GetInt64(spec.Name)

	case KindDuration:
		return cmd.Flags().GetDuration(spec.Name)

	default:
		return nil, fmt.Errorf(
			"unsupported parameter kind %d for %q",
			spec.Kind,
			spec.Name,
		)
	}
}
