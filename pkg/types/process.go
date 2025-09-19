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
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Command     string        `json:"command"`
	Args        []string      `json:"args"`
	WorkingDir  string        `json:"working_dir"`
	Status      ProcessStatus `json:"status"`
	PID         int           `json:"pid"`
	StartTime   time.Time     `json:"start_time"`
	Restarts    int           `json:"restarts"`
	AutoRestart bool          `json:"auto_restart"`
	MaxRestarts int           `json:"max_restarts"`
	LogFile     string        `json:"log_file"`
	Cmd         *exec.Cmd     `json:"-"`
}

type Config struct {
	Processes []Process `json:"processes"`
	LogDir    string    `json:"log_dir"`
}