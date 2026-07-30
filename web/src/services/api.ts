export interface LogEntry {
  id: string
  timestamp: string
  level: string
  message: string
}

export interface MonitoringResponse {
  cpu: number
  memory: number
  network: {
    in: number
    out: number
  }
}

export interface StatusResponse {
  state: string
  uptime: string
  config: string
}

const API_BASE = '/api'

export async function fetchStatus(): Promise<StatusResponse> {
  const res = await fetch(`${API_BASE}/status`)
  if (!res.ok) {
    throw new Error(`Failed to fetch status: ${res.statusText}`)
  }
  return res.json()
}

export async function fetchMonitoring(): Promise<MonitoringResponse> {
  const res = await fetch(`${API_BASE}/monitoring`)
  if (!res.ok) {
    throw new Error(`Failed to fetch monitoring: ${res.statusText}`)
  }
  return res.json()
}

export async function fetchLogs(limit = 50): Promise<LogEntry[]> {
  const res = await fetch(`${API_BASE}/logs?limit=${limit}`)
  if (!res.ok) {
    throw new Error(`Failed to fetch logs: ${res.statusText}`)
  }
  return res.json()
}
