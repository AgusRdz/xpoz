//go:build windows

package cli

import (
	"os/exec"
	"syscall"
)

// DETACHED_PROCESS creates the child without a console window so it runs
// independently of the parent's terminal session.
const detachedProcess = 0x00000008

func setSysProcAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: detachedProcess}
}
