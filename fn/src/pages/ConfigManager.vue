<template>
  <div class="config-manager">
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">Configuration Manager</h1>
        <p class="page-subtitle">Manage your GProc configuration files</p>
      </div>
      <n-space>
        <n-button @click="importConfig">
          <template #icon><span>📁</span></template>
          Import
        </n-button>
        <n-button type="primary" @click="saveConfig">
          <template #icon><span>💾</span></template>
          Save Config
        </n-button>
      </n-space>
    </div>

    <div class="config-tabs">
      <n-tabs v-model:value="activeTab" type="card" size="large">
        <n-tab-pane name="yaml" tab="📄 YAML Config">
          <div class="config-editor">
            <div class="editor-header">
              <div class="file-info">
                <span class="file-name">gproc.yaml</span>
                <n-tag size="small" type="success">Valid</n-tag>
              </div>
              <n-space>
                <n-button size="small" @click="formatConfig">Format</n-button>
                <n-button size="small" @click="validateConfig">Validate</n-button>
              </n-space>
            </div>
            <div class="editor-container">
              <textarea 
                v-model="yamlConfig" 
                class="config-textarea"
                placeholder="Enter your YAML configuration..."
                spellcheck="false"
              ></textarea>
            </div>
          </div>
        </n-tab-pane>
        
        <n-tab-pane name="json" tab="📄 JSON Config">
          <div class="config-editor">
            <div class="editor-header">
              <div class="file-info">
                <span class="file-name">gproc.json</span>
                <n-tag size="small" type="success">Valid</n-tag>
              </div>
              <n-space>
                <n-button size="small" @click="formatConfig">Format</n-button>
                <n-button size="small" @click="validateConfig">Validate</n-button>
              </n-space>
            </div>
            <div class="editor-container">
              <textarea 
                v-model="jsonConfig" 
                class="config-textarea"
                placeholder="Enter your JSON configuration..."
                spellcheck="false"
              ></textarea>
            </div>
          </div>
        </n-tab-pane>
        
        <n-tab-pane name="env" tab="🌍 Environment">
          <div class="env-manager">
            <div class="env-header">
              <h3>Environment Variables</h3>
              <n-button @click="addEnvVar">
                <template #icon><span>➕</span></template>
                Add Variable
              </n-button>
            </div>
            
            <div class="env-list">
              <div v-for="(env, index) in envVars" :key="index" class="env-item">
                <n-input v-model:value="env.key" placeholder="Variable name" style="flex: 1" />
                <span class="env-separator">=</span>
                <n-input v-model:value="env.value" placeholder="Variable value" style="flex: 2" />
                <n-button size="small" type="error" @click="removeEnvVar(index)">
                  <template #icon><span>🗑️</span></template>
                </n-button>
              </div>
            </div>
          </div>
        </n-tab-pane>
      </n-tabs>
    </div>

    <div class="config-templates">
      <h3 class="templates-title">Configuration Templates</h3>
      <div class="templates-grid">
        <div v-for="template in templates" :key="template.name" class="template-card" @click="loadTemplate(template)">
          <div class="template-icon">{{ template.icon }}</div>
          <div class="template-info">
            <h4 class="template-name">{{ template.name }}</h4>
            <p class="template-description">{{ template.description }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { NTabs, NTabPane, NButton, NSpace, NTag, NInput } from 'naive-ui'

const activeTab = ref('yaml')

const yamlConfig = ref(`# GProc Configuration
processes:
  - name: web-server
    command: ./server
    args: ["--port", "8080"]
    env:
      NODE_ENV: production
      PORT: "8080"
    auto_restart: true
    max_restarts: 5
    health_check:
      url: http://localhost:8080/health
      interval: 30s
      timeout: 5s
      retries: 3

  - name: worker
    command: ./worker
    auto_restart: true
    max_restarts: 3

groups:
  - name: web-stack
    processes: ["web-server", "worker"]

log_dir: ./logs
web_port: 3000`)

const jsonConfig = ref(`{
  "processes": [
    {
      "name": "web-server",
      "command": "./server",
      "args": ["--port", "8080"],
      "env": {
        "NODE_ENV": "production",
        "PORT": "8080"
      },
      "auto_restart": true,
      "max_restarts": 5,
      "health_check": {
        "url": "http://localhost:8080/health",
        "interval": "30s",
        "timeout": "5s",
        "retries": 3
      }
    }
  ],
  "log_dir": "./logs",
  "web_port": 3000
}`)

const envVars = ref([
  { key: 'NODE_ENV', value: 'production' },
  { key: 'PORT', value: '8080' },
  { key: 'DB_HOST', value: 'localhost' }
])

const templates = ref([
  {
    name: 'Web Application',
    description: 'Basic web server with health checks',
    icon: '🌐',
    config: 'web-app-template'
  },
  {
    name: 'Microservices',
    description: 'Multiple services with load balancing',
    icon: '🔍',
    config: 'microservices-template'
  },
  {
    name: 'Background Jobs',
    description: 'Worker processes and schedulers',
    icon: '⚙️',
    config: 'background-jobs-template'
  },
  {
    name: 'Development',
    description: 'Development environment setup',
    icon: '🛠️',
    config: 'development-template'
  }
])

const addEnvVar = () => {
  envVars.value.push({ key: '', value: '' })
}

const removeEnvVar = (index: number) => {
  envVars.value.splice(index, 1)
}

const loadTemplate = (template: any) => {
  console.log('Loading template:', template.name)
}

const saveConfig = () => {
  console.log('Saving configuration...')
}

const importConfig = () => {
  console.log('Importing configuration...')
}

const formatConfig = () => {
  console.log('Formatting configuration...')
}

const validateConfig = () => {
  console.log('Validating configuration...')
}
</script>

<style scoped>
.config-manager {
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

.config-tabs {
  margin-bottom: 32px;
}

.config-editor {
  background: var(--n-color-embedded);
  border: 1px solid var(--n-border-color);
  border-radius: 12px;
  overflow: hidden;
}

.editor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  background: var(--n-color);
  border-bottom: 1px solid var(--n-border-color);
}

.file-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.file-name {
  font-family: 'Consolas', 'Monaco', monospace;
  font-weight: 600;
  color: var(--n-text-color);
}

.editor-container {
  position: relative;
}

.config-textarea {
  width: 100%;
  height: 400px;
  padding: 20px;
  border: none;
  outline: none;
  resize: vertical;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.5;
  background: var(--n-color-embedded);
  color: var(--n-text-color);
}

.env-manager {
  padding: 20px;
}

.env-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.env-header h3 {
  margin: 0;
  color: var(--n-text-color);
}

.env-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.env-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: var(--n-color-embedded);
  border: 1px solid var(--n-border-color);
  border-radius: 8px;
}

.env-separator {
  font-family: monospace;
  font-weight: bold;
  color: var(--n-text-color-2);
}

.config-templates {
  margin-top: 32px;
}

.templates-title {
  font-size: 1.5rem;
  font-weight: 600;
  color: var(--n-text-color);
  margin: 0 0 20px 0;
}

.templates-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 20px;
}

.template-card {
  background: var(--n-color-embedded);
  border: 1px solid var(--n-border-color);
  border-radius: 12px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.template-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  border-color: var(--n-primary-color);
}

.template-icon {
  font-size: 2rem;
  opacity: 0.8;
}

.template-info {
  flex: 1;
}

.template-name {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--n-text-color);
  margin: 0 0 4px 0;
}

.template-description {
  font-size: 0.875rem;
  color: var(--n-text-color-2);
  margin: 0;
}

@media (max-width: 768px) {
  .config-manager {
    padding: 16px;
  }
  
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }
  
  .config-textarea {
    height: 300px;
  }
  
  .env-item {
    flex-direction: column;
    align-items: stretch;
  }
  
  .env-separator {
    display: none;
  }
}
</style>