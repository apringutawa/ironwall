import type { Metadata } from 'next'
import './globals.css'

export const metadata: Metadata = {
  title: 'IronWall Dashboard',
  description: 'Linux Server Hardening & Threat Prevention Toolkit',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  )
}
