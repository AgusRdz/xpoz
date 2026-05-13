// Package updater handles self-update with Ed25519 signature verification.
package updater

import (
	"fmt"
)

// publicKey is the Ed25519 public key used to verify release checksums.
// Generate the key pair with: make keygen
// The private key is stored as the SIGNING_KEY GitHub Actions secret.
const publicKey = "" // TODO: populate after running make keygen

// CheckLatest returns the latest available version tag from GitHub releases.
func CheckLatest() (string, error) {
	// TODO: implement in Phase 2
	return "", fmt.Errorf("updater: not implemented")
}

// Update downloads, verifies (SHA256 + Ed25519), and atomically replaces the binary.
func Update(current string) error {
	// TODO: implement in Phase 2
	return fmt.Errorf("updater: not implemented")
}

// ApplyPending applies a staged update binary if one is waiting from a previous run.
// Must be called at startup before any other logic.
func ApplyPending(current string) {
	// TODO: implement in Phase 2
}

// NotifyIfAvailable writes a one-line update hint to stderr if a newer version exists.
// It reads a cached check result to avoid a network call on every run.
func NotifyIfAvailable(current string) {
	// TODO: implement in Phase 2
}

// BackgroundCheck spawns a goroutine that checks for updates and caches the result.
// Call at the end of main, after the main command completes.
func BackgroundCheck(current string) {
	// TODO: implement in Phase 2
}
