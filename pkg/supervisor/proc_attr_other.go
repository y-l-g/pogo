//go:build !linux && !windows

package supervisor

import (
	"os/exec"
)

func configureCmd(cmd *exec.Cmd) {
	// No-op for systems without Pdeathsig or Job Objects (e.g. BSD, Darwin).
	// On these systems, we rely on the Pipe closure detection in the worker
	// loop to trigger exit.
}
