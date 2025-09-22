package healing

import (
	"context"
	"fmt"
	"time"

	"gproc/pkg/types"
)

type SelfHealingManager struct {
	config    *types.HealingConfig
	processes map[string]*HealingProcess
	policies  map[string]*HealingPolicy
}

type HealingProcess struct {
	ProcessID     string    `json:"process_id"`
	FailureCount  int       `json:"failure_count"`
	LastFailure   time.Time `json:"last_failure"`
	HealingAction string    `json:"healing_action"`
	Status        string    `json:"status"`
}

type HealingPolicy struct {
	Name            string        `json:"name"`
	Condition       string        `json:"condition"`
	Action          string        `json:"action"`
	MaxAttempts     int           `json:"max_attempts"`
	CooldownPeriod  time.Duration `json:"cooldown_period"`
	Enabled         bool          `json:"enabled"`
}

func NewSelfHealingManager(config *types.HealingConfig) *SelfHealingManager {
	return &SelfHealingManager{
		config:    config,
		processes: make(map[string]*HealingProcess),
		policies:  make(map[string]*HealingPolicy),
	}
}

func (s *SelfHealingManager) RegisterProcess(processID string) {
	s.processes[processID] = &HealingProcess{
		ProcessID:    processID,
		FailureCount: 0,
		Status:       "healthy",
	}
}

func (s *SelfHealingManager) HandleFailure(ctx context.Context, processID string, failure types.ProcessFailure) error {
	process, exists := s.processes[processID]
	if !exists {
		return fmt.Errorf("process %s not registered for healing", processID)
	}
	
	process.FailureCount++
	process.LastFailure = time.Now()
	
	// Find applicable healing policy
	policy := s.findPolicy(failure.Type)
	if policy == nil || !policy.Enabled {
		return nil
	}
	
	// Check cooldown period
	if time.Since(process.LastFailure) < policy.CooldownPeriod {
		return nil
	}
	
	// Check max attempts
	if process.FailureCount > policy.MaxAttempts {
		process.Status = "healing_failed"
		return fmt.Errorf("max healing attempts exceeded for process %s", processID)
	}
	
	// Execute healing action
	return s.executeHealingAction(ctx, processID, policy.Action)
}

func (s *SelfHealingManager) executeHealingAction(ctx context.Context, processID, action string) error {
	process := s.processes[processID]
	process.HealingAction = action
	process.Status = "healing"
	
	switch action {
	case "restart":
		return s.restartProcess(ctx, processID)
	case "replace":
		return s.replaceProcess(ctx, processID)
	case "scale_up":
		return s.scaleUpProcess(ctx, processID)
	case "migrate":
		return s.migrateProcess(ctx, processID)
	default:
		return fmt.Errorf("unknown healing action: %s", action)
	}
}

func (s *SelfHealingManager) restartProcess(ctx context.Context, processID string) error {
	fmt.Printf("Self-healing: Restarting process %s\n", processID)
	// Simulate process restart
	time.Sleep(2 * time.Second)
	s.processes[processID].Status = "healthy"
	return nil
}

func (s *SelfHealingManager) replaceProcess(ctx context.Context, processID string) error {
	fmt.Printf("Self-healing: Replacing process %s\n", processID)
	// Simulate process replacement
	time.Sleep(5 * time.Second)
	s.processes[processID].Status = "healthy"
	return nil
}

func (s *SelfHealingManager) scaleUpProcess(ctx context.Context, processID string) error {
	fmt.Printf("Self-healing: Scaling up process %s\n", processID)
	// Simulate scaling up
	time.Sleep(3 * time.Second)
	s.processes[processID].Status = "healthy"
	return nil
}

func (s *SelfHealingManager) migrateProcess(ctx context.Context, processID string) error {
	fmt.Printf("Self-healing: Migrating process %s to healthy node\n", processID)
	// Simulate process migration
	time.Sleep(10 * time.Second)
	s.processes[processID].Status = "healthy"
	return nil
}

func (s *SelfHealingManager) findPolicy(failureType string) *HealingPolicy {
	for _, policy := range s.policies {
		if policy.Condition == failureType {
			return policy
		}
	}
	return nil
}

func (s *SelfHealingManager) AddPolicy(policy *HealingPolicy) {
	s.policies[policy.Name] = policy
}