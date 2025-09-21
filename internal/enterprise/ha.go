package enterprise

import (
	"context"
	"fmt"
	"sync"
	"time"

	"gproc/pkg/types"
)

type HAManager struct {
	config     *types.HAConfig
	isActive   bool
	isPrimary  bool
	replicas   map[string]*ReplicaNode
	failover   *FailoverController
	mu         sync.RWMutex
}

type ReplicaNode struct {
	ID       string    `json:"id"`
	Address  string    `json:"address"`
	Status   string    `json:"status"` // active, standby, failed
	LastSeen time.Time `json:"last_seen"`
	Role     string    `json:"role"`   // primary, secondary
}

type FailoverController struct {
	config          *types.FailoverConfig
	healthChecks    map[string]*HealthChecker
	failoverInProgress bool
	mu              sync.Mutex
}

type HealthChecker struct {
	nodeID   string
	endpoint string
	interval time.Duration
	timeout  time.Duration
	failures int
	maxFailures int
}

func NewHAManager(config *types.HAConfig) *HAManager {
	ha := &HAManager{
		config:   config,
		replicas: make(map[string]*ReplicaNode),
	}

	if config.Failover != nil {
		ha.failover = &FailoverController{
			config:       config.Failover,
			healthChecks: make(map[string]*HealthChecker),
		}
	}

	return ha
}

func (ha *HAManager) Start(ctx context.Context) error {
	if !ha.config.Enabled {
		return nil
	}

	// Initialize as primary if no other nodes exist
	ha.mu.Lock()
	ha.isPrimary = len(ha.replicas) == 0
	ha.isActive = true
	ha.mu.Unlock()

	// Start health checking
	if ha.failover != nil {
		go ha.startHealthChecking(ctx)
	}

	// Start replica monitoring
	go ha.monitorReplicas(ctx)

	fmt.Printf("HA Manager started (mode: %s, primary: %v)\n", ha.config.Mode, ha.isPrimary)
	return nil
}

func (ha *HAManager) RegisterReplica(node *ReplicaNode) error {
	ha.mu.Lock()
	defer ha.mu.Unlock()

	ha.replicas[node.ID] = node

	// Setup health checker for this replica
	if ha.failover != nil {
		checker := &HealthChecker{
			nodeID:      node.ID,
			endpoint:    node.Address + "/health",
			interval:    30 * time.Second,
			timeout:     ha.failover.config.Timeout,
			maxFailures: ha.failover.config.Retries,
		}
		ha.failover.healthChecks[node.ID] = checker
	}

	fmt.Printf("Replica registered: %s at %s\n", node.ID, node.Address)
	return nil
}

func (ha *HAManager) IsPrimary() bool {
	ha.mu.RLock()
	defer ha.mu.RUnlock()
	return ha.isPrimary
}

func (ha *HAManager) IsActive() bool {
	ha.mu.RLock()
	defer ha.mu.RUnlock()
	return ha.isActive
}

func (ha *HAManager) GetActiveReplicas() []*ReplicaNode {
	ha.mu.RLock()
	defer ha.mu.RUnlock()

	var active []*ReplicaNode
	for _, replica := range ha.replicas {
		if replica.Status == "active" {
			active = append(active, replica)
		}
	}
	return active
}

func (ha *HAManager) TriggerFailover(failedNodeID string) error {
	if ha.failover == nil {
		return fmt.Errorf("failover not configured")
	}

	ha.failover.mu.Lock()
	defer ha.failover.mu.Unlock()

	if ha.failover.failoverInProgress {
		return fmt.Errorf("failover already in progress")
	}

	ha.failover.failoverInProgress = true
	defer func() { ha.failover.failoverInProgress = false }()

	fmt.Printf("Triggering failover for failed node: %s\n", failedNodeID)

	// Mark failed node as failed
	ha.mu.Lock()
	if replica, exists := ha.replicas[failedNodeID]; exists {
		replica.Status = "failed"
	}
	ha.mu.Unlock()

	// Select new primary based on mode
	switch ha.config.Mode {
	case "active-passive":
		return ha.performActivePassiveFailover(failedNodeID)
	case "active-active":
		return ha.performActiveActiveFailover(failedNodeID)
	default:
		return fmt.Errorf("unsupported HA mode: %s", ha.config.Mode)
	}
}

func (ha *HAManager) performActivePassiveFailover(failedNodeID string) error {
	// Find the best standby replica
	var bestReplica *ReplicaNode
	ha.mu.RLock()
	for _, replica := range ha.replicas {
		if replica.Status == "standby" && replica.ID != failedNodeID {
			if bestReplica == nil || replica.LastSeen.After(bestReplica.LastSeen) {
				bestReplica = replica
			}
		}
	}
	ha.mu.RUnlock()

	if bestReplica == nil {
		return fmt.Errorf("no healthy standby replica available")
	}

	// Promote standby to primary
	ha.mu.Lock()
	bestReplica.Status = "active"
	bestReplica.Role = "primary"
	
	// Update our state if we were the failed primary
	if failedNodeID == "self" {
		ha.isPrimary = false
		ha.isActive = false
	}
	ha.mu.Unlock()

	fmt.Printf("Promoted replica %s to primary\n", bestReplica.ID)
	return nil
}

func (ha *HAManager) performActiveActiveFailover(failedNodeID string) error {
	// In active-active mode, redistribute load among remaining nodes
	ha.mu.Lock()
	defer ha.mu.Unlock()

	activeCount := 0
	for _, replica := range ha.replicas {
		if replica.Status == "active" && replica.ID != failedNodeID {
			activeCount++
		}
	}

	if activeCount == 0 {
		return fmt.Errorf("no active replicas remaining")
	}

	fmt.Printf("Load redistributed among %d remaining active nodes\n", activeCount)
	return nil
}

func (ha *HAManager) startHealthChecking(ctx context.Context) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			ha.performHealthChecks()
		}
	}
}

func (ha *HAManager) performHealthChecks() {
	if ha.failover == nil {
		return
	}

	for nodeID, checker := range ha.failover.healthChecks {
		go func(id string, hc *HealthChecker) {
			if err := ha.checkNodeHealth(hc); err != nil {
				hc.failures++
				fmt.Printf("Health check failed for node %s: %v (failures: %d/%d)\n", 
					id, err, hc.failures, hc.maxFailures)

				if hc.failures >= hc.maxFailures {
					ha.TriggerFailover(id)
				}
			} else {
				hc.failures = 0 // Reset failure count on success
			}
		}(nodeID, checker)
	}
}

func (ha *HAManager) checkNodeHealth(checker *HealthChecker) error {
	// Simulate health check (in production, make HTTP request to health endpoint)
	// For now, randomly fail sometimes to demonstrate failover
	if time.Now().Unix()%30 < 2 { // Fail 2 seconds out of every 30
		return fmt.Errorf("simulated health check failure")
	}
	return nil
}

func (ha *HAManager) monitorReplicas(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			ha.updateReplicaStatus()
		}
	}
}

func (ha *HAManager) updateReplicaStatus() {
	ha.mu.Lock()
	defer ha.mu.Unlock()

	now := time.Now()
	for _, replica := range ha.replicas {
		// Mark replicas as failed if not seen for 2 minutes
		if now.Sub(replica.LastSeen) > 2*time.Minute && replica.Status != "failed" {
			replica.Status = "failed"
			fmt.Printf("Replica %s marked as failed (last seen: %v)\n", replica.ID, replica.LastSeen)
		}
	}
}

func (ha *HAManager) GetHAStatus() *HAStatus {
	ha.mu.RLock()
	defer ha.mu.RUnlock()

	status := &HAStatus{
		Enabled:   ha.config.Enabled,
		Mode:      ha.config.Mode,
		IsPrimary: ha.isPrimary,
		IsActive:  ha.isActive,
		Replicas:  make([]*ReplicaNode, 0, len(ha.replicas)),
	}

	for _, replica := range ha.replicas {
		status.Replicas = append(status.Replicas, replica)
	}

	if ha.failover != nil {
		status.FailoverEnabled = true
		status.FailoverInProgress = ha.failover.failoverInProgress
	}

	return status
}

type HAStatus struct {
	Enabled           bool           `json:"enabled"`
	Mode              string         `json:"mode"`
	IsPrimary         bool           `json:"is_primary"`
	IsActive          bool           `json:"is_active"`
	Replicas          []*ReplicaNode `json:"replicas"`
	FailoverEnabled   bool           `json:"failover_enabled"`
	FailoverInProgress bool          `json:"failover_in_progress"`
}