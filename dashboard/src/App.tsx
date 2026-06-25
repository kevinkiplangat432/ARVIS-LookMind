import { useEffect, useState, useCallback } from 'react'
import { fetchRequests, fetchAnomalies } from './api/client'
import type { Request, Anomaly } from './types'
import StatsBar from './components/StatsBar'
import RequestsTable from './components/RequestsTable'
import AnomaliesTable from './components/AnomaliesTable'

const POLL_INTERVAL = 30000

export default function App() {
  const [requests, setRequests] = useState<Request[]>([])
  const [anomalies, setAnomalies] = useState<Anomaly[]>([])
  const [error, setError] = useState<string | null>(null)
  const [lastUpdated, setLastUpdated] = useState<Date | null>(null)

  const refresh = useCallback(async () => {
    try {
      const [reqs, anoms] = await Promise.all([
        fetchRequests(),
        fetchAnomalies(),
      ])
      setRequests(reqs)
      setAnomalies(anoms)
      setLastUpdated(new Date())
      setError(null)
    } catch (err) {
      setError('Failed to reach ARVIS API. Is the server running?')
      console.error(err)
    }
  }, [])

  useEffect(() => {
    refresh()
    const interval = setInterval(refresh, POLL_INTERVAL)
    return () => clearInterval(interval)
  }, [refresh])

  return (
    <div className="min-h-screen bg-gray-950 text-white">
      <header className="border-b border-gray-800 px-8 py-4 flex items-center justify-between">
        <div>
          <h1 className="text-xl font-bold tracking-wide">ARVIS</h1>
          <p className="text-gray-500 text-xs mt-0.5">
            Automated Runtime Visibility & Intelligence System
          </p>
        </div>
        <div className="text-right">
          <button
            onClick={refresh}
            className="text-sm text-blue-400 hover:text-blue-300 transition-colors"
          >
            Refresh
          </button>
          {lastUpdated && (
            <p className="text-gray-600 text-xs mt-1">
              Updated {lastUpdated.toLocaleTimeString()}
            </p>
          )}
        </div>
      </header>

      <main className="px-8 py-6 max-w-7xl mx-auto">
        {error && (
          <div className="mb-6 bg-red-950 border border-red-800 text-red-400 px-4 py-3 rounded-lg text-sm">
            {error}
          </div>
        )}
        <StatsBar requests={requests} anomalies={anomalies} />
        <RequestsTable requests={requests} />
        <AnomaliesTable anomalies={anomalies} />
      </main>
    </div>
  )
}