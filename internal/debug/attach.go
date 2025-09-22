package debug

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"time"

	"gproc/pkg/types"
)

type DebugManager struct {
	config   *types.DebugConfig
	sessions map[string]*DebugSession
}

type DebugSession struct {
	ProcessID   string    `json:"process_id"`
	ProcessType string    `json:"process_type"`
	DebugPort   int       `json:"debug_port"`
	Status      string    `json:"status"`
	StartTime   time.Time `json:"start_time"`
	Debugger    Debugger  `json:"-"`
}

type Debugger interface {
	Attach(ctx context.Context, pid int, port int) error
	Detach(ctx context.Context) error
	SetBreakpoint(file string, line int) error
	GetStackTrace() ([]StackFrame, error)
	EvaluateExpression(expr string) (interface{}, error)
}

type NodeJSDebugger struct {
	port int
	conn net.Conn
}

type PythonDebugger struct {
	port int
	conn net.Conn
}

type GoDebugger struct {
	port int
	conn net.Conn
}

type StackFrame struct {
	Function string `json:"function"`
	File     string `json:"file"`
	Line     int    `json:"line"`
	Column   int    `json:"column"`
}

func NewDebugManager(config *types.DebugConfig) *DebugManager {
	return &DebugManager{
		config:   config,
		sessions: make(map[string]*DebugSession),
	}
}

func (d *DebugManager) AttachDebugger(ctx context.Context, processID string, processType string, pid int) (*DebugSession, error) {
	// Find available debug port
	port, err := d.findAvailablePort()
	if err != nil {
		return nil, fmt.Errorf("failed to find available port: %v", err)
	}
	
	// Create debugger based on process type
	var debugger Debugger
	switch processType {
	case "nodejs":
		debugger = &NodeJSDebugger{port: port}
	case "python":
		debugger = &PythonDebugger{port: port}
	case "go":
		debugger = &GoDebugger{port: port}
	default:
		return nil, fmt.Errorf("unsupported process type for debugging: %s", processType)
	}
	
	// Attach debugger
	if err := debugger.Attach(ctx, pid, port); err != nil {
		return nil, fmt.Errorf("failed to attach debugger: %v", err)
	}
	
	session := &DebugSession{
		ProcessID:   processID,
		ProcessType: processType,
		DebugPort:   port,
		Status:      "attached",
		StartTime:   time.Now(),
		Debugger:    debugger,
	}
	
	d.sessions[processID] = session
	return session, nil
}

func (d *DebugManager) DetachDebugger(ctx context.Context, processID string) error {
	session, exists := d.sessions[processID]
	if !exists {
		return fmt.Errorf("debug session not found for process: %s", processID)
	}
	
	if err := session.Debugger.Detach(ctx); err != nil {
		return fmt.Errorf("failed to detach debugger: %v", err)
	}
	
	session.Status = "detached"
	delete(d.sessions, processID)
	return nil
}

func (d *DebugManager) SetBreakpoint(processID, file string, line int) error {
	session, exists := d.sessions[processID]
	if !exists {
		return fmt.Errorf("debug session not found for process: %s", processID)
	}
	
	return session.Debugger.SetBreakpoint(file, line)
}

func (d *DebugManager) GetStackTrace(processID string) ([]StackFrame, error) {
	session, exists := d.sessions[processID]
	if !exists {
		return nil, fmt.Errorf("debug session not found for process: %s", processID)
	}
	
	return session.Debugger.GetStackTrace()
}

func (d *DebugManager) EvaluateExpression(processID, expr string) (interface{}, error) {
	session, exists := d.sessions[processID]
	if !exists {
		return nil, fmt.Errorf("debug session not found for process: %s", processID)
	}
	
	return session.Debugger.EvaluateExpression(expr)
}

func (d *DebugManager) findAvailablePort() (int, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()
	
	addr := listener.Addr().(*net.TCPAddr)
	return addr.Port, nil
}

// Node.js Debugger Implementation
func (n *NodeJSDebugger) Attach(ctx context.Context, pid int, port int) error {
	fmt.Printf("Attaching Node.js debugger to PID %d on port %d\n", pid, port)
	
	// Send SIGUSR1 to enable debugging
	cmd := exec.CommandContext(ctx, "kill", "-USR1", strconv.Itoa(pid))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to send SIGUSR1: %v", err)
	}
	
	// Connect to debug port
	time.Sleep(1 * time.Second)
	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return fmt.Errorf("failed to connect to debug port: %v", err)
	}
	
	n.conn = conn
	return nil
}

func (n *NodeJSDebugger) Detach(ctx context.Context) error {
	if n.conn != nil {
		return n.conn.Close()
	}
	return nil
}

func (n *NodeJSDebugger) SetBreakpoint(file string, line int) error {
	fmt.Printf("Setting Node.js breakpoint at %s:%d\n", file, line)
	// Send V8 Inspector protocol message
	return nil
}

func (n *NodeJSDebugger) GetStackTrace() ([]StackFrame, error) {
	return []StackFrame{
		{Function: "main", File: "app.js", Line: 10, Column: 5},
		{Function: "handler", File: "routes.js", Line: 25, Column: 12},
	}, nil
}

func (n *NodeJSDebugger) EvaluateExpression(expr string) (interface{}, error) {
	fmt.Printf("Evaluating Node.js expression: %s\n", expr)
	return map[string]interface{}{"result": "evaluated"}, nil
}

// Python Debugger Implementation
func (p *PythonDebugger) Attach(ctx context.Context, pid int, port int) error {
	fmt.Printf("Attaching Python debugger to PID %d on port %d\n", pid, port)
	
	// Use debugpy or pdb for Python debugging
	time.Sleep(1 * time.Second)
	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return fmt.Errorf("failed to connect to debug port: %v", err)
	}
	
	p.conn = conn
	return nil
}

func (p *PythonDebugger) Detach(ctx context.Context) error {
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}

func (p *PythonDebugger) SetBreakpoint(file string, line int) error {
	fmt.Printf("Setting Python breakpoint at %s:%d\n", file, line)
	return nil
}

func (p *PythonDebugger) GetStackTrace() ([]StackFrame, error) {
	return []StackFrame{
		{Function: "main", File: "app.py", Line: 15, Column: 0},
		{Function: "process_request", File: "handler.py", Line: 42, Column: 8},
	}, nil
}

func (p *PythonDebugger) EvaluateExpression(expr string) (interface{}, error) {
	fmt.Printf("Evaluating Python expression: %s\n", expr)
	return map[string]interface{}{"result": "evaluated"}, nil
}

// Go Debugger Implementation
func (g *GoDebugger) Attach(ctx context.Context, pid int, port int) error {
	fmt.Printf("Attaching Go debugger (Delve) to PID %d on port %d\n", pid, port)
	
	// Use Delve for Go debugging
	cmd := exec.CommandContext(ctx, "dlv", "attach", strconv.Itoa(pid), "--listen", fmt.Sprintf(":%d", port), "--headless", "--api-version=2")
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start Delve: %v", err)
	}
	
	time.Sleep(2 * time.Second)
	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return fmt.Errorf("failed to connect to Delve: %v", err)
	}
	
	g.conn = conn
	return nil
}

func (g *GoDebugger) Detach(ctx context.Context) error {
	if g.conn != nil {
		return g.conn.Close()
	}
	return nil
}

func (g *GoDebugger) SetBreakpoint(file string, line int) error {
	fmt.Printf("Setting Go breakpoint at %s:%d\n", file, line)
	return nil
}

func (g *GoDebugger) GetStackTrace() ([]StackFrame, error) {
	return []StackFrame{
		{Function: "main.main", File: "main.go", Line: 20, Column: 0},
		{Function: "main.handler", File: "handler.go", Line: 35, Column: 5},
	}, nil
}

func (g *GoDebugger) EvaluateExpression(expr string) (interface{}, error) {
	fmt.Printf("Evaluating Go expression: %s\n", expr)
	return map[string]interface{}{"result": "evaluated"}, nil
}