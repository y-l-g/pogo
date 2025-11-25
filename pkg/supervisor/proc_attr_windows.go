//go:build windows

package supervisor

import (
	"os/exec"
)

func configureCmd(cmd *exec.Cmd) {
	// Windows Process Management relies on the Worker detecting Pipe Closure.
	// We ensure the Protocol exits on EOF to prevent zombies.
}
