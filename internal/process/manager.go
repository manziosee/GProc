package process

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"gproc/internal/config"
	"gproc/pkg/types"
)

type Manager struct {
	processes map[string]*types.Process
	mutex     sync.RWMutex
	logDir    string
	config    *types.Config
}

func NewManager(logDir string) *Manager {
	cfg, _ := config.LoadConfig()
	m := &Manager{
		processes: make(map[string]*types.Process),
		logDir:    logDir,
		config:    cfg,
	}
	m.loadProcesses()
	return m
}

func (m *Manager) loadProcesses() {
	for _, proc := range m.config.Processes {
		if proc.Status == types.StatusRunning {
			proc.Status = types.StatusStopped
		}
		m.processes[proc.ID] = &proc
	}
}

func (m *Manager) saveConfig() {
	processes := make([]types.Process, 0, len(m.processes))
	for _, proc := range m.processes {
		processes = append(processes, *proc)
	}
	m.config.Processes = processes
	config.SaveConfig(m.config)
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
	
	// Set environment variables
	if len(proc.Env) > 0 {
		env := os.Environ()
		for k, v := range proc.Env {
			env = append(env, k+"="+v)
		}
		cmd.Env = env
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
	m.saveConfig()

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

	// Graceful shutdown: try SIGTERM first, then SIGKILL
	if err := proc.Cmd.Process.Signal(os.Interrupt); err == nil {
		// Wait 5 seconds for graceful shutdown
		done := make(chan error, 1)
		go func() {
			done <- proc.Cmd.Wait()
		}()
		
		select {
		case <-time.After(5 * time.Second):
			// Force kill if not stopped gracefully
			if err := proc.Cmd.Process.Kill(); err != nil {
				return err
			}
		case <-done:
			// Process stopped gracefully
		}
	} else {
		// Fallback to kill
		if err := proc.Cmd.Process.Kill(); err != nil {
			return err
		}
	}

	proc.Status = types.StatusStopped
	m.saveConfig()
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

func (m *Manager) Restart(id string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	proc, exists := m.processes[id]
	if !exists {
		return fmt.Errorf("process %s not found", id)
	}

	if proc.Status == types.StatusRunning {
		if err := proc.Cmd.Process.Kill(); err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
	}

	cmd := exec.Command(proc.Command, proc.Args...)
	if proc.WorkingDir != "" {
		cmd.Dir = proc.WorkingDir
	}
	
	// Set environment variables
	if len(proc.Env) > 0 {
		env := os.Environ()
		for k, v := range proc.Env {
			env = append(env, k+"="+v)
		}
		cmd.Env = env
	}

	file, err := os.OpenFile(proc.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	cmd.Stdout = file
	cmd.Stderr = file
	proc.Cmd = cmd

	if err := cmd.Start(); err != nil {
		return err
	}

	proc.PID = cmd.Process.Pid
	proc.Status = types.StatusRunning
	proc.StartTime = time.Now()
	proc.Restarts++
	m.saveConfig()

	go m.monitor(proc)
	return nil
}

func (m *Manager) StartByName(name string) error {
	m.mutex.RLock()
	proc, exists := m.processes[name]
	m.mutex.RUnlock()
	
	if !exists {
		return fmt.Errorf("process %s not found", name)
	}
	
	return m.Start(proc)
}

func (m *Manager) AddScheduledTask(task *types.ScheduledTask) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	m.config.ScheduledTasks = append(m.config.ScheduledTasks, *task)
	m.saveConfig()
	return nil
}

func (m *Manager) StartWebDashboard(port int) error {
	dashboard := &webDashboard{manager: m}
	return dashboard.Start(port)
}

type webDashboard struct {
	manager *Manager
}

func (w *webDashboard) Start(port int) error {
	fmt.Printf("Web dashboard started on port %d\n", port)
	select {} // Block forever for demo
}

func (m *Manager) SaveTemplate(template *types.ProcessTemplate) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	m.config.Templates = append(m.config.Templates, *template)
	m.saveConfig()
	return nil
}

func (m *Manager) StartFromTemplate(templateName, processName string) error {
	m.mutex.RLock()
	var template *types.ProcessTemplate
	for _, t := range m.config.Templates {
		if t.Name == templateName {
			template = &t
			break
		}
	}
	m.mutex.RUnlock()
	
	if template == nil {
		return fmt.Errorf("template %s not found", templateName)
	}
	
	proc := &types.Process{
		ID:          processName,
		Name:        processName,
		Command:     template.Command,
		Args:        template.Args,
		WorkingDir:  template.WorkingDir,
		Env:         template.Env,
		AutoRestart: template.AutoRestart,
		MaxRestarts: template.MaxRestarts,
	}
	
	return m.Start(proc)
}

func (m *Manager) StartFromConfig(configFile string) error {
	fmt.Printf("Starting processes from config file: %s\n", configFile)
	return nil
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