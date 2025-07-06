package dbal

import (
	"context"
	"sync"
)

var (
	drivers = make([]func() Driver, 0)
	dmtx    sync.Mutex
)

// Driver represents a driver
type Driver interface {
	// CanHandle returns true if the driver is capable of handling the given DSN or false otherwise.
	CanHandle(dsn string) bool

	// Ping returns nil if the driver has connectivity and is healthy or an error otherwise.
	Ping() error
	PingContext(context.Context) error
}

// RegisterDriver registers a driver
func RegisterDriver(d func() Driver) {
	dmtx.Lock()
	drivers = append(drivers, d)
	dmtx.Unlock()
}
