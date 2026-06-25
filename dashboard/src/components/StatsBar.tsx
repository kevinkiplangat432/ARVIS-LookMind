import type { Request, Anomaly } from '../types'

interface Props {
  requests: Request[]
  anomalies: Anomaly[]
}

export default function StatsBar({ requests, anomalies }: Props) {
  const avgLatency = requests.length
    ? Math.round(requests.reduce((sum, r) => sum + r.latency_ms, 0) / requests.length)
    : 0

  const stats = [
    { label: 'Total Requests', value: requests.length },
    { label: 'Total Anomalies', value: anomalies.length },
    { label: 'Avg Latency', value: `${avgLatency}ms` },
  ]

  return (
    <div className="grid grid-cols-3 gap-4 mb-8">
      {stats.map((s) => (
        <div key={s.label} className="bg-gray-900 border border-gray-800 rounded-lg p-4">
          <p className="text-gray-400 text-sm">{s.label}</p>
          <p className="text-white text-2xl font-bold mt-1">{s.value}</p>
        </div>
      ))}
    </div>
  )
}