import { defineStore } from 'pinia'
import { ref } from 'vue'
import axios from 'axios'

export interface ClusterNode {
  id: string
  address: string
  role: 'leader' | 'follower' | 'candidate'
  status: 'active' | 'inactive' | 'failed'
  lastSeen: string
  version: string
  processes: number
  cpu: number
  memory: number
  uptime: string
}

export interface ClusterStatus {
  leader: string
  nodes: number
  healthy: boolean
  lastElection: string
  raftTerm: number
  raftIndex: number
}

export const useClusterStore = defineStore('cluster', () => {
  const nodes = ref<ClusterNode[]>([])
  const status = ref<ClusterStatus | null>(null)
  const loading = ref(false)

  const fetchNodes = async () => {
    loading.value = true
    try {
      const response = await axios.get('/api/v1/cluster/nodes')
      nodes.value = response.data
    } catch (error) {
      console.error('Failed to fetch cluster nodes:', error)
    } finally {
      loading.value = false
    }
  }

  const fetchStatus = async () => {
    try {
      const response = await axios.get('/api/v1/cluster/status')
      status.value = response.data
    } catch (error) {
      console.error('Failed to fetch cluster status:', error)
    }
  }

  const promoteNode = async (nodeId: string) => {
    try {
      await axios.post(`/api/v1/cluster/nodes/${nodeId}/promote`)
      await fetchNodes()
      await fetchStatus()
      return { success: true }
    } catch (error: any) {
      return { success: false, error: error.response?.data?.message }
    }
  }

  const removeNode = async (nodeId: string) => {
    try {
      await axios.delete(`/api/v1/cluster/nodes/${nodeId}`)
      await fetchNodes()
      await fetchStatus()
      return { success: true }
    } catch (error: any) {
      return { success: false, error: error.response?.data?.message }
    }
  }

  return {
    nodes,
    status,
    loading,
    fetchNodes,
    fetchStatus,
    promoteNode,
    removeNode
  }
})