package networkx

import (
	"context"
	"embed"
	"errors"

	"github.com/huanggze/x/logrusx"
	"github.com/huanggze/x/otelx"
	"github.com/huanggze/x/sqlcon"
	"github.com/ory/pop/v6"
)

// Migrations of the network manager. Apply by merging with your local migrations using
// fsx.Merge() and then passing all to the migration box.
//
//go:embed migrations/sql/*.sql
var Migrations embed.FS

type Manager struct {
	c *pop.Connection
	l *logrusx.Logger
	t *otelx.Tracer
}

func NewManager(
	c *pop.Connection,
	l *logrusx.Logger,
	t *otelx.Tracer,
) *Manager {
	return &Manager{
		c: c,
		l: l,
		t: t,
	}
}

func (m *Manager) Determine(ctx context.Context) (*Network, error) {
	var p Network
	c := m.c.WithContext(ctx)
	if err := sqlcon.HandleError(c.Q().Order("created_at ASC").First(&p)); err != nil {
		if errors.Is(err, sqlcon.ErrNoRows) {
			np := NewNetwork()
			if err := c.Create(np); err != nil {
				return nil, err
			}
			return np, nil
		}
		return nil, err
	}
	return &p, nil
}
