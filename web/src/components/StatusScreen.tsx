import { useEffect, useState } from 'react'
import { fetchStatus } from '../services/api'

interface Status {
  state: string
  uptime: string
  config: string
}

export default function StatusScreen() {
  const [status, setStatus] = useState<Status | null>(null)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    fetchStatus()
      .then(setStatus)
      .catch(err => setError(err.message))
  }, [])

  if (error) {
    return <div><h2>Status</h2><p>Erro: {error}</p></div>
  }

  if (!status) {
    return <div><h2>Status</h2><p>Carregando...</p></div>
  }

  return (
    <div>
      <h2>Status</h2>
      <p>Estado: {status.state}</p>
      <p>Uptime: {status.uptime}</p>
      <p>Config: {status.config}</p>
    </div>
  )
}
