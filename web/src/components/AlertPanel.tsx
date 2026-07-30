import { useState } from 'react'

interface Alert {
  id: string
  severity: 'info' | 'warning' | 'critical'
  title: string
  message: string
  timestamp: string
  acknowledged: boolean
}

const MOCK_ALERTS: Alert[] = [
  { id: 'a1', severity: 'info', title: 'Provider started', message: 'WireGuard provider initialized', timestamp: new Date().toISOString(), acknowledged: false },
  { id: 'a2', severity: 'warning', title: 'High latency', message: 'Ping to peer exceeded 200ms', timestamp: new Date().toISOString(), acknowledged: false },
]

export default function AlertPanel() {
  const [alerts, setAlerts] = useState(MOCK_ALERTS)

  const ack = (id: string) => {
    setAlerts(prev => prev.map(a => a.id === id ? { ...a, acknowledged: true } : a))
  }

  const active = alerts.filter(a => !a.acknowledged)

  return (
    <div className="card">
      <h3>Alertas {active.length > 0 && <span className="badge">{active.length}</span>}</h3>
      {alerts.length === 0 ? (
        <p className="dim">Nenhum alerta</p>
      ) : (
        <div className="alert-list">
          {alerts.slice(0, 10).map(a => (
            <div key={a.id} className={`alert-row ${a.severity} ${a.acknowledged ? 'dim' : ''}`}>
              <span className="alert-icon">{a.severity === 'critical' ? '🔴' : a.severity === 'warning' ? '🟡' : '🔵'}</span>
              <div className="alert-body">
                <strong>{a.title}</strong>
                <span className="alert-msg">{a.message}</span>
              </div>
              {!a.acknowledged && (
                <button className="ack-btn" onClick={() => ack(a.id)}>OK</button>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
