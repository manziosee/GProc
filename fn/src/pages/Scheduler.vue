<template>
  <div class="scheduler">
    <n-card title="Cron Task Scheduler">
      <template #header-extra>
        <n-button type="primary" @click="showCreateModal = true">
          <template #icon><Clock /></template>
          New Task
        </n-button>
      </template>
      
      <!-- Tasks Table -->
      <n-data-table
        :columns="columns"
        :data="tasks"
        :loading="loading"
        :pagination="{ pageSize: 10 }"
      />
    </n-card>
    
    <!-- Create Task Modal -->
    <n-modal v-model:show="showCreateModal" title="Create Scheduled Task">
      <n-card style="width: 600px">
        <n-form :model="newTask" :rules="taskRules">
          <n-form-item path="name" label="Task Name">
            <n-input v-model:value="newTask.name" placeholder="backup-database" />
          </n-form-item>
          
          <n-form-item path="command" label="Command">
            <n-input v-model:value="newTask.command" placeholder="/usr/bin/backup.sh" />
          </n-form-item>
          
          <n-form-item label="Arguments">
            <n-dynamic-input
              v-model:value="newTask.args"
              placeholder="Add argument"
            />
          </n-form-item>
          
          <n-form-item path="cron" label="Cron Expression">
            <n-input v-model:value="newTask.cron" placeholder="0 2 * * *" />
            <template #feedback>
              <n-space vertical size="small">
                <span>{{ cronDescription }}</span>
                <n-space>
                  <n-button size="tiny" @click="setCronPreset('0 * * * *')">Hourly</n-button>
                  <n-button size="tiny" @click="setCronPreset('0 0 * * *')">Daily</n-button>
                  <n-button size="tiny" @click="setCronPreset('0 0 * * 0')">Weekly</n-button>
                  <n-button size="tiny" @click="setCronPreset('0 0 1 * *')">Monthly</n-button>
                </n-space>
              </n-space>
            </template>
          </n-form-item>
          
          <n-form-item label="Description">
            <n-input
              v-model:value="newTask.description"
              type="textarea"
              placeholder="Task description"
            />
          </n-form-item>
          
          <n-form-item label="Timeout">
            <n-input v-model:value="newTask.timeout" placeholder="1h" />
          </n-form-item>
          
          <n-form-item label="Working Directory">
            <n-input v-model:value="newTask.workingDir" placeholder="/app" />
          </n-form-item>
          
          <n-form-item label="Environment Variables">
            <n-dynamic-input
              v-model:value="envVars"
              placeholder="KEY=VALUE"
            />
          </n-form-item>
        </n-form>
        
        <template #footer>
          <n-space justify="end">
            <n-button @click="showCreateModal = false">Cancel</n-button>
            <n-button type="primary" @click="createTask">Create Task</n-button>
          </n-space>
        </template>
      </n-card>
    </n-modal>
    
    <!-- Task Executions Modal -->
    <n-modal v-model:show="showExecutionsModal" title="Task Executions">
      <n-card style="width: 800px">
        <n-data-table
          :columns="executionColumns"
          :data="executions"
          :pagination="{ pageSize: 10 }"
        />
      </n-card>
    </n-modal>
    
    <!-- Execution Logs Modal -->
    <n-modal v-model:show="showLogsModal" title="Execution Logs">
      <n-card style="width: 800px">
        <n-code
          :code="executionLogs"
          language="text"
          show-line-numbers
        />
      </n-card>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, h } from 'vue'
import { useSchedulerStore } from '../stores/scheduler'
import { useMessage } from 'naive-ui'
import { Clock, Play, Pause, Trash2, Eye, FileText } from 'lucide-vue-next'

const schedulerStore = useSchedulerStore()
const message = useMessage()

const loading = ref(false)
const showCreateModal = ref(false)
const showExecutionsModal = ref(false)
const showLogsModal = ref(false)
const executionLogs = ref('')
const envVars = ref<string[]>([])

const newTask = ref({
  name: '',
  command: '',
  args: [] as string[],
  cron: '',
  description: '',
  timeout: '1h',
  workingDir: ''
})

const { tasks, executions } = schedulerStore

const taskRules = {
  name: { required: true, message: 'Task name is required' },
  command: { required: true, message: 'Command is required' },
  cron: { required: true, message: 'Cron expression is required' }
}

const cronDescription = computed(() => {
  const cron = newTask.value.cron
  if (!cron) return 'Enter a cron expression'
  
  // Simple cron description logic
  const parts = cron.split(' ')
  if (parts.length !== 5) return 'Invalid cron expression'
  
  const [minute, hour, day, month, weekday] = parts
  
  if (cron === '0 * * * *') return 'Every hour'
  if (cron === '0 0 * * *') return 'Daily at midnight'
  if (cron === '0 0 * * 0') return 'Weekly on Sunday at midnight'
  if (cron === '0 0 1 * *') return 'Monthly on the 1st at midnight'
  
  return `At ${minute}:${hour} on day ${day} of month ${month}, weekday ${weekday}`
})

const columns = [
  { title: 'Name', key: 'name' },
  { title: 'Command', key: 'command' },
  { title: 'Cron', key: 'cron' },
  {
    title: 'Status',
    key: 'enabled',
    render: (row: any) => h('n-tag', {
      type: row.enabled ? 'success' : 'default'
    }, { default: () => row.enabled ? 'Enabled' : 'Disabled' })
  },
  {
    title: 'Next Run',
    key: 'nextRun',
    render: (row: any) => new Date(row.nextRun).toLocaleString()
  },
  {
    title: 'Last Run',
    key: 'lastRun',
    render: (row: any) => row.lastRun ? new Date(row.lastRun).toLocaleString() : 'Never'
  },
  {
    title: 'Result',
    key: 'lastResult',
    render: (row: any) => {
      if (!row.lastResult) return '-'
      return h('n-tag', {
        type: row.lastResult === 'success' ? 'success' : 'error'
      }, { default: () => row.lastResult })
    }
  },
  { title: 'Runs', key: 'runCount' },
  {
    title: 'Actions',
    key: 'actions',
    render: (row: any) => h('n-space', [
      h('n-button', {
        size: 'small',
        type: 'primary',
        onClick: () => runTaskNow(row.id)
      }, { default: () => 'Run Now', icon: () => h(Play) }),
      h('n-button', {
        size: 'small',
        onClick: () => toggleTask(row.id, !row.enabled)
      }, { 
        default: () => row.enabled ? 'Disable' : 'Enable',
        icon: () => h(row.enabled ? Pause : Play)
      }),
      h('n-button', {
        size: 'small',
        onClick: () => viewExecutions(row.id)
      }, { default: () => 'History', icon: () => h(Eye) }),
      h('n-button', {
        size: 'small',
        type: 'error',
        onClick: () => deleteTask(row.id)
      }, { default: () => 'Delete', icon: () => h(Trash2) })
    ])
  }
]

const executionColumns = [
  {
    title: 'Status',
    key: 'status',
    render: (row: any) => h('n-tag', {
      type: row.status === 'completed' ? 'success' : 
            row.status === 'failed' ? 'error' : 'info'
    }, { default: () => row.status })
  },
  {
    title: 'Start Time',
    key: 'startTime',
    render: (row: any) => new Date(row.startTime).toLocaleString()
  },
  {
    title: 'Duration',
    key: 'duration',
    render: (row: any) => {
      if (!row.endTime) return 'Running...'
      const duration = new Date(row.endTime).getTime() - new Date(row.startTime).getTime()
      return `${Math.round(duration / 1000)}s`
    }
  },
  { title: 'Exit Code', key: 'exitCode' },
  {
    title: 'Actions',
    key: 'actions',
    render: (row: any) => h('n-button', {
      size: 'small',
      onClick: () => viewExecutionLogs(row.id)
    }, { default: () => 'Logs', icon: () => h(FileText) })
  }
]

const setCronPreset = (cron: string) => {
  newTask.value.cron = cron
}

const createTask = async () => {
  const env = envVars.value.reduce((acc, envVar) => {
    const [key, value] = envVar.split('=')
    if (key && value) acc[key] = value
    return acc
  }, {} as Record<string, string>)

  const result = await schedulerStore.createTask({
    ...newTask.value,
    env
  })
  
  if (result.success) {
    message.success('Task created successfully')
    showCreateModal.value = false
    resetForm()
  } else {
    message.error(result.error)
  }
}

const runTaskNow = async (taskId: string) => {
  const result = await schedulerStore.runTaskNow(taskId)
  if (result.success) {
    message.success('Task started')
  } else {
    message.error(result.error)
  }
}

const toggleTask = async (taskId: string, enabled: boolean) => {
  const result = await schedulerStore.updateTask(taskId, { enabled })
  if (result.success) {
    message.success(`Task ${enabled ? 'enabled' : 'disabled'}`)
  } else {
    message.error(result.error)
  }
}

const deleteTask = async (taskId: string) => {
  const result = await schedulerStore.deleteTask(taskId)
  if (result.success) {
    message.success('Task deleted')
  } else {
    message.error(result.error)
  }
}

const viewExecutions = async (taskId: string) => {
  await schedulerStore.fetchExecutions(taskId)
  showExecutionsModal.value = true
}

const viewExecutionLogs = async (executionId: string) => {
  const result = await schedulerStore.getExecutionLogs(executionId)
  if (result.success) {
    executionLogs.value = result.logs
    showLogsModal.value = true
  } else {
    message.error(result.error)
  }
}

const resetForm = () => {
  newTask.value = {
    name: '',
    command: '',
    args: [],
    cron: '',
    description: '',
    timeout: '1h',
    workingDir: ''
  }
  envVars.value = []
}

onMounted(() => {
  schedulerStore.fetchTasks()
})
</script>

<style scoped>
.scheduler {
  padding: 1rem;
}
</style>