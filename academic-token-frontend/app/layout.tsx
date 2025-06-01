import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import './globals.css'
import { BlockchainProvider } from './context/BlockchainContext'

const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
  title: 'AcademicToken - Academic Token System',
  description: 'Blockchain platform for decentralized academic certification',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className={inter.className} suppressHydrationWarning>
        <BlockchainProvider>
          {children}
        </BlockchainProvider>
      </body>
    </html>
  )
}