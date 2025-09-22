<template>
  <div class="audit-logs">
    <n-card title="Audit Trail">
      <template #header-extra>
        <n-space>
          <n-date-picker v-model:value="dateRange" type="daterange" />
          <n-button @click="refreshLogs">
            <template #icon><RefreshCw /></template>
            Refresh
          </n-button>
        </n-space>
      </template>
      
      <n-data-table
        :columns="columns"
        :data="auditEvents"
        :loading="loading"
        :pagination="{ pageSize: 20 }"
      />
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import { useSecurityStore } from '../stores/security'
import { RefreshCw } from 'lucide-vue-next'

const securityStore = useSecurityStore()

const loading = ref(false)
const dateRange = ref<[number, number] | null>(null)

const { auditEvents } = securityStore

const columns = [
  {
    title: 'Timestamp',
    key: 'timestamp',
    render: (row: any) => new Date(row.timestamp).toLocaleString()
  },
  { title: 'User', key: 'username' },
  { title: 'Action', key: 'action' },
  { title: 'Resource', key: 'resource' },
  { title: 'Target', key: 'target' },
  {
    title: 'Result',
    key: 'result',
    render: (row: any) => h('n-tag', {
      type: row.result === 'success' ? 'success' : 'error'
    }, { default: () => row.result })
  },
  { title: 'IP', key: 'ip' },
  { title: 'Details', key: 'details' }
]

const refreshLogs = () => {
  securityStore.fetchAuditEvents()
}

onMounted(() => {
  refreshLogs()
})
</script>

<style scoped>
.audit-logs {
  padding: 1rem;
}
</style>