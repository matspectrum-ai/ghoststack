import { useEffect, useState } from 'react'
import { fetchLogs, LogEntry } from '../services/api'

export default function LogsScreen() {
  const [logs, setLogs] = useState<LogEntry[]>([])
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    fetchLogs()
      .then(setLogs)
      .catch(err => setError(err.message))
  }, [])

  if (error) {
    return <div><h2>Logs</h2><p>Erro: {error}</p></div>
  }

  return (
    <div>
      <h2>Logs</h2>
      {logs.length === 0 && <p>Nenhum log disponível.</p>}
      <ul>
        {logs.map(log => (
          <li key={log.id}>
            <strong>{log.timestamp}</strong> [{log.level}] {log.message}
          </li>
        ))}
      </ul>
    </div>
  )
}
