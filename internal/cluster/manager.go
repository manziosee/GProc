package cluster

import (
	"context"
	"fmt"
	"sync"
	"time"

	"gproc/pkg/types"
)

type ClusterManager struct {
	config    *types.ClusterConfig
	nodeID    string
	nodes     map[string]*types.ClusterNode
	isLeader  bool
	discovery DiscoveryService
	consensus ConsensusService
	mu        sync.RWMutex
}

func NewClusterManager(config *types.ClusterConfig) *ClusterManager {
	cm := &ClusterManager{
		config: config,
		nodeID: config.NodeID,
		nodes:  make(map[string]*types.ClusterNode),
	}
	
	// Initialize discovery service
	if config.Discovery != nil {
		cm.discovery = NewDiscoveryService(config.Discovery)
	}
	
	// Initialize consensus service
	if config.Consensus != nil {
		cm.consensus = NewConsensusService(config.Consensus)
	}
	
	return cm
}

func (cm *ClusterManager) Start(ctx context.Context) error {
	if !cm.config.Enabled {
		return nil
	}
	
	// Register this node
	node := &types.ClusterNode{
		ID:       cm.nodeID,
		Address:  "localhost:8080", // Would be configurable
		Role:     "follower",
		Status:   "active",
		LastSeen: time.Now(),
		Metadata: make(map[string]string),
	}
	
	cm.mu.Lock()
	cm.nodes[cm.nodeID] = node
	cm.mu.Unlock()
	
	// Start discovery
	if cm.discovery != nil {
		go cm.discovery.Start(ctx, cm.onNodeDiscovered)
	}
	
	// Start consensus
	if cm.consensus != nil {
		go cm.consensus.Start(ctx, cm.onLeaderElected)
	}
	
	// Start health checking
	go cm.healthCheckLoop(ctx)
	
	return nil
}

func (cm *ClusterManager) onNodeDiscovered(node *types.ClusterNode) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	cm.nodes[node.ID] = node
	fmt.Printf("Node discovered: %s at %s\n", node.ID, node.Address)
}

func (cm *ClusterManager) onLeaderElected(leaderID string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	cm.isLeader = (leaderID == cm.nodeID)
	
	if node, exists := cm.nodes[leaderID]; exists {
		node.Role = "leader"
	}
	
	fmt.Printf("Leader elected: %s (isLeader: %v)\n", leaderID, cm.isLeader)
}

func (cm *ClusterManager) healthCheckLoop(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			cm.performHealthCheck()
		}
	}
}

func (cm *ClusterManager) performHealthCheck() {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	now := time.Now()
	for id, node := range cm.nodes {
		if id == cm.nodeID {
			node.LastSeen = now
			continue
		}
		
		// Check if node is stale (no heartbeat for 2 minutes)
		if now.Sub(node.LastSeen) > 2*time.Minute {
			node.Status = "failed"
		}
	}
}

func (cm *ClusterManager) GetNodes() []*types.ClusterNode {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	
	nodes := make([]*types.ClusterNode, 0, len(cm.nodes))
	for _, node := range cm.nodes {
		nodes = append(nodes, node)
	}
	return nodes
}

func (cm *ClusterManager) IsLeader() bool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.isLeader
}

func (cm *ClusterManager) GetLeader() *types.ClusterNode {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	
	for _, node := range cm.nodes {
		if node.Role == "leader" {
			return node
		}
	}
	return nil
}

// Discovery Service Interface
type DiscoveryService interface {
	Start(ctx context.Context, onNodeDiscovered func(*types.ClusterNode)) error
	Register(node *types.ClusterNode) error
	Deregister(nodeID string) error
}

// Consensus Service Interface
type ConsensusService interface {
	Start(ctx context.Context, onLeaderElected func(string)) error
	ProposeLeader(nodeID string) error
	GetCurrentLeader() string
}

// Simple in-memory discovery service
type InMemoryDiscovery struct {
	config *types.DiscoveryConfig
	nodes  map[string]*types.ClusterNode
	mu     sync.RWMutex
}

func NewDiscoveryService(config *types.DiscoveryConfig) DiscoveryService {
	return &InMemoryDiscovery{
		config: config,
		nodes:  make(map[string]*types.ClusterNode),
	}
}

func (d *InMemoryDiscovery) Start(ctx context.Context, onNodeDiscovered func(*types.ClusterNode)) error {
	// Simulate node discovery
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// Simulate discovering nodes
				d.mu.RLock()
				for _, node := range d.nodes {
					onNodeDiscovered(node)
				}
				d.mu.RUnlock()
			}
		}
	}()
	
	return nil
}

func (d *InMemoryDiscovery) Register(node *types.ClusterNode) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	
	d.nodes[node.ID] = node
	return nil
}

func (d *InMemoryDiscovery) Deregister(nodeID string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	
	delete(d.nodes, nodeID)
	return nil
}

// Simple consensus service
type SimpleConsensus struct {
	config   *types.ConsensusConfig
	leaderID string
	mu       sync.RWMutex
}

func NewConsensusService(config *types.ConsensusConfig) ConsensusService {
	return &SimpleConsensus{
		config: config,
	}
}

func (c *SimpleConsensus) Start(ctx context.Context, onLeaderElected func(string)) error {
	// Simple leader election - first node becomes leader
	go func() {
		time.Sleep(5 * time.Second) // Wait for nodes to register
		
		c.mu.Lock()
		if c.leaderID == "" {
			c.leaderID = "node-1" // Simple election
			onLeaderElected(c.leaderID)
		}
		c.mu.Unlock()
	}()
	
	return nil
}

func (c *SimpleConsensus) ProposeLeader(nodeID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.leaderID = nodeID
	return nil
}

func (c *SimpleConsensus) GetCurrentLeader() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.leaderID
}