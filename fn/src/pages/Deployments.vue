<template>
  <div class="deployments">
    <n-card title="Zero-Downtime Deployments">
      <template #header-extra>
        <n-button type="primary" @click="showCreateModal = true">
          <template #icon><Rocket /></template>
          New Deployment
        </n-button>
      </template>
      
      <!-- Deployment History -->
      <n-data-table
        :columns="columns"
        :data="deployments"
        :loading="loading"
        :pagination="{ pageSize: 10 }"
      />
    </n-card>
    
    <!-- Create Deployment Modal -->
    <n-modal v-model:show="showCreateModal" title="Create Deployment">
      <n-card style="width: 600px">
        <n-form :model="newDeployment" :rules="deploymentRules">
          <n-form-item path="processName" label="Process Name">
            <n-select
              v-model:value="newDeployment.processName"
              :options="processOptions"
              placeholder="Select process"
            />
          </n-form-item>
          
          <n-form-item path="strategy" label="Deployment Strategy">
            <n-select
              v-model:value="newDeployment.strategy"
              :options="strategyOptions"
              placeholder="Select strategy"
            />
          </n-form-item>
          
          <n-form-item path="version" label="Version">
            <n-input v-model:value="newDeployment.version" placeholder="v1.2.3" />
          </n-form-item>
          
          <n-form-item label="Options">
            <n-space vertical>
              <n-checkbox v-model:checked="newDeployment.rollbackOnFail">
                Rollback on failure
              </n-checkbox>
              <n-checkbox v-model:checked="newDeployment.healthChecks">
                Enable health checks
              </n-checkbox>
            </n-space>
          </n-form-item>
          
          <!-- Strategy-specific options -->
          <div v-if="newDeployment.strategy === 'canary'">
            <n-form-item label="Canary Percentage">
              <n-slider v-model:value="canaryPercentage" :min="5" :max="50" :step="5" />
              <span>{{ canaryPercentage }}%</span>
            </n-form-item>
          </div>
          
          <div v-if="newDeployment.strategy === 'rolling'">
            <n-form-item label="Batch Size">
              <n-input-number v-model:value="batchSize" :min="1" :max="10" />
            </n-form-item>
          </div>
        </n-form>
        
        <template #footer>
          <n-space justify="end">
            <n-button @click="showCreateModal = false">Cancel</n-button>
            <n-button type="primary" @click="createDeployment">Deploy</n-button>
          </n-space>
        </template>
      </n-card>
    </n-modal>
    
    <!-- Deployment Details Modal -->
    <n-modal v-model:show="showDetailsModal" title="Deployment Details">
      <n-card style="width: 800px" v-if="selectedDeployment">
        <n-descriptions :column="2" bordered>
          <n-descriptions-item label="Process">{{ selectedDeployment.processName }}</n-descriptions-item>
          <n-descriptions-item label="Strategy">{{ selectedDeployment.strategy }}</n-descriptions-item>
          <n-descriptions-item label="Version">{{ selectedDeployment.version }}</n-descriptions-item>
          <n-descriptions-item label="Status">
            <n-tag :type="getStatusType(selectedDeployment.status)">
              {{ selectedDeployment.status }}
            </n-tag>
          </n-descriptions-item>
          <n-descriptions-item label="Progress">
            <n-progress :percentage="selectedDeployment.progress" />
          </n-descriptions-item>
          <n-descriptions-item label="Start Time">
            {{ new Date(selectedDeployment.startTime).toLocaleString() }}
          </n-descriptions-item>
        </n-descriptions>
        
        <!-- Deployment Logs -->
        <div class="deployment-logs">
          <h3>Deployment Logs</h3>
          <n-code
            :code="selectedDeployment.logs.join('\n')"
            language="text"
            show-line-numbers
          />
        </div>
        
        <template #footer>
          <n-space justify="end">
            <n-button
              v-if="selectedDeployment.status === 'in-progress'"
              type="error"
              @click="rollbackDeployment(selectedDeployment.id)"
            >
              Rollback
            </n-button>
            <n-button @click="showDetailsModal = false">Close</n-button>
          </n-space>
        </template>
      </n-card>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, h } from 'vue'
import { useDeploymentStore } from '../stores/deployment'
import { useProcessStore } from '../stores/processes'
import { useMessage } from 'naive-ui'
import { Rocket, Eye, RotateCcw } from 'lucide-vue-next'

const deploymentStore = useDeploymentStore()
const processStore = useProcessStore()
const message = useMessage()

const loading = ref(false)
const showCreateModal = ref(false)
const showDetailsModal = ref(false)
const selectedDeployment = ref(null)
const canaryPercentage = ref(10)
const batchSize = ref(2)

const newDeployment = ref({
  processName: '',
  strategy: '',
  version: '',
  rollbackOnFail: true,
  healthChecks: true
})

const { deployments } = deploymentStore
const { processes } = processStore

const processOptions = computed(() =>
  processes.value.map(p => ({ label: p.name, value: p.name }))
)

const strategyOptions = [
  { label: 'Blue-Green', value: 'blue-green' },
  { label: 'Rolling Update', value: 'rolling' },
  { label: 'Canary', value: 'canary' }
]

const deploymentRules = {
  processName: { required: true, message: 'Process name is required' },
  strategy: { required: true, message: 'Strategy is required' },
  version: { required: true, message: 'Version is required' }
}

const columns = [
  { title: 'Process', key: 'processName' },
  { title: 'Strategy', key: 'strategy' },
  { title: 'Version', key: 'version' },
  {
    title: 'Status',
    key: 'status',
    render: (row: any) => h('n-tag', {
      type: getStatusType(row.status)
    }, { default: () => row.status })
  },
  {
    title: 'Progress',
    key: 'progress',
    render: (row: any) => h('n-progress', { percentage: row.progress })
  },
  {
    title: 'Start Time',
    key: 'startTime',
    render: (row: any) => new Date(row.startTime).toLocaleString()
  },
  {
    title: 'Actions',
    key: 'actions',
    render: (row: any) => h('n-space', [
      h('n-button', {
        size: 'small',
        onClick: () => viewDeployment(row)
      }, { default: () => 'View', icon: () => h(Eye) }),
      h('n-button', {
        size: 'small',
        type: 'error',
        disabled: row.status !== 'in-progress',
        onClick: () => rollbackDeployment(row.id)
      }, { default: () => 'Rollback', icon: () => h(RotateCcw) })
    ])
  }
]

const getStatusType = (status: string) => {
  switch (status) {
    case 'completed': return 'success'
    case 'failed': return 'error'
    case 'in-progress': return 'info'
    case 'rolled-back': return 'warning'
    default: return 'default'
  }
}

const createDeployment = async () => {
  const result = await deploymentStore.createDeployment({
    ...newDeployment.value,
    canaryPercentage: canaryPercentage.value,
    batchSize: batchSize.value
  })
  
  if (result.success) {
    message.success('Deployment started successfully')
    showCreateModal.value = false
    resetForm()
  } else {
    message.error(result.error)
  }
}

const viewDeployment = (deployment: any) => {
  selectedDeployment.value = deployment
  showDetailsModal.value = true
}

const rollbackDeployment = async (deploymentId: string) => {
  const result = await deploymentStore.rollbackDeployment(deploymentId)
  if (result.success) {
    message.success('Rollback initiated')
  } else {
    message.error(result.error)
  }
}

const resetForm = () => {
  newDeployment.value = {
    processName: '',
    strategy: '',
    version: '',
    rollbackOnFail: true,
    healthChecks: true
  }
  canaryPercentage.value = 10
  batchSize.value = 2
}

onMounted(() => {
  deploymentStore.fetchDeployments()
  processStore.fetchProcesses()
})
</script>

<style scoped>
.deployments {
  padding: 1rem;
}

.deployment-logs {
  margin-top: 2rem;
}

.deployment-logs h3 {
  margin-bottom: 1rem;
}
</style>