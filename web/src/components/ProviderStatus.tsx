import { useWebSocketStatus, useWebSocketMetrics } from '../hooks/useWebSocket'

export default function ProviderStatus() {
  const status = useWebSocketStatus()
  const metrics = useWebSocketMetrics()

  return (
    <div className="card">
      <h3>Provider</h3>
      {status ? (
        <div className="status-grid">
          <div className="stat">
            <span className="stat-label">Estado</span>
            <span className={`stat-value ${status.state === 'running' ? 'green' : 'yellow'}`}>
              {status.state}
            </span>
          </div>
          <div className="stat">
            <span className="stat-label">Versão</span>
            <span className="stat-value">{status.version}</span>
          </div>
          <div className="stat">
            <span className="stat-label">Uptime</span>
            <span className="stat-value">{status.uptime}</span>
          </div>
        </div>
      ) : (
        <p className="dim">Conectando...</p>
      )}

      {metrics && (
        <div className="metrics-row">
          <div className="metric">
            <span className="num">{metrics.cpu.toFixed(1)}%</span>
            <span className="lbl">CPU</span>
          </div>
          <div className="metric">
            <span className="num">{fmtBytes(metrics.memory)}</span>
            <span className="lbl">RAM</span>
          </div>
          <div className="metric">
            <span className="num">{fmtBytes(metrics.rx_bytes)}/s</span>
            <span className="lbl">RX</span>
          </div>
          <div className="metric">
            <span className="num">{fmtBytes(metrics.tx_bytes)}/s</span>
            <span className="lbl">TX</span>
          </div>
        </div>
      )}
    </div>
  )
}

function fmtBytes(b: number): string {
  if (b >= 1 << 30) return (b / (1 << 30)).toFixed(1) + ' GB'
  if (b >= 1 << 20) return (b / (1 << 20)).toFixed(1) + ' MB'
  if (b >= 1 << 10) return (b / (1 << 10)).toFixed(0) + ' KB'
  return b + ' B'
}
