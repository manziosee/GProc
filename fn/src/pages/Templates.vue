<template>
  <div class="templates">
    <n-card title="Language Templates">
      <template #header-extra>
        <n-button type="primary" @click="showCreateModal = true">
          <template #icon><Plus /></template>
          Create Template
        </n-button>
      </template>
      
      <n-data-table
        :columns="columns"
        :data="templates"
        :loading="loading"
        :pagination="{ pageSize: 10 }"
      />
    </n-card>
    
    <!-- Create Template Modal -->
    <n-modal v-model:show="showCreateModal" title="Create Template">
      <n-card style="width: 600px">
        <n-form :model="newTemplate">
          <n-form-item label="Template Name">
            <n-input v-model:value="newTemplate.name" placeholder="nodejs-api" />
          </n-form-item>
          <n-form-item label="Language">
            <n-select v-model:value="newTemplate.language" :options="languageOptions" />
          </n-form-item>
          <n-form-item label="Command">
            <n-input v-model:value="newTemplate.command" placeholder="node server.js" />
          </n-form-item>
          <n-form-item label="Arguments">
            <n-dynamic-input v-model:value="newTemplate.args" />
          </n-form-item>
        </n-form>
        
        <template #footer>
          <n-space justify="end">
            <n-button @click="showCreateModal = false">Cancel</n-button>
            <n-button type="primary" @click="createTemplate">Create</n-button>
          </n-space>
        </template>
      </n-card>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, h } from 'vue'
import { useDeploymentStore } from '../stores/deployment'
import { useMessage } from 'naive-ui'
import { Plus, Eye, Trash2 } from 'lucide-vue-next'

const deploymentStore = useDeploymentStore()
const message = useMessage()

const loading = ref(false)
const showCreateModal = ref(false)

const newTemplate = ref({
  name: '',
  language: '',
  command: '',
  args: [] as string[]
})

const { templates } = deploymentStore

const languageOptions = [
  { label: 'Node.js', value: 'nodejs' },
  { label: 'Python', value: 'python' },
  { label: 'Java', value: 'java' },
  { label: 'Go', value: 'go' }
]

const columns = [
  { title: 'Name', key: 'name' },
  { title: 'Language', key: 'language' },
  { title: 'Command', key: 'command' },
  {
    title: 'Actions',
    key: 'actions',
    render: (row: any) => h('n-space', [
      h('n-button', {
        size: 'small',
        onClick: () => viewTemplate(row)
      }, { default: () => 'View', icon: () => h(Eye) }),
      h('n-button', {
        size: 'small',
        type: 'error',
        onClick: () => deleteTemplate(row.id)
      }, { default: () => 'Delete', icon: () => h(Trash2) })
    ])
  }
]

const createTemplate = async () => {
  const result = await deploymentStore.createTemplate(newTemplate.value)
  if (result.success) {
    message.success('Template created')
    showCreateModal.value = false
  } else {
    message.error(result.error)
  }
}

const viewTemplate = (template: any) => {
  message.info(`Viewing template: ${template.name}`)
}

const deleteTemplate = (templateId: string) => {
  message.info(`Deleting template: ${templateId}`)
}
</script>

<style scoped>
.templates {
  padding: 1rem;
}
</style>