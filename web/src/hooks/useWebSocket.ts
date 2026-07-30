import { useEffect, useState } from 'react'
import { wsService, WebSocketMessage } from '../services/websocket'

export function useWebSocketStatus() {
  const [status, setStatus] = useState<{ state: string; uptime: string; config: string } | null>(null)

  useEffect(() => {
    const wsUrl = `${window.location.protocol === 'https:' ? 'wss:' : 'ws:'}//${window.location.host}/ws`
    wsService.connect(wsUrl)

    const unsubscribe = wsService.on('status', (message: WebSocketMessage) => {
      setStatus(message.payload as { state: string; uptime: string; config: string })
    })

    return () => {
      unsubscribe()
      wsService.disconnect()
    }
  }, [])

  return status
}

export function useWebSocketMonitoring() {
  const [monitoring, setMonitoring] = useState<{ cpu: number; memory: number; network: { in: number; out: number } } | null>(null)

  useEffect(() => {
    const unsubscribe = wsService.on('monitoring', (message: WebSocketMessage) => {
      setMonitoring(message.payload as { cpu: number; memory: number; network: { in: number; out: number } })
    })

    return () => {
      unsubscribe()
    }
  }, [])

  return monitoring
}

export function useWebSocketLogs() {
  const [logs, setLogs] = useState<{ id: string; timestamp: string; level: string; message: string }[]>([])

  useEffect(() => {
    const unsubscribe = wsService.on('log', (message: WebSocketMessage) => {
      setLogs(prev => [message.payload as { id: string; timestamp: string; level: string; message: string }, ...prev].slice(0, 100))
    })

    return () => {
      unsubscribe()
    }
  }, [])

  return logs
}
