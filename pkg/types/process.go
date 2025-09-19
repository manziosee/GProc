package types

import (
	"os/exec"
	"time"
)

type ProcessStatus string

const (
	StatusRunning ProcessStatus = "running"
	StatusStopped ProcessStatus = "stopped"
	StatusFailed  ProcessStatus = "failed"
)

type Process struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	Command       string            `json:"command"`
	Args          []string          `json:"args"`
	WorkingDir    string            `json:"working_dir"`
	Env           map[string]string `json:"env"`
	Status        ProcessStatus     `json:"status"`
	PID           int               `json:"pid"`
	StartTime     time.Time         `json:"start_time"`
	Restarts      int               `json:"restarts"`
	AutoRestart   bool              `json:"auto_restart"`
	MaxRestarts   int               `json:"max_restarts"`
	LogFile       string            `json:"log_file"`
	Group         string            `json:"group"`
	HealthCheck   *HealthCheck      `json:"health_check"`
	LogRotation   *LogRotation      `json:"log_rotation"`
	ResourceLimit *ResourceLimit    `json:"resource_limit"`
	Notifications *Notifications    `json:"notifications"`
	Dependencies  []string          `json:"dependencies"`
	Priority      int               `json:"priority"`
	BlueGreen     *BlueGreenConfig  `json:"blue_green"`
	Metrics       *ProcessMetrics   `json:"metrics"`
	Cmd           *exec.Cmd         `json:"-"`
}

type BlueGreenConfig struct {
	Enabled     bool   `json:"enabled"`
	ActiveSlot  string `json:"active_slot"`
	BluePort    int    `json:"blue_port"`
	GreenPort   int    `json:"green_port"`
	HealthPath  string `json:"health_path"`
}

type ProcessMetrics struct {
	CPUUsage    float64   `json:"cpu_usage"`
	MemoryUsage int64     `json:"memory_usage"`
	Uptime      time.Duration `json:"uptime"`
	LastCheck   time.Time `json:"last_check"`
	History     []MetricPoint `json:"history"`
}

type MetricPoint struct {
	Timestamp time.Time `json:"timestamp"`
	CPU       float64   `json:"cpu"`
	Memory    int64     `json:"memory"`
}

type Snapshot struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
	Processes []Process `json:"processes"`
	Config    Config    `json:"config"`
}

type Alert struct {
	ID          string    `json:"id"`
	ProcessID   string    `json:"process_id"`
	Type        string    `json:"type"`
	Message     string    `json:"message"`
	Severity    string    `json:"severity"`
	Timestamp   time.Time `json:"timestamp"`
	Acknowledged bool     `json:"acknowledged"`
}

type Plugin struct {
	Name        string            `json:"name"`
	Path        string            `json:"path"`
	Events      []string          `json:"events"`
	Config      map[string]string `json:"config"`
	Enabled     bool              `json:"enabled"`
}

type HealthCheck struct {
	URL      string        `json:"url"`
	Interval time.Duration `json:"interval"`
	Timeout  time.Duration `json:"timeout"`
	Retries  int           `json:"retries"`
}

type LogRotation struct {
	MaxSize  string `json:"max_size"`
	MaxFiles int    `json:"max_files"`
}

type ResourceLimit struct {
	MemoryMB int     `json:"memory_mb"`
	CPULimit float64 `json:"cpu_limit"`
}

type Notifications struct {
	Email string `json:"email"`
	Slack string `json:"slack"`
}

type ProcessGroup struct {
	Name      string   `json:"name"`
	Processes []string `json:"processes"`
}

type ProcessTemplate struct {
	Name        string            `json:"name"`
	Command     string            `json:"command"`
	Args        []string          `json:"args"`
	WorkingDir  string            `json:"working_dir"`
	Env         map[string]string `json:"env"`
	AutoRestart bool              `json:"auto_restart"`
	MaxRestarts int               `json:"max_restarts"`
}

type ScheduledTask struct {
	Name     string `json:"name"`
	Command  string `json:"command"`
	Args     []string `json:"args"`
	Cron     string `json:"cron"`
	NextRun  time.Time `json:"next_run"`
}

type Config struct {
	Processes      []Process         `json:"processes"`
	Groups         []ProcessGroup    `json:"groups"`
	Templates      []ProcessTemplate `json:"templates"`
	ScheduledTasks []ScheduledTask   `json:"scheduled_tasks"`
	Snapshots      []Snapshot        `json:"snapshots"`
	Alerts         []Alert           `json:"alerts"`
	Plugins        []Plugin          `json:"plugins"`
	LogDir         string            `json:"log_dir"`
	WebPort        int               `json:"web_port"`
	Cluster        *ClusterConfig    `json:"cluster"`
	Security       *SecurityConfig   `json:"security"`
	Monitoring     *MonitoringConfig `json:"monitoring"`
}

type ClusterConfig struct {
	Enabled    bool     `json:"enabled"`
	Mode       string   `json:"mode"` // "master" or "agent"
	MasterAddr string   `json:"master_addr"`
	Agents     []string `json:"agents"`
	TLSEnabled bool     `json:"tls_enabled"`
}

type SecurityConfig struct {
	RBACEnabled   bool              `json:"rbac_enabled"`
	Users         []User            `json:"users"`
	AuditLog      bool              `json:"audit_log"`
	TLSCert       string            `json:"tls_cert"`
	TLSKey        string            `json:"tls_key"`
	SecretsVault  string            `json:"secrets_vault"`
}

type MonitoringConfig struct {
	MetricsEnabled   bool   `json:"metrics_enabled"`
	MetricsDB        string `json:"metrics_db"`
	RetentionDays    int    `json:"retention_days"`
	AlertingEnabled  bool   `json:"alerting_enabled"`
	ProfilingEnabled bool   `json:"profiling_enabled"`
}

type User struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	Roles    []string `json:"roles"`
	Enabled  bool     `json:"enabled"`
}