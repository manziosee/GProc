<template>
  <div class="scheduled-tasks">
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">Scheduled Tasks</h1>
        <p class="page-subtitle">Manage cron jobs and scheduled processes</p>
      </div>
      <n-button type="primary" size="large" @click="showAddModal = true">
        <template #icon><span>⏰</span></template>
        Add Task
      </n-button>
    </div>

    <div class="tasks-stats">
      <div class="stat-card active">
        <div class="stat-icon">🟢</div>
        <div class="stat-info">
          <div class="stat-value">{{ activeTasks }}</div>
          <div class="stat-label">Active Tasks</div>
        </div>
      </div>
      <div class="stat-card next">
        <div class="stat-icon">⏱️</div>
        <div class="stat-info">
          <div class="stat-value">{{ nextTaskTime }}</div>
          <div class="stat-label">Next Run</div>
        </div>
      </div>
      <div class="stat-card total">
        <div class="stat-icon">📊</div>
        <div class="stat-info">
          <div class="stat-value">{{ tasks.length }}</div>
          <div class="stat-label">Total Tasks</div>
        </div>
      </div>
    </div>

    <div class="tasks-container">
      <div class="tasks-header">
        <n-input placeholder="Search tasks..." style="width: 300px">
          <template #prefix>🔍</template>
        </n-input>
        <n-select v-model:value="statusFilter" :options="statusOptions" placeholder="Filter by status" style="width: 150px" />
      </div>

      <div class="tasks-list">
        <div v-for="task in filteredTasks" :key="task.id" class="task-card">
          <div class="task-status">
            <div class="status-indicator" :class="task.status"></div>
          </div>
          
          <div class="task-info">
            <div class="task-header-row">
              <h3 class="task-name">{{ task.name }}</h3>
              <n-tag :type="getStatusTagType(task.status)" size="small">{{ task.status }}</n-tag>
            </div>
            
            <div class="task-command">
              <span class="command-label">Command:</span>
              <code class="command-text">{{ task.command }}</code>
            </div>
            
            <div class="task-schedule">
              <div class="schedule-item">
                <span class="schedule-label">🕰️ Cron:</span>
                <code class="cron-expression">{{ task.cron }}</code>
                <span class="cron-description">({{ getCronDescription(task.cron) }})</span>
              </div>
              
              <div class="schedule-item">
                <span class="schedule-label">⏭️ Next:</span>
                <span class="next-run">{{ task.nextRun }}</span>
              </div>
              
              <div class="schedule-item">
                <span class="schedule-label">📅 Last:</span>
                <span class="last-run">{{ task.lastRun || 'Never' }}</span>
              </div>
            </div>
          </div>
          
          <div class="task-actions">
            <n-space vertical>
              <n-button size="small" :type="task.status === 'active' ? 'warning' : 'primary'" @click="toggleTask(task)">
                {{ task.status === 'active' ? 'Pause' : 'Resume' }}
              </n-button>
              <n-button size="small" @click="runNow(task)">Run Now</n-button>
              <n-dropdown :options="getTaskActions(task)" @select="handleTaskAction">
                <n-button size="small" quaternary>
                  <template #icon><span>⋮</span></template>
                </n-button>
              </n-dropdown>
            </n-space>
          </div>
        </div>
      </div>
    </div>

    <div class="cron-helper">
      <h3 class="helper-title">Cron Expression Helper</h3>
      <div class="cron-examples">
        <div class="cron-example" v-for="example in cronExamples" :key="example.expression">
          <code class="example-cron">{{ example.expression }}</code>
          <span class="example-description">{{ example.description }}</span>
        </div>
      </div>
    </div>

    <n-modal v-model:show="showAddModal" preset="card" title="Add Scheduled Task" style="width: 600px">
      <AddTaskForm @submit="handleAddTask" @cancel="showAddModal = false" />
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { NButton, NInput, NSelect, NTag, NSpace, NDropdown, NModal } from 'naive-ui'
import AddTaskForm from '../components/scheduler/AddTaskForm.vue'

const showAddModal = ref(false)
const statusFilter = ref('')

const tasks = ref([
  {
    id: '1',
    name: 'Database Backup',
    command: './backup.sh',
    cron: '0 2 * * *',
    status: 'active',
    nextRun: 'Today at 2:00 AM',
    lastRun: 'Yesterday at 2:00 AM'
  },
  {
    id: '2',
    name: 'Log Cleanup',
    command: './cleanup-logs.sh',
    cron: '0 0 */7 * *',
    status: 'active',
    nextRun: 'Sunday at 12:00 AM',
    lastRun: 'Last Sunday at 12:00 AM'
  },
  {
    id: '3',
    name: 'Health Check',
    command: './health-check.sh',
    cron: '*/5 * * * *',
    status: 'paused',
    nextRun: 'Paused',
    lastRun: '5 minutes ago'
  },
  {
    id: '4',
    name: 'Report Generation',
    command: './generate-report.py',
    cron: '0 9 1 * *',
    status: 'active',
    nextRun: '1st of next month at 9:00 AM',
    lastRun: '1st of this month at 9:00 AM'
  }
])

const statusOptions = [
  { label: 'All Status', value: '' },
  { label: 'Active', value: 'active' },
  { label: 'Paused', value: 'paused' },
  { label: 'Failed', value: 'failed' }
]

const cronExamples = [
  { expression: '0 2 * * *', description: 'Daily at 2:00 AM' },
  { expression: '0 9 * * 1-5', description: 'Weekdays at 9:00 AM' },
  { expression: '*/15 * * * *', description: 'Every 15 minutes' },
  { expression: '0 0 1 * *', description: 'First day of every month' },
  { expression: '0 0 * * 0', description: 'Every Sunday at midnight' },
  { expression: '30 14 * * 1', description: 'Every Monday at 2:30 PM' }
]

const activeTasks = computed(() => tasks.value.filter(t => t.status === 'active').length)
const nextTaskTime = computed(() => {
  const nextTask = tasks.value
    .filter(t => t.status === 'active')
    .sort((a, b) => a.nextRun.localeCompare(b.nextRun))[0]
  return nextTask ? '2:00 AM' : 'None'
})

const filteredTasks = computed(() => {
  return tasks.value.filter(task => {
    if (statusFilter.value && task.status !== statusFilter.value) return false
    return true
  })
})

const getStatusTagType = (status: string) => {
  switch (status) {
    case 'active': return 'success'
    case 'paused': return 'warning'
    case 'failed': return 'error'
    default: return 'default'
  }
}

const getCronDescription = (cron: string) => {
  const descriptions: { [key: string]: string } = {
    '0 2 * * *': 'Daily at 2:00 AM',
    '0 0 */7 * *': 'Every 7 days',
    '*/5 * * * *': 'Every 5 minutes',
    '0 9 1 * *': 'Monthly on 1st at 9:00 AM'
  }
  return descriptions[cron] || 'Custom schedule'
}

const getTaskActions = (task: any) => [
  { label: 'Edit Task', key: 'edit' },
  { label: 'View Logs', key: 'logs' },
  { label: 'Duplicate', key: 'duplicate' },
  { label: 'Delete Task', key: 'delete', props: { style: 'color: red' } }
]

const toggleTask = (task: any) => {
  task.status = task.status === 'active' ? 'paused' : 'active'
  console.log('Toggle task:', task.name)
}

const runNow = (task: any) => {
  console.log('Run task now:', task.name)
}

const handleTaskAction = (key: string) => {
  console.log('Task action:', key)
}

const handleAddTask = (data: any) => {
  console.log('Add task:', data)
  showAddModal.value = false
}
</script>

<style scoped>
.scheduled-tasks {
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

.tasks-stats {
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

.stat-card.active {
  --accent-color: #10b981;
  --accent-color-light: #34d399;
}

.stat-card.next {
  --accent-color: #f59e0b;
  --accent-color-light: #fbbf24;
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

.tasks-container {
  background: var(--n-color-embedded);
  border: 1px solid var(--n-border-color);
  border-radius: 16px;
  padding: 24px;
  margin-bottom: 32px;
}

.tasks-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  gap: 16px;
}

.tasks-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.task-card {
  background: var(--n-color);
  border: 1px solid var(--n-border-color);
  border-radius: 12px;
  padding: 20px;
  display: flex;
  align-items: flex-start;
  gap: 16px;
  transition: all 0.3s ease;
}

.task-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.task-status {
  flex-shrink: 0;
  padding-top: 4px;
}

.status-indicator {
  width: 12px;
  height: 12px;
  border-radius: 50%;
}

.status-indicator.active {
  background: #10b981;
  box-shadow: 0 0 0 2px rgba(16, 185, 129, 0.2);
}

.status-indicator.paused {
  background: #f59e0b;
  box-shadow: 0 0 0 2px rgba(245, 158, 11, 0.2);
}

.status-indicator.failed {
  background: #ef4444;
  box-shadow: 0 0 0 2px rgba(239, 68, 68, 0.2);
}

.task-info {
  flex: 1;
  min-width: 0;
}

.task-header-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.task-name {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--n-text-color);
  margin: 0;
}

.task-command {
  margin-bottom: 16px;
}

.command-label {
  font-size: 0.875rem;
  color: var(--n-text-color-2);
  margin-right: 8px;
}

.command-text {
  font-family: 'Consolas', 'Monaco', monospace;
  background: var(--n-color-embedded);
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 0.875rem;
}

.task-schedule {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.schedule-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.875rem;
}

.schedule-label {
  color: var(--n-text-color-2);
  min-width: 60px;
}

.cron-expression {
  font-family: 'Consolas', 'Monaco', monospace;
  background: var(--n-color-embedded);
  padding: 2px 6px;
  border-radius: 3px;
  font-size: 0.8rem;
}

.cron-description {
  color: var(--n-text-color-3);
  font-style: italic;
}

.next-run, .last-run {
  color: var(--n-text-color);
}

.task-actions {
  flex-shrink: 0;
}

.cron-helper {
  background: var(--n-color-embedded);
  border: 1px solid var(--n-border-color);
  border-radius: 16px;
  padding: 24px;
}

.helper-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--n-text-color);
  margin: 0 0 16px 0;
}

.cron-examples {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 12px;
}

.cron-example {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px;
  background: var(--n-color);
  border-radius: 6px;
}

.example-cron {
  font-family: 'Consolas', 'Monaco', monospace;
  background: var(--n-color-embedded);
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 0.875rem;
  min-width: 100px;
}

.example-description {
  color: var(--n-text-color-2);
  font-size: 0.875rem;
}

@media (max-width: 768px) {
  .scheduled-tasks {
    padding: 16px;
  }
  
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }
  
  .tasks-header {
    flex-direction: column;
    align-items: stretch;
  }
  
  .task-card {
    flex-direction: column;
  }
  
  .task-header-row {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
  
  .schedule-item {
    flex-wrap: wrap;
  }
}
</style>