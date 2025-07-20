package cmdx

import "github.com/spf13/pflag"

const (
	FlagQuiet = "quiet"
)

func RegisterNoiseFlags(flags *pflag.FlagSet) {
	flags.BoolP(FlagQuiet, FlagQuiet[:1], false, "Be quiet with output printing.")
}
