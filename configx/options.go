package configx

import (
	"fmt"
	"io"
	"os"

	"github.com/huanggze/x/jsonschemax"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/v2"
)

type (
	OptionModifier func(p *Provider)
)

//func WithImmutables(immutables ...string) OptionModifier {
//	return func(p *Provider) {
//		p.immutables = append(p.immutables, immutables...)
//	}
//}
//
//func WithExceptImmutables(exceptImmutables ...string) OptionModifier {
//	return func(p *Provider) {
//		p.exceptImmutables = append(p.exceptImmutables, exceptImmutables...)
//	}
//}
//
//func WithFlags(flags *pflag.FlagSet) OptionModifier {
//	return func(p *Provider) {
//		p.flags = flags
//	}
//}
//
//func WithLogger(l *logrusx.Logger) OptionModifier {
//	return func(p *Provider) {
//		p.logger = l
//	}
//}
//
//// DEPRECATED without replacement. This option is a no-op.
//func OmitKeysFromTracing(keys ...string) OptionModifier {
//	return func(*Provider) {}
//}
//
//func AttachWatcher(watcher func(event watcherx.Event, err error)) OptionModifier {
//	return func(p *Provider) {
//		p.onChanges = append(p.onChanges, watcher)
//	}
//}
//
//func WithLogrusWatcher(l *logrusx.Logger) OptionModifier {
//	return AttachWatcher(LogrusWatcher(l))
//}

func WithStderrValidationReporter() OptionModifier {
	return func(p *Provider) {
		p.onValidationError = func(k *koanf.Koanf, err error) {
			p.printHumanReadableValidationErrors(k, os.Stderr, err)
		}
	}
}

func (p *Provider) printHumanReadableValidationErrors(k *koanf.Koanf, w io.Writer, err error) {
	if err == nil {
		return
	}

	_, _ = fmt.Fprintln(os.Stderr, "")
	conf, innerErr := k.Marshal(json.Parser())
	if innerErr != nil {
		_, _ = fmt.Fprintf(w, "Unable to unmarshal configuration: %+v", innerErr)
	}

	jsonschemax.FormatValidationErrorForCLI(w, conf, err)
}
