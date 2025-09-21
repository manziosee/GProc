<template>
  <div class="stats-card" :class="color">
    <div class="stats-content">
      <div class="stats-icon">
        <component :is="getIcon()" :size="24" />
      </div>
      
      <div class="stats-info">
        <div class="stats-title">{{ title }}</div>
        <div class="stats-value">{{ value }}</div>
        
        <div class="stats-change" :class="trend">
          <component :is="trend === 'up' ? TrendingUp : TrendingDown" :size="14" />
          <span>{{ Math.abs(change) }}% from last hour</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Activity, AlertTriangle, Cpu, HardDrive, TrendingUp, TrendingDown } from 'lucide-vue-next'

const props = defineProps<{
  title: string
  value: string | number
  change: number
  trend: 'up' | 'down'
  icon: string
  color: 'success' | 'error' | 'warning' | 'info'
}>()

const getIcon = () => {
  const icons = {
    Activity,
    AlertTriangle,
    Cpu,
    HardDrive
  }
  return icons[props.icon as keyof typeof icons] || Activity
}
</script>

<style scoped>
.stats-card {
  background: var(--n-color-embedded);
  border: 1px solid var(--n-border-color);
  border-radius: 12px;
  padding: 24px;
  transition: all 0.3s ease;
}

.stats-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
}

.stats-content {
  display: flex;
  align-items: flex-start;
  gap: 16px;
}

.stats-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.stats-card.success .stats-icon {
  background: rgba(22, 163, 74, 0.1);
  color: #16a34a;
}

.stats-card.error .stats-icon {
  background: rgba(220, 38, 38, 0.1);
  color: #dc2626;
}

.stats-card.warning .stats-icon {
  background: rgba(234, 88, 12, 0.1);
  color: #ea580c;
}

.stats-card.info .stats-icon {
  background: rgba(37, 99, 235, 0.1);
  color: #2563eb;
}

.stats-info {
  flex: 1;
}

.stats-title {
  color: var(--n-text-color-2);
  font-size: 14px;
  font-weight: 500;
  margin-bottom: 4px;
}

.stats-value {
  color: var(--n-text-color);
  font-size: 28px;
  font-weight: 700;
  margin-bottom: 8px;
}

.stats-change {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  font-weight: 500;
}

.stats-change.up {
  color: #16a34a;
}

.stats-change.down {
  color: #dc2626;
}
</style>