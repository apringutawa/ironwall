import axios from 'axios'

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8001'

export const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

export const getStatus = async () => {
  const response = await api.get('/api/v1/status')
  return response.data
}

export const getModules = async () => {
  const response = await api.get('/api/v1/modules')
  return response.data
}

export const getEvents = async (params?: any) => {
  const response = await api.get('/api/v1/events', { params })
  return response.data
}

export const getFirewallStatus = async () => {
  const response = await api.get('/api/v1/firewall/status')
  return response.data
}

export const getBlockedIPs = async () => {
  const response = await api.get('/api/v1/firewall/blocked')
  return response.data
}

export const blockIP = async (ip: string, reason?: string) => {
  const response = await api.post('/api/v1/firewall/block', { ip, reason })
  return response.data
}

export const unblockIP = async (ip: string) => {
  const response = await api.post('/api/v1/firewall/unblock', null, { params: { ip } })
  return response.data
}

export const enableModule = async (moduleName: string) => {
  const response = await api.post(`/api/v1/modules/${moduleName}/enable`)
  return response.data
}

export const disableModule = async (moduleName: string) => {
  const response = await api.post(`/api/v1/modules/${moduleName}/disable`)
  return response.data
}

export const runScan = async (options?: any) => {
  const response = await api.post('/api/v1/scan', options)
  return response.data
}

export const protectSystem = async () => {
  const response = await api.post('/api/v1/protect')
  return response.data
}

export const unprotectSystem = async () => {
  const response = await api.post('/api/v1/unprotect')
  return response.data
}

export const rollbackSystem = async () => {
  const response = await api.post('/api/v1/rollback')
  return response.data
}

export const configureAlerts = async (config: any) => {
  const response = await api.post('/api/v1/alerts/configure', config)
  return response.data
}

export const getHardeningStatus = async () => {
  const response = await api.get('/api/v1/hardening/status')
  return response.data
}

export const applyKernelHardening = async () => {
  const response = await api.post('/api/v1/hardening/kernel/apply')
  return response.data
}

export const applyAccountHardening = async () => {
  const response = await api.post('/api/v1/hardening/accounts/apply')
  return response.data
}

export const initAIDE = async () => {
  const response = await api.post('/api/v1/hardening/aide/init')
  return response.data
}

export const checkAIDE = async () => {
  const response = await api.post('/api/v1/hardening/aide/check')
  return response.data
}

export const minimizeServices = async () => {
  const response = await api.post('/api/v1/hardening/services/minimize')
  return response.data
}

export const listServices = async () => {
  const response = await api.get('/api/v1/hardening/services/list')
  return response.data
}

export const runAudit = async () => {
  const response = await api.post('/api/v1/hardening/audit/run')
  return response.data
}

export const installAudit = async () => {
  const response = await api.post('/api/v1/hardening/audit/install')
  return response.data
}

export const getAuditStatus = async () => {
  const response = await api.get('/api/v1/hardening/audit/status')
  return response.data
}
