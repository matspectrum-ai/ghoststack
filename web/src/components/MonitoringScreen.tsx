import { useWebSocketMetrics } from '../hooks/useWebSocket'
import TrafficChart from './TrafficChart'
import PeerList from './PeerList'
import AlertPanel from './AlertPanel'

export default function MonitoringScreen() {
  const metrics = useWebSocketMetrics()

  return (
    <div className="screen">
      <h2>Monitoramento</h2>

      <div className="monitoring-grid">
        {metrics && (
          <div className="card">
            <h3>Recursos</h3>
            <div className="metrics-row">
              <div className="metric big">
                <span className="num">{metrics.cpu.toFixed(1)}%</span>
                <span className="lbl">CPU</span>
              </div>
              <div className="metric big">
                <span className="num">{fmtBytes(metrics.memory)}</span>
                <span className="lbl">RAM</span>
              </div>
              <div className="metric big">
                <span className="num">{fmtBytes(metrics.rx_bytes)}/s</span>
                <span className="lbl">Download</span>
              </div>
              <div className="metric big">
                <span className="num">{fmtBytes(metrics.tx_bytes)}/s</span>
                <span className="lbl">Upload</span>
              </div>
            </div>
          </div>
        )}

        <TrafficChart />
      </div>

      <PeerList />
      <AlertPanel />
    </div>
  )
}

function fmtBytes(b: number): string {
  if (b >= 1 << 30) return (b / (1 << 30)).toFixed(1) + ' GB'
  if (b >= 1 << 20) return (b / (1 << 20)).toFixed(1) + ' MB'
  if (b >= 1 << 10) return (b / (1 << 10)).toFixed(0) + ' KB'
  return b + ' B'
}
