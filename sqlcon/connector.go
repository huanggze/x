// Package sqlcon provides helpers for dealing with SQL connectivity.
package sqlcon

import (
	"runtime"
	"strings"
)

// GetDriverName returns the driver name of a given DSN.
func GetDriverName(dsn string) string {
	return strings.Split(dsn, "://")[0]
}
func maxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU
}
