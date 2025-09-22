package serverless

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gproc/pkg/types"
)

type ServerlessManager struct {
	config    *types.ServerlessConfig
	providers map[string]ServerlessProvider
	hooks     map[string]*Hook
}

type ServerlessProvider interface {
	InvokeFunction(ctx context.Context, functionName string, payload interface{}) (*InvocationResult, error)
	CreateFunction(ctx context.Context, config *FunctionConfig) error
	DeleteFunction(ctx context.Context, functionName string) error
}

type AWSLambdaProvider struct {
	region    string
	accessKey string
	secretKey string
}

type GCPFunctionsProvider struct {
	projectID   string
	region      string
	credentials string
}

type Hook struct {
	Name         string            `json:"name"`
	Event        string            `json:"event"`
	FunctionName string            `json:"function_name"`
	Provider     string            `json:"provider"`
	Enabled      bool              `json:"enabled"`
	Config       map[string]string `json:"config"`
	LastRun      time.Time         `json:"last_run"`
	RunCount     int               `json:"run_count"`
}

type FunctionConfig struct {
	Name        string            `json:"name"`
	Runtime     string            `json:"runtime"`
	Handler     string            `json:"handler"`
	Code        []byte            `json:"code"`
	Environment map[string]string `json:"environment"`
	Timeout     time.Duration     `json:"timeout"`
	Memory      int               `json:"memory_mb"`
}

type InvocationResult struct {
	StatusCode int                    `json:"status_code"`
	Payload    map[string]interface{} `json:"payload"`
	Duration   time.Duration          `json:"duration"`
	Error      string                 `json:"error,omitempty"`
}

type ProcessEvent struct {
	Type      string                 `json:"type"`
	ProcessID string                 `json:"process_id"`
	Process   *types.Process         `json:"process"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

func NewServerlessManager(config *types.ServerlessConfig) *ServerlessManager {
	manager := &ServerlessManager{
		config:    config,
		providers: make(map[string]ServerlessProvider),
		hooks:     make(map[string]*Hook),
	}
	
	// Initialize AWS Lambda provider
	if config.AWS != nil {
		manager.providers["aws"] = &AWSLambdaProvider{
			region:    config.AWS.Region,
			accessKey: config.AWS.AccessKey,
			secretKey: config.AWS.SecretKey,
		}
	}
	
	// Initialize GCP Functions provider
	if config.GCP != nil {
		manager.providers["gcp"] = &GCPFunctionsProvider{
			projectID:   config.GCP.ProjectID,
			region:      config.GCP.Region,
			credentials: config.GCP.Credentials,
		}
	}
	
	return manager
}

func (s *ServerlessManager) RegisterHook(hook *Hook) error {
	if _, exists := s.providers[hook.Provider]; !exists {
		return fmt.Errorf("provider %s not configured", hook.Provider)
	}
	
	s.hooks[hook.Name] = hook
	return nil
}

func (s *ServerlessManager) TriggerHooks(ctx context.Context, event *ProcessEvent) error {
	for _, hook := range s.hooks {
		if !hook.Enabled || hook.Event != event.Type {
			continue
		}
		
		go s.executeHook(ctx, hook, event)
	}
	
	return nil
}

func (s *ServerlessManager) executeHook(ctx context.Context, hook *Hook, event *ProcessEvent) {
	provider := s.providers[hook.Provider]
	
	hook.LastRun = time.Now()
	hook.RunCount++
	
	result, err := provider.InvokeFunction(ctx, hook.FunctionName, event)
	if err != nil {
		fmt.Printf("Hook %s failed: %v\n", hook.Name, err)
		return
	}
	
	fmt.Printf("Hook %s executed successfully: %d ms\n", hook.Name, result.Duration.Milliseconds())
}

func (s *ServerlessManager) CreateFunction(ctx context.Context, provider string, config *FunctionConfig) error {
	p, exists := s.providers[provider]
	if !exists {
		return fmt.Errorf("provider %s not configured", provider)
	}
	
	return p.CreateFunction(ctx, config)
}

// AWS Lambda Provider Implementation
func (a *AWSLambdaProvider) InvokeFunction(ctx context.Context, functionName string, payload interface{}) (*InvocationResult, error) {
	fmt.Printf("Invoking AWS Lambda function: %s\n", functionName)
	
	// Simulate AWS Lambda invocation
	payloadJSON, _ := json.Marshal(payload)
	fmt.Printf("Payload: %s\n", string(payloadJSON))
	
	// Mock response
	return &InvocationResult{
		StatusCode: 200,
		Payload: map[string]interface{}{
			"message": "Function executed successfully",
			"result":  "processed",
		},
		Duration: 150 * time.Millisecond,
	}, nil
}

func (a *AWSLambdaProvider) CreateFunction(ctx context.Context, config *FunctionConfig) error {
	fmt.Printf("Creating AWS Lambda function: %s\n", config.Name)
	// Simulate function creation via AWS SDK
	return nil
}

func (a *AWSLambdaProvider) DeleteFunction(ctx context.Context, functionName string) error {
	fmt.Printf("Deleting AWS Lambda function: %s\n", functionName)
	return nil
}

// GCP Functions Provider Implementation
func (g *GCPFunctionsProvider) InvokeFunction(ctx context.Context, functionName string, payload interface{}) (*InvocationResult, error) {
	fmt.Printf("Invoking GCP Cloud Function: %s\n", functionName)
	
	// Simulate GCP Functions invocation
	payloadJSON, _ := json.Marshal(payload)
	fmt.Printf("Payload: %s\n", string(payloadJSON))
	
	// Mock response
	return &InvocationResult{
		StatusCode: 200,
		Payload: map[string]interface{}{
			"status": "success",
			"data":   "function_result",
		},
		Duration: 120 * time.Millisecond,
	}, nil
}

func (g *GCPFunctionsProvider) CreateFunction(ctx context.Context, config *FunctionConfig) error {
	fmt.Printf("Creating GCP Cloud Function: %s\n", config.Name)
	// Simulate function creation via GCP SDK
	return nil
}

func (g *GCPFunctionsProvider) DeleteFunction(ctx context.Context, functionName string) error {
	fmt.Printf("Deleting GCP Cloud Function: %s\n", functionName)
	return nil
}

// Helper functions for common hooks
func (s *ServerlessManager) OnProcessStart(processID string) {
	event := &ProcessEvent{
		Type:      "process_start",
		ProcessID: processID,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"action": "started",
		},
	}
	s.TriggerHooks(context.Background(), event)
}

func (s *ServerlessManager) OnProcessStop(processID string) {
	event := &ProcessEvent{
		Type:      "process_stop",
		ProcessID: processID,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"action": "stopped",
		},
	}
	s.TriggerHooks(context.Background(), event)
}

func (s *ServerlessManager) OnProcessFailure(processID string, error string) {
	event := &ProcessEvent{
		Type:      "process_failure",
		ProcessID: processID,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"action": "failed",
			"error":  error,
		},
	}
	s.TriggerHooks(context.Background(), event)
}