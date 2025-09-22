package cost

import (
	"context"
	"fmt"
	"time"

	"gproc/pkg/types"
)

type CostTracker struct {
	config    *types.CostConfig
	rates     *ResourceRates
	usage     map[string]*ResourceUsage
	reports   map[string]*CostReport
}

type ResourceRates struct {
	CPUPerHour    float64 `json:"cpu_per_hour"`
	MemoryPerGB   float64 `json:"memory_per_gb_hour"`
	StoragePerGB  float64 `json:"storage_per_gb_hour"`
	NetworkPerGB  float64 `json:"network_per_gb"`
	ProcessPerHour float64 `json:"process_per_hour"`
}

type ResourceUsage struct {
	ProcessID     string        `json:"process_id"`
	TenantID      string        `json:"tenant_id"`
	StartTime     time.Time     `json:"start_time"`
	EndTime       time.Time     `json:"end_time"`
	CPUHours      float64       `json:"cpu_hours"`
	MemoryGBHours float64       `json:"memory_gb_hours"`
	StorageGBHours float64      `json:"storage_gb_hours"`
	NetworkGB     float64       `json:"network_gb"`
	ProcessHours  float64       `json:"process_hours"`
}

type CostReport struct {
	TenantID      string                 `json:"tenant_id"`
	Period        string                 `json:"period"`
	StartTime     time.Time              `json:"start_time"`
	EndTime       time.Time              `json:"end_time"`
	TotalCost     float64                `json:"total_cost"`
	Breakdown     *CostBreakdown         `json:"breakdown"`
	ProcessCosts  map[string]float64     `json:"process_costs"`
	Trends        *CostTrends            `json:"trends"`
}

type CostBreakdown struct {
	CPUCost     float64 `json:"cpu_cost"`
	MemoryCost  float64 `json:"memory_cost"`
	StorageCost float64 `json:"storage_cost"`
	NetworkCost float64 `json:"network_cost"`
	ProcessCost float64 `json:"process_cost"`
}

type CostTrends struct {
	DailyAverage   float64 `json:"daily_average"`
	WeeklyGrowth   float64 `json:"weekly_growth"`
	MonthlyGrowth  float64 `json:"monthly_growth"`
	Forecast       float64 `json:"forecast_next_month"`
}

type CostAlert struct {
	TenantID    string    `json:"tenant_id"`
	Type        string    `json:"type"`
	Threshold   float64   `json:"threshold"`
	CurrentCost float64   `json:"current_cost"`
	Message     string    `json:"message"`
	Timestamp   time.Time `json:"timestamp"`
}

func NewCostTracker(config *types.CostConfig) *CostTracker {
	return &CostTracker{
		config: config,
		rates: &ResourceRates{
			CPUPerHour:     0.05,   // $0.05 per CPU hour
			MemoryPerGB:    0.01,   // $0.01 per GB hour
			StoragePerGB:   0.001,  // $0.001 per GB hour
			NetworkPerGB:   0.02,   // $0.02 per GB
			ProcessPerHour: 0.001,  // $0.001 per process hour
		},
		usage:   make(map[string]*ResourceUsage),
		reports: make(map[string]*CostReport),
	}
}

func (c *CostTracker) StartTracking(processID, tenantID string) {
	usage := &ResourceUsage{
		ProcessID: processID,
		TenantID:  tenantID,
		StartTime: time.Now(),
	}
	
	c.usage[processID] = usage
}

func (c *CostTracker) StopTracking(processID string) {
	if usage, exists := c.usage[processID]; exists {
		usage.EndTime = time.Now()
		c.calculateUsage(usage)
	}
}

func (c *CostTracker) UpdateResourceUsage(processID string, metrics *types.ProcessMetrics) {
	usage, exists := c.usage[processID]
	if !exists {
		return
	}
	
	duration := time.Since(usage.StartTime).Hours()
	
	// Update cumulative usage
	usage.CPUHours += metrics.CPUUsage / 100.0 * (1.0 / 3600.0) // Convert to hours
	usage.MemoryGBHours += float64(metrics.MemoryUsage) / (1024 * 1024 * 1024) * (1.0 / 3600.0)
	usage.ProcessHours = duration
}

func (c *CostTracker) calculateUsage(usage *ResourceUsage) {
	if usage.EndTime.IsZero() {
		usage.EndTime = time.Now()
	}
	
	duration := usage.EndTime.Sub(usage.StartTime).Hours()
	usage.ProcessHours = duration
}

func (c *CostTracker) CalculateCost(usage *ResourceUsage) float64 {
	cost := 0.0
	
	cost += usage.CPUHours * c.rates.CPUPerHour
	cost += usage.MemoryGBHours * c.rates.MemoryPerGB
	cost += usage.StorageGBHours * c.rates.StoragePerGB
	cost += usage.NetworkGB * c.rates.NetworkPerGB
	cost += usage.ProcessHours * c.rates.ProcessPerHour
	
	return cost
}

func (c *CostTracker) GenerateReport(ctx context.Context, tenantID string, period string) (*CostReport, error) {
	var startTime, endTime time.Time
	now := time.Now()
	
	switch period {
	case "daily":
		startTime = now.AddDate(0, 0, -1)
		endTime = now
	case "weekly":
		startTime = now.AddDate(0, 0, -7)
		endTime = now
	case "monthly":
		startTime = now.AddDate(0, -1, 0)
		endTime = now
	default:
		return nil, fmt.Errorf("unsupported period: %s", period)
	}
	
	report := &CostReport{
		TenantID:     tenantID,
		Period:       period,
		StartTime:    startTime,
		EndTime:      endTime,
		ProcessCosts: make(map[string]float64),
		Breakdown:    &CostBreakdown{},
		Trends:       &CostTrends{},
	}
	
	// Calculate costs for all processes in the tenant
	totalCost := 0.0
	for _, usage := range c.usage {
		if usage.TenantID != tenantID {
			continue
		}
		
		if usage.StartTime.Before(endTime) && (usage.EndTime.IsZero() || usage.EndTime.After(startTime)) {
			processCost := c.CalculateCost(usage)
			report.ProcessCosts[usage.ProcessID] = processCost
			totalCost += processCost
			
			// Add to breakdown
			report.Breakdown.CPUCost += usage.CPUHours * c.rates.CPUPerHour
			report.Breakdown.MemoryCost += usage.MemoryGBHours * c.rates.MemoryPerGB
			report.Breakdown.StorageCost += usage.StorageGBHours * c.rates.StoragePerGB
			report.Breakdown.NetworkCost += usage.NetworkGB * c.rates.NetworkPerGB
			report.Breakdown.ProcessCost += usage.ProcessHours * c.rates.ProcessPerHour
		}
	}
	
	report.TotalCost = totalCost
	
	// Calculate trends
	c.calculateTrends(report)
	
	c.reports[fmt.Sprintf("%s-%s", tenantID, period)] = report
	return report, nil
}

func (c *CostTracker) calculateTrends(report *CostReport) {
	// Simplified trend calculation
	report.Trends.DailyAverage = report.TotalCost / float64(report.EndTime.Sub(report.StartTime).Hours()/24)
	report.Trends.WeeklyGrowth = 0.05  // 5% growth
	report.Trends.MonthlyGrowth = 0.15 // 15% growth
	report.Trends.Forecast = report.TotalCost * 1.15 // 15% increase forecast
}

func (c *CostTracker) CheckBudgetAlerts(ctx context.Context, tenantID string, budget float64) ([]*CostAlert, error) {
	var alerts []*CostAlert
	
	// Get current month cost
	report, err := c.GenerateReport(ctx, tenantID, "monthly")
	if err != nil {
		return nil, err
	}
	
	// Check if approaching budget (80% threshold)
	if report.TotalCost >= budget*0.8 {
		alert := &CostAlert{
			TenantID:    tenantID,
			Type:        "budget_warning",
			Threshold:   budget * 0.8,
			CurrentCost: report.TotalCost,
			Message:     fmt.Sprintf("Cost $%.2f is approaching budget limit $%.2f", report.TotalCost, budget),
			Timestamp:   time.Now(),
		}
		alerts = append(alerts, alert)
	}
	
	// Check if exceeded budget
	if report.TotalCost >= budget {
		alert := &CostAlert{
			TenantID:    tenantID,
			Type:        "budget_exceeded",
			Threshold:   budget,
			CurrentCost: report.TotalCost,
			Message:     fmt.Sprintf("Cost $%.2f has exceeded budget limit $%.2f", report.TotalCost, budget),
			Timestamp:   time.Now(),
		}
		alerts = append(alerts, alert)
	}
	
	return alerts, nil
}

func (c *CostTracker) GetCostEstimate(resources *types.ResourceRequest, duration time.Duration) float64 {
	hours := duration.Hours()
	
	estimate := 0.0
	estimate += resources.CPURequest * hours * c.rates.CPUPerHour
	estimate += float64(resources.MemoryRequest)/(1024*1024*1024) * hours * c.rates.MemoryPerGB
	estimate += float64(resources.ProcessCount) * hours * c.rates.ProcessPerHour
	
	return estimate
}

func (c *CostTracker) OptimizeCosts(tenantID string) (*CostOptimization, error) {
	report, err := c.GenerateReport(context.Background(), tenantID, "monthly")
	if err != nil {
		return nil, err
	}
	
	optimization := &CostOptimization{
		TenantID:      tenantID,
		CurrentCost:   report.TotalCost,
		Recommendations: []CostRecommendation{},
	}
	
	// Find expensive processes
	for processID, cost := range report.ProcessCosts {
		if cost > report.TotalCost*0.2 { // Process using >20% of budget
			recommendation := CostRecommendation{
				Type:        "right_size",
				ProcessID:   processID,
				Description: fmt.Sprintf("Process %s costs $%.2f (%.1f%% of total). Consider right-sizing resources.", processID, cost, cost/report.TotalCost*100),
				PotentialSavings: cost * 0.3, // Assume 30% savings
			}
			optimization.Recommendations = append(optimization.Recommendations, recommendation)
			optimization.PotentialSavings += recommendation.PotentialSavings
		}
	}
	
	return optimization, nil
}

type CostOptimization struct {
	TenantID          string                `json:"tenant_id"`
	CurrentCost       float64               `json:"current_cost"`
	PotentialSavings  float64               `json:"potential_savings"`
	Recommendations   []CostRecommendation  `json:"recommendations"`
}

type CostRecommendation struct {
	Type             string  `json:"type"`
	ProcessID        string  `json:"process_id"`
	Description      string  `json:"description"`
	PotentialSavings float64 `json:"potential_savings"`
}