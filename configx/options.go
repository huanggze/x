package configx

import (
	"github.com/spf13/pflag"
)

type (
	OptionModifier func(p *Provider)
)

func WithFlags(flags *pflag.FlagSet) OptionModifier {
	return func(p *Provider) {
		p.flags = flags
	}
}
