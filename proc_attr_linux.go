//go:build linux

package pogo

import (
	"os/exec"
	"syscall"
)

func configureCmd(cmd *exec.Cmd) {
	// Pdeathsig: SIGTERM (Graceful) instead of SIGKILL.
	// If the parent process dies, the kernel sends SIGTERM to the child,
	// allowing it to run shutdown handlers/destructors.
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Pdeathsig: syscall.SIGTERM,
	}
}
