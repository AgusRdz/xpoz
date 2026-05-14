//go:build !windows

package cli

import (
	"os/exec"
	"syscall"
)

// setSysProcAttr detaches the child process from the current session so it
// survives terminal close. Setsid creates a new process group and session.
func setSysProcAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
}
