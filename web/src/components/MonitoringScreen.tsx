import { useEffect, useState } from 'react'
import { fetchMonitoring } from '../services/api'

interface Monitoring {
  cpu: number
  memory: number
  network: { in: number; out: number }
}

export default function MonitoringScreen() {
  const [monitoring, setMonitoring] = useState<Monitoring | null>(null)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    fetchMonitoring()
      .then(setMonitoring)
      .catch(err => setError(err.message))
  }, [])

  if (error) {
    return <div><h2>Monitoramento</h2><p>Erro: {error}</p></div>
  }

  if (!monitoring) {
    return <div><h2>Monitoramento</h2><p>Carregando...</p></div>
  }

  return (
    <div>
      <h2>Monitoramento</h2>
      <p>CPU: {monitoring.cpu}%</p>
      <p>Memória: {monitoring.memory}%</p>
      <p>Rede: {monitoring.network.in} in / {monitoring.network.out} out</p>
    </div>
  )
}
