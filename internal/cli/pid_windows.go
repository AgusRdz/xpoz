//go:build windows

package cli

import "os"

// isRunning opens the process handle; an error means the process is gone.
// Note: PID reuse is a known edge case on Windows.
func isRunning(pid int) bool {
	_, err := os.FindProcess(pid)
	return err == nil
}
