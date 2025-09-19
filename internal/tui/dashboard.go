package tui

import (
	"fmt"
	"time"

	"gproc/pkg/types"
)

type TUIDashboard struct {
	processes []*types.Process
	running   bool
}

func NewTUIDashboard() *TUIDashboard {
	return &TUIDashboard{
		processes: []*types.Process{},
		running:   false,
	}
}

func (t *TUIDashboard) Start(processes []*types.Process) error {
	t.processes = processes
	t.running = true
	
	// Clear screen and hide cursor
	fmt.Print("\033[2J\033[H\033[?25l")
	
	for t.running {
		t.render()
		time.Sleep(1 * time.Second)
	}
	
	// Show cursor
	fmt.Print("\033[?25h")
	return nil
}

func (t *TUIDashboard) render() {
	// Move cursor to top
	fmt.Print("\033[H")
	
	// Header
	fmt.Println("┌─────────────────────────────────────────────────────────────────────┐")
	fmt.Println("│                        GProc Live Dashboard                        │")
	fmt.Println("├─────────────┬──────────┬─────────┬──────────┬─────────────────────┤")
	fmt.Println("│    NAME     │  STATUS  │   PID   │ RESTARTS │       UPTIME        │")
	fmt.Println("├─────────────┼──────────┼─────────┼──────────┼─────────────────────┤")
	
	// Process rows
	for _, proc := range t.processes {
		uptime := ""
		if proc.Status == types.StatusRunning {
			uptime = time.Since(proc.StartTime).Round(time.Second).String()
		}
		
		statusColor := getStatusColor(proc.Status)
		fmt.Printf("│ %-11s │ %s%-8s\033[0m │ %-7d │ %-8d │ %-19s │\n",
			truncate(proc.Name, 11),
			statusColor, proc.Status,
			proc.PID, proc.Restarts, uptime)
	}
	
	// Footer
	fmt.Println("└─────────────┴──────────┴─────────┴──────────┴─────────────────────┘")
	fmt.Printf("Processes: %d | Running: %d | Press Ctrl+C to exit\n", 
		len(t.processes), countRunning(t.processes))
	
	// Clear rest of screen
	fmt.Print("\033[J")
}

func getStatusColor(status types.ProcessStatus) string {
	switch status {
	case types.StatusRunning:
		return "\033[32m" // Green
	case types.StatusStopped:
		return "\033[33m" // Yellow
	case types.StatusFailed:
		return "\033[31m" // Red
	default:
		return "\033[0m"  // Default
	}
}

func truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length-3] + "..."
}

func countRunning(processes []*types.Process) int {
	count := 0
	for _, proc := range processes {
		if proc.Status == types.StatusRunning {
			count++
		}
	}
	return count
}