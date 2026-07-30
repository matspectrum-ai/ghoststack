interface Peer {
  id: string
  endpoint: string
  allowedIPs: string[]
  handshake: string
  transferRx: string
  transferTx: string
}

const MOCK_PEERS: Peer[] = [
  { id: 'peer-1', endpoint: '203.0.113.10:51820', allowedIPs: ['10.0.0.2/32'], handshake: '2s ago', transferRx: '1.2 GB', transferTx: '340 MB' },
  { id: 'peer-2', endpoint: '198.51.100.20:51820', allowedIPs: ['10.0.0.3/32'], handshake: '15s ago', transferRx: '450 MB', transferTx: '890 MB' },
]

export default function PeerList() {
  const peers = MOCK_PEERS

  return (
    <div className="card">
      <h3>Pares ({peers.length})</h3>
      <table className="peer-table">
        <thead>
          <tr>
            <th>Endpoint</th>
            <th>Allowed IPs</th>
            <th>Handshake</th>
            <th>RX</th>
            <th>TX</th>
          </tr>
        </thead>
        <tbody>
          {peers.map(p => (
            <tr key={p.id}>
              <td className="mono">{p.endpoint}</td>
              <td className="mono">{p.allowedIPs.join(', ')}</td>
              <td>{p.handshake}</td>
              <td>{p.transferRx}</td>
              <td>{p.transferTx}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
