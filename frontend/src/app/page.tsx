'use client'

import { useEffect, useState } from 'react'
import { Shield, Activity, AlertTriangle, Lock, Server, FileText } from 'lucide-react'
import Link from 'next/link'

export default function DashboardPage() {
  const [status, setStatus] = useState<any>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchStatus()
    const interval = setInterval(fetchStatus, 5000)
    return () => clearInterval(interval)
  }, [])

  const fetchStatus = async () => {
    try {
      const response = await fetch('http://localhost:8001/api/v1/status')
      const data = await response.json()
      setStatus(data)
      setLoading(false)
    } catch (error) {
      console.error('Failed to fetch status:', error)
      setLoading(false)
    }
  }

  if (loading) {
    return (
      <div className="min-h-screen bg-[var(--background)] flex items-center justify-center">
        <div className="text-center">
          <Shield className="w-16 h-16 text-[var(--primary)] mx-auto mb-4 animate-pulse" />
          <p className="text-gray-400">Loading dashboard...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-[var(--background)]">
      {/* Header */}
      <header className="bg-[var(--card-bg)] border-b border-[var(--border)] px-6 py-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <Shield className="w-8 h-8 text-[var(--primary)]" />
            <h1 className="text-2xl font-bold">IronWall</h1>
          </div>
          <div className="flex items-center gap-2">
            <div className="w-3 h-3 bg-[var(--primary)] rounded-full animate-pulse"></div>
            <span className="text-sm text-gray-400">Protected</span>
          </div>
        </div>
      </header>

      {/* Navigation */}
      <nav className="bg-[var(--card-bg)] border-b border-[var(--border)] px-6 py-3">
        <div className="flex gap-6">
          <Link href="/" className="text-[var(--primary)] font-medium">Dashboard</Link>
          <Link href="/events" className="text-gray-400 hover:text-white">Events</Link>
          <Link href="/modules" className="text-gray-400 hover:text-white">Modules</Link>
          <Link href="/firewall" className="text-gray-400 hover:text-white">Firewall</Link>
          <Link href="/settings" className="text-gray-400 hover:text-white">Settings</Link>
        </div>
      </nav>

      {/* Main Content */}
      <main className="p-6">
        {/* Stats Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-6">
          <StatCard
            icon={<Activity className="w-6 h-6" />}
            title="System Status"
            value={status?.status || 'Unknown'}
            color="primary"
          />
          <StatCard
            icon={<AlertTriangle className="w-6 h-6" />}
            title="Events (24h)"
            value={status?.events_24h || 0}
            color="warning"
          />
          <StatCard
            icon={<Lock className="w-6 h-6" />}
            title="Blocked IPs"
            value={status?.blocked_ips_count || 0}
            color="danger"
          />
          <StatCard
            icon={<Server className="w-6 h-6" />}
            title="CPU Usage"
            value={`${status?.cpu_usage || 0}%`}
            color="primary"
          />
        </div>

        {/* System Health */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
          <div className="bg-[var(--card-bg)] border border-[var(--border)] rounded-lg p-6">
            <h2 className="text-xl font-semibold mb-4 flex items-center gap-2">
              <Activity className="w-5 h-5 text-[var(--primary)]" />
              System Health
            </h2>
            <div className="space-y-3">
              <HealthItem label="Uptime" value={status?.uptime || 'N/A'} />
              <HealthItem label="CPU Usage" value={`${status?.cpu_usage || 0}%`} />
              <HealthItem label="Memory Usage" value={`${status?.memory_usage || 0}%`} />
            </div>
          </div>

          <div className="bg-[var(--card-bg)] border border-[var(--border)] rounded-lg p-6">
            <h2 className="text-xl font-semibold mb-4 flex items-center gap-2">
              <Shield className="w-5 h-5 text-[var(--primary)]" />
              Security Modules
            </h2>
            <div className="space-y-2">
              {status?.modules?.map((module: any) => (
                <ModuleStatus key={module.id} name={module.name} enabled={module.enabled} />
              ))}
            </div>
          </div>
        </div>

        {/* Recent Events */}
        <div className="bg-[var(--card-bg)] border border-[var(--border)] rounded-lg p-6">
          <h2 className="text-xl font-semibold mb-4 flex items-center gap-2">
            <FileText className="w-5 h-5 text-[var(--primary)]" />
            Recent Events
          </h2>
          <div className="text-gray-400 text-center py-8">
            No recent security events
          </div>
        </div>
      </main>
    </div>
  )
}

function StatCard({ icon, title, value, color }: any) {
  const colorClasses = {
    primary: 'text-[var(--primary)]',
    warning: 'text-[var(--warning)]',
    danger: 'text-[var(--danger)]'
  }

  return (
    <div className="bg-[var(--card-bg)] border border-[var(--border)] rounded-lg p-6">
      <div className={`${colorClasses[color]} mb-2`}>{icon}</div>
      <h3 className="text-sm text-gray-400 mb-1">{title}</h3>
      <p className="text-2xl font-bold">{value}</p>
    </div>
  )
}

function HealthItem({ label, value }: any) {
  return (
    <div className="flex justify-between items-center">
      <span className="text-gray-400">{label}</span>
      <span className="font-medium">{value}</span>
    </div>
  )
}

function ModuleStatus({ name, enabled }: any) {
  return (
    <div className="flex justify-between items-center py-2">
      <span className="capitalize">{name.replace(/_/g, ' ')}</span>
      <span className={`text-sm ${enabled ? 'text-[var(--primary)]' : 'text-gray-500'}`}>
        {enabled ? '✓ Active' : '○ Inactive'}
      </span>
    </div>
  )
}
