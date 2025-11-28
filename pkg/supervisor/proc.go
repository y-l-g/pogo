package supervisor

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// Process wraps os/exec.Cmd to provide a clean interface for
// spawning and managing PHP worker processes.
type Process struct {
	cmd *exec.Cmd

	// Pipes exposed for Transport
	ParentRead  *os.File
	ParentWrite *os.File

	// Child ends to close after start
	childRead  *os.File
	childWrite *os.File
}

// NewProcess prepares a PHP worker process command.
func NewProcess(ctx context.Context, entrypoint string, env map[string]string, extraFiles []*os.File) (*Process, error) {
	var bin string
	var args []string

	if testBin := os.Getenv("POGO_TEST_PHP_BINARY"); testBin != "" {
		bin = testBin
		args = []string{entrypoint}
	} else {
		ex, err := os.Executable()
		if err != nil {
			bin = "php"
		} else {
			bin = ex
		}
		args = []string{"php-cli", entrypoint}
	}

	cmd := exec.CommandContext(ctx, bin, args...)

	// Setup Pipes
	// FD 3: Input (Parent Write -> Child Read)
	// FD 4: Output (Child Write -> Parent Read)

	pRead, cWrite, err := os.Pipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create read pipe: %w", err)
	}

	cRead, pWrite, err := os.Pipe()
	if err != nil {
		_ = pRead.Close()
		_ = cWrite.Close()
		return nil, fmt.Errorf("failed to create write pipe: %w", err)
	}

	// Files to pass to child
	// 0, 1, 2 are Stdin, Stdout, Stderr (handled by cmd/exec defaults or explicitly set)
	// 3: cRead
	// 4: cWrite
	files := []*os.File{cRead, cWrite}
	files = append(files, extraFiles...)

	cmd.ExtraFiles = files
	cmd.Stderr = os.Stderr // Pass through stderr for logging

	// Setup Env
	cmd.Env = os.Environ()
	for k, v := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	// Platform specific attributes (pdeathsig etc)
	configureCmd(cmd)

	return &Process{
		cmd:         cmd,
		ParentRead:  pRead,
		ParentWrite: pWrite,
		childRead:   cRead,
		childWrite:  cWrite,
	}, nil
}

// Start launches the process.
// It closes the *child* ends of the pipes upon success,
// so the parent only holds its own ends.
func (p *Process) Start() error {
	if err := p.cmd.Start(); err != nil {
		p.Close() // Cleanup parent pipes
		// Also cleanup child pipes since we failed
		_ = p.childRead.Close()
		_ = p.childWrite.Close()
		return err
	}

	// Close child-side FDs in parent to prevent leaks/deadlocks.
	// IMPORTANT: Only close the pipes WE created. Do NOT close extraFiles (like SHM)
	// because they might be shared across workers!

	_ = p.childRead.Close()
	_ = p.childWrite.Close()

	return nil
}

// Wait blocks until the process exits.
func (p *Process) Wait() error {
	return p.cmd.Wait()
}

// Signal sends a signal to the process.
func (p *Process) Signal(sig syscall.Signal) error {
	if p.cmd.Process != nil {
		return p.cmd.Process.Signal(sig)
	}
	return nil
}

// Kill forcibly terminates the process.
func (p *Process) Kill() error {
	if p.cmd.Process != nil {
		return p.cmd.Process.Kill()
	}
	return nil
}

// Close closes the parent-side IPC pipes.
func (p *Process) Close() {
	if p.ParentRead != nil {
		_ = p.ParentRead.Close()
	}
	if p.ParentWrite != nil {
		_ = p.ParentWrite.Close()
	}
}
