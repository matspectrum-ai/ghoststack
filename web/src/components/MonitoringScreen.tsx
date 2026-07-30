import { useEffect, useState } from 'react'
import { fetchMonitoring } from '../services/api'
import { useWebSocketMonitoring } from '../hooks/useWebSocket'

interface Monitoring {
  cpu: number
  memory: number
  network: { in: number; out: number }
}

export default function MonitoringScreen() {
  const [monitoring, setMonitoring] = useState<Monitoring | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)
  const wsMonitoring = useWebSocketMonitoring()

  useEffect(() => {
    fetchMonitoring()
      .then(setMonitoring)
      .catch(err => setError(err.message))
      .finally(() => setLoading(false))
  }, [])

  useEffect(() => {
    if (wsMonitoring) {
      setMonitoring(wsMonitoring)
      setError(null)
    }
  }, [wsMonitoring])

  if (error && !monitoring) {
    return (
      <div className="screen">
        <h2>Monitoramento</h2>
        <div className="error">Erro: {error}</div>
      </div>
    )
  }

  if (loading && !monitoring) {
    return (
      <div className="screen">
        <h2>Monitoramento</h2>
        <div className="loading">Carregando...</div>
      </div>
    )
  }

  if (!monitoring) {
    return (
      <div className="screen">
        <h2>Monitoramento</h2>
        <div className="loading">Carregando...</div>
      </div>
    )
  }

  return (
    <div className="screen">
      <h2>Monitoramento</h2>
      <div className="monitoring-grid">
        <div className="monitoring-card">
          <span className="monitoring-label">CPU</span>
          <span className="monitoring-value">{monitoring.cpu.toFixed(1)}%</span>
        </div>
        <div className="monitoring-card">
          <span className="monitoring-label">Memória</span>
          <span className="monitoring-value">{monitoring.memory.toFixed(1)}%</span>
        </div>
        <div className="monitoring-card">
          <span className="monitoring-label">Rede</span>
          <span className="monitoring-value">{monitoring.network.in.toFixed(1)} / {monitoring.network.out.toFixed(1)} MB/s</span>
        </div>
      </div>
    </div>
  )
}
