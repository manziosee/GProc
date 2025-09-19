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
	Cmd           *exec.Cmd         `json:"-"`
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
	LogDir         string            `json:"log_dir"`
	WebPort        int               `json:"web_port"`
}