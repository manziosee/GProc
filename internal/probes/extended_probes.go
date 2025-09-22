package probes

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"gproc/pkg/types"
)

type ExtendedProbeManager struct {
	config *types.ExtendedProbeConfig
	probes map[string]ExtendedLanguageProbe
}

type ExtendedLanguageProbe interface {
	Detect(processPath string) bool
	GetMetrics(pid int) (*types.ProbeResult, error)
	GetHealthStatus(pid int) bool
}

type DotNetProbe struct{}
type RubyProbe struct{}
type ErlangProbe struct{}

func NewExtendedProbeManager(config *types.ExtendedProbeConfig) *ExtendedProbeManager {
	manager := &ExtendedProbeManager{
		config: config,
		probes: make(map[string]ExtendedLanguageProbe),
	}
	
	// Register extended language probes
	manager.probes["dotnet"] = &DotNetProbe{}
	manager.probes["ruby"] = &RubyProbe{}
	manager.probes["erlang"] = &ErlangProbe{}
	
	return manager
}

func (e *ExtendedProbeManager) ProbeProcess(pid int, processPath string) ([]*types.ProbeResult, error) {
	var results []*types.ProbeResult
	
	for language, probe := range e.probes {
		if probe.Detect(processPath) {
			metrics, err := probe.GetMetrics(pid)
			if err != nil {
				continue
			}
			
			metrics.Name = fmt.Sprintf("%s_metrics", language)
			metrics.Timestamp = time.Now()
			metrics.Healthy = probe.GetHealthStatus(pid)
			
			results = append(results, metrics)
		}
	}
	
	return results, nil
}

// .NET Probe Implementation
func (d *DotNetProbe) Detect(processPath string) bool {
	// Check if process is .NET application
	return strings.Contains(processPath, ".exe") || 
		   strings.Contains(processPath, "dotnet") ||
		   strings.Contains(processPath, ".dll")
}

func (d *DotNetProbe) GetMetrics(pid int) (*types.ProbeResult, error) {
	// Get .NET specific metrics using dotnet-counters or performance counters
	metrics := &types.ProbeResult{
		Value: map[string]interface{}{
			"gc_collections_gen0": d.getGCCollections(pid, 0),
			"gc_collections_gen1": d.getGCCollections(pid, 1),
			"gc_collections_gen2": d.getGCCollections(pid, 2),
			"heap_size_mb":        d.getHeapSize(pid),
			"thread_pool_threads": d.getThreadPoolThreads(pid),
			"exceptions_thrown":   d.getExceptionsThrown(pid),
			"jit_compiled_methods": d.getJITCompiledMethods(pid),
		},
		Unit: "various",
	}
	
	return metrics, nil
}

func (d *DotNetProbe) GetHealthStatus(pid int) bool {
	// Check if .NET process is responsive
	return d.getExceptionsThrown(pid) < 100 // Threshold for exceptions
}

func (d *DotNetProbe) getGCCollections(pid int, generation int) int {
	// Simulate getting GC collection count for specific generation
	cmd := exec.Command("dotnet-counters", "query", "--process-id", strconv.Itoa(pid), "--counter", fmt.Sprintf("System.Runtime:gc-gen-%d-collections", generation))
	output, err := cmd.Output()
	if err != nil {
		return 0
	}
	
	// Parse output (simplified)
	if count, err := strconv.Atoi(strings.TrimSpace(string(output))); err == nil {
		return count
	}
	return 0
}

func (d *DotNetProbe) getHeapSize(pid int) float64 {
	// Get managed heap size in MB
	cmd := exec.Command("dotnet-counters", "query", "--process-id", strconv.Itoa(pid), "--counter", "System.Runtime:gc-heap-size")
	output, err := cmd.Output()
	if err != nil {
		return 0.0
	}
	
	if size, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64); err == nil {
		return size / (1024 * 1024) // Convert to MB
	}
	return 0.0
}

func (d *DotNetProbe) getThreadPoolThreads(pid int) int {
	// Get thread pool thread count
	return 25 // Mock value
}

func (d *DotNetProbe) getExceptionsThrown(pid int) int {
	// Get exceptions thrown count
	return 5 // Mock value
}

func (d *DotNetProbe) getJITCompiledMethods(pid int) int {
	// Get JIT compiled methods count
	return 1500 // Mock value
}

// Ruby Probe Implementation
func (r *RubyProbe) Detect(processPath string) bool {
	return strings.Contains(processPath, "ruby") || 
		   strings.Contains(processPath, ".rb") ||
		   strings.Contains(processPath, "rails")
}

func (r *RubyProbe) GetMetrics(pid int) (*types.ProbeResult, error) {
	metrics := &types.ProbeResult{
		Value: map[string]interface{}{
			"object_count":        r.getObjectCount(pid),
			"gc_runs":            r.getGCRuns(pid),
			"heap_slots_live":    r.getHeapSlotsLive(pid),
			"heap_slots_free":    r.getHeapSlotsFree(pid),
			"thread_count":       r.getThreadCount(pid),
			"fiber_count":        r.getFiberCount(pid),
			"memory_usage_mb":    r.getMemoryUsage(pid),
		},
		Unit: "various",
	}
	
	return metrics, nil
}

func (r *RubyProbe) GetHealthStatus(pid int) bool {
	// Check Ruby process health
	return r.getObjectCount(pid) < 1000000 // Object count threshold
}

func (r *RubyProbe) getObjectCount(pid int) int {
	// Get Ruby object count using ObjectSpace
	return 50000 // Mock value
}

func (r *RubyProbe) getGCRuns(pid int) int {
	// Get GC run count
	return 150 // Mock value
}

func (r *RubyProbe) getHeapSlotsLive(pid int) int {
	// Get live heap slots
	return 25000 // Mock value
}

func (r *RubyProbe) getHeapSlotsFree(pid int) int {
	// Get free heap slots
	return 5000 // Mock value
}

func (r *RubyProbe) getThreadCount(pid int) int {
	// Get Ruby thread count
	return 8 // Mock value
}

func (r *RubyProbe) getFiberCount(pid int) int {
	// Get Ruby fiber count
	return 12 // Mock value
}

func (r *RubyProbe) getMemoryUsage(pid int) float64 {
	// Get Ruby memory usage in MB
	return 125.5 // Mock value
}

// Erlang/Elixir Probe Implementation
func (e *ErlangProbe) Detect(processPath string) bool {
	return strings.Contains(processPath, "erl") || 
		   strings.Contains(processPath, "elixir") ||
		   strings.Contains(processPath, "beam") ||
		   strings.Contains(processPath, ".ex")
}

func (e *ErlangProbe) GetMetrics(pid int) (*types.ProbeResult, error) {
	metrics := &types.ProbeResult{
		Value: map[string]interface{}{
			"process_count":       e.getProcessCount(pid),
			"message_queue_len":   e.getMessageQueueLength(pid),
			"memory_total_mb":     e.getMemoryTotal(pid),
			"memory_processes_mb": e.getMemoryProcesses(pid),
			"memory_system_mb":    e.getMemorySystem(pid),
			"reductions":          e.getReductions(pid),
			"run_queue_length":    e.getRunQueueLength(pid),
			"io_input_bytes":      e.getIOInput(pid),
			"io_output_bytes":     e.getIOOutput(pid),
		},
		Unit: "various",
	}
	
	return metrics, nil
}

func (e *ErlangProbe) GetHealthStatus(pid int) bool {
	// Check Erlang VM health
	return e.getProcessCount(pid) < 100000 && // Process count threshold
		   e.getRunQueueLength(pid) < 1000    // Run queue threshold
}

func (e *ErlangProbe) getProcessCount(pid int) int {
	// Get Erlang process count using erlang:system_info(process_count)
	return 2500 // Mock value
}

func (e *ErlangProbe) getMessageQueueLength(pid int) int {
	// Get total message queue length across all processes
	return 150 // Mock value
}

func (e *ErlangProbe) getMemoryTotal(pid int) float64 {
	// Get total memory usage in MB
	return 256.0 // Mock value
}

func (e *ErlangProbe) getMemoryProcesses(pid int) float64 {
	// Get memory used by processes in MB
	return 180.0 // Mock value
}

func (e *ErlangProbe) getMemorySystem(pid int) float64 {
	// Get memory used by system in MB
	return 76.0 // Mock value
}

func (e *ErlangProbe) getReductions(pid int) int64 {
	// Get total reductions (work units)
	return 1500000 // Mock value
}

func (e *ErlangProbe) getRunQueueLength(pid int) int {
	// Get scheduler run queue length
	return 25 // Mock value
}

func (e *ErlangProbe) getIOInput(pid int) int64 {
	// Get total IO input bytes
	return 1024000 // Mock value
}

func (e *ErlangProbe) getIOOutput(pid int) int64 {
	// Get total IO output bytes
	return 512000 // Mock value
}

func (e *ExtendedProbeManager) GetSupportedLanguages() []string {
	var languages []string
	for lang := range e.probes {
		languages = append(languages, lang)
	}
	return languages
}

func (e *ExtendedProbeManager) IsLanguageSupported(language string) bool {
	_, exists := e.probes[language]
	return exists
}