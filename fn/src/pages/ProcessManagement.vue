<template>
  <div class="process-management">
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">Process Management</h1>
        <p class="page-subtitle">Monitor and control your application processes</p>
      </div>
      <n-button type="primary" size="large" @click="showAddModal = true">
        <template #icon><span>➕</span></template>
        Add Process
      </n-button>
    </div>

    <div class="stats-overview">
      <div class="stat-card running">
        <div class="stat-icon">🟢</div>
        <div class="stat-info">
          <div class="stat-value">{{ runningCount }}</div>
          <div class="stat-label">Running</div>
        </div>
      </div>
      <div class="stat-card stopped">
        <div class="stat-icon">⚫</div>
        <div class="stat-info">
          <div class="stat-value">{{ stoppedCount }}</div>
          <div class="stat-label">Stopped</div>
        </div>
      </div>
      <div class="stat-card failed">
        <div class="stat-icon">🔴</div>
        <div class="stat-info">
          <div class="stat-value">{{ failedCount }}</div>
          <div class="stat-label">Failed</div>
        </div>
      </div>
      <div class="stat-card total">
        <div class="stat-icon">📊</div>
        <div class="stat-info">
          <div class="stat-value">{{ processes.length }}</div>
          <div class="stat-label">Total</div>
        </div>
      </div>
    </div>

    <div class="process-table-container">
      <n-data-table
        :columns="columns"
        :data="processes"
        :pagination="{ pageSize: 10 }"
        :bordered="false"
        striped
      />
    </div>

    <n-modal v-model:show="showAddModal" preset="card" title="Add New Process" style="width: 600px">
      <AddProcessForm @submit="handleAddProcess" @cancel="showAddModal = false" />
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, h } from 'vue'
import { NButton, NDataTable, NModal, NTag, NSpace } from 'naive-ui'
import AddProcessForm from '../components/process/AddProcessForm.vue'

const showAddModal = ref(false)

const processes = ref([
  {
    id: '1',
    name: 'web-server',
    command: './server',
    status: 'running',
    pid: 1234,
    cpu: '15.2%',
    memory: '128MB',
    uptime: '2h 15m',
    restarts: 0
  },
  {
    id: '2',
    name: 'api-service',
    command: './api',
    status: 'running',
    pid: 5678,
    cpu: '8.7%',
    memory: '64MB',
    uptime: '1h 45m',
    restarts: 1
  },
  {
    id: '3',
    name: 'worker',
    command: './worker',
    status: 'stopped',
    pid: null,
    cpu: '0%',
    memory: '0MB',
    uptime: '0m',
    restarts: 3
  },
  {
    id: '4',
    name: 'scheduler',
    command: './scheduler',
    status: 'failed',
    pid: null,
    cpu: '0%',
    memory: '0MB',
    uptime: '0m',
    restarts: 5
  }
])

const runningCount = computed(() => processes.value.filter(p => p.status === 'running').length)
const stoppedCount = computed(() => processes.value.filter(p => p.status === 'stopped').length)
const failedCount = computed(() => processes.value.filter(p => p.status === 'failed').length)

const columns = [
  {
    title: 'Process',
    key: 'name',
    render: (row: any) => h('div', { class: 'process-name' }, [
      h('div', { class: 'name' }, row.name),
      h('div', { class: 'command' }, row.command)
    ])
  },
  {
    title: 'Status',
    key: 'status',
    render: (row: any) => h(NTag, {
      type: row.status === 'running' ? 'success' : row.status === 'failed' ? 'error' : 'default',
      size: 'small'
    }, { default: () => row.status.toUpperCase() })
  },
  {
    title: 'PID',
    key: 'pid',
    render: (row: any) => row.pid || '-'
  },
  {
    title: 'CPU',
    key: 'cpu'
  },
  {
    title: 'Memory',
    key: 'memory'
  },
  {
    title: 'Uptime',
    key: 'uptime'
  },
  {
    title: 'Restarts',
    key: 'restarts'
  },
  {
    title: 'Actions',
    key: 'actions',
    render: (row: any) => h(NSpace, { size: 'small' }, {
      default: () => [
        h(NButton, {
          size: 'small',
          type: row.status === 'running' ? 'warning' : 'primary',
          onClick: () => toggleProcess(row)
        }, { default: () => row.status === 'running' ? 'Stop' : 'Start' }),
        h(NButton, {
          size: 'small',
          onClick: () => restartProcess(row)
        }, { default: () => 'Restart' }),
        h(NButton, {
          size: 'small',
          type: 'error',
          onClick: () => deleteProcess(row)
        }, { default: () => 'Delete' })
      ]
    })
  }
]

const toggleProcess = (process: any) => {
  console.log('Toggle process:', process.name)
}

const restartProcess = (process: any) => {
  console.log('Restart process:', process.name)
}

const deleteProcess = (process: any) => {
  console.log('Delete process:', process.name)
}

const handleAddProcess = (data: any) => {
  console.log('Add process:', data)
  showAddModal.value = false
}
</script>

<style scoped>
.process-management {
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

.header-content {
  flex: 1;
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

.stats-overview {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 24px;
  margin-bottom: 32px;
}

.stat-card {
  background: var(--n-color-embedded);
  border: 1px solid var(--n-border-color);
  border-radius: 16px;
  padding: 24px;
  display: flex;
  align-items: center;
  gap: 16px;
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
}

.stat-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, var(--accent-color), var(--accent-color-light));
}

.stat-card.running {
  --accent-color: #10b981;
  --accent-color-light: #34d399;
}

.stat-card.stopped {
  --accent-color: #6b7280;
  --accent-color-light: #9ca3af;
}

.stat-card.failed {
  --accent-color: #ef4444;
  --accent-color-light: #f87171;
}

.stat-card.total {
  --accent-color: #3b82f6;
  --accent-color-light: #60a5fa;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.1);
}

.stat-icon {
  font-size: 2rem;
  opacity: 0.8;
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 2rem;
  font-weight: 700;
  color: var(--n-text-color);
  line-height: 1;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 0.875rem;
  color: var(--n-text-color-2);
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.process-table-container {
  background: var(--n-color-embedded);
  border: 1px solid var(--n-border-color);
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.process-name .name {
  font-weight: 600;
  color: var(--n-text-color);
  margin-bottom: 2px;
}

.process-name .command {
  font-size: 0.875rem;
  color: var(--n-text-color-2);
  font-family: 'Consolas', 'Monaco', monospace;
}

@media (max-width: 768px) {
  .process-management {
    padding: 16px;
  }
  
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }
  
  .page-title {
    font-size: 2rem;
  }
  
  .stats-overview {
    grid-template-columns: repeat(2, 1fr);
    gap: 16px;
  }
}
</style>