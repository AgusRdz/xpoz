// Package store handles SQLite persistence for tunnel records.
package store

import (
	"context"
	"fmt"
	"time"
)

// Tunnel is a saved tunnel record.
type Tunnel struct {
	Name      string
	Subdomain string
	FullURL   string
	LocalPort int
	Active    bool
	CreatedAt time.Time
	LastUsed  time.Time
}

// Store is the interface for tunnel persistence.
// It is injected as an interface to allow testing without a real database.
type Store interface {
	// Get returns the tunnel with the given name, or ErrNotFound.
	Get(ctx context.Context, name string) (*Tunnel, error)
	// GetBySubdomain returns the tunnel with the given subdomain, or ErrNotFound.
	GetBySubdomain(ctx context.Context, subdomain string) (*Tunnel, error)
	// List returns all saved tunnels ordered by created_at descending.
	List(ctx context.Context) ([]*Tunnel, error)
	// Upsert creates or updates a tunnel record.
	// On conflict, created_at is preserved and all other fields are updated.
	Upsert(ctx context.Context, t *Tunnel) error
	// Delete removes the tunnel with the given name, or ErrNotFound.
	Delete(ctx context.Context, name string) error
	// SetActive updates the active flag of the named tunnel, or ErrNotFound.
	SetActive(ctx context.Context, name string, active bool) error
	// Close releases database resources.
	Close() error
}

// ErrNotFound is returned when a tunnel does not exist in the store.
var ErrNotFound = fmt.Errorf("tunnel not found")
