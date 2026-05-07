'use client'

import { useEffect, useState } from 'react'
import { Shield, AlertTriangle, Clock, Filter } from 'lucide-react'
import Link from 'next/link'
import { getEvents } from '@/lib/api'

export default function EventsPage() {
  const [events, setEvents] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [filter, setFilter] = useState('all')

  useEffect(() => {
    fetchEvents()
  }, [filter])

  const fetchEvents = async () => {
    try {
      const params = filter !== 'all' ? { severity: filter } : {}
      const data = await getEvents(params)
      setEvents(data)
      setLoading(false)
    } catch (error) {
      console.error('Failed to fetch events:', error)
      setLoading(false)
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
          <Link href="/events" className="text-[var(--primary)] font-medium">Events</Link>
          <Link href="/modules" className="text-gray-400 hover:text-white">Modules</Link>
          <Link href="/firewall" className="text-gray-400 hover:text-white">Firewall</Link>
          <Link href="/settings" className="text-gray-400 hover:text-white">Settings</Link>
        </div>
      </nav>

      <main className="p-6">
        <div className="bg-[var(--card-bg)] border border-[var(--border)] rounded-lg p-6">
          <div className="flex items-center justify-between mb-6">
            <h2 className="text-xl font-semibold flex items-center gap-2">
              <AlertTriangle className="w-5 h-5 text-[var(--primary)]" />
              Security Events
            </h2>
            <div className="flex gap-2">
              <button
                onClick={() => setFilter('all')}
                className={`px-4 py-2 rounded ${filter === 'all' ? 'bg-[var(--primary)] text-black' : 'bg-[var(--border)] text-gray-400'}`}
              >
                All
              </button>
              <button
                onClick={() => setFilter('critical')}
                className={`px-4 py-2 rounded ${filter === 'critical' ? 'bg-[var(--danger)] text-white' : 'bg-[var(--border)] text-gray-400'}`}
              >
                Critical
              </button>
              <button
                onClick={() => setFilter('warning')}
                className={`px-4 py-2 rounded ${filter === 'warning' ? 'bg-[var(--warning)] text-black' : 'bg-[var(--border)] text-gray-400'}`}
              >
                Warning
              </button>
            </div>
          </div>

          {loading ? (
            <div className="text-center py-8 text-gray-400">Loading events...</div>
          ) : events.length === 0 ? (
            <div className="text-center py-8 text-gray-400">No security events found</div>
          ) : (
            <div className="space-y-3">
              {events.map((event) => (
                <EventItem key={event.id} event={event} />
              ))}
            </div>
          )}
        </div>
      </main>
    </div>
  )
}

function EventItem({ event }: any) {
  const severityColors = {
    critical: 'text-[var(--danger)]',
    warning: 'text-[var(--warning)]',
    info: 'text-[var(--primary)]'
  }

  return (
    <div className="bg-[var(--background)] border border-[var(--border)] rounded p-4">
      <div className="flex items-start justify-between">
        <div className="flex-1">
          <div className="flex items-center gap-2 mb-1">
            <span className={`font-semibold ${severityColors[event.severity as keyof typeof severityColors]}`}>
              {event.severity.toUpperCase()}
            </span>
            <span className="text-gray-500">•</span>
            <span className="text-gray-400">{event.event_type}</span>
          </div>
          <p className="text-white mb-2">{event.description}</p>
          <div className="flex items-center gap-2 text-sm text-gray-500">
            <Clock className="w-4 h-4" />
            <span>{new Date(event.created_at).toLocaleString()}</span>
          </div>
        </div>
      </div>
    </div>
  )
}
