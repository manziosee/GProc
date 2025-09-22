package workflows

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"gproc/pkg/types"
)

type WorkflowEngine struct {
	config    *types.WorkflowConfig
	workflows map[string]*Workflow
	triggers  map[string][]*Workflow
}

type Workflow struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Trigger     *Trigger    `json:"trigger"`
	Steps       []Step      `json:"steps"`
	Variables   map[string]string `json:"variables"`
	Enabled     bool        `json:"enabled"`
	CreatedAt   time.Time   `json:"created_at"`
	LastRun     time.Time   `json:"last_run"`
	RunCount    int         `json:"run_count"`
}

type Trigger struct {
	Type       string            `json:"type"`
	Event      string            `json:"event"`
	Conditions map[string]string `json:"conditions"`
	Schedule   string            `json:"schedule,omitempty"`
}

type Step struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	Action      string            `json:"action"`
	Parameters  map[string]string `json:"parameters"`
	Condition   string            `json:"condition,omitempty"`
	OnSuccess   string            `json:"on_success,omitempty"`
	OnFailure   string            `json:"on_failure,omitempty"`
	Timeout     time.Duration     `json:"timeout,omitempty"`
}

type WorkflowExecution struct {
	ID         string                 `json:"id"`
	WorkflowID string                 `json:"workflow_id"`
	Status     string                 `json:"status"`
	StartTime  time.Time              `json:"start_time"`
	EndTime    time.Time              `json:"end_time"`
	Steps      []StepExecution        `json:"steps"`
	Variables  map[string]string      `json:"variables"`
	Error      string                 `json:"error,omitempty"`
}

type StepExecution struct {
	StepID    string            `json:"step_id"`
	Status    string            `json:"status"`
	StartTime time.Time         `json:"start_time"`
	EndTime   time.Time         `json:"end_time"`
	Output    map[string]string `json:"output"`
	Error     string            `json:"error,omitempty"`
}

func NewWorkflowEngine(config *types.WorkflowConfig) *WorkflowEngine {
	return &WorkflowEngine{
		config:    config,
		workflows: make(map[string]*Workflow),
		triggers:  make(map[string][]*Workflow),
	}
}

func (w *WorkflowEngine) RegisterWorkflow(workflow *Workflow) error {
	workflow.ID = fmt.Sprintf("workflow-%d", time.Now().Unix())
	workflow.CreatedAt = time.Now()
	workflow.Enabled = true
	
	w.workflows[workflow.ID] = workflow
	
	// Register trigger
	if workflow.Trigger != nil {
		triggerKey := fmt.Sprintf("%s:%s", workflow.Trigger.Type, workflow.Trigger.Event)
		w.triggers[triggerKey] = append(w.triggers[triggerKey], workflow)
	}
	
	return nil
}

func (w *WorkflowEngine) TriggerWorkflows(ctx context.Context, eventType, event string, data map[string]string) error {
	triggerKey := fmt.Sprintf("%s:%s", eventType, event)
	workflows := w.triggers[triggerKey]
	
	for _, workflow := range workflows {
		if !workflow.Enabled {
			continue
		}
		
		if w.matchesConditions(workflow.Trigger.Conditions, data) {
			go w.executeWorkflow(ctx, workflow, data)
		}
	}
	
	return nil
}

func (w *WorkflowEngine) executeWorkflow(ctx context.Context, workflow *Workflow, triggerData map[string]string) {
	execution := &WorkflowExecution{
		ID:         fmt.Sprintf("exec-%d", time.Now().Unix()),
		WorkflowID: workflow.ID,
		Status:     "running",
		StartTime:  time.Now(),
		Variables:  make(map[string]string),
		Steps:      []StepExecution{},
	}
	
	// Merge workflow variables with trigger data
	for k, v := range workflow.Variables {
		execution.Variables[k] = v
	}
	for k, v := range triggerData {
		execution.Variables[k] = v
	}
	
	workflow.LastRun = time.Now()
	workflow.RunCount++
	
	// Execute steps
	for _, step := range workflow.Steps {
		stepExec := w.executeStep(ctx, step, execution.Variables)
		execution.Steps = append(execution.Steps, stepExec)
		
		if stepExec.Status == "failed" {
			execution.Status = "failed"
			execution.Error = stepExec.Error
			break
		}
		
		// Merge step output into variables
		for k, v := range stepExec.Output {
			execution.Variables[k] = v
		}
	}
	
	if execution.Status != "failed" {
		execution.Status = "completed"
	}
	
	execution.EndTime = time.Now()
	fmt.Printf("Workflow %s execution %s: %s\n", workflow.Name, execution.ID, execution.Status)
}

func (w *WorkflowEngine) executeStep(ctx context.Context, step Step, variables map[string]string) StepExecution {
	stepExec := StepExecution{
		StepID:    step.ID,
		Status:    "running",
		StartTime: time.Now(),
		Output:    make(map[string]string),
	}
	
	// Check condition
	if step.Condition != "" && !w.evaluateCondition(step.Condition, variables) {
		stepExec.Status = "skipped"
		stepExec.EndTime = time.Now()
		return stepExec
	}
	
	// Execute step based on type
	switch step.Type {
	case "process":
		stepExec = w.executeProcessStep(ctx, step, variables)
	case "script":
		stepExec = w.executeScriptStep(ctx, step, variables)
	case "http":
		stepExec = w.executeHTTPStep(ctx, step, variables)
	case "notification":
		stepExec = w.executeNotificationStep(ctx, step, variables)
	default:
		stepExec.Status = "failed"
		stepExec.Error = fmt.Sprintf("unknown step type: %s", step.Type)
	}
	
	stepExec.EndTime = time.Now()
	return stepExec
}

func (w *WorkflowEngine) executeProcessStep(ctx context.Context, step Step, variables map[string]string) StepExecution {
	stepExec := StepExecution{
		StepID:    step.ID,
		Status:    "running",
		StartTime: time.Now(),
		Output:    make(map[string]string),
	}
	
	action := step.Action
	processName := w.substituteVariables(step.Parameters["process"], variables)
	
	switch action {
	case "start":
		fmt.Printf("Starting process: %s\n", processName)
		stepExec.Output["result"] = "started"
	case "stop":
		fmt.Printf("Stopping process: %s\n", processName)
		stepExec.Output["result"] = "stopped"
	case "restart":
		fmt.Printf("Restarting process: %s\n", processName)
		stepExec.Output["result"] = "restarted"
	default:
		stepExec.Status = "failed"
		stepExec.Error = fmt.Sprintf("unknown process action: %s", action)
		return stepExec
	}
	
	stepExec.Status = "completed"
	return stepExec
}

func (w *WorkflowEngine) executeScriptStep(ctx context.Context, step Step, variables map[string]string) StepExecution {
	stepExec := StepExecution{
		StepID:    step.ID,
		Status:    "running",
		StartTime: time.Now(),
		Output:    make(map[string]string),
	}
	
	script := w.substituteVariables(step.Parameters["script"], variables)
	shell := step.Parameters["shell"]
	if shell == "" {
		shell = "bash"
	}
	
	cmd := exec.CommandContext(ctx, shell, "-c", script)
	output, err := cmd.Output()
	
	if err != nil {
		stepExec.Status = "failed"
		stepExec.Error = err.Error()
	} else {
		stepExec.Status = "completed"
		stepExec.Output["output"] = string(output)
	}
	
	return stepExec
}

func (w *WorkflowEngine) executeHTTPStep(ctx context.Context, step Step, variables map[string]string) StepExecution {
	stepExec := StepExecution{
		StepID:    step.ID,
		Status:    "running",
		StartTime: time.Now(),
		Output:    make(map[string]string),
	}
	
	url := w.substituteVariables(step.Parameters["url"], variables)
	method := step.Parameters["method"]
	if method == "" {
		method = "GET"
	}
	
	fmt.Printf("HTTP %s request to: %s\n", method, url)
	
	// Simulate HTTP request
	stepExec.Status = "completed"
	stepExec.Output["status_code"] = "200"
	stepExec.Output["response"] = "success"
	
	return stepExec
}

func (w *WorkflowEngine) executeNotificationStep(ctx context.Context, step Step, variables map[string]string) StepExecution {
	stepExec := StepExecution{
		StepID:    step.ID,
		Status:    "running",
		StartTime: time.Now(),
		Output:    make(map[string]string),
	}
	
	message := w.substituteVariables(step.Parameters["message"], variables)
	channel := step.Parameters["channel"]
	
	fmt.Printf("Sending notification to %s: %s\n", channel, message)
	
	stepExec.Status = "completed"
	stepExec.Output["sent"] = "true"
	
	return stepExec
}

func (w *WorkflowEngine) matchesConditions(conditions map[string]string, data map[string]string) bool {
	for key, expectedValue := range conditions {
		if actualValue, exists := data[key]; !exists || actualValue != expectedValue {
			return false
		}
	}
	return true
}

func (w *WorkflowEngine) evaluateCondition(condition string, variables map[string]string) bool {
	// Simplified condition evaluation
	return true
}

func (w *WorkflowEngine) substituteVariables(text string, variables map[string]string) string {
	result := text
	for key, value := range variables {
		placeholder := fmt.Sprintf("${%s}", key)
		result = replaceAll(result, placeholder, value)
	}
	return result
}

func (w *WorkflowEngine) GetWorkflows() []*Workflow {
	workflows := make([]*Workflow, 0, len(w.workflows))
	for _, workflow := range w.workflows {
		workflows = append(workflows, workflow)
	}
	return workflows
}

func (w *WorkflowEngine) EnableWorkflow(workflowID string) error {
	workflow, exists := w.workflows[workflowID]
	if !exists {
		return fmt.Errorf("workflow not found: %s", workflowID)
	}
	
	workflow.Enabled = true
	return nil
}

func (w *WorkflowEngine) DisableWorkflow(workflowID string) error {
	workflow, exists := w.workflows[workflowID]
	if !exists {
		return fmt.Errorf("workflow not found: %s", workflowID)
	}
	
	workflow.Enabled = false
	return nil
}

func replaceAll(s, old, new string) string {
	// Simple string replacement
	result := ""
	for i := 0; i < len(s); i++ {
		if i+len(old) <= len(s) && s[i:i+len(old)] == old {
			result += new
			i += len(old) - 1
		} else {
			result += string(s[i])
		}
	}
	return result
}