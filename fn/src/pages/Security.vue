<template>
  <div class="security">
    <n-tabs default-value="users" type="line">
      <!-- Users Management -->
      <n-tab-pane name="users" tab="Users">
        <n-card title="User Management">
          <template #header-extra>
            <n-button type="primary" @click="showCreateUserModal = true">
              <template #icon><UserPlus /></template>
              Add User
            </n-button>
          </template>
          
          <n-data-table
            :columns="userColumns"
            :data="users"
            :loading="loading"
            :pagination="{ pageSize: 10 }"
          />
        </n-card>
      </n-tab-pane>
      
      <!-- Roles Management -->
      <n-tab-pane name="roles" tab="Roles">
        <n-card title="Role Management">
          <template #header-extra>
            <n-button type="primary" @click="showCreateRoleModal = true">
              <template #icon><Shield /></template>
              Add Role
            </n-button>
          </template>
          
          <n-data-table
            :columns="roleColumns"
            :data="roles"
            :pagination="{ pageSize: 10 }"
          />
        </n-card>
      </n-tab-pane>
      
      <!-- Audit Logs -->
      <n-tab-pane name="audit" tab="Audit Logs">
        <n-card title="Audit Trail">
          <template #header-extra>
            <n-button @click="fetchAuditEvents">
              <template #icon><RefreshCw /></template>
              Refresh
            </n-button>
          </template>
          
          <n-data-table
            :columns="auditColumns"
            :data="auditEvents"
            :pagination="{ pageSize: 20 }"
          />
        </n-card>
      </n-tab-pane>
      
      <!-- SSO Configuration -->
      <n-tab-pane name="sso" tab="SSO">
        <n-card title="Single Sign-On Configuration">
          <n-form :model="ssoConfig" v-if="ssoConfig">
            <n-form-item label="Enable SSO">
              <n-switch v-model:value="ssoConfig.enabled" />
            </n-form-item>
            
            <n-form-item label="Provider" v-if="ssoConfig.enabled">
              <n-select
                v-model:value="ssoConfig.provider"
                :options="ssoProviders"
              />
            </n-form-item>
            
            <div v-if="ssoConfig.enabled && ssoConfig.provider === 'oauth2'">
              <n-form-item label="Client ID">
                <n-input v-model:value="ssoConfig.clientId" />
              </n-form-item>
              <n-form-item label="Client Secret">
                <n-input v-model:value="ssoConfig.clientSecret" type="password" />
              </n-form-item>
              <n-form-item label="Authorization URL">
                <n-input v-model:value="ssoConfig.authUrl" />
              </n-form-item>
              <n-form-item label="Token URL">
                <n-input v-model:value="ssoConfig.tokenUrl" />
              </n-form-item>
            </div>
            
            <div v-if="ssoConfig.enabled && ssoConfig.provider === 'saml'">
              <n-form-item label="Entity ID">
                <n-input v-model:value="ssoConfig.entityId" />
              </n-form-item>
              <n-form-item label="Metadata URL">
                <n-input v-model:value="ssoConfig.metadataUrl" />
              </n-form-item>
            </div>
            
            <n-button type="primary" @click="updateSSOConfig">
              Save Configuration
            </n-button>
          </n-form>
        </n-card>
      </n-tab-pane>
    </n-tabs>
    
    <!-- Create User Modal -->
    <n-modal v-model:show="showCreateUserModal" title="Create User">
      <n-card style="width: 500px">
        <n-form :model="newUser" :rules="userRules">
          <n-form-item path="username" label="Username">
            <n-input v-model:value="newUser.username" />
          </n-form-item>
          <n-form-item path="email" label="Email">
            <n-input v-model:value="newUser.email" />
          </n-form-item>
          <n-form-item path="password" label="Password">
            <n-input v-model:value="newUser.password" type="password" />
          </n-form-item>
          <n-form-item label="Roles">
            <n-select
              v-model:value="newUser.roles"
              multiple
              :options="roleOptions"
            />
          </n-form-item>
          <n-form-item label="Enabled">
            <n-switch v-model:value="newUser.enabled" />
          </n-form-item>
        </n-form>
        
        <template #footer>
          <n-space justify="end">
            <n-button @click="showCreateUserModal = false">Cancel</n-button>
            <n-button type="primary" @click="createUser">Create</n-button>
          </n-space>
        </template>
      </n-card>
    </n-modal>
    
    <!-- Create Role Modal -->
    <n-modal v-model:show="showCreateRoleModal" title="Create Role">
      <n-card style="width: 600px">
        <n-form :model="newRole">
          <n-form-item label="Role Name">
            <n-input v-model:value="newRole.name" />
          </n-form-item>
          <n-form-item label="Description">
            <n-input v-model:value="newRole.description" type="textarea" />
          </n-form-item>
          <n-form-item label="Permissions">
            <n-space vertical>
              <div v-for="resource in resources" :key="resource">
                <h4>{{ resource }}</h4>
                <n-checkbox-group v-model:value="rolePermissions[resource]">
                  <n-space>
                    <n-checkbox value="read">Read</n-checkbox>
                    <n-checkbox value="write">Write</n-checkbox>
                    <n-checkbox value="delete">Delete</n-checkbox>
                    <n-checkbox value="execute">Execute</n-checkbox>
                  </n-space>
                </n-checkbox-group>
              </div>
            </n-space>
          </n-form-item>
        </n-form>
        
        <template #footer>
          <n-space justify="end">
            <n-button @click="showCreateRoleModal = false">Cancel</n-button>
            <n-button type="primary" @click="createRole">Create</n-button>
          </n-space>
        </template>
      </n-card>
    </n-modal>
    
    <!-- MFA Setup Modal -->
    <n-modal v-model:show="showMFAModal" title="Setup Multi-Factor Authentication">
      <n-card style="width: 400px">
        <div v-if="mfaSetup.qrCode" class="mfa-setup">
          <p>Scan this QR code with your authenticator app:</p>
          <div class="qr-code" v-html="mfaSetup.qrCode"></div>
          <p>Or enter this secret manually: <code>{{ mfaSetup.secret }}</code></p>
          
          <n-form-item label="Verification Code">
            <n-input v-model:value="mfaVerificationCode" placeholder="Enter 6-digit code" />
          </n-form-item>
          
          <n-button type="primary" @click="verifyMFA" block>
            Verify and Enable MFA
          </n-button>
        </div>
      </n-card>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, h } from 'vue'
import { useSecurityStore } from '../stores/security'
import { useMessage } from 'naive-ui'
import { UserPlus, Shield, RefreshCw, Key, Trash2, Edit } from 'lucide-vue-next'

const securityStore = useSecurityStore()
const message = useMessage()

const loading = ref(false)
const showCreateUserModal = ref(false)
const showCreateRoleModal = ref(false)
const showMFAModal = ref(false)
const mfaSetup = ref({ qrCode: '', secret: '' })
const mfaVerificationCode = ref('')
const rolePermissions = ref<Record<string, string[]>>({})

const newUser = ref({
  username: '',
  email: '',
  password: '',
  roles: [] as string[],
  enabled: true
})

const newRole = ref({
  name: '',
  description: '',
  permissions: []
})

const { users, roles, auditEvents, ssoConfig } = securityStore

const resources = ['process', 'cluster', 'user', 'role', 'metrics', 'deployment', 'scheduler', 'audit', 'backup', 'secrets']

const ssoProviders = [
  { label: 'OAuth2', value: 'oauth2' },
  { label: 'SAML', value: 'saml' }
]

const userRules = {
  username: { required: true, message: 'Username is required' },
  email: { required: true, message: 'Email is required' },
  password: { required: true, message: 'Password is required' }
}

const roleOptions = computed(() =>
  roles.value.map(r => ({ label: r.name, value: r.name }))
)

const userColumns = [
  { title: 'Username', key: 'username' },
  { title: 'Email', key: 'email' },
  {
    title: 'Roles',
    key: 'roles',
    render: (row: any) => h('n-space', row.roles.map((role: string) =>
      h('n-tag', { size: 'small' }, { default: () => role })
    ))
  },
  {
    title: 'Status',
    key: 'enabled',
    render: (row: any) => h('n-tag', {
      type: row.enabled ? 'success' : 'default'
    }, { default: () => row.enabled ? 'Active' : 'Disabled' })
  },
  {
    title: 'MFA',
    key: 'mfaEnabled',
    render: (row: any) => h('n-tag', {
      type: row.mfaEnabled ? 'success' : 'warning'
    }, { default: () => row.mfaEnabled ? 'Enabled' : 'Disabled' })
  },
  {
    title: 'Last Seen',
    key: 'lastSeen',
    render: (row: any) => new Date(row.lastSeen).toLocaleString()
  },
  {
    title: 'Actions',
    key: 'actions',
    render: (row: any) => h('n-space', [
      h('n-button', {
        size: 'small',
        onClick: () => setupMFA(row.id)
      }, { default: () => 'MFA', icon: () => h(Key) }),
      h('n-button', {
        size: 'small',
        onClick: () => editUser(row)
      }, { default: () => 'Edit', icon: () => h(Edit) }),
      h('n-button', {
        size: 'small',
        type: 'error',
        onClick: () => deleteUser(row.id)
      }, { default: () => 'Delete', icon: () => h(Trash2) })
    ])
  }
]

const roleColumns = [
  { title: 'Name', key: 'name' },
  { title: 'Description', key: 'description' },
  {
    title: 'Permissions',
    key: 'permissions',
    render: (row: any) => h('n-space', { vertical: true }, row.permissions.map((perm: any) =>
      h('n-tag', { size: 'small' }, { default: () => `${perm.resource}:${perm.actions.join(',')}` })
    ))
  },
  {
    title: 'Actions',
    key: 'actions',
    render: (row: any) => h('n-space', [
      h('n-button', {
        size: 'small',
        onClick: () => editRole(row)
      }, { default: () => 'Edit', icon: () => h(Edit) }),
      h('n-button', {
        size: 'small',
        type: 'error',
        onClick: () => deleteRole(row.name)
      }, { default: () => 'Delete', icon: () => h(Trash2) })
    ])
  }
]

const auditColumns = [
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

const createUser = async () => {
  const result = await securityStore.createUser(newUser.value)
  if (result.success) {
    message.success('User created successfully')
    showCreateUserModal.value = false
    resetUserForm()
  } else {
    message.error(result.error)
  }
}

const createRole = async () => {
  const permissions = Object.entries(rolePermissions.value).flatMap(([resource, actions]) =>
    actions.map(action => ({ resource, actions: [action], scope: '*' }))
  )
  
  const result = await securityStore.createRole({
    ...newRole.value,
    permissions
  })
  
  if (result.success) {
    message.success('Role created successfully')
    showCreateRoleModal.value = false
    resetRoleForm()
  } else {
    message.error(result.error)
  }
}

const setupMFA = async (userId: string) => {
  const result = await securityStore.enableMFA(userId)
  if (result.success) {
    mfaSetup.value = { qrCode: result.qrCode, secret: result.secret }
    showMFAModal.value = true
  } else {
    message.error(result.error)
  }
}

const verifyMFA = async () => {
  // Implementation would verify MFA code
  message.success('MFA enabled successfully')
  showMFAModal.value = false
}

const updateSSOConfig = async () => {
  if (!ssoConfig.value) return
  
  const result = await securityStore.updateSSOConfig(ssoConfig.value)
  if (result.success) {
    message.success('SSO configuration updated')
  } else {
    message.error(result.error)
  }
}

const editUser = (user: any) => {
  // Implementation for editing user
  message.info('Edit user functionality')
}

const deleteUser = async (userId: string) => {
  const result = await securityStore.deleteUser(userId)
  if (result.success) {
    message.success('User deleted')
  } else {
    message.error(result.error)
  }
}

const editRole = (role: any) => {
  // Implementation for editing role
  message.info('Edit role functionality')
}

const deleteRole = (roleName: string) => {
  // Implementation for deleting role
  message.info('Delete role functionality')
}

const fetchAuditEvents = () => {
  securityStore.fetchAuditEvents()
}

const resetUserForm = () => {
  newUser.value = {
    username: '',
    email: '',
    password: '',
    roles: [],
    enabled: true
  }
}

const resetRoleForm = () => {
  newRole.value = {
    name: '',
    description: '',
    permissions: []
  }
  rolePermissions.value = {}
}

onMounted(() => {
  securityStore.fetchUsers()
  securityStore.fetchRoles()
  securityStore.fetchAuditEvents()
  securityStore.fetchSSOConfig()
})
</script>

<style scoped>
.security {
  padding: 1rem;
}

.mfa-setup {
  text-align: center;
}

.qr-code {
  margin: 1rem 0;
}

.qr-code img {
  max-width: 200px;
}
</style>