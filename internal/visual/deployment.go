package visual

import (
	"encoding/json"
	"fmt"
	"time"

	"gproc/pkg/types"
)

type VisualDeploymentManager struct {
	config     *types.VisualConfig
	canvas     *DeploymentCanvas
	components map[string]*Component
	connections map[string]*Connection
}

type DeploymentCanvas struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	Width      int         `json:"width"`
	Height     int         `json:"height"`
	Components []Component `json:"components"`
	Connections []Connection `json:"connections"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

type Component struct {
	ID       string            `json:"id"`
	Type     string            `json:"type"`
	Name     string            `json:"name"`
	Position Position          `json:"position"`
	Size     Size              `json:"size"`
	Config   map[string]string `json:"config"`
	Status   string            `json:"status"`
	Ports    []Port            `json:"ports"`
}

type Connection struct {
	ID       string `json:"id"`
	FromPort string `json:"from_port"`
	ToPort   string `json:"to_port"`
	Type     string `json:"type"`
	Label    string `json:"label"`
}

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Size struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type Port struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Direction string `json:"direction"`
	Position  Position `json:"position"`
}

func NewVisualDeploymentManager(config *types.VisualConfig) *VisualDeploymentManager {
	return &VisualDeploymentManager{
		config:      config,
		components:  make(map[string]*Component),
		connections: make(map[string]*Connection),
		canvas: &DeploymentCanvas{
			ID:          "main-canvas",
			Name:        "Main Deployment Canvas",
			Width:       1920,
			Height:      1080,
			Components:  []Component{},
			Connections: []Connection{},
			CreatedAt:   time.Now(),
		},
	}
}

func (v *VisualDeploymentManager) CreateComponent(componentType, name string, x, y int, config map[string]string) (*Component, error) {
	component := &Component{
		ID:   fmt.Sprintf("comp-%d", time.Now().Unix()),
		Type: componentType,
		Name: name,
		Position: Position{X: x, Y: y},
		Size:     v.getDefaultSize(componentType),
		Config:   config,
		Status:   "inactive",
		Ports:    v.getDefaultPorts(componentType),
	}
	
	v.components[component.ID] = component
	v.canvas.Components = append(v.canvas.Components, *component)
	v.canvas.UpdatedAt = time.Now()
	
	return component, nil
}

func (v *VisualDeploymentManager) getDefaultSize(componentType string) Size {
	switch componentType {
	case "process":
		return Size{Width: 120, Height: 80}
	case "database":
		return Size{Width: 100, Height: 100}
	case "load_balancer":
		return Size{Width: 140, Height: 60}
	case "cache":
		return Size{Width: 80, Height: 80}
	default:
		return Size{Width: 100, Height: 80}
	}
}

func (v *VisualDeploymentManager) getDefaultPorts(componentType string) []Port {
	switch componentType {
	case "process":
		return []Port{
			{ID: "input", Name: "Input", Type: "data", Direction: "in", Position: Position{X: 0, Y: 40}},
			{ID: "output", Name: "Output", Type: "data", Direction: "out", Position: Position{X: 120, Y: 40}},
		}
	case "database":
		return []Port{
			{ID: "query", Name: "Query", Type: "data", Direction: "in", Position: Position{X: 50, Y: 0}},
			{ID: "result", Name: "Result", Type: "data", Direction: "out", Position: Position{X: 50, Y: 100}},
		}
	case "load_balancer":
		return []Port{
			{ID: "requests", Name: "Requests", Type: "http", Direction: "in", Position: Position{X: 0, Y: 30}},
			{ID: "backend1", Name: "Backend 1", Type: "http", Direction: "out", Position: Position{X: 140, Y: 20}},
			{ID: "backend2", Name: "Backend 2", Type: "http", Direction: "out", Position: Position{X: 140, Y: 40}},
		}
	default:
		return []Port{}
	}
}

func (v *VisualDeploymentManager) ConnectComponents(fromComponentID, fromPort, toComponentID, toPort, connectionType string) (*Connection, error) {
	fromComp, fromExists := v.components[fromComponentID]
	toComp, toExists := v.components[toComponentID]
	
	if !fromExists || !toExists {
		return nil, fmt.Errorf("component not found")
	}
	
	connection := &Connection{
		ID:       fmt.Sprintf("conn-%d", time.Now().Unix()),
		FromPort: fmt.Sprintf("%s:%s", fromComponentID, fromPort),
		ToPort:   fmt.Sprintf("%s:%s", toComponentID, toPort),
		Type:     connectionType,
		Label:    fmt.Sprintf("%s -> %s", fromComp.Name, toComp.Name),
	}
	
	v.connections[connection.ID] = connection
	v.canvas.Connections = append(v.canvas.Connections, *connection)
	v.canvas.UpdatedAt = time.Now()
	
	return connection, nil
}

func (v *VisualDeploymentManager) DeployCanvas() error {
	fmt.Printf("Deploying visual canvas: %s\n", v.canvas.Name)
	
	// Deploy components in dependency order
	deployOrder := v.calculateDeploymentOrder()
	
	for _, componentID := range deployOrder {
		component := v.components[componentID]
		if err := v.deployComponent(component); err != nil {
			return fmt.Errorf("failed to deploy component %s: %v", component.Name, err)
		}
	}
	
	// Establish connections
	for _, connection := range v.connections {
		if err := v.establishConnection(connection); err != nil {
			return fmt.Errorf("failed to establish connection %s: %v", connection.ID, err)
		}
	}
	
	return nil
}

func (v *VisualDeploymentManager) calculateDeploymentOrder() []string {
	// Simple topological sort based on connections
	var order []string
	deployed := make(map[string]bool)
	
	// Deploy components with no dependencies first
	for id, component := range v.components {
		if !v.hasDependencies(id) {
			order = append(order, id)
			deployed[id] = true
			component.Status = "deploying"
		}
	}
	
	// Deploy remaining components
	for len(deployed) < len(v.components) {
		for id, component := range v.components {
			if deployed[id] {
				continue
			}
			
			if v.dependenciesDeployed(id, deployed) {
				order = append(order, id)
				deployed[id] = true
				component.Status = "deploying"
			}
		}
	}
	
	return order
}

func (v *VisualDeploymentManager) hasDependencies(componentID string) bool {
	for _, connection := range v.connections {
		if v.getComponentFromPort(connection.ToPort) == componentID {
			return true
		}
	}
	return false
}

func (v *VisualDeploymentManager) dependenciesDeployed(componentID string, deployed map[string]bool) bool {
	for _, connection := range v.connections {
		if v.getComponentFromPort(connection.ToPort) == componentID {
			depID := v.getComponentFromPort(connection.FromPort)
			if !deployed[depID] {
				return false
			}
		}
	}
	return true
}

func (v *VisualDeploymentManager) getComponentFromPort(port string) string {
	// Extract component ID from port string (format: "componentID:portName")
	for i, char := range port {
		if char == ':' {
			return port[:i]
		}
	}
	return port
}

func (v *VisualDeploymentManager) deployComponent(component *Component) error {
	fmt.Printf("Deploying component: %s (%s)\n", component.Name, component.Type)
	
	switch component.Type {
	case "process":
		return v.deployProcess(component)
	case "database":
		return v.deployDatabase(component)
	case "load_balancer":
		return v.deployLoadBalancer(component)
	case "cache":
		return v.deployCache(component)
	default:
		return fmt.Errorf("unknown component type: %s", component.Type)
	}
}

func (v *VisualDeploymentManager) deployProcess(component *Component) error {
	command := component.Config["command"]
	if command == "" {
		return fmt.Errorf("process command not specified")
	}
	
	// Deploy process using GProc
	fmt.Printf("Starting process: %s\n", command)
	component.Status = "running"
	return nil
}

func (v *VisualDeploymentManager) deployDatabase(component *Component) error {
	dbType := component.Config["type"]
	if dbType == "" {
		dbType = "postgresql"
	}
	
	fmt.Printf("Starting database: %s\n", dbType)
	component.Status = "running"
	return nil
}

func (v *VisualDeploymentManager) deployLoadBalancer(component *Component) error {
	algorithm := component.Config["algorithm"]
	if algorithm == "" {
		algorithm = "round_robin"
	}
	
	fmt.Printf("Starting load balancer with algorithm: %s\n", algorithm)
	component.Status = "running"
	return nil
}

func (v *VisualDeploymentManager) deployCache(component *Component) error {
	cacheType := component.Config["type"]
	if cacheType == "" {
		cacheType = "redis"
	}
	
	fmt.Printf("Starting cache: %s\n", cacheType)
	component.Status = "running"
	return nil
}

func (v *VisualDeploymentManager) establishConnection(connection *Connection) error {
	fmt.Printf("Establishing connection: %s\n", connection.Label)
	
	// Configure networking, service discovery, etc.
	switch connection.Type {
	case "http":
		fmt.Printf("Configuring HTTP connection\n")
	case "tcp":
		fmt.Printf("Configuring TCP connection\n")
	case "database":
		fmt.Printf("Configuring database connection\n")
	}
	
	return nil
}

func (v *VisualDeploymentManager) GetCanvasJSON() (string, error) {
	data, err := json.MarshalIndent(v.canvas, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (v *VisualDeploymentManager) LoadCanvasFromJSON(jsonData string) error {
	var canvas DeploymentCanvas
	if err := json.Unmarshal([]byte(jsonData), &canvas); err != nil {
		return err
	}
	
	v.canvas = &canvas
	
	// Rebuild component and connection maps
	v.components = make(map[string]*Component)
	v.connections = make(map[string]*Connection)
	
	for i := range canvas.Components {
		comp := &canvas.Components[i]
		v.components[comp.ID] = comp
	}
	
	for i := range canvas.Connections {
		conn := &canvas.Connections[i]
		v.connections[conn.ID] = conn
	}
	
	return nil
}

func (v *VisualDeploymentManager) GetComponentStatus() map[string]string {
	status := make(map[string]string)
	for id, component := range v.components {
		status[id] = component.Status
	}
	return status
}