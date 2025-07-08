package configx

import (
	"github.com/knadh/koanf/v2"
	"github.com/spf13/pflag"
)

type Provider struct {
	flags             *pflag.FlagSet
	onValidationError func(k *koanf.Koanf, err error)
}
