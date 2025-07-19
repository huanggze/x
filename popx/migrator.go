package popx

import (
	"time"

	"github.com/huanggze/x/logrusx"
	"github.com/huanggze/x/otelx"
	"github.com/ory/pop/v6"
)

// NewMigrator returns a new "blank" migrator. It is recommended
// to use something like MigrationBox or FileMigrator. A "blank"
// Migrator should only be used as the basis for a new type of
// migration system.
func NewMigrator(c *pop.Connection, l *logrusx.Logger, tracer *otelx.Tracer, perMigrationTimeout time.Duration) *Migrator {
	return &Migrator{
		Connection: c,
		l:          l,
		Migrations: map[string]Migrations{
			"up":   {},
			"down": {},
		},
		tracer:              tracer,
		PerMigrationTimeout: perMigrationTimeout,
	}
}

// Migrator forms the basis of all migrations systems.
// It does the actual heavy lifting of running migrations.
// When building a new migration system, you should embed this
// type into your migrator.
type Migrator struct {
	Connection          *pop.Connection
	Migrations          map[string]Migrations
	l                   *logrusx.Logger
	PerMigrationTimeout time.Duration
	tracer              *otelx.Tracer

	// DumpMigrations if true will dump the migrations to a file called schema.sql
	DumpMigrations bool
}
