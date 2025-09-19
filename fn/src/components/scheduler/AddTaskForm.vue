<template>
  <div class="add-task-form">
    <n-form ref="formRef" :model="formData" :rules="rules" label-placement="top">
      <n-form-item label="Task Name" path="name">
        <n-input 
          v-model:value="formData.name" 
          placeholder="Enter task name"
        />
      </n-form-item>
      
      <n-form-item label="Description">
        <n-input 
          v-model:value="formData.description" 
          placeholder="Enter task description"
          type="textarea"
          :rows="2"
        />
      </n-form-item>
      
      <n-form-item label="Command" path="command">
        <n-input 
          v-model:value="formData.command" 
          placeholder="Enter command to execute"
          type="textarea"
          :rows="2"
        />
      </n-form-item>
      
      <n-form-item label="Schedule (Cron Expression)" path="schedule">
        <n-input 
          v-model:value="formData.schedule" 
          placeholder="0 2 * * *"
        />
        <template #feedback>
          <div class="cron-help">
            <span>Format: minute hour day month dayofweek</span>
            <a href="#" @click.prevent="showCronHelp = !showCronHelp">Examples</a>
          </div>
          
          <div v-if="showCronHelp" class="cron-examples">
            <div><code>0 2 * * *</code> - Daily at 2:00 AM</div>
            <div><code>0 9 * * 1-5</code> - Weekdays at 9:00 AM</div>
            <div><code>*/15 * * * *</code> - Every 15 minutes</div>
            <div><code>0 0 1 * *</code> - First day of every month</div>
          </div>
        </template>
      </n-form-item>
      
      <n-form-item label="Timezone">
        <n-select
          v-model:value="formData.timezone"
          :options="timezoneOptions"
          filterable
          placeholder="Select timezone"
        />
      </n-form-item>
      
      <n-form-item>
        <n-checkbox v-model:checked="formData.enabled">
          Enable task immediately
        </n-checkbox>
      </n-form-item>
      
      <div class="form-actions">
        <n-button @click="$emit('cancel')">Cancel</n-button>
        <n-button type="primary" @click="handleSubmit">Add Task</n-button>
      </div>
    </n-form>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { NForm, NFormItem, NInput, NSelect, NButton, NCheckbox } from 'naive-ui'

const emit = defineEmits<{
  submit: [data: any]
  cancel: []
}>()

const formRef = ref()
const showCronHelp = ref(false)

const formData = reactive({
  name: '',
  description: '',
  command: '',
  schedule: '',
  timezone: 'UTC',
  enabled: true
})

const rules = {
  name: {
    required: true,
    message: 'Task name is required'
  },
  command: {
    required: true,
    message: 'Command is required'
  },
  schedule: {
    required: true,
    message: 'Cron expression is required'
  }
}

const timezoneOptions = [
  { label: 'UTC', value: 'UTC' },
  { label: 'America/New_York', value: 'America/New_York' },
  { label: 'America/Los_Angeles', value: 'America/Los_Angeles' },
  { label: 'Europe/London', value: 'Europe/London' },
  { label: 'Europe/Paris', value: 'Europe/Paris' },
  { label: 'Asia/Tokyo', value: 'Asia/Tokyo' },
  { label: 'Australia/Sydney', value: 'Australia/Sydney' }
]

const handleSubmit = async () => {
  try {
    await formRef.value?.validate()
    emit('submit', { ...formData })
  } catch (error) {
    console.error('Validation failed:', error)
  }
}
</script>

<style scoped>
.add-task-form {
  padding: 20px;
  max-width: 500px;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 24px;
}

.cron-help {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
  color: var(--n-text-color-3);
  margin-top: 4px;
}

.cron-help a {
  color: var(--n-primary-color);
  text-decoration: none;
}

.cron-help a:hover {
  text-decoration: underline;
}

.cron-examples {
  margin-top: 8px;
  padding: 8px;
  background: var(--n-color);
  border-radius: 4px;
  font-size: 12px;
}

.cron-examples div {
  margin-bottom: 4px;
}

.cron-examples code {
  background: var(--n-color-embedded);
  padding: 2px 4px;
  border-radius: 2px;
  margin-right: 8px;
  font-family: monospace;
}
</style>