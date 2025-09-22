package dashboards

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gproc/pkg/types"
)

type GrafanaManager struct {
	config     *types.GrafanaConfig
	client     *GrafanaClient
	dashboards map[string]*Dashboard
}

type GrafanaClient struct {
	endpoint string
	apiKey   string
	orgID    int
}

type Dashboard struct {
	ID          int                    `json:"id"`
	UID         string                 `json:"uid"`
	Title       string                 `json:"title"`
	Tags        []string               `json:"tags"`
	Panels      []Panel                `json:"panels"`
	Templating  Templating             `json:"templating"`
	Time        TimeRange              `json:"time"`
	Refresh     string                 `json:"refresh"`
	Version     int                    `json:"version"`
	Meta        map[string]interface{} `json:"meta"`
}

type Panel struct {
	ID          int                    `json:"id"`
	Title       string                 `json:"title"`
	Type        string                 `json:"type"`
	GridPos     GridPos                `json:"gridPos"`
	Targets     []Target               `json:"targets"`
	Options     map[string]interface{} `json:"options"`
	FieldConfig FieldConfig            `json:"fieldConfig"`
}

type GridPos struct {
	H int `json:"h"`
	W int `json:"w"`
	X int `json:"x"`
	Y int `json:"y"`
}

type Target struct {
	Expr         string `json:"expr"`
	RefID        string `json:"refId"`
	Interval     string `json:"interval"`
	LegendFormat string `json:"legendFormat"`
}

type FieldConfig struct {
	Defaults map[string]interface{} `json:"defaults"`
}

type Templating struct {
	List []TemplateVariable `json:"list"`
}

type TemplateVariable struct {
	Name    string   `json:"name"`
	Type    string   `json:"type"`
	Query   string   `json:"query"`
	Options []Option `json:"options"`
}

type Option struct {
	Text  string `json:"text"`
	Value string `json:"value"`
}

type TimeRange struct {
	From string `json:"from"`
	To   string `json:"to"`
}

func NewGrafanaManager(config *types.GrafanaConfig) *GrafanaManager {
	return &GrafanaManager{
		config: config,
		client: &GrafanaClient{
			endpoint: config.Endpoint,
			apiKey:   config.APIKey,
			orgID:    config.OrgID,
		},
		dashboards: make(map[string]*Dashboard),
	}
}

func (g *GrafanaManager) CreateProcessDashboard(ctx context.Context, processID string) (*Dashboard, error) {
	dashboard := &Dashboard{
		UID:     fmt.Sprintf("gproc-process-%s", processID),
		Title:   fmt.Sprintf("GProc Process: %s", processID),
		Tags:    []string{"gproc", "process", processID},
		Refresh: "5s",
		Time: TimeRange{
			From: "now-1h",
			To:   "now",
		},
		Panels: []Panel{
			{
				ID:    1,
				Title: "CPU Usage",
				Type:  "stat",
				GridPos: GridPos{H: 8, W: 12, X: 0, Y: 0},
				Targets: []Target{
					{
						Expr:         fmt.Sprintf(`gproc_cpu_usage{process_id="%s"}`, processID),
						RefID:        "A",
						LegendFormat: "CPU %",
					},
				},
			},
			{
				ID:    2,
				Title: "Memory Usage",
				Type:  "stat",
				GridPos: GridPos{H: 8, W: 12, X: 12, Y: 0},
				Targets: []Target{
					{
						Expr:         fmt.Sprintf(`gproc_memory_usage{process_id="%s"}`, processID),
						RefID:        "B",
						LegendFormat: "Memory MB",
					},
				},
			},
			{
				ID:    3,
				Title: "Process Status",
				Type:  "timeseries",
				GridPos: GridPos{H: 8, W: 24, X: 0, Y: 8},
				Targets: []Target{
					{
						Expr:         fmt.Sprintf(`gproc_process_status{process_id="%s"}`, processID),
						RefID:        "C",
						LegendFormat: "Status",
					},
				},
			},
		},
	}
	
	if err := g.deployDashboard(ctx, dashboard); err != nil {
		return nil, err
	}
	
	g.dashboards[dashboard.UID] = dashboard
	return dashboard, nil
}

func (g *GrafanaManager) CreateClusterDashboard(ctx context.Context) (*Dashboard, error) {
	dashboard := &Dashboard{
		UID:     "gproc-cluster-overview",
		Title:   "GProc Cluster Overview",
		Tags:    []string{"gproc", "cluster", "overview"},
		Refresh: "10s",
		Time: TimeRange{
			From: "now-6h",
			To:   "now",
		},
		Templating: Templating{
			List: []TemplateVariable{
				{
					Name:  "node",
					Type:  "query",
					Query: "label_values(gproc_node_info, node_id)",
				},
			},
		},
		Panels: []Panel{
			{
				ID:    1,
				Title: "Total Processes",
				Type:  "stat",
				GridPos: GridPos{H: 4, W: 6, X: 0, Y: 0},
				Targets: []Target{
					{
						Expr:  `sum(gproc_processes_total)`,
						RefID: "A",
					},
				},
			},
			{
				ID:    2,
				Title: "Running Processes",
				Type:  "stat",
				GridPos: GridPos{H: 4, W: 6, X: 6, Y: 0},
				Targets: []Target{
					{
						Expr:  `sum(gproc_processes_running)`,
						RefID: "B",
					},
				},
			},
			{
				ID:    3,
				Title: "Failed Processes",
				Type:  "stat",
				GridPos: GridPos{H: 4, W: 6, X: 12, Y: 0},
				Targets: []Target{
					{
						Expr:  `sum(gproc_processes_failed)`,
						RefID: "C",
					},
				},
			},
			{
				ID:    4,
				Title: "Cluster Nodes",
				Type:  "stat",
				GridPos: GridPos{H: 4, W: 6, X: 18, Y: 0},
				Targets: []Target{
					{
						Expr:  `count(gproc_node_info)`,
						RefID: "D",
					},
				},
			},
		},
	}
	
	if err := g.deployDashboard(ctx, dashboard); err != nil {
		return nil, err
	}
	
	g.dashboards[dashboard.UID] = dashboard
	return dashboard, nil
}

func (g *GrafanaManager) deployDashboard(ctx context.Context, dashboard *Dashboard) error {
	fmt.Printf("Deploying Grafana dashboard: %s\n", dashboard.Title)
	
	// Simulate Grafana API call
	dashboardJSON, _ := json.MarshalIndent(dashboard, "", "  ")
	fmt.Printf("Dashboard JSON: %s\n", string(dashboardJSON))
	
	// In real implementation, POST to /api/dashboards/db
	dashboard.ID = int(time.Now().Unix())
	dashboard.Version = 1
	
	return nil
}

func (g *GrafanaManager) UpdateDashboard(ctx context.Context, uid string, updates map[string]interface{}) error {
	dashboard, exists := g.dashboards[uid]
	if !exists {
		return fmt.Errorf("dashboard not found: %s", uid)
	}
	
	// Apply updates
	if title, ok := updates["title"].(string); ok {
		dashboard.Title = title
	}
	if refresh, ok := updates["refresh"].(string); ok {
		dashboard.Refresh = refresh
	}
	
	dashboard.Version++
	return g.deployDashboard(ctx, dashboard)
}

func (g *GrafanaManager) DeleteDashboard(ctx context.Context, uid string) error {
	fmt.Printf("Deleting Grafana dashboard: %s\n", uid)
	
	// Simulate Grafana API call
	// In real implementation, DELETE /api/dashboards/uid/{uid}
	
	delete(g.dashboards, uid)
	return nil
}

func (g *GrafanaManager) GetDashboardURL(uid string) string {
	return fmt.Sprintf("%s/d/%s", g.config.Endpoint, uid)
}

func (g *GrafanaManager) ListDashboards() []*Dashboard {
	dashboards := make([]*Dashboard, 0, len(g.dashboards))
	for _, dashboard := range g.dashboards {
		dashboards = append(dashboards, dashboard)
	}
	return dashboards
}