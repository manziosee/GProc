<template>
  <div class="add-process-form">
    <n-form
      ref="formRef"
      :model="formData"
      :rules="rules"
      label-placement="top"
      size="medium"
    >
      <n-grid :cols="2" :x-gap="16">
        <n-form-item-gi label="Process Name" path="name">
          <n-input
            v-model:value="formData.name"
            placeholder="e.g., web-server"
            clearable
          />
        </n-form-item-gi>
        
        <n-form-item-gi label="Command" path="command">
          <n-input
            v-model:value="formData.command"
            placeholder="e.g., ./server.exe"
            clearable
          />
        </n-form-item-gi>
      </n-grid>
      
      <n-form-item label="Arguments" path="args">
        <n-dynamic-input
          v-model:value="formData.args"
          placeholder="Add argument"
          :min="0"
        />
      </n-form-item>
      
      <n-form-item label="Working Directory" path="workingDir">
        <n-input
          v-model:value="formData.workingDir"
          placeholder="e.g., /app or C:\\myapp"
          clearable
        />
      </n-form-item>
      
      <!-- Environment Variables -->
      <n-form-item label="Environment Variables">
        <div class="env-vars-container">
          <div
            v-for="(env, index) in formData.envVars"
            :key="index"
            class="env-var-row"
          >
            <n-input
              v-model:value="env.key"
              placeholder="Variable name"
              style="flex: 1"
            />
            <n-input
              v-model:value="env.value"
              placeholder="Variable value"
              style="flex: 1"
            />
            <n-button
              quaternary
              type="error"
              @click="removeEnvVar(index)"
            >
              <Trash2 :size="16" />
            </n-button>
          </div>
          <n-button
            dashed
            @click="addEnvVar"
            style="width: 100%"
          >
            <Plus :size="16" />
            Add Environment Variable
          </n-button>
        </div>
      </n-form-item>
      
      <!-- Advanced Settings -->
      <n-collapse>
        <n-collapse-item title="Advanced Settings" name="advanced">
          <n-grid :cols="2" :x-gap="16">
            <n-form-item-gi label="Process Group">
              <n-input
                v-model:value="formData.group"
                placeholder="e.g., web-services"
                clearable
              />
            </n-form-item-gi>
            
            <n-form-item-gi label="Priority">
              <n-input-number
                v-model:value="formData.priority"
                :min="-20"
                :max="20"
                placeholder="0"
              />
            </n-form-item-gi>
          </n-grid>
          
          <n-grid :cols="2" :x-gap="16">
            <n-form-item-gi label="Auto Restart">
              <n-switch v-model:value="formData.autoRestart" />
            </n-form-item-gi>
            
            <n-form-item-gi label="Max Restarts">
              <n-input-number
                v-model:value="formData.maxRestarts"
                :min="0"
                :max="100"
                :disabled="!formData.autoRestart"
              />
            </n-form-item-gi>
          </n-grid>
          
          <!-- Resource Limits -->
          <n-divider title-placement="left">Resource Limits</n-divider>
          
          <n-grid :cols="2" :x-gap="16">
            <n-form-item-gi label="Memory Limit (MB)">
              <n-input-number
                v-model:value="formData.memoryLimit"
                :min="0"
                placeholder="512"
              />
            </n-form-item-gi>
            
            <n-form-item-gi label="CPU Limit (%)">
              <n-input-number
                v-model:value="formData.cpuLimit"
                :min="0"
                :max="100"
                placeholder="50"
              />
            </n-form-item-gi>
          </n-grid>
          
          <!-- Health Check -->
          <n-divider title-placement="left">Health Check</n-divider>
          
          <n-form-item label="Health Check URL">
            <n-input
              v-model:value="formData.healthCheckUrl"
              placeholder="http://localhost:8080/health"
              clearable
            />
          </n-form-item>
          
          <n-grid :cols="2" :x-gap="16" v-if="formData.healthCheckUrl">
            <n-form-item-gi label="Check Interval (seconds)">
              <n-input-number
                v-model:value="formData.healthCheckInterval"
                :min="5"
                :max="3600"
                placeholder="30"
              />
            </n-form-item-gi>
            
            <n-form-item-gi label="Timeout (seconds)">
              <n-input-number
                v-model:value="formData.healthCheckTimeout"
                :min="1"
                :max="60"
                placeholder="5"
              />
            </n-form-item-gi>
          </n-grid>
          
          <!-- Notifications -->
          <n-divider title-placement="left">Notifications</n-divider>
          
          <n-grid :cols="2" :x-gap="16">
            <n-form-item-gi label="Email Notifications">
              <n-input
                v-model:value="formData.notifyEmail"
                placeholder="admin@company.com"
                clearable
              />
            </n-form-item-gi>
            
            <n-form-item-gi label="Slack Webhook">
              <n-input
                v-model:value="formData.notifySlack"
                placeholder="https://hooks.slack.com/..."
                clearable
              />
            </n-form-item-gi>
          </n-grid>
        </n-collapse-item>
      </n-collapse>
    </n-form>
    
    <div class="form-actions">
      <n-button @click="$emit('cancel')">
        Cancel
      </n-button>
      <n-button
        type="primary"
        @click="handleSubmit"
        :loading="submitting"
      >
        <Play :size="16" />
        Create Process
      </n-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import {
  NForm,
  NFormItem,
  NFormItemGi,
  NGrid,
  NInput,
  NInputNumber,
  NButton,
  NSwitch,
  NDynamicInput,
  NCollapse,
  NCollapseItem,
  NDivider,
  FormInst,
  FormRules
} from 'naive-ui'
import { Plus, Trash2, Play } from 'lucide-vue-next'

defineEmits<{
  submit: [data: any]
  cancel: []
}>()

const formRef = ref<FormInst | null>(null)
const submitting = ref(false)

const formData = reactive({
  name: '',
  command: '',
  args: [] as string[],
  workingDir: '',
  envVars: [] as Array<{ key: string; value: string }>,
  group: '',
  priority: 0,
  autoRestart: true,
  maxRestarts: 5,
  memoryLimit: null as number | null,
  cpuLimit: null as number | null,
  healthCheckUrl: '',
  healthCheckInterval: 30,
  healthCheckTimeout: 5,
  notifyEmail: '',
  notifySlack: ''
})

const rules: FormRules = {
  name: [
    { required: true, message: 'Process name is required' },
    { min: 2, max: 50, message: 'Name must be between 2 and 50 characters' }
  ],
  command: [
    { required: true, message: 'Command is required' }
  ]
}

const addEnvVar = () => {
  formData.envVars.push({ key: '', value: '' })
}

const removeEnvVar = (index: number) => {
  formData.envVars.splice(index, 1)
}

const handleSubmit = async () => {
  try {
    await formRef.value?.validate()
    submitting.value = true
    
    // Convert form data to API format
    const processData = {
      name: formData.name,
      command: formData.command,
      args: formData.args.filter(arg => arg.trim() !== ''),
      workingDir: formData.workingDir || undefined,
      env: formData.envVars.reduce((acc, env) => {
        if (env.key && env.value) {
          acc[env.key] = env.value
        }
        return acc
      }, {} as Record<string, string>),
      group: formData.group || undefined,
      priority: formData.priority,
      autoRestart: formData.autoRestart,
      maxRestarts: formData.maxRestarts,
      resourceLimit: (formData.memoryLimit || formData.cpuLimit) ? {
        memoryMB: formData.memoryLimit || 0,
        cpuLimit: formData.cpuLimit || 0
      } : undefined,
      healthCheck: formData.healthCheckUrl ? {
        url: formData.healthCheckUrl,
        interval: formData.healthCheckInterval,
        timeout: formData.healthCheckTimeout,
        retries: 3
      } : undefined,
      notifications: (formData.notifyEmail || formData.notifySlack) ? {
        email: formData.notifyEmail || undefined,
        slack: formData.notifySlack || undefined
      } : undefined
    }
    
    // Emit the form data
    $emit('submit', processData)
  } catch (error) {
    console.error('Form validation failed:', error)
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.add-process-form {
  max-width: 800px;
  margin: 0 auto;
}

.env-vars-container {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.env-var-row {
  display: flex;
  gap: 12px;
  align-items: center;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 24px;
  padding-top: 24px;
  border-top: 1px solid var(--n-border-color);
}

:deep(.n-collapse-item__header) {
  font-weight: 600;
}

:deep(.n-divider__title) {
  font-size: 14px;
  font-weight: 600;
  color: var(--n-text-color-2);
}
</style>