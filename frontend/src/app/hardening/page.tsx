'use client'

import { useEffect, useState } from 'react'
import { Shield, Cpu, Users, FolderCheck, Server, Search, ChevronDown, ChevronUp, CheckCircle, XCircle } from 'lucide-react'
import Link from 'next/link'
import {
  getHardeningStatus, applyKernelHardening, applyAccountHardening,
  initAIDE, checkAIDE, minimizeServices, listServices,
  runAudit, getAuditStatus
} from '@/lib/api'

export default function HardeningPage() {
  const [hardeningStatus, setHardeningStatus] = useState<any>(null)
  const [auditStatus, setAuditStatus] = useState<any>(null)
  const [services, setServices] = useState<any>(null)
  const [loading, setLoading] = useState(true)
  const [expandedSection, setExpandedSection] = useState<string | null>(null)

  useEffect(() => {
    fetchData()
  }, [])

  const fetchData = async () => {
    try {
      const [status, audit, svc] = await Promise.all([
        getHardeningStatus(),
        getAuditStatus(),
        listServices(),
      ])
      setHardeningStatus(status)
      setAuditStatus(audit)
      setServices(svc)
      setLoading(false)
    } catch (error) {
      console.error('Failed to fetch hardening data:', error)
      setLoading(false)
    }
  }

  const handleAction = async (action: () => Promise<any>, name: string) => {
    try {
      await action()
      fetchData()
    } catch (error) {
      console.error(`Failed to ${name}:`, error)
    }
  }

  if (loading) {
    return (
      <div className="min-h-screen bg-[var(--background)] flex items-center justify-center">
        <div className="text-center">
          <Shield className="w-16 h-16 text-[var(--primary)] mx-auto mb-4 animate-pulse" />
          <p className="text-gray-400">Loading hardening data...</p>
        </div>
      </div>
    )
  }

  const modules = [
    {
      id: 'kernel',
      icon: <Cpu className="w-6 h-6" />,
      title: 'Kernel Hardening',
      description: 'sysctl parameters: ASLR, anti-spoofing, SYN flood protection',
      status: hardeningStatus?.kernel_hardening?.status || 'inactive',
      actions: [{ label: 'Apply Kernel Hardening', action: () => handleAction(applyKernelHardening, 'kernel') }],
    },
    {
      id: 'accounts',
      icon: <Users className="w-6 h-6" />,
      title: 'Account Hardening',
      description: 'PAM policy, faillock lockout, sudo audit, password aging',
      status: hardeningStatus?.account_hardening?.status || 'inactive',
      actions: [{ label: 'Apply Account Hardening', action: () => handleAction(applyAccountHardening, 'accounts') }],
    },
    {
      id: 'aide',
      icon: <FolderCheck className="w-6 h-6" />,
      title: 'AIDE Integrity',
      description: 'File integrity monitoring with daily automated checks',
      status: hardeningStatus?.file_integrity?.status || 'inactive',
      actions: [
        { label: 'Initialize Database', action: () => handleAction(initAIDE, 'aide-init') },
        { label: 'Run Integrity Check', action: () => handleAction(checkAIDE, 'aide-check') },
      ],
    },
    {
      id: 'services',
      icon: <Server className="w-6 h-6" />,
      title: 'Service Minimization',
      description: 'Disable dangerous services: telnet, ftp, nfs, cups, etc.',
      status: hardeningStatus?.service_minimization?.status || 'inactive',
      actions: [{ label: 'Minimize Services', action: () => handleAction(minimizeServices, 'services') }],
    },
    {
      id: 'audit',
      icon: <Search className="w-6 h-6" />,
      title: 'Security Audit',
      description: 'Lynis security audit with hardening score',
      status: hardeningStatus?.security_audit?.status || 'inactive',
      actions: [
        { label: 'Install Lynis', action: () => handleAction(import('@/lib/api').then(m => m.installAudit), 'install-audit') },
        { label: 'Run Audit', action: () => handleAction(runAudit, 'audit') },
      ],
    },
  ]

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
          <Link href="/firewall" className="text-gray-400 hover:text-white">Firewall</Link>
          <Link href="/hardening" className="text-[var(--primary)] font-medium">Hardening</Link>
          <Link href="/settings" className="text-gray-400 hover:text-white">Settings</Link>
        </div>
      </nav>

      <main className="p-6">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-xl font-semibold flex items-center gap-2">
            <Shield className="w-5 h-5 text-[var(--primary)]" />
            System Hardening
          </h2>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {modules.map((mod) => (
            <div key={mod.id} className="bg-[var(--card-bg)] border border-[var(--border)] rounded-lg p-6">
              <div className="flex items-start justify-between mb-3">
                <div className="flex items-center gap-3">
                  <div className={`p-2 rounded-lg ${mod.status === 'active' ? 'bg-green-500/10 text-[var(--primary)]' : 'bg-gray-500/10 text-gray-400'}`}>
                    {mod.icon}
                  </div>
                  <div>
                    <h3 className="font-semibold">{mod.title}</h3>
                    <p className="text-sm text-gray-400">{mod.description}</p>
                  </div>
                </div>
                <div className={`flex items-center gap-1 text-sm ${mod.status === 'active' ? 'text-[var(--primary)]' : 'text-gray-500'}`}>
                  {mod.status === 'active' ? <CheckCircle className="w-4 h-4" /> : <XCircle className="w-4 h-4" />}
                  {mod.status === 'active' ? 'Active' : 'Inactive'}
                </div>
              </div>

              <div className="flex flex-wrap gap-2 mt-4">
                {mod.actions.map((action, i) => (
                  <button
                    key={i}
                    onClick={action.action}
                    className="px-3 py-1.5 text-sm bg-[var(--background)] border border-[var(--border)] rounded hover:bg-[var(--border)] transition-colors"
                  >
                    {action.label}
                  </button>
                ))}
              </div>
            </div>
          ))}
        </div>

        <div className="mt-6 bg-[var(--card-bg)] border border-[var(--border)] rounded-lg p-6">
          <div
            className="flex items-center justify-between cursor-pointer"
            onClick={() => setExpandedSection(expandedSection === 'services' ? null : 'services')}
          >
            <h3 className="font-semibold flex items-center gap-2">
              <Server className="w-5 h-5 text-[var(--primary)]" />
              Service Classification
            </h3>
            {expandedSection === 'services' ? <ChevronUp className="w-5 h-5" /> : <ChevronDown className="w-5 h-5" />}
          </div>

          {expandedSection === 'services' && services && (
            <div className="mt-4 space-y-4">
              {['dangerous', 'essential', 'optional'].map((category) => (
                <div key={category}>
                  <h4 className={`text-sm font-medium mb-2 capitalize ${
                    category === 'dangerous' ? 'text-red-400' : category === 'essential' ? 'text-green-400' : 'text-yellow-400'
                  }`}>
                    {category} Services
                  </h4>
                  <div className="space-y-1">
                    {services[category]?.map((svc: any, i: number) => (
                      <div key={i} className="flex justify-between text-sm py-1 px-2 rounded hover:bg-[var(--background)]">
                        <span>{svc.name}</span>
                        <span className={svc.status === 'disabled' ? 'text-gray-500' : 'text-[var(--primary)]'}>
                          {svc.status}
                        </span>
                      </div>
                    ))}
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </main>
    </div>
  )
}
