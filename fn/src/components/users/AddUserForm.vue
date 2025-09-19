<template>
  <div class="add-user-form">
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
      
      <n-form-item label="Password" path="password">
        <n-input 
          v-model:value="formData.password" 
          placeholder="Enter password"
          type="password"
        />
      </n-form-item>
      
      <n-form-item label="Confirm Password" path="confirmPassword">
        <n-input 
          v-model:value="formData.confirmPassword" 
          placeholder="Confirm password"
          type="password"
        />
      </n-form-item>
      
      <n-form-item>
        <n-checkbox v-model:checked="formData.sendWelcomeEmail">
          Send welcome email to user
        </n-checkbox>
      </n-form-item>
      
      <div class="form-actions">
        <n-button @click="$emit('cancel')">Cancel</n-button>
        <n-button type="primary" @click="handleSubmit">Add User</n-button>
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

const formData = reactive({
  fullName: '',
  username: '',
  email: '',
  role: 'user',
  password: '',
  confirmPassword: '',
  sendWelcomeEmail: true
})

const roleOptions = [
  { label: 'Administrator', value: 'admin' },
  { label: 'User', value: 'user' },
  { label: 'Viewer', value: 'viewer' }
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
  password: {
    required: true,
    min: 8,
    message: 'Password must be at least 8 characters'
  },
  confirmPassword: [
    {
      required: true,
      message: 'Please confirm password'
    },
    {
      validator: (rule: any, value: string) => value === formData.password,
      message: 'Passwords do not match'
    }
  ]
}

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
.add-user-form {
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