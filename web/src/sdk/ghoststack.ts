export interface GhostStackClientConfig {
  baseUrl: string
  token?: string
}

export interface StatusResponse {
  state: string
  uptime: string
  config: string
}

export interface MonitoringResponse {
  cpu: number
  memory: number
  network: {
    in: number
    out: number
  }
}

export interface LogEntry {
  id: string
  timestamp: string
  level: string
  message: string
}

export class GhostStackClient {
  private baseUrl: string
  private token?: string

  constructor(config: GhostStackClientConfig) {
    this.baseUrl = config.baseUrl.replace(/\/$/, '')
    this.token = config.token
  }

  async status(): Promise<StatusResponse> {
    const res = await fetch(`${this.baseUrl}/api/v1/status`, {
      headers: this.authHeaders(),
    })
    if (!res.ok) {
      throw new Error(`Failed to fetch status: ${res.statusText}`)
    }
    return res.json()
  }

  async monitoring(): Promise<MonitoringResponse> {
    const res = await fetch(`${this.baseUrl}/api/v1/monitoring`, {
      headers: this.authHeaders(),
    })
    if (!res.ok) {
      throw new Error(`Failed to fetch monitoring: ${res.statusText}`)
    }
    return res.json()
  }

  async logs(limit = 50): Promise<LogEntry[]> {
    const res = await fetch(`${this.baseUrl}/api/v1/logs?limit=${limit}`, {
      headers: this.authHeaders(),
    })
    if (!res.ok) {
      throw new Error(`Failed to fetch logs: ${res.statusText}`)
    }
    return res.json()
  }

  private authHeaders(): Record<string, string> {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
    }
    if (this.token) {
      headers['Authorization'] = `Bearer ${this.token}`
    }
    return headers
  }
}
