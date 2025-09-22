import { defineStore } from 'pinia'
import { ref } from 'vue'

type WS = WebSocket | null

export const useWebSocketStore = defineStore('websocket', () => {
  const socket = ref<WS>(null)
  const connected = ref(false)
  const realTimeData = ref<any>({})

  const connect = (_token?: string) => {
    // Backend WS is Gorilla WebSocket at /api/v1/ws
    const url = `ws://${location.hostname}:8080/api/v1/ws`

    if (socket.value) {
      socket.value.close()
    }

    const ws = new WebSocket(url)
    socket.value = ws

    ws.onopen = () => {
      connected.value = true
      console.log('WebSocket connected')
    }

    ws.onclose = () => {
      connected.value = false
      console.log('WebSocket disconnected')
    }

    ws.onerror = (e) => {
      console.warn('WebSocket error', e)
    }

    ws.onmessage = (evt) => {
      try {
        const msg = JSON.parse(evt.data)
        const type = msg.type || 'unknown'
        const data = msg.data || msg

        switch (type) {
          case 'process_update':
            realTimeData.value.processes = data
            break
          case 'metrics_update':
            realTimeData.value.metrics = data
            break
          case 'alert':
            realTimeData.value.alerts = realTimeData.value.alerts || []
            realTimeData.value.alerts.unshift(data)
            break
          case 'cluster_update':
            realTimeData.value.cluster = data
            break
          case 'deployment_update':
            realTimeData.value.deployments = data
            break
          case 'log_stream':
            realTimeData.value.logs = realTimeData.value.logs || {}
            if (data?.processId && data?.line) {
              const pid = data.processId
              realTimeData.value.logs[pid] = realTimeData.value.logs[pid] || []
              realTimeData.value.logs[pid].push(data.line)
              if (realTimeData.value.logs[pid].length > 1000) {
                realTimeData.value.logs[pid] = realTimeData.value.logs[pid].slice(-1000)
              }
            }
            break
          default:
            // Unknown type; ignore
            break
        }
      } catch (_) {
        // Ignore malformed messages
      }
    }
  }

  const disconnect = () => {
    if (socket.value) {
      try { socket.value.close() } catch {}
      socket.value = null
      connected.value = false
    }
  }

  // No-ops kept for API compatibility
  const subscribeToLogs = (_processId: string) => {}
  const unsubscribeFromLogs = (_processId: string) => {}

  return {
    socket,
    connected,
    realTimeData,
    connect,
    disconnect,
    subscribeToLogs,
    unsubscribeFromLogs
  }
})