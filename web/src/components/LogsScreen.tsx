import { useEffect, useState } from 'react'
import { fetchLogs } from '../services/api'
import { useWebSocketLogs } from '../hooks/useWebSocket'

export default function LogsScreen() {
  const [logs, setLogs] = useState<{ id: string; timestamp: string; level: string; message: string }[]>([])
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)
  const wsLogs = useWebSocketLogs()

  useEffect(() => {
    fetchLogs()
      .then(setLogs)
      .catch(err => setError(err.message))
      .finally(() => setLoading(false))
  }, [])

  useEffect(() => {
    if (wsLogs.length > 0) {
      setLogs(prev => [...wsLogs, ...prev].slice(0, 100))
      setError(null)
    }
  }, [wsLogs])

  if (error && logs.length === 0) {
    return (
      <div className="screen">
        <h2>Logs</h2>
        <div className="error">Erro: {error}</div>
      </div>
    )
  }

  if (loading && logs.length === 0) {
    return (
      <div className="screen">
        <h2>Logs</h2>
        <div className="loading">Carregando...</div>
      </div>
    )
  }

  return (
    <div className="screen">
      <h2>Logs</h2>
      {logs.length === 0 && <p className="empty">Nenhum log disponível.</p>}
      <div className="logs-list">
        {logs.map(log => (
          <div key={log.id} className={`log-entry log-${log.level.toLowerCase()}`}>
            <span className="log-time">{log.timestamp}</span>
            <span className="log-level">{log.level}</span>
            <span className="log-message">{log.message}</span>
          </div>
        ))}
      </div>
    </div>
  )
}
