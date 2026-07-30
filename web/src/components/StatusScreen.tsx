import { useEffect, useState } from 'react'
import { fetchStatus } from '../services/api'
import { useWebSocketStatus } from '../hooks/useWebSocket'

interface Status {
  state: string
  uptime: string
  config: string
}

export default function StatusScreen() {
  const [status, setStatus] = useState<Status | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)
  const wsStatus = useWebSocketStatus()

  useEffect(() => {
    fetchStatus()
      .then(setStatus)
      .catch(err => setError(err.message))
      .finally(() => setLoading(false))
  }, [])

  useEffect(() => {
    if (wsStatus) {
      setStatus(wsStatus)
      setError(null)
    }
  }, [wsStatus])

  if (error && !status) {
    return (
      <div className="screen">
        <h2>Status</h2>
        <div className="error">Erro: {error}</div>
      </div>
    )
  }

  if (loading && !status) {
    return (
      <div className="screen">
        <h2>Status</h2>
        <div className="loading">Carregando...</div>
      </div>
    )
  }

  if (!status) {
    return (
      <div className="screen">
        <h2>Status</h2>
        <div className="loading">Carregando...</div>
      </div>
    )
  }

  return (
    <div className="screen">
      <h2>Status</h2>
      <div className="status-grid">
        <div className="status-card">
          <span className="status-label">Estado</span>
          <span className={`status-value ${status.state === 'running' ? 'success' : 'warning'}`}>
            {status.state}
          </span>
        </div>
        <div className="status-card">
          <span className="status-label">Uptime</span>
          <span className="status-value">{status.uptime}</span>
        </div>
        <div className="status-card">
          <span className="status-label">Config</span>
          <span className="status-value">{status.config}</span>
        </div>
      </div>
    </div>
  )
}
