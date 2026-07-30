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
const MAX_RETRIES = 2

async function fetchWithRetry(input: RequestInfo, init?: RequestInit): Promise<Response> {
  let lastError: Error | null = null
  for (let attempt = 0; attempt <= MAX_RETRIES; attempt++) {
    try {
      const response = await fetch(input, init)
      if (!response.ok) {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`)
      }
      return response
    } catch (error) {
      lastError = error as Error
      if (attempt < MAX_RETRIES) {
        await new Promise(resolve => setTimeout(resolve, 500 * (attempt + 1)))
      }
    }
  }
  throw lastError
}

export async function fetchStatus(): Promise<StatusResponse> {
  const res = await fetchWithRetry(`${API_BASE}/status`)
  return res.json()
}

export async function fetchMonitoring(): Promise<MonitoringResponse> {
  const res = await fetchWithRetry(`${API_BASE}/monitoring`)
  return res.json()
}

export async function fetchLogs(limit = 50): Promise<LogEntry[]> {
  const res = await fetchWithRetry(`${API_BASE}/logs?limit=${limit}`)
  return res.json()
}
