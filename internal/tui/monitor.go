package tui

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"gproc/pkg/types"
)

type MonitorDashboard struct {
	processes map[string]*types.Process
	metrics   map[string]*ProcessMetrics
	width     int
	height    int
}

type ProcessMetrics struct {
	CPU        float64
	Memory     int64
	Uptime     time.Duration
	Restarts   int
	Status     string
	LastUpdate time.Time
}

func NewMonitorDashboard() *MonitorDashboard {
	return &MonitorDashboard{
		processes: make(map[string]*types.Process),
		metrics:   make(map[string]*ProcessMetrics),
		width:     120,
		height:    30,
	}
}

func (md *MonitorDashboard) Start() {
	// Clear screen and hide cursor
	fmt.Print("\033[2J\033[?25l")
	
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			md.refresh()
		}
	}
}

func (md *MonitorDashboard) refresh() {
	// Move cursor to top-left
	fmt.Print("\033[H")
	
	// Update metrics
	md.updateMetrics()
	
	// Render dashboard
	md.renderHeader()
	md.renderProcessTable()
	md.renderSystemStats()
	md.renderFooter()
}

func (md *MonitorDashboard) renderHeader() {
	now := time.Now().Format("2006-01-02 15:04:05")
	title := "GProc Live Monitor"
	
	// Top border
	fmt.Print("┌" + strings.Repeat("─", md.width-2) + "┐\n")
	
	// Title line
	padding := (md.width - len(title) - len(now) - 4) / 2
	fmt.Printf("│ %s%s%s%s │\n", 
		title,
		strings.Repeat(" ", padding),
		now,
		strings.Repeat(" ", md.width-len(title)-len(now)-padding-4))
	
	// Separator
	fmt.Print("├" + strings.Repeat("─", md.width-2) + "┤\n")
}

func (md *MonitorDashboard) renderProcessTable() {
	// Table header
	fmt.Printf("│ %-20s %-10s %-8s %-10s %-8s %-12s %-20s │\n",
		"NAME", "STATUS", "CPU%", "MEMORY", "UPTIME", "RESTARTS", "LAST UPDATE")
	
	fmt.Print("├" + strings.Repeat("─", md.width-2) + "┤\n")
	
	// Sort processes by name
	var names []string
	for name := range md.processes {
		names = append(names, name)
	}
	sort.Strings(names)
	
	// Render process rows
	for _, name := range names {
		_ = md.processes[name]
		metrics := md.metrics[name]
		
		if metrics == nil {
			metrics = &ProcessMetrics{Status: "unknown"}
		}
		
		// Format values
		cpuStr := fmt.Sprintf("%.1f", metrics.CPU)
		memoryStr := md.formatMemory(metrics.Memory)
		uptimeStr := md.formatDuration(metrics.Uptime)
		restartsStr := fmt.Sprintf("%d", metrics.Restarts)
		lastUpdateStr := metrics.LastUpdate.Format("15:04:05")
		
		// Color coding for status
		statusColor := md.getStatusColor(metrics.Status)
		
		fmt.Printf("│ %-20s %s%-10s\033[0m %-8s %-10s %-8s %-12s %-20s │\n",
			name, statusColor, metrics.Status, cpuStr, memoryStr, 
			uptimeStr, restartsStr, lastUpdateStr)
	}
	
	// Fill remaining rows
	usedRows := len(names) + 3 // header + separator + title
	for i := usedRows; i < md.height-5; i++ {
		fmt.Printf("│%s│\n", strings.Repeat(" ", md.width-2))
	}
}

func (md *MonitorDashboard) renderSystemStats() {
	fmt.Print("├" + strings.Repeat("─", md.width-2) + "┤\n")
	
	// System statistics
	totalProcesses := len(md.processes)
	runningCount := 0
	stoppedCount := 0
	failedCount := 0
	
	for _, metrics := range md.metrics {
		switch metrics.Status {
		case "running":
			runningCount++
		case "stopped":
			stoppedCount++
		case "failed":
			failedCount++
		}
	}
	
	statsLine := fmt.Sprintf("Total: %d | Running: %d | Stopped: %d | Failed: %d",
		totalProcesses, runningCount, stoppedCount, failedCount)
	
	padding := md.width - len(statsLine) - 4
	fmt.Printf("│ %s%s │\n", statsLine, strings.Repeat(" ", padding))
}

func (md *MonitorDashboard) renderFooter() {
	fmt.Print("├" + strings.Repeat("─", md.width-2) + "┤\n")
	
	helpText := "Press 'q' to quit, 'r' to refresh, 'h' for help"
	padding := md.width - len(helpText) - 4
	fmt.Printf("│ %s%s │\n", helpText, strings.Repeat(" ", padding))
	
	// Bottom border
	fmt.Print("└" + strings.Repeat("─", md.width-2) + "┘\n")
}

func (md *MonitorDashboard) updateMetrics() {
	// This would integrate with the actual process manager
	// For now, simulate some data
	
	for name, process := range md.processes {
		if md.metrics[name] == nil {
			md.metrics[name] = &ProcessMetrics{}
		}
		
		metrics := md.metrics[name]
		metrics.Status = string(process.Status)
		metrics.Restarts = process.Restarts
		metrics.LastUpdate = time.Now()
		
		// Simulate changing metrics
		if process.Status == types.StatusRunning {
			metrics.CPU = 15.0 + float64(time.Now().Second()%10)
			metrics.Memory = 1024*1024*100 + int64(time.Now().Second()*1024*1024)
			metrics.Uptime = time.Since(process.StartTime)
		}
	}
}

func (md *MonitorDashboard) getStatusColor(status string) string {
	switch status {
	case "running":
		return "\033[32m" // Green
	case "stopped":
		return "\033[33m" // Yellow
	case "failed":
		return "\033[31m" // Red
	default:
		return "\033[37m" // White
	}
}

func (md *MonitorDashboard) formatMemory(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%dB", bytes)
	} else if bytes < 1024*1024 {
		return fmt.Sprintf("%.1fKB", float64(bytes)/1024)
	} else if bytes < 1024*1024*1024 {
		return fmt.Sprintf("%.1fMB", float64(bytes)/(1024*1024))
	} else {
		return fmt.Sprintf("%.1fGB", float64(bytes)/(1024*1024*1024))
	}
}

func (md *MonitorDashboard) formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	} else if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	} else if d < 24*time.Hour {
		return fmt.Sprintf("%dh", int(d.Hours()))
	} else {
		return fmt.Sprintf("%dd", int(d.Hours()/24))
	}
}

func (md *MonitorDashboard) AddProcess(process *types.Process) {
	md.processes[process.Name] = process
}

func (md *MonitorDashboard) RemoveProcess(name string) {
	delete(md.processes, name)
	delete(md.metrics, name)
}

func (md *MonitorDashboard) UpdateProcess(process *types.Process) {
	md.processes[process.Name] = process
}