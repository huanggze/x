package contextx

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/huanggze/x/configx"
)

type (
	Contextualizer interface {
		// Network returns the network id for the given context.
		Network(ctx context.Context, network uuid.UUID) uuid.UUID
		// Config returns the config for the given context.
		Config(ctx context.Context, config *configx.Provider) *configx.Provider
	}
	Provider interface {
		Contextualizer() Contextualizer
	}
)
