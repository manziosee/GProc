package types

import (
	"os/exec"
	"time"
)

type ProcessStatus string

const (
	StatusRunning ProcessStatus = "running"
	StatusStopped ProcessStatus = "stopped"
	StatusFailed  ProcessStatus = "failed"
)

type Process struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	Command       string            `json:"command"`
	Args          []string          `json:"args"`
	WorkingDir    string            `json:"working_dir"`
	Env           map[string]string `json:"env"`
	Status        ProcessStatus     `json:"status"`
	PID           int               `json:"pid"`
	StartTime     time.Time         `json:"start_time"`
	Restarts      int               `json:"restarts"`
	AutoRestart   bool              `json:"auto_restart"`
	MaxRestarts   int               `json:"max_restarts"`
	LogFile       string            `json:"log_file"`
	Group         string            `json:"group"`
	HealthCheck   *HealthCheck      `json:"health_check"`
	LogRotation   *LogRotation      `json:"log_rotation"`
	ResourceLimit *ResourceLimit    `json:"resource_limit"`
	Notifications *Notifications    `json:"notifications"`
	Cmd           *exec.Cmd         `json:"-"`
}

type HealthCheck struct {
	URL      string        `json:"url"`
	Interval time.Duration `json:"interval"`
	Timeout  time.Duration `json:"timeout"`
	Retries  int           `json:"retries"`
}

type LogRotation struct {
	MaxSize  string `json:"max_size"`
	MaxFiles int    `json:"max_files"`
}

type ResourceLimit struct {
	MemoryMB int     `json:"memory_mb"`
	CPULimit float64 `json:"cpu_limit"`
}

type Notifications struct {
	Email string `json:"email"`
	Slack string `json:"slack"`
}

type ProcessGroup struct {
	Name      string   `json:"name"`
	Processes []string `json:"processes"`
}

type ProcessTemplate struct {
	Name        string            `json:"name"`
	Command     string            `json:"command"`
	Args        []string          `json:"args"`
	WorkingDir  string            `json:"working_dir"`
	Env         map[string]string `json:"env"`
	AutoRestart bool              `json:"auto_restart"`
	MaxRestarts int               `json:"max_restarts"`
}

type ScheduledTask struct {
	Name        string            `json:"name"`
	Command     string            `json:"command"`
	Args        []string          `json:"args"`
	Cron        string            `json:"cron"`
	NextRun     time.Time         `json:"next_run"`
	LastRun     time.Time         `json:"last_run,omitempty"`
	Enabled     bool              `json:"enabled"`
	Timeout     time.Duration     `json:"timeout,omitempty"`
	Retries     int               `json:"retries,omitempty"`
	Env         map[string]string `json:"env,omitempty"`
	WorkingDir  string            `json:"working_dir,omitempty"`
	Description string            `json:"description,omitempty"`
}

type Config struct {
	Processes      []Process         `json:"processes"`
	Groups         []ProcessGroup    `json:"groups"`
	Templates      []ProcessTemplate `json:"templates"`
	ScheduledTasks []ScheduledTask   `json:"scheduled_tasks"`
	LogDir         string            `json:"log_dir"`
	WebPort        int               `json:"web_port"`
	Security       *SecurityConfig   `json:"security,omitempty"`
	Cluster        *ClusterConfig    `json:"cluster,omitempty"`
	Observability  *ObservabilityConfig `json:"observability,omitempty"`
	API            *APIConfig        `json:"api,omitempty"`
	Enterprise     *EnterpriseConfig `json:"enterprise,omitempty"`
	Container      *ContainerConfig  `json:"container,omitempty"`
	Kubernetes     *KubernetesConfig `json:"kubernetes,omitempty"`
}

// Enterprise & Security Types
type SecurityConfig struct {
	RBAC         *RBACConfig    `json:"rbac,omitempty"`
	Auth         *AuthConfig    `json:"auth,omitempty"`
	TLS          *TLSConfig     `json:"tls,omitempty"`
	Encryption   *EncryptionConfig `json:"encryption,omitempty"`
	AuditLog     *AuditConfig   `json:"audit_log,omitempty"`
	Secrets      *SecretsConfig `json:"secrets,omitempty"`
}

type RBACConfig struct {
	Enabled bool   `json:"enabled"`
	Roles   []Role `json:"roles"`
	Users   []User `json:"users"`
}

type Role struct {
	Name        string       `json:"name"`
	Permissions []Permission `json:"permissions"`
}

type Permission struct {
	Resource string   `json:"resource"` // process, cluster, config
	Actions  []string `json:"actions"`  // read, write, delete, execute
	Scope    string   `json:"scope"`    // *, group:name, process:name
}

type User struct {
	ID       string    `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"-"` // Don't serialize password
	Email    string    `json:"email"`
	Roles    []string  `json:"roles"`
	Created  time.Time `json:"created"`
	LastSeen time.Time `json:"last_seen"`
	Enabled  bool      `json:"enabled,omitempty"`
}

type AuthConfig struct {
	JWT    *JWTConfig    `json:"jwt,omitempty"`
	OAuth2 *OAuth2Config `json:"oauth2,omitempty"`
	LDAP   *LDAPConfig   `json:"ldap,omitempty"`
	SSO    *SSOConfig    `json:"sso,omitempty"`
}

type JWTConfig struct {
	Secret     string        `json:"secret"`
	Expiration time.Duration `json:"expiration"`
	Issuer     string        `json:"issuer"`
}

type OAuth2Config struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURL  string `json:"redirect_url"`
	Provider     string `json:"provider"` // google, github, azure
}

type LDAPConfig struct {
	Server   string `json:"server"`
	BaseDN   string `json:"base_dn"`
	BindUser string `json:"bind_user"`
	BindPass string `json:"bind_pass"`
}

type SSOConfig struct {
	Enabled     bool   `json:"enabled"`
	Provider    string `json:"provider"`
	ClientId    string `json:"client_id"`
	EntityID    string `json:"entity_id"`
	MetadataUrl string `json:"metadata_url"`
}

type MFAConfig struct {
	Enabled     bool   `json:"enabled"`
	Issuer      string `json:"issuer"`
	WindowSize  int    `json:"window_size"`
	BackupCodes int    `json:"backup_codes"`
}

type TLSConfig struct {
	Enabled  bool   `json:"enabled"`
	CertFile string `json:"cert_file"`
	KeyFile  string `json:"key_file"`
	CAFile   string `json:"ca_file"`
}

type EncryptionConfig struct {
	Enabled    bool   `json:"enabled"`
	KeyFile    string `json:"key_file"`
	Algorithm  string `json:"algorithm"` // AES-256-GCM
}

type AuditConfig struct {
	Enabled bool   `json:"enabled"`
	LogFile string `json:"log_file"`
	Format  string `json:"format"` // json, syslog
}

type SecretsConfig struct {
	Provider string            `json:"provider"` // vault, aws-kms, azure-kv
	Config   map[string]string `json:"config"`
}

// Distributed & Multi-Node Types
type ClusterConfig struct {
	Enabled     bool         `json:"enabled"`
	NodeID      string       `json:"node_id"`
	Nodes       []ClusterNode `json:"nodes"`
	Discovery   *DiscoveryConfig `json:"discovery,omitempty"`
	Consensus   *ConsensusConfig `json:"consensus,omitempty"`
	Replication *ReplicationConfig `json:"replication,omitempty"`
}

type ClusterNode struct {
	ID       string    `json:"id"`
	Address  string    `json:"address"`
	Role     string    `json:"role"` // leader, follower, agent
	Status   string    `json:"status"` // active, inactive, failed
	LastSeen time.Time `json:"last_seen"`
	Metadata map[string]string `json:"metadata"`
}

type DiscoveryConfig struct {
	Provider string            `json:"provider"` // etcd, consul, dns
	Config   map[string]string `json:"config"`
}

type ConsensusConfig struct {
	Algorithm string `json:"algorithm"` // raft
	DataDir   string `json:"data_dir"`
}

type ReplicationConfig struct {
	Factor int `json:"factor"` // replication factor
}

// Cloud-Native Types
type ContainerConfig struct {
	Enabled    bool              `json:"enabled"`
	Runtime    string            `json:"runtime"` // docker, containerd
	Registries []ContainerRegistry `json:"registries"`
	Networks   []ContainerNetwork `json:"networks"`
}

type ContainerRegistry struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type ContainerNetwork struct {
	Name   string `json:"name"`
	Driver string `json:"driver"`
	Subnet string `json:"subnet"`
}

type KubernetesConfig struct {
	Enabled    bool   `json:"enabled"`
	Kubeconfig string `json:"kubeconfig"`
	Namespace  string `json:"namespace"`
}

// Observability Types
type ObservabilityConfig struct {
	Metrics *MetricsConfig `json:"metrics,omitempty"`
	Tracing *TracingConfig `json:"tracing,omitempty"`
	Logging *LoggingConfig `json:"logging,omitempty"`
	Alerting *AlertingConfig `json:"alerting,omitempty"`
}

type MetricsConfig struct {
	Enabled    bool   `json:"enabled"`
	Prometheus *PrometheusConfig `json:"prometheus,omitempty"`
	Grafana    *GrafanaConfig `json:"grafana,omitempty"`
}

type PrometheusConfig struct {
	Endpoint string `json:"endpoint"`
	Port     int    `json:"port"`
	Path     string `json:"path"`
}

type GrafanaConfig struct {
	Endpoint   string   `json:"endpoint"`
	APIKey     string   `json:"api_key"`
	OrgID      int      `json:"org_id"`
	Dashboards []string `json:"dashboards"`
}

type TracingConfig struct {
	Enabled      bool   `json:"enabled"`
	Provider     string `json:"provider"` // jaeger, zipkin, otel
	Endpoint     string `json:"endpoint"`
	SampleRate   float64 `json:"sample_rate"`
}

type LoggingConfig struct {
	Level       string `json:"level"`
	Format      string `json:"format"` // json, text
	Output      string `json:"output"` // stdout, file, syslog
	Aggregation *LogAggregationConfig `json:"aggregation,omitempty"`
}

type LogAggregationConfig struct {
	Provider string            `json:"provider"` // elasticsearch, splunk
	Config   map[string]string `json:"config"`
}

type AlertingConfig struct {
	Enabled   bool            `json:"enabled"`
	Providers []AlertProvider `json:"providers"`
	Rules     []AlertRule     `json:"rules"`
}

type AlertProvider struct {
	Name   string            `json:"name"`
	Type   string            `json:"type"` // slack, email, pagerduty, webhook
	Config map[string]string `json:"config"`
}

type AlertRule struct {
	Name        string  `json:"name"`
	Condition   string  `json:"condition"`
	Threshold   float64 `json:"threshold"`
	Duration    time.Duration `json:"duration"`
	Severity    string  `json:"severity"`
	Providers   []string `json:"providers"`
}

// API Types
type APIConfig struct {
	REST      *RESTConfig      `json:"rest,omitempty"`
	GRPC      *GRPCConfig      `json:"grpc,omitempty"`
	WebSocket *WebSocketConfig `json:"websocket,omitempty"`
	Plugins   *PluginConfig    `json:"plugins,omitempty"`
}

type RESTConfig struct {
	Enabled bool   `json:"enabled"`
	Port    int    `json:"port"`
	Prefix  string `json:"prefix"`
}

type GRPCConfig struct {
	Enabled bool `json:"enabled"`
	Port    int  `json:"port"`
}

type WebSocketConfig struct {
	Enabled bool   `json:"enabled"`
	Path    string `json:"path"`
}

type PluginConfig struct {
	Enabled bool     `json:"enabled"`
	Dir     string   `json:"dir"`
	Plugins []Plugin `json:"plugins"`
}

type Plugin struct {
	Name    string   `json:"name"`
	Path    string   `json:"path"`
	Enabled bool     `json:"enabled"`
	Events  []string `json:"events,omitempty"`
}

// Enterprise Ops Types
type EnterpriseConfig struct {
	HA          *HAConfig          `json:"ha,omitempty"`
	Backup      *BackupConfig      `json:"backup,omitempty"`
	MultiTenant *MultiTenantConfig `json:"multi_tenant,omitempty"`
	Quotas      *QuotaConfig       `json:"quotas,omitempty"`
}

type HAConfig struct {
	Enabled   bool   `json:"enabled"`
	Mode      string `json:"mode"` // active-passive, active-active
	Replicas  int    `json:"replicas"`
	Failover  *FailoverConfig `json:"failover,omitempty"`
}

type FailoverConfig struct {
	Timeout    time.Duration `json:"timeout"`
	Retries    int          `json:"retries"`
	HealthCheck string      `json:"health_check"`
}

type BackupConfig struct {
	Enabled   bool          `json:"enabled"`
	Interval  time.Duration `json:"interval"`
	Retention int          `json:"retention"`
	Storage   *StorageConfig `json:"storage,omitempty"`
}

type StorageConfig struct {
	Provider string            `json:"provider"` // s3, gcs, azure
	Config   map[string]string `json:"config"`
}

type MultiTenantConfig struct {
	Enabled    bool       `json:"enabled"`
	Tenants    []Tenant   `json:"tenants"`
	Namespaces []Namespace `json:"namespaces"`
}

type Tenant struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Namespaces  []string  `json:"namespaces"`
	Quotas      *TenantQuota `json:"quotas,omitempty"`
	Created     time.Time `json:"created"`
}

type Namespace struct {
	Name     string `json:"name"`
	TenantID string `json:"tenant_id"`
	Quotas   *NamespaceQuota `json:"quotas,omitempty"`
}

type QuotaConfig struct {
	Enabled bool    `json:"enabled"`
	Global  *GlobalQuota `json:"global,omitempty"`
}

type GlobalQuota struct {
	MaxProcesses int     `json:"max_processes"`
	MaxCPU       float64 `json:"max_cpu"`
	MaxMemory    int64   `json:"max_memory"`
}

type TenantQuota struct {
	MaxProcesses int     `json:"max_processes"`
	MaxCPU       float64 `json:"max_cpu"`
	MaxMemory    int64   `json:"max_memory"`
	MaxNamespaces int    `json:"max_namespaces"`
}

type NamespaceQuota struct {
	MaxProcesses int     `json:"max_processes"`
	MaxCPU       float64 `json:"max_cpu"`
	MaxMemory    int64   `json:"max_memory"`
}

// Missing types for compilation
type Snapshot struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
	Processes []Process `json:"processes"`
	Config    *Config   `json:"config"`
}

type BlueGreenConfig struct {
	Enabled   bool   `json:"enabled"`
	BluePort  int    `json:"blue_port"`
	GreenPort int    `json:"green_port"`
	Active    string `json:"active"` // blue or green
}

type Alert struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Message      string    `json:"message"`
	Severity     string    `json:"severity"`
	Type         string    `json:"type"`
	Timestamp    time.Time `json:"timestamp"`
	ProcessID    string    `json:"process_id"`
	Resolved     bool      `json:"resolved"`
	Acknowledged bool      `json:"acknowledged"`
}

type ProcessMetrics struct {
	CPUUsage    float64       `json:"cpu_usage"`
	MemoryUsage int64         `json:"memory_usage"`
	Uptime      time.Duration `json:"uptime"`
	Restarts    int           `json:"restarts"`
}

type MetricPoint struct {
	Timestamp time.Time `json:"timestamp"`
	CPU       float64   `json:"cpu"`
	Memory    float64   `json:"memory"`
}

// Language-specific probe results
type ProbeResult struct {
	Name      string      `json:"name"`
	Value     interface{} `json:"value"`
	Unit      string      `json:"unit"`
	Timestamp time.Time   `json:"timestamp"`
	Healthy   bool        `json:"healthy"`
}

// Deployment types
type DeploymentStrategy string

const (
	BlueGreenStrategy DeploymentStrategy = "blue-green"
	RollingStrategy   DeploymentStrategy = "rolling"
	CanaryStrategy    DeploymentStrategy = "canary"
)

type DeploymentConfig struct {
	ProcessName    string             `json:"process_name"`
	Strategy       DeploymentStrategy `json:"strategy"`
	NewVersion     string             `json:"new_version"`
	Command        string             `json:"command"`
	Args           []string           `json:"args"`
	HealthCheck    *HealthCheck       `json:"health_check"`
	RollbackOnFail bool               `json:"rollback_on_fail"`
	Timeout        time.Duration      `json:"timeout"`
}

type DeploymentStatus struct {
	Strategy    string    `json:"strategy"`
	Status      string    `json:"status"` // pending, in-progress, completed, failed
	Progress    int       `json:"progress"` // 0-100
	StartTime   time.Time `json:"start_time"`
	CompletedAt time.Time `json:"completed_at"`
	Error       string    `json:"error,omitempty"`
}

// Language template types
type LanguageTemplate struct {
	Name        string            `json:"name"`
	Extensions  []string          `json:"extensions"`
	Command     string            `json:"command"`
	Args        []string          `json:"args"`
	HealthCheck string            `json:"health_check"`
	EnvVars     map[string]string `json:"env_vars"`
	Probes      []string          `json:"probes"`
}

// Enhanced Plugin type (replacing the basic one)
type PluginExtended struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Path        string            `json:"path"`
	Enabled     bool              `json:"enabled"`
	Events      []string          `json:"events,omitempty"`
	Config      map[string]string `json:"config,omitempty"`
	Description string            `json:"description,omitempty"`
}

type PluginEvent struct {
	Type      string                 `json:"type"`
	ProcessID string                 `json:"process_id"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// Notification types
type NotificationChannel struct {
	Name     string            `json:"name"`
	Type     string            `json:"type"` // email, slack, teams, discord, webhook
	Enabled  bool              `json:"enabled"`
	Config   map[string]string `json:"config"`
	Filters  []string          `json:"filters,omitempty"`
}

type NotificationRule struct {
	Name      string   `json:"name"`
	Condition string   `json:"condition"`
	Channels  []string `json:"channels"`
	Enabled   bool     `json:"enabled"`
	Cooldown  time.Duration `json:"cooldown"`
}

type ExtendedProbeConfig struct {
	Enabled   bool     `json:"enabled"`
	Languages []string `json:"languages"`
	Interval  time.Duration `json:"interval"`
}

type ComposeConfig struct {
	Enabled     bool              `json:"enabled"`
	File        string            `json:"file"`
	DefaultPath string            `json:"default_path"`
	Project     string            `json:"project"`
	Services    []string          `json:"services"`
	Env         map[string]string `json:"env"`
}

// Missing config types
type CostConfig struct {
	Enabled  bool   `json:"enabled"`
	Provider string `json:"provider"`
	Region   string `json:"region"`
}

type ResourceRequest struct {
	CPU          float64 `json:"cpu"`
	Memory       int64   `json:"memory"`
	Disk         int64   `json:"disk"`
	CPURequest   float64 `json:"cpu_request"`
	MemoryRequest int64  `json:"memory_request"`
	ProcessCount int     `json:"process_count"`
}

type DebugConfig struct {
	Enabled bool   `json:"enabled"`
	Port    int    `json:"port"`
	Path    string `json:"path"`
}

type GitOpsConfig struct {
	Enabled      bool     `json:"enabled"`
	Repository   string   `json:"repository"`
	Repositories []string `json:"repositories"`
	Branch       string   `json:"branch"`
	Path         string   `json:"path"`
}

type HealingConfig struct {
	Enabled     bool          `json:"enabled"`
	MaxRetries  int           `json:"max_retries"`
	Cooldown    time.Duration `json:"cooldown"`
	MLEnabled   bool          `json:"ml_enabled"`
}

type ProcessFailure struct {
	ProcessID string    `json:"process_id"`
	Type      string    `json:"type"`
	Reason    string    `json:"reason"`
	Timestamp time.Time `json:"timestamp"`
	Retries   int       `json:"retries"`
}

type HotReloadConfig struct {
	Enabled   bool     `json:"enabled"`
	Languages []string `json:"languages"`
	WatchDirs []string `json:"watch_dirs"`
}

type MarketplaceConfig struct {
	Enabled bool   `json:"enabled"`
	URL     string `json:"url"`
	APIKey  string `json:"api_key"`
}

type ServiceMeshConfig struct {
	Enabled   bool              `json:"enabled"`
	Provider  string            `json:"provider"`
	Namespace string            `json:"namespace"`
	Endpoint  string            `json:"endpoint"`
	Config    map[string]string `json:"config"`
}

type AnomalyConfig struct {
	Enabled     bool          `json:"enabled"`
	Threshold   float64       `json:"threshold"`
	Sensitivity float64       `json:"sensitivity"`
	Window      time.Duration `json:"window"`
}

type ProfilingConfig struct {
	Enabled bool   `json:"enabled"`
	Port    int    `json:"port"`
	Path    string `json:"path"`
}

type SDKConfig struct {
	Enabled   bool     `json:"enabled"`
	Languages []string `json:"languages"`
	OutputDir string   `json:"output_dir"`
}

type ServerlessConfig struct {
	Enabled   bool              `json:"enabled"`
	Provider  string            `json:"provider"`
	AWS       *AWSConfig        `json:"aws"`
	GCP       *GCPConfig        `json:"gcp"`
	Functions map[string]string `json:"functions"`
}

type AWSConfig struct {
	Region    string `json:"region"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Role      string `json:"role"`
}

type GCPConfig struct {
	Project     string `json:"project"`
	ProjectID   string `json:"project_id"`
	Region      string `json:"region"`
	Credentials string `json:"credentials"`
}

// Repository type for GitOps
type Repository struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Branch string `json:"branch"`
	Path   string `json:"path"`
}

type VisualConfig struct {
	Enabled bool   `json:"enabled"`
	Theme   string `json:"theme"`
	Layout  string `json:"layout"`
}

type TopologyConfig struct {
	Enabled bool `json:"enabled"`
	Layout  string `json:"layout"`
}

type WorkflowConfig struct {
	Enabled   bool              `json:"enabled"`
	Workflows map[string]string `json:"workflows"`
}

// Federation types
type FederationConfig struct {
	Enabled        bool   `json:"enabled"`
	RoutingStrategy string `json:"routing_strategy"`
	Clusters       []string `json:"clusters"`
}

type FederatedProcessConfig struct {
	ProcessID    string          `json:"process_id"`
	Replicas     int             `json:"replicas"`
	ClusterCount int             `json:"cluster_count"`
	Regions      []string        `json:"regions"`
	Resources    ResourceRequest `json:"resources"`
}