'use client'

import { useState } from 'react'
import { Shield, Bell, Save, TestTube } from 'lucide-react'
import Link from 'next/link'
import { configureAlerts } from '@/lib/api'

export default function SettingsPage() {
  const [telegramToken, setTelegramToken] = useState('')
  const [telegramChatId, setTelegramChatId] = useState('')
  const [discordWebhook, setDiscordWebhook] = useState('')
  const [slackWebhook, setSlackWebhook] = useState('')
  const [email, setEmail] = useState('')
  const [saving, setSaving] = useState(false)
  const [message, setMessage] = useState('')

  const handleSave = async () => {
    setSaving(true)
    setMessage('')
    try {
      await configureAlerts({
        telegram_token: telegramToken || undefined,
        telegram_chat_id: telegramChatId || undefined,
        discord_webhook: discordWebhook || undefined,
        slack_webhook: slackWebhook || undefined,
        email: email || undefined
      })
      setMessage('Settings saved successfully!')
      setTimeout(() => setMessage(''), 3000)
    } catch (error) {
      setMessage('Failed to save settings')
    } finally {
      setSaving(false)
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
          <Link href="/firewall" className="text-gray-400 hover:text-white">Firewall</Link>
          <Link href="/settings" className="text-[var(--primary)] font-medium">Settings</Link>
        </div>
      </nav>

      <main className="p-6 max-w-4xl">
        <div className="bg-[var(--card-bg)] border border-[var(--border)] rounded-lg p-6">
          <h2 className="text-xl font-semibold mb-6 flex items-center gap-2">
            <Bell className="w-5 h-5 text-[var(--primary)]" />
            Alert Notifications
          </h2>

          {message && (
            <div className={`mb-4 p-3 rounded ${message.includes('success') ? 'bg-[var(--primary)] bg-opacity-20 text-[var(--primary)]' : 'bg-[var(--danger)] bg-opacity-20 text-[var(--danger)]'}`}>
              {message}
            </div>
          )}

          <div className="space-y-6">
            <div>
              <h3 className="text-lg font-medium mb-4">Telegram</h3>
              <div className="space-y-3">
                <div>
                  <label className="block text-sm text-gray-400 mb-2">Bot Token</label>
                  <input
                    type="text"
                    value={telegramToken}
                    onChange={(e) => setTelegramToken(e.target.value)}
                    placeholder="1234567890:ABCdefGHIjklMNOpqrsTUVwxyz"
                    className="w-full bg-[var(--background)] border border-[var(--border)] rounded p-3 text-white"
                  />
                </div>
                <div>
                  <label className="block text-sm text-gray-400 mb-2">Chat ID</label>
                  <input
                    type="text"
                    value={telegramChatId}
                    onChange={(e) => setTelegramChatId(e.target.value)}
                    placeholder="-1001234567890"
                    className="w-full bg-[var(--background)] border border-[var(--border)] rounded p-3 text-white"
                  />
                </div>
              </div>
            </div>

            <div className="border-t border-[var(--border)] pt-6">
              <h3 className="text-lg font-medium mb-4">Discord</h3>
              <div>
                <label className="block text-sm text-gray-400 mb-2">Webhook URL</label>
                <input
                  type="text"
                  value={discordWebhook}
                  onChange={(e) => setDiscordWebhook(e.target.value)}
                  placeholder="https://discord.com/api/webhooks/..."
                  className="w-full bg-[var(--background)] border border-[var(--border)] rounded p-3 text-white"
                />
              </div>
            </div>

            <div className="border-t border-[var(--border)] pt-6">
              <h3 className="text-lg font-medium mb-4">Slack</h3>
              <div>
                <label className="block text-sm text-gray-400 mb-2">Webhook URL</label>
                <input
                  type="text"
                  value={slackWebhook}
                  onChange={(e) => setSlackWebhook(e.target.value)}
                  placeholder="https://hooks.slack.com/services/..."
                  className="w-full bg-[var(--background)] border border-[var(--border)] rounded p-3 text-white"
                />
              </div>
            </div>

            <div className="border-t border-[var(--border)] pt-6">
              <h3 className="text-lg font-medium mb-4">Email</h3>
              <div>
                <label className="block text-sm text-gray-400 mb-2">Email Address</label>
                <input
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  placeholder="admin@example.com"
                  className="w-full bg-[var(--background)] border border-[var(--border)] rounded p-3 text-white"
                />
              </div>
            </div>

            <div className="flex gap-3 pt-6">
              <button
                onClick={handleSave}
                disabled={saving}
                className="flex items-center gap-2 px-6 py-3 bg-[var(--primary)] text-black rounded font-medium hover:opacity-90 disabled:opacity-50"
              >
                <Save className="w-4 h-4" />
                {saving ? 'Saving...' : 'Save Settings'}
              </button>
              <button
                className="flex items-center gap-2 px-6 py-3 bg-[var(--background)] border border-[var(--border)] rounded hover:bg-[var(--border)]"
              >
                <TestTube className="w-4 h-4" />
                Test Alert
              </button>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
