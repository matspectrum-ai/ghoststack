import { useState, useEffect } from 'react'
import { wsService, WSMessage } from '../services/websocket'

interface Point {
  t: number
  rx: number
  tx: number
}

export default function TrafficChart() {
  const [points, setPoints] = useState<Point[]>([])

  useEffect(() => {
    const wsUrl = `${window.location.protocol === 'https:' ? 'wss:' : 'ws:'}//${window.location.host}/api/ws`
    wsService.connect(wsUrl)

    const unsub = wsService.on('metrics', (msg: WSMessage) => {
      const data = msg.data as { rx_bytes: number; tx_bytes: number }
      setPoints(prev => {
        const next = [...prev, { t: Date.now(), rx: data.rx_bytes, tx: data.tx_bytes }]
        return next.slice(-30)
      })
    })

    return () => {
      unsub()
    }
  }, [])

  if (points.length < 2) {
    return (
      <div className="card">
        <h3>Tráfego</h3>
        <p className="dim">Aguardando dados...</p>
      </div>
    )
  }

  const maxRx = Math.max(...points.map(p => p.rx), 1)
  const maxTx = Math.max(...points.map(p => p.tx), 1)
  const max = Math.max(maxRx, maxTx)

  const w = 280
  const h = 80

  const toPath = (key: 'rx' | 'tx') => {
    const pts = points.map((p, i) => {
      const x = (i / (points.length - 1)) * w
      const y = h - (p[key] / max) * h
      return `${i === 0 ? 'M' : 'L'}${x.toFixed(0)},${y.toFixed(0)}`
    })
    return pts.join(' ')
  }

  return (
    <div className="card">
      <h3>Tráfego (tempo real)</h3>
      <svg viewBox={`0 0 ${w} ${h}`} className="chart">
        <path d={toPath('rx')} fill="none" stroke="#22c55e" strokeWidth="2" />
        <path d={toPath('tx')} fill="none" stroke="#3b82f6" strokeWidth="2" />
      </svg>
      <div className="chart-legend">
        <span><span className="dot green" /> RX</span>
        <span><span className="dot blue" /> TX</span>
      </div>
    </div>
  )
}
