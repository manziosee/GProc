<template>
  <div class="cluster-management">
    <n-card title="Cluster Management">
      <template #header-extra>
        <n-button @click="refreshData" :loading="loading">
          <template #icon><RefreshCw /></template>
          Refresh
        </n-button>
      </template>
      
      <!-- Cluster Status Overview -->
      <n-grid :cols="4" :x-gap="16" class="mb-4">
        <n-grid-item>
          <n-statistic label="Leader Node" :value="status?.leader || 'Unknown'" />
        </n-grid-item>
        <n-grid-item>
          <n-statistic label="Total Nodes" :value="status?.nodes || 0" />
        </n-grid-item>
        <n-grid-item>
          <n-statistic label="Health Status">
            <template #default>
              <n-tag :type="status?.healthy ? 'success' : 'error'">
                {{ status?.healthy ? 'Healthy' : 'Unhealthy' }}
              </n-tag>
            </template>
          </n-statistic>
        </n-grid-item>
        <n-grid-item>
          <n-statistic label="Raft Term" :value="status?.raftTerm || 0" />
        </n-grid-item>
      </n-grid>
      
      <!-- Nodes Table -->
      <n-data-table
        :columns="columns"
        :data="nodes"
        :loading="loading"
        :pagination="{ pageSize: 10 }"
      />
    </n-card>
    
    <!-- Add Node Modal -->
    <n-modal v-model:show="showAddModal" title="Add Cluster Node">
      <n-card style="width: 500px">
        <n-form :model="newNode" :rules="nodeRules">
          <n-form-item path="id" label="Node ID">
            <n-input v-model:value="newNode.id" placeholder="node-3" />
          </n-form-item>
          <n-form-item path="address" label="Address">
            <n-input v-model:value="newNode.address" placeholder="10.0.1.12:9090" />
          </n-form-item>
          <n-form-item path="role" label="Role">
            <n-select v-model:value="newNode.role" :options="roleOptions" />
          </n-form-item>
        </n-form>
        
        <template #footer>
          <n-space justify="end">
            <n-button @click="showAddModal = false">Cancel</n-button>
            <n-button type="primary" @click="addNode">Add Node</n-button>
          </n-space>
        </template>
      </n-card>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import { useClusterStore } from '../stores/cluster'
import { useMessage } from 'naive-ui'
import { RefreshCw, Crown, Users, Trash2 } from 'lucide-vue-next'

const clusterStore = useClusterStore()
const message = useMessage()

const loading = ref(false)
const showAddModal = ref(false)
const newNode = ref({ id: '', address: '', role: 'follower' })

const { nodes, status } = clusterStore

const nodeRules = {
  id: { required: true, message: 'Node ID is required' },
  address: { required: true, message: 'Address is required' },
  role: { required: true, message: 'Role is required' }
}

const roleOptions = [
  { label: 'Leader', value: 'leader' },
  { label: 'Follower', value: 'follower' },
  { label: 'Candidate', value: 'candidate' }
]

const columns = [
  { title: 'Node ID', key: 'id' },
  { title: 'Address', key: 'address' },
  {
    title: 'Role',
    key: 'role',
    render: (row: any) => h('n-tag', {
      type: row.role === 'leader' ? 'warning' : 'default'
    }, { default: () => row.role })
  },
  {
    title: 'Status',
    key: 'status',
    render: (row: any) => h('n-tag', {
      type: row.status === 'active' ? 'success' : 'error'
    }, { default: () => row.status })
  },
  { title: 'Processes', key: 'processes' },
  { title: 'CPU %', key: 'cpu', render: (row: any) => `${row.cpu}%` },
  { title: 'Memory %', key: 'memory', render: (row: any) => `${row.memory}%` },
  { title: 'Uptime', key: 'uptime' },
  {
    title: 'Actions',
    key: 'actions',
    render: (row: any) => h('n-space', [
      h('n-button', {
        size: 'small',
        type: 'warning',
        disabled: row.role === 'leader',
        onClick: () => promoteNode(row.id)
      }, { default: () => 'Promote', icon: () => h(Crown) }),
      h('n-button', {
        size: 'small',
        type: 'error',
        disabled: row.role === 'leader',
        onClick: () => removeNode(row.id)
      }, { default: () => 'Remove', icon: () => h(Trash2) })
    ])
  }
]

const refreshData = async () => {
  loading.value = true
  await Promise.all([
    clusterStore.fetchNodes(),
    clusterStore.fetchStatus()
  ])
  loading.value = false
}

const promoteNode = async (nodeId: string) => {
  const result = await clusterStore.promoteNode(nodeId)
  if (result.success) {
    message.success('Node promoted successfully')
  } else {
    message.error(result.error)
  }
}

const removeNode = async (nodeId: string) => {
  const result = await clusterStore.removeNode(nodeId)
  if (result.success) {
    message.success('Node removed successfully')
  } else {
    message.error(result.error)
  }
}

const addNode = async () => {
  // Implementation would call cluster API to add node
  message.info('Add node functionality would be implemented')
  showAddModal.value = false
}

onMounted(() => {
  refreshData()
})
</script>

<style scoped>
.cluster-management {
  padding: 1rem;
}

.mb-4 {
  margin-bottom: 1rem;
}
</style>