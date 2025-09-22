package marketplace

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gproc/pkg/types"
)

type MarketplaceManager struct {
	config   *types.MarketplaceConfig
	registry *PluginRegistry
	store    *PluginStore
}

type PluginRegistry struct {
	plugins map[string]*MarketplacePlugin
	authors map[string]*Author
}

type MarketplacePlugin struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Author      string            `json:"author"`
	Description string            `json:"description"`
	Category    string            `json:"category"`
	Tags        []string          `json:"tags"`
	Downloads   int               `json:"downloads"`
	Rating      float64           `json:"rating"`
	Reviews     []Review          `json:"reviews"`
	Repository  string            `json:"repository"`
	Homepage    string            `json:"homepage"`
	License     string            `json:"license"`
	Dependencies []string         `json:"dependencies"`
	Config      map[string]string `json:"config"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Verified    bool              `json:"verified"`
}

type Author struct {
	ID       string    `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Website  string    `json:"website"`
	Plugins  []string  `json:"plugins"`
	Verified bool      `json:"verified"`
	JoinedAt time.Time `json:"joined_at"`
}

type Review struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

type PluginStore struct {
	installed map[string]*InstalledPlugin
}

type InstalledPlugin struct {
	ID          string    `json:"id"`
	Version     string    `json:"version"`
	InstallPath string    `json:"install_path"`
	Enabled     bool      `json:"enabled"`
	InstalledAt time.Time `json:"installed_at"`
	Config      map[string]string `json:"config"`
}

func NewMarketplaceManager(config *types.MarketplaceConfig) *MarketplaceManager {
	return &MarketplaceManager{
		config: config,
		registry: &PluginRegistry{
			plugins: make(map[string]*MarketplacePlugin),
			authors: make(map[string]*Author),
		},
		store: &PluginStore{
			installed: make(map[string]*InstalledPlugin),
		},
	}
}

func (m *MarketplaceManager) SearchPlugins(ctx context.Context, query string, category string) ([]*MarketplacePlugin, error) {
	var results []*MarketplacePlugin
	
	for _, plugin := range m.registry.plugins {
		if m.matchesSearch(plugin, query, category) {
			results = append(results, plugin)
		}
	}
	
	return results, nil
}

func (m *MarketplaceManager) matchesSearch(plugin *MarketplacePlugin, query, category string) bool {
	if category != "" && plugin.Category != category {
		return false
	}
	
	if query == "" {
		return true
	}
	
	// Search in name, description, tags
	if contains(plugin.Name, query) || contains(plugin.Description, query) {
		return true
	}
	
	for _, tag := range plugin.Tags {
		if contains(tag, query) {
			return true
		}
	}
	
	return false
}

func (m *MarketplaceManager) InstallPlugin(ctx context.Context, pluginID, version string) error {
	plugin, exists := m.registry.plugins[pluginID]
	if !exists {
		return fmt.Errorf("plugin not found: %s", pluginID)
	}
	
	// Check dependencies
	for _, dep := range plugin.Dependencies {
		if _, installed := m.store.installed[dep]; !installed {
			return fmt.Errorf("missing dependency: %s", dep)
		}
	}
	
	// Download and install
	installPath := fmt.Sprintf("/plugins/%s", pluginID)
	
	installed := &InstalledPlugin{
		ID:          pluginID,
		Version:     version,
		InstallPath: installPath,
		Enabled:     true,
		InstalledAt: time.Now(),
		Config:      make(map[string]string),
	}
	
	m.store.installed[pluginID] = installed
	plugin.Downloads++
	
	fmt.Printf("Plugin %s installed successfully\n", pluginID)
	return nil
}

func (m *MarketplaceManager) UninstallPlugin(ctx context.Context, pluginID string) error {
	installed, exists := m.store.installed[pluginID]
	if !exists {
		return fmt.Errorf("plugin not installed: %s", pluginID)
	}
	
	// Check if other plugins depend on this
	for _, other := range m.store.installed {
		if other.ID == pluginID {
			continue
		}
		
		plugin := m.registry.plugins[other.ID]
		for _, dep := range plugin.Dependencies {
			if dep == pluginID {
				return fmt.Errorf("plugin %s depends on %s", other.ID, pluginID)
			}
		}
	}
	
	delete(m.store.installed, pluginID)
	fmt.Printf("Plugin %s uninstalled from %s\n", pluginID, installed.InstallPath)
	return nil
}

func (m *MarketplaceManager) PublishPlugin(ctx context.Context, plugin *MarketplacePlugin) error {
	// Validate plugin
	if err := m.validatePlugin(plugin); err != nil {
		return fmt.Errorf("plugin validation failed: %v", err)
	}
	
	plugin.ID = fmt.Sprintf("%s-%s", plugin.Name, plugin.Author)
	plugin.CreatedAt = time.Now()
	plugin.UpdatedAt = time.Now()
	plugin.Downloads = 0
	plugin.Rating = 0.0
	plugin.Verified = false
	
	m.registry.plugins[plugin.ID] = plugin
	
	// Add to author's plugins
	if author, exists := m.registry.authors[plugin.Author]; exists {
		author.Plugins = append(author.Plugins, plugin.ID)
	}
	
	fmt.Printf("Plugin %s published to marketplace\n", plugin.ID)
	return nil
}

func (m *MarketplaceManager) validatePlugin(plugin *MarketplacePlugin) error {
	if plugin.Name == "" {
		return fmt.Errorf("plugin name is required")
	}
	if plugin.Version == "" {
		return fmt.Errorf("plugin version is required")
	}
	if plugin.Author == "" {
		return fmt.Errorf("plugin author is required")
	}
	return nil
}

func (m *MarketplaceManager) AddReview(ctx context.Context, pluginID, userID, username string, rating int, comment string) error {
	plugin, exists := m.registry.plugins[pluginID]
	if !exists {
		return fmt.Errorf("plugin not found: %s", pluginID)
	}
	
	review := Review{
		ID:        fmt.Sprintf("review-%d", time.Now().Unix()),
		UserID:    userID,
		Username:  username,
		Rating:    rating,
		Comment:   comment,
		CreatedAt: time.Now(),
	}
	
	plugin.Reviews = append(plugin.Reviews, review)
	
	// Recalculate rating
	totalRating := 0
	for _, r := range plugin.Reviews {
		totalRating += r.Rating
	}
	plugin.Rating = float64(totalRating) / float64(len(plugin.Reviews))
	
	return nil
}

func (m *MarketplaceManager) GetPopularPlugins(limit int) []*MarketplacePlugin {
	plugins := make([]*MarketplacePlugin, 0, len(m.registry.plugins))
	for _, plugin := range m.registry.plugins {
		plugins = append(plugins, plugin)
	}
	
	// Sort by downloads
	for i := 0; i < len(plugins)-1; i++ {
		for j := i + 1; j < len(plugins); j++ {
			if plugins[i].Downloads < plugins[j].Downloads {
				plugins[i], plugins[j] = plugins[j], plugins[i]
			}
		}
	}
	
	if limit > 0 && len(plugins) > limit {
		plugins = plugins[:limit]
	}
	
	return plugins
}

func (m *MarketplaceManager) GetInstalledPlugins() []*InstalledPlugin {
	plugins := make([]*InstalledPlugin, 0, len(m.store.installed))
	for _, plugin := range m.store.installed {
		plugins = append(plugins, plugin)
	}
	return plugins
}

func (m *MarketplaceManager) EnablePlugin(pluginID string) error {
	installed, exists := m.store.installed[pluginID]
	if !exists {
		return fmt.Errorf("plugin not installed: %s", pluginID)
	}
	
	installed.Enabled = true
	fmt.Printf("Plugin %s enabled\n", pluginID)
	return nil
}

func (m *MarketplaceManager) DisablePlugin(pluginID string) error {
	installed, exists := m.store.installed[pluginID]
	if !exists {
		return fmt.Errorf("plugin not installed: %s", pluginID)
	}
	
	installed.Enabled = false
	fmt.Printf("Plugin %s disabled\n", pluginID)
	return nil
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0)
}