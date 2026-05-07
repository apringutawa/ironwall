'use client'

import { useEffect, useState } from 'react'
import { Shield, Power, RefreshCw } from 'lucide-react'
import Link from 'next/link'
import { getModules, enableModule, disableModule } from '@/lib/api'

export default function ModulesPage() {
  const [modules, setModules] = useState<any[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchModules()
  }, [])

  const fetchModules = async () => {
    try {
      const data = await getModules()
      setModules(data)
      setLoading(false)
    } catch (error) {
      console.error('Failed to fetch modules:', error)
      setLoading(false)
    }
  }

  const toggleModule = async (name: string, enabled: boolean) => {
    try {
      if (enabled) {
        await disableModule(name)
      } else {
        await enableModule(name)
      }
      fetchModules()
    } catch (error) {
      console.error('Failed to toggle module:', error)
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
          <Link href="/modules" className="text-[var(--primary)] font-medium">Modules</Link>
          <Link href="/firewall" className="text-gray-400 hover:text-white">Firewall</Link>
          <Link href="/settings" className="text-gray-400 hover:text-white">Settings</Link>
        </div>
      </nav>

      <main className="p-6">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-xl font-semibold flex items-center gap-2">
            <Power className="w-5 h-5 text-[var(--primary)]" />
            Security Modules
          </h2>
          <button
            onClick={fetchModules}
            className="flex items-center gap-2 px-4 py-2 bg-[var(--card-bg)] border border-[var(--border)] rounded hover:bg-[var(--border)]"
          >
            <RefreshCw className="w-4 h-4" />
            Refresh
          </button>
        </div>

        {loading ? (
          <div className="text-center py-8 text-gray-400">Loading modules...</div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {modules.map((module) => (
              <ModuleCard key={module.id} module={module} onToggle={toggleModule} />
            ))}
          </div>
        )}
      </main>
    </div>
  )
}

function ModuleCard({ module, onToggle }: any) {
  const moduleInfo: Record<string, { icon: string; description: string }> = {
    ssh_hardening: {
      icon: '🔐',
      description: 'SSH security configuration and protection'
    },
    firewall: {
      icon: '🛡️',
      description: 'Network firewall rules and filtering'
    },
    fail2ban: {
      icon: '🚫',
      description: 'Brute force attack prevention'
    },
    malware_scanner: {
      icon: '🦠',
      description: 'Malware and virus detection'
    },
    cron_protection: {
      icon: '⏰',
      description: 'Cron job monitoring and protection'
    },
    file_integrity: {
      icon: '📁',
      description: 'Critical file monitoring'
    }
  }

  const info = moduleInfo[module.name] || { icon: '📦', description: 'Security module' }

  return (
    <div className="bg-[var(--card-bg)] border border-[var(--border)] rounded-lg p-6">
      <div className="flex items-start justify-between mb-4">
        <div className="text-3xl">{info.icon}</div>
        <button
          onClick={() => onToggle(module.name, module.enabled)}
          className={`relative w-12 h-6 rounded-full transition-colors ${
            module.enabled ? 'bg-[var(--primary)]' : 'bg-[var(--border)]'
          }`}
        >
          <div
            className={`absolute top-1 w-4 h-4 rounded-full bg-white transition-transform ${
              module.enabled ? 'right-1' : 'left-1'
            }`}
          />
        </button>
      </div>
      <h3 className="text-lg font-semibold mb-2 capitalize">
        {module.name.replace(/_/g, ' ')}
      </h3>
      <p className="text-sm text-gray-400 mb-4">{info.description}</p>
      <div className="flex items-center gap-2 text-sm">
        <div className={`w-2 h-2 rounded-full ${module.enabled ? 'bg-[var(--primary)]' : 'bg-gray-500'}`} />
        <span className={module.enabled ? 'text-[var(--primary)]' : 'text-gray-500'}>
          {module.enabled ? 'Active' : 'Inactive'}
        </span>
      </div>
    </div>
  )
}
