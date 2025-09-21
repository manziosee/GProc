package observability

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gproc/pkg/types"
)

type MetricsManager struct {
	config    *types.MetricsConfig
	registry  *prometheus.Registry
	server    *http.Server
	
	// Process metrics
	processCount    *prometheus.GaugeVec
	processRestarts *prometheus.CounterVec
	processUptime   *prometheus.GaugeVec
	processCPU      *prometheus.GaugeVec
	processMemory   *prometheus.GaugeVec
	
	// System metrics
	systemCPU    prometheus.Gauge
	systemMemory prometheus.Gauge
	systemLoad   prometheus.Gauge
	
	mu sync.RWMutex
}

func NewMetricsManager(config *types.MetricsConfig) *MetricsManager {
	registry := prometheus.NewRegistry()
	
	mm := &MetricsManager{
		config:   config,
		registry: registry,
	}
	
	mm.initMetrics()
	return mm
}

func (mm *MetricsManager) initMetrics() {
	// Process metrics
	mm.processCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gproc_processes_total",
			Help: "Total number of processes by status",
		},
		[]string{"status"},
	)
	
	mm.processRestarts = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gproc_process_restarts_total",
			Help: "Total number of process restarts",
		},
		[]string{"process_name"},
	)
	
	mm.processUptime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gproc_process_uptime_seconds",
			Help: "Process uptime in seconds",
		},
		[]string{"process_name"},
	)
	
	mm.processCPU = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gproc_process_cpu_percent",
			Help: "Process CPU usage percentage",
		},
		[]string{"process_name"},
	)
	
	mm.processMemory = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gproc_process_memory_bytes",
			Help: "Process memory usage in bytes",
		},
		[]string{"process_name"},
	)
	
	// System metrics
	mm.systemCPU = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "gproc_system_cpu_percent",
			Help: "System CPU usage percentage",
		},
	)
	
	mm.systemMemory = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "gproc_system_memory_percent",
			Help: "System memory usage percentage",
		},
	)
	
	mm.systemLoad = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "gproc_system_load_average",
			Help: "System load average",
		},
	)
	
	// Register metrics
	mm.registry.MustRegister(
		mm.processCount,
		mm.processRestarts,
		mm.processUptime,
		mm.processCPU,
		mm.processMemory,
		mm.systemCPU,
		mm.systemMemory,
		mm.systemLoad,
	)
}

func (mm *MetricsManager) Start(ctx context.Context) error {
	if !mm.config.Enabled {
		return nil
	}
	
	// Start Prometheus HTTP server
	if mm.config.Prometheus != nil {
		mux := http.NewServeMux()
		mux.Handle(mm.config.Prometheus.Path, promhttp.HandlerFor(mm.registry, promhttp.HandlerOpts{}))
		
		mm.server = &http.Server{
			Addr:    fmt.Sprintf(":%d", mm.config.Prometheus.Port),
			Handler: mux,
		}
		
		go func() {
			if err := mm.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				fmt.Printf("Metrics server error: %v\n", err)
			}
		}()
		
		fmt.Printf("Metrics server started on :%d%s\n", 
			mm.config.Prometheus.Port, mm.config.Prometheus.Path)
	}
	
	// Start metrics collection
	go mm.collectMetrics(ctx)
	
	return nil
}

func (mm *MetricsManager) Stop() error {
	if mm.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return mm.server.Shutdown(ctx)
	}
	return nil
}

func (mm *MetricsManager) collectMetrics(ctx context.Context) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			mm.updateSystemMetrics()
		}
	}
}

func (mm *MetricsManager) updateSystemMetrics() {
	// Simulate system metrics (in production, use actual system calls)
	mm.systemCPU.Set(float64(35 + (time.Now().Unix() % 30)))
	mm.systemMemory.Set(float64(60 + (time.Now().Unix() % 20)))
	mm.systemLoad.Set(float64(1.2 + float64(time.Now().Unix()%10)/10))
}

func (mm *MetricsManager) RecordProcessStart(processName string) {
	mm.processCount.WithLabelValues("running").Inc()
	mm.processUptime.WithLabelValues(processName).Set(0)
}

func (mm *MetricsManager) RecordProcessStop(processName string) {
	mm.processCount.WithLabelValues("running").Dec()
	mm.processCount.WithLabelValues("stopped").Inc()
	mm.processUptime.WithLabelValues(processName).Set(0)
}

func (mm *MetricsManager) RecordProcessRestart(processName string) {
	mm.processRestarts.WithLabelValues(processName).Inc()
}

func (mm *MetricsManager) UpdateProcessMetrics(processName string, cpu float64, memory float64, uptime time.Duration) {
	mm.processCPU.WithLabelValues(processName).Set(cpu)
	mm.processMemory.WithLabelValues(processName).Set(memory)
	mm.processUptime.WithLabelValues(processName).Set(uptime.Seconds())
}

// Alerting Manager
type AlertManager struct {
	config    *types.AlertingConfig
	providers map[string]AlertProvider
	rules     []*types.AlertRule
	mu        sync.RWMutex
}

func NewAlertManager(config *types.AlertingConfig) *AlertManager {
	am := &AlertManager{
		config:    config,
		providers: make(map[string]AlertProvider),
		rules:     config.Rules,
	}
	
	// Initialize alert providers
	for _, providerConfig := range config.Providers {
		provider := NewAlertProvider(providerConfig)
		am.providers[providerConfig.Name] = provider
	}
	
	return am
}

func (am *AlertManager) Start(ctx context.Context) error {
	if !am.config.Enabled {
		return nil
	}
	
	// Start rule evaluation
	go am.evaluateRules(ctx)
	
	return nil
}

func (am *AlertManager) evaluateRules(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			am.checkRules()
		}
	}
}

func (am *AlertManager) checkRules() {
	am.mu.RLock()
	defer am.mu.RUnlock()
	
	for _, rule := range am.rules {
		// Simulate rule evaluation
		if am.shouldTriggerAlert(rule) {
			am.sendAlert(rule)
		}
	}
}

func (am *AlertManager) shouldTriggerAlert(rule *types.AlertRule) bool {
	// Simple simulation - trigger alert randomly
	return time.Now().Unix()%60 < 5 // 5 seconds out of every minute
}

func (am *AlertManager) sendAlert(rule *types.AlertRule) {
	alert := Alert{
		Name:      rule.Name,
		Severity:  rule.Severity,
		Message:   fmt.Sprintf("Alert triggered: %s", rule.Condition),
		Timestamp: time.Now(),
	}
	
	for _, providerName := range rule.Providers {
		if provider, exists := am.providers[providerName]; exists {
			go provider.SendAlert(alert)
		}
	}
}

// Alert Provider Interface
type AlertProvider interface {
	SendAlert(alert Alert) error
}

type Alert struct {
	Name      string
	Severity  string
	Message   string
	Timestamp time.Time
}

// Slack Alert Provider
type SlackProvider struct {
	webhookURL string
}

func NewAlertProvider(config types.AlertProvider) AlertProvider {
	switch config.Type {
	case "slack":
		return &SlackProvider{
			webhookURL: config.Config["webhook_url"],
		}
	case "email":
		return &EmailProvider{
			smtpServer: config.Config["smtp_server"],
			from:       config.Config["from"],
			to:         config.Config["to"],
		}
	default:
		return &LogProvider{}
	}
}

func (s *SlackProvider) SendAlert(alert Alert) error {
	// In production, send actual Slack webhook
	fmt.Printf("SLACK ALERT [%s]: %s - %s\n", alert.Severity, alert.Name, alert.Message)
	return nil
}

// Email Alert Provider
type EmailProvider struct {
	smtpServer string
	from       string
	to         string
}

func (e *EmailProvider) SendAlert(alert Alert) error {
	// In production, send actual email
	fmt.Printf("EMAIL ALERT [%s]: %s - %s\n", alert.Severity, alert.Name, alert.Message)
	return nil
}

// Log Alert Provider
type LogProvider struct{}

func (l *LogProvider) SendAlert(alert Alert) error {
	fmt.Printf("LOG ALERT [%s]: %s - %s\n", alert.Severity, alert.Name, alert.Message)
	return nil
}