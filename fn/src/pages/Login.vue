<template>
  <div class="login">
    <div class="panel left">
      <div class="brand">
        <div class="logo">🚀</div>
        <div class="name">GProc</div>
      </div>
      <h1>Operate your services with confidence</h1>
      <p>Beautiful PM2‑style console, role‑based access, and real‑time insights. Everything you need to run production.</p>
      <ul class="highlights">
        <li>📊 Live dashboard with KPIs and process table</li>
        <li>🔐 RBAC with audit logging</li>
        <li>📈 Prometheus metrics and health checks</li>
        <li>🚦 Zero‑downtime restarts</li>
      </ul>
      <div class="backdrop" aria-hidden="true"></div>
    </div>

    <div class="panel right">
      <div class="card">
        <div class="card-header">
          <div class="badge">{{ isRegisterMode ? 'Get started' : 'Welcome back' }}</div>
          <h2>{{ isRegisterMode ? 'Create your account' : 'Sign in to your console' }}</h2>
          <p>{{ isRegisterMode ? 'Join GProc to manage your processes' : 'Use your GProc credentials to continue' }}</p>
        </div>

        <n-form ref="formRef" :model="form" :rules="rules" @submit.prevent="isRegisterMode ? handleRegister : handleLogin">
          <n-form-item path="username" label="Username">
            <n-input v-model:value="form.username" placeholder="Enter username" size="large">
              <template #prefix>👤</template>
            </n-input>
          </n-form-item>

          <n-form-item v-if="isRegisterMode" path="email" label="Email">
            <n-input v-model:value="form.email" placeholder="Enter email" size="large">
              <template #prefix>📧</template>
            </n-input>
          </n-form-item>

          <n-form-item path="password" label="Password">
            <n-input v-model:value="form.password" type="password" placeholder="Enter password" size="large">
              <template #prefix>🔒</template>
            </n-input>
          </n-form-item>

          <n-form-item v-if="!isRegisterMode && showMFA" path="mfaCode" label="MFA Code">
            <n-input v-model:value="form.mfaCode" placeholder="Enter 6‑digit code" size="large">
              <template #prefix>🧩</template>
            </n-input>
          </n-form-item>

          <div class="row" v-if="!isRegisterMode">
            <n-checkbox v-model:checked="remember">Remember me</n-checkbox>
            <a class="link" href="#">Forgot password?</a>
          </div>

          <n-button type="primary" block size="large" :loading="loading" @click="isRegisterMode ? handleRegister : handleLogin">
            {{ isRegisterMode ? 'Create Account' : (showMFA ? 'Verify MFA' : 'Sign in') }}
          </n-button>
        </n-form>

        <div class="toggle">
          <span>{{ isRegisterMode ? 'Already have an account?' : "Don't have an account?" }}</span>
          <a class="link" @click="isRegisterMode = !isRegisterMode">
            {{ isRegisterMode ? 'Sign in' : 'Register' }}
          </a>
        </div>

        <div class="sso" v-if="ssoEnabled">
          <n-divider>Or continue with</n-divider>
          <n-button block strong secondary size="large" @click="handleSSOLogin">Single Sign‑On</n-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useMessage, NForm, NFormItem, NInput, NButton, NDivider, NCheckbox } from 'naive-ui'

const router = useRouter()
const authStore = useAuthStore()
const message = useMessage()

const loading = ref(false)
const showMFA = ref(false)
const ssoEnabled = ref(true)
const isRegisterMode = ref(false)

// match template ref
const formRef = ref()
const remember = ref(true)

const form = reactive({
  username: '',
  email: '',
  password: '',
  mfaCode: ''
})

const rules = reactive({
  username: { required: true, message: 'Username is required' },
  email: { required: () => isRegisterMode.value, message: 'Email is required' },
  password: { required: true, message: 'Password is required' },
  mfaCode: { required: () => showMFA.value, message: 'MFA code is required' }
})

const handleRegister = async () => {
  loading.value = true
  try {
    const result = await authStore.register(form.username, form.password, form.email)
    
    if (result.success) {
      message.success('Registration successful! Please sign in.')
      isRegisterMode.value = false
      form.email = ''
    } else {
      message.error(result.error || 'Registration failed')
    }
  } finally {
    loading.value = false
  }
}

const handleLogin = async () => {
  loading.value = true
  try {
    const result = await authStore.login(form.username, form.password, form.mfaCode)
    
    if (result.success) {
      message.success('Login successful')
      router.push('/')
    } else if (result.requiresMFA) {
      showMFA.value = true
      message.info('Please enter your MFA code')
    } else {
      message.error(result.error || 'Login failed')
    }
  } finally {
    loading.value = false
  }
}

const handleSSOLogin = () => {
  window.location.href = 'https://gproc-backend-demo.fly.dev/api/v1/auth/sso/login'
}
</script>

<style scoped>
.login { min-height: 100vh; display: grid; grid-template-columns: 1.2fr 1fr; }
.panel { position: relative; }
.panel.left { padding: 56px; background: linear-gradient(180deg, #0b1220 0%, #0c1223 100%); color: #e5e7eb; overflow: hidden; }
.panel.right { display: flex; align-items: center; justify-content: center; background: linear-gradient(180deg, #0f172a 0%, #111827 100%); }

.brand { display: flex; align-items: center; gap: 10px; margin-bottom: 24px; }
.brand .logo { font-size: 28px; }
.brand .name { font-weight: 700; letter-spacing: .3px; }

.left h1 { font-size: 44px; margin: 0 0 10px; letter-spacing: -.02em; }
.left p { margin: 0 0 16px; color: #9fb2d2; max-width: 540px; }
.highlights { list-style: none; padding: 0; margin: 18px 0 0; display: grid; gap: 10px; color: #9fb2d2; }
.backdrop { position: absolute; inset: -20%; background: radial-gradient(closest-side, rgba(96,165,250,.2), transparent 60%), radial-gradient(closest-side, rgba(167,139,250,.18), transparent 60%); filter: blur(60px); opacity: .8; z-index: 0; }

.card { position: relative; background: linear-gradient(180deg, rgba(30,41,59,.8), rgba(17,24,39,.85)); border: 1px solid rgba(148,163,184,.25); border-radius: 16px; padding: 28px; width: 100%; max-width: 460px; box-shadow: 0 30px 70px rgba(0,0,0,.35); backdrop-filter: blur(6px); }
.card-header { text-align: left; margin-bottom: 16px; }
.card-header .badge { display: inline-block; background: rgba(99,102,241,.15); color: #c7d2fe; border: 1px solid rgba(99,102,241,.35); padding: 4px 8px; border-radius: 999px; font-size: 12px; margin-bottom: 6px; }
.card-header h2 { margin: 0; color: var(--n-text-color); }
.card-header p { margin: 6px 0 0; color: var(--n-text-color-3); }

.row { display: flex; align-items: center; justify-content: space-between; margin: 8px 0 16px; }
.link { color: #8fb4ff; text-decoration: none; font-size: 13px; }
.link:hover { text-decoration: underline; }

/* stronger input contrast in dark 
   (Naive UI variables differ by theme, so use deep selector) */
:deep(.n-input) {
  background-color: rgba(15,23,42,.6) !important;
}
:deep(.n-input .n-input__border) {
  border-color: rgba(148,163,184,.35) !important;
}
:deep(.n-input--focus .n-input__state-border) {
  box-shadow: 0 0 0 2px rgba(99,102,241,.35) inset;
}

.sso { margin-top: 16px; }
.toggle { margin-top: 16px; text-align: center; color: var(--n-text-color-3); font-size: 14px; }
.toggle .link { margin-left: 4px; cursor: pointer; }

@media (max-width: 980px) {
  .login { grid-template-columns: 1fr; }
  .panel.left { display: none; }
}
</style>