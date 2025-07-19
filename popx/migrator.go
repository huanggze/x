package popx

import "github.com/huanggze/x/logrusx"

// Migrator forms the basis of all migrations systems.
// It does the actual heavy lifting of running migrations.
// When building a new migration system, you should embed this
// type into your migrator.
type Migrator struct {
	Migrations map[string]Migrations
	l          *logrusx.Logger
}
