package contextx

import (
	"context"

	"github.com/huanggze/x/configx"
)

type Default struct{}

var _ Contextualizer = (*Default)(nil)

func (d *Default) Config(ctx context.Context, config *configx.Provider) *configx.Provider {
	return config
}
