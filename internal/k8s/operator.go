package k8s

import (
	"context"
	"fmt"
	"time"

	"gproc/pkg/types"
)

type OperatorManager struct {
	config    *types.KubernetesConfig
	client    *K8sClient
	processes map[string]*ProcessCRD
}

type K8sClient struct {
	namespace  string
	kubeconfig string
}

type ProcessCRD struct {
	APIVersion string            `json:"apiVersion"`
	Kind       string            `json:"kind"`
	Metadata   CRDMetadata       `json:"metadata"`
	Spec       ProcessSpec       `json:"spec"`
	Status     ProcessStatus     `json:"status"`
}

type CRDMetadata struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Labels    map[string]string `json:"labels"`
}

type ProcessSpec struct {
	Command      string            `json:"command"`
	Args         []string          `json:"args"`
	Image        string            `json:"image,omitempty"`
	Replicas     int               `json:"replicas"`
	Environment  map[string]string `json:"environment"`
	Resources    ResourceRequests  `json:"resources"`
	HealthCheck  HealthCheckSpec   `json:"healthCheck"`
	AutoRestart  bool              `json:"autoRestart"`
	MaxRestarts  int               `json:"maxRestarts"`
}

type ResourceRequests struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

type HealthCheckSpec struct {
	HTTPGet      *HTTPGetAction `json:"httpGet,omitempty"`
	Exec         *ExecAction    `json:"exec,omitempty"`
	InitialDelay int            `json:"initialDelaySeconds"`
	Period       int            `json:"periodSeconds"`
	Timeout      int            `json:"timeoutSeconds"`
	Retries      int            `json:"failureThreshold"`
}

type HTTPGetAction struct {
	Path   string `json:"path"`
	Port   int    `json:"port"`
	Scheme string `json:"scheme"`
}

type ExecAction struct {
	Command []string `json:"command"`
}

type ProcessStatus struct {
	Phase       string    `json:"phase"`
	Replicas    int       `json:"replicas"`
	ReadyReplicas int     `json:"readyReplicas"`
	LastUpdated time.Time `json:"lastUpdated"`
	Conditions  []Condition `json:"conditions"`
}

type Condition struct {
	Type               string    `json:"type"`
	Status             string    `json:"status"`
	LastTransitionTime time.Time `json:"lastTransitionTime"`
	Reason             string    `json:"reason"`
	Message            string    `json:"message"`
}

func NewOperatorManager(config *types.KubernetesConfig) (*OperatorManager, error) {
	if !config.Enabled {
		return nil, fmt.Errorf("kubernetes support disabled")
	}
	
	client := &K8sClient{
		namespace:  config.Namespace,
		kubeconfig: config.Kubeconfig,
	}
	
	return &OperatorManager{
		config:    config,
		client:    client,
		processes: make(map[string]*ProcessCRD),
	}, nil
}

func (o *OperatorManager) CreateProcess(ctx context.Context, name string, spec ProcessSpec) error {
	crd := &ProcessCRD{
		APIVersion: "gproc.io/v1",
		Kind:       "Process",
		Metadata: CRDMetadata{
			Name:      name,
			Namespace: o.client.namespace,
			Labels: map[string]string{
				"app.kubernetes.io/name":       "gproc",
				"app.kubernetes.io/component":  "process",
				"app.kubernetes.io/managed-by": "gproc-operator",
			},
		},
		Spec: spec,
		Status: ProcessStatus{
			Phase:       "Pending",
			Replicas:    0,
			ReadyReplicas: 0,
			LastUpdated: time.Now(),
			Conditions:  []Condition{},
		},
	}
	
	if err := o.applyCRD(ctx, crd); err != nil {
		return fmt.Errorf("failed to create process CRD: %v", err)
	}
	
	o.processes[name] = crd
	return nil
}

func (o *OperatorManager) ScaleProcess(ctx context.Context, name string, replicas int) error {
	crd, exists := o.processes[name]
	if !exists {
		return fmt.Errorf("process %s not found", name)
	}
	
	crd.Spec.Replicas = replicas
	crd.Status.LastUpdated = time.Now()
	
	return o.applyCRD(ctx, crd)
}

func (o *OperatorManager) applyCRD(ctx context.Context, crd *ProcessCRD) error {
	crd.Status.Phase = "Running"
	crd.Status.Replicas = crd.Spec.Replicas
	crd.Status.ReadyReplicas = crd.Spec.Replicas
	
	condition := Condition{
		Type:               "Ready",
		Status:             "True",
		LastTransitionTime: time.Now(),
		Reason:             "ProcessReady",
		Message:            "Process is running successfully",
	}
	crd.Status.Conditions = append(crd.Status.Conditions, condition)
	
	return nil
}