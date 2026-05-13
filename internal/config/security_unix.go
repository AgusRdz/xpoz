//go:build !windows

package config

import (
	"fmt"
	"os"
)

// IsSecure returns true if path is readable only by its owner (mode & 0077 == 0).
func IsSecure(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode().Perm()&0o077 == 0
}

// EnforcePermissions sets path to 0600 permissions.
func EnforcePermissions(path string) error {
	if err := os.Chmod(path, 0600); err != nil {
		return fmt.Errorf("config: setting permissions on %s: %w", path, err)
	}
	return nil
}
