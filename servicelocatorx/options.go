package servicelocatorx

import "github.com/huanggze/x/contextx"

type (
	Options struct {
		contextualizer contextx.Contextualizer
	}
	Option func(o *Options)
)

func NewOptions(options ...Option) *Options {
	o := &Options{
		contextualizer: &contextx.Default{},
	}
	for _, opt := range options {
		opt(o)
	}
	return o
}
