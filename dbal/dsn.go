package dbal

import "regexp"

var dsnRegex = regexp.MustCompile(`^(sqlite://file:(?:.+)\?((\w+=\w+)(&\w+=\w+)*)?(&?mode=memory)(&\w+=\w+)*)$|(?:sqlite://(file:)?:memory:(?:\?\w+=\w+)?(?:&\w+=\w+)*)|^(?:(?::memory:)|(?:memory))$`)

// IsMemorySQLite returns true if a given DSN string is pointing to a SQLite database.
//
// SQLite can be written in different styles depending on the use case
// - just in memory
// - shared connection
// - shared but unique in the same process
// see: https://sqlite.org/inmemorydb.html
func IsMemorySQLite(dsn string) bool {
	return dsnRegex.MatchString(dsn)
}
