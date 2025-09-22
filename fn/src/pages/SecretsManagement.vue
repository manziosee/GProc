<template>
  <div class="secrets-management">
    <n-card title="Secrets Management">
      <template #header-extra>
        <n-space>
          <n-button @click="refreshSecrets">
            <template #icon><RefreshCw /></template>
            Refresh
          </n-button>
          <n-button type="primary" @click="showCreateModal = true">
            <template #icon><Plus /></template>
            Add Secret
          </n-button>
        </n-space>
      </template>
      
      <!-- Secrets Table -->
      <n-data-table
        :columns="columns"
        :data="secrets"
        :loading="loading"
        :pagination="{ pageSize: 10 }"
      />
    </n-card>
    
    <!-- Create Secret Modal -->
    <n-modal v-model:show="showCreateModal" title="Create Secret">
      <n-card style="width: 500px">
        <n-form :model="newSecret" :rules="secretRules">
          <n-form-item path="key" label="Secret Key">
            <n-input v-model:value="newSecret.key" placeholder="DATABASE_PASSWORD" />
          </n-form-item>
          
          <n-form-item path="value" label="Secret Value">
            <n-input
              v-model:value="newSecret.value"
              type="password"
              placeholder="Enter secret value"
              show-password-on="click"
            />
          </n-form-item>
          
          <n-form-item label="Description">
            <n-input
              v-model:value="newSecret.description"
              type="textarea"
              placeholder="Optional description"
            />
          </n-form-item>
          
          <n-form-item label="Tags">
            <n-dynamic-tags v-model:value="newSecret.tags" />
          </n-form-item>
          
          <n-form-item label="Expiration">
            <n-date-picker
              v-model:value="newSecret.expiresAt"
              type="datetime"
              placeholder="Optional expiration date"
            />
          </n-form-item>
        </n-form>
        
        <template #footer>
          <n-space justify="end">
            <n-button @click="showCreateModal = false">Cancel</n-button>
            <n-button type="primary" @click="createSecret">Create Secret</n-button>
          </n-space>
        </template>
      </n-card>
    </n-modal>
    
    <!-- View Secret Modal -->
    <n-modal v-model:show="showViewModal" title="Secret Details">
      <n-card style="width: 500px" v-if="selectedSecret">
        <n-descriptions :column="1" bordered>
          <n-descriptions-item label="Key">{{ selectedSecret.key }}</n-descriptions-item>
          <n-descriptions-item label="Value">
            <n-space align="center">
              <span v-if="!showSecretValue">••••••••••••</span>
              <span v-else>{{ selectedSecret.value }}</span>
              <n-button size="small" @click="toggleSecretVisibility">
                <template #icon>
                  <component :is="showSecretValue ? EyeOff : Eye" />
                </template>
              </n-button>
            </n-space>
          </n-descriptions-item>
          <n-descriptions-item label="Description">
            {{ selectedSecret.description || 'No description' }}
          </n-descriptions-item>
          <n-descriptions-item label="Tags">
            <n-space>
              <n-tag v-for="tag in selectedSecret.tags" :key="tag" size="small">
                {{ tag }}
              </n-tag>
            </n-space>
          </n-descriptions-item>
          <n-descriptions-item label="Created">
            {{ new Date(selectedSecret.createdAt).toLocaleString() }}
          </n-descriptions-item>
          <n-descriptions-item label="Last Updated">
            {{ new Date(selectedSecret.updatedAt).toLocaleString() }}
          </n-descriptions-item>
          <n-descriptions-item label="Expires">
            {{ selectedSecret.expiresAt ? new Date(selectedSecret.expiresAt).toLocaleString() : 'Never' }}
          </n-descriptions-item>
          <n-descriptions-item label="Version">{{ selectedSecret.version }}</n-descriptions-item>
        </n-descriptions>
      </n-card>
    </n-modal>
    
    <!-- Vault Configuration -->
    <n-card title="Vault Configuration" class="mt-4">
      <n-form :model="vaultConfig">
        <n-grid :cols="2" :x-gap="16">
          <n-grid-item>
            <n-form-item label="Vault Address">
              <n-input v-model:value="vaultConfig.address" placeholder="https://vault.example.com" />
            </n-form-item>
          </n-grid-item>
          <n-grid-item>
            <n-form-item label="Auth Method">
              <n-select
                v-model:value="vaultConfig.authMethod"
                :options="authMethods"
              />
            </n-form-item>
          </n-grid-item>
          <n-grid-item>
            <n-form-item label="Mount Path">
              <n-input v-model:value="vaultConfig.mountPath" placeholder="secret/" />
            </n-form-item>
          </n-grid-item>
          <n-grid-item>
            <n-form-item label="Namespace">
              <n-input v-model:value="vaultConfig.namespace" placeholder="Optional namespace" />
            </n-form-item>
          </n-grid-item>
        </n-grid>
        
        <n-space>
          <n-button @click="testVaultConnection" :loading="testingConnection">
            Test Connection
          </n-button>
          <n-button type="primary" @click="saveVaultConfig">
            Save Configuration
          </n-button>
        </n-space>
      </n-form>
    </n-card>
    
    <!-- AWS KMS Configuration -->
    <n-card title="AWS KMS Configuration" class="mt-4">
      <n-form :model="kmsConfig">
        <n-grid :cols="2" :x-gap="16">
          <n-grid-item>
            <n-form-item label="AWS Region">
              <n-select
                v-model:value="kmsConfig.region"
                :options="awsRegions"
              />
            </n-form-item>
          </n-grid-item>
          <n-grid-item>
            <n-form-item label="KMS Key ID">
              <n-input v-model:value="kmsConfig.keyId" placeholder="arn:aws:kms:..." />
            </n-form-item>
          </n-grid-item>
          <n-grid-item>
            <n-form-item label="Access Key ID">
              <n-input v-model:value="kmsConfig.accessKeyId" />
            </n-form-item>
          </n-grid-item>
          <n-grid-item>
            <n-form-item label="Secret Access Key">
              <n-input v-model:value="kmsConfig.secretAccessKey" type="password" />
            </n-form-item>
          </n-grid-item>
        </n-grid>
        
        <n-space>
          <n-button @click="testKMSConnection" :loading="testingKMS">
            Test Connection
          </n-button>
          <n-button type="primary" @click="saveKMSConfig">
            Save Configuration
          </n-button>
        </n-space>
      </n-form>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import { useMessage } from 'naive-ui'
import { RefreshCw, Plus, Eye, EyeOff, Key, Trash2, Edit } from 'lucide-vue-next'

const message = useMessage()

const loading = ref(false)
const showCreateModal = ref(false)
const showViewModal = ref(false)
const selectedSecret = ref(null)
const showSecretValue = ref(false)
const testingConnection = ref(false)
const testingKMS = ref(false)

const secrets = ref([
  {
    id: '1',
    key: 'DATABASE_PASSWORD',
    value: 'super-secret-password',
    description: 'Main database password',
    tags: ['database', 'production'],
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z',
    expiresAt: null,
    version: 1
  }
])

const newSecret = ref({
  key: '',
  value: '',
  description: '',
  tags: [] as string[],
  expiresAt: null
})

const vaultConfig = ref({
  address: '',
  authMethod: 'token',
  mountPath: 'secret/',
  namespace: ''
})

const kmsConfig = ref({
  region: 'us-east-1',
  keyId: '',
  accessKeyId: '',
  secretAccessKey: ''
})

const secretRules = {
  key: { required: true, message: 'Secret key is required' },
  value: { required: true, message: 'Secret value is required' }
}

const authMethods = [
  { label: 'Token', value: 'token' },
  { label: 'AWS IAM', value: 'aws' },
  { label: 'Kubernetes', value: 'kubernetes' },
  { label: 'LDAP', value: 'ldap' }
]

const awsRegions = [
  { label: 'US East 1', value: 'us-east-1' },
  { label: 'US West 2', value: 'us-west-2' },
  { label: 'EU West 1', value: 'eu-west-1' },
  { label: 'AP Southeast 1', value: 'ap-southeast-1' }
]

const columns = [
  { title: 'Key', key: 'key' },
  {
    title: 'Description',
    key: 'description',
    render: (row: any) => row.description || 'No description'
  },
  {
    title: 'Tags',
    key: 'tags',
    render: (row: any) => h('n-space', row.tags.map((tag: string) =>
      h('n-tag', { size: 'small' }, { default: () => tag })
    ))
  },
  {
    title: 'Expires',
    key: 'expiresAt',
    render: (row: any) => row.expiresAt ? new Date(row.expiresAt).toLocaleDateString() : 'Never'
  },
  {
    title: 'Version',
    key: 'version'
  },
  {
    title: 'Updated',
    key: 'updatedAt',
    render: (row: any) => new Date(row.updatedAt).toLocaleDateString()
  },
  {
    title: 'Actions',
    key: 'actions',
    render: (row: any) => h('n-space', [
      h('n-button', {
        size: 'small',
        onClick: () => viewSecret(row)
      }, { default: () => 'View', icon: () => h(Eye) }),
      h('n-button', {
        size: 'small',
        onClick: () => editSecret(row)
      }, { default: () => 'Edit', icon: () => h(Edit) }),
      h('n-button', {
        size: 'small',
        type: 'error',
        onClick: () => deleteSecret(row.id)
      }, { default: () => 'Delete', icon: () => h(Trash2) })
    ])
  }
]

const refreshSecrets = async () => {
  loading.value = true
  try {
    // Implementation would fetch secrets from API
    await new Promise(resolve => setTimeout(resolve, 1000))
    message.success('Secrets refreshed')
  } catch (error) {
    message.error('Failed to refresh secrets')
  } finally {
    loading.value = false
  }
}

const createSecret = async () => {
  try {
    // Implementation would create secret via API
    const secret = {
      id: Date.now().toString(),
      ...newSecret.value,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      version: 1
    }
    
    secrets.value.push(secret)
    message.success('Secret created successfully')
    showCreateModal.value = false
    resetForm()
  } catch (error) {
    message.error('Failed to create secret')
  }
}

const viewSecret = (secret: any) => {
  selectedSecret.value = secret
  showSecretValue.value = false
  showViewModal.value = true
}

const editSecret = (secret: any) => {
  message.info('Edit secret functionality would be implemented')
}

const deleteSecret = async (secretId: string) => {
  try {
    secrets.value = secrets.value.filter(s => s.id !== secretId)
    message.success('Secret deleted successfully')
  } catch (error) {
    message.error('Failed to delete secret')
  }
}

const toggleSecretVisibility = () => {
  showSecretValue.value = !showSecretValue.value
}

const testVaultConnection = async () => {
  testingConnection.value = true
  try {
    // Implementation would test Vault connection
    await new Promise(resolve => setTimeout(resolve, 2000))
    message.success('Vault connection successful')
  } catch (error) {
    message.error('Vault connection failed')
  } finally {
    testingConnection.value = false
  }
}

const saveVaultConfig = async () => {
  try {
    // Implementation would save Vault config
    message.success('Vault configuration saved')
  } catch (error) {
    message.error('Failed to save Vault configuration')
  }
}

const testKMSConnection = async () => {
  testingKMS.value = true
  try {
    // Implementation would test KMS connection
    await new Promise(resolve => setTimeout(resolve, 2000))
    message.success('KMS connection successful')
  } catch (error) {
    message.error('KMS connection failed')
  } finally {
    testingKMS.value = false
  }
}

const saveKMSConfig = async () => {
  try {
    // Implementation would save KMS config
    message.success('KMS configuration saved')
  } catch (error) {
    message.error('Failed to save KMS configuration')
  }
}

const resetForm = () => {
  newSecret.value = {
    key: '',
    value: '',
    description: '',
    tags: [],
    expiresAt: null
  }
}

onMounted(() => {
  refreshSecrets()
})
</script>

<style scoped>
.secrets-management {
  padding: 1rem;
}

.mt-4 {
  margin-top: 2rem;
}
</style>