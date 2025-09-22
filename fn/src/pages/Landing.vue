<template>
  <div class="landing">
    <header class="landing-header">
      <div class="brand">
        <span class="logo">🚀</span>
        <span class="name">GProc</span>
      </div>
      <nav class="nav">
        <a href="#features">Features</a>
        <a href="#why">Why GProc</a>
        <a href="#security">Security</a>
        <n-button secondary size="large" @click="toggleTheme">{{ isDark ? 'Light' : 'Dark' }}</n-button>
        <n-button type="primary" strong size="large" @click="goLogin">Launch Console</n-button>
      </nav>
    </header>

    <section class="hero">
      <div class="anim" aria-hidden="true">
        <div class="blob b1"></div>
        <div class="blob b2"></div>
        <div class="blob b3"></div>
      </div>
      <div class="copy">
        <div class="badge"><span>New</span> Enterprise dashboard</div>
        <h1>Scale processes with confidence</h1>
        <p>Ship faster with zero‑downtime deploys, observability, and RBAC — all in a beautiful PM2‑style console.</p>
        <div class="cta">
          <n-button type="primary" size="large" strong @click="goLogin">Open Dashboard</n-button>
          <n-button size="large" tertiary @click="scrollTo('#features')">Explore features</n-button>
        </div>

        <div class="trust">
          <span>Trusted by ops teams</span>
          <div class="dots">
            <i></i><i></i><i></i><i></i>
          </div>
        </div>
      </div>

      <div class="preview">
        <div class="preview-inner">
          <div class="preview-bar">
            <i></i><i></i><i></i>
            <span>gproc — processes</span>
          </div>
          <pre class="preview-code"><code>{{ demo }}</code></pre>
        </div>
      </div>
    </section>

    <section id="features" class="features">
      <div class="feature">
        <div class="icon">📊</div>
        <h3>PM2‑style Dashboard</h3>
        <p>Live KPIs, process list, restarts, uptime, CPU & memory in a clean dark UI.</p>
      </div>
      <div class="feature">
        <div class="icon">🚦</div>
        <h3>Zero‑downtime</h3>
        <p>Blue/green and graceful restarts keep your services available while you deploy.</p>
      </div>
      <div class="feature">
        <div class="icon">🔐</div>
        <h3>RBAC & Audit</h3>
        <p>Granular permissions with comprehensive audit logs for compliance.</p>
      </div>
      <div class="feature">
        <div class="icon">📈</div>
        <h3>Metrics & Health</h3>
        <p>Prometheus metrics, health checks, and alerts — insights at your fingertips.</p>
      </div>
    </section>

    <section id="why" class="why">
      <div class="why-card">
        <h4>Operate at scale</h4>
        <p>Unified process management, observability, and security — from one console.</p>
      </div>
      <div class="why-card">
        <h4>Minutes to value</h4>
        <p>Run the server, open the console, and start managing processes instantly.</p>
      </div>
      <div class="why-card">
        <h4>Built for teams</h4>
        <p>Roles, permissions, and audit trails designed for real‑world operations.</p>
      </div>
    </section>

    <footer class="landing-footer">
      <span>© {{ year }} GProc. All rights reserved.</span>
    </footer>
  </div>
  
</template>

<script setup lang="ts">
// import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useLocalStorage } from '@vueuse/core'
import { NButton } from 'naive-ui'

const router = useRouter()
const year = new Date().getFullYear()
const isDark = useLocalStorage('gproc-theme', true)

const demo = `App name      id  mode   pid    status   restart  uptime  memory  watching\napi-service   0   cluster 27387  online   0        2m      35.5MB  disabled\nweb-frontend  1   fork    27390  online   0        2m      26.1MB  enabled\nworker        2   fork    27292  online   0        2m      24.9MB  disabled`

const goLogin = () => router.push('/login')
const scrollTo = (sel: string) => {
  document.querySelector(sel)?.scrollIntoView({ behavior: 'smooth' })
}
const toggleTheme = () => {
  isDark.value = !isDark.value
  if (isDark.value) document.documentElement.setAttribute('data-theme', 'dark')
  else document.documentElement.removeAttribute('data-theme')
}
</script>

<style scoped>
.landing { min-height: 100vh; background: radial-gradient(1200px 600px at 80% -10%, rgba(99,102,241,.25), transparent), linear-gradient(180deg, #0b1220 0%, #0b1220 60%, #0f172a 100%); color: #e5e7eb; }
.landing-header { max-width: 1200px; margin: 0 auto; padding: 20px 24px; display: flex; align-items: center; justify-content: space-between; }
.brand { display: flex; align-items: center; gap: 10px; font-weight: 700; letter-spacing: .5px; }
.logo { font-size: 24px; }
.name { font-size: 18px; color: #fff; }
.nav { display: flex; align-items: center; gap: 18px; }
.nav a { color: #cbd5e1; text-decoration: none; font-size: 14px; }

.hero { max-width: 1200px; margin: 40px auto; padding: 24px; display: grid; grid-template-columns: 1.2fr .8fr; gap: 32px; align-items: center; }
.anim { position: absolute; inset: 0; pointer-events: none; overflow: hidden; }
.blob { position: absolute; width: 520px; height: 520px; filter: blur(60px); opacity: .25; border-radius: 50%; mix-blend-mode: screen; }
.b1 { background: #6366f1; top: -120px; left: -100px; animation: float 14s ease-in-out infinite alternate; }
.b2 { background: #22d3ee; right: -160px; top: -80px; animation: float 18s ease-in-out infinite alternate-reverse; }
.b3 { background: #a78bfa; bottom: -160px; left: 20%; animation: float 22s ease-in-out infinite alternate; }
.badge { display: inline-flex; align-items: center; gap: 8px; background: rgba(99,102,241,.15); color: #c7d2fe; border: 1px solid rgba(99,102,241,.3); padding: 6px 10px; border-radius: 999px; font-size: 12px; margin-bottom: 12px; }
.badge span { display: inline-block; background: #6366f1; color: #fff; padding: 2px 8px; border-radius: 999px; font-weight: 700; }
.copy h1 { font-size: 52px; line-height: 1.05; margin: 0 0 12px; color: #fff; letter-spacing: -.02em; }
.copy p { color: #a3b2cd; margin: 0 0 20px; font-size: 18px; }
.cta { display: flex; gap: 12px; }
.trust { display: flex; align-items: center; gap: 12px; margin-top: 18px; color: #91a4c9; font-size: 14px; }
.trust .dots { display: flex; gap: 6px; }
.trust .dots i { width: 6px; height: 6px; border-radius: 999px; background: rgba(148,163,184,.5); display: inline-block; }

.preview { position: relative; height: 380px; border-radius: 16px; background: linear-gradient(180deg, rgba(2,6,23,.75), rgba(2,6,23,.4)); border: 1px solid rgba(148,163,184,.15); overflow: hidden; box-shadow: 0 25px 60px rgba(0,0,0,.35); }
.preview-inner { position: absolute; inset: 0; display: flex; flex-direction: column; }
.preview-bar { height: 40px; display: flex; align-items: center; gap: 8px; padding: 0 12px; border-bottom: 1px solid rgba(148,163,184,.12); color: #96a6c7; font-size: 12px; }
.preview-bar i { width: 10px; height: 10px; border-radius: 999px; background: rgba(248,113,113,.9); box-shadow: 16px 0 0 rgba(250,204,21,.9), 32px 0 0 rgba(74,222,128,.9); }
.preview-code { flex: 1; margin: 0; padding: 16px; color: #cbd5e1; font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace; font-size: 13px; }

.features { max-width: 1200px; margin: 12px auto 60px; padding: 0 24px; display: grid; grid-template-columns: repeat(auto-fit, minmax(240px, 1fr)); gap: 16px; }
.feature { background: linear-gradient(180deg, rgba(30,41,59,.6), rgba(15,23,42,.5)); border: 1px solid rgba(148,163,184,.12); border-radius: 14px; padding: 18px; }
.feature .icon { font-size: 20px; margin-bottom: 8px; }
.feature h3 { margin: 4px 0 6px; color: #e2e8f0; font-weight: 600; letter-spacing: .2px; }
.feature p { margin: 0; color: #94a3b8; font-size: 14px; }

.why { max-width: 1200px; margin: 0 auto 64px; padding: 0 24px; display: grid; grid-template-columns: repeat(auto-fit, minmax(260px, 1fr)); gap: 16px; }
.why-card { background: linear-gradient(180deg, rgba(30,41,59,.55), rgba(15,23,42,.45)); border: 1px solid rgba(148,163,184,.12); border-radius: 14px; padding: 16px 18px; }
.why-card h4 { margin: 0 0 6px; color: #e2e8f0; }
.why-card p { margin: 0; color: #8aa0bf; font-size: 14px; }

.landing-footer { max-width: 1200px; margin: 0 auto; padding: 24px; color: #8aa0bf; font-size: 13px; text-align: center; border-top: 1px solid rgba(148,163,184,.12); }

@keyframes float { to { transform: translateY(20px) scale(1.03); } }
@keyframes wave { 0% { transform: translateY(0) } 100% { transform: translateY(10px) } }

@media (max-width: 900px) { .hero { grid-template-columns: 1fr; } }
</style>
