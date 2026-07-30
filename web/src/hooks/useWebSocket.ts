import { useEffect, useState, useCallback } from 'react'
import { wsService, WSMessage } from '../services/websocket'

interface ProviderStatus {
  state: string
  uptime: string
  version: string
}

interface MetricsData {
  cpu: number
  memory: number
  rx_bytes: number
  tx_bytes: number
  uptime: number
}

export function useWebSocketStatus() {
  const [status, setStatus] = useState<ProviderStatus | null>(null)

  useEffect(() => {
    const wsUrl = `${window.location.protocol === 'https:' ? 'wss:' : 'ws:'}//${window.location.host}/api/ws`
    wsService.connect(wsUrl)

    const unsub = wsService.on('provider_status', (msg: WSMessage) => {
      setStatus(msg.data as ProviderStatus)
    })

    return () => {
      unsub()
      wsService.disconnect()
    }
  }, [])

  return status
}

export function useWebSocketMetrics() {
  const [metrics, setMetrics] = useState<MetricsData | null>(null)

  useEffect(() => {
    const wsUrl = `${window.location.protocol === 'https:' ? 'wss:' : 'ws:'}//${window.location.host}/api/ws`
    wsService.connect(wsUrl)

    const unsub = wsService.on('metrics', (msg: WSMessage) => {
      setMetrics(msg.data as MetricsData)
    })

    return () => {
      unsub()
    }
  }, [])

  return metrics
}

export function useWebSocketLogs() {
  const [logs, setLogs] = useState<{ id: string; timestamp: string; level: string; message: string }[]>([])

  useEffect(() => {
    const wsUrl = `${window.location.protocol === 'https:' ? 'wss:' : 'ws:'}//${window.location.host}/api/ws`
    wsService.connect(wsUrl)

    const unsub = wsService.on('log', (msg: WSMessage) => {
      const entry = msg.data as { id: string; timestamp: string; level: string; message: string }
      setLogs(prev => [entry, ...prev].slice(0, 100))
    })

    return () => {
      unsub()
    }
  }, [])

  return logs
}

export function useWebSocketAll() {
  const [lastMessage, setLastMessage] = useState<WSMessage | null>(null)

  useEffect(() => {
    const wsUrl = `${window.location.protocol === 'https:' ? 'wss:' : 'ws:'}//${window.location.host}/api/ws`
    wsService.connect(wsUrl)

    const unsub = wsService.on('*', (msg: WSMessage) => {
      setLastMessage(msg)
    })

    return () => {
      unsub()
    }
  }, [])

  return lastMessage
}

export function usePolling<T>(fetcher: () => Promise<T>, interval = 3000) {
  const [data, setData] = useState<T | null>(null)
  const [error, setError] = useState<string | null>(null)

  const fetchData = useCallback(() => {
    fetcher()
      .then(setData)
      .catch(e => setError(e.message))
  }, [fetcher])

  useEffect(() => {
    fetchData()
    const timer = setInterval(fetchData, interval)
    return () => clearInterval(timer)
  }, [fetchData, interval])

  return { data, error }
}
