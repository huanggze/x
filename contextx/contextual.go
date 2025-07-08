package contextx

import (
	"context"

	"github.com/huanggze/x/configx"
)

type (
	Contextualizer interface {
		// Config returns the config for the given context.
		Config(ctx context.Context, config *configx.Provider) *configx.Provider
	}
)
