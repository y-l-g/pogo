//go:build !linux

package pogo

import (
	"os/exec"
)

func configureCmd(cmd *exec.Cmd) {
	// No-op for non-Linux systems.
	// On these systems, we rely on the Pipe closure detection in the worker
	// loop to trigger exit.
}
