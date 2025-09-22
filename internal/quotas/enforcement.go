package quotas

import (
	"context"
	"fmt"
	"sync"

	"gproc/pkg/types"
)

type QuotaManager struct {
	config     *types.QuotaConfig
	tenants    map[string]*TenantQuota
	namespaces map[string]*NamespaceQuota
	usage      map[string]*ResourceUsage
	mutex      sync.RWMutex
}

type TenantQuota struct {
	TenantID      string  `json:"tenant_id"`
	MaxProcesses  int     `json:"max_processes"`
	MaxCPU        float64 `json:"max_cpu"`
	MaxMemory     int64   `json:"max_memory"`
	MaxNamespaces int     `json:"max_namespaces"`
}

type NamespaceQuota struct {
	Namespace    string  `json:"namespace"`
	TenantID     string  `json:"tenant_id"`
	MaxProcesses int     `json:"max_processes"`
	MaxCPU       float64 `json:"max_cpu"`
	MaxMemory    int64   `json:"max_memory"`
}

type ResourceUsage struct {
	ProcessCount int     `json:"process_count"`
	CPUUsage     float64 `json:"cpu_usage"`
	MemoryUsage  int64   `json:"memory_usage"`
}

func NewQuotaManager(config *types.QuotaConfig) *QuotaManager {
	return &QuotaManager{
		config:     config,
		tenants:    make(map[string]*TenantQuota),
		namespaces: make(map[string]*NamespaceQuota),
		usage:      make(map[string]*ResourceUsage),
	}
}

func (q *QuotaManager) SetTenantQuota(tenantID string, quota *TenantQuota) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	
	quota.TenantID = tenantID
	q.tenants[tenantID] = quota
}

func (q *QuotaManager) SetNamespaceQuota(namespace, tenantID string, quota *NamespaceQuota) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	
	quota.Namespace = namespace
	quota.TenantID = tenantID
	q.namespaces[namespace] = quota
}

func (q *QuotaManager) CheckQuota(ctx context.Context, tenantID, namespace string, request *ResourceRequest) error {
	q.mutex.RLock()
	defer q.mutex.RUnlock()
	
	// Check tenant quota
	if err := q.checkTenantQuota(tenantID, request); err != nil {
		return err
	}
	
	// Check namespace quota
	if err := q.checkNamespaceQuota(namespace, request); err != nil {
		return err
	}
	
	return nil
}

func (q *QuotaManager) checkTenantQuota(tenantID string, request *ResourceRequest) error {
	quota, exists := q.tenants[tenantID]
	if !exists {
		return nil // No quota set
	}
	
	usage := q.getTenantUsage(tenantID)
	
	// Check process count
	if usage.ProcessCount+request.ProcessCount > quota.MaxProcesses {
		return fmt.Errorf("tenant %s would exceed process quota: %d/%d", 
			tenantID, usage.ProcessCount+request.ProcessCount, quota.MaxProcesses)
	}
	
	// Check CPU
	if usage.CPUUsage+request.CPURequest > quota.MaxCPU {
		return fmt.Errorf("tenant %s would exceed CPU quota: %.2f/%.2f", 
			tenantID, usage.CPUUsage+request.CPURequest, quota.MaxCPU)
	}
	
	// Check memory
	if usage.MemoryUsage+request.MemoryRequest > quota.MaxMemory {
		return fmt.Errorf("tenant %s would exceed memory quota: %d/%d MB", 
			tenantID, (usage.MemoryUsage+request.MemoryRequest)/1024/1024, quota.MaxMemory/1024/1024)
	}
	
	return nil
}

func (q *QuotaManager) checkNamespaceQuota(namespace string, request *ResourceRequest) error {
	quota, exists := q.namespaces[namespace]
	if !exists {
		return nil // No quota set
	}
	
	usage := q.getNamespaceUsage(namespace)
	
	// Check process count
	if usage.ProcessCount+request.ProcessCount > quota.MaxProcesses {
		return fmt.Errorf("namespace %s would exceed process quota: %d/%d", 
			namespace, usage.ProcessCount+request.ProcessCount, quota.MaxProcesses)
	}
	
	// Check CPU
	if usage.CPUUsage+request.CPURequest > quota.MaxCPU {
		return fmt.Errorf("namespace %s would exceed CPU quota: %.2f/%.2f", 
			namespace, usage.CPUUsage+request.CPURequest, quota.MaxCPU)
	}
	
	// Check memory
	if usage.MemoryUsage+request.MemoryRequest > quota.MaxMemory {
		return fmt.Errorf("namespace %s would exceed memory quota: %d/%d MB", 
			namespace, (usage.MemoryUsage+request.MemoryRequest)/1024/1024, quota.MaxMemory/1024/1024)
	}
	
	return nil
}

func (q *QuotaManager) UpdateUsage(tenantID, namespace string, usage *ResourceUsage) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	
	tenantKey := "tenant:" + tenantID
	namespaceKey := "namespace:" + namespace
	
	q.usage[tenantKey] = usage
	q.usage[namespaceKey] = usage
}

func (q *QuotaManager) getTenantUsage(tenantID string) *ResourceUsage {
	key := "tenant:" + tenantID
	if usage, exists := q.usage[key]; exists {
		return usage
	}
	return &ResourceUsage{}
}

func (q *QuotaManager) getNamespaceUsage(namespace string) *ResourceUsage {
	key := "namespace:" + namespace
	if usage, exists := q.usage[key]; exists {
		return usage
	}
	return &ResourceUsage{}
}

func (q *QuotaManager) GetQuotaStatus(tenantID, namespace string) *QuotaStatus {
	q.mutex.RLock()
	defer q.mutex.RUnlock()
	
	tenantQuota := q.tenants[tenantID]
	namespaceQuota := q.namespaces[namespace]
	tenantUsage := q.getTenantUsage(tenantID)
	namespaceUsage := q.getNamespaceUsage(namespace)
	
	return &QuotaStatus{
		TenantID:       tenantID,
		Namespace:      namespace,
		TenantQuota:    tenantQuota,
		NamespaceQuota: namespaceQuota,
		TenantUsage:    tenantUsage,
		NamespaceUsage: namespaceUsage,
	}
}

type ResourceRequest struct {
	ProcessCount  int     `json:"process_count"`
	CPURequest    float64 `json:"cpu_request"`
	MemoryRequest int64   `json:"memory_request"`
}

type QuotaStatus struct {
	TenantID       string           `json:"tenant_id"`
	Namespace      string           `json:"namespace"`
	TenantQuota    *TenantQuota     `json:"tenant_quota"`
	NamespaceQuota *NamespaceQuota  `json:"namespace_quota"`
	TenantUsage    *ResourceUsage   `json:"tenant_usage"`
	NamespaceUsage *ResourceUsage   `json:"namespace_usage"`
}