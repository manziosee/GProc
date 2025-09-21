<template>
  <div class="user-management">
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">User Management</h1>
        <p class="page-subtitle">Manage users and their permissions</p>
      </div>
      <n-button type="primary" size="large" @click="showAddModal = true">
        <template #icon><span>👥</span></template>
        Add User
      </n-button>
    </div>

    <div class="user-stats">
      <div class="stat-card active">
        <div class="stat-icon">🟢</div>
        <div class="stat-info">
          <div class="stat-value">{{ activeUsers }}</div>
          <div class="stat-label">Active Users</div>
        </div>
      </div>
      <div class="stat-card admin">
        <div class="stat-icon">🔑</div>
        <div class="stat-info">
          <div class="stat-value">{{ adminUsers }}</div>
          <div class="stat-label">Administrators</div>
        </div>
      </div>
      <div class="stat-card total">
        <div class="stat-icon">📊</div>
        <div class="stat-info">
          <div class="stat-value">{{ users.length }}</div>
          <div class="stat-label">Total Users</div>
        </div>
      </div>
    </div>

    <div class="users-container">
      <div class="users-header">
        <n-input placeholder="Search users..." style="width: 300px">
          <template #prefix>🔍</template>
        </n-input>
        <n-select v-model:value="roleFilter" :options="roleOptions" placeholder="Filter by role" style="width: 150px" />
      </div>

      <div class="users-grid">
        <div v-for="user in filteredUsers" :key="user.id" class="user-card">
          <div class="user-avatar">
            <div class="avatar-circle" :class="user.role">
              {{ user.name.charAt(0).toUpperCase() }}
            </div>
            <div class="status-indicator" :class="user.status"></div>
          </div>
          
          <div class="user-info">
            <h3 class="user-name">{{ user.name }}</h3>
            <p class="user-email">{{ user.email }}</p>
            <div class="user-meta">
              <n-tag :type="getRoleTagType(user.role)" size="small">{{ user.role }}</n-tag>
              <span class="last-seen">Last seen {{ user.lastSeen }}</span>
            </div>
          </div>
          
          <div class="user-actions">
            <n-dropdown :options="getUserActions(user)" @select="handleUserAction">
              <n-button size="small" quaternary>
                <template #icon><span>⋮</span></template>
              </n-button>
            </n-dropdown>
          </div>
        </div>
      </div>
    </div>

    <n-modal v-model:show="showAddModal" preset="card" title="Add New User" style="width: 500px">
      <AddUserForm @submit="handleAddUser" @cancel="showAddModal = false" />
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { NButton, NInput, NSelect, NTag, NDropdown, NModal } from 'naive-ui'
import AddUserForm from '../components/users/AddUserForm.vue'

const showAddModal = ref(false)
const roleFilter = ref('')

const users = ref([
  {
    id: '1',
    name: 'John Doe',
    email: 'john@company.com',
    role: 'admin',
    status: 'online',
    lastSeen: '2 minutes ago'
  },
  {
    id: '2',
    name: 'Jane Smith',
    email: 'jane@company.com',
    role: 'user',
    status: 'online',
    lastSeen: '5 minutes ago'
  },
  {
    id: '3',
    name: 'Mike Johnson',
    email: 'mike@company.com',
    role: 'viewer',
    status: 'offline',
    lastSeen: '2 hours ago'
  },
  {
    id: '4',
    name: 'Sarah Wilson',
    email: 'sarah@company.com',
    role: 'admin',
    status: 'online',
    lastSeen: '1 minute ago'
  }
])

const roleOptions = [
  { label: 'All Roles', value: '' },
  { label: 'Administrator', value: 'admin' },
  { label: 'User', value: 'user' },
  { label: 'Viewer', value: 'viewer' }
]

const activeUsers = computed(() => users.value.filter(u => u.status === 'online').length)
const adminUsers = computed(() => users.value.filter(u => u.role === 'admin').length)

const filteredUsers = computed(() => {
  return users.value.filter(user => {
    if (roleFilter.value && user.role !== roleFilter.value) return false
    return true
  })
})

const getRoleTagType = (role: string) => {
  switch (role) {
    case 'admin': return 'error'
    case 'user': return 'primary'
    case 'viewer': return 'default'
    default: return 'default'
  }
}

const getUserActions = (user: any) => [
  { label: 'Edit User', key: 'edit' },
  { label: 'Reset Password', key: 'reset' },
  { label: 'Suspend User', key: 'suspend' },
  { label: 'Delete User', key: 'delete', props: { style: 'color: red' } }
]

const handleUserAction = (key: string) => {
  console.log('User action:', key)
}

const handleAddUser = (data: any) => {
  console.log('Add user:', data)
  showAddModal.value = false
}
</script>

<style scoped>
.user-management {
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

.user-stats {
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

.stat-card.admin {
  --accent-color: #ef4444;
  --accent-color-light: #f87171;
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

.users-container {
  background: var(--n-color-embedded);
  border: 1px solid var(--n-border-color);
  border-radius: 16px;
  padding: 24px;
}

.users-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  gap: 16px;
}

.users-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 20px;
}

.user-card {
  background: var(--n-color);
  border: 1px solid var(--n-border-color);
  border-radius: 12px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
  transition: all 0.3s ease;
}

.user-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.user-avatar {
  position: relative;
  flex-shrink: 0;
}

.avatar-circle {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 1.25rem;
  color: white;
}

.avatar-circle.admin {
  background: linear-gradient(135deg, #ef4444, #dc2626);
}

.avatar-circle.user {
  background: linear-gradient(135deg, #3b82f6, #2563eb);
}

.avatar-circle.viewer {
  background: linear-gradient(135deg, #6b7280, #4b5563);
}

.status-indicator {
  position: absolute;
  bottom: 2px;
  right: 2px;
  width: 12px;
  height: 12px;
  border-radius: 50%;
  border: 2px solid var(--n-color);
}

.status-indicator.online {
  background: #10b981;
}

.status-indicator.offline {
  background: #6b7280;
}

.user-info {
  flex: 1;
  min-width: 0;
}

.user-name {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--n-text-color);
  margin: 0 0 4px 0;
}

.user-email {
  font-size: 0.875rem;
  color: var(--n-text-color-2);
  margin: 0 0 8px 0;
}

.user-meta {
  display: flex;
  align-items: center;
  gap: 12px;
}

.last-seen {
  font-size: 0.75rem;
  color: var(--n-text-color-3);
}

.user-actions {
  flex-shrink: 0;
}

@media (max-width: 768px) {
  .user-management {
    padding: 16px;
  }
  
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }
  
  .users-header {
    flex-direction: column;
    align-items: stretch;
  }
  
  .users-grid {
    grid-template-columns: 1fr;
  }
}
</style>