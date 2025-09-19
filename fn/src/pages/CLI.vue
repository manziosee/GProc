<template>
  <div class="cli-page">
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">CLI Reference</h1>
        <p class="page-subtitle">Complete command-line interface documentation</p>
      </div>
      <n-button type="primary" @click="copyInstallCommand">
        <template #icon><span>📋</span></template>
        Copy Install
      </n-button>
    </div>

    <div class="cli-sections">
      <div class="install-section">
        <h2 class="section-title">🚀 Installation</h2>
        <div class="code-block">
          <div class="code-header">
            <span class="code-title">Build from source</span>
            <n-button size="small" text @click="copyToClipboard(installCommand)">Copy</n-button>
          </div>
          <pre class="code-content">{{ installCommand }}</pre>
        </div>
      </div>

      <div class="commands-grid">
        <div v-for="category in commandCategories" :key="category.name" class="command-category">
          <div class="category-header">
            <h3 class="category-title">{{ category.icon }} {{ category.name }}</h3>
            <span class="category-count">{{ category.commands.length }} commands</span>
          </div>
          
          <div class="commands-list">
            <div v-for="cmd in category.commands" :key="cmd.name" class="command-item">
              <div class="command-header">
                <code class="command-name">{{ cmd.name }}</code>
                <n-tag size="small" :type="cmd.type || 'default'">{{ cmd.category }}</n-tag>
              </div>
              <p class="command-description">{{ cmd.description }}</p>
              
              <div class="command-examples">
                <div v-for="example in cmd.examples" :key="example.command" class="example-item">
                  <div class="example-header">
                    <span class="example-label">{{ example.label }}</span>
                    <n-button size="tiny" text @click="copyToClipboard(example.command)">Copy</n-button>
                  </div>
                  <code class="example-command">{{ example.command }}</code>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="quick-reference">
      <h2 class="section-title">⚡ Quick Reference</h2>
      <div class="reference-grid">
        <div class="reference-card">
          <h4>Process Management</h4>
          <div class="reference-commands">
            <code>gproc start &lt;name&gt; &lt;command&gt;</code>
            <code>gproc stop &lt;name&gt;</code>
            <code>gproc restart &lt;name&gt;</code>
            <code>gproc list</code>
          </div>
        </div>
        
        <div class="reference-card">
          <h4>Monitoring</h4>
          <div class="reference-commands">
            <code>gproc logs &lt;name&gt;</code>
            <code>gproc status &lt;name&gt;</code>
            <code>gproc web --port 3000</code>
          </div>
        </div>
        
        <div class="reference-card">
          <h4>Advanced</h4>
          <div class="reference-commands">
            <code>gproc cluster start &lt;name&gt;</code>
            <code>gproc schedule &lt;name&gt; --cron</code>
            <code>gproc daemon</code>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { NButton, NTag } from 'naive-ui'

const installCommand = `# Clone and build GProc
git clone https://github.com/manziosee/GProc.git
cd GProc
go build -o gproc.exe cmd/main.go cmd/daemon.go

# Add to PATH (Windows)
set PATH=%PATH%;%CD%

# Verify installation
gproc --version`

const commandCategories = ref([
  {
    name: 'Process Management',
    icon: '⚙️',
    commands: [
      {
        name: 'gproc start',
        category: 'Core',
        type: 'primary',
        description: 'Start a new process with optional configuration',
        examples: [
          { label: 'Basic start', command: 'gproc start myapp ./myapp.exe' },
          { label: 'With environment', command: 'gproc start webapp ./server.exe --env "NODE_ENV=production"' },
          { label: 'With health check', command: 'gproc start api ./api.exe --health-check "http://localhost:8080/health"' }
        ]
      },
      {
        name: 'gproc stop',
        category: 'Core',
        type: 'warning',
        description: 'Gracefully stop a running process',
        examples: [
          { label: 'Stop process', command: 'gproc stop myapp' },
          { label: 'Force stop', command: 'gproc stop myapp --force' }
        ]
      },
      {
        name: 'gproc restart',
        category: 'Core',
        type: 'info',
        description: 'Restart a process (stop then start)',
        examples: [
          { label: 'Restart process', command: 'gproc restart myapp' },
          { label: 'Restart all', command: 'gproc restart all' }
        ]
      }
    ]
  },
  {
    name: 'Monitoring',
    icon: '📊',
    commands: [
      {
        name: 'gproc list',
        category: 'Info',
        type: 'success',
        description: 'List all processes with their status',
        examples: [
          { label: 'List all', command: 'gproc list' },
          { label: 'JSON format', command: 'gproc list --json' }
        ]
      },
      {
        name: 'gproc logs',
        category: 'Info',
        type: 'info',
        description: 'View process logs with real-time tailing',
        examples: [
          { label: 'View logs', command: 'gproc logs myapp' },
          { label: 'Last 50 lines', command: 'gproc logs myapp --lines 50' },
          { label: 'Follow logs', command: 'gproc logs myapp --follow' }
        ]
      },
      {
        name: 'gproc status',
        category: 'Info',
        type: 'default',
        description: 'Show detailed status of a process',
        examples: [
          { label: 'Process status', command: 'gproc status myapp' },
          { label: 'All processes', command: 'gproc status --all' }
        ]
      }
    ]
  },
  {
    name: 'Advanced Features',
    icon: '🚀',
    commands: [
      {
        name: 'gproc cluster',
        category: 'Scale',
        type: 'primary',
        description: 'Manage process clusters for load balancing',
        examples: [
          { label: 'Start cluster', command: 'gproc cluster start webapp ./server.exe --instances 4' },
          { label: 'Scale cluster', command: 'gproc cluster scale webapp 8' }
        ]
      },
      {
        name: 'gproc schedule',
        category: 'Cron',
        type: 'warning',
        description: 'Schedule processes with cron expressions',
        examples: [
          { label: 'Daily backup', command: 'gproc schedule backup ./backup.sh --cron "0 2 * * *"' },
          { label: 'Hourly task', command: 'gproc schedule cleanup ./cleanup.sh --cron "0 * * * *"' }
        ]
      },
      {
        name: 'gproc web',
        category: 'UI',
        type: 'success',
        description: 'Start the web dashboard interface',
        examples: [
          { label: 'Default port', command: 'gproc web' },
          { label: 'Custom port', command: 'gproc web --port 3000' }
        ]
      }
    ]
  }
])

const copyToClipboard = async (text: string) => {
  try {
    await navigator.clipboard.writeText(text)
    console.log('Copied to clipboard:', text)
  } catch (err) {
    console.error('Failed to copy:', err)
  }
}

const copyInstallCommand = () => {
  copyToClipboard(installCommand)
}
</script>

<style scoped>
.cli-page {
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

.cli-sections {
  margin-bottom: 40px;
}

.install-section {
  margin-bottom: 40px;
}

.section-title {
  font-size: 1.5rem;
  font-weight: 600;
  color: var(--n-text-color);
  margin: 0 0 20px 0;
}

.code-block {
  background: var(--n-color-embedded);
  border: 1px solid var(--n-border-color);
  border-radius: 12px;
  overflow: hidden;
}

.code-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: var(--n-color);
  border-bottom: 1px solid var(--n-border-color);
}

.code-title {
  font-weight: 500;
  color: var(--n-text-color);
}

.code-content {
  padding: 16px;
  margin: 0;
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 0.875rem;
  line-height: 1.5;
  color: var(--n-text-color);
  background: transparent;
  overflow-x: auto;
}

.commands-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 32px;
  margin-bottom: 40px;
}

.command-category {
  background: var(--n-color-embedded);
  border: 1px solid var(--n-border-color);
  border-radius: 16px;
  padding: 24px;
}

.category-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 16px;
  border-bottom: 1px solid var(--n-border-color);
}

.category-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--n-text-color);
  margin: 0;
}

.category-count {
  font-size: 0.875rem;
  color: var(--n-text-color-3);
}

.commands-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.command-item {
  padding: 16px;
  background: var(--n-color);
  border: 1px solid var(--n-border-color);
  border-radius: 8px;
}

.command-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.command-name {
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--gproc-primary);
  background: rgba(16, 185, 129, 0.1);
  padding: 4px 8px;
  border-radius: 4px;
}

.command-description {
  color: var(--n-text-color-2);
  margin: 0 0 12px 0;
  font-size: 0.875rem;
}

.command-examples {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.example-item {
  background: var(--n-color-embedded);
  border-radius: 6px;
  padding: 8px;
}

.example-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}

.example-label {
  font-size: 0.75rem;
  color: var(--n-text-color-3);
  font-weight: 500;
}

.example-command {
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 0.8rem;
  color: var(--n-text-color);
  display: block;
  word-break: break-all;
}

.quick-reference {
  background: var(--n-color-embedded);
  border: 1px solid var(--n-border-color);
  border-radius: 16px;
  padding: 24px;
}

.reference-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
}

.reference-card {
  background: var(--n-color);
  border: 1px solid var(--n-border-color);
  border-radius: 8px;
  padding: 16px;
}

.reference-card h4 {
  margin: 0 0 12px 0;
  color: var(--n-text-color);
  font-weight: 600;
}

.reference-commands {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.reference-commands code {
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 0.8rem;
  color: var(--n-text-color-2);
  background: var(--n-color-embedded);
  padding: 4px 6px;
  border-radius: 3px;
}

@media (max-width: 768px) {
  .cli-page {
    padding: 16px;
  }
  
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }
  
  .commands-grid {
    grid-template-columns: 1fr;
    gap: 24px;
  }
  
  .reference-grid {
    grid-template-columns: 1fr;
  }
}
</style>