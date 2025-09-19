package process

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"gproc/pkg/types"
)

type Manager struct {
	processes map[string]*types.Process
	mutex     sync.RWMutex
	logDir    string
}

func NewManager(logDir string) *Manager {
	return &Manager{
		processes: make(map[string]*types.Process),
		logDir:    logDir,
	}
}

func (m *Manager) Start(proc *types.Process) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if existing, exists := m.processes[proc.ID]; exists && existing.Status == types.StatusRunning {
		return fmt.Errorf("process %s is already running", proc.ID)
	}

	cmd := exec.Command(proc.Command, proc.Args...)
	if proc.WorkingDir != "" {
		cmd.Dir = proc.WorkingDir
	}

	logFile := filepath.Join(m.logDir, proc.ID+".log")
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	cmd.Stdout = file
	cmd.Stderr = file
	proc.LogFile = logFile
	proc.Cmd = cmd

	if err := cmd.Start(); err != nil {
		return err
	}

	proc.PID = cmd.Process.Pid
	proc.Status = types.StatusRunning
	proc.StartTime = time.Now()
	m.processes[proc.ID] = proc

	go m.monitor(proc)
	return nil
}

func (m *Manager) Stop(id string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	proc, exists := m.processes[id]
	if !exists {
		return fmt.Errorf("process %s not found", id)
	}

	if proc.Status != types.StatusRunning {
		return fmt.Errorf("process %s is not running", id)
	}

	if err := proc.Cmd.Process.Kill(); err != nil {
		return err
	}

	proc.Status = types.StatusStopped
	return nil
}

func (m *Manager) List() []*types.Process {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	processes := make([]*types.Process, 0, len(m.processes))
	for _, proc := range m.processes {
		processes = append(processes, proc)
	}
	return processes
}

func (m *Manager) monitor(proc *types.Process) {
	proc.Cmd.Wait()
	
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if proc.Status == types.StatusStopped {
		return
	}

	proc.Status = types.StatusFailed
	
	if proc.AutoRestart && proc.Restarts < proc.MaxRestarts {
		proc.Restarts++
		go func() {
			time.Sleep(2 * time.Second)
			m.Start(proc)
		}()
	}
}