'use client'

import { useEffect, useState } from 'react'
import { Shield, Lock, Ban, Plus, Trash2, RefreshCw } from 'lucide-react'
import Link from 'next/link'
import { getFirewallStatus, getBlockedIPs, blockIP, unblockIP } from '@/lib/api'

export default function FirewallPage() {
  const [status, setStatus] = useState<any>(null)
  const [blockedIPs, setBlockedIPs] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [newIP, setNewIP] = useState('')
  const [reason, setReason] = useState('')

  useEffect(() => {
    fetchData()
  }, [])

  const fetchData = async () => {
    try {
      const [statusData, blockedData] = await Promise.all([
        getFirewallStatus(),
        getBlockedIPs()
      ])
      setStatus(statusData)
      setBlockedIPs(blockedData.blocked_ips || [])
      setLoading(false)
    } catch (error) {
      console.error('Failed to fetch firewall data:', error)
      setLoading(false)
    }
  }

  const handleBlockIP = async () => {
    if (!newIP) return
    try {
      await blockIP(newIP, reason)
      setNewIP('')
      setReason('')
      fetchData()
    } catch (error) {
      console.error('Failed to block IP:', error)
    }
  }

  const handleUnblockIP = async (ip: string) => {
    try {
      await unblockIP(ip)
      fetchData()
    } catch (error) {
      console.error('Failed to unblock IP:', error)
    }
  }

  return (
    <div className="min-h-screen bg-[var(--background)]">
      <header className="bg-[var(--card-bg)] border-b border-[var(--border)] px-6 py-4">
        <div className="flex items-center gap-3">
          <Shield className="w-8 h-8 text-[var(--primary)]" />
          <h1 className="text-2xl font-bold">IronWall</h1>
        </div>
      </header>

      <nav className="bg-[var(--card-bg)] border-b border-[var(--border)] px-6 py-3">
        <div className="flex gap-6">
          <Link href="/" className="text-gray-400 hover:text-white">Dashboard</Link>
          <Link href="/events" className="text-gray-400 hover:text-white">Events</Link>
          <Link href="/modules" className="text-gray-400 hover:text-white">Modules</Link>
          <Link href="/firewall" className="text-[var(--primary)] font-medium">Firewall</Link>
          <Link href="/settings" className="text-gray-400 hover:text-white">Settings</Link>
        </div>
      </nav>

      <main className="p-6">
        {loading ? (
          <div className="text-center py-8 text-gray-400">Loading firewall data...</div>
        ) : (
          <>
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
              <div className="bg-[var(--card-bg)] border border-[var(--border)] rounded-lg p-6">
                <h2 className="text-xl font-semibold mb-4 flex items-center gap-2">
                  <Lock className="w-5 h-5 text-[var(--primary)]" />
                  Firewall Status
                </h2>
                <div className="space-y-3">
                  <div className="flex justify-between">
                    <span className="text-gray-400">Status</span>
                    <span className="text-[var(--primary)] font-medium">{status?.status}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-400">Backend</span>
                    <span>{status?.backend}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-400">Rate Limiting</span>
                    <span className={status?.rate_limiting ? 'text-[var(--primary)]' : 'text-gray-500'}>
                      {status?.rate_limiting ? 'Enabled' : 'Disabled'}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-400">Anti Port Scan</span>
                    <span className={status?.anti_port_scan ? 'text-[var(--primary)]' : 'text-gray-500'}>
                      {status?.anti_port_scan ? 'Enabled' : 'Disabled'}
                    </span>
                  </div>
                </div>

                <h3 className="text-lg font-semibold mt-6 mb-3">Allowed Ports</h3>
                <div className="space-y-2">
                  {status?.allowed_ports?.map((port: any) => (
                    <div key={port.port} className="flex justify-between items-center bg-[var(--background)] rounded p-2">
                      <span>{port.port}/{port.protocol}</span>
                      <span className="text-gray-400">{port.service}</span>
                    </div>
                  ))}
                </div>
              </div>

              <div className="bg-[var(--card-bg)] border border-[var(--border)] rounded-lg p-6">
                <h2 className="text-xl font-semibold mb-4 flex items-center gap-2">
                  <Ban className="w-5 h-5 text-[var(--danger)]" />
                  Block IP Address
                </h2>
                <div className="space-y-4">
                  <div>
                    <label className="block text-sm text-gray-400 mb-2">IP Address</label>
                    <input
                      type="text"
                      value={newIP}
                      onChange={(e) => setNewIP(e.target.value)}
                      placeholder="e.g., 192.168.1.100"
                      className="w-full bg-[var(--background)] border border-[var(--border)] rounded p-3 text-white"
                    />
                  </div>
                  <div>
                    <label className="block text-sm text-gray-400 mb-2">Reason (optional)</label>
                    <input
                      type="text"
                      value={reason}
                      onChange={(e) => setReason(e.target.value)}
                      placeholder="e.g., Brute force attempt"
                      className="w-full bg-[var(--background)] border border-[var(--border)] rounded p-3 text-white"
                    />
                  </div>
                  <button
                    onClick={handleBlockIP}
                    disabled={!newIP}
                    className="w-full bg-[var(--danger)] text-white py-3 rounded font-medium hover:opacity-90 disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    Block IP
                  </button>
                </div>
              </div>
            </div>

            <div className="bg-[var(--card-bg)] border border-[var(--border)] rounded-lg p-6">
              <div className="flex items-center justify-between mb-4">
                <h2 className="text-xl font-semibold">Blocked IPs ({blockedIPs.length})</h2>
                <button
                  onClick={fetchData}
                  className="flex items-center gap-2 px-4 py-2 bg-[var(--background)] border border-[var(--border)] rounded hover:bg-[var(--border)]"
                >
                  <RefreshCw className="w-4 h-4" />
                  Refresh
                </button>
              </div>

              {blockedIPs.length === 0 ? (
                <div className="text-center py-8 text-gray-400">No blocked IPs</div>
              ) : (
                <div className="overflow-x-auto">
                  <table className="w-full">
                    <thead>
                      <tr className="border-b border-[var(--border)]">
                        <th className="text-left py-3 px-4 text-gray-400 font-medium">IP Address</th>
                        <th className="text-left py-3 px-4 text-gray-400 font-medium">Reason</th>
                        <th className="text-left py-3 px-4 text-gray-400 font-medium">Blocked At</th>
                        <th className="text-right py-3 px-4 text-gray-400 font-medium">Action</th>
                      </tr>
                    </thead>
                    <tbody>
                      {blockedIPs.map((item: any) => (
                        <tr key={item.ip} className="border-b border-[var(--border)]">
                          <td className="py-3 px-4 font-mono">{item.ip}</td>
                          <td className="py-3 px-4 text-gray-400">{item.reason}</td>
                          <td className="py-3 px-4 text-gray-400">{new Date(item.blocked_at).toLocaleString()}</td>
                          <td className="py-3 px-4 text-right">
                            <button
                              onClick={() => handleUnblockIP(item.ip)}
                              className="text-[var(--primary)] hover:underline"
                            >
                              Unblock
                            </button>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              )}
            </div>
          </>
        )}
      </main>
    </div>
  )
}
