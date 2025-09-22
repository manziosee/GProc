package federation

import (
	"context"
	"fmt"
	"sync"
	"time"

	"gproc/pkg/types"
)

type FederationManager struct {
	config   *types.FederationConfig
	clusters map[string]*ClusterConnection
	router   *RequestRouter
	mutex    sync.RWMutex
}

type ClusterConnection struct {
	ClusterID   string    `json:"cluster_id"`
	Name        string    `json:"name"`
	Endpoint    string    `json:"endpoint"`
	Region      string    `json:"region"`
	Status      string    `json:"status"`
	LastPing    time.Time `json:"last_ping"`
	Latency     time.Duration `json:"latency"`
	ProcessCount int      `json:"process_count"`
	Capacity    *ClusterCapacity `json:"capacity"`
}

type ClusterCapacity struct {
	MaxProcesses int     `json:"max_processes"`
	MaxCPU       float64 `json:"max_cpu"`
	MaxMemory    int64   `json:"max_memory"`
	UsedCPU      float64 `json:"used_cpu"`
	UsedMemory   int64   `json:"used_memory"`
}

type RequestRouter struct {
	strategy string
	clusters map[string]*ClusterConnection
}

type FederatedProcess struct {
	ProcessID   string `json:"process_id"`
	ClusterID   string `json:"cluster_id"`
	Region      string `json:"region"`
	Status      string `json:"status"`
	Replicas    int    `json:"replicas"`
	Distribution map[string]int `json:"distribution"`
}

func NewFederationManager(config *types.FederationConfig) *FederationManager {
	return &FederationManager{
		config:   config,
		clusters: make(map[string]*ClusterConnection),
		router:   &RequestRouter{
			strategy: config.RoutingStrategy,
			clusters: make(map[string]*ClusterConnection),
		},
	}
}

func (f *FederationManager) RegisterCluster(ctx context.Context, cluster *ClusterConnection) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	
	// Test connection
	if err := f.testClusterConnection(ctx, cluster); err != nil {
		return fmt.Errorf("failed to connect to cluster %s: %v", cluster.ClusterID, err)
	}
	
	cluster.Status = "connected"
	cluster.LastPing = time.Now()
	
	f.clusters[cluster.ClusterID] = cluster
	f.router.clusters[cluster.ClusterID] = cluster
	
	// Start health monitoring
	go f.monitorCluster(ctx, cluster.ClusterID)
	
	return nil
}

func (f *FederationManager) DeployFederatedProcess(ctx context.Context, config *types.FederatedProcessConfig) (*FederatedProcess, error) {
	// Select target clusters based on strategy
	targetClusters, err := f.selectClusters(config)
	if err != nil {
		return nil, err
	}
	
	process := &FederatedProcess{
		ProcessID:    config.ProcessID,
		Status:       "deploying",
		Replicas:     config.Replicas,
		Distribution: make(map[string]int),
	}
	
	// Distribute replicas across clusters
	replicasPerCluster := config.Replicas / len(targetClusters)
	remainder := config.Replicas % len(targetClusters)
	
	for i, cluster := range targetClusters {
		replicas := replicasPerCluster
		if i < remainder {
			replicas++
		}
		
		if err := f.deployToCluster(ctx, cluster.ClusterID, config, replicas); err != nil {
			return nil, fmt.Errorf("failed to deploy to cluster %s: %v", cluster.ClusterID, err)
		}
		
		process.Distribution[cluster.ClusterID] = replicas
	}
	
	process.Status = "running"
	return process, nil
}

func (f *FederationManager) selectClusters(config *types.FederatedProcessConfig) ([]*ClusterConnection, error) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	
	var candidates []*ClusterConnection
	
	// Filter by region if specified
	for _, cluster := range f.clusters {
		if cluster.Status != "connected" {
			continue
		}
		
		if len(config.Regions) > 0 {
			regionMatch := false
			for _, region := range config.Regions {
				if cluster.Region == region {
					regionMatch = true
					break
				}
			}
			if !regionMatch {
				continue
			}
		}
		
		// Check capacity
		if f.hasCapacity(cluster, config) {
			candidates = append(candidates, cluster)
		}
	}
	
	if len(candidates) == 0 {
		return nil, fmt.Errorf("no suitable clusters found")
	}
	
	// Apply selection strategy
	switch f.router.strategy {
	case "round_robin":
		return f.selectRoundRobin(candidates, config.ClusterCount)
	case "least_loaded":
		return f.selectLeastLoaded(candidates, config.ClusterCount)
	case "latency":
		return f.selectByLatency(candidates, config.ClusterCount)
	default:
		return candidates[:min(len(candidates), config.ClusterCount)], nil
	}
}

func (f *FederationManager) hasCapacity(cluster *ClusterConnection, config *types.FederatedProcessConfig) bool {
	if cluster.Capacity == nil {
		return true // Assume capacity if not specified
	}
	
	requiredCPU := config.Resources.CPU * float64(config.Replicas)
	requiredMemory := config.Resources.Memory * int64(config.Replicas)
	
	availableCPU := cluster.Capacity.MaxCPU - cluster.Capacity.UsedCPU
	availableMemory := cluster.Capacity.MaxMemory - cluster.Capacity.UsedMemory
	
	return availableCPU >= requiredCPU && availableMemory >= requiredMemory
}

func (f *FederationManager) selectRoundRobin(candidates []*ClusterConnection, count int) []*ClusterConnection, error) {
	if count > len(candidates) {
		count = len(candidates)
	}
	return candidates[:count], nil
}

func (f *FederationManager) selectLeastLoaded(candidates []*ClusterConnection, count int) []*ClusterConnection, error) {
	// Sort by CPU usage
	for i := 0; i < len(candidates)-1; i++ {
		for j := i + 1; j < len(candidates); j++ {
			if candidates[i].Capacity.UsedCPU > candidates[j].Capacity.UsedCPU {
				candidates[i], candidates[j] = candidates[j], candidates[i]
			}
		}
	}
	
	if count > len(candidates) {
		count = len(candidates)
	}
	return candidates[:count], nil
}

func (f *FederationManager) selectByLatency(candidates []*ClusterConnection, count int) []*ClusterConnection, error) {
	// Sort by latency
	for i := 0; i < len(candidates)-1; i++ {
		for j := i + 1; j < len(candidates); j++ {
			if candidates[i].Latency > candidates[j].Latency {
				candidates[i], candidates[j] = candidates[j], candidates[i]
			}
		}
	}
	
	if count > len(candidates) {
		count = len(candidates)
	}
	return candidates[:count], nil
}

func (f *FederationManager) deployToCluster(ctx context.Context, clusterID string, config *types.FederatedProcessConfig, replicas int) error {
	fmt.Printf("Deploying %d replicas of process %s to cluster %s\n", replicas, config.ProcessID, clusterID)
	
	// Simulate deployment to remote cluster
	time.Sleep(100 * time.Millisecond)
	
	return nil
}

func (f *FederationManager) testClusterConnection(ctx context.Context, cluster *ClusterConnection) error {
	fmt.Printf("Testing connection to cluster %s at %s\n", cluster.ClusterID, cluster.Endpoint)
	
	// Simulate connection test
	start := time.Now()
	time.Sleep(10 * time.Millisecond)
	cluster.Latency = time.Since(start)
	
	return nil
}

func (f *FederationManager) monitorCluster(ctx context.Context, clusterID string) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			f.pingCluster(ctx, clusterID)
		}
	}
}

func (f *FederationManager) pingCluster(ctx context.Context, clusterID string) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	
	cluster, exists := f.clusters[clusterID]
	if !exists {
		return
	}
	
	start := time.Now()
	if err := f.testClusterConnection(ctx, cluster); err != nil {
		cluster.Status = "disconnected"
		fmt.Printf("Cluster %s is disconnected: %v\n", clusterID, err)
	} else {
		cluster.Status = "connected"
		cluster.LastPing = time.Now()
		cluster.Latency = time.Since(start)
	}
}

func (f *FederationManager) GetClusterStatus() map[string]*ClusterConnection {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	
	status := make(map[string]*ClusterConnection)
	for id, cluster := range f.clusters {
		status[id] = cluster
	}
	
	return status
}

func (f *FederationManager) MigrateProcess(ctx context.Context, processID, fromCluster, toCluster string) error {
	fmt.Printf("Migrating process %s from cluster %s to cluster %s\n", processID, fromCluster, toCluster)
	
	// Simulate process migration
	time.Sleep(500 * time.Millisecond)
	
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}