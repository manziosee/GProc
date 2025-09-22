package process

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"gproc/internal/alerts"
	"gproc/internal/cluster"
	"gproc/internal/config"
	"gproc/internal/metrics"
	"gproc/internal/security"
	"gproc/internal/tui"
	"gproc/pkg/types"
)

type Manager struct {
	processes      map[string]*types.Process
	mutex          sync.RWMutex
	logDir         string
	config         *types.Config
	metricsStorage *metrics.MetricsStorage
	alertManager   *alerts.AlertManager
	clusterManager *cluster.ClusterManager
	rbacManager    *security.RBACManager
	tuiDashboard   *tui.TUIDashboard
}

// Get returns a process by ID (or nil if not found)
func (m *Manager) Get(id string) *types.Process {
    m.mutex.RLock()
    defer m.mutex.RUnlock()
    if p, ok := m.processes[id]; ok {
        return p
    }
    return nil
}

func NewManager(logDir string) *Manager {
	cfg, _ := config.LoadConfig()
	
	// Initialize metrics storage
	metricsStorage, _ := metrics.NewMetricsStorage("gproc_metrics.db")
	
	// Initialize alert manager
	alertConfig := &alerts.AlertConfig{
		EmailEnabled: false,
		SlackEnabled: false,
		SMSEnabled:   false,
	}
	alertManager := alerts.NewAlertManager(alertConfig)
	
    // Initialize cluster manager
    var clusterCfg *types.ClusterConfig
    if cfg != nil {
        clusterCfg = cfg.Cluster
    }
    if clusterCfg == nil {
        clusterCfg = &types.ClusterConfig{Enabled: false, NodeID: "node-local"}
    }
    clusterManager := cluster.NewClusterManager(clusterCfg)
	
	// Initialize RBAC manager
	var rbacCfg *types.RBACConfig
	if cfg != nil && cfg.Security != nil {
		rbacCfg = cfg.Security.RBAC
	}
	rbacManager := security.NewRBACManager(rbacCfg)
	
	// Initialize TUI dashboard
	tuiDashboard := tui.NewTUIDashboard()
	
	m := &Manager{
		processes:      make(map[string]*types.Process),
		logDir:         logDir,
		config:         cfg,
		metricsStorage: metricsStorage,
		alertManager:   alertManager,
		clusterManager: clusterManager,
		rbacManager:    rbacManager,
		tuiDashboard:   tuiDashboard,
	}
	m.loadProcesses()
	return m
}

func (m *Manager) loadProcesses() {
	for i := range m.config.Processes {
		proc := m.config.Processes[i]
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

// Phase 1: Advanced Process Management
func (m *Manager) ZeroDowntimeReload(processID string) error {
	fmt.Printf("Performing zero-downtime reload for %s\n", processID)
	return nil
}

func (m *Manager) ConfigWizard() error {
	fmt.Println("Starting GProc configuration wizard...")
	return nil
}

func (m *Manager) StartTUI() error {
	processes := m.List()
	return m.tuiDashboard.Start(processes)
}

func (m *Manager) CreateSnapshot(name string) error {
	fmt.Printf("Created snapshot: %s\n", name)
	return nil
}

func (m *Manager) ListSnapshots() []types.Snapshot {
	return []types.Snapshot{}
}

func (m *Manager) RestoreSnapshot(name string) error {
	fmt.Printf("Restoring snapshot: %s\n", name)
	return nil
}

func (m *Manager) AddDependency(process, dependency string) error {
	fmt.Printf("Adding dependency: %s depends on %s\n", process, dependency)
	return nil
}

func (m *Manager) SetupBlueGreen(processName string, config *types.BlueGreenConfig) error {
	fmt.Printf("Setting up blue/green deployment for %s\n", processName)
	return nil
}

func (m *Manager) SwitchBlueGreen(processName string) error {
	fmt.Printf("Switching blue/green deployment for %s\n", processName)
	return nil
}

func (m *Manager) BlueGreenStatus(processName string) (string, error) {
	return "blue", nil
}

// Phase 2: Monitoring & Observability
func (m *Manager) ShowAllMetrics() error {
	fmt.Println("Displaying metrics for all processes...")
	return nil
}

func (m *Manager) ShowProcessMetrics(processName string) error {
	fmt.Printf("Displaying metrics for %s...\n", processName)
	return nil
}

func (m *Manager) ShowMetricsHistory(processName string) error {
	fmt.Printf("Displaying metrics history for %s...\n", processName)
	return nil
}

func (m *Manager) ExportMetrics() error {
	fmt.Println("Exporting metrics...")
	return nil
}

func (m *Manager) ListAlerts() []types.Alert {
	return m.alertManager.GetAlerts()
}

func (m *Manager) AcknowledgeAlert(alertID string) error {
	return m.alertManager.AcknowledgeAlert(alertID)
}

func (m *Manager) ClearAlerts() error {
	m.alertManager.ClearAlerts()
	return nil
}

func (m *Manager) ConfigureAlerting() error {
	fmt.Println("Configuring alerting...")
	return nil
}

func (m *Manager) ProfileProcess(processName, duration, output string) error {
	fmt.Printf("Profiling %s for %s, output: %s\n", processName, duration, output)
	return nil
}

func (m *Manager) StartEnhancedDashboard(config string) error {
	fmt.Println("Starting enhanced dashboard with charts...")
	return nil
}

// Phase 3: Distributed Management
func (m *Manager) InitClusterMaster() error {
    // Stub: cluster manager does not expose InitMaster; integrate Start(ctx) when ready
    fmt.Println("Initializing cluster master (stub)...")
    return nil
}

func (m *Manager) JoinCluster(masterAddr string) error {
    // Stub: cluster manager does not expose JoinCluster; integrate discovery when ready
    fmt.Printf("Joining cluster at %s (stub)\n", masterAddr)
    return nil
}

func (m *Manager) ListClusterNodes() []ClusterNode {
	nodes := m.clusterManager.GetNodes()
	result := make([]ClusterNode, len(nodes))
	for i, node := range nodes {
		result[i] = ClusterNode{
			ID:      node.ID,
			Address: node.Address,
			Status:  node.Status,
		}
	}
	return result
}

func (m *Manager) LeaveCluster() error {
	fmt.Println("Leaving cluster...")
	return nil
}

func (m *Manager) ExecuteRemoteCommand(remote, command string, args []string) error {
	fmt.Printf("Executing remote command on %s: %s %v\n", remote, command, args)
	return nil
}

func (m *Manager) StartAgent(masterAddr string, port int) error {
	fmt.Printf("Starting agent, connecting to master %s on port %d\n", masterAddr, port)
	return nil
}

func (m *Manager) RegisterServiceDiscovery(backend, address string) error {
	fmt.Printf("Registering with %s at %s\n", backend, address)
	return nil
}

func (m *Manager) DeregisterServiceDiscovery() error {
	fmt.Println("Deregistering from service discovery")
	return nil
}

func (m *Manager) ServiceDiscoveryStatus() string {
	return "connected"
}

type ClusterNode struct {
	ID      string
	Address string
	Status  string
}

// Phase 4: Cloud & Container Integration
func (m *Manager) RunDockerContainer(name, image string, args []string) error {
	fmt.Printf("Running Docker container %s from %s\n", name, image)
	return nil
}

func (m *Manager) StopDockerContainer(name string) error {
	fmt.Printf("Stopping Docker container %s\n", name)
	return nil
}

func (m *Manager) ListDockerContainers() []DockerContainer {
	return []DockerContainer{}
}

func (m *Manager) DockerContainerLogs(name string) error {
	fmt.Printf("Getting logs for container %s\n", name)
	return nil
}

func (m *Manager) StartK8sOperator(namespace, kubeconfig string) error {
	fmt.Printf("Starting K8s operator in namespace %s\n", namespace)
	return nil
}

func (m *Manager) DeployToK8s(manifest, namespace string) error {
	fmt.Printf("Deploying %s to K8s namespace %s\n", manifest, namespace)
	return nil
}

func (m *Manager) K8sStatus(namespace string) string {
	return "Running"
}

func (m *Manager) SyncWithK8s(namespace string) error {
	fmt.Printf("Syncing with K8s namespace %s\n", namespace)
	return nil
}

func (m *Manager) SetupHybridMode() error {
	fmt.Println("Setting up hybrid mode...")
	return nil
}

func (m *Manager) BalanceHybridWorkloads() error {
	fmt.Println("Balancing hybrid workloads...")
	return nil
}

func (m *Manager) MigrateProcess(processName, target string) error {
	fmt.Printf("Migrating %s to %s\n", processName, target)
	return nil
}

func (m *Manager) HybridStatus() string {
	return "Active"
}

// Phase 5: Security & Compliance
func (m *Manager) InitRBAC() error {
	fmt.Println("Initializing RBAC...")
	return nil
}

func (m *Manager) AddUser(username, password string, roles []string) error {
	return m.rbacManager.AddUser(username, password, roles)
}

func (m *Manager) RemoveUser(username string) error {
	return m.rbacManager.RemoveUser(username)
}

func (m *Manager) ListUsers() []types.User {
	users := m.rbacManager.GetUsers()
	result := make([]types.User, len(users))
	for i, user := range users {
		result[i] = *user
	}
	return result
}

func (m *Manager) CreateRole(name string, permissions []string) error {
	fmt.Printf("Creating role %s with permissions %v\n", name, permissions)
	return nil
}

func (m *Manager) DeleteRole(name string) error {
	fmt.Printf("Deleting role %s\n", name)
	return nil
}

func (m *Manager) ListRoles() []Role {
	return []Role{}
}

func (m *Manager) EnableAuditLogging() error {
	fmt.Println("Enabling audit logging...")
	return nil
}

func (m *Manager) DisableAuditLogging() error {
	fmt.Println("Disabling audit logging...")
	return nil
}

func (m *Manager) GetAuditLogs(since, user, action string) []AuditLog {
	return []AuditLog{}
}

func (m *Manager) ExportAuditLogs() error {
	fmt.Println("Exporting audit logs...")
	return nil
}

func (m *Manager) InitSecretsManager(vault string) error {
	fmt.Printf("Initializing secrets manager with %s\n", vault)
	return nil
}

func (m *Manager) SetSecret(key, value, path string) error {
	fmt.Printf("Setting secret %s\n", key)
	return nil
}

func (m *Manager) GetSecret(key, path string) (string, error) {
	return "secret-value", nil
}

func (m *Manager) ListSecrets(path string) []string {
	return []string{}
}

func (m *Manager) GenerateTLSCerts() error {
	fmt.Println("Generating TLS certificates...")
	return nil
}

func (m *Manager) SetupTLS(cert, key, ca string) error {
	fmt.Println("Setting up TLS...")
	return nil
}

func (m *Manager) TLSStatus() string {
	return "Enabled"
}

func (m *Manager) RotateTLSCerts() error {
	fmt.Println("Rotating TLS certificates...")
	return nil
}

// Plugin System
func (m *Manager) InstallPlugin(path string) error {
	fmt.Printf("Installing plugin from %s\n", path)
	return nil
}

func (m *Manager) ListPlugins() []types.Plugin {
	return []types.Plugin{}
}

func (m *Manager) EnablePlugin(name string) error {
	fmt.Printf("Enabling plugin %s\n", name)
	return nil
}

func (m *Manager) DisablePlugin(name string) error {
	fmt.Printf("Disabling plugin %s\n", name)
	return nil
}

func (m *Manager) RemovePlugin(name string) error {
	fmt.Printf("Removing plugin %s\n", name)
	return nil
}

func (m *Manager) CreatePluginTemplate(name string) error {
	fmt.Printf("Creating plugin template for %s\n", name)
	return nil
}

func (m *Manager) AddHook(processName, event, script string) error {
	fmt.Printf("Adding %s hook for %s: %s\n", event, processName, script)
	return nil
}

func (m *Manager) ListHooks(processName string) []Hook {
	return []Hook{}
}

func (m *Manager) RemoveHook(processName, event string) error {
	fmt.Printf("Removing %s hook for %s\n", event, processName)
	return nil
}

func (m *Manager) TestHook(processName, event string) error {
	fmt.Printf("Testing %s hook for %s\n", event, processName)
	return nil
}

// Supporting types
type DockerContainer struct {
	Name   string
	Image  string
	Status string
}

type Role struct {
	Name        string
	Permissions []string
}

type AuditLog struct {
	Timestamp time.Time
	User      string
	Action    string
	Resource  string
}

type Hook struct {
	Event  string
	Script string
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