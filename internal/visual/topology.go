package visual

import (
	"encoding/json"
	"fmt"
	"time"

	"gproc/pkg/types"
)

type TopologyManager struct {
	config *types.TopologyConfig
	graph  *ProcessGraph
}

type ProcessGraph struct {
	Nodes []ProcessNode `json:"nodes"`
	Edges []ProcessEdge `json:"edges"`
	Layout string       `json:"layout"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProcessNode struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Type     string            `json:"type"`
	Status   string            `json:"status"`
	Position Position          `json:"position"`
	Metrics  *NodeMetrics      `json:"metrics"`
	Config   map[string]string `json:"config"`
	Group    string            `json:"group"`
}

type ProcessEdge struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
	Type   string `json:"type"`
	Label  string `json:"label"`
	Weight int    `json:"weight"`
}

type NodeMetrics struct {
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage int64   `json:"memory_usage"`
	Connections int     `json:"connections"`
	Requests    int     `json:"requests"`
	Errors      int     `json:"errors"`
	Uptime      int64   `json:"uptime"`
}

func NewTopologyManager(config *types.TopologyConfig) *TopologyManager {
	return &TopologyManager{
		config: config,
		graph: &ProcessGraph{
			Nodes:     []ProcessNode{},
			Edges:     []ProcessEdge{},
			Layout:    "force-directed",
			UpdatedAt: time.Now(),
		},
	}
}

func (t *TopologyManager) AddProcess(processID, name, processType, group string) *ProcessNode {
	node := ProcessNode{
		ID:     processID,
		Name:   name,
		Type:   processType,
		Status: "running",
		Position: Position{
			X: t.calculateNodeX(len(t.graph.Nodes)),
			Y: t.calculateNodeY(len(t.graph.Nodes)),
		},
		Metrics: &NodeMetrics{},
		Config:  make(map[string]string),
		Group:   group,
	}
	
	t.graph.Nodes = append(t.graph.Nodes, node)
	t.graph.UpdatedAt = time.Now()
	
	return &node
}

func (t *TopologyManager) AddConnection(sourceID, targetID, connectionType, label string) *ProcessEdge {
	edge := ProcessEdge{
		ID:     fmt.Sprintf("%s-%s", sourceID, targetID),
		Source: sourceID,
		Target: targetID,
		Type:   connectionType,
		Label:  label,
		Weight: 1,
	}
	
	t.graph.Edges = append(t.graph.Edges, edge)
	t.graph.UpdatedAt = time.Now()
	
	return &edge
}

func (t *TopologyManager) UpdateNodeMetrics(nodeID string, metrics *NodeMetrics) error {
	for i := range t.graph.Nodes {
		if t.graph.Nodes[i].ID == nodeID {
			t.graph.Nodes[i].Metrics = metrics
			t.graph.UpdatedAt = time.Now()
			return nil
		}
	}
	return fmt.Errorf("node not found: %s", nodeID)
}

func (t *TopologyManager) UpdateNodeStatus(nodeID, status string) error {
	for i := range t.graph.Nodes {
		if t.graph.Nodes[i].ID == nodeID {
			t.graph.Nodes[i].Status = status
			t.graph.UpdatedAt = time.Now()
			return nil
		}
	}
	return fmt.Errorf("node not found: %s", nodeID)
}

func (t *TopologyManager) RemoveProcess(processID string) error {
	// Remove node
	for i, node := range t.graph.Nodes {
		if node.ID == processID {
			t.graph.Nodes = append(t.graph.Nodes[:i], t.graph.Nodes[i+1:]...)
			break
		}
	}
	
	// Remove associated edges
	var newEdges []ProcessEdge
	for _, edge := range t.graph.Edges {
		if edge.Source != processID && edge.Target != processID {
			newEdges = append(newEdges, edge)
		}
	}
	t.graph.Edges = newEdges
	
	t.graph.UpdatedAt = time.Now()
	return nil
}

func (t *TopologyManager) GetTopologyJSON() (string, error) {
	data, err := json.MarshalIndent(t.graph, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (t *TopologyManager) GetD3ForceLayout() (string, error) {
	// Generate D3.js force-directed layout configuration
	d3Config := map[string]interface{}{
		"nodes": t.graph.Nodes,
		"links": t.convertEdgesToLinks(),
		"simulation": map[string]interface{}{
			"forces": map[string]interface{}{
				"link": map[string]interface{}{
					"distance": 100,
					"strength": 0.1,
				},
				"charge": map[string]interface{}{
					"strength": -300,
				},
				"center": map[string]interface{}{
					"x": 400,
					"y": 300,
				},
			},
		},
	}
	
	data, err := json.MarshalIndent(d3Config, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (t *TopologyManager) convertEdgesToLinks() []map[string]interface{} {
	var links []map[string]interface{}
	
	for _, edge := range t.graph.Edges {
		link := map[string]interface{}{
			"source": edge.Source,
			"target": edge.Target,
			"type":   edge.Type,
			"label":  edge.Label,
			"weight": edge.Weight,
		}
		links = append(links, link)
	}
	
	return links
}

func (t *TopologyManager) GetCytoscapeLayout() (string, error) {
	// Generate Cytoscape.js layout configuration
	elements := []map[string]interface{}{}
	
	// Add nodes
	for _, node := range t.graph.Nodes {
		element := map[string]interface{}{
			"data": map[string]interface{}{
				"id":     node.ID,
				"label":  node.Name,
				"type":   node.Type,
				"status": node.Status,
				"group":  node.Group,
			},
			"position": map[string]interface{}{
				"x": node.Position.X,
				"y": node.Position.Y,
			},
			"classes": t.getNodeClasses(node),
		}
		elements = append(elements, element)
	}
	
	// Add edges
	for _, edge := range t.graph.Edges {
		element := map[string]interface{}{
			"data": map[string]interface{}{
				"id":     edge.ID,
				"source": edge.Source,
				"target": edge.Target,
				"label":  edge.Label,
				"type":   edge.Type,
				"weight": edge.Weight,
			},
			"classes": t.getEdgeClasses(edge),
		}
		elements = append(elements, element)
	}
	
	cytoscapeConfig := map[string]interface{}{
		"elements": elements,
		"style": t.getCytoscapeStyles(),
		"layout": map[string]interface{}{
			"name": "cose",
			"idealEdgeLength": 100,
			"nodeOverlap": 20,
			"refresh": 20,
			"fit": true,
			"padding": 30,
		},
	}
	
	data, err := json.MarshalIndent(cytoscapeConfig, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (t *TopologyManager) getNodeClasses(node ProcessNode) string {
	classes := node.Type
	
	switch node.Status {
	case "running":
		classes += " running"
	case "stopped":
		classes += " stopped"
	case "failed":
		classes += " failed"
	}
	
	if node.Metrics != nil && node.Metrics.CPUUsage > 80 {
		classes += " high-cpu"
	}
	
	return classes
}

func (t *TopologyManager) getEdgeClasses(edge ProcessEdge) string {
	return edge.Type
}

func (t *TopologyManager) getCytoscapeStyles() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"selector": "node",
			"style": map[string]interface{}{
				"background-color": "#666",
				"label":           "data(label)",
				"width":           "60px",
				"height":          "60px",
				"text-valign":     "center",
				"text-halign":     "center",
				"color":           "#fff",
				"font-size":       "12px",
			},
		},
		{
			"selector": "node.running",
			"style": map[string]interface{}{
				"background-color": "#4CAF50",
			},
		},
		{
			"selector": "node.stopped",
			"style": map[string]interface{}{
				"background-color": "#9E9E9E",
			},
		},
		{
			"selector": "node.failed",
			"style": map[string]interface{}{
				"background-color": "#F44336",
			},
		},
		{
			"selector": "node.high-cpu",
			"style": map[string]interface{}{
				"border-width": "3px",
				"border-color": "#FF9800",
			},
		},
		{
			"selector": "edge",
			"style": map[string]interface{}{
				"width":              "2px",
				"line-color":         "#ccc",
				"target-arrow-color": "#ccc",
				"target-arrow-shape": "triangle",
				"curve-style":        "bezier",
				"label":              "data(label)",
				"font-size":          "10px",
				"text-rotation":      "autorotate",
			},
		},
		{
			"selector": "edge.http",
			"style": map[string]interface{}{
				"line-color":         "#2196F3",
				"target-arrow-color": "#2196F3",
			},
		},
		{
			"selector": "edge.database",
			"style": map[string]interface{}{
				"line-color":         "#9C27B0",
				"target-arrow-color": "#9C27B0",
			},
		},
	}
}

func (t *TopologyManager) calculateNodeX(index int) int {
	// Simple grid layout calculation
	cols := 5
	col := index % cols
	return 100 + (col * 150)
}

func (t *TopologyManager) calculateNodeY(index int) int {
	// Simple grid layout calculation
	cols := 5
	row := index / cols
	return 100 + (row * 120)
}

func (t *TopologyManager) GetHealthOverview() map[string]interface{} {
	totalNodes := len(t.graph.Nodes)
	runningNodes := 0
	failedNodes := 0
	highCPUNodes := 0
	
	for _, node := range t.graph.Nodes {
		switch node.Status {
		case "running":
			runningNodes++
		case "failed":
			failedNodes++
		}
		
		if node.Metrics != nil && node.Metrics.CPUUsage > 80 {
			highCPUNodes++
		}
	}
	
	return map[string]interface{}{
		"total_processes":   totalNodes,
		"running_processes": runningNodes,
		"failed_processes":  failedNodes,
		"high_cpu_processes": highCPUNodes,
		"total_connections": len(t.graph.Edges),
		"last_updated":     t.graph.UpdatedAt,
	}
}

func (t *TopologyManager) GetProcessDependencies(processID string) ([]string, []string, error) {
	var dependencies []string // processes this depends on
	var dependents []string   // processes that depend on this
	
	for _, edge := range t.graph.Edges {
		if edge.Target == processID {
			dependencies = append(dependencies, edge.Source)
		}
		if edge.Source == processID {
			dependents = append(dependents, edge.Target)
		}
	}
	
	return dependencies, dependents, nil
}