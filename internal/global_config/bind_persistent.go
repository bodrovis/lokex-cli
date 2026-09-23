package global_config

import (
	"github.com/bodrovis/lokex-cli/internal/params"
	"github.com/spf13/pflag"
)

func BindPersistentFlags(
	fs *pflag.FlagSet,
	userAgent string,
) {
	params.BindFlagSet(
		fs,
		paramSpecs(userAgent),
	)
}
