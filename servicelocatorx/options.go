package servicelocatorx

import (
	"github.com/huanggze/x/contextx"
	"github.com/huanggze/x/logrusx"
)

type (
	Options struct {
		logger         *logrusx.Logger
		contextualizer contextx.Contextualizer
	}
	Option func(o *Options)
)

func (o *Options) Logger() *logrusx.Logger {
	return o.logger
}

func (o *Options) Contextualizer() contextx.Contextualizer {
	return o.contextualizer
}

func NewOptions(options ...Option) *Options {
	o := &Options{
		contextualizer: &contextx.Default{},
	}
	for _, opt := range options {
		opt(o)
	}
	return o
}
