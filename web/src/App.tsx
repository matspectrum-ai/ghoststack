import { useState } from 'react'
import StatusScreen from './components/StatusScreen'
import MonitoringScreen from './components/MonitoringScreen'
import ConfigScreen from './components/ConfigScreen'
import LogsScreen from './components/LogsScreen'

type Tab = 'status' | 'monitoring' | 'config' | 'logs'

export default function App() {
  const [tab, setTab] = useState<Tab>('status')

  return (
    <div style={{ padding: 24 }}>
      <h1>GhostStack Dashboard</h1>
      <nav style={{ marginBottom: 16 }}>
        <button onClick={() => setTab('status')}>Status</button>
        <button onClick={() => setTab('monitoring')}>Monitoramento</button>
        <button onClick={() => setTab('config')}>Configurações</button>
        <button onClick={() => setTab('logs')}>Logs</button>
      </nav>
      {tab === 'status' && <StatusScreen />}
      {tab === 'monitoring' && <MonitoringScreen />}
      {tab === 'config' && <ConfigScreen />}
      {tab === 'logs' && <LogsScreen />}
    </div>
  )
}
