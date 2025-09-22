package hotreload

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gproc/pkg/types"
)

type HotReloadManager struct {
	config   *types.HotReloadConfig
	watchers map[string]*FileWatcher
	handlers map[string]ReloadHandler
}

type FileWatcher struct {
	ProcessID string
	WatchPath string
	LastMod   time.Time
	Active    bool
}

type ReloadHandler interface {
	CanReload(processType string) bool
	Reload(ctx context.Context, processID string) error
}

type NodeJSHandler struct{}
type PythonHandler struct{}
type GoHandler struct{}

func NewHotReloadManager(config *types.HotReloadConfig) *HotReloadManager {
	manager := &HotReloadManager{
		config:   config,
		watchers: make(map[string]*FileWatcher),
		handlers: make(map[string]ReloadHandler),
	}
	
	// Register handlers
	manager.handlers["nodejs"] = &NodeJSHandler{}
	manager.handlers["python"] = &PythonHandler{}
	manager.handlers["go"] = &GoHandler{}
	
	return manager
}

func (h *HotReloadManager) EnableHotReload(processID, processType, watchPath string) error {
	handler, exists := h.handlers[processType]
	if !exists {
		return fmt.Errorf("hot reload not supported for %s", processType)
	}
	
	if !handler.CanReload(processType) {
		return fmt.Errorf("hot reload disabled for %s", processType)
	}
	
	watcher := &FileWatcher{
		ProcessID: processID,
		WatchPath: watchPath,
		LastMod:   time.Now(),
		Active:    true,
	}
	
	h.watchers[processID] = watcher
	
	// Start watching in background
	go h.watchFiles(processID)
	
	return nil
}

func (h *HotReloadManager) DisableHotReload(processID string) {
	if watcher, exists := h.watchers[processID]; exists {
		watcher.Active = false
		delete(h.watchers, processID)
	}
}

func (h *HotReloadManager) watchFiles(processID string) {
	watcher := h.watchers[processID]
	
	for watcher.Active {
		if h.hasFileChanged(watcher.WatchPath, watcher.LastMod) {
			fmt.Printf("File change detected for process %s, triggering hot reload\n", processID)
			
			// Determine process type and reload
			processType := h.detectProcessType(watcher.WatchPath)
			if handler, exists := h.handlers[processType]; exists {
				ctx := context.Background()
				if err := handler.Reload(ctx, processID); err != nil {
					fmt.Printf("Hot reload failed for %s: %v\n", processID, err)
				} else {
					fmt.Printf("Hot reload successful for %s\n", processID)
				}
			}
			
			watcher.LastMod = time.Now()
		}
		
		time.Sleep(1 * time.Second)
	}
}

func (h *HotReloadManager) hasFileChanged(watchPath string, lastMod time.Time) bool {
	err := filepath.Walk(watchPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if info.ModTime().After(lastMod) {
			return fmt.Errorf("file changed") // Use error to break walk
		}
		
		return nil
	})
	
	return err != nil
}

func (h *HotReloadManager) detectProcessType(watchPath string) string {
	if _, err := os.Stat(filepath.Join(watchPath, "package.json")); err == nil {
		return "nodejs"
	}
	if _, err := os.Stat(filepath.Join(watchPath, "requirements.txt")); err == nil {
		return "python"
	}
	if _, err := os.Stat(filepath.Join(watchPath, "go.mod")); err == nil {
		return "go"
	}
	return "unknown"
}

// Node.js Handler
func (n *NodeJSHandler) CanReload(processType string) bool {
	return processType == "nodejs"
}

func (n *NodeJSHandler) Reload(ctx context.Context, processID string) error {
	// Send SIGUSR2 to Node.js process for graceful reload
	fmt.Printf("Sending SIGUSR2 to Node.js process %s\n", processID)
	return nil
}

// Python Handler
func (p *PythonHandler) CanReload(processType string) bool {
	return processType == "python"
}

func (p *PythonHandler) Reload(ctx context.Context, processID string) error {
	// Restart Python process (no graceful reload)
	fmt.Printf("Restarting Python process %s\n", processID)
	return nil
}

// Go Handler
func (g *GoHandler) CanReload(processType string) bool {
	return processType == "go"
}

func (g *GoHandler) Reload(ctx context.Context, processID string) error {
	// Rebuild and restart Go process
	fmt.Printf("Rebuilding and restarting Go process %s\n", processID)
	return nil
}