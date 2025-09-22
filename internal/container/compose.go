package container

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"gproc/pkg/types"
)

type ComposeManager struct {
	config    *types.ComposeConfig
	stacks    map[string]*ComposeStack
	dockerCmd string
}

type ComposeStack struct {
	Name         string                 `json:"name"`
	FilePath     string                 `json:"file_path"`
	Services     map[string]*Service    `json:"services"`
	Status       string                 `json:"status"`
	CreatedAt    time.Time              `json:"created_at"`
	Environment  map[string]string      `json:"environment"`
}

type Service struct {
	Name      string            `json:"name"`
	Image     string            `json:"image"`
	Status    string            `json:"status"`
	Ports     []string          `json:"ports"`
	Volumes   []string          `json:"volumes"`
	Env       map[string]string `json:"env"`
	DependsOn []string          `json:"depends_on"`
}

type ComposeConfig struct {
	Enabled     bool   `json:"enabled"`
	DefaultPath string `json:"default_path"`
	ProjectName string `json:"project_name"`
}

func NewComposeManager(config *types.ComposeConfig) (*ComposeManager, error) {
	dockerCmd, err := exec.LookPath("docker")
	if err != nil {
		return nil, fmt.Errorf("docker command not found: %v", err)
	}
	
	return &ComposeManager{
		config:    config,
		stacks:    make(map[string]*ComposeStack),
		dockerCmd: dockerCmd,
	}, nil
}

func (c *ComposeManager) DeployStack(ctx context.Context, name, filePath string, env map[string]string) error {
	if !c.config.Enabled {
		return fmt.Errorf("docker compose support disabled")
	}
	
	// Validate compose file exists
	if !filepath.IsAbs(filePath) {
		filePath = filepath.Join(c.config.DefaultPath, filePath)
	}
	
	// Create stack entry
	stack := &ComposeStack{
		Name:        name,
		FilePath:    filePath,
		Services:    make(map[string]*Service),
		Status:      "deploying",
		CreatedAt:   time.Now(),
		Environment: env,
	}
	
	c.stacks[name] = stack
	
	// Deploy using docker-compose
	if err := c.executeComposeCommand(ctx, "up", "-d", "--project-name", name, "-f", filePath); err != nil {
		stack.Status = "failed"
		return fmt.Errorf("failed to deploy stack %s: %v", name, err)
	}
	
	// Update stack status and services
	if err := c.updateStackStatus(ctx, name); err != nil {
		return err
	}
	
	stack.Status = "running"
	return nil
}

func (c *ComposeManager) StopStack(ctx context.Context, name string) error {
	stack, exists := c.stacks[name]
	if !exists {
		return fmt.Errorf("stack %s not found", name)
	}
	
	if err := c.executeComposeCommand(ctx, "stop", "--project-name", name, "-f", stack.FilePath); err != nil {
		return fmt.Errorf("failed to stop stack %s: %v", name, err)
	}
	
	stack.Status = "stopped"
	return nil
}

func (c *ComposeManager) RemoveStack(ctx context.Context, name string) error {
	stack, exists := c.stacks[name]
	if !exists {
		return fmt.Errorf("stack %s not found", name)
	}
	
	if err := c.executeComposeCommand(ctx, "down", "--project-name", name, "-f", stack.FilePath); err != nil {
		return fmt.Errorf("failed to remove stack %s: %v", name, err)
	}
	
	delete(c.stacks, name)
	return nil
}

func (c *ComposeManager) ListStacks() []*ComposeStack {
	stacks := make([]*ComposeStack, 0, len(c.stacks))
	for _, stack := range c.stacks {
		stacks = append(stacks, stack)
	}
	return stacks
}

func (c *ComposeManager) GetStackStatus(ctx context.Context, name string) (*ComposeStack, error) {
	stack, exists := c.stacks[name]
	if !exists {
		return nil, fmt.Errorf("stack %s not found", name)
	}
	
	if err := c.updateStackStatus(ctx, name); err != nil {
		return nil, err
	}
	
	return stack, nil
}

func (c *ComposeManager) ScaleService(ctx context.Context, stackName, serviceName string, replicas int) error {
	stack, exists := c.stacks[stackName]
	if !exists {
		return fmt.Errorf("stack %s not found", stackName)
	}
	
	scaleArg := fmt.Sprintf("%s=%d", serviceName, replicas)
	if err := c.executeComposeCommand(ctx, "up", "-d", "--scale", scaleArg, "--project-name", stackName, "-f", stack.FilePath); err != nil {
		return fmt.Errorf("failed to scale service %s: %v", serviceName, err)
	}
	
	return c.updateStackStatus(ctx, stackName)
}

func (c *ComposeManager) GetServiceLogs(ctx context.Context, stackName, serviceName string, lines int) ([]string, error) {
	stack, exists := c.stacks[stackName]
	if !exists {
		return nil, fmt.Errorf("stack %s not found", stackName)
	}
	
	cmd := exec.CommandContext(ctx, c.dockerCmd, "compose", 
		"--project-name", stackName, 
		"-f", stack.FilePath,
		"logs", "--tail", fmt.Sprintf("%d", lines), serviceName)
	
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get logs: %v", err)
	}
	
	lines_output := strings.Split(string(output), "\n")
	return lines_output, nil
}

func (c *ComposeManager) executeComposeCommand(ctx context.Context, args ...string) error {
	cmd := exec.CommandContext(ctx, c.dockerCmd, append([]string{"compose"}, args...)...)
	
	// Set environment variables
	cmd.Env = append(cmd.Env, "COMPOSE_DOCKER_CLI_BUILD=1")
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("command failed: %v, output: %s", err, string(output))
	}
	
	return nil
}

func (c *ComposeManager) updateStackStatus(ctx context.Context, name string) error {
	stack := c.stacks[name]
	
	// Get service status using docker-compose ps
	cmd := exec.CommandContext(ctx, c.dockerCmd, "compose", 
		"--project-name", name, 
		"-f", stack.FilePath, 
		"ps", "--format", "json")
	
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get stack status: %v", err)
	}
	
	// Parse service status (simplified)
	services := parseComposePS(string(output))
	stack.Services = services
	
	// Update overall stack status
	allRunning := true
	for _, service := range services {
		if service.Status != "running" {
			allRunning = false
			break
		}
	}
	
	if allRunning && len(services) > 0 {
		stack.Status = "running"
	} else if len(services) == 0 {
		stack.Status = "stopped"
	} else {
		stack.Status = "partial"
	}
	
	return nil
}

func parseComposePS(output string) map[string]*Service {
	services := make(map[string]*Service)
	
	// Simplified parsing - in real implementation, parse JSON output
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		
		// Mock service parsing
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			serviceName := parts[0]
			status := "running"
			if len(parts) > 1 && strings.Contains(parts[1], "Exit") {
				status = "stopped"
			}
			
			services[serviceName] = &Service{
				Name:   serviceName,
				Status: status,
				Ports:  []string{},
				Env:    make(map[string]string),
			}
		}
	}
	
	return services
}