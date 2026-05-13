//go:build windows

package config

import "os"

// IsSecure returns true if path exists.
// Full NTFS DACL validation is a Phase 2 enhancement.
func IsSecure(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// EnforcePermissions is a no-op on Windows;
// permissions are controlled by the NTFS ACL set at directory creation.
func EnforcePermissions(path string) error {
	return nil
}
