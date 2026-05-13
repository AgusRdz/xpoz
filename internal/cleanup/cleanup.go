// Package cleanup handles xpoz uninstall and data removal.
package cleanup

import "fmt"

// Options controls what gets removed during uninstall.
type Options struct {
	KeepData bool // preserve tunnels.db when true
}

// Run performs the full uninstall sequence:
// stops running tunnels, removes config, removes database (unless KeepData),
// and removes the binary.
func Run(opts Options) error {
	// TODO: implement in Phase 2
	return fmt.Errorf("cleanup: not implemented")
}

// Reset clears the tunnel database and log files without removing
// the binary, config, or service registration.
func Reset() error {
	// TODO: implement in Phase 2
	return fmt.Errorf("cleanup: not implemented")
}
