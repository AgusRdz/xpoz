// Package updater handles self-update with Ed25519 signature verification.
package updater

import (
	"context"
	"fmt"
	"os"
	"time"
)

const repo = "AgusRdz/xpoz"

// CheckLatest returns the latest available version tag from GitHub releases.
func CheckLatest(ctx context.Context) (string, error) {
	return fetchLatestTag(ctx)
}

// Update downloads, verifies (SHA256 + Ed25519), and atomically replaces the binary.
func Update(ctx context.Context, current string) error {
	tag, err := CheckLatest(ctx)
	if err != nil {
		return err
	}
	if tag == current {
		return fmt.Errorf("already on the latest version (%s)", current)
	}
	return installRelease(ctx, tag)
}

// NotifyIfAvailable writes a one-line update hint to stderr if a newer version exists.
// It reads a cached check result to avoid a network call on every run.
func NotifyIfAvailable(current string) {
	tag, at, err := readCache()
	if err != nil || time.Since(at) > 24*time.Hour {
		return
	}
	if tag != "" && tag != current {
		fmt.Fprintf(os.Stderr, "\n  Update available: %s → %s\n  Run 'xpoz update' to upgrade.\n\n", current, tag)
	}
}

// BackgroundCheck spawns a goroutine that checks for updates and caches the result.
// Call at the end of main, after the main command completes.
func BackgroundCheck(current string) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		tag, err := CheckLatest(ctx)
		if err != nil {
			return
		}
		_ = writeCache(tag)
	}()
}

// ApplyPending applies a staged update binary (bin+".new") if one is waiting.
// Must be called at startup before any other logic.
func ApplyPending(current string) {
	bin, err := os.Executable()
	if err != nil {
		return
	}
	pending := bin + ".new"
	if _, err := os.Stat(pending); err != nil {
		return
	}
	if err := applyPendingUpdate(bin, pending); err != nil {
		fmt.Fprintf(os.Stderr, "⚠ failed to apply pending update: %v\n", err)
	}
}

// applyPendingUpdate atomically replaces bin with pending, rolling back on failure.
func applyPendingUpdate(bin, pending string) error {
	old := bin + ".old"
	_ = os.Rename(bin, old)
	if err := os.Rename(pending, bin); err != nil {
		// Try to restore the old binary.
		_ = os.Rename(old, bin)
		return fmt.Errorf("replacing binary: %w", err)
	}
	_ = os.Remove(old)
	return nil
}
