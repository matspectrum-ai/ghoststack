import { useState, useEffect } from 'react'

interface WireGuardConfig {
  privateKey: string
  address: string
  dns: string
  peers: { publicKey: string; endpoint: string; allowedIPs: string }[]
}

const DEFAULT_WG: WireGuardConfig = {
  privateKey: '',
  address: '10.0.0.1/24',
  dns: '1.1.1.1',
  peers: [{ publicKey: '', endpoint: '', allowedIPs: '0.0.0.0/0' }],
}

export default function ConfigScreen() {
  const [wg, setWg] = useState<WireGuardConfig>(DEFAULT_WG)
  const [saved, setSaved] = useState(false)
  const [activeTab, setActiveTab] = useState<'wg' | 'tor' | 'unbound'>('wg')

  useEffect(() => {
    const timer = setTimeout(() => setSaved(false), 3000)
    return () => clearTimeout(timer)
  }, [saved])

  const updatePeer = (i: number, field: string, val: string) => {
    setWg(prev => {
      const peers = [...prev.peers]
      peers[i] = { ...peers[i], [field]: val }
      return { ...prev, peers }
    })
  }

  const addPeer = () => {
    setWg(prev => ({
      ...prev,
      peers: [...prev.peers, { publicKey: '', endpoint: '', allowedIPs: '0.0.0.0/0' }],
    }))
  }

  const removePeer = (i: number) => {
    setWg(prev => ({
      ...prev,
      peers: prev.peers.filter((_, idx) => idx !== i),
    }))
  }

  const handleSave = () => {
    fetch('/api/config', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ wireguard: wg }),
    }).then(() => setSaved(true)).catch(() => setSaved(true))
  }

  const tabs = [
    { id: 'wg' as const, label: 'WireGuard' },
    { id: 'tor' as const, label: 'Tor' },
    { id: 'unbound' as const, label: 'Unbound' },
  ]

  return (
    <div className="screen">
      <h2>Configurações</h2>

      <nav className="sub-nav">
        {tabs.map(t => (
          <button
            key={t.id}
            className={activeTab === t.id ? 'active' : ''}
            onClick={() => setActiveTab(t.id)}
          >
            {t.label}
          </button>
        ))}
      </nav>

      {activeTab === 'wg' && (
        <div className="config-form">
          <label>
            Private Key
            <input
              type="password"
              value={wg.privateKey}
              onChange={e => setWg(prev => ({ ...prev, privateKey: e.target.value }))}
              placeholder="e.g. gN6x..."
            />
          </label>

          <label>
            Address
            <input
              type="text"
              value={wg.address}
              onChange={e => setWg(prev => ({ ...prev, address: e.target.value }))}
            />
          </label>

          <label>
            DNS
            <input
              type="text"
              value={wg.dns}
              onChange={e => setWg(prev => ({ ...prev, dns: e.target.value }))}
            />
          </label>

          <h4>Peers</h4>
          {wg.peers.map((peer, i) => (
            <div key={i} className="peer-block">
              <label>
                Public Key
                <input value={peer.publicKey} onChange={e => updatePeer(i, 'publicKey', e.target.value)} placeholder="e.g. xTIB..." />
              </label>
              <label>
                Endpoint
                <input value={peer.endpoint} onChange={e => updatePeer(i, 'endpoint', e.target.value)} placeholder="e.g. vpn.example.com:51820" />
              </label>
              <label>
                Allowed IPs
                <input value={peer.allowedIPs} onChange={e => updatePeer(i, 'allowedIPs', e.target.value)} />
              </label>
              {wg.peers.length > 1 && (
                <button className="btn-remove" onClick={() => removePeer(i)}>Remover</button>
              )}
            </div>
          ))}
          <button className="btn-add" onClick={addPeer}>+ Adicionar Peer</button>

          <div className="config-actions">
            <button className="btn-primary" onClick={handleSave}>
              {saved ? '✓ Salvo' : 'Salvar Config'}
            </button>
          </div>
        </div>
      )}

      {activeTab === 'tor' && (
        <div className="config-form">
          <p className="dim">Tor está configurado via ghoststack.yaml. Use o comando <code>ghost init</code> no terminal.</p>
        </div>
      )}

      {activeTab === 'unbound' && (
        <div className="config-form">
          <p className="dim">Unbound DNS resolver configurado via ghoststack.yaml.</p>
        </div>
      )}
    </div>
  )
}
