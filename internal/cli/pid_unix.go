//go:build !windows

package cli

import (
	"os"
	"syscall"
)

// isRunning sends signal 0 to check process existence without affecting it.
func isRunning(pid int) bool {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	return proc.Signal(syscall.Signal(0)) == nil
}
