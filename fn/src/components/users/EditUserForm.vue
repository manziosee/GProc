<template>
  <div class="edit-user-form">
    <n-form ref="formRef" :model="formData" :rules="rules" label-placement="top">
      <n-form-item label="Full Name" path="fullName">
        <n-input 
          v-model:value="formData.fullName" 
          placeholder="Enter full name"
        />
      </n-form-item>
      
      <n-form-item label="Username" path="username">
        <n-input 
          v-model:value="formData.username" 
          placeholder="Enter username"
        />
      </n-form-item>
      
      <n-form-item label="Email" path="email">
        <n-input 
          v-model:value="formData.email" 
          placeholder="Enter email address"
          type="email"
        />
      </n-form-item>
      
      <n-form-item label="Role" path="role">
        <n-select
          v-model:value="formData.role"
          :options="roleOptions"
          placeholder="Select role"
        />
      </n-form-item>
      
      <n-form-item label="Status" path="status">
        <n-select
          v-model:value="formData.status"
          :options="statusOptions"
          placeholder="Select status"
        />
      </n-form-item>
      
      <div class="form-actions">
        <n-button @click="$emit('cancel')">Cancel</n-button>
        <n-button type="primary" @click="handleSubmit">Update User</n-button>
      </div>
    </n-form>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { NForm, NFormItem, NInput, NSelect, NButton } from 'naive-ui'

interface User {
  id: string
  username: string
  email: string
  fullName: string
  role: string
  status: string
}

const props = defineProps<{
  user: User
}>()

const emit = defineEmits<{
  submit: [data: any]
  cancel: []
}>()

const formRef = ref()

const formData = reactive({
  fullName: '',
  username: '',
  email: '',
  role: 'user',
  status: 'active'
})

const roleOptions = [
  { label: 'Administrator', value: 'admin' },
  { label: 'User', value: 'user' },
  { label: 'Viewer', value: 'viewer' }
]

const statusOptions = [
  { label: 'Active', value: 'active' },
  { label: 'Inactive', value: 'inactive' },
  { label: 'Suspended', value: 'suspended' }
]

const rules = {
  fullName: {
    required: true,
    message: 'Full name is required'
  },
  username: {
    required: true,
    message: 'Username is required'
  },
  email: {
    required: true,
    type: 'email',
    message: 'Valid email is required'
  },
  role: {
    required: true,
    message: 'Role is required'
  },
  status: {
    required: true,
    message: 'Status is required'
  }
}

const handleSubmit = async () => {
  try {
    await formRef.value?.validate()
    emit('submit', { ...formData })
  } catch (error) {
    console.error('Validation failed:', error)
  }
}

onMounted(() => {
  // Initialize form with user data
  formData.fullName = props.user.fullName
  formData.username = props.user.username
  formData.email = props.user.email
  formData.role = props.user.role
  formData.status = props.user.status
})
</script>

<style scoped>
.edit-user-form {
  padding: 20px;
  max-width: 500px;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 24px;
}
</style>