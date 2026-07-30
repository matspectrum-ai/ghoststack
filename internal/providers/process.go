package providers

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

type ProcessState string

const (
	ProcessStopped ProcessState = "stopped"
	ProcessRunning ProcessState = "running"
	ProcessFailed  ProcessState = "failed"
)

type ProcessManager struct {
	mu     sync.RWMutex
	cmd    *exec.Cmd
	state  ProcessState
	pid    int
	cancel context.CancelFunc
	done   chan struct{}
	err    error
	stopFn func()
	stdout io.Writer
	stderr io.Writer
}

type ProcessConfig struct {
	Name   string
	Args   []string
	Env    []string
	Stdout io.Writer
	Stderr io.Writer
}

func NewProcessManager() *ProcessManager {
	return &ProcessManager{
		state: ProcessStopped,
	}
}

func (pm *ProcessManager) Start(ctx context.Context, cfg ProcessConfig) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if pm.state == ProcessRunning {
		return fmt.Errorf("process already running (pid %d)", pm.pid)
	}

	ctx, cancel := context.WithCancel(ctx)
	cmd := exec.CommandContext(ctx, cfg.Name, cfg.Args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Env = append(os.Environ(), cfg.Env...)

	if cfg.Stdout != nil {
		cmd.Stdout = cfg.Stdout
	}
	if cfg.Stderr != nil {
		cmd.Stderr = cfg.Stderr
	}

	if err := cmd.Start(); err != nil {
		cancel()
		return fmt.Errorf("start process %s: %w", cfg.Name, err)
	}

	pm.cmd = cmd
	pm.pid = cmd.Process.Pid
	pm.state = ProcessRunning
	pm.cancel = cancel
	pm.err = nil
	pm.done = make(chan struct{})

	go pm.wait(cmd)

	return nil
}

func (pm *ProcessManager) wait(cmd *exec.Cmd) {
	err := cmd.Wait()

	pm.mu.Lock()
	if pm.state != ProcessStopped {
		pm.state = ProcessFailed
	}
	pm.err = err
	pm.pid = 0
	close(pm.done)
	pm.mu.Unlock()
}

func (pm *ProcessManager) Stop() error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if pm.state == ProcessStopped {
		return fmt.Errorf("process not running")
	}

	if pm.cancel != nil {
		pm.cancel()
	}

	done := pm.done
	if done != nil {
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			if pm.cmd != nil && pm.cmd.Process != nil {
				pm.cmd.Process.Signal(syscall.SIGKILL)
			}
		}
	}

	pm.state = ProcessStopped
	return nil
}

func (pm *ProcessManager) Signal(sig os.Signal) error {
	pm.mu.RLock()
	cmd := pm.cmd
	pm.mu.RUnlock()

	if cmd == nil || cmd.Process == nil {
		return fmt.Errorf("process not running")
	}
	return cmd.Process.Signal(sig)
}

func (pm *ProcessManager) PID() int {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.pid
}

func (pm *ProcessManager) Running() bool {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.state == ProcessRunning
}

func (pm *ProcessManager) State() ProcessState {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.state
}

func (pm *ProcessManager) Err() error {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.err
}
