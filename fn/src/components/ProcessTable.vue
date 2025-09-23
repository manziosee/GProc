<template>
  <div class="process-table">
    <div class="toolbar">
      <n-input v-model:value="query" placeholder="Search processes..." clearable size="large" />
      <div class="spacer" />
      <n-button @click="$emit('refresh')" secondary>Refresh</n-button>
    </div>

    <n-data-table :columns="columns" :data="filtered" :bordered="false" :striped="true" :max-height="560" />
  </div>
</template>

<script setup lang="ts">
import { computed, ref, h } from 'vue'
import { NDataTable, NTag, NButton, NInput } from 'naive-ui'

export interface Proc {
  id: string
  name: string
  pid: number
  status: string
  restarts: number
  start_time?: string
  log_file?: string
}

const props = defineProps<{
  items: Proc[]
  onStart: (id: string) => void
  onStop: (id: string) => void
  onRestart: (id: string) => void
}>()

const query = ref('')

const columns = [
  {
    title: 'App name',
    key: 'name',
    className: 'col-name'
  },
  {
    title: 'pid',
    key: 'pid',
    width: 90
  },
  {
    title: 'status',
    key: 'status',
    width: 120,
    render: (row: Proc) => {
      const type = row.status === 'running' ? 'success' : row.status === 'failed' ? 'error' : 'warning'
      return h(NTag, { size: 'small', type: type as any, bordered: true }, () => row.status)
    }
  },
  {
    title: 'restart',
    key: 'restarts',
    width: 90
  },
  {
    title: 'uptime',
    key: 'uptime',
    width: 160,
    render: (row: Proc) => formatUptime(row.start_time)
  },
  {
    title: 'actions',
    key: 'actions',
    width: 260,
    render: (row: Proc) => {
      return h('div', { class: 'actions' }, [
        h(NButton, { size: 'small', tertiary: true, onClick: () => props.onStart(row.id) }, () => 'Start'),
        h(NButton, { size: 'small', tertiary: true, onClick: () => props.onStop(row.id) }, () => 'Stop'),
        h(NButton, { size: 'small', tertiary: true, type: 'primary', onClick: () => props.onRestart(row.id) }, () => 'Restart')
      ])
    }
  }
]

const filtered = computed(() => {
  if (!query.value) return props.items
  const q = query.value.toLowerCase()
  return props.items.filter(p => p.name.toLowerCase().includes(q) || p.id.toLowerCase().includes(q))
})

function formatUptime(start?: string) {
  if (!start) return '—'
  const t = new Date(start).getTime()
  if (!t) return '—'
  const secs = Math.max(0, Math.floor((Date.now() - t) / 1000))
  const d = Math.floor(secs / 86400)
  const h = Math.floor((secs % 86400) / 3600)
  const m = Math.floor((secs % 3600) / 60)
  return `${d}d ${h}h ${m}m`
}
</script>

<style scoped>
.process-table { background: var(--n-color-embedded); border: 1px solid var(--n-border-color); border-radius: 16px; padding: 16px; }
.toolbar { display: flex; align-items: center; gap: 12px; margin-bottom: 12px; }
.toolbar .spacer { flex: 1; }
.actions { display: flex; gap: 8px; }
.col-name { font-weight: 600; }
</style>
