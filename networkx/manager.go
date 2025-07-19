package networkx

import "embed"

// Migrations of the network manager. Apply by merging with your local migrations using
// fsx.Merge() and then passing all to the migration box.
//
//go:embed migrations/sql/*.sql
var Migrations embed.FS
