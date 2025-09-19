package monitor

import (
	"fmt"
	"runtime"
)

type ResourceUsage struct {
	CPUPercent    float64 `json:"cpu_percent"`
	MemoryMB      float64 `json:"memory_mb"`
	MemoryPercent float64 `json:"memory_percent"`
}

func GetProcessResources(pid int) (*ResourceUsage, error) {
	if runtime.GOOS == "windows" {
		return &ResourceUsage{
			CPUPercent:    0,
			MemoryMB:      0,
			MemoryPercent: 0,
		}, nil
	}
	return &ResourceUsage{}, fmt.Errorf("resource monitoring not implemented for %s", runtime.GOOS)
}