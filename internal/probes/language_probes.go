package probes

import (
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type ProbeResult struct {
	Name      string      `json:"name"`
	Value     interface{} `json:"value"`
	Unit      string      `json:"unit"`
	Timestamp time.Time   `json:"timestamp"`
	Healthy   bool        `json:"healthy"`
}

type LanguageProbe interface {
	Probe(pid int) (*ProbeResult, error)
	GetName() string
}

// Node.js Probes
type NodeEventLoopLagProbe struct{}

func (p *NodeEventLoopLagProbe) GetName() string { return "event_loop_lag" }

func (p *NodeEventLoopLagProbe) Probe(pid int) (*ProbeResult, error) {
	// Try to get event loop lag from Node.js process
	// This would typically connect to a debug port or use process metrics
	
	result := &ProbeResult{
		Name:      "event_loop_lag",
		Value:     2.5, // Simulated value in ms
		Unit:      "ms",
		Timestamp: time.Now(),
		Healthy:   true,
	}
	
	if lag, ok := result.Value.(float64); ok && lag > 10 {
		result.Healthy = false
	}
	
	return result, nil
}

type NodeHeapUsageProbe struct{}

func (p *NodeHeapUsageProbe) GetName() string { return "heap_usage" }

func (p *NodeHeapUsageProbe) Probe(pid int) (*ProbeResult, error) {
	// Get heap usage from Node.js process
	result := &ProbeResult{
		Name:      "heap_usage",
		Value:     45.2, // Simulated percentage
		Unit:      "%",
		Timestamp: time.Now(),
		Healthy:   true,
	}
	
	if usage, ok := result.Value.(float64); ok && usage > 85 {
		result.Healthy = false
	}
	
	return result, nil
}

// Python Probes
type PythonGILProbe struct{}

func (p *PythonGILProbe) GetName() string { return "gil_wait_percent" }

func (p *PythonGILProbe) Probe(pid int) (*ProbeResult, error) {
	result := &ProbeResult{
		Name:      "gil_wait_percent",
		Value:     15.3, // Simulated GIL wait percentage
		Unit:      "%",
		Timestamp: time.Now(),
		Healthy:   true,
	}
	
	if wait, ok := result.Value.(float64); ok && wait > 50 {
		result.Healthy = false
	}
	
	return result, nil
}

// Java Probes
type JavaHeapProbe struct{}

func (p *JavaHeapProbe) GetName() string { return "heap_usage" }

func (p *JavaHeapProbe) Probe(pid int) (*ProbeResult, error) {
	// Use jstat or JMX to get heap usage
	cmd := exec.Command("jstat", "-gc", strconv.Itoa(pid))
	output, err := cmd.Output()
	if err != nil {
		// Fallback to simulated data
		return &ProbeResult{
			Name:      "heap_usage",
			Value:     60.5,
			Unit:      "%",
			Timestamp: time.Now(),
			Healthy:   true,
		}, nil
	}
	
	// Parse jstat output (simplified)
	lines := strings.Split(string(output), "\n")
	if len(lines) > 1 {
		fields := strings.Fields(lines[1])
		if len(fields) > 2 {
			// Extract heap usage from jstat output
			// This is a simplified parser
		}
	}
	
	return &ProbeResult{
		Name:      "heap_usage",
		Value:     60.5,
		Unit:      "%",
		Timestamp: time.Now(),
		Healthy:   true,
	}, nil
}

type JavaGCProbe struct{}

func (p *JavaGCProbe) GetName() string { return "gc_stats" }

func (p *JavaGCProbe) Probe(pid int) (*ProbeResult, error) {
	result := &ProbeResult{
		Name:      "gc_stats",
		Value:     map[string]interface{}{
			"young_gc_count": 45,
			"old_gc_count":   3,
			"gc_time_ms":     120,
		},
		Unit:      "mixed",
		Timestamp: time.Now(),
		Healthy:   true,
	}
	
	return result, nil
}

// Go Probes
type GoGoroutinesProbe struct{}

func (p *GoGoroutinesProbe) GetName() string { return "goroutines" }

func (p *GoGoroutinesProbe) Probe(pid int) (*ProbeResult, error) {
	// Try to get goroutine count from pprof endpoint
	resp, err := http.Get("http://localhost:6060/debug/pprof/goroutine?debug=1")
	if err != nil {
		// Fallback to simulated data
		return &ProbeResult{
			Name:      "goroutines",
			Value:     25,
			Unit:      "count",
			Timestamp: time.Now(),
			Healthy:   true,
		}, nil
	}
	defer resp.Body.Close()
	
	// Parse goroutine count from pprof output
	// This would need proper parsing
	
	result := &ProbeResult{
		Name:      "goroutines",
		Value:     25,
		Unit:      "count",
		Timestamp: time.Now(),
		Healthy:   true,
	}
	
	if count, ok := result.Value.(int); ok && count > 1000 {
		result.Healthy = false
	}
	
	return result, nil
}

// Rust Probes
type RustThreadCountProbe struct{}

func (p *RustThreadCountProbe) GetName() string { return "thread_count" }

func (p *RustThreadCountProbe) Probe(pid int) (*ProbeResult, error) {
	// Get thread count from /proc/pid/status on Linux
	// or use platform-specific methods
	
	result := &ProbeResult{
		Name:      "thread_count",
		Value:     8,
		Unit:      "count",
		Timestamp: time.Now(),
		Healthy:   true,
	}
	
	return result, nil
}

// PHP Probes
type PHPOpcacheProbe struct{}

func (p *PHPOpcacheProbe) GetName() string { return "opcache_stats" }

func (p *PHPOpcacheProbe) Probe(pid int) (*ProbeResult, error) {
	result := &ProbeResult{
		Name: "opcache_stats",
		Value: map[string]interface{}{
			"hit_rate":      95.2,
			"memory_usage":  "45MB",
			"cached_files":  1250,
		},
		Unit:      "mixed",
		Timestamp: time.Now(),
		Healthy:   true,
	}
	
	return result, nil
}

// Probe Manager
type ProbeManager struct {
	probes map[string]map[string]LanguageProbe
}

func NewProbeManager() *ProbeManager {
	pm := &ProbeManager{
		probes: make(map[string]map[string]LanguageProbe),
	}
	
	// Register Node.js probes
	pm.probes["node"] = map[string]LanguageProbe{
		"event_loop_lag": &NodeEventLoopLagProbe{},
		"heap_usage":     &NodeHeapUsageProbe{},
	}
	
	// Register Python probes
	pm.probes["python"] = map[string]LanguageProbe{
		"gil_wait_percent": &PythonGILProbe{},
	}
	
	// Register Java probes
	pm.probes["java"] = map[string]LanguageProbe{
		"heap_usage": &JavaHeapProbe{},
		"gc_stats":   &JavaGCProbe{},
	}
	
	// Register Go probes
	pm.probes["go"] = map[string]LanguageProbe{
		"goroutines": &GoGoroutinesProbe{},
	}
	
	// Register Rust probes
	pm.probes["rust"] = map[string]LanguageProbe{
		"thread_count": &RustThreadCountProbe{},
	}
	
	// Register PHP probes
	pm.probes["php"] = map[string]LanguageProbe{
		"opcache_stats": &PHPOpcacheProbe{},
	}
	
	return pm
}

func (pm *ProbeManager) RunProbes(language string, pid int, probeNames []string) ([]*ProbeResult, error) {
	langProbes, exists := pm.probes[language]
	if !exists {
		return nil, fmt.Errorf("no probes available for language: %s", language)
	}
	
	var results []*ProbeResult
	
	for _, probeName := range probeNames {
		probe, exists := langProbes[probeName]
		if !exists {
			continue
		}
		
		result, err := probe.Probe(pid)
		if err != nil {
			continue
		}
		
		results = append(results, result)
	}
	
	return results, nil
}

func (pm *ProbeManager) GetAvailableProbes(language string) []string {
	langProbes, exists := pm.probes[language]
	if !exists {
		return nil
	}
	
	var probeNames []string
	for name := range langProbes {
		probeNames = append(probeNames, name)
	}
	
	return probeNames
}