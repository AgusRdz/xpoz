// Package cleanup handles xpoz uninstall and data removal.
package cleanup

import (
	"fmt"
	"os"

	"github.com/AgusRdz/xpoz/internal/config"
)

// Options controls what gets removed during uninstall.
type Options struct {
	KeepData bool // preserve tunnels.db when true
}

// Run performs the full uninstall sequence:
// stops running tunnels, removes config, removes database (unless KeepData),
// and removes the binary.
func Run(opts Options) error {
	if !opts.KeepData {
		if err := removeIgnoreNotExist(config.DBPath()); err != nil {
			return fmt.Errorf("cleanup: removing database: %w", err)
		}
		if err := removeIgnoreNotExist(config.Path()); err != nil {
			return fmt.Errorf("cleanup: removing config: %w", err)
		}
		if err := os.RemoveAll(config.Dir()); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("cleanup: removing data directory: %w", err)
		}
		return nil
	}
	if err := removeIgnoreNotExist(config.Path()); err != nil {
		return fmt.Errorf("cleanup: removing config: %w", err)
	}
	return nil
}

// Reset clears the tunnel database without removing
// the binary, config, or service registration.
func Reset() error {
	if err := removeIgnoreNotExist(config.DBPath()); err != nil {
		return fmt.Errorf("cleanup: resetting database: %w", err)
	}
	return nil
}

// removeIgnoreNotExist removes path and returns nil if the file does not exist.
func removeIgnoreNotExist(path string) error {
	err := os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
