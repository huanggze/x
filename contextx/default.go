package contextx

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/huanggze/x/configx"
)

type Default struct{}

var _ Contextualizer = (*Default)(nil)

func (d *Default) Network(ctx context.Context, network uuid.UUID) uuid.UUID {
	if network == uuid.Nil {
		panic("nid must be not nil")
	}
	return network
}

func (d *Default) Config(ctx context.Context, config *configx.Provider) *configx.Provider {
	return config
}
