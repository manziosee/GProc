package gitops

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"

	"gproc/pkg/types"
)

type GitOpsManager struct {
	config      *types.GitOpsConfig
	repositories map[string]*Repository
	syncInterval time.Duration
}

type Repository struct {
	Name       string    `json:"name"`
	URL        string    `json:"url"`
	Branch     string    `json:"branch"`
	LocalPath  string    `json:"local_path"`
	LastSync   time.Time `json:"last_sync"`
	LastCommit string    `json:"last_commit"`
}

type GitOpsConfig struct {
	Enabled      bool              `json:"enabled"`
	Repositories []Repository      `json:"repositories"`
	SyncInterval time.Duration     `json:"sync_interval"`
	AutoDeploy   bool              `json:"auto_deploy"`
	ConfigPaths  map[string]string `json:"config_paths"`
}

func NewGitOpsManager(config *types.GitOpsConfig) *GitOpsManager {
	return &GitOpsManager{
		config:       config,
		repositories: make(map[string]*Repository),
		syncInterval: 30 * time.Second,
	}
}

func (g *GitOpsManager) Start(ctx context.Context) error {
	if !g.config.Enabled {
		return nil
	}
	
	// Initialize repositories
	for _, repo := range g.config.Repositories {
		if err := g.initRepository(ctx, &repo); err != nil {
			return fmt.Errorf("failed to initialize repository %s: %v", repo.Name, err)
		}
		g.repositories[repo.Name] = &repo
	}
	
	// Start sync loop
	go g.syncLoop(ctx)
	
	return nil
}

func (g *GitOpsManager) initRepository(ctx context.Context, repo *Repository) error {
	// Clone repository if not exists
	if !g.repositoryExists(repo.LocalPath) {
		if err := g.cloneRepository(ctx, repo); err != nil {
			return err
		}
	}
	
	// Get current commit
	commit, err := g.getCurrentCommit(repo.LocalPath)
	if err != nil {
		return err
	}
	repo.LastCommit = commit
	
	return nil
}

func (g *GitOpsManager) syncLoop(ctx context.Context) {
	ticker := time.NewTicker(g.syncInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			g.syncAllRepositories(ctx)
		}
	}
}

func (g *GitOpsManager) syncAllRepositories(ctx context.Context) {
	for _, repo := range g.repositories {
		if err := g.syncRepository(ctx, repo); err != nil {
			fmt.Printf("Failed to sync repository %s: %v\n", repo.Name, err)
		}
	}
}

func (g *GitOpsManager) syncRepository(ctx context.Context, repo *Repository) error {
	// Pull latest changes
	if err := g.pullRepository(ctx, repo); err != nil {
		return err
	}
	
	// Check for new commits
	currentCommit, err := g.getCurrentCommit(repo.LocalPath)
	if err != nil {
		return err
	}
	
	if currentCommit != repo.LastCommit {
		fmt.Printf("New commit detected in %s: %s -> %s\n", repo.Name, repo.LastCommit, currentCommit)
		
		// Apply configuration changes
		if err := g.applyConfigChanges(ctx, repo); err != nil {
			return err
		}
		
		repo.LastCommit = currentCommit
		repo.LastSync = time.Now()
	}
	
	return nil
}

func (g *GitOpsManager) applyConfigChanges(ctx context.Context, repo *Repository) error {
	// Find configuration files
	configFiles, err := g.findConfigFiles(repo.LocalPath)
	if err != nil {
		return err
	}
	
	// Apply each configuration
	for _, configFile := range configFiles {
		if err := g.applyConfig(ctx, configFile); err != nil {
			fmt.Printf("Failed to apply config %s: %v\n", configFile, err)
		} else {
			fmt.Printf("Applied configuration: %s\n", configFile)
		}
	}
	
	return nil
}

func (g *GitOpsManager) findConfigFiles(repoPath string) ([]string, error) {
	var configFiles []string
	
	// Look for gproc.yaml files
	matches, err := filepath.Glob(filepath.Join(repoPath, "**", "gproc.yaml"))
	if err != nil {
		return nil, err
	}
	configFiles = append(configFiles, matches...)
	
	// Look for process definitions
	matches, err = filepath.Glob(filepath.Join(repoPath, "**", "processes.yaml"))
	if err != nil {
		return nil, err
	}
	configFiles = append(configFiles, matches...)
	
	return configFiles, nil
}

func (g *GitOpsManager) applyConfig(ctx context.Context, configFile string) error {
	// Parse and apply configuration
	fmt.Printf("Applying configuration from %s\n", configFile)
	
	// In real implementation:
	// 1. Parse YAML configuration
	// 2. Validate configuration
	// 3. Apply process changes
	// 4. Update running processes
	
	return nil
}

func (g *GitOpsManager) repositoryExists(path string) bool {
	gitDir := filepath.Join(path, ".git")
	if _, err := filepath.Glob(gitDir); err != nil {
		return false
	}
	return true
}

func (g *GitOpsManager) cloneRepository(ctx context.Context, repo *Repository) error {
	cmd := exec.CommandContext(ctx, "git", "clone", "-b", repo.Branch, repo.URL, repo.LocalPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone repository: %v", err)
	}
	return nil
}

func (g *GitOpsManager) pullRepository(ctx context.Context, repo *Repository) error {
	cmd := exec.CommandContext(ctx, "git", "-C", repo.LocalPath, "pull", "origin", repo.Branch)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to pull repository: %v", err)
	}
	return nil
}

func (g *GitOpsManager) getCurrentCommit(repoPath string) (string, error) {
	cmd := exec.Command("git", "-C", repoPath, "rev-parse", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output[:7]), nil // Return short commit hash
}