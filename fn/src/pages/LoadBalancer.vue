<template>
  <div class="load-balancer">
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">Load Balancer</h1>
        <p class="page-subtitle">Manage process clusters and load distribution</p>
      </div>
      <n-button type="primary" size="large" @click="showCreateModal = true">
        <template #icon><span>🚀</span></template>
        Create Cluster
      </n-button>
    </div>

    <div class="clusters-grid">
      <div v-for="cluster in clusters" :key="cluster.name" class="cluster-card">
        <div class="cluster-header">
          <div class="cluster-info">
            <h3 class="cluster-name">{{ cluster.name }}</h3>
            <n-tag :type="getStatusType(cluster.status)" size="small">{{ cluster.status }}</n-tag>
          </div>
          <div class="cluster-actions">
            <n-button size="small" @click="scaleCluster(cluster)">Scale</n-button>
            <n-button size="small" type="warning" @click="stopCluster(cluster)">Stop</n-button>
          </div>
        </div>
        
        <div class="cluster-stats">
          <div class="stat-item">
            <span class="stat-label">Instances</span>
            <span class="stat-value">{{ cluster.instances }}/{{ cluster.maxInstances }}</span>
          </div>
          <div class="stat-item">
            <span class="stat-label">Load</span>
            <span class="stat-value">{{ cluster.load }}%</span>
          </div>
          <div class="stat-item">
            <span class="stat-label">Requests/min</span>
            <span class="stat-value">{{ cluster.requestsPerMin }}</span>
          </div>
        </div>

        <div class="instances-list">
          <div v-for="instance in cluster.instanceList" :key="instance.id" class="instance-item">
            <div class="instance-status" :class="instance.status"></div>
            <span class="instance-name">{{ instance.name }}</span>
            <span class="instance-port">:{{ instance.port }}</span>
            <span class="instance-load">{{ instance.load }}%</span>
          </div>
        </div>
      </div>
    </div>

    <n-modal v-model:show="showCreateModal" preset="card" title="Create Cluster" style="width: 600px">
      <n-form :model="clusterForm" label-placement="top">
        <n-form-item label="Cluster Name">
          <n-input v-model:value="clusterForm.name" placeholder="webapp-cluster" />
        </n-form-item>
        <n-form-item label="Command">
          <n-input v-model:value="clusterForm.command" placeholder="./server.exe" />
        </n-form-item>
        <n-form-item label="Instances">
          <n-input-number v-model:value="clusterForm.instances" :min="1" :max="10" />
        </n-form-item>
        <n-form-item label="Base Port">
          <n-input-number v-model:value="clusterForm.basePort" :min="3000" :max="9999" />
        </n-form-item>
        <div style="display: flex; justify-content: flex-end; gap: 12px; margin-top: 24px;">
          <n-button @click="showCreateModal = false">Cancel</n-button>
          <n-button type="primary" @click="createCluster">Create Cluster</n-button>
        </div>
      </n-form>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { NButton, NTag, NModal, NForm, NFormItem, NInput, NInputNumber } from 'naive-ui'

const showCreateModal = ref(false)
const clusterForm = ref({
  name: '',
  command: '',
  instances: 3,
  basePort: 8080
})

const clusters = ref([
  {
    name: 'webapp-cluster',
    status: 'running',
    instances: 4,
    maxInstances: 6,
    load: 65,
    requestsPerMin: 1250,
    instanceList: [
      { id: '1', name: 'webapp-1', port: 8080, status: 'healthy', load: 70 },
      { id: '2', name: 'webapp-2', port: 8081, status: 'healthy', load: 60 },
      { id: '3', name: 'webapp-3', port: 8082, status: 'healthy', load: 65 },
      { id: '4', name: 'webapp-4', port: 8083, status: 'warning', load: 85 }
    ]
  },
  {
    name: 'api-cluster',
    status: 'running',
    instances: 2,
    maxInstances: 4,
    load: 45,
    requestsPerMin: 800,
    instanceList: [
      { id: '1', name: 'api-1', port: 9000, status: 'healthy', load: 40 },
      { id: '2', name: 'api-2', port: 9001, status: 'healthy', load: 50 }
    ]
  }
])

const getStatusType = (status: string) => {
  switch (status) {
    case 'running': return 'success'
    case 'stopped': return 'default'
    case 'error': return 'error'
    default: return 'warning'
  }
}

const createCluster = () => {
  console.log('Creating cluster:', clusterForm.value)
  showCreateModal.value = false
}

const scaleCluster = (cluster: any) => {
  console.log('Scaling cluster:', cluster.name)
}

const stopCluster = (cluster: any) => {
  console.log('Stopping cluster:', cluster.name)
}
</script>

<style scoped>
.load-balancer {
  padding: 32px;
  max-width: 1400px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  margin-bottom: 32px;
  padding-bottom: 24px;
  border-bottom: 1px solid var(--n-border-color);
}

.page-title {
  font-size: 2.5rem;
  font-weight: 700;
  color: var(--n-text-color);
  margin: 0 0 8px 0;
  letter-spacing: -0.025em;
}

.page-subtitle {
  font-size: 1.125rem;
  color: var(--n-text-color-2);
  margin: 0;
  font-weight: 400;
}

.clusters-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 24px;
}

.cluster-card {
  background: var(--n-color-embedded);
  border: 1px solid var(--n-border-color);
  border-radius: 16px;
  padding: 24px;
  transition: all 0.3s ease;
}

.cluster-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-lg);
}

.cluster-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.cluster-name {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--n-text-color);
  margin: 0 0 4px 0;
}

.cluster-actions {
  display: flex;
  gap: 8px;
}

.cluster-stats {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  margin-bottom: 20px;
  padding: 16px;
  background: var(--n-color);
  border-radius: 8px;
}

.stat-item {
  text-align: center;
}

.stat-label {
  display: block;
  font-size: 0.75rem;
  color: var(--n-text-color-3);
  margin-bottom: 4px;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.stat-value {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--n-text-color);
}

.instances-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.instance-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 12px;
  background: var(--n-color);
  border-radius: 6px;
  font-size: 0.875rem;
}

.instance-status {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}

.instance-status.healthy {
  background: var(--gproc-success);
}

.instance-status.warning {
  background: var(--gproc-warning);
}

.instance-status.error {
  background: var(--gproc-error);
}

.instance-name {
  font-weight: 500;
  color: var(--n-text-color);
}

.instance-port {
  color: var(--n-text-color-2);
  font-family: monospace;
}

.instance-load {
  margin-left: auto;
  color: var(--n-text-color-2);
  font-weight: 500;
}

@media (max-width: 768px) {
  .load-balancer {
    padding: 16px;
  }
  
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }
  
  .clusters-grid {
    grid-template-columns: 1fr;
  }
  
  .cluster-stats {
    grid-template-columns: 1fr;
    gap: 12px;
  }
}
</style>