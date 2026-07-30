import './global.css'
import { useState } from 'react'
import StatusScreen from './components/StatusScreen'
import MonitoringScreen from './components/MonitoringScreen'
import ConfigScreen from './components/ConfigScreen'
import LogsScreen from './components/LogsScreen'
import RegisterSW from './components/RegisterSW'

type Tab = 'status' | 'monitoring' | 'config' | 'logs'

const NAV_ITEMS: { id: Tab; label: string }[] = [
  { id: 'status', label: 'Status' },
  { id: 'monitoring', label: 'Monitoramento' },
  { id: 'config', label: 'Configurações' },
  { id: 'logs', label: 'Logs' },
]

export default function App() {
  const [tab, setTab] = useState<Tab>('status')

  return (
    <div className="app">
      <header className="app-header">
        <h1>GhostStack</h1>
        <span className="app-subtitle">Privacy Orchestration</span>
      </header>

      <nav className="app-nav">
        {NAV_ITEMS.map(item => (
          <button
            key={item.id}
            className={tab === item.id ? 'active' : ''}
            onClick={() => setTab(item.id)}
          >
            {item.label}
          </button>
        ))}
      </nav>

      <main className="app-main">
        {tab === 'status' && <StatusScreen />}
        {tab === 'monitoring' && <MonitoringScreen />}
        {tab === 'config' && <ConfigScreen />}
        {tab === 'logs' && <LogsScreen />}
      </main>

      <RegisterSW />
    </div>
  )
}
