package logrusx

import (
	"cmp"
	"os"
	"time"

	gelf "github.com/seatgeek/logrus-gelf-formatter"
	"github.com/sirupsen/logrus"

	"github.com/huanggze/x/stringsx"
)

type (
	options struct {
		level     *logrus.Level
		formatter logrus.Formatter
		format    string
		c         configurator
	}
	Option           func(*options)
	nullConfigurator struct{}
	configurator     interface {
		Bool(key string) bool
		String(key string) string
	}
)

func setLevel(l *logrus.Logger, o *options) {
	if o.level != nil {
		l.Level = *o.level
	} else {
		var err error
		l.Level, err = logrus.ParseLevel(cmp.Or(
			o.c.String("log.level"),
			os.Getenv("LOG_LEVEL")))
		if err != nil {
			l.Level = logrus.InfoLevel
		}
	}
}

func setFormatter(l *logrus.Logger, o *options) {
	if o.formatter != nil {
		l.Formatter = o.formatter
	} else {
		var unknownFormat bool // we first have to set the formatter before we can complain about the unknown format

		format := stringsx.SwitchExact(cmp.Or(o.format, o.c.String("log.format"), os.Getenv("LOG_FORMAT")))
		switch {
		case format.AddCase("json"):
			l.Formatter = &logrus.JSONFormatter{PrettyPrint: false, TimestampFormat: time.RFC3339Nano, DisableHTMLEscape: true}
		case format.AddCase("json_pretty"):
			l.Formatter = &logrus.JSONFormatter{PrettyPrint: true, TimestampFormat: time.RFC3339Nano, DisableHTMLEscape: true}
		case format.AddCase("gelf"):
			l.Formatter = new(gelf.GelfFormatter)
		default:
			unknownFormat = true
			fallthrough
		case format.AddCase("text", ""):
			l.Formatter = &logrus.TextFormatter{
				DisableQuote:     true,
				DisableTimestamp: false,
				FullTimestamp:    true,
			}
		}

		if unknownFormat {
			l.WithError(format.ToUnknownCaseErr()).Warn("got unknown \"log.format\", falling back to \"text\"")
		}
	}
}

func WithConfigurator(c configurator) Option {
	return func(o *options) {
		o.c = c
	}
}

func (c *nullConfigurator) Bool(_ string) bool {
	return false
}

func (c *nullConfigurator) String(_ string) string {
	return ""
}

func newOptions(opts []Option) *options {
	o := new(options)
	o.c = new(nullConfigurator)
	for _, f := range opts {
		f(o)
	}
	return o
}

func (l *Logger) UseConfig(c configurator) {
	l.leakSensitive = l.leakSensitive || c.Bool("log.leak_sensitive_values")
	l.redactionText = cmp.Or(c.String("log.redaction_text"), l.redactionText)
	o := newOptions(append(l.opts, WithConfigurator(c)))
	setLevel(l.Entry.Logger, o)
	setFormatter(l.Entry.Logger, o)
}
