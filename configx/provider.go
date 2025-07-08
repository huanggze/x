package configx

import (
	"github.com/knadh/koanf/v2"
	"github.com/spf13/pflag"

	"github.com/huanggze/x/logrusx"
	"github.com/huanggze/x/watcherx"
)

type Provider struct {
	*koanf.Koanf
	immutables, exceptImmutables []string

	flags             *pflag.FlagSet
	onChanges         []func(watcherx.Event, error)
	onValidationError func(k *koanf.Koanf, err error)

	logger *logrusx.Logger

	userProviders []koanf.Provider
}
