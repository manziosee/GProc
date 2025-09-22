package ml

import (
	"context"
	"fmt"
	"math"
	"time"

	"gproc/pkg/types"
)

type AnomalyDetector struct {
	config   *types.AnomalyConfig
	models   map[string]*AnomalyModel
	baseline map[string]*Baseline
}

type AnomalyModel struct {
	ProcessID   string    `json:"process_id"`
	Algorithm   string    `json:"algorithm"`
	Sensitivity float64   `json:"sensitivity"`
	TrainedAt   time.Time `json:"trained_at"`
}

type Baseline struct {
	CPUMean    float64 `json:"cpu_mean"`
	CPUStdDev  float64 `json:"cpu_stddev"`
	MemMean    float64 `json:"mem_mean"`
	MemStdDev  float64 `json:"mem_stddev"`
	SampleSize int     `json:"sample_size"`
}

type Anomaly struct {
	ProcessID   string    `json:"process_id"`
	Type        string    `json:"type"`
	Severity    string    `json:"severity"`
	Score       float64   `json:"score"`
	Description string    `json:"description"`
	DetectedAt  time.Time `json:"detected_at"`
	Metrics     map[string]float64 `json:"metrics"`
}

func NewAnomalyDetector(config *types.AnomalyConfig) *AnomalyDetector {
	return &AnomalyDetector{
		config:   config,
		models:   make(map[string]*AnomalyModel),
		baseline: make(map[string]*Baseline),
	}
}

func (a *AnomalyDetector) TrainModel(ctx context.Context, processID string, metrics []types.ProcessMetrics) error {
	if len(metrics) < 10 {
		return fmt.Errorf("insufficient data for training")
	}
	
	baseline := a.calculateBaseline(metrics)
	a.baseline[processID] = baseline
	
	model := &AnomalyModel{
		ProcessID:   processID,
		Algorithm:   "statistical",
		Sensitivity: a.config.Sensitivity,
		TrainedAt:   time.Now(),
	}
	a.models[processID] = model
	
	return nil
}

func (a *AnomalyDetector) DetectAnomalies(ctx context.Context, processID string, metrics *types.ProcessMetrics) ([]Anomaly, error) {
	baseline, exists := a.baseline[processID]
	if !exists {
		return nil, fmt.Errorf("no baseline for process %s", processID)
	}
	
	var anomalies []Anomaly
	
	// CPU anomaly detection
	cpuScore := a.calculateZScore(metrics.CPUUsage, baseline.CPUMean, baseline.CPUStdDev)
	if math.Abs(cpuScore) > a.config.Threshold {
		anomalies = append(anomalies, Anomaly{
			ProcessID:   processID,
			Type:        "cpu_spike",
			Severity:    a.getSeverity(cpuScore),
			Score:       cpuScore,
			Description: fmt.Sprintf("CPU usage %.2f%% deviates from baseline", metrics.CPUUsage),
			DetectedAt:  time.Now(),
			Metrics:     map[string]float64{"cpu": metrics.CPUUsage},
		})
	}
	
	// Memory anomaly detection
	memScore := a.calculateZScore(float64(metrics.MemoryUsage), baseline.MemMean, baseline.MemStdDev)
	if math.Abs(memScore) > a.config.Threshold {
		anomalies = append(anomalies, Anomaly{
			ProcessID:   processID,
			Type:        "memory_leak",
			Severity:    a.getSeverity(memScore),
			Score:       memScore,
			Description: fmt.Sprintf("Memory usage %d MB deviates from baseline", metrics.MemoryUsage/1024/1024),
			DetectedAt:  time.Now(),
			Metrics:     map[string]float64{"memory": float64(metrics.MemoryUsage)},
		})
	}
	
	return anomalies, nil
}

func (a *AnomalyDetector) calculateBaseline(metrics []types.ProcessMetrics) *Baseline {
	var cpuSum, memSum float64
	
	for _, m := range metrics {
		cpuSum += m.CPUUsage
		memSum += float64(m.MemoryUsage)
	}
	
	cpuMean := cpuSum / float64(len(metrics))
	memMean := memSum / float64(len(metrics))
	
	var cpuVariance, memVariance float64
	for _, m := range metrics {
		cpuVariance += math.Pow(m.CPUUsage-cpuMean, 2)
		memVariance += math.Pow(float64(m.MemoryUsage)-memMean, 2)
	}
	
	cpuStdDev := math.Sqrt(cpuVariance / float64(len(metrics)))
	memStdDev := math.Sqrt(memVariance / float64(len(metrics)))
	
	return &Baseline{
		CPUMean:    cpuMean,
		CPUStdDev:  cpuStdDev,
		MemMean:    memMean,
		MemStdDev:  memStdDev,
		SampleSize: len(metrics),
	}
}

func (a *AnomalyDetector) calculateZScore(value, mean, stddev float64) float64 {
	if stddev == 0 {
		return 0
	}
	return (value - mean) / stddev
}

func (a *AnomalyDetector) getSeverity(score float64) string {
	absScore := math.Abs(score)
	if absScore > 3.0 {
		return "critical"
	} else if absScore > 2.0 {
		return "high"
	} else if absScore > 1.5 {
		return "medium"
	}
	return "low"
}