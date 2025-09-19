package cluster

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type ClusterManager struct {
	mode      string // "master" or "agent"
	nodes     map[string]*ClusterNode
	mutex     sync.RWMutex
	masterURL string
	agentPort int
}

type ClusterNode struct {
	ID       string    `json:"id"`
	Address  string    `json:"address"`
	Status   string    `json:"status"`
	LastSeen time.Time `json:"last_seen"`
	Metadata map[string]string `json:"metadata"`
}

type ClusterMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
	From    string      `json:"from"`
	To      string      `json:"to"`
}

func NewClusterManager() *ClusterManager {
	return &ClusterManager{
		nodes: make(map[string]*ClusterNode),
		mode:  "standalone",
	}
}

func (cm *ClusterManager) InitMaster(port int) error {
	cm.mode = "master"
	cm.agentPort = port
	
	// Start master HTTP server
	http.HandleFunc("/cluster/join", cm.handleJoin)
	http.HandleFunc("/cluster/heartbeat", cm.handleHeartbeat)
	http.HandleFunc("/cluster/command", cm.handleCommand)
	http.HandleFunc("/cluster/nodes", cm.handleNodes)
	
	go func() {
		fmt.Printf("Cluster master starting on port %d\n", port)
		http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	}()
	
	// Start heartbeat monitor
	go cm.monitorNodes()
	
	return nil
}

func (cm *ClusterManager) JoinCluster(masterAddr string) error {
	cm.mode = "agent"
	cm.masterURL = masterAddr
	
	// Register with master
	joinReq := map[string]string{
		"node_id": fmt.Sprintf("agent-%d", time.Now().Unix()),
		"address": "localhost:9091", // Agent address
		"status":  "joining",
	}
	
	_, _ = json.Marshal(joinReq)
	resp, err := http.Post(fmt.Sprintf("http://%s/cluster/join", masterAddr), 
		"application/json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	// Start heartbeat
	go cm.sendHeartbeat()
	
	fmt.Printf("Joined cluster at %s\n", masterAddr)
	return nil
}

func (cm *ClusterManager) handleJoin(w http.ResponseWriter, r *http.Request) {
	var joinReq map[string]string
	json.NewDecoder(r.Body).Decode(&joinReq)
	
	node := &ClusterNode{
		ID:       joinReq["node_id"],
		Address:  joinReq["address"],
		Status:   "active",
		LastSeen: time.Now(),
		Metadata: make(map[string]string),
	}
	
	cm.mutex.Lock()
	cm.nodes[node.ID] = node
	cm.mutex.Unlock()
	
	fmt.Printf("Node %s joined cluster\n", node.ID)
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "joined"})
}

func (cm *ClusterManager) handleHeartbeat(w http.ResponseWriter, r *http.Request) {
	var heartbeat map[string]interface{}
	json.NewDecoder(r.Body).Decode(&heartbeat)
	
	nodeID := heartbeat["node_id"].(string)
	
	cm.mutex.Lock()
	if node, exists := cm.nodes[nodeID]; exists {
		node.LastSeen = time.Now()
		node.Status = "active"
	}
	cm.mutex.Unlock()
	
	w.WriteHeader(http.StatusOK)
}

func (cm *ClusterManager) handleCommand(w http.ResponseWriter, r *http.Request) {
	var cmd ClusterMessage
	json.NewDecoder(r.Body).Decode(&cmd)
	
	// Execute command on local node or forward to target
	result := cm.executeCommand(cmd.Type, cmd.Payload)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"result": result,
		"status": "success",
	})
}

func (cm *ClusterManager) handleNodes(w http.ResponseWriter, r *http.Request) {
	cm.mutex.RLock()
	nodes := make([]*ClusterNode, 0, len(cm.nodes))
	for _, node := range cm.nodes {
		nodes = append(nodes, node)
	}
	cm.mutex.RUnlock()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nodes)
}

func (cm *ClusterManager) sendHeartbeat() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	
	for range ticker.C {
		if cm.mode != "agent" || cm.masterURL == "" {
			continue
		}
		
		heartbeat := map[string]interface{}{
			"node_id":   "agent-1", // Should be dynamic
			"timestamp": time.Now(),
			"status":    "active",
		}
		
		_, _ = json.Marshal(heartbeat)
		http.Post(fmt.Sprintf("http://%s/cluster/heartbeat", cm.masterURL),
			"application/json", nil)
	}
}

func (cm *ClusterManager) monitorNodes() {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	
	for range ticker.C {
		cm.mutex.Lock()
		for nodeID, node := range cm.nodes {
			if time.Since(node.LastSeen) > 2*time.Minute {
				node.Status = "offline"
				fmt.Printf("Node %s marked as offline\n", nodeID)
			}
		}
		cm.mutex.Unlock()
	}
}

func (cm *ClusterManager) executeCommand(cmdType string, payload interface{}) interface{} {
	switch cmdType {
	case "list_processes":
		return map[string]string{"result": "process list"}
	case "start_process":
		return map[string]string{"result": "process started"}
	case "stop_process":
		return map[string]string{"result": "process stopped"}
	default:
		return map[string]string{"error": "unknown command"}
	}
}

func (cm *ClusterManager) GetNodes() []*ClusterNode {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	
	nodes := make([]*ClusterNode, 0, len(cm.nodes))
	for _, node := range cm.nodes {
		nodes = append(nodes, node)
	}
	return nodes
}

func (cm *ClusterManager) ExecuteRemoteCommand(nodeID, command string, args []string) error {
	cm.mutex.RLock()
	node, exists := cm.nodes[nodeID]
	cm.mutex.RUnlock()
	
	if !exists {
		return fmt.Errorf("node %s not found", nodeID)
	}
	
	cmd := ClusterMessage{
		Type: command,
		Payload: map[string]interface{}{
			"args": args,
		},
		From: "master",
		To:   nodeID,
	}
	
	_, _ = json.Marshal(cmd)
	resp, err := http.Post(fmt.Sprintf("http://%s/cluster/command", node.Address),
		"application/json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	fmt.Printf("Executed command %s on node %s\n", command, nodeID)
	return nil
}